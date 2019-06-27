package message

// Type enumerates the message types.
type Type byte

const (
	None Type = iota
	SendKeyShare
	AskForSigShare
	DeleteKeyShare
)

// TypeToString transforms a message type into a string. Useful for debugging.
var TypeToString = map[Type]string{
	None:           "Undefined type",
	SendKeyShare:   "Send Key Share",
	AskForSigShare: "Ask for Signature Share",
	DeleteKeyShare: "Delete Key Share",
}

func (mType Type) String() string {
	if name, ok := TypeToString[mType]; ok {
		return name
	} else {
		return "Unknown Message"
	}
}
