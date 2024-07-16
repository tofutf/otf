package workspace

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/resource"
	"github.com/tofutf/tofutf/internal/sql"
	"github.com/tofutf/tofutf/internal/sql/pggen"
)

type (
	// pgdb is a workspace database on postgres
	pgdb struct {
		*sql.Pool // provides access to generated SQL queries
	}

	// pgresult represents the result of a database query for a workspace.
	pgresult struct {
		WorkspaceID                pgtype.Text            `json:"workspace_id"`
		CreatedAt                  pgtype.Timestamptz     `json:"created_at"`
		UpdatedAt                  pgtype.Timestamptz     `json:"updated_at"`
		AllowDestroyPlan           pgtype.Bool            `json:"allow_destroy_plan"`
		AutoApply                  pgtype.Bool            `json:"auto_apply"`
		CanQueueDestroyPlan        pgtype.Bool            `json:"can_queue_destroy_plan"`
		Description                pgtype.Text            `json:"description"`
		Environment                pgtype.Text            `json:"environment"`
		ExecutionMode              pgtype.Text            `json:"execution_mode"`
		GlobalRemoteState          pgtype.Bool            `json:"global_remote_state"`
		MigrationEnvironment       pgtype.Text            `json:"migration_environment"`
		Name                       pgtype.Text            `json:"name"`
		QueueAllRuns               pgtype.Bool            `json:"queue_all_runs"`
		SpeculativeEnabled         pgtype.Bool            `json:"speculative_enabled"`
		SourceName                 pgtype.Text            `json:"source_name"`
		SourceURL                  pgtype.Text            `json:"source_url"`
		StructuredRunOutputEnabled pgtype.Bool            `json:"structured_run_output_enabled"`
		TerraformVersion           pgtype.Text            `json:"terraform_version"`
		TriggerPrefixes            []string               `json:"trigger_prefixes"`
		WorkingDirectory           pgtype.Text            `json:"working_directory"`
		LockRunID                  pgtype.Text            `json:"lock_run_id"`
		LatestRunID                pgtype.Text            `json:"latest_run_id"`
		OrganizationName           pgtype.Text            `json:"organization_name"`
		Branch                     pgtype.Text            `json:"branch"`
		LockUsername               pgtype.Text            `json:"lock_username"`
		CurrentStateVersionID      pgtype.Text            `json:"current_state_version_id"`
		TriggerPatterns            []string               `json:"trigger_patterns"`
		VCSTagsRegex               pgtype.Text            `json:"vcs_tags_regex"`
		AllowCLIApply              pgtype.Bool            `json:"allow_cli_apply"`
		AgentPoolID                pgtype.Text            `json:"agent_pool_id"`
		Tags                       []string               `json:"tags"`
		LatestRunStatus            pgtype.Text            `json:"latest_run_status"`
		UserLock                   *pggen.Users           `json:"user_lock"`
		RunLock                    *pggen.Runs            `json:"run_lock"`
		WorkspaceConnection        *pggen.RepoConnections `json:"workspace_connection"`
	}
)

