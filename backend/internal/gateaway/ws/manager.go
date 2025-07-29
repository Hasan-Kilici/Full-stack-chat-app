package ws

import (
	"sync"
)

type room struct {
	clients map[*Client]bool
	sync.RWMutex
}

type manager struct {
	rooms map[string]*room
	sync.RWMutex
}

var Hub = &manager{
	rooms: make(map[string]*room),
}

func (m *manager) getOrCreateRoom(name string) *room {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.rooms[name]; !ok {
		m.rooms[name] = &room{clients: make(map[*Client]bool)}
	}
	return m.rooms[name]
}

func (m *manager) Join(c *Client) {
	r := m.getOrCreateRoom(c.Room)
	r.Lock()
	defer r.Unlock()
	r.clients[c] = true
}

func (m *manager) Leave(c *Client) {
	m.RLock()
	room, ok := m.rooms[c.Room]
	m.RUnlock()
	if !ok {
		return
	}
	room.Lock()
	defer room.Unlock()
	delete(room.clients, c)
	c.Conn.Close()
}

func (m *manager) Broadcast(c *Client, msg Message) {
	m.RLock()
	room, ok := m.rooms[c.Room]
	m.RUnlock()
	if !ok {
		return
	}
	room.RLock()
	defer room.RUnlock()
	for client := range room.clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			continue
		}
	}
}
