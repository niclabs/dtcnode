package message

import (
	"bytes"
	"encoding/gob"
	"github.com/niclabs/tcecdsa"
	"math/big"
)

type Signature struct {
	R, S *big.Int
}

// EncodeECDSAKeyShare encodes a keyshare struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSAKeyShare(share *tcecdsa.KeyShare) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeECDSAKeyMeta encodes a keymeta struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSAKeyMeta(meta *tcecdsa.KeyMeta) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(meta); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}


// EncodeECDSAKeyInitMessage encodes a KeyInitMessage struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSAKeyInitMessage(share *tcecdsa.KeyInitMessage) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeECDSAKeyInitMessageList encodes a KeyInitMessageList struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSAKeyInitMessageList(share tcecdsa.KeyInitMessageList) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeECDSARound1Message encodes a Round1Message struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSARound1Message(share *tcecdsa.Round1Message) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeECDSARound1MessageList encodes a Round1MessageList struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSARound1MessageList(share tcecdsa.Round1MessageList) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeECDSARound2Message encodes a Round2Message struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSARound2Message(share *tcecdsa.Round2Message) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeECDSARound2MessageList encodes a Round2MessageList struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSARound2MessageList(share tcecdsa.Round2MessageList) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeECDSARound3Message encodes a Round3Message struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSARound3Message(share *tcecdsa.Round3Message) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeECDSARound3MessageList encodes a Round3MessageList struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSARound3MessageList(share tcecdsa.Round3MessageList) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(share); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeECDSASignature encodes a Signature struct into an array of bytes, using the golang gob encoder. It returns an error if it cannot encode the struct.
func EncodeECDSASignature(r, s *big.Int) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(&Signature{r, s}); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DecodeECDSAKeyShare decodes an array of bytes into a keyshare struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSAKeyShare(byteShare []byte) (*tcecdsa.KeyShare, error) {
	var keyShare *tcecdsa.KeyShare
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&keyShare); err != nil {
		return nil, err
	}
	return keyShare, nil
}


// DecodeECDSAKeyMeta decodes an array of bytes into a keymeta struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSAKeyMeta(byteShare []byte) (*tcecdsa.KeyMeta, error) {
	var keyMeta tcecdsa.KeyMeta
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&keyMeta); err != nil {
		return nil, err
	}
	return &keyMeta, nil
}

// DecodeECDSAKeyInitMessage decodes an array of bytes into a KeyInitMessage struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSAKeyInitMessage(byteShare []byte) (*tcecdsa.KeyInitMessage, error) {
	var keyInitMsg tcecdsa.KeyInitMessage
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&keyInitMsg); err != nil {
		return nil, err
	}
	return &keyInitMsg, nil
}

// DecodeECDSAKeyInitMessageList decodes an array of bytes into a KeyInitMessageList struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSAKeyInitMessageList(byteShare []byte) (tcecdsa.KeyInitMessageList, error) {
	var keyInitMsg tcecdsa.KeyInitMessageList
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&keyInitMsg); err != nil {
		return nil, err
	}
	return keyInitMsg, nil
}

// DecodeECDSARound1Message decodes an array of bytes into a Round1Message struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSARound1Message(byteShare []byte) (*tcecdsa.Round1Message, error) {
	var round1Msg tcecdsa.Round1Message
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&round1Msg); err != nil {
		return nil, err
	}
	return &round1Msg, nil
}

// DecodeECDSARound1MessageList decodes an array of bytes into a Round1MessageList struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSARound1MessageList(byteShare []byte) (tcecdsa.Round1MessageList, error) {
	var round1Msg tcecdsa.Round1MessageList
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&round1Msg); err != nil {
		return nil, err
	}
	return round1Msg, nil
}

// DecodeECDSARound2Message decodes an array of bytes into a Round2Message struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSARound2Message(byteShare []byte) (*tcecdsa.Round2Message, error) {
	var round2Msg tcecdsa.Round2Message
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&round2Msg); err != nil {
		return nil, err
	}
	return &round2Msg, nil
}

// DecodeECDSARound2Message decodes an array of bytes into a Round2Message struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSARound2MessageList(byteShare []byte) (tcecdsa.Round2MessageList, error) {
	var round2Msg tcecdsa.Round2MessageList
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&round2Msg); err != nil {
		return nil, err
	}
	return round2Msg, nil
}

// DecodeECDSARound3Message decodes an array of bytes into a Round3Message struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSARound3Message(byteShare []byte) (*tcecdsa.Round3Message, error) {
	var round3Msg tcecdsa.Round3Message
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&round3Msg); err != nil {
		return nil, err
	}
	return &round3Msg, nil
}

// DecodeECDSARound3MessageList decodes an array of bytes into a Round3MessageList struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSARound3MessageList(byteShare []byte) (tcecdsa.Round3MessageList, error) {
	var round3Msg tcecdsa.Round3MessageList
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&round3Msg); err != nil {
		return nil, err
	}
	return round3Msg, nil
}

// DecodeECDSASignature decodes an array of bytes into a signature struct, using the golang gob decode. It returns an error if it cannot decode the struct.
func DecodeECDSASignature(byteShare []byte) (*big.Int, *big.Int, error) {
	var sig *Signature
	buffer := bytes.NewBuffer(byteShare)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&sig); err != nil {
		return nil, nil, err
	}
	return sig.R, sig.S, nil
}

