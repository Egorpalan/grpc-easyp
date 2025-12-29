package notes

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/Egorpalan/grpc-easyp/internal/lib/postgres"
	"github.com/Egorpalan/grpc-easyp/internal/model/notes"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Querier interface {
	Create(ctx context.Context, note *notes.Note) error
	GetByID(ctx context.Context, id string) (*notes.Note, error)
	List(ctx context.Context) ([]*notes.Note, error)
	Update(ctx context.Context, note *notes.Note) error
	Delete(ctx context.Context, id string) error
}

type Query struct {
	ctx         context.Context
	queryEngine postgres.QueryEngine
}

func NewQuery(ctx context.Context, queryEngine postgres.QueryEngine) Querier {
	return &Query{
		ctx:         ctx,
		queryEngine: queryEngine,
	}
}

func (q *Query) Create(ctx context.Context, note *notes.Note) error {
	if note.ID == nil {
		return errors.New("note ID is required")
	}

	qb := psql.Insert(notesTable).
		Columns("id", "title", "description", "created_at", "updated_at").
		Values(note.ID, note.Title, note.Description, squirrel.Expr("NOW()"), squirrel.Expr("NOW()"))

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = q.queryEngine.Exec(ctx, query, args...)
	return err
}

func (q *Query) GetByID(ctx context.Context, id string) (*notes.Note, error) {
	qb := psql.Select("id", "title", "description", "created_at", "updated_at").
		From(notesTable).
		Where(squirrel.Eq{"id": id})

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var result Note
	err = q.queryEngine.QueryRow(ctx, query, args...).Scan(
		&result.ID,
		&result.Title,
		&result.Description,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result.ToServiceModel(), nil
}

func (q *Query) List(ctx context.Context) ([]*notes.Note, error) {
	qb := psql.Select("id", "title", "description", "created_at", "updated_at").
		From(notesTable).
		OrderBy("created_at DESC")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := q.queryEngine.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notesList []*notes.Note
	for rows.Next() {
		var data Note
		if err := rows.Scan(
			&data.ID,
			&data.Title,
			&data.Description,
			&data.CreatedAt,
			&data.UpdatedAt,
		); err != nil {
			return nil, err
		}
		notesList = append(notesList, data.ToServiceModel())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notesList, nil
}

func (q *Query) Update(ctx context.Context, note *notes.Note) error {
	if note.ID == nil {
		return errors.New("note ID is required")
	}

	qb := psql.Update(notesTable).
		Set("title", note.Title).
		Set("description", note.Description).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": *note.ID})

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	result, err := q.queryEngine.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}

	return nil
}

func (q *Query) Delete(ctx context.Context, id string) error {
	qb := psql.Delete(notesTable).
		Where(squirrel.Eq{"id": id})

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	result, err := q.queryEngine.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}

	return nil
}
