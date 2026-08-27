package api

import (
	"rest/dto"

	"net/http"

	"github.com/gin-gonic/gin"
)

type createClientRequest struct {
	Name string `json:"name" binding:"required"`
	ContactPhone int32  `json:"contact_phone" binding:"required"`
	ContactEmail string `json:"contact_email" binding:"required"`
}

type updateClientRequest struct {
	Name string `json:"name" binding:"required"`
	ContactPhone int32  `json:"contact_phone" binding:"required"`
	ContactEmail string `json:"contact_email" binding:"required"`
	ClientID int32 `json:"client_id" binding:"required"`
}

type deleteClientRequestJSON struct {
	ClientID int32 `json:"client_id" binding:"required"`
}

type deleteClientRequestURI struct {
	ClientID int32 `uri:"client_id" binding:"required"`
}

// SELECT * FROM proveedores WHERE proveedor_id=?;
type getClientRequest struct {
	ClientID int32 `uri:"client_id" binding:"required"`
}

// createClient es el handler para crear un nuevo client.
func (server *Server) createClient(ctx *gin.Context) {
	var req createClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := dto.CreateClientParams{
		Name:  req.Name,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
	}
	client, err := server.dbtx.CreateClient(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var lastId, _ = client.LastInsertId()
	ctx.JSON(http.StatusOK, gin.H{"generated_id": lastId})
}

// updateProveedor es el handler para actualizar un proveedor o suppliers.
func (server *Server) updateClient(ctx *gin.Context) {
	var req updateClientRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := dto.UpdateClientParams{
		Name:  req.Name,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		ClientID:   req.ClientID,
	}

	_, err := server.dbtx.UpdateClient(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cliente actualizado correctamente",
	})
}

// deleteProveedor elimina un proveedor por JSON.
func (server *Server) deleteClientJSON(ctx *gin.Context) {
	var req deleteClientRequestJSON

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.dbtx.DeleteClient(ctx, req.ClientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cliente eliminado correctamente",
	})
}

// deleteProveedor elimina un proveedor por URI.
func (server *Server) deleteClientURI(ctx *gin.Context) {
	var req deleteClientRequestURI

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.dbtx.DeleteClient(ctx, req.ClientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cliente eliminado correctamente",
	})
}

// getAllProveedores es el handler para obtener todos los proveedores o suppliers.
func (server *Server) getAllClients(ctx *gin.Context) {
	clients, err := server.dbtx.GetAllClients(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, clients)
}

// getByIDProveedores es el handler para obtener un proveedor por su ID.
func (server *Server) getClientsByID(ctx *gin.Context) {
	var req getClientRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	client, err := server.dbtx.GetClientById(ctx, req.ClientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, client)
}
