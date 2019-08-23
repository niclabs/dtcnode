package message

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"github.com/niclabs/tcrsa"
)

// GetRandomHexString returns a random hexadecimal string. It returns an error if it has any problem with the local PRNG.
func GetRandomHexString(len int) (string, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// EncodeKeyShare encodes a keyshare struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeKeyShare(share *tcrsa.KeyShare) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeKeyMeta encodes a keymeta struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeKeyMeta(meta *tcrsa.KeyMeta) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(meta); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeSigShare encodes a sigshare struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeSigShare(share *tcrsa.SigShare) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DecodeKeyShare decodes an array of bytes into a keyshare struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeKeyShare(byteShare []byte) (*tcrsa.KeyShare, error) {
	var keyShare tcrsa.KeyShare
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&keyShare); err != nil {
		return nil, err
	}
	return &keyShare, nil
}

// DecodeKeyMeta decodes an array of bytes into a keymeta struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeKeyMeta(byteShare []byte) (*tcrsa.KeyMeta, error) {
	var keyMeta tcrsa.KeyMeta
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&keyMeta); err != nil {
		return nil, err
	}
	return &keyMeta, nil
}

// DecodeSigShare decodes an array of bytes into a sigshare struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeSigShare(byteShare []byte) (*tcrsa.SigShare, error) {
	var sigShare tcrsa.SigShare
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&sigShare); err != nil {
		return nil, err
	}
	return &sigShare, nil
}
