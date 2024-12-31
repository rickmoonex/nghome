package database

import (
	ti "github.com/thingsdb/go-thingsdb"
)

// Client represents a wrapper around a ThingsDB Conn
type Client struct {
	*ti.Conn
}

// GetMigrationClient returns a new MigrationClient
func (c *Client) GetMigrationClient() *MigrationClient {
	return &MigrationClient{Client: c}
}
