package postgres

import (
	"database/sql"
	"errors"

	"github.com/gbih/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires, mv1, mv2, mv3 string) (int, error) {

	var id int

	stmt := `INSERT INTO snippets
	(title, content, created, expires, mv1, mv2, mv3)
	VALUES
	($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 DAY' * $3, $4, $5, $6)
	RETURNING id`

	err := m.DB.QueryRow(stmt, title, content, expires, mv1, mv2, mv3).Scan(&id)

	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires, mv1, mv2, mv3 FROM snippets
	WHERE expires > CURRENT_TIMESTAMP AND id = $1`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.MV1, &s.MV2, &s.MV3)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires, mv1, mv2, mv3 FROM snippets
	WHERE expires > CURRENT_TIMESTAMP ORDER BY created DESC LIMIT 3`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.MV1, &s.MV2, &s.MV3)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
