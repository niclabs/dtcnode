package main

type Config struct {
	PublicKey  string
	PrivateKey string
	Host       string
	Port       uint16
	Servers    []*ServerConfig
}

func (config *Config) GetServerByID(id string) *ServerConfig {
	for _, serverConfig := range config.Servers {
		if serverConfig.PublicKey == id {
			return serverConfig
		}
	}
	return nil
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
	PublicKey string
	IP        string
	Port      uint16
	Keys      []*KeyConfig
}

func (serverConfig *ServerConfig) GetKeyByID(id string) *KeyConfig {
	for _, key := range serverConfig.Keys {
		if key.ID == id {
			return key
		}
	}
	return nil
}

type KeyConfig struct {
	ID          string
	KeyShare    string
	KeyMetaInfo string
}


