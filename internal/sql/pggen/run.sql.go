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
var _ RegisterConn = (*pgx.Conn)(nil)

const insertRunSQL = `INSERT INTO runs (
    run_id,
    created_at,
    is_destroy,
    position_in_queue,
    refresh,
    refresh_only,
    source,
    status,
    replace_addrs,
    target_addrs,
    auto_apply,
    plan_only,
    configuration_version_id,
    workspace_id,
    created_by,
    terraform_version,
    allow_empty_apply
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17
);`

type InsertRunParams struct {
	ID                     pgtype.Text        `json:"id"`
	CreatedAt              pgtype.Timestamptz `json:"created_at"`
	IsDestroy              pgtype.Bool        `json:"is_destroy"`
	PositionInQueue        pgtype.Int4        `json:"position_in_queue"`
	Refresh                pgtype.Bool        `json:"refresh"`
	RefreshOnly            pgtype.Bool        `json:"refresh_only"`
	Source                 pgtype.Text        `json:"source"`
	Status                 pgtype.Text        `json:"status"`
	ReplaceAddrs           []string           `json:"replace_addrs"`
	TargetAddrs            []string           `json:"target_addrs"`
	AutoApply              pgtype.Bool        `json:"auto_apply"`
	PlanOnly               pgtype.Bool        `json:"plan_only"`
	ConfigurationVersionID pgtype.Text        `json:"configuration_version_id"`
	WorkspaceID            pgtype.Text        `json:"workspace_id"`
	CreatedBy              pgtype.Text        `json:"created_by"`
	TerraformVersion       pgtype.Text        `json:"terraform_version"`
	AllowEmptyApply        pgtype.Bool        `json:"allow_empty_apply"`
}

// InsertRun implements Querier.InsertRun.
func (q *DBQuerier) InsertRun(ctx context.Context, params InsertRunParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertRun")
	cmdTag, err := q.conn.Exec(ctx, insertRunSQL, params.ID, params.CreatedAt, params.IsDestroy, params.PositionInQueue, params.Refresh, params.RefreshOnly, params.Source, params.Status, params.ReplaceAddrs, params.TargetAddrs, params.AutoApply, params.PlanOnly, params.ConfigurationVersionID, params.WorkspaceID, params.CreatedBy, params.TerraformVersion, params.AllowEmptyApply)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertRun: %w", err)
	}
	return cmdTag, err
}

const insertRunStatusTimestampSQL = `INSERT INTO run_status_timestamps (
    run_id,
    status,
    timestamp
) VALUES (
    $1,
    $2,
    $3
);`

type InsertRunStatusTimestampParams struct {
	ID        pgtype.Text        `json:"id"`
	Status    pgtype.Text        `json:"status"`
	Timestamp pgtype.Timestamptz `json:"timestamp"`
}

// InsertRunStatusTimestamp implements Querier.InsertRunStatusTimestamp.
func (q *DBQuerier) InsertRunStatusTimestamp(ctx context.Context, params InsertRunStatusTimestampParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertRunStatusTimestamp")
	cmdTag, err := q.conn.Exec(ctx, insertRunStatusTimestampSQL, params.ID, params.Status, params.Timestamp)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertRunStatusTimestamp: %w", err)
	}
	return cmdTag, err
}

const insertRunVariableSQL = `INSERT INTO run_variables (
    run_id,
    key,
    value
) VALUES (
    $1,
    $2,
    $3
);`

type InsertRunVariableParams struct {
	RunID pgtype.Text `json:"run_id"`
	Key   pgtype.Text `json:"key"`
	Value pgtype.Text `json:"value"`
}

// InsertRunVariable implements Querier.InsertRunVariable.
func (q *DBQuerier) InsertRunVariable(ctx context.Context, params InsertRunVariableParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertRunVariable")
	cmdTag, err := q.conn.Exec(ctx, insertRunVariableSQL, params.RunID, params.Key, params.Value)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertRunVariable: %w", err)
	}
	return cmdTag, err
}

