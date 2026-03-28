package websocket

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHub_Lifecycle(t *testing.T) {
	hub := InitHub()
	
	// Create a dummy client
	client := &Client{
		hub:  hub,
		send: make(chan []byte, 10),
	}

	// Test Register
	hub.Register(client)
	time.Sleep(10 * time.Millisecond) // Give time for the run loop to handle registration

	hub.mu.Lock()
	_, ok := hub.clients[client]
	hub.mu.Unlock()
	assert.True(t, ok)

	// Test BroadcastEvent
	payload := map[string]string{"id": "123"}
	hub.BroadcastEvent(EventOrderCreated, payload)

	select {
	case msg := <-client.send:
		var event Event
		err := json.Unmarshal(msg, &event)
		assert.NoError(t, err)
		assert.Equal(t, EventOrderCreated, event.Type)
		
		payloadMap := event.Payload.(map[string]interface{})
		assert.Equal(t, "123", payloadMap["id"])
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timeout waiting for broadcast")
	}

	// Test Unregister
	hub.Unregister(client)
	time.Sleep(10 * time.Millisecond)

	hub.mu.Lock()
	_, ok = hub.clients[client]
	hub.mu.Unlock()
	assert.False(t, ok)
}

func TestHub_BroadcastFull(t *testing.T) {
	hub := InitHub()
	
	// Create a client with NO buffer space to test blocking/dropping
	client := &Client{
		hub:  hub,
		send: make(chan []byte), // Unbuffered channel
	}

	hub.Register(client)
	time.Sleep(10 * time.Millisecond)

	// Since the client isn't reading, the broadcast should eventually trigger the default case
	// and unregister the client.
	hub.BroadcastEvent("TEST_DROP", nil)
	time.Sleep(50 * time.Millisecond)

	hub.mu.Lock()
	_, ok := hub.clients[client]
	hub.mu.Unlock()
	assert.False(t, ok, "Client should have been unregistered due to blocked send")
}

func TestHub_BroadcastError(t *testing.T) {
	hub := InitHub()
	// Channel cannot be marshaled to JSON
	hub.BroadcastEvent("ERROR", make(chan int))
}

func TestNewClient(t *testing.T) {
	hub := InitHub()
	client := NewClient(hub, nil)
	assert.NotNil(t, client)
	assert.Equal(t, hub, client.hub)
}