func (r pgresult) toWorkspace() (*Workspace, error) {
	ws := Workspace{
		ID:                         r.WorkspaceID.String,
		CreatedAt:                  r.CreatedAt.Time.UTC(),
		UpdatedAt:                  r.UpdatedAt.Time.UTC(),
		AllowDestroyPlan:           r.AllowDestroyPlan.Bool,
		AutoApply:                  r.AutoApply.Bool,
		CanQueueDestroyPlan:        r.CanQueueDestroyPlan.Bool,
		Description:                r.Description.String,
		Environment:                r.Environment.String,
		ExecutionMode:              ExecutionMode(r.ExecutionMode.String),
		GlobalRemoteState:          r.GlobalRemoteState.Bool,
		MigrationEnvironment:       r.MigrationEnvironment.String,
		Name:                       r.Name.String,
		QueueAllRuns:               r.QueueAllRuns.Bool,
		SpeculativeEnabled:         r.SpeculativeEnabled.Bool,
		StructuredRunOutputEnabled: r.StructuredRunOutputEnabled.Bool,
		SourceName:                 r.SourceName.String,
		SourceURL:                  r.SourceURL.String,
		TerraformVersion:           r.TerraformVersion.String,
		TriggerPrefixes:            r.TriggerPrefixes,
		TriggerPatterns:            r.TriggerPatterns,
		WorkingDirectory:           r.WorkingDirectory.String,
		Organization:               r.OrganizationName.String,
		Tags:                       r.Tags,
	}
	if r.AgentPoolID.Valid {
		ws.AgentPoolID = &r.AgentPoolID.String
	}

	if r.WorkspaceConnection != nil {
		ws.Connection = &Connection{
			AllowCLIApply: r.AllowCLIApply.Bool,
			VCSProviderID: r.WorkspaceConnection.VCSProviderID.String,
			Repo:          r.WorkspaceConnection.RepoPath.String,
			Branch:        r.Branch.String,
		}
		if r.VCSTagsRegex.Valid {
			ws.Connection.TagsRegex = r.VCSTagsRegex.String
		}
	}

	if r.LatestRunID.Valid && r.LatestRunStatus.Valid {
		ws.LatestRun = &LatestRun{
			ID:     r.LatestRunID.String,
			Status: runStatus(r.LatestRunStatus.String),
		}
	}

	if r.UserLock != nil {
		ws.Lock = &Lock{
			id:       r.UserLock.Username.String,
			LockKind: UserLock,
		}
	} else if r.RunLock != nil && r.RunLock.RunID.Valid {
		ws.Lock = &Lock{
			id:       r.RunLock.RunID.String,
			LockKind: RunLock,
		}
	}

	return &ws, nil
}

