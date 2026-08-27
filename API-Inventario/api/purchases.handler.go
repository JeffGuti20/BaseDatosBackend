package api

import (
	"net/http"
	"rest/dto"
	"time"

	"github.com/gin-gonic/gin"
)

// Requests para movimientos
type createPurchaseDetailRequest struct {
	SupplyID int32 `json:"supply_id" binding:"required"`
	Quantity int32 `json:"quantity" binding:"required"`
	//UnitPrice float64 `json:"unit_price" binding:"required"`
	Discount *float64 `json:"discount"`
	//Subtotal float64 `json:"subtotal"`
}

type createPurchaseRequest struct {
	UserID     int32 `json:"user_id" binding:"required"`
	SupplierID int32 `json:"supplier_id" binding:"required"`
	//Subtotal float64                        `json:"subtotal"`
	//Discount float64                        `json:"discount"`
	IVA float64 `json:"iva" binding:"required"`
	//Total   float64                        `json:"total"`
	Details []createPurchaseDetailRequest `json:"details" binding:"required"`
}

type updatePurchaseDetailRequest struct {
	SupplyID int32    `json:"supply_id" binding:"required"`
	Quantity int32    `json:"quantity" binding:"required"`
	Discount *float64 `json:"discount"`
}

type updatePurchaseRequest struct {
	PurchaseID int32                         `json:"purchase_id" binding:"required"`
	UserID     int32                         `json:"user_id" binding:"required"`
	SupplierID int32                         `json:"supplier_id"`
	IVA        float64                       `json:"iva" binding:"required"`
	Details    []updatePurchaseDetailRequest `json:"details" binding:"required"`
}

type deletePurchaseRequestJSON struct {
	PurchaseID int32 `json:"purchase_id" binding:"required"`
}

type deletePurchaseRequestURI struct {
	PurchaseID int32 `uri:"purchase_id" binding:"required"`
}

type getPurchaseByIDRequest struct {
	PurchaseID int32 `uri:"purchase_id" binding:"required"`
}

// Response para movimientos
type PurchaseDetailResponse struct {
	SupplyID  int32    `json:"supply_id"`
	Name      string   `json:"supply_name"`
	Tipo      string   `json:"supply_type"`
	Quantity  int32    `json:"quantity" `
	UnitPrice float64  `json:"unit_price" `
	Discount  *float64 `json:"discount"`
	Subtotal  float64  `json:"subtotal"`
}

type UserResponse struct {
	UserID int32  `json:"user_id"`
	Name   string `json:"name"`
}

type PurchaseResponse struct {
	PurchaseID int32                    `json:"purchase_id"`
	User       *UserResponse            `json:"user"`
	Supplier   *dto.Supplier            `json:"supplier"`
	Date	  	*time.Time                   `json:"date"`
	Subtotal   float64                  `json:"subtotal"`
	Discount   *float64                  `json:"discount"`
	IVA        float64                  `json:"iva"`
	Total      float64                  `json:"total"`
	Detalles   []PurchaseDetailResponse `json:"details"`
}

// createMovement es el handler para crear un nuevo movimiento.
func (server *Server) createPurchases(ctx *gin.Context) {
	var req createPurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// INICIAR TRANSACCIÓN
	tx, errTran := server.conn.BeginTx(ctx, nil)
	if errTran != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errTran))
		return
	}
	// ROLLBACK :P
	defer tx.Rollback()

	qtx := server.dbtx.WithTx(tx)

	args := dto.CreatePurchaseParams{
		UserID:     req.UserID,
		SupplierID: req.SupplierID,
		Iva:        req.IVA,
	}

	purchase, err := qtx.CreatePurchase(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	purchaseID, err := purchase.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var subtotal float64 = 0
	var descuentos float64 = 0
	for _, detalle := range req.Details {

		supply, errr := qtx.GetSupplyById(ctx, detalle.SupplyID)
		if errr != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "insumo no encontrado"})
			return
		}

		err := qtx.AddStock(ctx, dto.AddStockParams{
			SupplyID: detalle.SupplyID,
			Stock:    detalle.Quantity,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if detalle.Quantity <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "cantidad debe ser mayor a 0",
			})
			return
		}

		if detalle.Discount != nil &&
			(*detalle.Discount < 0 || *detalle.Discount > 100) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "El descuento debe estar entre 0 y 100",
			})
			return
		}

		subtotal_detalle := (supply.Price * float64(detalle.Quantity))
		argsDetail := dto.CreatePurchaseDetailParams{
			PurchaseID: int32(purchaseID),
			SupplyID:   detalle.SupplyID,
			Quantity:   detalle.Quantity,
			UnitPrice:  supply.Price,
			Discount:   detalle.Discount,
			Subtotal:   subtotal_detalle,
		}
		subtotal = subtotal + subtotal_detalle
		if detalle.Discount != nil {
			descuentoLinea := subtotal_detalle * (*detalle.Discount / 100)
			descuentos += descuentoLinea
		}
		_, err = qtx.CreatePurchaseDetail(ctx, argsDetail)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	base := subtotal - descuentos
	ivaMonto := base * (req.IVA / 100.0)
	total := base + ivaMonto

	var descuentosPtr *float64
	if descuentos > 0 {
		descuentosPtr = &descuentos
	}
	argsUpdate := dto.UpdatePurchaseParams{
		PurchaseID: int32(purchaseID),
		UserID:     req.UserID,
		SupplierID: req.SupplierID,
		Iva:        req.IVA,
		Subtotal:   subtotal,
		Discount:   descuentosPtr,
		Total:      total,
	}
	_, errrr := qtx.UpdatePurchase(ctx, argsUpdate)
	if errrr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errrr))
		return
	}

	// COMMIT :P
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"generated_id": purchaseID})
}

