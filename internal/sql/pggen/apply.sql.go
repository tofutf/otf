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

const insertApplySQL = `INSERT INTO applies (
    run_id,
    status
) VALUES (
    $1,
    $2
);`

// InsertApply implements Querier.InsertApply.
func (q *DBQuerier) InsertApply(ctx context.Context, runID pgtype.Text, status pgtype.Text) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertApply")
	cmdTag, err := q.conn.Exec(ctx, insertApplySQL, runID, status)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertApply: %w", err)
	}
	return cmdTag, err
}

const updateAppliedChangesByIDSQL = `UPDATE applies
SET resource_report = (
    $1,
    $2,
    $3
)
WHERE run_id = $4
RETURNING run_id
;`

type UpdateAppliedChangesByIDParams struct {
	Additions    pgtype.Int4 `json:"additions"`
	Changes      pgtype.Int4 `json:"changes"`
	Destructions pgtype.Int4 `json:"destructions"`
	RunID        pgtype.Text `json:"run_id"`
}

// UpdateAppliedChangesByID implements Querier.UpdateAppliedChangesByID.
func (q *DBQuerier) UpdateAppliedChangesByID(ctx context.Context, params UpdateAppliedChangesByIDParams) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateAppliedChangesByID")
	rows, err := q.conn.Query(ctx, updateAppliedChangesByIDSQL, params.Additions, params.Changes, params.Destructions, params.RunID)
	if err != nil {
		return pgtype.Text{}, fmt.Errorf("query UpdateAppliedChangesByID: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (pgtype.Text, error) {
		var item pgtype.Text
		if err := row.Scan(&item); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const updateApplyStatusByIDSQL = `UPDATE applies
SET status = $1
WHERE run_id = $2
RETURNING run_id
;`

// UpdateApplyStatusByID implements Querier.UpdateApplyStatusByID.
func (q *DBQuerier) UpdateApplyStatusByID(ctx context.Context, status pgtype.Text, runID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateApplyStatusByID")
	rows, err := q.conn.Query(ctx, updateApplyStatusByIDSQL, status, runID)
	if err != nil {
		return pgtype.Text{}, fmt.Errorf("query UpdateApplyStatusByID: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (pgtype.Text, error) {
		var item pgtype.Text
		if err := row.Scan(&item); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
