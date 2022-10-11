package protocols

type ProtocolMetadata struct {
	Version   int    `json:"version"`
	From      string `json:"from"`
	Direction string `json:"direction"`
	Type      string `json:"type"`
}

type Protocol struct {
	Metadata ProtocolMetadata `json:"metadata"`
}

const (
	TYPE_MESSAGE = "message"
	TYPE_ROOM    = "room"
	TYPE_USER    = "user"
)

const (
	DIRECTION_PULL     = "pull"
	DIRECTION_PUSH     = "push"
	DIRECTION_JOIN     = "join"
	DIRECTION_GREETING = "greeting"
	DIRECTION_UPDATE   = "update"
)

const (
	V1 = 1
)
