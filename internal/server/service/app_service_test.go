package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go-starter-kit/internal/log"
	"go-starter-kit/internal/server/config"
	"go-starter-kit/internal/server/model"
	"go-starter-kit/internal/server/repository"
	mocks "go-starter-kit/mocks/repository"
	"testing"
	"time"
)

type appServiceTestSuite struct {
	suite.Suite
	service     AppService
	mockAppRepo *mocks.MockAppRepository
}

func newTestAppService(paths []string, appRepo repository.AppRepository) (AppService, error) {
	conf, err := config.NewExampleConfig(paths)
	if err != nil {
		return nil, err
	}

	logger, err := log.NewLogger(conf)
	if err != nil {
		return nil, err
	}

	return &appService{logger: logger, appRepo: appRepo}, nil
}

func (s *appServiceTestSuite) SetupSuite() {
	s.mockAppRepo = mocks.NewMockAppRepository(gomock.NewController(s.T()))
	service, err := newTestAppService([]string{"../../.."}, s.mockAppRepo)
	s.Require().NoError(err)
	s.service = service
}

func (s *appServiceTestSuite) TearDownSuite() {
}

func TestAppServiceTestSuite(t *testing.T) {
	suite.Run(t, new(appServiceTestSuite))
}

func (s *appServiceTestSuite) TestGetAppByID() {
	now := time.Now().UnixMilli()
	type args struct {
		appID string
	}
	tests := []struct {
		name        string
		args        args
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "successful get app",
			args: args{
				appID: "550e8400-e29b-41d4-a716-446655440000",
			},
			mockSetup: func() {
				s.mockAppRepo.EXPECT().GetByAppID(context.Background(), "550e8400-e29b-41d4-a716-446655440000").Return(&model.App{
					ID:        "550e8400-e29b-41d4-a716-446655440000",
					Name:      "test",
					OrgID:     "1234",
					APIToken:  "1234",
					CreatedAt: now,
					UpdatedAt: now,
					IsRemoved: true,
				}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "app not found",
			args: args{
				appID: "550e8400-e29b-41d4-a716-446655440000",
			},
			mockSetup: func() {
				s.mockAppRepo.EXPECT().GetByAppID(context.Background(), gomock.Any()).Return(nil, errors.New("app not found"))
			},
			expectedErr: errors.New("app not found"),
		},
		{
			name: "database error",
			args: args{
				appID: "550e8400-e29b-41d4-a716-446655440000",
			},
			mockSetup: func() {
				s.mockAppRepo.EXPECT().GetByAppID(context.Background(), gomock.Any()).Return(nil, errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			app, err := s.service.GetAppByID(context.Background(), tt.args.appID)
			if tt.expectedErr != nil {
				s.Require().Error(err)
				s.Equal(tt.expectedErr, err)
			} else {
				s.Require().NoError(err)
				s.Equal(tt.args.appID, app.ID)
			}
		})
	}
}

func (s *appServiceTestSuite) TestCreateApp() {
	now := time.Now().UnixMilli()
	type args struct {
		app model.App
	}
	tests := []struct {
		name        string
		args        args
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "successful create app",
			args: args{
				app: model.App{
					ID:        "550e8400-e29b-41d4-a716-446655440000",
					Name:      "test",
					OrgID:     "1234",
					APIToken:  "1234",
					CreatedAt: now,
					UpdatedAt: now,
					IsRemoved: true,
				},
			},
			mockSetup: func() {
				s.mockAppRepo.EXPECT().GetByAppID(context.Background(), gomock.Any()).Return(nil, nil)
				s.mockAppRepo.EXPECT().CreateApp(context.Background(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "get app error",
			args: args{
				app: model.App{
					ID: "123",
				},
			},
			mockSetup: func() {
				s.mockAppRepo.EXPECT().GetByAppID(context.Background(), gomock.Any()).Return(nil, errors.New("error database error"))
			},
			expectedErr: errors.New("error database error"),
		},
		{
			name: "app already exists",
			args: args{
				app: model.App{
					ID: "123",
				},
			},
			mockSetup: func() {
				s.mockAppRepo.EXPECT().GetByAppID(context.Background(), gomock.Any()).Return(&model.App{}, nil)
			},
			expectedErr: errors.New("app already exists"),
		},
		{
			name: "create app error",
			args: args{
				app: model.App{
					ID: "123",
				},
			},
			mockSetup: func() {
				s.mockAppRepo.EXPECT().GetByAppID(context.Background(), gomock.Any()).Return(nil, nil)
				s.mockAppRepo.EXPECT().CreateApp(context.Background(), gomock.Any()).Return(errors.New("error database error"))
			},
			expectedErr: errors.New("error database error"),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := s.service.CreateApp(context.Background(), tt.args.app)
			if tt.expectedErr != nil {
				s.Require().Error(err)
				s.Equal(tt.expectedErr, err)
			} else {
				s.Require().NoError(err)
			}

		})
	}
}
