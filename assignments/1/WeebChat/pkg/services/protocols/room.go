package protocols

import "WeebChat/pkg/models"

type ProtocolRoom struct {
	Metadata ProtocolMetadata
	Data     []models.Room
}
