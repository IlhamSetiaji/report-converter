package dto

import (
	"github.com/IlhamSetiaji/report-converter/config"
	"github.com/IlhamSetiaji/report-converter/entity"
	"github.com/IlhamSetiaji/report-converter/logger"
	"github.com/IlhamSetiaji/report-converter/response"
)

type ITemplateDTO interface {
	ConvertEntityToResponse(ent *entity.Template) *response.TemplateResponse
}

type TemplateDTO struct {
	config config.Config
	logger logger.Logger
}

func NewTemplateDTO(config config.Config, logger logger.Logger) ITemplateDTO {
	return &TemplateDTO{
		config: config,
		logger: logger,
	}
}

func (t *TemplateDTO) ConvertEntityToResponse(ent *entity.Template) *response.TemplateResponse {
	return &response.TemplateResponse{
		ID:           ent.ID.String(),
		Name:         ent.Name,
		TemplateType: string(ent.TemplateType),
		Path:         config.GetConfig().Server.Url + "/" + ent.Path,
		PathOriginal: ent.Path,
	}
}
