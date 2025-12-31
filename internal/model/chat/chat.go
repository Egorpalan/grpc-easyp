package chat

import "time"

type Message struct {
	CorrelationID string
	Text          string
	Timestamp     time.Time
}

type Error struct {
	Code    int32
	Message string
	Details []interface{}
}
