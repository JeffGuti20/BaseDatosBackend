package api

import (
	"net/http"
	"rest/dto"
	"time"

	"github.com/gin-gonic/gin"
)

// Requests para movimientos
type createMovementDetailRequest struct {
	InsumoID int32 `json:"supply_id" binding:"required"`
	Cantidad int32 `json:"quantity" binding:"required"`
}

type createMovementRequest struct {
	UsuarioID   int32                         `json:"user_id" binding:"required"`
	ProveedorID *int32                        `json:"supplier_id"`
	Detalles    []createMovementDetailRequest `json:"details" binding:"required"`
}

type updateMovementDetailRequest struct {
	InsumoID int32 `json:"supply_id" binding:"required"`
	Cantidad int32 `json:"quantity" binding:"required"`
}

type updateMovementRequest struct {
	MovimientoID int32                         `json:"movement_id" binding:"required"`
	UsuarioID    int32                         `json:"user_id" binding:"required"`
	ProveedorID  *int32                        `json:"supplier_id"`
	Detalles     []updateMovementDetailRequest `json:"details" binding:"required"`
}

type deleteMovementRequestJSON struct {
	MovimientoID int32 `json:"movement_id" binding:"required"`
}

type deleteMovementRequestURI struct {
	MovimientoID int32 `uri:"movement_id" binding:"required"`
}

type getMovementByIDRequest struct {
	MovimientoID int32 `uri:"movement_id" binding:"required"`
}

// Response para movimientos
type MovementDetailResponse struct {
	InsumoID int32  `json:"supply_id"`
	Nombre   string `json:"supply_name"`
	Tipo     string `json:"supply_type"`
	Cantidad int32  `json:"quantity"`
}

type UserResponse struct {
	UsuarioID int32  `json:"user_id"`
	Nombre    string `json:"name"`
}

type MovementResponse struct {
	MovimientoID int32                    `json:"movement_id"`
	Usuario      *UserResponse            `json:"user"`
	Tipo         string                   `json:"type"`
	Proveedor    *dto.Proveedore          `json:"supplier"`
	Fecha 		*time.Time 				  `json:"date"`
	Detalles     []MovementDetailResponse `json:"details"`
}

