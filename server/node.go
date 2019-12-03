package server

import (
	"encoding/base64"
	"fmt"
	"github.com/niclabs/dtcnode/v2/config"
	"github.com/niclabs/dtcnode/v2/message"
	"github.com/niclabs/tcrsa"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"log"
	"net"
	"sync"
)

// The domain of the ZMQ connection. This value must be the same in the server, or it will not work.
const TchsmDomain = "tchsm"

// The protocol used for the ZMQ connection. TCP is the best for this usage cases.
const TchsmProtocol = "tcp"

// Node represents a node in the distributed TCHSM application. It saves zero or more keys from a configured server.
type Node struct {
	ID          string         // Node ID (random string)
	privKey     string         // The private key for the node, used in ZMQ CURVE Auth.
	pubKey      string         // The public key for the node, used in ZMQ CURVE Auth.
	host        *net.IPAddr    // A string representing the IP the node is going to use to listen to requests.
	port        uint16         // a int representing the port the node is going to use to listen to requests
	config      *config.Config // A pointer to the struct which saves the configuration of the node.
	context     *zmq4.Context  // The context used by zmq connections.
	clients     []*Client      // A list of clients. Currently the configuration allows only one server at a time.
	configMutex sync.Mutex     // A mutex used for config editing.
	socket      *zmq4.Socket   // The socket where the messages are received and sent to the server.
}

func init() {
	zmq4.AuthSetVerbose(true)
}

// InitNode inits the node using the configuration provided. Returns a started node or an error if the function fails.
func InitNode(config *config.Config) (*Node, error) {
	ip, err := net.ResolveIPAddr("ip", config.Host)
	if err != nil {
		return nil, err
	}
	nodeID, err := message.GetRandomHexString(8)
	if err != nil {
		return nil, err
	}
	node := &Node{
		ID:      nodeID,
		pubKey:  config.PublicKey,
		privKey: config.PrivateKey,
		host:    ip,
		port:    config.Port,
		config:  config,
		clients: make([]*Client, 0),
	}
	log.Printf("Creating node with ID: %s", node.GetID())
	context, err := zmq4.NewContext()
	if err != nil {
		return nil, err
	}
	node.context = context

	ips, err := config.GetClientIPs()
	if err != nil {
		return nil, err
	}
	zmq4.AuthAllow(TchsmDomain, ips...)
	zmq4.AuthCurveAdd(TchsmDomain, config.GetClientPubKeys()...)

	s, err := context.NewSocket(zmq4.REP)
	if err != nil {
		return nil, err
	}
	node.socket = s

	if err := node.socket.SetIdentity(node.GetID()); err != nil {
		return nil, err
	}

	if err := node.socket.ServerAuthCurve(TchsmDomain, node.privKey); err != nil {
		return nil, err
	}

	log.Printf("Listening messages in %s", node.GetConnString())
	if err := node.socket.Bind(node.GetConnString()); err != nil {
		return nil, err
	}

	serverConfig := config.Client

	serverIP, err := net.ResolveIPAddr("ip", serverConfig.Host)
	if err != nil {
		return nil, err
	}
	server := &Client{
		pubKey: serverConfig.PublicKey,
		host:   serverIP,
		node:   node,
		keys:   make(map[string]*Key, len(serverConfig.Keys)),
	}

	for _, key := range serverConfig.Keys {
		var keyShare *tcrsa.KeyShare
		var keyMeta *tcrsa.KeyMeta
		if key.KeyShare != "" && key.KeyMetaInfo != "" {
			keyShareByte, err := base64.StdEncoding.DecodeString(key.KeyShare)
			if err != nil {
				return nil, err
			}
			keyShare, err = message.DecodeKeyShare(keyShareByte)
			if err != nil {
				return nil, err
			}
			keyMetaByte, err := base64.StdEncoding.DecodeString(key.KeyMetaInfo)
			if err != nil {
				return nil, err
			}
			keyMeta, err = message.DecodeKeyMeta(keyMetaByte)
			if err != nil {
				return nil, err
			}
		}
		server.keys[key.ID] = &Key{
			ID:    key.ID,
			Meta:  keyMeta,
			Share: keyShare,
		}
	}

	node.clients = append(node.clients, server)

	return node, nil
}

// GetID returns the ID of the node.
func (node *Node) GetID() string {
	return node.ID
}

// FindServer returns a server with the provided ID, or nil if it doesn't exist.
func (node *Node) FindServer(name string) *Client {
	for _, server := range node.clients {
		if server.pubKey == name {
			return server
		}
	}
	return nil
}

// GetConnString returns the string that is used to bind the node to a port.
func (node *Node) GetConnString() string {
	return fmt.Sprintf("%s://%s:%d", TchsmProtocol, node.host, node.port)
}

// SaveConfigKeys saves the currently received keys into memory.
func (node *Node) SaveConfigKeys() error {
	node.configMutex.Lock()
	defer node.configMutex.Unlock()
	for _, client := range node.clients {
		serverConfig := node.config.GetClientByID(client.GetID())
		if serverConfig == nil {
			return fmt.Errorf("error encoding keys: client config not found")
		}
		serverConfig.Keys = make([]*config.KeyConfig, 0)
		for _, key := range client.keys {
			keyShareBytes, err := message.EncodeKeyShare(key.Share)
			if err != nil {
				return fmt.Errorf("error encoding keys: %s", err)
			}
			keyMetaBytes, err := message.EncodeKeyMeta(key.Meta)
			if err != nil {
				return fmt.Errorf("error encoding keys: %s", err)
			}
			keyShareB64 := base64.StdEncoding.EncodeToString(keyShareBytes)
			keyMetaB64 := base64.StdEncoding.EncodeToString(keyMetaBytes)
			serverConfig.Keys = append(serverConfig.Keys, &config.KeyConfig{
				ID:          key.ID,
				KeyMetaInfo: keyMetaB64,
				KeyShare:    keyShareB64,
			})
		}
	}
	viper.Set("config", node.config)
	return viper.WriteConfig()
}

// Listen starts all the server listening subroutines, and waits for a message received in the input socket. It checks and parses the messages to Message objects and sends them to a channel, that is used by the subroutines.
func (node *Node) Listen() {
	for _, client := range node.clients {
		client.Listen()
	}
	select {}
}