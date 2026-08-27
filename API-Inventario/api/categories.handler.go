package api

import (
	"rest/dto"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Requests para categorias
type createCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type updateCategoryRequest struct {
	CategoryID int32  `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

type deleteCategoryRequestJSON struct {
	CategoryID int32 `json:"category_id" binding:"required"`
}

type deleteCategoryRequestURI struct {
	CategoryID int32 `uri:"category_id" binding:"required"`
}

type getCategoryByIDRequest struct {
	CategoryID int32 `uri:"category_id" binding:"required"`
}

// createCategory es el handler para crear una nueva categoria.
func (server *Server) createCategory(c *gin.Context) {
	var req createCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.dbtx.CreateCategory(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, result)
}

// updateCategory es el handler para actualizar una categoria.
func (server *Server) updateCategory(c *gin.Context) {
	var req updateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := dto.UpdateCategoryParams{
		Name:       req.Name,
		CategoryID: req.CategoryID,
	}

	result, err := server.dbtx.UpdateCategory(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

// deleteCategoryJSON elimina una categoria por JSON.
func (server *Server) deleteCategoryJSON(c *gin.Context) {
	var req deleteCategoryRequestJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.dbtx.DeleteCategory(c, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

// deleteCategoryURI elimina una categoria por URI.
func (server *Server) deleteCategoryURI(c *gin.Context) {
	var req deleteCategoryRequestURI
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.dbtx.DeleteCategory(c, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

// getAllCategories es el handler para obtener todas las categorias.
func (server *Server) getAllCategories(c *gin.Context) {
	categories, err := server.dbtx.GetAllCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, categories)
}

// getCategoryByID es el handler para obtener una categoria por ID.
func (server *Server) getCategoryByID(c *gin.Context) {
	var req getCategoryByIDRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.dbtx.GetCategoryById(c, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, category)
}
