package database

import (
	"errors"

	ti "github.com/thingsdb/go-thingsdb"
)

// globalClient represents the client that is globally accessible
var globalClient *Client

// Client represents a wrapper around a ThingsDB Conn
type Client struct {
	*ti.Conn
}

// InitializeClient initializes the globalClient
func InitializeClient(host string, port uint16, token string) (*Client, error) {
	conn := ti.NewConn(host, port, nil)

	if err := conn.Connect(); err != nil {
		return nil, err
	}

	if err := conn.AuthToken(token); err != nil {
		return nil, err
	}

	globalClient = &Client{Conn: conn}
	return globalClient, nil
}

// GetClient returns the global client when initialized
func GetClient() (*Client, error) {
	if globalClient == nil {
		return nil, errors.New("global client has not been initialized yet")
	}
	return globalClient, nil
}

// GetMigrationClient returns a new MigrationClient
func (c *Client) GetMigrationClient() *MigrationClient {
	return &MigrationClient{Client: c}
}
