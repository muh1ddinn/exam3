package models

type Users struct {
	Id         string `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"Last_name"`
	Mail       string `json:"mail"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Login      string `json:"login"`
	Sex        string `json:"sex"`
	Active     bool   `json:"active"`
	Created_at string `json:"createdAt"`
	Updated_at string `json:"updatedAt"`
	Delete_at  string `json:"deleteAt"`
}

type GetAllusersResponse struct {
	Users []Users `json:"users"`
	Count int64   `json:"count"`
}

type GetAllUsersRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllUser struct {
	Id         string `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Mail       string `json:"mail"`
	Phone      string `json:"phone"`
	Is_blocked bool   `json:"isblocked"`
	Created_at string `json:"createdAt"`
	Updated_at string `json:"updatedAt"`
}

type Changepasswor struct {
	Id          string `json:"id"`
	Mail        string `json:"mail"`
	OldPassword string `json:"password"`
	NewPassword string `json:"newpassword"`
}

type Checklogin struct {
	Id    string `json:"id"`
	Phone string `json:"phone"`
	Mail  string `json:"mail"`
}

type UserRegisterRequest struct {
	Mail string `json:"mail`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLoginRequest struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type UserRegisterConf struct {
	Mail string `json:"mail"`
	Otp  string `json:"otp"`
	User Users  `json:"user"`
}

type Updatestatus struct {
	Active bool   `json:"active"`
	Id     string `json:"id"`
}
