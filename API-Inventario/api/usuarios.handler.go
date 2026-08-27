package api

import (
	"database/sql"
	"net/http"
	"rest/dto"
	"rest/security"
	"time"
	"github.com/go-sql-driver/mysql"
	"fmt"
	"io"

	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// Cosas del Login
type loginRequest struct {
	Email      *string `json:"email" binding:"required"`
	Contraseña string  `json:"password" binding:"required"`
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
	RolID      int32   `json:"role_id" binding:"required"`
	Nombre     string  `json:"name" binding:"required"`
	Contraseña string  `json:"password" binding:"required"`
	Image      string `json:"image"`
	Email      *string `json:"email"`

}

type updateUserRequest struct {
	UsuarioID  int32   `json:"user_id" binding:"required"`
	RolID      int32   `json:"role_id" binding:"required"`
	Nombre     string  `json:"name" binding:"required"`
	Contraseña string  `json:"password" binding:"required"`
	Image      string `json:"image"`
	Email      *string `json:"email"`
}

type updateUserRoleRequest struct {
	UsuarioID int32 `uri:"user_id" binding:"required"`
	RolID     int32 `uri:"role_id" binding:"required"`
}

type deleteUserRequest struct {
	UsuarioID int32 `json:"user_id" binding:"required"`
}

type deleteUserRequestURI struct {
	UsuarioID int32 `uri:"user_id" binding:"required"`
}

type getUserByIDRequest struct {
	UsuarioID int32 `uri:"user_id" binding:"required"`
}

type getUserByNameRequest struct {
	Nombre string `uri:"name" binding:"required"`
}

// Responses para usuarios
type userResponse struct {
	UsuarioID int32     `json:"user_id"`
	Nombre    string    `json:"name"`
	Email     *string   `json:"email"`
	Image     string   `json:"image"`
	Rol       *dto.Role `json:"role"`
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
	if user.Contraseña != req.Contraseña {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Usuario no autorizado"})
		return
	}
	//GENERACION DEL TOKEN
	accessToken, err := server.tokenBuilder.CreateToken(
		user.UsuarioID,
		*user.RoleID,
		user.Nombre,
		time.Hour*1,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := LoginResponse{
		AccessToken: accessToken,
		Payload: payload{
			ID:   user.UsuarioID,
			Name: user.Nombre,
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
		RolID:      req.RolID,
		Nombre:     req.Nombre,
		Contraseña: req.Contraseña,
		Email:      req.Email,
		Image:      req.Image,
	}

	result, err := server.dbtx.CreateUser(ctx, arg)
	if err != nil {
    if mysqlErr, ok := err.(*mysql.MySQLError); ok {
        if mysqlErr.Number == 1062 {
            ctx.JSON(http.StatusConflict, gin.H{
                "error": "El correo ya está registrado",
            })
            return
        }

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
    }

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
	

	authPayload := ctx.MustGet("authorized").(*security.Payload)

	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req updateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if authPayload.UserID == req.UsuarioID && authPayload.Role != req.RolID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "No puedes cambiar tu propio rol",
		})
		return
	}

	arg := dto.UpdateUserParams{
		RolID:      req.RolID,
		Nombre:     req.Nombre,
		Contraseña: req.Contraseña,
		Email:      req.Email,
		UsuarioID:  req.UsuarioID,
		Image:      req.Image,
	}

	_, err := server.dbtx.UpdateUser(ctx, arg)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "El correo ya está registrado",
				})
				return
			}
		}

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


	if authPayload.UserID == uriReq.UsuarioID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "No puedes cambiar tu propio rol",
		})
		return
	}

	arg := dto.UpdateUserRoleParams{
		UsuarioID: uriReq.UsuarioID,
		RolID:     uriReq.RolID,
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

	if authPayload.UserID == req.UsuarioID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "No puedes eliminar tu propia cuenta",
		})
		return
	}

	_, err := server.dbtx.DeleteUser(ctx, req.UsuarioID)

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

	if authPayload.UserID == req.UsuarioID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "No puedes eliminar tu propia cuenta",
		})
		return
	}

	_, err := server.dbtx.DeleteUser(ctx, req.UsuarioID)

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
				Nombre: *user.RolNombre,
			}
		}
		response = append(response, userResponse{
			UsuarioID: user.UsuarioID,
			Nombre:    user.Nombre,

			Email: user.Email,
			Rol:   role,
		})

	}

	ctx.JSON(http.StatusOK, response)
}


// getUserByName es el handler para obtener un usuario por su nombre.
func (server *Server) getUserByName(ctx *gin.Context) {
	var req getUserByNameRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.dbtx.GetUserByName(ctx, req.Nombre)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
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
	user, err := server.dbtx.GetUserRoleById(ctx, request.UsuarioID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var role *dto.Role
	if user.RoleID != nil {
		role = &dto.Role{
			RoleID: *user.RoleID,
			Nombre: *user.RolNombre,
		}
	}
	response := userResponse{
		UsuarioID: user.UsuarioID,
		Nombre:    user.Nombre,
		Email:     user.Email,
		Rol:       role,
	}
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) uploadUserImg(ctx *gin.Context) {

	fileHeader, err := ctx.FormFile("file0")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer file.Close()

	wd, _ := os.Getwd()
	fmt.Println("Working directory:", wd)
	upDir := "/data/users"

	err = os.MkdirAll(upDir, os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	filename := uuid.New().String() + "_" + filepath.Base(fileHeader.Filename)

	destinationFile, err := os.Create(filepath.Join(upDir, filename))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"message":  "Imagen cargada exitosamente",
	})
}

func (server *Server) downloadUserImg(ctx *gin.Context) {

	var req userImageRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fileUrl := "/data/users/" + req.Filename

	ctx.File(fileUrl)
}

func (server *Server) deleteUserImg(ctx *gin.Context) {

    var req userImageRequest
    if err := ctx.ShouldBindUri(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    filePath := filepath.Join("/data/users", req.Filename)

    err := os.Remove(filePath)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "No se pudo eliminar la imagen",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Imagen eliminada correctamente",
    })
}
