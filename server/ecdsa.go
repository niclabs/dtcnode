package server

import (
	"encoding/base64"
	"fmt"
	"github.com/niclabs/dtcnode/v3/config"
	"github.com/niclabs/dtcnode/v3/message"
	"github.com/niclabs/tcecdsa"
	"log"
)

type ecdsa struct {
	keys           map[string]*ecdsaKey
	currentKey     string
	currentSession *tcecdsa.SigSession
}

// ecdsaKey represents a keyshare managed by the node and used by the server for signing documents.
type ecdsaKey struct {
	ID        string
	Completed bool
	Share     *tcecdsa.KeyShare
	Meta      *tcecdsa.KeyMeta
}

func (client *Client) dispatchECDSA(msg *message.Message) *message.Message {
	resp := msg.NewResponse(client.node.GetID(), message.Ok)
	switch msg.Type {
	case message.SendECDSAKeyShare:
		keyID := string(msg.Data[0])
		log.Printf("Client %s is sending us a new incomplete ECDSA KeyShare with id=%s", client.GetConnString(), keyID)
		keyShare, err := message.DecodeECDSAKeyShare(msg.Data[1])
		if err != nil {
			log.Printf("error decoding ECDSA KeyShare message: %s", err)
			resp.Error = message.DecodingError
			break
		}
		keyMeta, err := message.DecodeECDSAKeyMeta(msg.Data[2])
		if err != nil {
			log.Printf("error decoding ECDSA KeyMeta message: %s", err)
			resp.Error = message.DecodingError
			break
		}
		keyInitMsg, err := keyShare.Init(keyMeta)
		encodedKeyInit, err := message.EncodeECDSAKeyInitMessage(keyInitMsg)
		if err != nil {
			log.Printf("error encoding ECDSA KeyInit message: %s", err)
			resp.Error = message.EncodingError
			break
		}
		resp.AddMessage(encodedKeyInit)
		log.Printf("Saving incomplete keyshare for keyid=%s", keyID)
		if err := client.SaveECDSAKey(keyID, keyShare, keyMeta); err != nil {
			log.Printf("Error with incomplete ECDSA keyshare saving process: %s", err)
			resp.Error = message.InternalError
			break
		}
		log.Printf("incomplete ECDSA Keyshare saved for keyid=%s", keyID)
	case message.ECDSAInitKeys:
		keyID := string(msg.Data[0])
		log.Printf("Client %s is sending us key init params for key %s", client.GetConnString(), keyID)
		keyInitMessages, err := message.DecodeECDSAKeyInitMessageList(msg.Data[1])
		if err != nil {
			log.Printf("error decoding ECDSA KeyInit message list: %s", err)
			resp.Error = message.DecodingError
			break
		}
		key, ok := client.ecdsa.keys[keyID]
		if !ok {
			log.Printf("error finding ECDSA key with id: %s", keyID)
			resp.Error = message.KeyNotFoundError
			break
		}
		err = key.Share.SetKey(key.Meta, keyInitMessages)
		if err != nil {
			log.Printf("error setting ECDSA Key: %s", err)
			resp.Error = message.InternalError
			break
		}
		key.Completed = true
		log.Printf("Saving complete keyshare for keyid=%s", keyID)
		if err := client.SaveECDSAKey(keyID, key.Share, key.Meta); err != nil {
			log.Printf("Error with complete ECDSA keyshare saving process: %s", err)
			resp.Error = message.InternalError
			break
		}
		log.Printf("complete ECDSA Keyshare saved for keyid=%s", keyID)
	case message.ECDSARound1:
		keyID := string(msg.Data[0])
		key, ok := client.ecdsa.keys[keyID]
		if !ok {
			log.Printf("error finding ECDSA key with id: %s. Keys available:", keyID)
			for k, _ := range client.ecdsa.keys {
				log.Printf("%s", k)
			}
		}
		client.ecdsa.currentKey = keyID
		h := msg.Data[1]
		log.Printf("Starting Round1 in signing document with key %s as asked by client %s", keyID, client.GetConnString())
		session, err := key.Share.NewSigSession(key.Meta, h)
		if err != nil {
			resp.Error = message.InternalError
			break
		}
		client.ecdsa.currentSession = session
		round1Msg, err := session.Round1()
		if err != nil {
			log.Printf("cannot execute round 1: %s", err)
			resp.Error = message.InternalError
			break
		}
		encoded, err := message.EncodeECDSARound1Message(round1Msg)
		if err != nil {
			log.Printf("cannot encode ECDSA round 1 message: %s", err)
			resp.Error = message.EncodingError
			break
		}
		resp.AddMessage(encoded)
	case message.ECDSARound2:
		keyID := client.ecdsa.currentKey
		if keyID == "" {
			log.Printf("Error: currentKey is not set %s", keyID)
			resp.Error = message.InternalError
			break
		}
		log.Printf("Starting Round2 in signing document with key %s as asked by client %s", keyID, client.GetConnString())
		round1Messages, err := message.DecodeECDSARound1MessageList(msg.Data[0])
		if err != nil {
			log.Printf("cannot decode ECDSA round 1 message list: %s", err)
			resp.Error = message.DecodingError
			break
		}
		round2Msg, err := client.ecdsa.currentSession.Round2(round1Messages)
		if err != nil {
			log.Printf("cannot execute round 2: %s", err)
			resp.Error = message.InternalError
			break
		}
		encoded, err := message.EncodeECDSARound2Message(round2Msg)
		if err != nil {
			log.Printf("cannot encode ECDSA round 2 message: %s", err)
			resp.Error = message.EncodingError
			break
		}
		resp.AddMessage(encoded)
	case message.ECDSARound3:
		keyID := client.ecdsa.currentKey
		if keyID == "" {
			log.Printf("Error: currentKey is not set %s", keyID)
			resp.Error = message.InternalError
			break
		}
		log.Printf("Starting Round2 in signing document with key %s as asked by client %s", keyID, client.GetConnString())
		round2Messages, err := message.DecodeECDSARound2MessageList(msg.Data[0])
		if err != nil {
			log.Printf("cannot decode ECDSA round 2 message List: %s", err)
			resp.Error = message.DecodingError
			break
		}
		round3Msg, err := client.ecdsa.currentSession.Round3(round2Messages)
		if err != nil {
			log.Printf("cannot execute round 3: %s", err)
			resp.Error = message.InternalError
			break
		}
		encoded, err := message.EncodeECDSARound3Message(round3Msg)
		if err != nil {
			log.Printf("cannot encode ECDSA round 3 message: %s", err)
			resp.Error = message.EncodingError
			break
		}
		resp.AddMessage(encoded)
	case message.ECDSAGetSignature:
		keyID := client.ecdsa.currentKey
		if keyID == "" {
			log.Printf("Error: currentKey is not set %s", keyID)
			resp.Error = message.InternalError
			break
		}
		log.Printf("Starting Round3 in signing document with key %s as asked by client %s", keyID, client.GetConnString())
		round3Messages, err := message.DecodeECDSARound3MessageList(msg.Data[0])
		if err != nil {
			log.Printf("cannot decode ECDSA round 3 message list: %s", err)
			resp.Error = message.DecodingError
			break
		}
		r, s, err := client.ecdsa.currentSession.GetSignature(round3Messages)
		if err != nil {
			log.Printf("error getting signature: %s", err)
			resp.Error = message.InternalError
			break
		}
		encoded, err := message.EncodeECDSASignature(r, s)
		if err != nil {
			log.Printf("cannot encode ECDSA signature: %s", err)
			resp.Error = message.EncodingError
			break
		}
		resp.AddMessage(encoded)
	case message.DeleteECDSAKeyShare:
		log.Printf("Client %s is asking us to delete a ECDSA KeyShare", client.GetConnString())
		keyID := string(msg.Data[0])
		log.Printf("Deleting keyshare for keyid=%s", keyID)
		if err := client.DeleteECDSAKey(keyID); err != nil {
			log.Printf("Error with key deleting: %s", err)
			resp.Error = message.InternalError
			break
		}
		log.Printf("Keyshare deleted for keyid=%s", keyID)
	}
	return resp
}

