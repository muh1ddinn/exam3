package handler

import (
	"context"
	models "exam3/api/model"
	"exam3/pkg/check"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Createuser godoc
// @Security ApiKeyAuth
// @Router      /user [POST]
// @Summary     Create a user
// @Description Create a new user
// @Tags        user
// @Accept      json
// @Produce 	json
// @Param 		user body models.Users true "user"
// @Success 	200  {object}  string
// @Failure		400  {object}  models.Responsee
// @Failure		404  {object}  models.Responsee
// @Failure		500  {object}  models.Responsee
func (h Handler) Createcus(c *gin.Context) {
	cus := models.Users{}

	if err := c.ShouldBindJSON(&cus); err != nil {
		return
	}

	id, err := h.Services.User().CreateCus(context.Background(), cus)
	if err != nil {
		handleResponseLog(c, h.Log, "error while creating user", http.StatusBadRequest, err.Error())
		return
	}
	handleResponseLog(c, h.Log, "Created successfully", http.StatusOK, id)
}

// GetAllUSer godoc
// @Security ApiKeyAuth
// @Router 			/user [GET]
// @Summary 		Get all user
// @Description		Retrieves information about all user.
// @Tags 			user
// @Accept 			json
// @Produce 		json
// @Param 			search query string true "user"
// @Param 			page query uint64 false "page"
// @Param 			limit query uint64 false "limit"
// @Success 		200 {object} models.GetAllusersResponse
// @Failure 		400 {object} models.Responsee
// @Failure 		500 {object} models.Responsee
func (h Handler) Getalluser(c *gin.Context) {

	var request = models.GetAllUsersRequest{}

	request.Search = c.Query("search")

	page, err := ParsePageQueryParam(c)

	if err != nil {
		handleResponseLog(c, h.Log, "error while parsing page", http.StatusBadRequest, err.Error())

	}

	limit, err := ParseLimitQueryParam(c)
	if err != nil {
		handleResponseLog(c, h.Log, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("page: ", page)
	fmt.Println("limit: ", limit)

	request.Page = page
	request.Limit = limit

	customer, err := h.Services.User().GetAllCus(context.Background(), request)
	if err != nil {

		handleResponseLog(c, h.Log, "error while getting user", http.StatusBadRequest, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "", http.StatusOK, customer)
}

// GetUserById godoc
// @Security ApiKeyAuth
// @Router		/user/{id} [GET]
// @Summary		get a user by its id
// @Description This api gets a user by its id and returns its info
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		id path string true "user"
// @Success		200  {object}  models.Users
// @Failure		400  {object}  models.Responsee
// @Failure		404  {object}  models.Responsee
// @Failure		500  {object}  models.Responsee
func (h Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		handleResponseLog(c, h.Log, "missing car ID", http.StatusBadRequest, id)
		return
	}

	customer, err := h.Services.User().GetByID(context.Background(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while getting user by ID", http.StatusBadRequest, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "user was successfully gotten by Id", http.StatusOK, customer)
}

// DeleteUser godoc
// @Security ApiKeyAuth
// @Router		/user/{id} [DELETE]
// @Summary		delete a user by its id
// @Description This api deletes a user by its id and returns error or nil
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		id path string true "user ID"
// @Success		200  {object}  nil
// @Failure		400  {object}  models.Responsee
// @Failure		404  {object}  models.Responsee
// @Failure		500  {object}  models.Responsee
func (h Handler) Deleteuser(c *gin.Context) {

	id := c.Param("id")
	fmt.Println("id: ", id)

	err := uuid.Validate(id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while validating id ", http.StatusBadRequest, err.Error())
		return
	}

	id, err = h.Services.User().Deletuser(context.Background(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while deleting user", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "all good ", http.StatusOK, id)

}

// Createcus godoc
// @Security ApiKeyAuth
// @Router      /user [PUT]
// @Summary     Update a USER
// @Description Update a new USER
// @Tags        USER
// @Accept      json
// @Produce 	json
// @Param 		USER body models.Users true "USER"
// @Success 	200  {object}  string
// @Failure		400  {object}  models.Responsee
// @Failure		404  {object}  models.Responsee
// @Failure		500  {object}  models.Responsee
func (h Handler) Update(c *gin.Context) {
	cus := models.Users{}

	if err := c.ShouldBindJSON(&cus); err != nil {
		return
	}

	id, err := h.Services.User().UpdateUSER(context.Background(), cus)
	if err != nil {
		handleResponseLog(c, h.Log, "error while creating user", http.StatusBadRequest, err.Error())
		return
	}
	handleResponseLog(c, h.Log, "Created successfully", http.StatusOK, id)
}

// Create godoc
// @Security ApiKeyAuth
// @Router      /user/password [PUT]
// @Summary     Update a password
// @Description Update a password
// @Tags        USER
// @Accept      json
// @Produce 	json
// @Param 		USER body models.Changepasswor true "USER"
// @Success 	200  {object}  string
// @Failure		400  {object}  models.Responsee
// @Failure		404  {object}  models.Responsee
// @Failure		500  {object}  models.Responsee
func (h Handler) Changepassword(c *gin.Context) {
	pass := models.Changepasswor{}

	if err := c.ShouldBindJSON(&pass); err != nil {
		return
	}

	if err := check.ValidatePassword(pass.OldPassword); err != nil {
		handleResponseLog(c, h.Log, "error while validating password", http.StatusBadRequest, err.Error())
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass.OldPassword), bcrypt.DefaultCost)
	if err != nil {
		handleResponseLog(c, h.Log, "error while hashing new password", http.StatusInternalServerError, err.Error())
		return
	}
	pass.OldPassword = string(hashedPassword)

	msg, err := h.Services.User().ChangePassword(context.Background(), pass)
	if err != nil {
		handleResponseLog(c, h.Log, "error while updating customer", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Customer password was successfully updated", http.StatusOK, msg)
}

// Updatestatus godoc
// @Security ApiKeyAuth
// @Router      /status [PATCH]
// @Summary     Update a user's status
// @Description Update a user's status
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       user body models.Updatestatus true "User"
// @Success     200 {string} string "OK"
// @Failure     400 {object} models.Responsee "Bad Request"
// @Failure     404 {object} models.Responsee "Not Found"
// @Failure     500 {object} models.Responsee "Internal Server Error"
func (h *Handler) Updatestatus(c *gin.Context) {
	cus := models.Updatestatus{}

	if err := c.ShouldBindJSON(&cus); err != nil {
		handleResponseLog(c, h.Log, "error while binding JSON", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.Services.User().Updatestatus(c.Request.Context(), cus)
	if err != nil {
		handleResponseLog(c, h.Log, "error while updating user status", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Updated user status successfully", http.StatusOK, id)
}
