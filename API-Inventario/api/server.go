package api

/*
Puse algunas importaciones comentadas porque no se están usando, pero las dejo ahí
por si las necesitamos después.
Vamos a necesitarlas despues, pero solo las descomentamos cuando las necesitemos,
porque go no compila si hay importaciones que no se usan.
De momento para hacer los handlers y probarlos con postman,
no las necesitamos, pero luego si vamos a necesitar.
Y con ese codigo es suficiente para hacer los handlers y luego hacer la URL.
*/

import (
	"database/sql"

	"rest/dto"
	"rest/security"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	cors "github.com/itsjamie/gin-cors"
)

type Server struct {
	conn         *sql.DB
	dbtx         *dto.Queries
	router       *gin.Engine
	tokenBuilder security.Builder
}

func NewServer(dbtx *dto.Queries, conn *sql.DB, secret string) (*Server, error) {

	builder, err := security.NewPasetoBuilder(secret)
	if err != nil {
		return nil, err
	}

	server := &Server{
		conn:         conn,
		dbtx:         dbtx,
		tokenBuilder: builder,
	}
	router := gin.Default()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET,POST,PUT,DELETE",
		RequestHeaders:  "Origin,Authorization,Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	//RUTAS SIN MIDDLEWARE
	router.POST("api/v1/login", server.login)
	router.POST("api/v1/register", server.createUser)

	//RUTAS CON MIDDLEWARE
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenBuilder))

	//ROUTES CATEGORIES
	authRoutes.GET("api/v1/categories", server.getAllCategories)
	authRoutes.GET("api/v1/categories/:category_id", server.getCategoryByID)

	authRoutes.POST("api/v1/categories", server.createCategory)
	authRoutes.PUT("api/v1/categories", server.updateCategory)
	authRoutes.DELETE("api/v1/categories", server.deleteCategoryJSON)
	authRoutes.DELETE("api/v1/categories/:category_id", server.deleteCategoryURI)

	//ROUTES SUPPLIERS
	authRoutes.GET("api/v1/suppliers", server.getAllSuppliers)
	authRoutes.GET("api/v1/suppliers/:supplier_id", server.getSuppliersByID)

	authRoutes.POST("api/v1/suppliers", server.createSupplier)
	authRoutes.PUT("api/v1/suppliers", server.updateSupplier)
	authRoutes.DELETE("api/v1/suppliers", server.deleteSupplierJSON)
	authRoutes.DELETE("api/v1/suppliers/:supplier_id", server.deleteSupplierURI)

	//ROUTES CLIENTS
	authRoutes.GET("api/v1/clients", server.getAllClients)
	authRoutes.GET("api/v1/clients/:client_id", server.getClientsByID)

	authRoutes.POST("api/v1/clients", server.createClient)
	authRoutes.PUT("api/v1/clients", server.updateClient)
	authRoutes.DELETE("api/v1/clients", server.deleteClientJSON)
	authRoutes.DELETE("api/v1/clients/:client_id", server.deleteClientURI)

	// ROUTES SUPPLIES
	authRoutes.GET("api/v1/supplies", server.getAllSupplies)
	authRoutes.GET("api/v1/supplies/:supply_id", server.getSupplyByID)
	authRoutes.GET("api/v1/supplies/:supply_id/category", server.getSuppliesByCategory)

	authRoutes.POST("api/v1/supplies", server.createSupply)
	authRoutes.PUT("api/v1/supplies", server.updateSupply)
	authRoutes.DELETE("api/v1/supplies", server.deleteSupplyJSON)
	authRoutes.DELETE("api/v1/supplies/:supply_id", server.deleteSupplyURI)

	authRoutes.POST("api/v1/supplies/upload_image", server.uploadSupplyImg)
	authRoutes.GET("api/v1/supplies/download_image/:filename", server.downloadSupplyImg)
	authRoutes.DELETE("api/v1/supplies/delete_image/:filename", server.deleteSupplyImg)

	//ROUTES ROLES
	authRoutes.GET("api/v1/roles", server.getAllRoles)
	authRoutes.GET("api/v1/roles/:role_id", server.getRoleByID)

	authRoutes.POST("api/v1/roles", server.createRole)
	authRoutes.PUT("api/v1/roles", server.updateRole)
	authRoutes.DELETE("api/v1/roles", server.deleteRoleJSON)
	authRoutes.DELETE("api/v1/roles/:role_id", server.deleteRoleURI)

	//ROUTES PURCHASES
	authRoutes.GET("api/v1/purchases", server.getAllPurchase)

	authRoutes.POST("api/v1/purchases", server.createPurchases)
	authRoutes.PUT("api/v1/purchases", server.updatePurchases)
	authRoutes.DELETE("api/v1/purchases", server.deletePurchaseJSON)
	authRoutes.DELETE("api/v1/purchases/:purchase_id", server.deletePurchaseURI)

	/*
		//ROUTES MOVEMENTS
		authRoutes.GET("api/v1/movements", server.getAllMovement)
		authRoutes.GET("api/v1/movements/:movement_id", server.getMovementByID)
		authRoutes.GET("api/v1/movements/:movement_id/supplies", server.getMovementsBySupply)
		authRoutes.GET("api/v1/movements/:movement_id/supplier", server.getMovementsBySupplier)
		authRoutes.GET("api/v1/movements/:movement_id/user", server.getMovementsByUser)

		authRoutes.POST("api/v1/movements", server.createMovement)
		authRoutes.PUT("api/v1/movements", server.updateMovement)
		authRoutes.DELETE("api/v1/movements", server.deleteMovementJSON)
		authRoutes.DELETE("api/v1/movements/:movement_id", server.deleteMovementURI)

	*/
	//ROUTES USERS
	authRoutes.GET("api/v1/users", server.getAllUsers)
	authRoutes.GET("api/v1/users/:user_id", server.getUserByID)

	authRoutes.PUT("api/v1/users/:user_id/role/:role_id", server.updateUserRole)
	authRoutes.PUT("api/v1/users", server.updateUser)
	authRoutes.PUT("api/v1/profile", server.updateProfile)
	authRoutes.DELETE("api/v1/users", server.deleteUserJSON)
	authRoutes.DELETE("api/v1/users/:user_id", server.deleteUserURI)

	server.router = router
	return server, nil
}

func (server *Server) Start(url string) error {
	return server.router.Run(url)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