// SaveECDSAKey updates the key array of the server and asks the node to save the ecdsaKeys into the config file.
func (client *Client) SaveECDSAKey(id string, keyShare *tcecdsa.KeyShare, keyMeta *tcecdsa.KeyMeta) error {
	key, ok := client.ecdsa.keys[id]
	if !ok {
		key = &ecdsaKey{}
		client.ecdsa.keys[id] = key
	}
	key.ID = id
	key.Meta = keyMeta
	key.Share = keyShare
	return client.node.SaveConfigKeys()
}

// SaveECDSAKey deletes a key from the array of the server and asks the node to save the new key array into the config file.
func (client *Client) DeleteECDSAKey(id string) error {
	delete(client.ecdsa.keys, id)
	return client.node.SaveConfigKeys()
}

func parseECDSAKeys(conf []*config.ECDSAKeyConfig) (map[string]*ecdsaKey, error) {
	keys := make(map[string]*ecdsaKey)
	for _, key := range conf {
		var keyShare *tcecdsa.KeyShare
		var keyMeta *tcecdsa.KeyMeta
		if key.KeyShare != "" && key.KeyMetaInfo != "" {
			keyShareByte, err := base64.StdEncoding.DecodeString(key.KeyShare)
			if err != nil {
				return nil, err
			}
			keyShare, err = message.DecodeECDSAKeyShare(keyShareByte)
			if err != nil {
				return nil, err
			}
			keyMetaByte, err := base64.StdEncoding.DecodeString(key.KeyMetaInfo)
			if err != nil {
				return nil, err
			}
			keyMeta, err = message.DecodeECDSAKeyMeta(keyMetaByte)
			if err != nil {
				return nil, err
			}
		}
		keys[key.ID] = &ecdsaKey{
			ID:    key.ID,
			Meta:  keyMeta,
			Share: keyShare,
		}
	}
	return keys, nil
}

func saveECDSAKeys(keys map[string]*ecdsaKey) ([]*config.ECDSAKeyConfig, error) {
	keysConfig := make([]*config.ECDSAKeyConfig, 0)
	for _, key := range keys {
		keyShareBytes, err := message.EncodeECDSAKeyShare(key.Share)
		if err != nil {
			return nil, fmt.Errorf("error encoding ecdsaKeys: %s", err)
		}
		keyMetaBytes, err := message.EncodeECDSAKeyMeta(key.Meta)
		if err != nil {
			return nil, fmt.Errorf("error encoding ecdsaKeys: %s", err)
		}
		keyShareB64 := base64.StdEncoding.EncodeToString(keyShareBytes)
		keyMetaB64 := base64.StdEncoding.EncodeToString(keyMetaBytes)
		keysConfig = append(keysConfig, &config.ECDSAKeyConfig{
			ID:          key.ID,
			KeyMetaInfo: keyMetaB64,
			KeyShare:    keyShareB64,
		})
	}
	return keysConfig, nil
}
