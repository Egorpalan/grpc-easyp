package notes

import (
	"context"
	"log/slog"

	"github.com/Egorpalan/grpc-easyp/internal/model/event"
	"github.com/Egorpalan/grpc-easyp/internal/model/exception"
	"github.com/Egorpalan/grpc-easyp/internal/model/notes"
	"github.com/Egorpalan/grpc-easyp/internal/repository/postgresql"
	"github.com/google/uuid"
)

type Service interface {
	CreateOrUpdateNote(ctx context.Context, id *string, title, description string) (*notes.Note, error)
	GetNote(ctx context.Context, id string) (*notes.Note, error)
	ListNotes(ctx context.Context) ([]*notes.Note, error)
	DeleteNote(ctx context.Context, id string) error
}

type service struct {
	logger         *slog.Logger
	repo           postgresql.Repository
	eventPublisher EventPublisher
}

func NewService(logger *slog.Logger, repo postgresql.Repository, eventPublisher EventPublisher) Service {
	return &service{
		logger:         logger,
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}

func (s *service) CreateOrUpdateNote(ctx context.Context, id *string, title, description string) (*notes.Note, error) {
	if id == nil || *id == "" {
		return s.createNote(ctx, title, description)
	}

	return s.updateNote(ctx, *id, title, description)
}

func (s *service) createNote(ctx context.Context, title, description string) (*notes.Note, error) {
	newID := uuid.New().String()
	note := &notes.Note{
		ID:          &newID,
		Title:       title,
		Description: description,
	}

	var createdNote *notes.Note
	err := s.repo.RunInTransaction(ctx, func(ctx context.Context) error {
		querier := s.repo.NewNotesQuery(ctx)

		if err := querier.Create(ctx, note); err != nil {
			s.logger.ErrorContext(ctx,
				"failed to create note",
				"error", err,
				"note_id", newID,
			)
			return err
		}

		var getErr error
		createdNote, getErr = querier.GetByID(ctx, newID)
		if getErr != nil {
			s.logger.ErrorContext(ctx,
				"failed to get created note",
				"error", getErr,
				"note_id", newID,
			)
			return getErr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if s.eventPublisher != nil && createdNote != nil && createdNote.ID != nil {
		s.eventPublisher.Publish(ctx, event.Event{
			Type:   event.MessageTypeNoteCreated,
			NoteID: *createdNote.ID,
		})
	}

	return createdNote, nil
}

func (s *service) updateNote(ctx context.Context, id, title, description string) (*notes.Note, error) {
	var updatedNote *notes.Note

	err := s.repo.RunInTransaction(ctx, func(ctx context.Context) error {
		querier := s.repo.NewNotesQuery(ctx)

		note, err := querier.GetByID(ctx, id)
		if err != nil {
			s.logger.ErrorContext(ctx,
				"failed to get note for update",
				"error", err,
				"note_id", id,
			)
			return err
		}
		if note == nil {
			return exception.ErrNoteNotFound
		}

		note.Title = title
		note.Description = description

		if updateErr := querier.Update(ctx, note); updateErr != nil {
			s.logger.ErrorContext(ctx,
				"failed to update note",
				"error", updateErr,
				"note_id", id,
			)
			return updateErr
		}

		var getErr error
		updatedNote, getErr = querier.GetByID(ctx, id)
		if getErr != nil {
			s.logger.ErrorContext(ctx,
				"failed to get updated note",
				"error", getErr,
				"note_id", id,
			)
			return getErr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedNote, nil
}

func (s *service) GetNote(ctx context.Context, id string) (*notes.Note, error) {
	if id == "" {
		return nil, exception.ErrInvalidInput
	}

	note, err := s.repo.NewNotesQuery(ctx).GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx,
			"failed to get note",
			"error", err,
			"note_id", id,
		)
		return nil, err
	}
	if note == nil {
		return nil, exception.ErrNoteNotFound
	}

	return note, nil
}

func (s *service) ListNotes(ctx context.Context) ([]*notes.Note, error) {
	notesList, err := s.repo.NewNotesQuery(ctx).List(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx,
			"failed to list notes",
			"error", err,
		)
		return nil, err
	}

	return notesList, nil
}

func (s *service) DeleteNote(ctx context.Context, id string) error {
	if id == "" {
		return exception.ErrInvalidInput
	}

	querier := s.repo.NewNotesQuery(ctx)

	note, err := querier.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx,
			"failed to get note for deletion",
			"error", err,
			"note_id", id,
		)
		return err
	}
	if note == nil {
		return exception.ErrNoteNotFound
	}

	if err := querier.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx,
			"failed to delete note",
			"error", err,
			"note_id", id,
		)
		return err
	}

	return nil
}
