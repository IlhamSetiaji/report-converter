package validator

import (
	"github.com/IlhamSetiaji/report-converter/config"
	"github.com/go-playground/validator/v10"
)

type validatorV10 struct {
	ValidatorV10 *validator.Validate
}

func NewValidatorV10(conf *config.Config) Validator {
	validate := validator.New()
	validate.RegisterValidation("template_type", func(fl validator.FieldLevel) bool {
		templateType := fl.Field().String()
		switch templateType {
		case "excel", "pdf":
			return true
		default:
			return false
		}
	})
	return &validatorV10{
		ValidatorV10: validate,
	}
}

func (v *validatorV10) GetValidator() *validator.Validate {
	return v.ValidatorV10
}
