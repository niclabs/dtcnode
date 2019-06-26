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
	log.Printf("Listening messages from server %s in %s", server.GetConnString(), server.client.GetConnString())
	for msg := range server.channel {
		resp := msg.CopyWithoutData(message.Ok)
		switch msg.Type {
		case message.SendKeyShare:
			log.Printf("Server %s is sending us a new KeyShare", server.GetConnString())
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
			log.Printf("Saving keyshare for keyid=%s", keyID)
			if err := server.SaveKey(keyID, keyShare, keyMeta); err != nil {
				log.Printf("Error with key saving: %s", err)
				break
			}
			log.Printf("Keyshare saved for keyid=%s", keyID)
		case message.AskForSigShare:
			if len(msg.Data) != 2 {
				resp.Error = message.InvalidMessageError
				break
			}
			keyID := string(msg.Data[0])
			log.Printf("Server %s is asking us for a signature share using key %s", server.GetConnString(), keyID)
			key, ok := server.keys[keyID]
			if !ok {
				resp.Error = message.NotInitializedError
				break
			}
			doc := msg.Data[1]
			log.Printf("Signing document hash %s with key %s as asked by server %s", doc, keyID, server.GetConnString())
			sigShare, err := key.Share.Sign(doc, crypto.SHA256, key.Meta)
			if err != nil {
				resp.Error = message.DocSignError
				break
			}
			log.Printf("The document hash %s was signed succesfully with key %s as asked by server %s", doc, keyID, server.GetConnString())
			encodedSigShare, err := message.EncodeSigShare(sigShare)
			if err != nil {
				resp.Error = message.SigShareEncodeError
				break
			}
			resp.AddMessage(encodedSigShare)
		default:
			log.Printf("invalid message received from server %s", server.GetConnString())

			resp.Error = message.InvalidMessageError
		}
		if resp.Error != message.Ok {
			log.Printf("%s", resp.Error.Error())
		}
		log.Printf("sending response to server %s", server.GetConnString())
		if _, err := server.socket.SendMessage(resp.GetBytesLists()...); err != nil {
			log.Printf("%s", err.Error())
		}
		log.Printf("A response to server %s was sent", server.GetConnString())
	}
}

func (server *Server) SaveKey(id string, keyShare *tcrsa.KeyShare, keyMeta *tcrsa.KeyMeta) error {
	key, ok := server.keys[id]
	if !ok {
		key = &Key{}
		server.keys[id] = key
	}
	key.Meta = keyMeta
	key.Share = keyShare
	return server.client.SaveConfigKeys()
}

type Key struct {
	ID    string
	Share *tcrsa.KeyShare
	Meta  *tcrsa.KeyMeta
}
