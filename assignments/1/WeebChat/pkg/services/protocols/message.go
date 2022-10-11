package protocols

import "WeebChat/pkg/models"

type ProtocolMessage struct {
	Metadata ProtocolMetadata  `json:"metadata"`
	Data     []*models.Message `json:"data"`
}
