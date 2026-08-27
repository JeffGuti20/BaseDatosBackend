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
	Name       string  `json:"name" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Stock      int32   `json:"stock" binding:"required"`
	CategoryID int32   `json:"category_id" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	Image      string  `json:"image" binding:"required"`
}

type updateSupplyRequest struct {
	SupplyID   int32   `json:"supply_id" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Stock      int32   `json:"stock" binding:"required"`
	CategoryID int32   `json:"category_id" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	Image      string  `json:"image" binding:"required"`
}

type deleteSupplyRequestJSON struct {
	SupplyID int32 `json:"supply_id" binding:"required"`
}

type deleteSupplyRequestURI struct {
	SupplyID int32 `uri:"supply_id" binding:"required"`
}

type getSupplyByIDRequest struct {
	SupplyID int32 `uri:"supply_id" binding:"required"`
}

// Responses para insumos
type categoryResponse struct {
	CategoryID int32  `json:"category_id"`
	Name       string `json:"name"`
}

type supplyResponse struct {
	SupplyID int32            `json:"supply_id"`
	Name     string           `json:"name"`
	Type     dto.SuppliesType `json:"type"`
	Stock    int32            `json:"stock"`
	Price    float64          `json:"price"`
	Category categoryResponse `json:"category"`
	Image    string           `json:"image"`
}

type supplyImageRequest struct {
	Filename string `uri:"filename" binding:"required"`
}

// createSupply es el handler para crear un nuevo movimiento.
func (server *Server) createSupply(ctx *gin.Context) {
	var req createSupplyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := dto.CreateSupplyParams{
		Name:       req.Name,
		Type:       dto.SuppliesType(req.Type),
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
		Price:      req.Price,
		Image:      req.Image,
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
		Name:       req.Name,
		Type:       dto.SuppliesType(req.Type),
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
		Price:      req.Price,
		SupplyID:   req.SupplyID,
		Image:      req.Image,
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

	result, err := server.dbtx.DeleteSupply(ctx, req.SupplyID)
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

	result, err := server.dbtx.DeleteSupply(ctx, req.SupplyID)
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

		category := categoryResponse{
			CategoryID: supply.CategoryID,
			Name:       supply.CategoryName.String,
		}

		response = append(response, supplyResponse{
			SupplyID: supply.SupplyID,
			Name:     supply.Name,
			Type:     supply.Type,
			Stock:    supply.Stock,
			Price:    supply.Price,
			Category: category,
			Image: supply.Image,
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

	supply, err := server.dbtx.GetSupplyById(ctx, req.SupplyID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	category := categoryResponse{
		CategoryID: supply.CategoryID,
		Name:       supply.CategoryName.String,
	}

	response := supplyResponse{
		SupplyID: supply.SupplyID,
		Name:     supply.Name,
		Type:     supply.Type,
		Stock:    supply.Stock,
		Price:    supply.Price,
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

	supplies, err := server.dbtx.GetSupplyById(ctx, req.SupplyID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	response := categoryResponse{
		CategoryID: supplies.CategoryID,
		Name:       supplies.CategoryName.String,
	}

	ctx.JSON(http.StatusOK, response)
}

/*
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
	upDir := filepath.Join(wd, "utils", "data", "supplies")

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

	wd, _ := os.Getwd()

	fileUrl := filepath.Join(
		wd,
		"utils",
		"data",
		"supplies",
		req.Filename,
	)

	ctx.File(fileUrl)
}

func (server *Server) deleteSupplyImg(ctx *gin.Context) {

	var req supplyImageRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	wd, _ := os.Getwd()

	filePath := filepath.Join(
		wd,
		"utils",
		"data",
		"supplies",
		req.Filename,
	)

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
*/

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