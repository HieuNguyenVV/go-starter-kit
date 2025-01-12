package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	log2 "go-starter-kit/internal/log"
	"go-starter-kit/internal/pkg/database"
	"go-starter-kit/internal/server/config"
	"go-starter-kit/internal/server/model"
	"testing"
	"time"
)

// gen: mockgen -source=./internal/server/repository/app_repository.go -destination=./mocks/repository/mock_app_repository.go -package=mocks

type appRepositoryTestSuite struct {
	suite.Suite
	pg   *database.Postgres
	repo AppRepository
}

func TestAppRepository(t *testing.T) {
	suite.Run(t, new(appRepositoryTestSuite))
}

func newTestAppRepository(paths []string) (AppRepository, *database.Postgres, error) {
	conf, err := config.NewExampleConfig(paths)
	if err != nil {
		return nil, nil, err
	}
	logger, err := log2.NewLogger(conf)
	if err != nil {
		return nil, nil, err
	}
	db, err := database.NewPostgres(conf, logger)
	if err != nil {
		return nil, nil, err
	}
	return &appRepository{postgres: db}, db, nil
}

func (s *appRepositoryTestSuite) SetupSuite() {
	repo, db, err := newTestAppRepository([]string{"../../.."})
	s.Require().NoError(err)
	s.repo = repo
	s.pg = db
}

func (s *appRepositoryTestSuite) TearDownSuite() {
	s.pg.GetMasterDB().MustExec("TRUNCATE TABLE apps CASCADE")
}

func (s *appRepositoryTestSuite) TestGetByAppID() {
	now := time.Now().UnixMilli()
	type args struct {
		appID string
	}

	tests := []struct {
		name      string
		args      args
		want      *model.App
		wantError bool
	}{
		{
			name: "successful get app",
			args: args{
				appID: "550e8400-e29b-41d4-a716-446655440000",
			},
			want: &model.App{
				ID:        "550e8400-e29b-41d4-a716-446655440000",
				Name:      "test",
				OrgID:     "1234",
				APIToken:  "1234",
				CreatedAt: now,
				UpdatedAt: now,
				IsRemoved: true,
			},
			wantError: false,
		},
		{
			name: "app not found",
			args: args{
				appID: "550e8400-e29b-41d4-a716-446655440001",
			},
			want:      nil,
			wantError: true,
		},
	}

	s.pg.GetMasterDB().MustExec(`
		INSERT INTO apps (
					id,
					org_id,
					name,
					api_token,
					created_at,
					updated_at,
					is_removed
				) VALUES (
					'550e8400-e29b-41d4-a716-446655440000',
					'1234',
					'test',
					'1234',
					$1,
					$1,
					true
		)`, now)

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.GetByAppID(context.Background(), tt.args.appID)
			if tt.wantError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tt.want, got)
			}
		})
	}
	s.pg.GetMasterDB().MustExec(`TRUNCATE TABLE apps CASCADE`)
}

func (s *appRepositoryTestSuite) TestCreateApp() {
	now := time.Now().UnixMilli()
	id := uuid.New().String()
	type args struct {
		app model.App
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
		count     int
	}{
		{
			name: "successfully create app",
			args: args{
				app: model.App{
					ID:        id,
					Name:      "test",
					OrgID:     "123",
					APIToken:  "123",
					CreatedAt: now,
					UpdatedAt: now,
					IsRemoved: false,
				},
			},
			wantError: false,
			count:     1,
		},
		{
			name: "duplicate app row",
			args: args{
				app: model.App{
					ID:        id,
					Name:      "test",
					OrgID:     "123",
					APIToken:  "123",
					CreatedAt: now,
					UpdatedAt: now,
					IsRemoved: false,
				},
			},
			wantError: true,
			count:     0,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repo.CreateApp(context.Background(), tt.args.app)
			if tt.wantError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				sql := "SELECT COUNT(*) FROM apps"
				var count int
				err = s.pg.GetMasterDB().GetContext(context.Background(), &count, sql)
				s.Require().NoError(err)
				s.Require().Equal(tt.count, count)
				if tt.count == 1 {
					var app app
					err = s.pg.GetMasterDB().Get(&app, "SELECT * FROM apps WHERE id = $1", id)
					s.Require().NoError(err)
					appModel := model.App{
						ID:        id,
						Name:      app.Name,
						OrgID:     app.OrgID,
						APIToken:  app.APIToken,
						CreatedAt: app.CreateAt,
						UpdatedAt: app.UpdateAt,
						IsRemoved: app.IsRemoved,
					}
					s.Require().Equal(tt.args.app, appModel)
				}
			}
		})
	}
}
