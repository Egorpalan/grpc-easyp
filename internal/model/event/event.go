package event

type MessageType string

const (
	MessageTypeNoteCreated MessageType = "note_created"
)

type Event struct {
	Type   MessageType
	NoteID string
}
