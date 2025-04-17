package repository

import (
	"errors"

	"github.com/IlhamSetiaji/report-converter/database"
	"github.com/IlhamSetiaji/report-converter/entity"
	"github.com/IlhamSetiaji/report-converter/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITemplateRepository interface {
	CreateTemplate(template *entity.Template) (*entity.Template, error)
	FindAllTemplate() ([]entity.Template, error)
	FindTemplateByID(id uuid.UUID) (*entity.Template, error)
	DeleteTemplateByID(id uuid.UUID) error
}

type TemplateRepository struct {
	db     database.Database
	logger logger.Logger
}

func NewTemplateRepository(db database.Database, logger logger.Logger) ITemplateRepository {
	return &TemplateRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TemplateRepository) CreateTemplate(template *entity.Template) (*entity.Template, error) {
	err := r.db.GetDb().Create(template).Error
	if err != nil {
		r.logger.GetLogger().Error("Failed to create template", err)
		return nil, err
	}
	return template, nil
}

func (r *TemplateRepository) FindAllTemplate() ([]entity.Template, error) {
	var templates []entity.Template
	err := r.db.GetDb().Find(&templates).Error
	if err != nil {
		r.logger.GetLogger().Error("Failed to find all templates", err)
		return nil, err
	}
	return templates, nil
}

func (r *TemplateRepository) FindTemplateByID(id uuid.UUID) (*entity.Template, error) {
	var template entity.Template
	err := r.db.GetDb().First(&template, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.GetLogger().Error("Template not found", err)
			return nil, nil
		}
		r.logger.GetLogger().Error("Failed to find template by ID", err)
	}
	return &template, nil
}

func (r *TemplateRepository) DeleteTemplateByID(id uuid.UUID) error {
	var template entity.Template
	err := r.db.GetDb().First(&template, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			r.logger.GetLogger().Error("Template not found", err)
			return nil
		}
		r.logger.GetLogger().Error("Failed to find template by ID", err)
		return err
	}

	err = r.db.GetDb().Delete(&template).Error
	if err != nil {
		r.logger.GetLogger().Error("Failed to delete template", err)
		return err
	}
	return nil
}