func (db *pgdb) create(ctx context.Context, ws *Workspace) error {
	err := db.Query(ctx, func(ctx context.Context, q pggen.Querier) error {
		params := pggen.InsertWorkspaceParams{
			ID:                         sql.String(ws.ID),
			CreatedAt:                  sql.Timestamptz(ws.CreatedAt),
			UpdatedAt:                  sql.Timestamptz(ws.UpdatedAt),
			AgentPoolID:                sql.StringPtr(ws.AgentPoolID),
			AllowCLIApply:              sql.Bool(false),
			AllowDestroyPlan:           sql.Bool(ws.AllowDestroyPlan),
			AutoApply:                  sql.Bool(ws.AutoApply),
			Branch:                     sql.String(""),
			CanQueueDestroyPlan:        sql.Bool(ws.CanQueueDestroyPlan),
			Description:                sql.String(ws.Description),
			Environment:                sql.String(ws.Environment),
			ExecutionMode:              sql.String(string(ws.ExecutionMode)),
			GlobalRemoteState:          sql.Bool(ws.GlobalRemoteState),
			MigrationEnvironment:       sql.String(ws.MigrationEnvironment),
			Name:                       sql.String(ws.Name),
			QueueAllRuns:               sql.Bool(ws.QueueAllRuns),
			SpeculativeEnabled:         sql.Bool(ws.SpeculativeEnabled),
			SourceName:                 sql.String(ws.SourceName),
			SourceURL:                  sql.String(ws.SourceURL),
			StructuredRunOutputEnabled: sql.Bool(ws.StructuredRunOutputEnabled),
			TerraformVersion:           sql.String(ws.TerraformVersion),
			TriggerPrefixes:            ws.TriggerPrefixes,
			TriggerPatterns:            ws.TriggerPatterns,
			VCSTagsRegex:               sql.StringPtr(nil),
			WorkingDirectory:           sql.String(ws.WorkingDirectory),
			OrganizationName:           sql.String(ws.Organization),
		}
		if ws.Connection != nil {
			params.AllowCLIApply = sql.Bool(ws.Connection.AllowCLIApply)
			params.Branch = sql.String(ws.Connection.Branch)
			params.VCSTagsRegex = sql.String(ws.Connection.TagsRegex)
		}
		_, err := q.InsertWorkspace(ctx, params)
		return sql.Error(err)
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *pgdb) update(ctx context.Context, workspaceID string, fn func(*Workspace) error) (*Workspace, error) {
	ws, err := sql.Tx(ctx, db.Pool, func(ctx context.Context, q pggen.Querier) (*Workspace, error) {
		var err error
		// retrieve workspace
		result, err := q.FindWorkspaceByIDForUpdate(ctx, sql.String(workspaceID))
		if err != nil {
			return nil, sql.Error(err)
		}
		ws, err := pgresult(result).toWorkspace()
		if err != nil {
			return nil, err
		}

		// update workspace
		if err := fn(ws); err != nil {
			return nil, err
		}
		// persist update
		params := pggen.UpdateWorkspaceByIDParams{
			AgentPoolID:                sql.StringPtr(ws.AgentPoolID),
			AllowDestroyPlan:           sql.Bool(ws.AllowDestroyPlan),
			AllowCLIApply:              sql.Bool(false),
			AutoApply:                  sql.Bool(ws.AutoApply),
			Branch:                     sql.String(""),
			Description:                sql.String(ws.Description),
			ExecutionMode:              sql.String(string(ws.ExecutionMode)),
			GlobalRemoteState:          sql.Bool(ws.GlobalRemoteState),
			Name:                       sql.String(ws.Name),
			QueueAllRuns:               sql.Bool(ws.QueueAllRuns),
			SpeculativeEnabled:         sql.Bool(ws.SpeculativeEnabled),
			StructuredRunOutputEnabled: sql.Bool(ws.StructuredRunOutputEnabled),
			TerraformVersion:           sql.String(ws.TerraformVersion),
			TriggerPrefixes:            ws.TriggerPrefixes,
			TriggerPatterns:            ws.TriggerPatterns,
			VCSTagsRegex:               sql.StringPtr(nil),
			WorkingDirectory:           sql.String(ws.WorkingDirectory),
			UpdatedAt:                  sql.Timestamptz(ws.UpdatedAt),
			ID:                         sql.String(ws.ID),
		}
		if ws.Connection != nil {
			params.AllowCLIApply = sql.Bool(ws.Connection.AllowCLIApply)
			params.Branch = sql.String(ws.Connection.Branch)
			params.VCSTagsRegex = sql.String(ws.Connection.TagsRegex)
		}
		_, err = q.UpdateWorkspaceByID(ctx, params)
		return ws, err
	})
	if err != nil {
		return nil, err
	}

	return ws, nil
}

// setCurrentRun sets the ID of the current run for the specified workspace.
func (db *pgdb) setCurrentRun(ctx context.Context, workspaceID, runID string) (*Workspace, error) {
	ws, err := sql.Query(ctx, db.Pool, func(ctx context.Context, q pggen.Querier) (*Workspace, error) {
		_, err := q.UpdateWorkspaceLatestRun(ctx, sql.String(runID), sql.String(workspaceID))
		if err != nil {
			return nil, sql.Error(err)
		}

		return db.get(ctx, workspaceID)
	})
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return ws, nil
}

func (db *pgdb) list(ctx context.Context, opts ListOptions) (*resource.Page[*Workspace], error) {
	ws, err := sql.Query(ctx, db.Pool, func(ctx context.Context, q pggen.Querier) (*resource.Page[*Workspace], error) {

		// Organization name filter is optional - if not provided use a % which in
		// SQL means match any organization.
		organization := "%"
		if opts.Organization != nil {
			organization = *opts.Organization
		}
		tags := []string{}
		if len(opts.Tags) > 0 {
			tags = opts.Tags
		}

		rows, err := q.FindWorkspaces(ctx, pggen.FindWorkspacesParams{
			OrganizationNames: []string{organization},
			Search:            sql.String(opts.Search),
			Tags:              tags,
			Limit:             opts.GetLimit(),
			Offset:            opts.GetOffset(),
		})
		if err != nil {
			return nil, err
		}
		count, err := q.CountWorkspaces(ctx, pggen.CountWorkspacesParams{
			Search:            sql.String(opts.Search),
			OrganizationNames: []string{organization},
			Tags:              tags,
		})
		if err != nil {
			return nil, err
		}

		items := make([]*Workspace, len(rows))
		for i, r := range rows {
			ws, err := pgresult(r).toWorkspace()
			if err != nil {
				return nil, err
			}
			items[i] = ws
		}

		return resource.NewPage(items, opts.PageOptions, internal.Int64(count.Int64)), nil
	})
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (db *pgdb) listByConnection(ctx context.Context, vcsProviderID, repoPath string) ([]*Workspace, error) {
	ws, err := sql.Query(ctx, db.Pool, func(ctx context.Context, q pggen.Querier) ([]*Workspace, error) {
		rows, err := q.FindWorkspacesByConnection(ctx, sql.String(vcsProviderID), sql.String(repoPath))
		if err != nil {
			return nil, err
		}

		items := make([]*Workspace, len(rows))
		for i, r := range rows {
			ws, err := pgresult(r).toWorkspace()
			if err != nil {
				return nil, err
			}
			items[i] = ws
		}

		return items, nil
	})
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (db *pgdb) listByUsername(ctx context.Context, username string, organization string, opts resource.PageOptions) (*resource.Page[*Workspace], error) {
	rp, err := sql.Query(ctx, db.Pool, func(ctx context.Context, q pggen.Querier) (*resource.Page[*Workspace], error) {
		rows, err := q.FindWorkspacesByUsername(ctx, pggen.FindWorkspacesByUsernameParams{
			OrganizationName: sql.String(organization),
			Username:         sql.String(username),
			Limit:            opts.GetLimit(),
			Offset:           opts.GetOffset(),
		})
		if err != nil {
			return nil, err
		}

		count, err := q.CountWorkspacesByUsername(ctx, sql.String(organization), sql.String(username))
		if err != nil {
			return nil, err
		}

		items := make([]*Workspace, len(rows))
		for i, r := range rows {
			ws, err := pgresult(r).toWorkspace()
			if err != nil {
				return nil, err
			}
			items[i] = ws
		}

		return resource.NewPage(items, opts, internal.Int64(count.Int64)), nil
	})
	if err != nil {
		return nil, err
	}

	return rp, nil
}

func (db *pgdb) get(ctx context.Context, workspaceID string) (*Workspace, error) {
	ws, err := sql.Query(ctx, db.Pool, func(ctx context.Context, q pggen.Querier) (*Workspace, error) {
		result, err := q.FindWorkspaceByID(ctx, sql.String(workspaceID))
		if err != nil {
			return nil, sql.Error(err)
		}

		ws, err := pgresult(result).toWorkspace()
		if err != nil {
			return nil, err
		}

		return ws, nil
	})
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (db *pgdb) getByName(ctx context.Context, organization, workspace string) (*Workspace, error) {
	ws, err := sql.Query(ctx, db.Pool, func(ctx context.Context, q pggen.Querier) (*Workspace, error) {
		result, err := q.FindWorkspaceByName(ctx, sql.String(workspace), sql.String(organization))
		if err != nil {
			return nil, sql.Error(err)
		}

		ws, err := pgresult(result).toWorkspace()
		if err != nil {
			return nil, err
		}

		return ws, nil
	})
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (db *pgdb) delete(ctx context.Context, workspaceID string) error {
	err := db.Query(ctx, func(ctx context.Context, q pggen.Querier) error {
		_, err := q.DeleteWorkspaceByID(ctx, sql.String(workspaceID))
		if err != nil {
			return sql.Error(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
