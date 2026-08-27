package api

import (
	"rest/dto"

	"github.com/gin-gonic/gin"
	"net/http"
)

// Request para categorias
type createCategoryRequest struct {
	Nombre string `json:"name" binding:"required"`
}

type updateCategoryRequest struct {
	CategoriaID int32  `json:"category_id" binding:"required"`
	Nombre      string `json:"name" binding:"required"`
}

type deleteCategoryRequestJSON struct {
	CategoriaID int32 `json:"category_id" binding:"required"`
}

type deleteCategoryRequestURI struct {
	CategoriaID int32 `uri:"category_id" binding:"required"`
}

type getCategoryByIDRequest struct {
	CategoriaID int32 `uri:"category_id" binding:"required"`
}

// createCategory es el handler para crear una nueva categoria.
func (server *Server) createCategory(c *gin.Context) {
	var req createCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.dbtx.CreateCategory(c, req.Nombre)
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
		CategoriaID: req.CategoriaID,
		Nombre:      req.Nombre,
	}

	result, err := server.dbtx.UpdateCategory(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

// deleteCategory elimina una categoria por JSON.
func (server *Server) deleteCategoryJSON(c *gin.Context) {
	var req deleteCategoryRequestJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.dbtx.DeleteCategory(c, req.CategoriaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

// deleteCategory elimina una categoria por URI.
func (server *Server) deleteCategoryURI(c *gin.Context) {
	var req deleteCategoryRequestURI
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.dbtx.DeleteCategory(c, req.CategoriaID)
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

	category, err := server.dbtx.GetCategoryById(c, req.CategoriaID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, category)
}
