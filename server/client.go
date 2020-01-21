package server

import (
	"fmt"
	"github.com/niclabs/dtcnode/v3/message"
	"log"
	"net"
)

// Client represents the connection with the Distributed TCHSM server.
// It saves its connection values, its public key, and the keyshares and keymetainfo sent by the server.
type Client struct {
	host   *net.IPAddr // IP where the server is listening.
	pubKey string      // Public key of the server. Used for SMQ CURVE auth.
	rsa    rsa         // struct with RSA structures, as keys.
	ecdsa  ecdsa       // struct with ECDSA structures, as keys and the active currentSession.
	node   *Node       // A pointer to the node that manages this server subroutine.
}

// GetID returns the id of the server.
func (client *Client) GetID() string {
	return client.pubKey
}

// GetConnString returns the string that identifies the client.
func (client *Client) GetConnString() string {
	return fmt.Sprintf("%s://%s", TchsmProtocol, client.host)
}

// Listen is the subroutine that keeps waiting for message on its channel. Then it acts depending on each message.
func (client *Client) Listen() {
	for {
		log.Printf("Waiting for message...")
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
		var resp *message.Message
		if !msg.ValidClientDataLength() {
			resp.Error = message.InvalidMessageError
		} else {
			if msg.Type.IsRSA() {
				resp = client.dispatchRSA(msg)
			} else if msg.Type.IsECDSA() {
				resp = client.dispatchECDSA(msg)
			} else {
				log.Printf("Unknown message of type %d. Ignored.", msg.Type)
				continue
			}
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
