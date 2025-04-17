package server

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/report-converter/config"
	"github.com/IlhamSetiaji/report-converter/database"
	"github.com/IlhamSetiaji/report-converter/dto"
	"github.com/IlhamSetiaji/report-converter/handler"
	"github.com/IlhamSetiaji/report-converter/logger"
	"github.com/IlhamSetiaji/report-converter/repository"
	"github.com/IlhamSetiaji/report-converter/usecase"
	"github.com/IlhamSetiaji/report-converter/validator"
	"github.com/gin-gonic/gin"
)

type ginServer struct {
	app       *gin.Engine
	db        database.Database
	conf      config.Config
	log       logger.Logger
	validator validator.Validator
}

func NewGinServer(db database.Database, conf config.Config, log logger.Logger, validator validator.Validator) Server {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	return &ginServer{
		app:       app,
		db:        db,
		conf:      conf,
		log:       log,
		validator: validator,
	}
}

func (g *ginServer) Start() {
	g.app.Static("/storage", "./storage")
	g.app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("App-Name", g.conf.Server.Name)
	})
	g.app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the User Service API",
			"status":  "OK",
		})
	})
	g.app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Service is running",
			"status":  "OK",
		})
	})

	g.initializeTemplateHandler()

	g.log.GetLogger().Info("Server started on port " + strconv.Itoa(g.conf.Server.Port))
	g.app.Run(":" + strconv.Itoa(g.conf.Server.Port))
}

func (g *ginServer) GetApp() *gin.Engine {
	return g.app
}

func (g *ginServer) initializeTemplateHandler() {
	templateRepository := repository.NewTemplateRepository(g.db, g.log)
	templateDTO := dto.NewTemplateDTO(g.conf, g.log)
	templateUseCase := usecase.NewTemplateUseCase(templateRepository, templateDTO)
	templateHandler := handler.NewTemplateHandler(templateUseCase, g.log, g.validator, g.conf)

	templateRoutes := g.app.Group("/api/v1/templates")
	templateRoutes.POST("/store", templateHandler.CreateTemplate)
	templateRoutes.POST("/generate-pdf", templateHandler.GeneratePDF)
	templateRoutes.GET("/", templateHandler.FindAllTemplate)
	templateRoutes.GET("/:id", templateHandler.FindTemplateByID)
	templateRoutes.DELETE("/:id", templateHandler.DeleteTemplateByID)
}
