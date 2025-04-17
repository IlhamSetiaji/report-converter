package response

type TemplateResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	TemplateType string `json:"template_type"`
	Path         string `json:"path"`
	PathOriginal string `json:"path_original"`
}
