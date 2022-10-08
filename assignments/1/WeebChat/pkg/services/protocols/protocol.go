package protocols

type ProtocolMetadata struct {
	Version   int
	From      string
	Direction string
	Type      string
}

type Protocol struct {
	Metadata ProtocolMetadata
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
)

const (
	V1 = 1
)