// updateMovement es el handler para actualizar un movimiento existente.
func (server *Server) updatePurchases(ctx *gin.Context) {
	var req updatePurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// INICIAR TRANSACCIÓN
	tx, errTran := server.conn.BeginTx(ctx, nil)
	if errTran != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errTran))
		return
	}

	// ROLLBACK :P
	defer tx.Rollback()
	qtx := server.dbtx.WithTx(tx)

	oldDetails, err := qtx.GetPurchaseDetailsByPurchaseId(ctx, req.PurchaseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, detail := range oldDetails {
		err := qtx.RemoveStock(ctx, dto.RemoveStockParams{
			SupplyID: detail.SupplyID,
			Stock:    detail.Quantity,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	_, errDel := qtx.DeletePurchaseDetail(ctx, req.PurchaseID)
	if errDel != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errDel))
		return
	}

	var subtotal float64 = 0
	var descuentos float64 = 0
	for _, detalle := range req.Details {
		supply, errr := qtx.GetSupplyById(ctx, detalle.SupplyID)
		if errr != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "insumo no encontrado"})
			return
		}

		err := qtx.AddStock(ctx, dto.AddStockParams{
			SupplyID: detalle.SupplyID,
			Stock:    detalle.Quantity,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if detalle.Quantity <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "cantidad debe ser mayor a 0",
			})
			return
		}

		if detalle.Discount != nil &&
			(*detalle.Discount < 0 || *detalle.Discount > 100) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "El descuento debe estar entre 0 y 100",
			})
			return
		}

		subtotal_detalle := (supply.Price * float64(detalle.Quantity))
		argsDetail := dto.CreatePurchaseDetailParams{
			PurchaseID: req.PurchaseID,
			SupplyID:   detalle.SupplyID,
			Quantity:   detalle.Quantity,
			UnitPrice:  supply.Price,
			Discount:   detalle.Discount,
			Subtotal:   subtotal_detalle,
		}
		subtotal = subtotal + subtotal_detalle
		if detalle.Discount != nil {
			descuentoLinea := subtotal_detalle * (*detalle.Discount / 100)
			descuentos += descuentoLinea
		}
		_, err = qtx.CreatePurchaseDetail(ctx, argsDetail)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	base := subtotal - descuentos
	ivaMonto := base * (req.IVA / 100.0)
	total := base + ivaMonto

	var descuentosPtr *float64
	if descuentos > 0 {
		descuentosPtr = &descuentos
	}
	argsUpdate := dto.UpdatePurchaseParams{
		PurchaseID: req.PurchaseID,
		UserID:     req.UserID,
		SupplierID: req.SupplierID,
		Iva:        req.IVA,
		Subtotal:   subtotal,
		Discount:   descuentosPtr,
		Total:      total,
	}
	_, errrr := qtx.UpdatePurchase(ctx, argsUpdate)
	if errrr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errrr))
		return
	}

	// COMMIT :P
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Compra actualizado correctamente"})
}

