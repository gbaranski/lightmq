package client

import (
	"crypto/ed25519"
	"fmt"
)

// Config ...
type Config struct {
	// ClientID
	//
	// Required
	ClientID string

	// Ed25519 Public key
	//
	// Required
	PublicKey ed25519.PublicKey

	// Ed25519 Private key
	//
	// Required
	PrivateKey ed25519.PrivateKey

	// Hostname of where broker should listen
	//
	// Default: "localhost"
	Hostname string

	// Port of where broker should listen
	//
	// Default: "997"
	Port uint32
}

// Parse parses options and set defaults
func (cfg Config) Parse() (Config, error) {
	if cfg.Hostname == "" {
		cfg.Hostname = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 997
	}
	if cfg.ClientID == "" {
		return cfg, fmt.Errorf("ClientID cannot be empty")
	}

	return cfg, nil
}
