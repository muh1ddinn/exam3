package service

import (
	"context"
	models "exam3/api/model"
	"exam3/pkg/logger"
	"exam3/pkg/password"
	"exam3/storage"
	"fmt"
)

type userser struct {
	storage storage.IStorage
	logger  logger.ILogger
}

func NewUserSeer(storage storage.IStorage, logger logger.ILogger) userser {

	return userser{
		storage: storage,
		logger:  logger,
	}
}
func (u userser) CreateCus(ctx context.Context, users models.Users) (string, error) {

	pKey, err := u.storage.Users().Create(ctx, users)
	if err != nil {
		u.logger.Error("ERROR in service layer while creating users", logger.Error(err))
		return "", err
	}

	return pKey, nil
}

func (s userser) GetByID(ctx context.Context, id string) (models.Users, error) {
	users, err := s.storage.Users().GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get users by ID", logger.Error(err))
		return users, err
	}

	return users, nil
}

func (s userser) GetAllCus(ctx context.Context, userss models.GetAllUsersRequest) ([]models.Users, error) {
	users, err := s.storage.Users().GetAlluser(ctx, userss)
	if err != nil {
		s.logger.Error("error in service layer while getting allcars: ", logger.Error(err))
		return nil, err
	}
	return users.Users, nil
}

func (s userser) UpdateUSER(ctx context.Context, userss models.Users) (string, error) {
	usersid, err := s.storage.Users().Update(ctx, userss)
	if err != nil {
		s.logger.Error("error in service layer while getting allcars: ", logger.Error(err))
		return usersid, err
	}
	return usersid, nil
}

func (s userser) Deletuser(ctx context.Context, id string) (string, error) {
	usersid, err := s.storage.Users().Delete(ctx, id)
	if err != nil {
		s.logger.Error("error in service layer while getting allcars: ", logger.Error(err))
		return usersid, err
	}
	return usersid, nil
}

func (s userser) Login(ctx context.Context, req models.Changepasswor) (string, error) {

	hashedPswd, err := s.storage.Users().GetPassword(ctx, req.Mail)
	if err != nil {
		s.logger.Error("error while getting users password", logger.Error(err))
		return "", err
	}

	err = password.CompareHashAndPassword(hashedPswd, req.OldPassword)
	if err != nil {
		s.logger.Error("incorrect password", logger.Error(err))
		return "", err
	}
	return "Login successfully", nil
}

func (s userser) ChangePassword(ctx context.Context, pass models.Changepasswor) (string, error) {
	msg, err := s.storage.Users().ChangePassword(ctx, pass)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("failed to change users password", logger.Error(err))
		return "", err
	}
	return msg, nil
}

func (s userser) Updatestatus(ctx context.Context, pass models.Updatestatus) (string, error) {
	msg, err := s.storage.Users().Updatestatus(ctx, pass)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("failed to change users password", logger.Error(err))
		return "", err
	}
	return msg, nil
}
