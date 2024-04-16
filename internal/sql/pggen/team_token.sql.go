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

const insertTeamTokenSQL = `INSERT INTO team_tokens (
    team_token_id,
    created_at,
    team_id,
    expiry
) VALUES (
    $1,
    $2,
    $3,
    $4
) ON CONFLICT (team_id) DO UPDATE
  SET team_token_id = $1,
      created_at    = $2,
      expiry        = $4;`

type InsertTeamTokenParams struct {
	TeamTokenID pgtype.Text        `json:"team_token_id"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	TeamID      pgtype.Text        `json:"team_id"`
	Expiry      pgtype.Timestamptz `json:"expiry"`
}

// InsertTeamToken implements Querier.InsertTeamToken.
func (q *DBQuerier) InsertTeamToken(ctx context.Context, params InsertTeamTokenParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertTeamToken")
	cmdTag, err := q.conn.Exec(ctx, insertTeamTokenSQL, params.TeamTokenID, params.CreatedAt, params.TeamID, params.Expiry)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertTeamToken: %w", err)
	}
	return cmdTag, err
}

const findTeamTokensByIDSQL = `SELECT *
FROM team_tokens
WHERE team_id = $1
;`

type FindTeamTokensByIDRow struct {
	TeamTokenID pgtype.Text        `json:"team_token_id"`
	Description pgtype.Text        `json:"description"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	TeamID      pgtype.Text        `json:"team_id"`
	Expiry      pgtype.Timestamptz `json:"expiry"`
}

// FindTeamTokensByID implements Querier.FindTeamTokensByID.
func (q *DBQuerier) FindTeamTokensByID(ctx context.Context, teamID pgtype.Text) ([]FindTeamTokensByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindTeamTokensByID")
	rows, err := q.conn.Query(ctx, findTeamTokensByIDSQL, teamID)
	if err != nil {
		return nil, fmt.Errorf("query FindTeamTokensByID: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindTeamTokensByIDRow, error) {
		var item FindTeamTokensByIDRow
		if err := row.Scan(&item.TeamTokenID, // 'team_token_id', 'TeamTokenID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Description, // 'description', 'Description', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.CreatedAt,   // 'created_at', 'CreatedAt', 'pgtype.Timestamptz', 'github.com/jackc/pgx/v5/pgtype', 'Timestamptz'
			&item.TeamID,      // 'team_id', 'TeamID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Expiry,      // 'expiry', 'Expiry', 'pgtype.Timestamptz', 'github.com/jackc/pgx/v5/pgtype', 'Timestamptz'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const deleteTeamTokenByIDSQL = `DELETE
FROM team_tokens
WHERE team_id = $1
RETURNING team_token_id
;`

// DeleteTeamTokenByID implements Querier.DeleteTeamTokenByID.
func (q *DBQuerier) DeleteTeamTokenByID(ctx context.Context, teamID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteTeamTokenByID")
	rows, err := q.conn.Query(ctx, deleteTeamTokenByIDSQL, teamID)
	if err != nil {
		return pgtype.Text{}, fmt.Errorf("query DeleteTeamTokenByID: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (pgtype.Text, error) {
		var item pgtype.Text
		if err := row.Scan(&item); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
