package main

import (
	"github.com/IlhamSetiaji/report-converter/config"
	"github.com/IlhamSetiaji/report-converter/database"
	"github.com/IlhamSetiaji/report-converter/entity"
	"github.com/IlhamSetiaji/report-converter/logger"
)

func main() {
	config := config.GetConfig()
	logger := logger.NewLogger()
	db := database.NewPostgresDatabase(config)

	// Initialize the database connection
	if err := db.GetDb().AutoMigrate(&entity.Template{}); err != nil {
		logger.GetLogger().Fatal("Failed to migrate database", err)
	}
}
