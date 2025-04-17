package main

import (
	"github.com/IlhamSetiaji/report-converter/config"
	"github.com/IlhamSetiaji/report-converter/database"
	"github.com/IlhamSetiaji/report-converter/logger"
	"github.com/IlhamSetiaji/report-converter/server"
	"github.com/IlhamSetiaji/report-converter/validator"
)

func main() {
	// Initialize the application components (config, logger, database, server)
	config := config.GetConfig()
	logger := logger.NewLogger()
	db := database.NewPostgresDatabase(config)
	validator := validator.NewValidatorV10(config)
	server := server.NewGinServer(db, *config, logger, validator)

	// Start the server
	server.Start()
}
