package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"rest/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Requests para insumos
type createSupplyRequest struct {
	Nombre      string `json:"name" binding:"required"`
	Tipo        string `json:"type" binding:"required"`
	Cantidad    int32  `json:"quantity" binding:"required"`
	Image       string `json:"image" binding:"required"`
	CategoriaID *int32 `json:"category_id" binding:"required"`
}

type updateSupplyRequest struct {
	InsumoID    int32  `json:"supply_id" binding:"required"`
	Nombre      string `json:"name" binding:"required"`
	Tipo        string `json:"type" binding:"required"`
	Cantidad    int32  `json:"quantity" binding:"required"`
	Image       string `json:"image" binding:"required"`
	CategoriaID *int32 `json:"category_id" binding:"required"`
}

type deleteSupplyRequestJSON struct {
	InsumoID int32 `json:"supply_id" binding:"required"`
}

type deleteSupplyRequestURI struct {
	InsumoID int32 `uri:"supply_id" binding:"required"`
}

type getSupplyByIDRequest struct {
	InsumoID int32 `uri:"supply_id" binding:"required"`
}

type supplyImageRequest struct {
	Filename string `uri:"filename" binding:"required"`
}

// Responses para insumos
type categoryResponse struct {
	CategoriaID int32  `json:"category_id"`
	Nombre      string `json:"name"`
}

type supplyResponse struct {
	InsumoID int32            `json:"supply_id"`
	Nombre   string           `json:"name"`
	Tipo     dto.InsumosTipo  `json:"type"`
	Cantidad int32            `json:"quantity"`
	Image    string           `json:"image"`
	Category categoryResponse `json:"category"`
}

// createSupply es el handler para crear un nuevo movimiento.
func (server *Server) createSupply(ctx *gin.Context) {
	var req createSupplyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := dto.CreateSupplyParams{
		Nombre:      req.Nombre,
		Tipo:        dto.InsumosTipo(req.Tipo),
		Cantidad:    req.Cantidad,
		CategoriaID: req.CategoriaID,
		Image:       req.Image,
	}

	result, err := server.dbtx.CreateSupply(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// updateSupply es el handler para actualizar un movimiento.
func (server *Server) updateSupply(ctx *gin.Context) {
	var req updateSupplyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := dto.UpdateSupplyParams{
		InsumoID:    req.InsumoID,
		Nombre:      req.Nombre,
		Tipo:        dto.InsumosTipo(req.Tipo),
		Cantidad:    req.Cantidad,
		Image:       req.Image,
		CategoriaID: req.CategoriaID,
	}

	result, err := server.dbtx.UpdateSupply(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// deleteSupply elimina un insumo por JSON.
func (server *Server) deleteSupplyJSON(ctx *gin.Context) {
	var req deleteSupplyRequestJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.dbtx.DeleteSupply(ctx, req.InsumoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// deleteSupply elimina un insumo por URI.
func (server *Server) deleteSupplyURI(ctx *gin.Context) {
	var req deleteSupplyRequestURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.dbtx.DeleteSupply(ctx, req.InsumoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// getAllSupplies obtiene todos los insumos.
func (server *Server) getAllSupplies(ctx *gin.Context) {
	supplies, err := server.dbtx.GetAllSupplies(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := []supplyResponse{}
	for _, supply := range supplies {

		var category categoryResponse
		if supply.CategoriaID != nil {
			category = categoryResponse{
				CategoriaID: *supply.CategoriaID,
				Nombre:      supply.CategoriaNombre,
			}
		}

		response = append(response, supplyResponse{
			InsumoID: supply.InsumoID,
			Nombre:   supply.Nombre,
			Tipo:     supply.Tipo,
			Cantidad: supply.Cantidad,
			Image:    supply.Image,
			Category: category,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// getSupplyByID obtiene un insumo por su ID.
func (server *Server) getSupplyByID(ctx *gin.Context) {
	var req getSupplyByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supply, err := server.dbtx.GetSupplyById(ctx, req.InsumoID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var category categoryResponse
	if supply.CategoriaID != nil {
		category = categoryResponse{
			CategoriaID: *supply.CategoriaID,
			Nombre:      supply.CategoriaNombre,
		}
	}

	response := supplyResponse{
		InsumoID: supply.InsumoID,
		Nombre:   supply.Nombre,
		Tipo:     supply.Tipo,
		Cantidad: supply.Cantidad,
		Image:    supply.Image,
		Category: category,
	}

	ctx.JSON(http.StatusOK, response)
}

// getSuppliesByCategory obtiene los insumos por su categoría.
func (server *Server) getSuppliesByCategory(ctx *gin.Context) {
	var req getSupplyByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supplies, err := server.dbtx.GetSupplyById(ctx, req.InsumoID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	response := categoryResponse{
		CategoriaID: *supplies.CategoriaID,
		Nombre:      supplies.CategoriaNombre,
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) uploadSupplyImg(ctx *gin.Context) {

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
	upDir := "/data/supplies"

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

func (server *Server) downloadSupplyImg(ctx *gin.Context) {

	var req supplyImageRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fileUrl := "/data/supplies/" + req.Filename

	ctx.File(fileUrl)
}

func (server *Server) deleteSupplyImg(ctx *gin.Context) {

    var req supplyImageRequest
    if err := ctx.ShouldBindUri(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    filePath := filepath.Join("/data/supplies", req.Filename)

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