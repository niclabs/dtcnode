package main

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"github.com/niclabs/dtcnode/message"
	"github.com/niclabs/tcrsa"
	"log"
	"net"
)

// Server represents the connection with the Distributed TCHSM server.
// It saves its connection values, its public key, and the keyshares and keymetainfo sent by the server.
type Server struct {
	host   *net.IPAddr     // IP where the server is listening.
	pubKey string          // Public key of the server. Used for SMQ CURVE auth.
	keys   map[string]*Key // Dictionary with key shares created by this server.
	client *Node           // A pointer to the node that manages this server subroutine.
}

// Key represents a keyshare managed by the node and used by the server for signing documents.
type Key struct {
	ID    string
	Share *tcrsa.KeyShare
	Meta  *tcrsa.KeyMeta
}

// GetID returns the id of the server.
func (server *Server) GetID() string {
	return server.pubKey
}

// GetConnString returns the string that is used for connecting to the server.
func (server *Server) GetConnString() string {
	return fmt.Sprintf("%s://%s", TchsmProtocol, server.host)
}

// Listen is the subroutine that keeps waiting for messages on its channel. Then it acts depending on each message.
func (server *Server) Listen() {
	log.Printf("Listening messages in %s", server.client.GetConnString())
	for {
		rawMsg, err := server.client.socket.RecvMessageBytes(0)
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
				resp.Error = message.InternalError
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
			b64doc := base64.StdEncoding.EncodeToString(doc)
			log.Printf("Signing document hash %s with key %s as asked by server %s", b64doc, keyID, server.GetConnString())
			sigShare, err := key.Share.Sign(doc, crypto.SHA256, key.Meta)
			if err != nil {
				resp.Error = message.DocSignError
				break
			}
			// Verify sigshare locally
			if err := sigShare.Verify(doc, key.Meta); err != nil {
				resp.Error = message.DocSignError
				break
			}

			log.Printf("The document %s was signed succesfully with key %s as asked by server %s", b64doc, keyID, server.GetConnString())
			encodedSigShare, err := message.EncodeSigShare(sigShare)
			if err != nil {
				resp.Error = message.SigShareEncodeError
				break
			}
			resp.AddMessage(encodedSigShare)
		case message.DeleteKeyShare:
			log.Printf("Server %s is asking us to delete a KeyShare", server.GetConnString())
			if len(msg.Data) != 1 { // keyID
				resp.Error = message.InvalidMessageError
				break
			}
			keyID := string(msg.Data[0])
			log.Printf("Deleting keyshare for keyid=%s", keyID)
			if err := server.DeleteKey(keyID); err != nil {
				log.Printf("Error with key saving: %s", err)
				resp.Error = message.InternalError
				break
			}
			log.Printf("Keyshare deleted for keyid=%s", keyID)
		default:
			log.Printf("invalid message received from server %s", server.GetConnString())

			resp.Error = message.InvalidMessageError
		}
		if resp.Error != message.Ok {
			log.Printf("Error processing message: %s", resp.Error.Error())
		}
		log.Printf("sending response to server %s", server.GetConnString())
		if _, err := server.client.socket.SendMessage(resp.GetBytesLists()...); err != nil {
			log.Printf("%s", err.Error())
		}
		log.Printf("A response to server %s was sent", server.GetConnString())
	}
}

// SaveKey updates the key array of the server and asks the node to save the keys into the config file.
func (server *Server) SaveKey(id string, keyShare *tcrsa.KeyShare, keyMeta *tcrsa.KeyMeta) error {
	key, ok := server.keys[id]
	if !ok {
		key = &Key{}
		server.keys[id] = key
	}
	key.ID = id
	key.Meta = keyMeta
	key.Share = keyShare
	return server.client.SaveConfigKeys()
}

// SaveKey deletes a key from the array of the server and asks the node to save the new key array into the config file.
func (server *Server) DeleteKey(id string) error {
	delete(server.keys, id)
	return server.client.SaveConfigKeys()
}
