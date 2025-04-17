package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateType string

const (
	TemplateTypeExcel TemplateType = "excel"
	TemplateTypeDocx  TemplateType = "docx"
)

type Template struct {
	gorm.Model   `json:"-"`
	ID           uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	Name         string       `json:"name" gorm:"type:varchar(255);not null"`
	TemplateType TemplateType `json:"template_type" gorm:"type:varchar(255);not null"`
	Path         string       `json:"path" gorm:"type:text;not null"`
}

func (t *Template) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t.CreatedAt = time.Now().In(loc)
	t.UpdatedAt = time.Now().In(loc)
	return nil
}

func (t *Template) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t.UpdatedAt = time.Now().In(loc)
	return nil
}

func (Template) TableName() string {
	return "templates"
}
