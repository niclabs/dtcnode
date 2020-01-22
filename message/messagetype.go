package message

// Type enumerates the message types.
type Type byte

const (
	None Type = iota
	SendRSAKeyShare
	GetRSASigShare
	DeleteRSAKeyShare
	SendECDSAKeyShare
	ECDSAInitKeys
	ECDSARound1
	ECDSARound2
	ECDSARound3
	ECDSAGetSignature
	DeleteECDSAKeyShare
)

// TypeToString transforms a message type into a string. Useful for debugging.
var TypeToString = map[Type]string{
	None:                "Undefined type",
	SendRSAKeyShare:     "RSA Send Key Share",
	GetRSASigShare:      "RSA Ask for Signature Share",
	DeleteRSAKeyShare:   "RSA Delete Key Share",
	SendECDSAKeyShare:   "ECDSA Send Key Share",
	ECDSAInitKeys:       "ECDSA Initialize RSAKeys",
	ECDSARound1:         "ECDSA Round 1",
	ECDSARound2:         "ECDSA Round 2",
	ECDSARound3:         "ECDSA Round 3",
	ECDSAGetSignature:   "ECDSA Get Signature",
	DeleteECDSAKeyShare: "ECDSA Delete Key Share",
}

var TypeToClientDataLength = map[Type]int{
	None:                0,
	SendRSAKeyShare:     3, // keyID, keyShare, keyMeta -> {}
	GetRSASigShare:      2, // keyID, hash -> sigShare
	DeleteRSAKeyShare:   1, // keyID -> {}
	SendECDSAKeyShare:   3, // keyID, keyShare, keyMeta -> InitKeyMessage
	ECDSAInitKeys:       2, // keyID, InitKeyMessageList -> {}
	ECDSARound1:         2, // keyID, hash -> Round1Message
	ECDSARound2:         1, // Round1MessageList -> Round2Message
	ECDSARound3:         1, // Round2MessageList -> Round3Message
	ECDSAGetSignature:   1, // Round3MessageList -> r, s
	DeleteECDSAKeyShare: 1, // keyID -> {}
}

var TypeToNodeDataLength = map[Type]int{
	None:                0,
	SendRSAKeyShare:     0, // keyID, keyShare, keyMeta -> {}
	GetRSASigShare:      1, // keyID, hash -> sigShare
	DeleteRSAKeyShare:   0, // keyID -> {}
	SendECDSAKeyShare:   1, // keyID, keyShare, keyMeta -> InitKeyMessage
	ECDSAInitKeys:       0, // keyID, InitKeyMessageList -> {}
	ECDSARound1:         1, // keyID, hash -> Round1Message
	ECDSARound2:         1, // Round1MessageList -> Round2Message
	ECDSARound3:         1, // Round2MessageList -> Round3Message
	ECDSAGetSignature:   1, // Round3MessageList -> (r, s)
	DeleteECDSAKeyShare: 0, // keyID -> {}
}

func (mType Type) String() string {
	if name, ok := TypeToString[mType]; ok {
		return name
	} else {
		return "Unknown Message"
	}
}

// Returns true if the message is of type RSA, and false if it is not.
func (mType Type) IsRSA() bool {
	return mType >= SendRSAKeyShare && mType <= DeleteRSAKeyShare
}

// IsECDSA returns true if the message is of type ECDSA, and false if it is not.
func (mType Type) IsECDSA() bool {
	return mType >= SendECDSAKeyShare && mType <= DeleteECDSAKeyShare
}

func (mType Type) ClientDataLength() int {
	if length, ok := TypeToClientDataLength[mType]; ok {
		return length
	}
	return 0
}

func (mType Type) NodeDataLength() int {
	if length, ok := TypeToNodeDataLength[mType]; ok {
		return length
	}
	return 0
}