// createMovement es el handler para crear un nuevo movimiento.
func (server *Server) createMovement(ctx *gin.Context) {
	var req createMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var tipo dto.MovimientosTipo
	if req.ProveedorID == nil {
		tipo = dto.MovimientosTipo("salida")
	} else {
		tipo = dto.MovimientosTipo("entrada")
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

	args := dto.CreateMovementParams{
		UsuarioID:   req.UsuarioID,
		ProveedorID: req.ProveedorID,
		Tipo:        tipo,
	}

	movement, err := qtx.CreateMovement(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	movimientoID, err := movement.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, detalle := range req.Detalles {

		supply, errr := qtx.GetSupplyById(ctx, detalle.InsumoID)
    	if errr != nil {
        	ctx.JSON(http.StatusNotFound, gin.H{"error": "insumo no encontrado"})
        	return
    	}
		

		if tipo == "entrada" {
			err := qtx.AddStock(ctx, dto.AddStockParams{
				InsumoID: detalle.InsumoID,
				Cantidad: detalle.Cantidad,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		if tipo == "salida" {
			if supply.Cantidad < detalle.Cantidad {
            	ctx.JSON(http.StatusBadRequest, gin.H{"error": "stock insuficiente",})
            		return
        	}	

			err := qtx.RemoveStock(ctx, dto.RemoveStockParams{
				InsumoID: detalle.InsumoID,
				Cantidad: detalle.Cantidad,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}

		if detalle.Cantidad <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "cantidad debe ser mayor a 0",
			})
			return
		}
		argsDetail := dto.CreateMovementDetailParams{
			MovimientoID: int32(movimientoID),
			InsumoID:     detalle.InsumoID,
			Cantidad:     detalle.Cantidad,
		}

		_, err := qtx.CreateMovementDetail(ctx, argsDetail)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	// COMMIT :P
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"generated_id": movimientoID})
}

// updateMovement es el handler para actualizar un movimiento existente.
func (server *Server) updateMovement(ctx *gin.Context) {
	var req updateMovementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var tipo dto.MovimientosTipo
	if req.ProveedorID == nil {
		tipo = dto.MovimientosTipo("salida")
	} else {
		tipo = dto.MovimientosTipo("entrada")
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

	args := dto.UpdateMovementParams{
		MovimientoID: req.MovimientoID,
		UsuarioID:    req.UsuarioID,
		ProveedorID:  req.ProveedorID,
		Tipo:         tipo,
	}
	_, err := qtx.UpdateMovement(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	oldMovement, err := qtx.GetMovementsById(ctx, req.MovimientoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	oldDetails, err := qtx.GetDetailById(ctx, req.MovimientoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, detail := range oldDetails {
		if oldMovement.Tipo == "entrada" {
			err := qtx.RemoveStock(ctx, dto.RemoveStockParams{
				InsumoID: detail.InsumoID,
				Cantidad: detail.Cantidad,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			err := qtx.AddStock(ctx, dto.AddStockParams{
				InsumoID: detail.InsumoID,
				Cantidad: detail.Cantidad,
			})
				if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
	}

	_, errDel := qtx.DeleteMovementDetail(ctx, req.MovimientoID)
	if errDel != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errDel))
		return
	}

	for _, detalle := range req.Detalles {
		supply, errr := qtx.GetSupplyById(ctx, detalle.InsumoID)
    	if errr != nil {
        	ctx.JSON(http.StatusInternalServerError, errorResponse(errr))
        	return
    	}

		if tipo == "entrada" {
			err := qtx.AddStock(ctx, dto.AddStockParams{
				InsumoID: detalle.InsumoID,
				Cantidad: detalle.Cantidad,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		if tipo == "salida" {
			if supply.Cantidad < detalle.Cantidad {
            	ctx.JSON(http.StatusBadRequest, gin.H{"error": "stock insuficiente",})
            		return
        	}	

			err := qtx.RemoveStock(ctx, dto.RemoveStockParams{
				InsumoID: detalle.InsumoID,
				Cantidad: detalle.Cantidad,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		
		if detalle.Cantidad <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "cantidad debe ser mayor a 0",
			})
			return
		}
		argsDetail := dto.CreateMovementDetailParams{
			MovimientoID: req.MovimientoID,
			InsumoID:     detalle.InsumoID,
			Cantidad:     detalle.Cantidad,
		}

		_, err := qtx.CreateMovementDetail(ctx, argsDetail)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// COMMIT :P
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Movimiento actualizado correctamente"})
}

// deleteMovementJSON es el handler para eliminar un movimiento existente por medio del JSON.
func (server *Server) deleteMovementJSON(ctx *gin.Context) {
	var req deleteMovementRequestJSON
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

	oldMovement, errr := qtx.GetMovementsById(ctx, req.MovimientoID)
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errr))
		return
	}
	
	oldDetails, errr := qtx.GetDetailById(ctx, req.MovimientoID)
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errr))
		return
	}

	for _, detail := range oldDetails {
		if oldMovement.Tipo == "entrada" {
			err := qtx.RemoveStock(ctx, dto.RemoveStockParams{
				InsumoID: detail.InsumoID,
				Cantidad: detail.Cantidad,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			err := qtx.AddStock(ctx, dto.AddStockParams{
				InsumoID: detail.InsumoID,
				Cantidad: detail.Cantidad,
			})
				if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
	}

	_, errDel := qtx.DeleteMovementDetail(ctx, req.MovimientoID)
	if errDel != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errDel))
		return
	}

	_, err := qtx.DeleteMovement(ctx, req.MovimientoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// COMMIT :P
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Movimiento eliminado correctamente"})
}

// deleteMovementURI es el handler para eliminar un movimiento existente por medio de la URI.
func (server *Server) deleteMovementURI(ctx *gin.Context) {
	var req deleteMovementRequestURI
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

	oldMovement, errr := qtx.GetMovementsById(ctx, req.MovimientoID)
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errr))
		return
	}
	
	oldDetails, errr := qtx.GetDetailById(ctx, req.MovimientoID)
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errr))
		return
	}

	for _, detail := range oldDetails {
		if oldMovement.Tipo == "entrada" {
			err := qtx.RemoveStock(ctx, dto.RemoveStockParams{
				InsumoID: detail.InsumoID,
				Cantidad: detail.Cantidad,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			err := qtx.AddStock(ctx, dto.AddStockParams{
				InsumoID: detail.InsumoID,
				Cantidad: detail.Cantidad,
			})
				if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
	}

	_, errDel := qtx.DeleteMovementDetail(ctx, req.MovimientoID)
	if errDel != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errDel))
		return
	}

	_, err := qtx.DeleteMovement(ctx, req.MovimientoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// COMMIT :P
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Movimiento eliminado correctamente"})
}

// getAllMovement es el handler para obtener todos los movimientos con sus detalles y proveedores asociados.
func (server *Server) getAllMovement(ctx *gin.Context) {
	movements, err := server.dbtx.GetAllMovements(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	details, err := server.dbtx.GetAllMovementDetails(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []MovementResponse
	for _, movement := range movements {

		var usuario *UserResponse
		if movement.UsuarioID != nil {
			usuario = &UserResponse{
				UsuarioID: *movement.UsuarioID,
				Nombre:    *movement.UsuarioNombre,
			}
		}

		var proveedor *dto.Proveedore
		if movement.ProveedorID != nil {
			proveedor = &dto.Proveedore{
				ProveedorID: *movement.ProveedorID,
				Nombre:      *movement.Nombre,
				Telefono:    movement.Telefono,
				Email:       movement.Email,
			}
		}

		var detalles []MovementDetailResponse
		for _, detail := range details {
			if detail.MovimientoID == movement.MovimientoID {
				detalles = append(detalles, MovementDetailResponse{
					InsumoID: detail.InsumoID,
					Nombre:   *detail.Nombre,
					Tipo:     string(*detail.Tipo),
					Cantidad: detail.Cantidad,
				})
			}
		}

		response = append(response, MovementResponse{
			MovimientoID: movement.MovimientoID,
			Usuario:      usuario,
			Tipo:         string(movement.Tipo),
			Proveedor:    proveedor,
			Fecha: movement.Fecha,
			Detalles:     detalles,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// getMovementByID es el handler para obtener un movimiento específico por su ID, incluyendo sus detalles y proveedor asociado.
func (server *Server) getMovementByID(ctx *gin.Context) {
	var req getMovementByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	movement, err := server.dbtx.GetMovementsById(ctx, req.MovimientoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	details, err := server.dbtx.GetDetailById(ctx, req.MovimientoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var usuario *UserResponse
	if movement.UsuarioID != nil {
		usuario = &UserResponse{
			UsuarioID: *movement.UsuarioID,
			Nombre:    *movement.UsuarioNombre,
		}
	}

	var proveedor *dto.Proveedore
	if movement.ProveedorID != nil {
		proveedor = &dto.Proveedore{
			ProveedorID: *movement.ProveedorID,
			Nombre:      *movement.Nombre,
			Telefono:    movement.Telefono,
			Email:       movement.Email,
		}
	}

	var detalles []MovementDetailResponse
	for _, detail := range details {
		detalles = append(detalles, MovementDetailResponse{
			InsumoID: detail.InsumoID,
			Nombre:   *detail.Nombre,
			Tipo:     string(*detail.Tipo),
			Cantidad: detail.Cantidad,
		})
	}

	response := MovementResponse{
		MovimientoID: movement.MovimientoID,
		Usuario:      usuario,
		Tipo:         string(movement.Tipo),
		Proveedor:    proveedor,
		Fecha: movement.Fecha,
		Detalles:     detalles,
	}
	ctx.JSON(http.StatusOK, response)
}

// GET MOVEMENTS BY SUPPLY, SUPPLIER, USER
// Estas funciones son para obtener los movimientos asociados a un insumo, proveedor o usuario específico. Se pueden implementar de manera similar a getAllMovement, pero filtrando por el ID correspondiente.
func (server *Server) getMovementsBySupply(ctx *gin.Context) {
	var req getMovementByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	details, err := server.dbtx.GetDetailById(ctx, req.MovimientoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []MovementDetailResponse
	for _, detail := range details {
		if detail.MovimientoID == req.MovimientoID {
			response = append(response, MovementDetailResponse{
				InsumoID: detail.InsumoID,
				Nombre:   *detail.Nombre,
				Tipo:     string(*detail.Tipo),
				Cantidad: detail.Cantidad,
			})
		}
	}
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) getMovementsBySupplier(ctx *gin.Context) {
	var req getMovementByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	movement, err := server.dbtx.GetMovementsById(ctx, req.MovimientoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var suppliers *dto.Proveedore
	if movement.ProveedorID == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "No hay proveedor asociado a este movimiento"})
		return
	}
	if movement.ProveedorID != nil {
		suppliers = &dto.Proveedore{
			ProveedorID: *movement.ProveedorID,
			Nombre:      *movement.Nombre,
			Telefono:    movement.Telefono,
			Email:       movement.Email,
		}
	}

	ctx.JSON(http.StatusOK, suppliers)

}

func (server *Server) getMovementsByUser(ctx *gin.Context) {
	var req getMovementByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	movement, err := server.dbtx.GetMovementsById(ctx, req.MovimientoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var usuario *UserResponse
	if movement.UsuarioID == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "No hay usuario asociado a este movimiento"})
		return
	}
	if movement.UsuarioID != nil {
		usuario = &UserResponse{
			UsuarioID: *movement.UsuarioID,
			Nombre:    *movement.UsuarioNombre,
		}
	}
	ctx.JSON(http.StatusOK, usuario)

}
