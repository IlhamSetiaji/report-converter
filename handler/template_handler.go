package handler

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/IlhamSetiaji/report-converter/config"
	"github.com/IlhamSetiaji/report-converter/entity"
	"github.com/IlhamSetiaji/report-converter/logger"
	"github.com/IlhamSetiaji/report-converter/request"
	"github.com/IlhamSetiaji/report-converter/usecase"
	"github.com/IlhamSetiaji/report-converter/utils"
	"github.com/IlhamSetiaji/report-converter/validator"
	"github.com/gin-gonic/gin"
	"github.com/nguyenthenguyen/docx"
)

type ITemplateHandler interface {
	CreateTemplate(ctx *gin.Context)
	FindAllTemplate(ctx *gin.Context)
	FindTemplateByID(ctx *gin.Context)
	DeleteTemplateByID(ctx *gin.Context)
	GeneratePDF(ctx *gin.Context)
}

type TemplateHandler struct {
	templateUseCase usecase.ITemplateUseCase
	logger          logger.Logger
	validator       validator.Validator
	config          config.Config
}

func NewTemplateHandler(
	templateUseCase usecase.ITemplateUseCase,
	logger logger.Logger,
	validator validator.Validator,
	config config.Config,
) ITemplateHandler {
	return &TemplateHandler{
		templateUseCase: templateUseCase,
		logger:          logger,
		validator:       validator,
		config:          config,
	}
}

