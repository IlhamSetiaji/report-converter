package usecase

import (
	"github.com/IlhamSetiaji/report-converter/dto"
	"github.com/IlhamSetiaji/report-converter/entity"
	"github.com/IlhamSetiaji/report-converter/repository"
	"github.com/IlhamSetiaji/report-converter/request"
	"github.com/IlhamSetiaji/report-converter/response"
	"github.com/google/uuid"
)

type ITemplateUseCase interface {
	CreateTemplate(template *request.TemplateRequest) (*response.TemplateResponse, error)
	FindAllTemplate() ([]*response.TemplateResponse, error)
	FindTemplateByID(id string) (*response.TemplateResponse, error)
	DeleteTemplateByID(id string) error
}

type TemplateUseCase struct {
	templateRepository repository.ITemplateRepository
	templateDTO        dto.ITemplateDTO
}

func NewTemplateUseCase(templateRepository repository.ITemplateRepository, templateDTO dto.ITemplateDTO) ITemplateUseCase {
	return &TemplateUseCase{
		templateRepository: templateRepository,
		templateDTO:        templateDTO,
	}
}

func (t *TemplateUseCase) CreateTemplate(template *request.TemplateRequest) (*response.TemplateResponse, error) {
	ent := &entity.Template{
		Name:         template.Name,
		TemplateType: entity.TemplateType(template.TemplateType),
		Path:         template.Path,
	}

	createdTemplate, err := t.templateRepository.CreateTemplate(ent)
	if err != nil {
		return nil, err
	}

	return t.templateDTO.ConvertEntityToResponse(createdTemplate), nil
}

func (t *TemplateUseCase) FindAllTemplate() ([]*response.TemplateResponse, error) {
	templates, err := t.templateRepository.FindAllTemplate()
	if err != nil {
		return nil, err
	}

	var templateResponses []*response.TemplateResponse
	for _, template := range templates {
		templateResponses = append(templateResponses, t.templateDTO.ConvertEntityToResponse(&template))
	}

	return templateResponses, nil
}

func (t *TemplateUseCase) FindTemplateByID(id string) (*response.TemplateResponse, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	ent, err := t.templateRepository.FindTemplateByID(parsedId)
	if err != nil {
		return nil, err
	}
	if ent == nil {
		return nil, nil
	}

	return t.templateDTO.ConvertEntityToResponse(ent), nil
}

func (t *TemplateUseCase) DeleteTemplateByID(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = t.templateRepository.DeleteTemplateByID(parsedId)
	if err != nil {
		return err
	}
	return nil
}
