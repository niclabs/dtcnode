package config

import "net"

// Config represents the main config of a node.
type Config struct {
	PublicKey  string        // Node public key
	PrivateKey string        // Node private key
	Host       string        // Node host
	Port       uint16        // Node port
	Server     *ServerConfig // List of servers
}

// ServerConfig represents a server configuration.
type ServerConfig struct {
	PublicKey string       // Server public key
	Host      string       // Server hostname or IP
	Keys      []*KeyConfig // List of keys
}

// KeyConfig represents a key on a server.
type KeyConfig struct {
	ID          string // Key UUID
	KeyShare    string // Keyshare
	KeyMetaInfo string // Key Metainformation
}

// Returns a server, given its ID
func (config *Config) GetServerByID(id string) *ServerConfig {
	if config.Server.PublicKey == id {
		return config.Server
	}
	return nil
}

// Returns a list of IPs. If a server has a hostname instead of an IP, it resolves it.
func (config *Config) GetServerIPs() ([]string, error) {
	ips := make([]string, 1)
	// try to parse as IP
	ip, err := net.ResolveIPAddr("ip", config.Server.Host)
	if err != nil {
		return nil, err
	}
	ips[0] = ip.String()
	return ips, nil
}

// Returns the list of public keys of the servers
func (config *Config) GetServerPubKeys() []string {
	pubkeys := make([]string, 1)
	pubkeys[0] = config.Server.PublicKey
	return pubkeys
}

// Returns a key in a server, based on its ID
func (serverConfig *ServerConfig) GetKeyByID(id string) *KeyConfig {
	for _, key := range serverConfig.Keys {
		if key.ID == id {
			return key
		}
	}
	return nil
}
