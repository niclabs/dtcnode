package main

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"os"
)

func main() {
	pk, sk, err := zmq4.NewCurveKeypair()
	if err != nil {
		fmt.Fprintf(os.Stderr, "CANNOT GENERATE CURVE KEY PAIR!\n")
	}
	fmt.Printf("PUBLIC KEY: %s\n", pk)
	fmt.Printf("SECRET KEY: %s\n", sk)
}
