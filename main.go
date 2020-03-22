package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	authHandler "go-echo-api/auth/delivery/http"
	authService "go-echo-api/auth/usecase"
	"go-echo-api/config"
	"go-echo-api/infrastructure/database"
	"go-echo-api/infrastructure/validator"
	jwtMiddleware "go-echo-api/middleware"
	//"go-echo-api/oauth"
	"go-echo-api/services"
	userHandler "go-echo-api/user/delivery/http"
	userService "go-echo-api/user/usecase"

	"net/http"
	"os"
	"path/filepath"
)

func init() {
	fileExecutable, _ := os.Executable()
	basePath, _ := filepath.Split(fileExecutable)
	if os.Getenv("APP_ENV") != "production" {
		basePath = ""
	}
	_ = godotenv.Load(basePath + ".env")
}
//var (
//	// OauthService ...
//	OauthService oauth.ServiceInterface
//)

func main() {

	e := echo.New()
	e.Validator = validator.NewValidator()
	db := database.New()
	database.AutoMigrate(db)
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	api := e.Group("/api")
	v1 := api.Group("/v1")

	cnf := config.NewConfig(true, false, "consul")
	//OauthService = oauth.NewService(cnf, db)

	// start oauth2 services
	services.Init(cnf, db)
	defer services.Close()

	services.OauthService.RegisterRoutes(e, "/v1/oauth")

	//AuthController
	authController := authHandler.NewAuthController(authService.NewAuthService(db))
	auth := v1.Group("/auth")
	auth.POST("/token", authController.Login)
	auth.POST("/register", authController.Register)
	auth.POST("/refresh-token", authController.RefreshToken)

	//UserController
	userController := userHandler.NewUserController(userService.NewUserService(db))
	user := v1.Group("/user")
	user.GET("", userController.FindAll, jwtMiddleware.IsLoggedIn)
	user.GET("/:id", userController.FindById, jwtMiddleware.IsLoggedIn)
	user.POST("", userController.Store, jwtMiddleware.IsLoggedIn)
	user.PUT("/:id", userController.Update, jwtMiddleware.IsLoggedIn)
	user.DELETE("/:id", userController.Delete, jwtMiddleware.IsLoggedIn)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(os.Getenv("APP_PORT")))
}
