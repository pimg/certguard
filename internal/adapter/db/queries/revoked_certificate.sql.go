// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: revoked_certificate.sql

package queries

import (
	"context"
	"time"
)

const createRevokedCertificates = `-- name: CreateRevokedCertificates :exec
INSERT INTO revoked_certificate(
    serialnumber,
    revocation_date,
    reason,
    revocation_list
) VALUES (
          ?,?,?,?
)
ON CONFLICT DO NOTHING
`

type CreateRevokedCertificatesParams struct {
	Serialnumber   string
	RevocationDate time.Time
	Reason         string
	RevocationList int64
}

func (q *Queries) CreateRevokedCertificates(ctx context.Context, arg CreateRevokedCertificatesParams) error {
	_, err := q.db.ExecContext(ctx, createRevokedCertificates,
		arg.Serialnumber,
		arg.RevocationDate,
		arg.Reason,
		arg.RevocationList,
	)
	return err
}

const getRevokedCertificate = `-- name: GetRevokedCertificate :one
SELECT cert.id, cert.serialnumber, cert.reason, DATETIME(cert.revocation_date) as revocation_date, crl.name AS revoked_by
FROM revoked_certificate as cert
JOIN certificate_revocation_list AS crl ON crl.id = cert.revocation_list
WHERE serialnumber = ?
`

type GetRevokedCertificateRow struct {
	ID             int64
	Serialnumber   string
	Reason         string
	RevocationDate interface{}
	RevokedBy      string
}

func (q *Queries) GetRevokedCertificate(ctx context.Context, serialnumber string) (GetRevokedCertificateRow, error) {
	row := q.db.QueryRowContext(ctx, getRevokedCertificate, serialnumber)
	var i GetRevokedCertificateRow
	err := row.Scan(
		&i.ID,
		&i.Serialnumber,
		&i.Reason,
		&i.RevocationDate,
		&i.RevokedBy,
	)
	return i, err
}

const getRevokedCertificatesByRevocationList = `-- name: GetRevokedCertificatesByRevocationList :many
SELECT id, serialnumber, DATETIME(revocation_date) as revocation_date, reason, revocation_list
FROM revoked_certificate
WHERE revocation_list = ?
ORDER BY revocation_date
`

type GetRevokedCertificatesByRevocationListRow struct {
	ID             int64
	Serialnumber   string
	RevocationDate interface{}
	Reason         string
	RevocationList int64
}

func (q *Queries) GetRevokedCertificatesByRevocationList(ctx context.Context, revocationList int64) ([]GetRevokedCertificatesByRevocationListRow, error) {
	rows, err := q.db.QueryContext(ctx, getRevokedCertificatesByRevocationList, revocationList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRevokedCertificatesByRevocationListRow
	for rows.Next() {
		var i GetRevokedCertificatesByRevocationListRow
		if err := rows.Scan(
			&i.ID,
			&i.Serialnumber,
			&i.RevocationDate,
			&i.Reason,
			&i.RevocationList,
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
