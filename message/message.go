package message

import (
	"fmt"
)

// Message represents a generic message which is sent between server and nodes.
type Message struct {
	From       string    // Identification for the sender node.
	ResponseOf string    // Identification for the original "from" field if the message is a response.
	ID         string    // Random hex ID for the message. Useful to do follow ups
	Type       Type      // Type of the message.
	Error      NodeError // An error code. It is 0 if the message is ok.
	Data       [][]byte  // A list of byte arrays with the binary data of the message.
}

// FromBytes transforms a raw array of array of bytes into a message, or returns an error if it can't transform the message.
func FromBytes(rawMsg [][]byte) (*Message, error) {
	if len(rawMsg) < 5 { // header is dealer ID, rest is message struct.
		return nil, fmt.Errorf("bad byte array length: %d instead of 5", len(rawMsg))
	}
	return &Message{
		From:  string(rawMsg[0]),
		ResponseOf: string(rawMsg[1]),
		ID:    string(rawMsg[2]),
		Type:  Type(rawMsg[3][0]),
		Error: NodeError(rawMsg[4][0]),
		Data:  rawMsg[5:],
	}, nil
}

// NewMessage creates a new message using the arguments provided, or returns an error if it cannot create the message object
//(related currently to a problem in the generation of message IDs)
func NewMessage(rType Type, from string, msgs ...[]byte) (*Message, error) {
	id, err := GetRandomHexString(6)
	if err != nil {
		return nil, err
	}
	req := &Message{
		From: from,
		ID:   id,
		Type: rType,
		Data: make([][]byte, 0),
	}
	req.Data = append(req.Data, msgs...)
	return req, nil
}

// GetBytesLists transforms a message into an array of arrays of bytes, useful to send the message to the other end.
func (message *Message) GetBytesLists() []interface{} {
	b := []interface{}{
		[]byte(message.From),
		[]byte(message.ResponseOf),
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

// NewResponse creates a new message with some fields copied from another message. This method is useful to create replies quickly. It receives a default status code as argument and the new Node ID.
func (message *Message) NewResponse(ourID string, status NodeError) *Message {
	return &Message{
		From:       ourID,
		ResponseOf: message.From,
		ID:         message.ID,
		Type:       message.Type,
		Error:      status,
		Data:       make([][]byte, 0),
	}
}

// ValidClientDataLength returns true if the number of data fields is equal to the expected in a message sent by the client.
func (message *Message) ValidClientDataLength() bool {
	return len(message.Data) == message.Type.ClientDataLength()
}

// ValidNodeDataLength returns true if the number of data fields is equal to the expected in a message sent by the node.
func (message *Message) ValidNodeDataLength() bool {
	return len(message.Data) == message.Type.NodeDataLength()
}

// ResponseOK returns true if a response matches with the message that generated it. It also checks for the length of data fields on the message. If it is less than the minDataLen argument, it returns an error.
func (message *Message) ResponseOK(message2 *Message) error {
	if message.ID != message2.ID {
		return fmt.Errorf("ID mismatch: got: %s, expected: %s", message.ID, message2.ID)
	}
	// Note: From should not be the same as the response
	if message.Type != message2.Type {
		return fmt.Errorf("type mismatch: got: %s, expected: %s", message.Type, message2.Type)
	}
	if message.Error != Ok {
		return fmt.Errorf("response has error: %s", message.Error.Error())
	}
	if !message.ValidNodeDataLength() {
		return fmt.Errorf("data length mismatch: got: %d, expected: %d", len(message.Data), message.Type.NodeDataLength())
	}
	return nil
}
