package lightmq

import "crypto/ed25519"

// Config ...
type Config struct {
	// Hostname of where broker should listen
	//
	// Default: "0.0.0.0"
	Hostname string

	// Port of where broker should listen
	//
	// Default: "1883"
	Port uint32

	// Ed25519 Private key
	//
	// Required
	PrivateKey ed25519.PrivateKey
}

// Parse parses options and set defaults
func (cfg Config) Parse() Config {
	if cfg.Hostname == "" {
		cfg.Hostname = "0.0.0.0"
	}
	if cfg.Port == 0 {
		cfg.Port = 1883
	}

	return cfg
}
