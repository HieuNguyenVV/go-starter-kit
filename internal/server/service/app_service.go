package service

import (
	"context"
	"errors"
	"go-starter-kit/internal/log"
	"go-starter-kit/internal/server/model"
	"go-starter-kit/internal/server/repository"
)

type AppService interface {
	GetAppByID(ctx context.Context, appID string) (*model.App, error)
	CreateApp(ctx context.Context, appModel model.App) error
}
type appService struct {
	logger  log.Logger
	appRepo repository.AppRepository
}

func NewAppService(logger log.Logger, appRepo repository.AppRepository) AppService {
	return &appService{logger: logger, appRepo: appRepo}
}

func (s *appService) GetAppByID(ctx context.Context, appID string) (*model.App, error) {
	app, err := s.appRepo.GetByAppID(ctx, appID)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (s *appService) CreateApp(ctx context.Context, appModel model.App) error {
	app, err := s.appRepo.GetByAppID(ctx, appModel.ID)
	if err != nil {
		return err
	}
	if app != nil {
		return errors.New("app already exists")
	}

	err = s.appRepo.CreateApp(ctx, appModel)
	if err != nil {
		return err
	}
	return nil
}
