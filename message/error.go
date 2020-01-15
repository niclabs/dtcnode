package message

import "fmt"

// NodeError represents an error code from the nodes.
type NodeError uint8

const (
	Ok NodeError = iota
	// Invalid message
	InvalidMessageError
	// Network Errors
	ReceiveMessageError
	ParseMessageError
	SendResponseError
	// Encryption/Decryption errors
	DecodingError
	EncodingError
	// Signing Errors
	KeyNotFoundError
	DocSignError
	// Internal Errors (I/O)
	InternalError
	// Invalid error number (keep at the end)
	UnknownError = NodeError(1<<8 - 1)
)

// ErrorToString maps the error codes to string message. Useful for debugging.
var ErrorToString = map[NodeError]string{
	Ok:                  "not an error",
	InvalidMessageError: "invalid message",
	ReceiveMessageError: "cannot receive message",
	ParseMessageError:   "cannot parse received message",
	SendResponseError:   "cannot send response",
	EncodingError:       "cannot encode a struct to a message",
	DecodingError:       "cannot decode received struct",
	KeyNotFoundError:    "key not found in the node",
	DocSignError:        "cannot sign the document",
	InternalError:       "internal input/output error",
	UnknownError:        "unknown error",
}

func (err NodeError) Error() string {
	if int(err) >= len(ErrorToString) {
		return ErrorToString[UnknownError]
	}
	return ErrorToString[err]
}

// Composes a nodeError with another error thrown by a routine the node uses.
func (err NodeError) ComposeError(err2 error) string {
	return fmt.Sprintf("%s: %s", err.Error(), err2.Error())
}
