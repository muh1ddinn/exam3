package handler

import (
	models "exam3/api/model"
	"exam3/pkg/check"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UsersLogin godoc
// @Router       /users/login [POST]
// @Summary      users login
// @Description  users login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.UserLoginRequest true "login"
// @Success      201 {object} models.UserLoginResponse
// @Failure      400 {object} models.Responsee
// @Failure      404 {object} models.Responsee
// @Failure      500 {object} models.Responsee
func (h *Handler) UsersLogin(c *gin.Context) {

	loginReq := models.UserLoginRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handleResponseLog(c, h.Log, "error while binding body", http.StatusBadRequest, err)
		return
	}
	fmt.Println("loginReq:", loginReq)

	if err := check.ValidatePassword(loginReq.Password); err != nil {

		handleResponseLog(c, h.Log, "error whilr decoding request body", http.StatusBadRequest, err.Error())

		return

	}

	loginResp, err := h.Services.Auth().UserLOgin(c.Request.Context(), loginReq)

	if err != nil {
		handleResponseLog(c, h.Log, "unauthorized", http.StatusUnauthorized, err)
		return
	}

	handleResponseLog(c, h.Log, "Succes", http.StatusOK, loginResp)

}

// UsersRegister godoc
// @Router       /users/register [POST]
// @Summary      users register
// @Description  users register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.UserRegisterRequest true "register"
// @Success      200 {object} models.Responsee
// @Failure      400 {object} models.Responsee
// @Failure      404 {object} models.Responsee
// @Failure      500 {object} models.Responsee
func (h *Handler) UserRegister(c *gin.Context) {
	loginReq := models.UserRegisterRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handleResponseLog(c, h.Log, "error while binding body", http.StatusBadRequest, err)
		return
	}

	if err := check.Validategmail(loginReq.Mail); err != nil {

		handleResponseLog(c, h.Log, "error while putting mail", http.StatusBadRequest, err)

		fmt.Println("loginReq: ", loginReq)

		return

	}

	err := h.Services.Auth().UsersRegister(c.Request.Context(), loginReq)

	if err != nil {
		fmt.Println(err, "")
		handleResponseLog(c, h.Log, "", http.StatusInternalServerError, err)
		return
	}

	handleResponseLog(c, h.Log, "Otp sent successfull", http.StatusOK, "")

}

// UserRegister godoc
// @Router       /users/register-confirm [POST]
// @Summary      users register
// @Description  users register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.UserRegisterConf true "register"
// @Success      201  {object}  models.UserLoginResponse
// @Failure      400  {object}  models.Responsee
// @Failure      404  {object}  models.Responsee
// @Failure      500  {object}  models.Responsee
func (h *Handler) UsersRegisterConfirm(c *gin.Context) {
	req := models.UserRegisterConf{}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponseLog(c, h.Log, "error while binding body", http.StatusBadRequest, err)
		return
	}
	fmt.Println("req: ", req)

	//TODO: need validate login & password

	confResp, err := h.Services.Auth().UsersRegisterConfirm(c.Request.Context(), req)
	if err != nil {
		fmt.Println(confResp, "eeeeee")
		handleResponseLog(c, h.Log, "error while confirming", http.StatusUnauthorized, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Succes", http.StatusOK, confResp)

}

// UsersLogin godoc
// @Router       /users/login_with_otp [POST]
// @Summary      users loginwith_otp
// @Description  users loginwith_otp
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.UserLoginRequest true "login"
// @Success      201 {object} models.UserLoginResponse
// @Failure      400 {object} models.Responsee
// @Failure      404 {object} models.Responsee
// @Failure      500 {object} models.Responsee
func (h *Handler) UsersLoginwithotp(c *gin.Context) {

	loginReq := models.UserLoginRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handleResponseLog(c, h.Log, "error while binding body", http.StatusBadRequest, err)
		return
	}
	fmt.Println("loginReq:", loginReq)

	loginResp, err := h.Services.Auth().Usersend(c.Request.Context(), loginReq)

	if err != nil {
		handleResponseLog(c, h.Log, "unauthorized", http.StatusUnauthorized, err)
		return
	}

	handleResponseLog(c, h.Log, "Succes", http.StatusOK, loginResp)

}

// UsersRegister godoc
// @Router       /users/login_for_otp [POST]
// @Summary      users sendloginotp
// @Description  users sendloginotp
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.UserRegisterRequest true "register"
// @Success      200 {object} models.Responsee
// @Failure      400 {object} models.Responsee
// @Failure      404 {object} models.Responsee
// @Failure      500 {object} models.Responsee
func (h *Handler) Userloginwith_otp(c *gin.Context) {
	loginReq := models.UserRegisterRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handleResponseLog(c, h.Log, "error while binding body", http.StatusBadRequest, err)
		return
	}

	err := h.Services.Auth().UsersRegister(c.Request.Context(), loginReq)

	if err != nil {
		fmt.Println(err, "")
		handleResponseLog(c, h.Log, "", http.StatusInternalServerError, err)
		return
	}

	handleResponseLog(c, h.Log, "Otp sent successfull", http.StatusOK, "")

}
