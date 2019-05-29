package main


type Config struct {
	PublicKey string
	PrivateKey string
	IP string
	RouterPort uint16
	SubPort uint16
	Servers []Server
}


type Server struct {
	PublicKey string
	IP string
	KeyShare string
	KeyMetaInfo string
}