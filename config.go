package main


type Config struct {
	PublicKey  string
	PrivateKey string
	IP         string
	Port       uint16
	Servers    []ServerConfig
}

func (config *Config) GetServerIPs() []string {
	ips := make([]string, len(config.Servers))
	for i, server := range config.Servers {
		ips[i] = server.IP
	}
	return ips
}

func (config *Config) GetServerPubKeys() []string {
	pubkeys := make([]string, len(config.Servers))
	for i, server := range config.Servers {
		pubkeys[i] = server.PublicKey
	}
	return pubkeys
}


type ServerConfig struct {
	PublicKey   string
	IP          string
	Port        uint16
	KeyShare    string
	KeyMetaInfo string
}
