// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: certificate_revocation_list.sql

package queries

import (
	"context"
	"database/sql"
	"time"
)

const createCertificateRevocationList = `-- name: CreateCertificateRevocationList :one
INSERT INTO certificate_revocation_list(
    name,
    signature,
    this_update,
    next_update,
    url,
    raw
) VALUES (?,?,?,?,?,?)
  ON CONFLICT DO UPDATE SET
    this_update = excluded.this_update,
    next_update = excluded.next_update
RETURNING id
`

type CreateCertificateRevocationListParams struct {
	Name       string
	Signature  []byte
	ThisUpdate time.Time
	NextUpdate sql.NullTime
	Url        sql.NullString
	Raw        []byte
}

func (q *Queries) CreateCertificateRevocationList(ctx context.Context, arg CreateCertificateRevocationListParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createCertificateRevocationList,
		arg.Name,
		arg.Signature,
		arg.ThisUpdate,
		arg.NextUpdate,
		arg.Url,
		arg.Raw,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getCertificateRevocationList = `-- name: GetCertificateRevocationList :one
SELECT id, name, signature, DATETIME(this_update) as this_update, DATETIME(next_update) as next_update, url, raw FROM certificate_revocation_list
WHERE name = ?
`

type GetCertificateRevocationListRow struct {
	ID         int64
	Name       string
	Signature  []byte
	ThisUpdate interface{}
	NextUpdate interface{}
	Url        sql.NullString
	Raw        []byte
}

func (q *Queries) GetCertificateRevocationList(ctx context.Context, name string) (GetCertificateRevocationListRow, error) {
	row := q.db.QueryRowContext(ctx, getCertificateRevocationList, name)
	var i GetCertificateRevocationListRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Signature,
		&i.ThisUpdate,
		&i.NextUpdate,
		&i.Url,
		&i.Raw,
	)
	return i, err
}

const listCertificateRevocationLists = `-- name: ListCertificateRevocationLists :many
SELECT id, name, signature, DATETIME(this_update) as this_update, DATETIME(next_update) as next_update, url, raw  FROM certificate_revocation_list
ORDER BY id
`

type ListCertificateRevocationListsRow struct {
	ID         int64
	Name       string
	Signature  []byte
	ThisUpdate interface{}
	NextUpdate interface{}
	Url        sql.NullString
	Raw        []byte
}

func (q *Queries) ListCertificateRevocationLists(ctx context.Context) ([]ListCertificateRevocationListsRow, error) {
	rows, err := q.db.QueryContext(ctx, listCertificateRevocationLists)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListCertificateRevocationListsRow
	for rows.Next() {
		var i ListCertificateRevocationListsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Signature,
			&i.ThisUpdate,
			&i.NextUpdate,
			&i.Url,
			&i.Raw,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCertificateRevocationList = `-- name: UpdateCertificateRevocationList :one
UPDATE certificate_revocation_list
set this_update = ?,
    next_update = ?,
    raw = ?
WHERE name = ?
RETURNING id, name, signature, this_update, next_update, url, raw
`

type UpdateCertificateRevocationListParams struct {
	ThisUpdate time.Time
	NextUpdate sql.NullTime
	Raw        []byte
	Name       string
}

func (q *Queries) UpdateCertificateRevocationList(ctx context.Context, arg UpdateCertificateRevocationListParams) (CertificateRevocationList, error) {
	row := q.db.QueryRowContext(ctx, updateCertificateRevocationList,
		arg.ThisUpdate,
		arg.NextUpdate,
		arg.Raw,
		arg.Name,
	)
	var i CertificateRevocationList
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Signature,
		&i.ThisUpdate,
		&i.NextUpdate,
		&i.Url,
		&i.Raw,
	)
	return i, err
}
