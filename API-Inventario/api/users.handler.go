package api

import (
	"database/sql"
	"net/http"

	"rest/dto"
	"rest/security"
	"time"

	"github.com/gin-gonic/gin"

)

// Cosas del Login
type loginRequest struct {
	Email    *string `json:"email" binding:"required"`
	Password string  `json:"password" binding:"required"`
}
type LoginResponse struct {
	AccessToken string  `json:"access_token"`
	Payload     payload `json:"payload"`
}
type payload struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Role int32  `json:"role"`
}

// Requests para usuarios
type createUserRequest struct {
	RoleID   int32   `json:"role_id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Email    *string `json:"email"`
}

type updateUserRequest struct {
	UserID   int32   `json:"user_id" binding:"required"`
	RoleID   int32   `json:"role_id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Email    *string `json:"email"`
}

type updateProfileRequest struct {
	UserID   int32   `json:"user_id"`
	Name     string  `json:"name"`
	Email    *string `json:"email"`
}

type updateUserRoleRequest struct {
	UserID int32 `uri:"user_id" binding:"required"`
	RoleID     int32 `uri:"role_id" binding:"required"`
}

type deleteUserRequest struct {
	UserID int32 `json:"user_id" binding:"required"`
}

type deleteUserRequestURI struct {
	UserID int32 `uri:"user_id" binding:"required"`
}

type getUserByIDRequest struct {
	UserID int32 `uri:"user_id" binding:"required"`
}

// Responses para usuarios
type userResponse struct {
	UserID int32     `json:"user_id"`
	Name   string    `json:"name"`
	Email  *string   `json:"email"`
	Rol    *dto.Role `json:"role"`
}

type userImageRequest struct {
	Filename string `uri:"filename" binding:"required"`
}

/*Funcion de Login*/
func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.dbtx.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if user.Password != req.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Usuario no autorizado"})
		return
	}
	//GENERACION DEL TOKEN
	accessToken, err := server.tokenBuilder.CreateToken(*user.RoleID, user.UserID, user.Name, time.Hour*1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := LoginResponse{
		AccessToken: accessToken,
		Payload: payload{
			ID:   user.UserID,
			Name: user.Name,
			Role: *user.RoleID,
		},
	}
	ctx.JSON(http.StatusOK, resp)

}

// createUser es el handler para crear un nuevo usuario.
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := dto.CreateUserParams{
		RoleID:   req.RoleID,
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
	}

	result, err := server.dbtx.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	id, _ := result.LastInsertId()

	ctx.JSON(http.StatusOK, gin.H{
		"generated_id": id,
	})
}

// updateUser es el handler para actualizar un usuario.
func (server *Server) updateUser(ctx *gin.Context) {



	var req updateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := dto.UpdateUserParams{
		RoleID:   req.RoleID,
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		UserID:   req.UserID,
	}

	_, err := server.dbtx.UpdateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuario actualizado correctamente",
	})
}

// updateUserRole actualiza únicamente el rol de un usuario.
func (server *Server) updateUserRole(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)

	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var uriReq updateUserRoleRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}


	if authPayload.UserID == uriReq.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "No puedes cambiar tu propio rol",
		})
		return
	}

	arg := dto.UpdateUserRoleParams{
		UserID: uriReq.UserID,
		RoleID:     uriReq.RoleID,
	}

	if _, err := server.dbtx.UpdateUserRole(ctx, arg); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Rol de usuario actualizado correctamente",
	})
}

// deleteUser elimina un usuario por JSON.
func (server *Server) deleteUserJSON(ctx *gin.Context) {

	authPayload := ctx.MustGet("authorized").(*security.Payload)

	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req deleteUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.dbtx.DeleteUser(ctx, req.UserID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuario eliminado correctamente",
	})
}

// PETICION DE COCA COLA MASTER
// deleteUserURI LO que hace es que elimina un usuario por medio de la URI.
func (server *Server) deleteUserURI(ctx *gin.Context) {

	authPayload := ctx.MustGet("authorized").(*security.Payload)

	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req deleteUserRequestURI

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.dbtx.DeleteUser(ctx, req.UserID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuario eliminado correctamente",
	})
}

// getAllUsers es el handler para obtener todos los usuarios.
func (server *Server) getAllUsers(ctx *gin.Context) {

	authPayload := ctx.MustGet("authorized").(*security.Payload)

	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	users, err := server.dbtx.GetAllUserRole(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []userResponse
	for _, user := range users {
		var role *dto.Role
		if user.RoleID != nil {
			role = &dto.Role{
				RoleID: *user.RoleID,
				Name:   *user.RoleName,
			}
		}
		response = append(response, userResponse{
			UserID: user.UserID,
			Name:   user.Name,
			Email:  user.Email,
			Rol:    role,
		})

	}

	ctx.JSON(http.StatusOK, response)
}

// getUserRoleByID es el handler para obtener un usuario por ID, incluyendo su rol.
func (server *Server) getUserRoleByID(ctx *gin.Context) {

	authPayload := ctx.MustGet("authorized").(*security.Payload)

	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var request getUserByIDRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.dbtx.GetUserRoleById(ctx, request.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var role *dto.Role
	if user.RoleID != nil {
		role = &dto.Role{
			RoleID: *user.RoleID,
			Name:   *user.RoleName,
		}
	}
	response := userResponse{
		UserID: user.UserID,
		Name:   user.Name,
		Email:  user.Email,
		Rol:    role,
	}
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) getUserByID(ctx *gin.Context) {
	var request getUserByIDRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.dbtx.GetUserById(ctx, request.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := userResponse{
		UserID: user.UserID,
		Name:   user.Name,
		Email:  user.Email,
	}
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) updateProfile(ctx *gin.Context) {

	var req updateProfileRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}


	arg := dto.UpdateProfileParams{
		UserID: req.UserID,
		Name:   req.Name,
		Email:  req.Email,
	}

	_, err := server.dbtx.UpdateProfile(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Perfil actualizado correctamente",
	})
}
