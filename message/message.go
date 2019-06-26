package message

import (
	"fmt"
)

// Message represents a generic message which is sent between server and nodes.
type Message struct {
	NodeID string    // Identification for the sender node. It usually is the public key of the node.
	ID     string    // Random hex ID for the message. Useful to do follow ups
	Type   Type      // Type of the message.
	Error  NodeError // An error code. It is 0 if the message is ok.
	Data   [][]byte  // A list of byte arrays with the binary data of the message.
}

// FromBytes transforms a raw array of array of bytes into a message, or returns an error if it can't transform the message.
func FromBytes(rawMsg [][]byte) (*Message, error) {
	if len(rawMsg) < 4 { // header is dealer ID, rest is message struct.
		return nil, fmt.Errorf("bad byte array length")
	}
	return &Message{
		NodeID: string(rawMsg[0]), // Provided by
		ID:     string(rawMsg[1]),
		Type:   Type(rawMsg[2][0]),
		Error:  NodeError(rawMsg[3][0]),
		Data:   rawMsg[4:],
	}, nil
}

// NewMessage creates a new message using the arguments provided, or returns an error if it cannot create the message object (related currently to a problem in the generation of message IDs)
func NewMessage(rType Type, nodeID string, msgs ...[]byte) (*Message, error) {
	id, err := GetRandomHexString(6)
	if err != nil {
		return nil, err
	}
	req := &Message{
		NodeID: nodeID,
		ID:     id,
		Type:   rType,
		Data:   make([][]byte, 0),
	}
	req.Data = append(req.Data, msgs...)
	return req, nil
}

// GetBytesLists transforms a message into an array of arrays of bytes, useful to send the message to the other end.
func (message *Message) GetBytesLists() []interface{} {
	b := []interface{}{
		[]byte(message.ID),
		[]byte{byte(message.Type)},
		[]byte{byte(message.Error)},
	}
	for _, datum := range message.Data {
		b = append(b, datum)
	}
	return b
}

// AddMessage appends a data field to the message.
func (message *Message) AddMessage(data []byte) {
	message.Data = append(message.Data, data)
}

// CopyWithoutData creates a new message with some fields copied from another message. This method is useful to create replies quickly. It receives a default status code as argument. This status code is used in the new message.
func (message *Message) CopyWithoutData(status NodeError) *Message {
	return &Message{
		ID:     message.ID,
		NodeID: message.NodeID,
		Type:   message.Type,
		Error:  status,
		Data:   make([][]byte, 0),
	}
}

// Ok returns true if a response matches with the message that generated it. It also checks for the length of data fields on the message. If it is less than the minDataLen argument, it returns an error.
func (message *Message) Ok(message2 *Message, minDataLen int) error {
	if message.ID != message2.ID {
		return fmt.Errorf("ID mismatch: got: %s, expected: %s", message.ID, message2.ID)
	}
	if message.NodeID != message2.NodeID {
		return fmt.Errorf("node ID mismatch: got: %s, expected: %s", message.NodeID, message2.NodeID)
	}
	if message.Type != message2.Type {
		return fmt.Errorf("type mismatch: got: %s, expected: %s", message.Type, message2.Type)
	}
	if message.Error != Ok {
		return fmt.Errorf("response has error: %s", message.Error.Error())
	}
	if len(message.Data) < minDataLen {
		return fmt.Errorf("data length mismatch: got: %d, expected at least: %d", len(message.Data), minDataLen)
	}
	return nil
}
