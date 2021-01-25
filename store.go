package lightmq

import (
	"sync"
)

// ClientStore ...
type ClientStore struct {
	c  map[string]Client
	mu sync.RWMutex
}

// NewClientStore creates store of clients
func NewClientStore() *ClientStore {
	return &ClientStore{
		c:  make(map[string]Client),
		mu: sync.RWMutex{},
	}
}

// Add adds new client
func (l *ClientStore) Add(c Client) {
	l.mu.Lock()
	l.c[c.ClientID] = c
	l.mu.Unlock()
}
