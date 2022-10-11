package protocols

import "WeebChat/pkg/models"

type ProtocolRoom struct {
	Metadata ProtocolMetadata `json:"metadata"`
	Data     []models.Room    `json:"data"`
}
