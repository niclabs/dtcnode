package main

import (
	"dtcnode/message"
	"fmt"
	"github.com/pebbe/zmq4"
	"net"
	"os"
	"sync"
)

const TchsmDomain = "tchsm"
const TchsmProtocol = "tcp"

type Client struct {
	privKey     string
	pubKey      string
	ip          net.IP
	port        uint16
	config      *Config
	context     *zmq4.Context
	socket      *zmq4.Socket
	servers     map[string]*Server
	configMutex sync.Mutex
}

func InitClient(config *Config) (*Client, error) {

	ip := net.ParseIP(config.IP)
	if ip == nil {
		return nil, fmt.Errorf("invalid ip format")
	}
	node := &Client{
		pubKey:  config.PublicKey,
		privKey: config.PrivateKey,
		ip:      ip,
		port:    config.Port,
		config:  config,
		servers: make(map[string]*Server, len(config.Servers)),
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

	for _, serverConfig := range config.Servers {
		serverIP := net.ParseIP(config.IP)
		server := &Server{
			pubKey: config.PublicKey,
			ip:     &serverIP,
			port:   config.Port,
			client: node,
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
		if err := in.Connect(server.GetConnString()); err != nil {
			return nil, err
		}
		server.socket = in
		node.servers[serverConfig.PublicKey] = server
	}
	node.socket = in

	return node, nil
}

func (client *Client) GetID() string {
	return client.pubKey
}

func (client *Client) GetConnString() string {
	return fmt.Sprintf("%s://%s:%d", TchsmProtocol, client.ip, client.port)
}

func (client *Client) SaveConfig() string {
	client.configMutex.Lock()
	defer client.configMutex.Unlock()

}

func (client *Client) Listen() {
	for _, server := range client.servers {
		go server.Listen()
	}

	for {
		rawMsg, err := client.socket.RecvMessageBytes(0)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", message.ReceiveMessageError.ComposeError(err))
			continue
		}
		_, _ = fmt.Fprintf(os.Stderr, "message from client %s\n", rawMsg[0])
		_, _ = fmt.Fprintf(os.Stderr, "parsing message...\n")
		msg, err := message.FromBytes(rawMsg)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", message.ParseMessageError.ComposeError(err))
			continue
		}

		if server, ok := client.servers[msg.NodeID]; ok {
			server.channel <- msg
		}
	}
}
