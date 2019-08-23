package server

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"github.com/niclabs/dtcnode/message"
	"github.com/niclabs/tcrsa"
	"log"
	"net"
)

// Client represents the connection with the Distributed TCHSM server.
// It saves its connection values, its public key, and the keyshares and keymetainfo sent by the server.
type Client struct {
	host   *net.IPAddr     // IP where the server is listening.
	pubKey string          // Public key of the server. Used for SMQ CURVE auth.
	keys   map[string]*Key // Dictionary with key shares created by this server.
	node   *Node           // A pointer to the node that manages this server subroutine.
}

// Key represents a keyshare managed by the node and used by the server for signing documents.
type Key struct {
	ID    string
	Share *tcrsa.KeyShare
	Meta  *tcrsa.KeyMeta
}

// GetID returns the id of the server.
func (client *Client) GetID() string {
	return client.pubKey
}

// GetConnString returns the string that identifies the client.
func (client *Client) GetConnString() string {
	return fmt.Sprintf("%s://%s", TchsmProtocol, client.host)
}

// Listen is the subroutine that keeps waiting for messages on its channel. Then it acts depending on each message.
func (client *Client) Listen() {
	for {
		log.Printf("Waiting for messages...")
		rawMsg, err := client.node.socket.RecvMessageBytes(0)
		if err != nil {
			log.Printf("%s", message.ReceiveMessageError.ComposeError(err))
			continue
		}
		log.Printf("message from node %s", rawMsg[0])
		log.Printf("parsing message")
		msg, err := message.FromBytes(rawMsg)
		if err != nil {
			log.Printf("%s", message.ParseMessageError.ComposeError(err))
			continue
		}

		resp := msg.CopyWithoutData(client.node.GetID(), message.Ok)

		switch msg.Type {
		case message.SendKeyShare:
			log.Printf("Client %s is sending us a new KeyShare", client.GetConnString())
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
			if err := client.SaveKey(keyID, keyShare, keyMeta); err != nil {
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
			log.Printf("Client %s is asking us for a signature share using key %s", client.GetConnString(), keyID)
			key, ok := client.keys[keyID]
			if !ok {
				resp.Error = message.NotInitializedError
				break
			}
			doc := msg.Data[1]
			b64doc := base64.StdEncoding.EncodeToString(doc)
			log.Printf("Signing document hash %s with key %s as asked by client %s", b64doc, keyID, client.GetConnString())
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

			log.Printf("The document %s was signed succesfully with key %s as asked by client %s", b64doc, keyID, client.GetConnString())
			encodedSigShare, err := message.EncodeSigShare(sigShare)
			if err != nil {
				resp.Error = message.SigShareEncodeError
				break
			}
			resp.AddMessage(encodedSigShare)
		case message.DeleteKeyShare:
			log.Printf("Client %s is asking us to delete a KeyShare", client.GetConnString())
			if len(msg.Data) != 1 { // keyID
				resp.Error = message.InvalidMessageError
				break
			}
			keyID := string(msg.Data[0])
			log.Printf("Deleting keyshare for keyid=%s", keyID)
			if err := client.DeleteKey(keyID); err != nil {
				log.Printf("Error with key saving: %s", err)
				resp.Error = message.InternalError
				break
			}
			log.Printf("Keyshare deleted for keyid=%s", keyID)
		default:
			log.Printf("invalid message received from client %s", client.GetConnString())

			resp.Error = message.InvalidMessageError
		}
		if resp.Error != message.Ok {
			log.Printf("Error processing message: %s", resp.Error.Error())
		}
		log.Printf("sending response to client %s", client.GetConnString())
		if _, err := client.node.socket.SendMessage(resp.GetBytesLists()...); err != nil {
			log.Printf("%s", err.Error())
		}
		log.Printf("A response to client %s was sent", client.GetConnString())
	}
}

// SaveKey updates the key array of the server and asks the node to save the keys into the config file.
func (client *Client) SaveKey(id string, keyShare *tcrsa.KeyShare, keyMeta *tcrsa.KeyMeta) error {
	key, ok := client.keys[id]
	if !ok {
		key = &Key{}
		client.keys[id] = key
	}
	key.ID = id
	key.Meta = keyMeta
	key.Share = keyShare
	return client.node.SaveConfigKeys()
}

// SaveKey deletes a key from the array of the server and asks the node to save the new key array into the config file.
func (client *Client) DeleteKey(id string) error {
	delete(client.keys, id)
	return client.node.SaveConfigKeys()
}
