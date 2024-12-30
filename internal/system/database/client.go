package database

import (
	"errors"

	ti "github.com/thingsdb/go-thingsdb"
)

// Client represents a wrapper around a ThingsDB Conn
type Client struct {
	*ti.Conn
}

// ensureCollection takes a collection name and makes sure it exists
func (c *Client) ensureCollection(name string) error {
	vars := map[string]interface{}{
		"name": name,
	}

	res, err := c.Query("@thingsdb", "has_collection(name);", vars)
	if err != nil {
		return err
	}

	exists, ok := res.(bool)
	if !ok {
		return errors.New("could not cast `has_collection` results into bool")
	}

	if exists {
		return nil
	}

	_, err = c.Query("@thingsdb", "new_collection(name);", vars)
	return err
}
