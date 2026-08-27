package api

import (
	"net/http"
	"rest/dto"
	"rest/security"

	"github.com/gin-gonic/gin"
)

// Request para roles
type createRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type updateRoleRequest struct {
	RoleID int32  `json:"role_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

type deleteRoleRequestJSON struct {
	RoleID int32 `json:"role_id" binding:"required"`
}

type deleteRoleRequestURI struct {
	RoleID int32 `uri:"role_id" binding:"required"`
}

type getRoleByIDRequest struct {
	RoleID int32 `uri:"role_id" binding:"required"`
}

// createRole es el handler para crear un nuevo rol.
func (server *Server) createRole(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req createRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.dbtx.CreateRole(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	lastID, _ := result.LastInsertId()
	ctx.JSON(http.StatusOK, gin.H{
		"generated_id": lastID,
	})
}

// updateRole es el handler para actualizar un rol.
func (server *Server) updateRole(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req updateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := dto.UpdateRoleParams{
		RoleID: req.RoleID,
		Name:   req.Name,
	}

	_, err := server.dbtx.UpdateRole(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Rol actualizado correctamente",
	})
}

// deleteRole elimina un rol por JSON.
func (server *Server) deleteRoleJSON(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req deleteRoleRequestJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.dbtx.DeleteRole(ctx, req.RoleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Rol eliminado correctamente",
	})
}

// deleteRole elimina un rol por URI.
func (server *Server) deleteRoleURI(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req deleteRoleRequestURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.dbtx.DeleteRole(ctx, req.RoleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Rol eliminado correctamente",
	})
}

// getAllRoles obtiene todos los roles.
func (server *Server) getAllRoles(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	roles, err := server.dbtx.GetAllRoles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

// getRoleByID es el handler para obtener un rol por ID.
func (server *Server) getRoleByID(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req getRoleByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	role, err := server.dbtx.GetRoleById(ctx, req.RoleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, role)
}
