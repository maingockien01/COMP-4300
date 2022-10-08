package protocols

import "WeebChat/pkg/models"

type ProtocolUser struct {
	Metadata ProtocolMetadata
	Data     models.User
}
