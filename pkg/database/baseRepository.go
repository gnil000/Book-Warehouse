package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type BaseRepository[T any] struct {
	DB *sql.DB
}

func (repo *BaseRepository[T]) SelectMultiple(mapRow func(*sql.Rows, *T) error, query string, args ...any) ([]*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			return
		}
	}(rows)

	var list []*T

	for rows.Next() {
		var t T
		if err := mapRow(rows, &t); err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (repo *BaseRepository[T]) SelectSingle(mapRow func(*sql.Row, *T) error, query string, args ...any) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := repo.DB.QueryRowContext(ctx, query, args...)
	var t T
	if err := mapRow(row, &t); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, err
	}
	return &t, nil
}

func (repo *BaseRepository[T]) Insert(query string, args ...any) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	query += " RETURNING id"

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *BaseRepository[T]) ExecuteQuery(query string, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
