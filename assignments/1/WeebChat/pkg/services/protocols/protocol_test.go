package protocols

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	text := "{\"metadata\": {\"version\": 1, \"from\": \"user1\", \"direction\": \"greeting\", \"type\": \"user\"}, \"data\": {\"id\": \"user1\", \"name\": \"Kevin\"}}"
	var protocolMetadata Protocol

	err := json.Unmarshal([]byte(text), &protocolMetadata)

	if err != nil {
		t.Fatal("Failed to unmarhsal")
	}

	if protocolMetadata.Metadata.Direction != "greeting" || protocolMetadata.Metadata.From != "user1" || protocolMetadata.Metadata.Type != "user" {
		fmt.Println(protocolMetadata)
		t.Fatal("Fail to unmarshal type", protocolMetadata)
	}
}
