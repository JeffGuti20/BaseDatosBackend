package api

import (
	"net/http"
	"rest/dto"

	"github.com/gin-gonic/gin"
)

type createProveedorRequest struct {
	Nombre   string  `json:"name" binding:"required"`
	Telefono *int32  `json:"phone" `
	Email    *string `json:"email" `
}

type updateProveedorRequest struct {
	ProveedorID int32   `json:"supplier_id" binding:"required"`
	Nombre      string  `json:"name" binding:"required"`
	Telefono    *int32  `json:"phone"`
	Email       *string `json:"email"`
}

type deleteProveedorRequestJSON struct {
	ProveedorID int32 `json:"supplier_id" binding:"required"`
}

type deleteProveedorRequestURI struct {
	ProveedorID int32 `uri:"supplier_id" binding:"required"`
}

// SELECT * FROM proveedores WHERE proveedor_id=?;
type getProveedorRequest struct {
	ProveedorID int32 `uri:"supplier_id" binding:"required"`
}


// createProveedor es el handler para crear un nuevo proveedor o supplier.
func (server *Server) createProveedor(ctx *gin.Context) {
	var req createProveedorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := dto.CreateSupplierParams{
		Nombre:   req.Nombre,
		Telefono: req.Telefono,
		Email:    req.Email,
	}
	supplier, err := server.dbtx.CreateSupplier(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var lastId, _ = supplier.LastInsertId()
	ctx.JSON(http.StatusOK, gin.H{"generated_id": lastId})
}

// updateProveedor es el handler para actualizar un proveedor o suppliers.
func (server *Server) updateProveedor(ctx *gin.Context) {
	var req updateProveedorRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := dto.UpdateSupplierParams{
		Nombre:      req.Nombre,
		Telefono:    req.Telefono,
		Email:       req.Email,
		ProveedorID: req.ProveedorID,
	}

	_, err := server.dbtx.UpdateSupplier(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Proveedor actualizado correctamente",
	})
}

// deleteProveedor elimina un proveedor por JSON.
func (server *Server) deleteProveedorJSON(ctx *gin.Context) {
	var req deleteProveedorRequestJSON

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.dbtx.DeleteSupplier(ctx, req.ProveedorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Proveedor eliminado correctamente",
	})
}

// deleteProveedor elimina un proveedor por URI.
func (server *Server) deleteProveedorURI(ctx *gin.Context) {
	var req deleteProveedorRequestURI

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.dbtx.DeleteSupplier(ctx, req.ProveedorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Proveedor eliminado correctamente",
	})
}

// getAllProveedores es el handler para obtener todos los proveedores o suppliers.
func (server *Server) getAllProveedores(ctx *gin.Context) {
	suppliers, err := server.dbtx.GetAllSuppliers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, suppliers)
}

// getByIDProveedores es el handler para obtener un proveedor por su ID.
func (server *Server) getProveedoresByID(ctx *gin.Context) {
	var req getProveedorRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supplier, err := server.dbtx.GetSupplierById(ctx, req.ProveedorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, supplier)
}
