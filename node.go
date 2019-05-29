package main

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"net"
)

type Node struct {
	privKey      string
	pubKey       string
	ip           net.IP
	context      *zmq4.Context
	routerSocket *zmq4.Socket
	subSocket    *zmq4.Socket
}

func NewNode(pubkey, privkey, ipstr string, routerPort, subPort uint16) (*Node, error) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return nil, fmt.Errorf("ip not well formatted")
	}
	node := &Node{
		pubKey:  pubkey,
		privKey: privkey,
		ip:      ip,
	}
	context, err := zmq4.NewContext()
	if err != nil {
		return nil, err
	}
	node.context = context

	// SUB socket
	subSocket, err := context.NewSocket(zmq4.SUB)
	if err != nil {
		return nil, err
	}
	if err := subSocket.SetIdentity("node"); err != nil {
		return nil, err
	}
	if err := subSocket.ServerAuthCurve("node", node.privKey); err != nil {
		return nil, err
	}
	if err := subSocket.Bind(fmt.Sprintf("tcp://%s:%d", node.ip, subPort)); err != nil {
		return nil, err
	}
	node.subSocket = subSocket

	// ROUTER socket
	routerSocket, err := context.NewSocket(zmq4.ROUTER)
	if err != nil {
		return nil, err
	}
	if err := routerSocket.SetIdentity("node"); err != nil {
		return nil, err
	}
	if err := routerSocket.ServerAuthCurve("node", node.privKey); err != nil {
		return nil, err
	}
	if err := subSocket.Bind(fmt.Sprintf("tcp://%s:%d", node.ip, subPort)); err != nil {
		return nil, err
	}
	node.routerSocket = routerSocket

	return node, nil
}


func (n *Node) Listen() {
	// Infinite loop to check for messages from SUB
	// In case of receiving a message, send the response to ROUTER
}