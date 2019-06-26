package main

import (
	"encoding/base64"
	"fmt"
	"github.com/niclabs/dtcnode/config"
	"github.com/niclabs/dtcnode/message"
	"github.com/niclabs/tcrsa"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"log"
	"sync"
)

// The domain of the ZMQ connection. This value must be the same in the server, or it will not work.
const TchsmDomain = "tchsm"

// The protocol used for the ZMQ connection. TCP is the best for this usage cases.
const TchsmProtocol = "tcp"

// Node represents a node in the distributed TCHSM application. It saves zero or more keys from a configured server.
type Node struct {
	privKey     string         // The private key for the node, used in ZMQ CURVE Auth.
	pubKey      string         // The public key for the node, used in ZMQ CURVE Auth.
	host        string         // A string representing the IP the node is going to use to listen to requests.
	port        uint16         // a int representing the port the node is going to use to listen to requests
	config      *config.Config // A pointer to the struct which saves the configuration of the node.
	context     *zmq4.Context  // The context used by zmq connections.
	socket      *zmq4.Socket   //  The socket used to receive messages.
	servers     []*Server      // A list of servers. Currently the configuration allows only one server at a time.
	configMutex sync.Mutex     // A mutex used for config editing.
}

// InitNode inits the node using the configuration provided. Returns a started node or an error if the function fails.
func InitNode(config *config.Config) (*Node, error) {

	node := &Node{
		pubKey:  config.PublicKey,
		privKey: config.PrivateKey,
		host:    config.Host,
		port:    config.Port,
		config:  config,
		servers: make([]*Server, 0),
	}

	context, err := zmq4.NewContext()
	if err != nil {
		return nil, err
	}
	node.context = context

	zmq4.AuthAllow(TchsmDomain, config.GetServerIPs()...)
	zmq4.AuthCurveAdd(TchsmDomain, config.GetServerPubKeys()...)

	in, err := context.NewSocket(zmq4.ROUTER)
	if err != nil {
		return nil, err
	}
	if err := in.SetIdentity(node.GetID()); err != nil {
		return nil, err
	}
	if err := in.ServerAuthCurve(TchsmDomain, node.privKey); err != nil {
		return nil, err
	}
	if err := in.Bind(node.GetConnString()); err != nil {
		return nil, err
	}
	node.socket = in

	serverConfig := config.Server
	server := &Server{
		pubKey:  serverConfig.PublicKey,
		host:    serverConfig.Host,
		port:    serverConfig.Port,
		client:  node,
		channel: make(chan *message.Message),
		keys:    make(map[string]*Key, len(serverConfig.Keys)),
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

	out, err := context.NewSocket(zmq4.DEALER)
	if err != nil {
		return nil, err
	}
	if err := out.SetIdentity(node.GetID()); err != nil {
		return nil, err
	}
	if err := out.ClientAuthCurve(serverConfig.PublicKey, node.pubKey, node.privKey); err != nil {
		return nil, err
	}
	if err := out.Connect(server.GetConnString()); err != nil {
		return nil, err
	}
	server.socket = out
	node.servers = append(node.servers, server)

	return node, nil
}

// GetID returns the ID of the node.
func (client *Node) GetID() string {
	return client.pubKey
}

// FindServer returns a server with the provided ID, or nil if it doesn't exist.
func (client *Node) FindServer(name string) *Server {
	for _, server := range client.servers {
		if server.pubKey == name {
			return server
		}
	}
	return nil
}

// GetConnString returns the string that is used to bind the client to a port.
func (client *Node) GetConnString() string {
	return fmt.Sprintf("%s://%s:%d", TchsmProtocol, client.host, client.port)
}

// SaveConfigKeys saves the currently received keys into memory.
func (client *Node) SaveConfigKeys() error {
	client.configMutex.Lock()
	defer client.configMutex.Unlock()
	for _, server := range client.servers {
		serverConfig := client.config.GetServerByID(server.GetID())
		if serverConfig == nil {
			return fmt.Errorf("error encoding keys: server config not found")
		}
		for _, key := range server.keys {
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
			keyConfig := serverConfig.GetKeyByID(key.ID)
			if keyConfig == nil {
				serverConfig.Keys = append(serverConfig.Keys, &config.KeyConfig{
					ID:          key.ID,
					KeyMetaInfo: keyMetaB64,
					KeyShare:    keyShareB64,
				})
			}
		}
	}
	viper.Set("config", client.config)
	return viper.WriteConfig()
}

// Listen starts all the server listening subroutines, and waits for a message received in the input socket. It checks and parses the messages to Message objects and sends them to a channel, that is used by the subroutines.
func (client *Node) Listen() {
	for _, server := range client.servers {
		go server.Listen()
	}

	for {
		rawMsg, err := client.socket.RecvMessageBytes(0)
		if err != nil {
			log.Printf("%s", message.ReceiveMessageError.ComposeError(err))
			continue
		}
		log.Printf("message from client %s", rawMsg[0])
		log.Printf("parsing message")
		msg, err := message.FromBytes(rawMsg)
		if err != nil {
			log.Printf("%s", message.ParseMessageError.ComposeError(err))
			continue
		}
		if server := client.FindServer(msg.NodeID); server != nil {
			log.Printf("nodeID %s is valid, sending to channel", msg.NodeID)
			server.channel <- msg
		} else {
			log.Printf("could not send message to server listener %s: not found", msg.NodeID)
		}
	}
}
