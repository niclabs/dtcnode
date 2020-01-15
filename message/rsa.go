package message

import (
	"bytes"
	"encoding/gob"
	"github.com/niclabs/tcrsa"
)

// EncodeRSAKeyShare encodes a keyshare struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeRSAKeyShare(share *tcrsa.KeyShare) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeRSAKeyMeta encodes a keymeta struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeRSAKeyMeta(meta *tcrsa.KeyMeta) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(meta); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeRSASigShare encodes a sigshare struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeRSASigShare(share *tcrsa.SigShare) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DecodeRSAKeyShare decodes an array of bytes into a keyshare struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeRSAKeyShare(byteShare []byte) (*tcrsa.KeyShare, error) {
	var keyShare tcrsa.KeyShare
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&keyShare); err != nil {
		return nil, err
	}
	return &keyShare, nil
}

// DecodeRSAKeyMeta decodes an array of bytes into a keymeta struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeRSAKeyMeta(byteShare []byte) (*tcrsa.KeyMeta, error) {
	var keyMeta tcrsa.KeyMeta
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&keyMeta); err != nil {
		return nil, err
	}
	return &keyMeta, nil
}

// DecodeRSASigShare decodes an array of bytes into a sigshare struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeRSASigShare(byteShare []byte) (*tcrsa.SigShare, error) {
	var sigShare tcrsa.SigShare
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&sigShare); err != nil {
		return nil, err
	}
	return &sigShare, nil
}
