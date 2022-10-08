package protocols

import "WeebChat/pkg/models"

type ProtocolMessage struct {
	Metadata ProtocolMetadata
	Messages []models.Message
}
