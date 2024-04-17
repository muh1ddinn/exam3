package service

import (
	"context"
	"errors"
	models "exam3/api/model"
	"exam3/config"
	"exam3/pkg"
	"exam3/pkg/jwt"
	"exam3/pkg/logger"
	"exam3/pkg/password"
	smtp "exam3/pkg/stmp"
	"exam3/storage"
	"fmt"
	"time"
)

type authService struct {
	storage storage.IStorage
	log     logger.ILogger
	redis   storage.IREdisStorage
}

func NewAuthService(storage storage.IStorage, log logger.ILogger, redis storage.IREdisStorage) authService {
	return authService{
		storage: storage,
		log:     log,
		redis:   redis,
	}
}

func (a authService) UserLOgin(ctx context.Context, loginRequest models.UserLoginRequest) (models.UserLoginResponse, error) {
	fmt.Println(" loginRequest.Login: ", loginRequest.Mail)
	user, err := a.storage.Users().GetByID(ctx, loginRequest.Mail)
	if err != nil {
		a.log.Error("error while getting user credentials by login", logger.Error(err))
		return models.UserLoginResponse{}, err
	}

	if err = password.CompareHashAndPassword(user.Password, loginRequest.Password); err != nil {
		a.log.Error("error while comparing password", logger.Error(err))
		return models.UserLoginResponse{}, err
	}

	m := make(map[interface{}]interface{})

	m["user_id"] = user.Id
	m["user_role"] = config.USER_ROLE

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for users login", logger.Error(err))
		return models.UserLoginResponse{}, err
	}

	return models.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a authService) UsersRegister(ctx context.Context, loginRequest models.UserRegisterRequest) error {
	fmt.Println(" loginRequest.Login: ", loginRequest.Mail)

	otpCode := pkg.GenerateOTP()

	msg := fmt.Sprintf("Your otp code is: %v, for registering . Don't give it to anyone", otpCode)

	err := a.redis.SetX(ctx, loginRequest.Mail, otpCode, time.Minute*2)
	if err != nil {

		a.log.Error("error while setting otpCode to redis users register", logger.Error(err))
		return err
	}

	err = smtp.Sendmail(loginRequest.Mail, msg)
	if err != nil {
		a.log.Error("error while sending otp code to users register", logger.Error(err))
		return err
	}
	fmt.Println(msg)

	return nil
}

func (a authService) UsersRegisterConfirm(ctx context.Context, req models.UserRegisterConf) (models.UserLoginResponse, error) {
	resp := models.UserLoginResponse{}

	otp, err := a.redis.Get(ctx, req.Mail)
	if err != nil {
		a.log.Error("error while getting otp code for users register confirm", logger.Error(err))
		return resp, err
	}
	if req.Otp != otp {

		a.log.Error("incorrect otp code for users register confirm", logger.Error(err))
		return resp, errors.New("incorrect otp code")
	}
	fmt.Println(otp)

	req.User.Mail = req.Mail
	id, err := a.storage.Users().Createconf(ctx, req.User)
	if err != nil {
		a.log.Error("error while creating users", logger.Error(err))
		return resp, err
	}
	var m = make(map[interface{}]interface{})

	m["user_id"] = id
	m["user_role"] = config.USER_ROLE

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for users register confirm", logger.Error(err))
		return resp, err
	}
	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken
	fmt.Println(accessToken)

	return resp, nil
}
func (a authService) Usersend(ctx context.Context, loginRequest models.UserLoginRequest) (models.UserLoginResponse, error) {
	storedPassword, err := a.redis.Get(ctx, loginRequest.Mail)
	if err != nil {
		fmt.Println(storedPassword,"000000000")
		a.log.Error("error while getting password from Redis", logger.Error(err))
		return models.UserLoginResponse{}, err
	}

	if storedPassword != loginRequest.Password {
		a.log.Error("password does not match")
		return models.UserLoginResponse{}, errors.New("incorrect password")
	}

	m := make(map[interface{}]interface{})
	m["user_id"] = loginRequest.Mail
	m["user_role"] = config.USER_ROLE

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for users login", logger.Error(err))
		return models.UserLoginResponse{}, err
	}

	return models.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
