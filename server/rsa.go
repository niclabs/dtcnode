package server

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"github.com/niclabs/dtcnode/v3/config"
	"github.com/niclabs/dtcnode/v3/message"
	"github.com/niclabs/tcrsa"
	"log"
)

// rsa represents the data related to rsa signing processes
type rsa struct {
	keys map[string]*rsaKey
}

// rsaKey represents a keyshare managed by the node and used by the server for signing documents.
type rsaKey struct {
	ID    string
	Share *tcrsa.KeyShare
	Meta  *tcrsa.KeyMeta
}

func (client *Client) dispatchRSA(msg *message.Message) *message.Message {
	resp := msg.NewResponse(client.node.GetID(), message.Ok)
	switch msg.Type {
	case message.SendRSAKeyShare:
		log.Printf("Client %s is sending us a new RSA KeyShare", client.GetConnString())
		keyID := string(msg.Data[0])
		keyShare, err := message.DecodeRSAKeyShare(msg.Data[1])
		if err != nil {
			resp.Error = message.DecodingError
			break
		}
		keyMeta, err := message.DecodeRSAKeyMeta(msg.Data[2])
		if err != nil {
			resp.Error = message.DecodingError
			break
		}
		log.Printf("Saving keyshare for keyid=%s", keyID)
		if err := client.SaveRSAKey(keyID, keyShare, keyMeta); err != nil {
			log.Printf("Error with RSA keyshare saving process: %s", err)
			resp.Error = message.InternalError
			break
		}
		log.Printf("Keyshare saved for keyid=%s", keyID)
	case message.GetRSASigShare:
		keyID := string(msg.Data[0])
		log.Printf("Client %s is asking us for a RSA signature share using key %s", client.GetConnString(), keyID)
		key, ok := client.rsa.keys[keyID]
		if !ok {
			resp.Error = message.KeyNotFoundError
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
		encodedSigShare, err := message.EncodeRSASigShare(sigShare)
		if err != nil {
			resp.Error = message.EncodingError
			break
		}
		resp.AddMessage(encodedSigShare)
	case message.DeleteRSAKeyShare:
		log.Printf("Client %s is asking us to delete a RSA KeyShare", client.GetConnString())
		keyID := string(msg.Data[0])
		log.Printf("Deleting keyshare for keyid=%s", keyID)
		if err := client.DeleteRSAKey(keyID); err != nil {
			log.Printf("Error with key deleting: %s", err)
			resp.Error = message.InternalError
			break
		}
		log.Printf("Keyshare deleted for keyid=%s", keyID)
	default:
		log.Printf("invalid message received from client %s", client.GetConnString())

		resp.Error = message.InvalidMessageError
	}
	return resp
}

// SaveRSAKey updates the key array of the server and asks the node to save the rsaKeys into the config file.
func (client *Client) SaveRSAKey(id string, keyShare *tcrsa.KeyShare, keyMeta *tcrsa.KeyMeta) error {
	key, ok := client.rsa.keys[id]
	if !ok {
		key = &rsaKey{}
		client.rsa.keys[id] = key
	}
	key.ID = id
	key.Meta = keyMeta
	key.Share = keyShare
	return client.node.SaveConfigKeys()
}

// SaveRSAKey deletes a key from the array of the server and asks the node to save the new key array into the config file.
func (client *Client) DeleteRSAKey(id string) error {
	log.Printf("deleting ecdsa key with id %s", id)
	delete(client.rsa.keys, id)
	return client.node.SaveConfigKeys()
}

func parseRSAKeys(conf []*config.RSAKeyConfig) (map[string]*rsaKey, error) {
	keys := make(map[string]*rsaKey)
	for _, key := range conf {
		var keyShare *tcrsa.KeyShare
		var keyMeta *tcrsa.KeyMeta
		if key.KeyShare != "" && key.KeyMetaInfo != "" {
			keyShareByte, err := base64.StdEncoding.DecodeString(key.KeyShare)
			if err != nil {
				return nil, err
			}
			keyShare, err = message.DecodeRSAKeyShare(keyShareByte)
			if err != nil {
				return nil, err
			}
			keyMetaByte, err := base64.StdEncoding.DecodeString(key.KeyMetaInfo)
			if err != nil {
				return nil, err
			}
			keyMeta, err = message.DecodeRSAKeyMeta(keyMetaByte)
			if err != nil {
				return nil, err
			}
		}
		keys[key.ID] = &rsaKey{
			ID:    key.ID,
			Meta:  keyMeta,
			Share: keyShare,
		}
	}
	return keys, nil
}

func saveRSAKeys(keys map[string]*rsaKey) ([]*config.RSAKeyConfig, error) {
	keysConfig := make([]*config.RSAKeyConfig, 0)
	for _, key := range keys {
		keyShareBytes, err := message.EncodeRSAKeyShare(key.Share)
		if err != nil {
			return nil, fmt.Errorf("error encoding rsaKeys: %s", err)
		}
		keyMetaBytes, err := message.EncodeRSAKeyMeta(key.Meta)
		if err != nil {
			return nil, fmt.Errorf("error encoding rsaKeys: %s", err)
		}
		keyShareB64 := base64.StdEncoding.EncodeToString(keyShareBytes)
		keyMetaB64 := base64.StdEncoding.EncodeToString(keyMetaBytes)
		keysConfig = append(keysConfig, &config.RSAKeyConfig{
			ID:          key.ID,
			KeyMetaInfo: keyMetaB64,
			KeyShare:    keyShareB64,
		})
	}
	return keysConfig, nil
}