// deleteMovementJSON es el handler para eliminar un movimiento existente por medio del JSON.
func (server *Server) deletePurchaseJSON(ctx *gin.Context) {
	var req deletePurchaseRequestJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// INICIAR TRANSACCIÓN
	tx, errTran := server.conn.BeginTx(ctx, nil)
	if errTran != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errTran))
		return
	}

	// ROLLBACK :P
	defer tx.Rollback()
	qtx := server.dbtx.WithTx(tx)

	oldDetails, errr := qtx.GetPurchaseDetailsByPurchaseId(ctx, req.PurchaseID)
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errr))
		return
	}

	for _, detail := range oldDetails {

		err := qtx.RemoveStock(ctx, dto.RemoveStockParams{
			SupplyID: detail.SupplyID,
			Stock:    detail.Quantity,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	_, errDel := qtx.DeletePurchaseDetail(ctx, req.PurchaseID)
	if errDel != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errDel))
		return
	}

	_, err := qtx.DeletePurchase(ctx, req.PurchaseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// COMMIT :P
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Compra eliminado correctamente"})
}

// deleteMovementJSON es el handler para eliminar un movimiento existente por medio del JSON.
func (server *Server) deletePurchaseURI(ctx *gin.Context) {
	var req deletePurchaseRequestURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// INICIAR TRANSACCIÓN
	tx, errTran := server.conn.BeginTx(ctx, nil)
	if errTran != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errTran))
		return
	}

	// ROLLBACK :P
	defer tx.Rollback()
	qtx := server.dbtx.WithTx(tx)

	oldDetails, errr := qtx.GetPurchaseDetailsByPurchaseId(ctx, req.PurchaseID)
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errr))
		return
	}

	for _, detail := range oldDetails {

		err := qtx.RemoveStock(ctx, dto.RemoveStockParams{
			SupplyID: detail.SupplyID,
			Stock:    detail.Quantity,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	_, errDel := qtx.DeletePurchaseDetail(ctx, req.PurchaseID)
	if errDel != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errDel))
		return
	}

	_, err := qtx.DeletePurchase(ctx, req.PurchaseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// COMMIT :P
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Compra eliminado correctamente"})
}

// getAllMovement es el handler para obtener todos los movimientos con sus detalles y proveedores asociados.
func (server *Server) getAllPurchase(ctx *gin.Context) {
	purchases, err := server.dbtx.GetAllPurchases(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	details, err := server.dbtx.GetAllPurchaseDetails(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []PurchaseResponse
	for _, purchase := range purchases {

		var user *UserResponse
		if purchase.UserID != nil {
			user = &UserResponse{
				UserID: *purchase.UserID,
				Name:   *purchase.UserName,
			}
		}
		/*
			SupplierID   int32          `json:"supplier_id"`
			CompanyName  string         `json:"company_name"`
			CompanyPhone int32          `json:"company_phone"`
			CompanyEmail string         `json:"company_email"`
			ContactName  *string `json:"contact_name"`
			ContactPhone *int32  `json:"contact_phone"`
			ContactEmail *string `json:"contact_email"`
		*/
		var supplier *dto.Supplier
		if purchase.SupplierID != nil {
			supplier = &dto.Supplier{
				SupplierID:   *purchase.SupplierID,
				CompanyName:  *purchase.CompanyName,
				CompanyPhone: *purchase.CompanyPhone,
				CompanyEmail: *purchase.CompanyEmail,
				ContactName:  purchase.ContactName,
				ContactPhone: purchase.ContactPhone,
				ContactEmail: purchase.ContactEmail,
			}
		}

		var detalles []PurchaseDetailResponse
		for _, detail := range details {
			if detail.PurchaseID == purchase.PurchaseID {
				detalles = append(detalles, PurchaseDetailResponse{
					SupplyID:  detail.SupplyID,
					Name:      *detail.Name,
					Quantity:  detail.Quantity,
					UnitPrice: detail.UnitPrice,
					Discount:  detail.Discount,
					Subtotal:  detail.Subtotal,
				})
			}
		}

		response = append(response, PurchaseResponse{
			PurchaseID: purchase.PurchaseID,
			User:       user,
			Supplier:   supplier,
			Date: 		purchase.Date, 
			Detalles:   detalles,
			Subtotal:   purchase.Subtotal,
			Discount:   purchase.Discount,
			IVA:        purchase.Iva,
			Total:      purchase.Total,

	
		})
	}

	ctx.JSON(http.StatusOK, response)
}
