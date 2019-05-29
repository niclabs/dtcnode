package main

import "fmt"

type NodeError uint8

const (
	// c'est ne pas un err
	Ok NodeError = iota
	// Network Errors
	ReceiveMessageError
	ParseMessageError
	SendResponseError
	// Signature Reception Errors
	KeyShareDecodeError
	KeyMetaDecodeError
	// Signing Errors
	NotInitializedError
	DocDecodeError
	DocSignError
	SigShareEncodeError
	// Invalid error number (keep at the end)
	UnknownError = NodeError(1<<8 - 1)
)

var ErrorToString = map[NodeError]string{
	Ok:                  "not an error",
	ReceiveMessageError: "cannot receive message",
	ParseMessageError:   "cannot parse received message",
	SendResponseError:   "cannot send response",
	KeyShareDecodeError: "cannot decode received keyshare",
	KeyMetaDecodeError:  "cannot decode received keymeta",
	NotInitializedError: "node not initialized with the server",
	DocDecodeError:      "cannot decode the document to sign",
	SigShareEncodeError: "cannot encode the signature to a message",
	DocSignError:        "cannot sign the document",
	UnknownError:        "unknown error",
}

func (err NodeError) Error() string {
	if int(err) >= len(ErrorToString) {
		return ErrorToString[UnknownError]
	}
	return ErrorToString[err]
}

func (err NodeError) ComposeError(err2 error) string {
	return fmt.Sprintf("%s: %s", err.Error(), err2.Error())
}
