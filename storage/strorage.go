package storage

import (
	"context"
	models "exam3/api/model"
	"time"
)

type IStorage interface {
	CloseDB()
	Users() IUsersStorage
}

type IUsersStorage interface {
	Create(context.Context, models.Users) (string, error)
	GetAlluser(context.Context, models.GetAllUsersRequest) (models.GetAllusersResponse, error)
	Update(context.Context, models.Users) (string, error)
	Delete(context.Context, string) (string, error)
	GetByID(context.Context, string) (models.Users, error)
	Login(context.Context, models.Changepasswor) (string, error)
	GetPassword(ctx context.Context, phone string) (string, error)
	ChangePassword(ctx context.Context, pass models.Changepasswor) (string, error)
	GetByMail(context.Context, string) (models.Users, error)
	Checklogin(context.Context, string) (models.UserRegisterRequest, error)
	Createconf(context.Context, models.Users) (string, error)
	Updatestatus(context.Context, models.Updatestatus) (string, error)
}

type IREdisStorage interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Del(ctx context.Context, key string) error
}
