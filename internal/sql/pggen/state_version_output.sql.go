// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ genericConn = (*pgx.Conn)(nil)

const insertStateVersionOutputSQL = `INSERT INTO state_version_outputs (
    state_version_output_id,
    name,
    sensitive,
    type,
    value,
    state_version_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);`

type InsertStateVersionOutputParams struct {
	ID             pgtype.Text `json:"id"`
	Name           pgtype.Text `json:"name"`
	Sensitive      pgtype.Bool `json:"sensitive"`
	Type           pgtype.Text `json:"type"`
	Value          []byte      `json:"value"`
	StateVersionID pgtype.Text `json:"state_version_id"`
}

// InsertStateVersionOutput implements Querier.InsertStateVersionOutput.
func (q *DBQuerier) InsertStateVersionOutput(ctx context.Context, params InsertStateVersionOutputParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertStateVersionOutput")
	cmdTag, err := q.conn.Exec(ctx, insertStateVersionOutputSQL, params.ID, params.Name, params.Sensitive, params.Type, params.Value, params.StateVersionID)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertStateVersionOutput: %w", err)
	}
	return cmdTag, err
}

const findStateVersionOutputByIDSQL = `SELECT *
FROM state_version_outputs
WHERE state_version_output_id = $1
;`

type FindStateVersionOutputByIDRow struct {
	StateVersionOutputID pgtype.Text `json:"state_version_output_id"`
	Name                 pgtype.Text `json:"name"`
	Sensitive            pgtype.Bool `json:"sensitive"`
	Type                 pgtype.Text `json:"type"`
	Value                []byte      `json:"value"`
	StateVersionID       pgtype.Text `json:"state_version_id"`
}

// FindStateVersionOutputByID implements Querier.FindStateVersionOutputByID.
func (q *DBQuerier) FindStateVersionOutputByID(ctx context.Context, id pgtype.Text) (FindStateVersionOutputByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindStateVersionOutputByID")
	rows, err := q.conn.Query(ctx, findStateVersionOutputByIDSQL, id)
	if err != nil {
		return FindStateVersionOutputByIDRow{}, fmt.Errorf("query FindStateVersionOutputByID: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (FindStateVersionOutputByIDRow, error) {
		var item FindStateVersionOutputByIDRow
		if err := row.Scan(&item.StateVersionOutputID, // 'state_version_output_id', 'StateVersionOutputID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Name,           // 'name', 'Name', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Sensitive,      // 'sensitive', 'Sensitive', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.Type,           // 'type', 'Type', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Value,          // 'value', 'Value', '[]byte', '', '[]byte'
			&item.StateVersionID, // 'state_version_id', 'StateVersionID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
