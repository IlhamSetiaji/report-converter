package request

import "mime/multipart"

type TemplateRequest struct {
	Name         string                `form:"name" validate:"required"`
	TemplateType string                `form:"template_type" validate:"required"`
	File         *multipart.FileHeader `form:"file" validate:"required"`
	Path         string                `form:"path" validate:"omitempty"`
}

type GeneratePDFRequest struct {
	TemplateID string                 `json:"template_id" validate:"required"`
	Data       map[string]interface{} `json:"data" validate:"required"`
}
