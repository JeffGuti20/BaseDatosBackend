package api

import (
	"net/http"
	"rest/dto"

	"github.com/gin-gonic/gin"
)

type createSupplierRequest struct {
	CompanyName  string  `json:"company_name" binding:"required"`
	CompanyPhone int32   `json:"company_phone" binding:"required"`
	CompanyEmail string  `json:"company_email" binding:"required"`
	ContactName  *string `json:"contact_name" `
	ContactPhone *int32  `json:"contact_phone"`
	ContactEmail *string `json:"contact_email"`
}

type updateSupplierRequest struct {
	CompanyName  string  `json:"company_name" binding:"required"`
	CompanyPhone int32   `json:"company_phone" binding:"required"`
	CompanyEmail string  `json:"company_email" binding:"required"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *int32  `json:"contact_phone"`
	ContactEmail *string `json:"contact_email"`
	SupplierID   int32   `json:"supplier_id" binding:"required"`
}

type deleteSupplierRequestJSON struct {
	SupplierID int32 `json:"supplier_id" binding:"required"`
}

type deleteSupplierRequestURI struct {
	SupplierID int32 `uri:"supplier_id" binding:"required"`
}

// SELECT * FROM proveedores WHERE proveedor_id=?;
type getSupplierRequest struct {
	SupplierID int32 `uri:"supplier_id" binding:"required"`
}

// createProveedor es el handler para crear un nuevo proveedor o supplier.
func (server *Server) createSupplier(ctx *gin.Context) {
	var req createSupplierRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := dto.CreateSupplierParams{
		CompanyName:  req.CompanyName,
		CompanyPhone: req.CompanyPhone,
		CompanyEmail: req.CompanyEmail,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
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
func (server *Server) updateSupplier(ctx *gin.Context) {
	var req updateSupplierRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := dto.UpdateSupplierParams{
		CompanyName:  req.CompanyName,
		CompanyPhone: req.CompanyPhone,
		CompanyEmail: req.CompanyEmail,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		SupplierID:   req.SupplierID,
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
func (server *Server) deleteSupplierJSON(ctx *gin.Context) {
	var req deleteSupplierRequestJSON

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.dbtx.DeleteSupplier(ctx, req.SupplierID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Proveedor eliminado correctamente",
	})
}

// deleteProveedor elimina un proveedor por URI.
func (server *Server) deleteSupplierURI(ctx *gin.Context) {
	var req deleteSupplierRequestURI

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.dbtx.DeleteSupplier(ctx, req.SupplierID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Proveedor eliminado correctamente",
	})
}

// getAllProveedores es el handler para obtener todos los proveedores o suppliers.
func (server *Server) getAllSuppliers(ctx *gin.Context) {
	suppliers, err := server.dbtx.GetAllSuppliers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, suppliers)
}

// getByIDProveedores es el handler para obtener un proveedor por su ID.
func (server *Server) getSuppliersByID(ctx *gin.Context) {
	var req getSupplierRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supplier, err := server.dbtx.GetSupplierById(ctx, req.SupplierID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, supplier)
}
