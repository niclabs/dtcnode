package main

import (
	"bytes"
	"crypto"
	"dtcnode/message"
	"encoding/gob"
	"fmt"
	"github.com/niclabs/tcrsa"
	"github.com/pebbe/zmq4"
	"log"
	"net"
	"os"
)

type Server struct {
	ip       *net.IP
	port     uint16
	pubKey   string
	keyShare *tcrsa.KeyShare
	keyMeta  *tcrsa.KeyMeta
	client   *Client
	socket   *zmq4.Socket
	channel  chan *message.Message
}


func (server *Server) GetID() string {
	return server.pubKey
}

func (server *Server) GetConnString() string {
	return fmt.Sprintf("%s://%s:%d", TchsmProtocol, server.ip, server.port)
}


func (server *Server) Listen() string {
	for msg := range server.channel {
		resp := msg.CopyWithoutData(message.Ok)
		switch msg.Type {
		case message.SendKeyShare:
			var keyShare tcrsa.KeyShare
			var keyMeta tcrsa.KeyMeta
			if server.keyShare != nil || server.keyMeta != nil {
				resp.Error = message.AlreadyInitializedError
				break
			}
			if len(msg.Data) != 2 {
				resp.Error = message.KeyShareDecodeError
				break
			}
			keyShareBuffer := bytes.NewBuffer(msg.Data[0])
			keyShareDecoder := gob.NewDecoder(keyShareBuffer)
			if err := keyShareDecoder.Decode(&keyShare); err != nil {
				resp.Error = message.KeyShareDecodeError
				break
			}
			keyMetaBuffer := bytes.NewBuffer(msg.Data[1])
			keyMetaDecoder := gob.NewDecoder(keyMetaBuffer)
			if err := keyMetaDecoder.Decode(&keyMeta); err != nil {
				server.keyShare = nil
				resp.Error = message.KeyMetaDecodeError
				break
			}
			server.SaveKey(&keyShare, &keyMeta)
		case message.AskForSigShare:
			if server.keyShare == nil || server.keyMeta == nil {
				resp.Error = message.NotInitializedError
				break
			}
			doc := msg.Data[0]
			sigShare, err := server.keyShare.Sign(doc, crypto.SHA256, &server.keyMeta)
			if err != nil {
				resp.Error = message.DocSignError
				break
			}
			var keyBuffer bytes.Buffer
			if err := gob.NewEncoder(&keyBuffer).Encode(sigShare); err != nil {
				resp.Error = message.SigShareEncodeError
				break
			}
			resp.AddMessage(keyBuffer.Bytes())
		default:
			resp.Error = message.InvalidMessageError
		}
		if resp.Error != message.Ok {
			log.Printf("%s", resp.Error.Error())
		}
		if _, err := server.socket.SendMessage(resp.GetBytesLists()...); err != nil {
			log.Printf("%s", err.Error())
		}
	}
}

func (server *Server) SaveKey(keyShare *tcrsa.KeyShare, keyMeta *tcrsa.KeyMeta) {
	server.keyShare = keyShare
	server.keyMeta = keyMeta
	server.client.SaveConfig()
	// Encode keyshare and keymeta and save them in client config
}
