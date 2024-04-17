package service

import (
	"exam3/pkg/logger"
	"exam3/storage"
)

type IServiceMangaer interface {
	User() userser
	Auth() authService
}

type Service struct {
	customerSer userser
	logger      logger.ILogger
	authService authService
}

func New(storage storage.IStorage, log logger.ILogger, redis storage.IREdisStorage) Service {
	return Service{

		customerSer: NewUserSeer(storage, log),
		logger:      log,
		authService: NewAuthService(storage, log, redis),

	}
}

func (s Service) User() userser {
	return s.customerSer
}

func (s Service) Auth() authService {
	return s.authService

}
