package main

type Config struct {
	PublicKey  string
	PrivateKey string
	Host       string
	Port       uint16
	Server     *ServerConfig
}

func (config *Config) GetServerByID(id string) *ServerConfig {
	if config.Server.PublicKey == id {
		return config.Server
	}
	return nil
}

func (config *Config) GetServerIPs() []string {
	ips := make([]string, 1)
	ips[0] = config.Server.Host
	return ips
}

func (config *Config) GetServerPubKeys() []string {
	pubkeys := make([]string, 1)
	pubkeys[0] = config.Server.PublicKey
	return pubkeys
}

type ServerConfig struct {
	PublicKey string
	Host      string
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