func (h *TemplateHandler) CreateTemplate(ctx *gin.Context) {
	h.logger.GetLogger().Info("Creating template")
	var req request.TemplateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.GetLogger().Error("Failed to bind JSON", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.validator.GetValidator().Struct(req); err != nil {
		h.logger.GetLogger().Error("Validation error", err)
		utils.BadRequestResponse(ctx, "Validation error", err.Error())
		return
	}

	if req.File != nil {
		timestamp := time.Now().UnixNano()
		filePath := "storage/templates/" + strconv.FormatInt(timestamp, 10) + "_" + req.File.Filename
		if err := ctx.SaveUploadedFile(req.File, filePath); err != nil {
			h.logger.GetLogger().Error("failed to save cover file: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save cover file", err.Error())
			return
		}

		req.File = nil
		req.Path = filePath
	}

	templateResponse, err := h.templateUseCase.CreateTemplate(&req)
	if err != nil {
		h.logger.GetLogger().Error("Failed to create template", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create template", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Template created successfully", templateResponse)
}

func (h *TemplateHandler) FindAllTemplate(ctx *gin.Context) {
	h.logger.GetLogger().Info("Finding all templates")
	templates, err := h.templateUseCase.FindAllTemplate()
	if err != nil {
		h.logger.GetLogger().Error("Failed to find all templates", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all templates", err.Error())
		return
	}

	if len(templates) == 0 {
		utils.SuccessResponse(ctx, http.StatusOK, "No templates found", nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Templates found successfully", templates)
}

func (h *TemplateHandler) FindTemplateByID(ctx *gin.Context) {
	h.logger.GetLogger().Info("Finding template by ID")
	id := ctx.Param("id")
	template, err := h.templateUseCase.FindTemplateByID(id)
	if err != nil {
		h.logger.GetLogger().Error("Failed to find template by ID", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find template by ID", err.Error())
		return
	}

	if template == nil {
		utils.SuccessResponse(ctx, http.StatusOK, "Template not found", nil)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Template found successfully", template)
}

func (h *TemplateHandler) DeleteTemplateByID(ctx *gin.Context) {
	h.logger.GetLogger().Info("Deleting template by ID")
	id := ctx.Param("id")
	err := h.templateUseCase.DeleteTemplateByID(id)
	if err != nil {
		h.logger.GetLogger().Error("Failed to delete template by ID", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete template by ID", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Template deleted successfully", nil)
}

func (h *TemplateHandler) GeneratePDF(c *gin.Context) {
	var req request.GeneratePDFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.GetLogger().Error("Failed to bind JSON", err)
		utils.BadRequestResponse(c, "Invalid request", err.Error())
		return
	}

	template, err := h.templateUseCase.FindTemplateByID(req.TemplateID)
	if err != nil {
		h.logger.GetLogger().Error("Failed to find template by ID", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find template", err.Error())
		return
	}

	if template == nil {
		h.logger.GetLogger().Error("Template not found")
		utils.ErrorResponse(c, http.StatusNotFound, "Template not found", "Template not found")
		return
	}

	if template.TemplateType != string(entity.TemplateTypeDocx) {
		h.logger.GetLogger().Error("Invalid template type")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid template type", "Invalid template type")
		return
	}

	templatePath := template.PathOriginal
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		h.logger.GetLogger().Error("Template file does not exist", err)
		utils.ErrorResponse(c, http.StatusNotFound, "Template file not found", "Template file not found")
		return
	}

	data := make(map[string]string)
	for key, value := range req.Data {
		if strValue, ok := value.(string); ok {
			data[key] = strValue
		} else if intValue, ok := value.(float64); ok {
			data[key] = strconv.FormatFloat(intValue, 'f', -1, 64)
		} else {
			h.logger.GetLogger().Error("Invalid data type for key", key)
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid data type", fmt.Sprintf("Invalid data type for key %s", key))
			return
		}
	}

	// Process the document
	pdfPath, err := h.processDocument(templatePath, data)
	if err != nil {
		h.logger.GetLogger().Error("Failed to process document ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process document", err.Error())
		return
	}
	defer os.Remove(pdfPath)

	// Send the PDF as response
	c.File(pdfPath)
}

func (h *TemplateHandler) processDocument(templatePath string, data map[string]string) (string, error) {
	// Read the docx file
	r, err := docx.ReadDocxFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read document: %v", err)
	}

	docxContent := r.Editable()

	// Replace all variables in the content
	content := docxContent.GetContent()
	for key, value := range data {
		templateVar := "{{." + key + "}}"
		content = strings.ReplaceAll(content, templateVar, value)
	}

	// Update the content
	docxContent.SetContent(content)

	// Ensure the directory for generated PDFs exists
	generatedPDFDir := "storage/generated_pdf"
	if err := os.MkdirAll(generatedPDFDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", generatedPDFDir, err)
	}

	// Save the modified DOCX
	modifiedNamePath := "modified_" + filepath.Base(templatePath)
	modifiedDocxPath := filepath.Join(generatedPDFDir, modifiedNamePath)
	err = docxContent.WriteToFile(modifiedDocxPath)
	if err != nil {
		return "", fmt.Errorf("failed to save modified document: %v", err)
	}
	defer os.Remove(modifiedDocxPath)

	// Convert to PDF using LibreOffice
	h.logger.GetLogger().Info("Converting to PDF: ", modifiedDocxPath)
	err = convertToPDF(modifiedDocxPath, generatedPDFDir)
	if err != nil {
		return "", fmt.Errorf("failed to convert to PDF: %v", err)
	}

	// Construct the expected PDF file path
	pdfFileName := strings.TrimSuffix(filepath.Base(modifiedDocxPath), filepath.Ext(modifiedDocxPath)) + ".pdf"
	pdfPath := filepath.Join(generatedPDFDir, pdfFileName)

	// Verify the PDF was created
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return "", fmt.Errorf("PDF file was not created")
	}

	return pdfPath, nil
}

func convertToPDF(docxPath, outputDir string) error {
	// Try both direct command and container-specific paths
	loPaths := []string{
		"/usr/bin/soffice", // Linux default
		"/usr/local/bin/soffice",
		"/opt/libreoffice/program/soffice",
	}

	if runtime.GOOS == "windows" {
		loPaths = append(loPaths, []string{
			`C:\Program Files\LibreOffice\program\soffice.exe`,
			`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
		}...)
	}

	var loPath string
	for _, path := range loPaths {
		if _, err := os.Stat(path); err == nil {
			loPath = path
			break
		}
	}

	if loPath == "" {
		return fmt.Errorf("LibreOffice not found in standard locations")
	}

	// Add container-specific environment variables
	env := os.Environ()
	env = append(env, "HOME=/tmp") // LibreOffice needs a home directory

	cmd := exec.Command(
		loPath,
		"--headless",
		"--convert-to", "pdf",
		"--outdir", outputDir,
		docxPath,
	)
	cmd.Env = env

	// Log the paths being used
	fmt.Printf("Using LibreOffice at: %s\n", loPath)
	fmt.Printf("Input DOCX path: %s\n", docxPath)
	fmt.Printf("Output directory: %s\n", outputDir)

	// Set timeout for the conversion
	done := make(chan error, 1)
	go func() {
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("LibreOffice output: %s\n", string(output))
			done <- fmt.Errorf("PDF conversion failed: %v, output: %s", err, string(output))
			return
		}
		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			return err
		}
	case <-time.After(30 * time.Second):
		// Kill the process if it takes too long
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return fmt.Errorf("PDF conversion timed out after 30 seconds")
	}

	return nil
}
