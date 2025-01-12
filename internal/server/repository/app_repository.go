package repository

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"go-starter-kit/internal/pkg/database"
	"go-starter-kit/internal/server/model"
)

type app struct {
	ID        string `db:"id"`
	OrgID     string `db:"org_id"`
	Name      string `db:"name"`
	APIToken  string `db:"api_token"`
	CreateAt  int64  `db:"created_at"`
	UpdateAt  int64  `db:"updated_at"`
	IsRemoved bool   `db:"is_removed"`
}

func (a app) toEntity() *model.App {
	return &model.App{
		ID:        a.ID,
		Name:      a.Name,
		OrgID:     a.OrgID,
		APIToken:  a.APIToken,
		CreatedAt: a.CreateAt,
		UpdatedAt: a.UpdateAt,
		IsRemoved: a.IsRemoved,
	}
}

func (a *app) FromEntity(e model.App) {
	a.ID = e.ID
	a.Name = e.Name
	a.OrgID = e.OrgID
	a.APIToken = e.APIToken
	a.CreateAt = e.CreatedAt
	a.UpdateAt = e.UpdatedAt
	a.IsRemoved = e.IsRemoved
}

type AppRepository interface {
	CreateApp(ctx context.Context, app model.App) error
	GetByAppID(ctx context.Context, appID string) (*model.App, error)
}

type appRepository struct {
	postgres *database.Postgres
}

func NewAppRepository(postgres *database.Postgres) AppRepository {
	return &appRepository{
		postgres: postgres,
	}
}

func (r *appRepository) CreateApp(ctx context.Context, appModel model.App) error {
	conn, err := r.postgres.GetWriteConnection(ctx)
	if err != nil {
		return errors.Wrap(err, "get write connection failed")
	}

	query := `
		INSERT INTO apps (
		    id,
		    org_id,
		    name,
		    api_token,
		    created_at,
		    updated_at,
		    is_removed
		) VALUES (
		    :id,
		    :org_id,
			:name,
		    :api_token,
		    :created_at,
		    :updated_at,
		    :is_removed
		)`
	app := new(app)
	app.FromEntity(appModel)

	_, err = conn.NamedExec(query, &app)
	if err != nil {
		return errors.Wrap(err, "create app failed")
	}
	return nil
}

func (r *appRepository) GetByAppID(ctx context.Context, appID string) (*model.App, error) {
	conn, err := r.postgres.GetReadConnection(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get read connection failed")
	}
	stmt := conn.Rebind(
		`SELECT
            a.id,
            a.org_id,
            a.name,
            a.is_removed,
            a.api_token,
            a.created_at,
            a.updated_at
        FROM apps a
        WHERE a.id = ?
        ORDER BY a.created_at DESC
        LIMIT 1`)

	var app app
	if err := conn.Get(&app, stmt, appID); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "get app failed")
	}
	return app.toEntity(), nil
}
