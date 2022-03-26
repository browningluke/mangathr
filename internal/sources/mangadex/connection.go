package mangadex

import "fmt"

type Connection struct {
}

func NewConnection() *Connection {
	fmt.Println("Created a mangadex connection")
	return &Connection{}
}

func (c *Connection) GetConnectionName() string {
	return "Mangadex"
}
