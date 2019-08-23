package config

import "net"

// Config represents the main config of a node.
type Config struct {
	PublicKey  string        // Node public key
	PrivateKey string        // Node private key
	Host       string        // Node host
	Port       uint16        // Node port
	Client     *ClientConfig // List of servers
}

// ClientConfig represents a client configuration.
type ClientConfig struct {
	PublicKey string       // Client public key
	Host      string       // Client hostname or IP
	Keys      []*KeyConfig // List of keys
}

// KeyConfig represents a key share on the node.
type KeyConfig struct {
	ID          string // Key UUID
	KeyShare    string // Keyshare
	KeyMetaInfo string // Key Metainformation
}

// Returns a client, given its ID
func (config *Config) GetClientByID(id string) *ClientConfig {
	if config.Client.PublicKey == id {
		return config.Client
	}
	return nil
}

// Returns a list of IPs. If a server has a hostname instead of an IP, it resolves it.
func (config *Config) GetClientIPs() ([]string, error) {
	ips := make([]string, 1)
	// try to parse as IP
	ip, err := net.ResolveIPAddr("ip", config.Client.Host)
	if err != nil {
		return nil, err
	}
	ips[0] = ip.String()
	return ips, nil
}

// Returns the list of public keys of the servers
func (config *Config) GetClientPubKeys() []string {
	pubkeys := make([]string, 1)
	pubkeys[0] = config.Client.PublicKey
	return pubkeys
}

// Returns a key in a server, based on its ID
func (serverConfig *ClientConfig) GetKeyByID(id string) *KeyConfig {
	for _, key := range serverConfig.Keys {
		if key.ID == id {
			return key
		}
	}
	return nil
}
