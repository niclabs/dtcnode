package main

import (
	"crypto"
	"fmt"
	"github.com/niclabs/dtcnode/message"
	"github.com/niclabs/tcrsa"
	"github.com/pebbe/zmq4"
	"log"
)

type Server struct {
	host    string
	port    uint16
	pubKey  string
	keys    map[string]*Key
	client  *Client
	socket  *zmq4.Socket
	channel chan *message.Message
}

func (server *Server) GetID() string {
	return server.pubKey
}

func (server *Server) GetConnString() string {
	return fmt.Sprintf("%s://%s:%d", TchsmProtocol, server.host, server.port)
}

func (server *Server) Listen() {
	log.Printf("Listening messages from server %s", server.GetConnString())
	for msg := range server.channel {
		resp := msg.CopyWithoutData(message.Ok)
		switch msg.Type {
		case message.SendKeyShare:
			log.Printf("SendKeyShare message received from server %s", server.GetConnString())
			if len(msg.Data) != 3 { // keyID, keyshare, sigshare
				resp.Error = message.InvalidMessageError
				break
			}

			keyID := string(msg.Data[0])
			keyShare, err := message.DecodeKeyShare(msg.Data[1])
			if err != nil {
				resp.Error = message.KeyShareDecodeError
				break
			}
			keyMeta, err := message.DecodeKeyMeta(msg.Data[2])
			if err != nil {
				resp.Error = message.KeyMetaDecodeError
				break
			}
			server.SaveKey(keyID, keyShare, keyMeta)
		case message.AskForSigShare:
			log.Printf("AskForSigShare message received from server %s", server.GetConnString())

			if len(msg.Data) != 2 {
				resp.Error = message.InvalidMessageError
				break
			}
			keyID := string(msg.Data[0])
			key, ok := server.keys[keyID]
			if !ok {
				resp.Error = message.NotInitializedError
				break
			}
			doc := msg.Data[1]
			sigShare, err := key.Share.Sign(doc, crypto.SHA256, key.Meta)
			if err != nil {
				resp.Error = message.DocSignError
				break
			}
			encodedSigShare, err := message.EncodeSigShare(sigShare)
			if err != nil {
				resp.Error = message.SigShareEncodeError
				break
			}
			resp.AddMessage(encodedSigShare)
		default:
			resp.Error = message.InvalidMessageError
		}
		if resp.Error != message.Ok {
			log.Printf("%s", resp.Error.Error())
		}
		log.Printf("sending response to server %s", server.GetConnString())
		if _, err := server.socket.SendMessage(resp.GetBytesLists()...); err != nil {
			log.Printf("%s", err.Error())
		}
		log.Printf("sent response to server %s", server.GetConnString())
	}
}

func (server *Server) SaveKey(id string, keyShare *tcrsa.KeyShare, keyMeta *tcrsa.KeyMeta) {
	key, ok := server.keys[id]
	if !ok {
		key = &Key{}
		server.keys[id] = key
	}
	key.Meta = keyMeta
	key.Share = keyShare
	server.client.SaveConfigKeys()
}

type Key struct {
	ID    string
	Share *tcrsa.KeyShare
	Meta  *tcrsa.KeyMeta
}