const findRunsSQL = `SELECT
    runs.run_id,
    runs.created_at,
    runs.cancel_signaled_at,
    runs.is_destroy,
    runs.position_in_queue,
    runs.refresh,
    runs.refresh_only,
    runs.source,
    runs.status,
    plans.status      AS plan_status,
    applies.status      AS apply_status,
    runs.replace_addrs,
    runs.target_addrs,
    runs.auto_apply,
    plans.resource_report AS plan_resource_report,
    plans.output_report AS plan_output_report,
    applies.resource_report AS apply_resource_report,
    runs.configuration_version_id,
    runs.workspace_id,
    runs.plan_only,
    runs.created_by,
    runs.terraform_version,
    runs.allow_empty_apply,
    workspaces.execution_mode AS execution_mode,
    CASE WHEN workspaces.latest_run_id = runs.run_id THEN true
         ELSE false
    END AS latest,
    workspaces.organization_name,
    organizations.cost_estimation_enabled,
    (ia.*)::"ingress_attributes" AS ingress_attributes,
    (
        SELECT array_agg(rst.*) AS run_status_timestamps
        FROM run_status_timestamps rst
        WHERE rst.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = plans.run_id
        AND   st.phase = 'plan'
        GROUP BY run_id, phase
    ) AS plan_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = applies.run_id
        AND   st.phase = 'apply'
        GROUP BY run_id, phase
    ) AS apply_status_timestamps,
    (
        SELECT array_agg(v.*) AS run_variables
        FROM run_variables v
        WHERE v.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_variables
FROM runs
JOIN plans USING (run_id)
JOIN applies USING (run_id)
JOIN (configuration_versions LEFT JOIN ingress_attributes ia USING (configuration_version_id)) USING (configuration_version_id)
JOIN workspaces ON runs.workspace_id = workspaces.workspace_id
JOIN organizations ON workspaces.organization_name = organizations.name
WHERE
    workspaces.organization_name LIKE ANY($1)
AND workspaces.workspace_id      LIKE ANY($2)
AND workspaces.name              LIKE ANY($3)
AND runs.source                  LIKE ANY($4)
AND runs.status                  LIKE ANY($5)
AND runs.plan_only::text         LIKE ANY($6)
AND (($7::text IS NULL) OR ia.commit_sha = $7)
AND (($8::text IS NULL) OR ia.sender_username = $8)
ORDER BY runs.created_at DESC
LIMIT $9 OFFSET $10
;`

type FindRunsParams struct {
	OrganizationNames []string    `json:"organization_names"`
	WorkspaceIds      []string    `json:"workspace_ids"`
	WorkspaceNames    []string    `json:"workspace_names"`
	Sources           []string    `json:"sources"`
	Statuses          []string    `json:"statuses"`
	PlanOnly          []string    `json:"plan_only"`
	CommitSHA         pgtype.Text `json:"commit_sha"`
	VCSUsername       pgtype.Text `json:"vcs_username"`
	Limit             pgtype.Int8 `json:"limit"`
	Offset            pgtype.Int8 `json:"offset"`
}

type FindRunsRow struct {
	RunID                  pgtype.Text              `json:"run_id"`
	CreatedAt              pgtype.Timestamptz       `json:"created_at"`
	CancelSignaledAt       pgtype.Timestamptz       `json:"cancel_signaled_at"`
	IsDestroy              pgtype.Bool              `json:"is_destroy"`
	PositionInQueue        pgtype.Int4              `json:"position_in_queue"`
	Refresh                pgtype.Bool              `json:"refresh"`
	RefreshOnly            pgtype.Bool              `json:"refresh_only"`
	Source                 pgtype.Text              `json:"source"`
	Status                 pgtype.Text              `json:"status"`
	PlanStatus             pgtype.Text              `json:"plan_status"`
	ApplyStatus            pgtype.Text              `json:"apply_status"`
	ReplaceAddrs           []string                 `json:"replace_addrs"`
	TargetAddrs            []string                 `json:"target_addrs"`
	AutoApply              pgtype.Bool              `json:"auto_apply"`
	PlanResourceReport     *Report                  `json:"plan_resource_report"`
	PlanOutputReport       *Report                  `json:"plan_output_report"`
	ApplyResourceReport    *Report                  `json:"apply_resource_report"`
	ConfigurationVersionID pgtype.Text              `json:"configuration_version_id"`
	WorkspaceID            pgtype.Text              `json:"workspace_id"`
	PlanOnly               pgtype.Bool              `json:"plan_only"`
	CreatedBy              pgtype.Text              `json:"created_by"`
	TerraformVersion       pgtype.Text              `json:"terraform_version"`
	AllowEmptyApply        pgtype.Bool              `json:"allow_empty_apply"`
	ExecutionMode          pgtype.Text              `json:"execution_mode"`
	Latest                 pgtype.Bool              `json:"latest"`
	OrganizationName       pgtype.Text              `json:"organization_name"`
	CostEstimationEnabled  pgtype.Bool              `json:"cost_estimation_enabled"`
	IngressAttributes      *IngressAttributes       `json:"ingress_attributes"`
	RunStatusTimestamps    []*RunStatusTimestamps   `json:"run_status_timestamps"`
	PlanStatusTimestamps   []*PhaseStatusTimestamps `json:"plan_status_timestamps"`
	ApplyStatusTimestamps  []*PhaseStatusTimestamps `json:"apply_status_timestamps"`
	RunVariables           []*RunVariables          `json:"run_variables"`
}

