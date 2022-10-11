package protocols

import "WeebChat/pkg/models"

type ProtocolUser struct {
	Metadata ProtocolMetadata `json:"metadata"`
	Data     models.User      `json:"data"`
}