// FindRuns implements Querier.FindRuns.
func (q *DBQuerier) FindRuns(ctx context.Context, params FindRunsParams) ([]FindRunsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRuns")
	rows, err := q.conn.Query(ctx, findRunsSQL, params.OrganizationNames, params.WorkspaceIds, params.WorkspaceNames, params.Sources, params.Statuses, params.PlanOnly, params.CommitSHA, params.VCSUsername, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query FindRuns: %w", err)
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindRunsRow, error) {
		var item FindRunsRow
		if err := row.Scan(&item.RunID, // 'run_id', 'RunID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.CreatedAt,              // 'created_at', 'CreatedAt', 'pgtype.Timestamptz', 'github.com/jackc/pgx/v5/pgtype', 'Timestamptz'
			&item.CancelSignaledAt,       // 'cancel_signaled_at', 'CancelSignaledAt', 'pgtype.Timestamptz', 'github.com/jackc/pgx/v5/pgtype', 'Timestamptz'
			&item.IsDestroy,              // 'is_destroy', 'IsDestroy', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.PositionInQueue,        // 'position_in_queue', 'PositionInQueue', 'pgtype.Int4', 'github.com/jackc/pgx/v5/pgtype', 'Int4'
			&item.Refresh,                // 'refresh', 'Refresh', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.RefreshOnly,            // 'refresh_only', 'RefreshOnly', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.Source,                 // 'source', 'Source', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Status,                 // 'status', 'Status', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.PlanStatus,             // 'plan_status', 'PlanStatus', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.ApplyStatus,            // 'apply_status', 'ApplyStatus', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.ReplaceAddrs,           // 'replace_addrs', 'ReplaceAddrs', '[]string', '', '[]string'
			&item.TargetAddrs,            // 'target_addrs', 'TargetAddrs', '[]string', '', '[]string'
			&item.AutoApply,              // 'auto_apply', 'AutoApply', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.PlanResourceReport,     // 'plan_resource_report', 'PlanResourceReport', '*Report', '', '*Report'
			&item.PlanOutputReport,       // 'plan_output_report', 'PlanOutputReport', '*Report', '', '*Report'
			&item.ApplyResourceReport,    // 'apply_resource_report', 'ApplyResourceReport', '*Report', '', '*Report'
			&item.ConfigurationVersionID, // 'configuration_version_id', 'ConfigurationVersionID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.WorkspaceID,            // 'workspace_id', 'WorkspaceID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.PlanOnly,               // 'plan_only', 'PlanOnly', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.CreatedBy,              // 'created_by', 'CreatedBy', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.TerraformVersion,       // 'terraform_version', 'TerraformVersion', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.AllowEmptyApply,        // 'allow_empty_apply', 'AllowEmptyApply', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.ExecutionMode,          // 'execution_mode', 'ExecutionMode', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Latest,                 // 'latest', 'Latest', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.OrganizationName,       // 'organization_name', 'OrganizationName', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.CostEstimationEnabled,  // 'cost_estimation_enabled', 'CostEstimationEnabled', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.IngressAttributes,      // 'ingress_attributes', 'IngressAttributes', '*IngressAttributes', '', '*IngressAttributes'
			&item.RunStatusTimestamps,    // 'run_status_timestamps', 'RunStatusTimestamps', '[]*RunStatusTimestamps', '', '[]*RunStatusTimestamps'
			&item.PlanStatusTimestamps,   // 'plan_status_timestamps', 'PlanStatusTimestamps', '[]*PhaseStatusTimestamps', '', '[]*PhaseStatusTimestamps'
			&item.ApplyStatusTimestamps,  // 'apply_status_timestamps', 'ApplyStatusTimestamps', '[]*PhaseStatusTimestamps', '', '[]*PhaseStatusTimestamps'
			&item.RunVariables,           // 'run_variables', 'RunVariables', '[]*RunVariables', '', '[]*RunVariables'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const countRunsSQL = `SELECT count(*)
FROM runs
JOIN workspaces USING(workspace_id)
JOIN (configuration_versions LEFT JOIN ingress_attributes ia USING (configuration_version_id)) USING (configuration_version_id)
WHERE
    workspaces.organization_name LIKE ANY($1)
AND workspaces.workspace_id      LIKE ANY($2)
AND workspaces.name              LIKE ANY($3)
AND runs.source                  LIKE ANY($4)
AND runs.status                  LIKE ANY($5)
AND runs.plan_only::text         LIKE ANY($6)
AND (($7::text IS NULL) OR ia.commit_sha = $7)
AND (($8::text IS NULL) OR ia.sender_username = $8)
;`

type CountRunsParams struct {
	OrganizationNames []string    `json:"organization_names"`
	WorkspaceIds      []string    `json:"workspace_ids"`
	WorkspaceNames    []string    `json:"workspace_names"`
	Sources           []string    `json:"sources"`
	Statuses          []string    `json:"statuses"`
	PlanOnly          []string    `json:"plan_only"`
	CommitSHA         pgtype.Text `json:"commit_sha"`
	VCSUsername       pgtype.Text `json:"vcs_username"`
}

// CountRuns implements Querier.CountRuns.
func (q *DBQuerier) CountRuns(ctx context.Context, params CountRunsParams) (pgtype.Int8, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CountRuns")
	rows, err := q.conn.Query(ctx, countRunsSQL, params.OrganizationNames, params.WorkspaceIds, params.WorkspaceNames, params.Sources, params.Statuses, params.PlanOnly, params.CommitSHA, params.VCSUsername)
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("query CountRuns: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (pgtype.Int8, error) {
		var item pgtype.Int8
		if err := row.Scan(&item); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const findRunByIDSQL = `SELECT
    runs.run_id,
    runs.created_at,
    runs.cancel_signaled_at,
    runs.is_destroy,
    runs.position_in_queue,
    runs.refresh,
    runs.refresh_only,
    runs.source,
    runs.status,
    plans.status      AS plan_status,
    applies.status      AS apply_status,
    runs.replace_addrs,
    runs.target_addrs,
    runs.auto_apply,
    plans.resource_report AS plan_resource_report,
    plans.output_report AS plan_output_report,
    applies.resource_report AS apply_resource_report,
    runs.configuration_version_id,
    runs.workspace_id,
    runs.plan_only,
    runs.created_by,
    runs.terraform_version,
    runs.allow_empty_apply,
    workspaces.execution_mode AS execution_mode,
    CASE WHEN workspaces.latest_run_id = runs.run_id THEN true
         ELSE false
    END AS latest,
    workspaces.organization_name,
    organizations.cost_estimation_enabled,
    (ia.*)::"ingress_attributes" AS ingress_attributes,
    (
        SELECT array_agg(rst.*) AS run_status_timestamps
        FROM run_status_timestamps rst
        WHERE rst.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = plans.run_id
        AND   st.phase = 'plan'
        GROUP BY run_id, phase
    ) AS plan_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = applies.run_id
        AND   st.phase = 'apply'
        GROUP BY run_id, phase
    ) AS apply_status_timestamps,
    (
        SELECT array_agg(v.*) AS run_variables
        FROM run_variables v
        WHERE v.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_variables
FROM runs
JOIN plans USING (run_id)
JOIN applies USING (run_id)
JOIN (configuration_versions LEFT JOIN ingress_attributes ia USING (configuration_version_id)) USING (configuration_version_id)
JOIN workspaces ON runs.workspace_id = workspaces.workspace_id
JOIN organizations ON workspaces.organization_name = organizations.name
WHERE runs.run_id = $1
;`

type FindRunByIDRow struct {
	RunID                  pgtype.Text              `json:"run_id"`
	CreatedAt              pgtype.Timestamptz       `json:"created_at"`
	CancelSignaledAt       pgtype.Timestamptz       `json:"cancel_signaled_at"`
	IsDestroy              pgtype.Bool              `json:"is_destroy"`
	PositionInQueue        pgtype.Int4              `json:"position_in_queue"`
	Refresh                pgtype.Bool              `json:"refresh"`
	RefreshOnly            pgtype.Bool              `json:"refresh_only"`
	Source                 pgtype.Text              `json:"source"`
	Status                 pgtype.Text              `json:"status"`
	PlanStatus             pgtype.Text              `json:"plan_status"`
	ApplyStatus            pgtype.Text              `json:"apply_status"`
	ReplaceAddrs           []string                 `json:"replace_addrs"`
	TargetAddrs            []string                 `json:"target_addrs"`
	AutoApply              pgtype.Bool              `json:"auto_apply"`
	PlanResourceReport     *Report                  `json:"plan_resource_report"`
	PlanOutputReport       *Report                  `json:"plan_output_report"`
	ApplyResourceReport    *Report                  `json:"apply_resource_report"`
	ConfigurationVersionID pgtype.Text              `json:"configuration_version_id"`
	WorkspaceID            pgtype.Text              `json:"workspace_id"`
	PlanOnly               pgtype.Bool              `json:"plan_only"`
	CreatedBy              pgtype.Text              `json:"created_by"`
	TerraformVersion       pgtype.Text              `json:"terraform_version"`
	AllowEmptyApply        pgtype.Bool              `json:"allow_empty_apply"`
	ExecutionMode          pgtype.Text              `json:"execution_mode"`
	Latest                 pgtype.Bool              `json:"latest"`
	OrganizationName       pgtype.Text              `json:"organization_name"`
	CostEstimationEnabled  pgtype.Bool              `json:"cost_estimation_enabled"`
	IngressAttributes      *IngressAttributes       `json:"ingress_attributes"`
	RunStatusTimestamps    []*RunStatusTimestamps   `json:"run_status_timestamps"`
	PlanStatusTimestamps   []*PhaseStatusTimestamps `json:"plan_status_timestamps"`
	ApplyStatusTimestamps  []*PhaseStatusTimestamps `json:"apply_status_timestamps"`
	RunVariables           []*RunVariables          `json:"run_variables"`
}

// FindRunByID implements Querier.FindRunByID.
func (q *DBQuerier) FindRunByID(ctx context.Context, runID pgtype.Text) (FindRunByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunByID")
	rows, err := q.conn.Query(ctx, findRunByIDSQL, runID)
	if err != nil {
		return FindRunByIDRow{}, fmt.Errorf("query FindRunByID: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (FindRunByIDRow, error) {
		var item FindRunByIDRow
		if err := row.Scan(&item.RunID, // 'run_id', 'RunID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.CreatedAt,              // 'created_at', 'CreatedAt', 'pgtype.Timestamptz', 'github.com/jackc/pgx/v5/pgtype', 'Timestamptz'
			&item.CancelSignaledAt,       // 'cancel_signaled_at', 'CancelSignaledAt', 'pgtype.Timestamptz', 'github.com/jackc/pgx/v5/pgtype', 'Timestamptz'
			&item.IsDestroy,              // 'is_destroy', 'IsDestroy', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.PositionInQueue,        // 'position_in_queue', 'PositionInQueue', 'pgtype.Int4', 'github.com/jackc/pgx/v5/pgtype', 'Int4'
			&item.Refresh,                // 'refresh', 'Refresh', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.RefreshOnly,            // 'refresh_only', 'RefreshOnly', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.Source,                 // 'source', 'Source', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Status,                 // 'status', 'Status', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.PlanStatus,             // 'plan_status', 'PlanStatus', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.ApplyStatus,            // 'apply_status', 'ApplyStatus', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.ReplaceAddrs,           // 'replace_addrs', 'ReplaceAddrs', '[]string', '', '[]string'
			&item.TargetAddrs,            // 'target_addrs', 'TargetAddrs', '[]string', '', '[]string'
			&item.AutoApply,              // 'auto_apply', 'AutoApply', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.PlanResourceReport,     // 'plan_resource_report', 'PlanResourceReport', '*Report', '', '*Report'
			&item.PlanOutputReport,       // 'plan_output_report', 'PlanOutputReport', '*Report', '', '*Report'
			&item.ApplyResourceReport,    // 'apply_resource_report', 'ApplyResourceReport', '*Report', '', '*Report'
			&item.ConfigurationVersionID, // 'configuration_version_id', 'ConfigurationVersionID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.WorkspaceID,            // 'workspace_id', 'WorkspaceID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.PlanOnly,               // 'plan_only', 'PlanOnly', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.CreatedBy,              // 'created_by', 'CreatedBy', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.TerraformVersion,       // 'terraform_version', 'TerraformVersion', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.AllowEmptyApply,        // 'allow_empty_apply', 'AllowEmptyApply', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.ExecutionMode,          // 'execution_mode', 'ExecutionMode', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Latest,                 // 'latest', 'Latest', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.OrganizationName,       // 'organization_name', 'OrganizationName', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.CostEstimationEnabled,  // 'cost_estimation_enabled', 'CostEstimationEnabled', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.IngressAttributes,      // 'ingress_attributes', 'IngressAttributes', '*IngressAttributes', '', '*IngressAttributes'
			&item.RunStatusTimestamps,    // 'run_status_timestamps', 'RunStatusTimestamps', '[]*RunStatusTimestamps', '', '[]*RunStatusTimestamps'
			&item.PlanStatusTimestamps,   // 'plan_status_timestamps', 'PlanStatusTimestamps', '[]*PhaseStatusTimestamps', '', '[]*PhaseStatusTimestamps'
			&item.ApplyStatusTimestamps,  // 'apply_status_timestamps', 'ApplyStatusTimestamps', '[]*PhaseStatusTimestamps', '', '[]*PhaseStatusTimestamps'
			&item.RunVariables,           // 'run_variables', 'RunVariables', '[]*RunVariables', '', '[]*RunVariables'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const findRunByIDForUpdateSQL = `SELECT
    runs.run_id,
    runs.created_at,
    runs.cancel_signaled_at,
    runs.is_destroy,
    runs.position_in_queue,
    runs.refresh,
    runs.refresh_only,
    runs.source,
    runs.status,
    plans.status        AS plan_status,
    applies.status      AS apply_status,
    runs.replace_addrs,
    runs.target_addrs,
    runs.auto_apply,
    plans.resource_report AS plan_resource_report,
    plans.output_report AS plan_output_report,
    applies.resource_report AS apply_resource_report,
    runs.configuration_version_id,
    runs.workspace_id,
    runs.plan_only,
    runs.created_by,
    runs.terraform_version,
    runs.allow_empty_apply,
    workspaces.execution_mode AS execution_mode,
    CASE WHEN workspaces.latest_run_id = runs.run_id THEN true
         ELSE false
    END AS latest,
    workspaces.organization_name,
    organizations.cost_estimation_enabled,
    (ia.*)::"ingress_attributes" AS ingress_attributes,
    (
        SELECT array_agg(rst.*) AS run_status_timestamps
        FROM run_status_timestamps rst
        WHERE rst.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = plans.run_id
        AND   st.phase = 'plan'
        GROUP BY run_id, phase
    ) AS plan_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = applies.run_id
        AND   st.phase = 'apply'
        GROUP BY run_id, phase
    ) AS apply_status_timestamps,
    (
        SELECT array_agg(v.*) AS run_variables
        FROM run_variables v
        WHERE v.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_variables
FROM runs
JOIN plans USING (run_id)
JOIN applies USING (run_id)
JOIN (configuration_versions LEFT JOIN ingress_attributes ia USING (configuration_version_id)) USING (configuration_version_id)
JOIN workspaces ON runs.workspace_id = workspaces.workspace_id
JOIN organizations ON workspaces.organization_name = organizations.name
WHERE runs.run_id = $1
FOR UPDATE OF runs, plans, applies
;`

type FindRunByIDForUpdateRow struct {
	RunID                  pgtype.Text              `json:"run_id"`
	CreatedAt              pgtype.Timestamptz       `json:"created_at"`
	CancelSignaledAt       pgtype.Timestamptz       `json:"cancel_signaled_at"`
	IsDestroy              pgtype.Bool              `json:"is_destroy"`
	PositionInQueue        pgtype.Int4              `json:"position_in_queue"`
	Refresh                pgtype.Bool              `json:"refresh"`
	RefreshOnly            pgtype.Bool              `json:"refresh_only"`
	Source                 pgtype.Text              `json:"source"`
	Status                 pgtype.Text              `json:"status"`
	PlanStatus             pgtype.Text              `json:"plan_status"`
	ApplyStatus            pgtype.Text              `json:"apply_status"`
	ReplaceAddrs           []string                 `json:"replace_addrs"`
	TargetAddrs            []string                 `json:"target_addrs"`
	AutoApply              pgtype.Bool              `json:"auto_apply"`
	PlanResourceReport     *Report                  `json:"plan_resource_report"`
	PlanOutputReport       *Report                  `json:"plan_output_report"`
	ApplyResourceReport    *Report                  `json:"apply_resource_report"`
	ConfigurationVersionID pgtype.Text              `json:"configuration_version_id"`
	WorkspaceID            pgtype.Text              `json:"workspace_id"`
	PlanOnly               pgtype.Bool              `json:"plan_only"`
	CreatedBy              pgtype.Text              `json:"created_by"`
	TerraformVersion       pgtype.Text              `json:"terraform_version"`
	AllowEmptyApply        pgtype.Bool              `json:"allow_empty_apply"`
	ExecutionMode          pgtype.Text              `json:"execution_mode"`
	Latest                 pgtype.Bool              `json:"latest"`
	OrganizationName       pgtype.Text              `json:"organization_name"`
	CostEstimationEnabled  pgtype.Bool              `json:"cost_estimation_enabled"`
	IngressAttributes      *IngressAttributes       `json:"ingress_attributes"`
	RunStatusTimestamps    []*RunStatusTimestamps   `json:"run_status_timestamps"`
	PlanStatusTimestamps   []*PhaseStatusTimestamps `json:"plan_status_timestamps"`
	ApplyStatusTimestamps  []*PhaseStatusTimestamps `json:"apply_status_timestamps"`
	RunVariables           []*RunVariables          `json:"run_variables"`
}

// FindRunByIDForUpdate implements Querier.FindRunByIDForUpdate.
func (q *DBQuerier) FindRunByIDForUpdate(ctx context.Context, runID pgtype.Text) (FindRunByIDForUpdateRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunByIDForUpdate")
	rows, err := q.conn.Query(ctx, findRunByIDForUpdateSQL, runID)
	if err != nil {
		return FindRunByIDForUpdateRow{}, fmt.Errorf("query FindRunByIDForUpdate: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (FindRunByIDForUpdateRow, error) {
		var item FindRunByIDForUpdateRow
		if err := row.Scan(&item.RunID, // 'run_id', 'RunID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.CreatedAt,              // 'created_at', 'CreatedAt', 'pgtype.Timestamptz', 'github.com/jackc/pgx/v5/pgtype', 'Timestamptz'
			&item.CancelSignaledAt,       // 'cancel_signaled_at', 'CancelSignaledAt', 'pgtype.Timestamptz', 'github.com/jackc/pgx/v5/pgtype', 'Timestamptz'
			&item.IsDestroy,              // 'is_destroy', 'IsDestroy', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.PositionInQueue,        // 'position_in_queue', 'PositionInQueue', 'pgtype.Int4', 'github.com/jackc/pgx/v5/pgtype', 'Int4'
			&item.Refresh,                // 'refresh', 'Refresh', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.RefreshOnly,            // 'refresh_only', 'RefreshOnly', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.Source,                 // 'source', 'Source', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Status,                 // 'status', 'Status', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.PlanStatus,             // 'plan_status', 'PlanStatus', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.ApplyStatus,            // 'apply_status', 'ApplyStatus', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.ReplaceAddrs,           // 'replace_addrs', 'ReplaceAddrs', '[]string', '', '[]string'
			&item.TargetAddrs,            // 'target_addrs', 'TargetAddrs', '[]string', '', '[]string'
			&item.AutoApply,              // 'auto_apply', 'AutoApply', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.PlanResourceReport,     // 'plan_resource_report', 'PlanResourceReport', '*Report', '', '*Report'
			&item.PlanOutputReport,       // 'plan_output_report', 'PlanOutputReport', '*Report', '', '*Report'
			&item.ApplyResourceReport,    // 'apply_resource_report', 'ApplyResourceReport', '*Report', '', '*Report'
			&item.ConfigurationVersionID, // 'configuration_version_id', 'ConfigurationVersionID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.WorkspaceID,            // 'workspace_id', 'WorkspaceID', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.PlanOnly,               // 'plan_only', 'PlanOnly', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.CreatedBy,              // 'created_by', 'CreatedBy', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.TerraformVersion,       // 'terraform_version', 'TerraformVersion', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.AllowEmptyApply,        // 'allow_empty_apply', 'AllowEmptyApply', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.ExecutionMode,          // 'execution_mode', 'ExecutionMode', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.Latest,                 // 'latest', 'Latest', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.OrganizationName,       // 'organization_name', 'OrganizationName', 'pgtype.Text', 'github.com/jackc/pgx/v5/pgtype', 'Text'
			&item.CostEstimationEnabled,  // 'cost_estimation_enabled', 'CostEstimationEnabled', 'pgtype.Bool', 'github.com/jackc/pgx/v5/pgtype', 'Bool'
			&item.IngressAttributes,      // 'ingress_attributes', 'IngressAttributes', '*IngressAttributes', '', '*IngressAttributes'
			&item.RunStatusTimestamps,    // 'run_status_timestamps', 'RunStatusTimestamps', '[]*RunStatusTimestamps', '', '[]*RunStatusTimestamps'
			&item.PlanStatusTimestamps,   // 'plan_status_timestamps', 'PlanStatusTimestamps', '[]*PhaseStatusTimestamps', '', '[]*PhaseStatusTimestamps'
			&item.ApplyStatusTimestamps,  // 'apply_status_timestamps', 'ApplyStatusTimestamps', '[]*PhaseStatusTimestamps', '', '[]*PhaseStatusTimestamps'
			&item.RunVariables,           // 'run_variables', 'RunVariables', '[]*RunVariables', '', '[]*RunVariables'
		); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const putLockFileSQL = `UPDATE runs
SET lock_file = $1
WHERE run_id = $2
RETURNING run_id
;`

// PutLockFile implements Querier.PutLockFile.
func (q *DBQuerier) PutLockFile(ctx context.Context, lockFile []byte, runID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "PutLockFile")
	rows, err := q.conn.Query(ctx, putLockFileSQL, lockFile, runID)
	if err != nil {
		return pgtype.Text{}, fmt.Errorf("query PutLockFile: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (pgtype.Text, error) {
		var item pgtype.Text
		if err := row.Scan(&item); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const getLockFileByIDSQL = `SELECT lock_file
FROM runs
WHERE run_id = $1
;`

// GetLockFileByID implements Querier.GetLockFileByID.
func (q *DBQuerier) GetLockFileByID(ctx context.Context, runID pgtype.Text) ([]byte, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetLockFileByID")
	rows, err := q.conn.Query(ctx, getLockFileByIDSQL, runID)
	if err != nil {
		return nil, fmt.Errorf("query GetLockFileByID: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) ([]byte, error) {
		var item []byte
		if err := row.Scan(&item); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const updateRunStatusSQL = `UPDATE runs
SET
    status = $1
WHERE run_id = $2
RETURNING run_id
;`

// UpdateRunStatus implements Querier.UpdateRunStatus.
func (q *DBQuerier) UpdateRunStatus(ctx context.Context, status pgtype.Text, id pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateRunStatus")
	rows, err := q.conn.Query(ctx, updateRunStatusSQL, status, id)
	if err != nil {
		return pgtype.Text{}, fmt.Errorf("query UpdateRunStatus: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (pgtype.Text, error) {
		var item pgtype.Text
		if err := row.Scan(&item); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const updateCancelSignaledAtSQL = `UPDATE runs
SET
    cancel_signaled_at = $1
WHERE run_id = $2
RETURNING run_id
;`

// UpdateCancelSignaledAt implements Querier.UpdateCancelSignaledAt.
func (q *DBQuerier) UpdateCancelSignaledAt(ctx context.Context, cancelSignaledAt pgtype.Timestamptz, id pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateCancelSignaledAt")
	rows, err := q.conn.Query(ctx, updateCancelSignaledAtSQL, cancelSignaledAt, id)
	if err != nil {
		return pgtype.Text{}, fmt.Errorf("query UpdateCancelSignaledAt: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (pgtype.Text, error) {
		var item pgtype.Text
		if err := row.Scan(&item); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}

const deleteRunByIDSQL = `DELETE
FROM runs
WHERE run_id = $1
RETURNING run_id
;`

// DeleteRunByID implements Querier.DeleteRunByID.
func (q *DBQuerier) DeleteRunByID(ctx context.Context, runID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteRunByID")
	rows, err := q.conn.Query(ctx, deleteRunByIDSQL, runID)
	if err != nil {
		return pgtype.Text{}, fmt.Errorf("query DeleteRunByID: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (pgtype.Text, error) {
		var item pgtype.Text
		if err := row.Scan(&item); err != nil {
			return item, fmt.Errorf("failed to scan: %w", err)
		}
		return item, nil
	})
}
