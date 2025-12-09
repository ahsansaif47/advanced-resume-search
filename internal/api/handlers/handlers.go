package handlers

import (
	"github.com/ahsansaif47/advanced-resume/internal/api/controllers"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service controllers.IWeaviateService
}

// UploadResume godoc
// @Summary 	Add a new resume
// @Description Stores a single resume object in the database/vector DB
// @Tags 		Resume
// @Accept 		json
// @Produce 	json
// @Param 		resume 	body 	object 	true 	"Resume data (dynamic fields)"
// @Success 	201 	{object} 	map[string]any 		"ID of inserted resume"
// @Failure 	400 	{object} 	map[string]string 	"Invalid request body"
// @Failure 	500 	{object} 	map[string]string 	"Internal Server Error"
// @Router 		/upload [post]
func (h *Handler) UploadResume(ctx *fiber.Ctx) error {

	file, err := ctx.FormFile("document")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	status, id, err := h.service.AddResumeToDB(ctx, file) // 10s

	return ctx.Status(status).JSON(fiber.Map{
		"ID": id,
	})
}

// BatchUploadResume 	godoc
// @Summary 		Add multiple resumes (batch)
// @Description 	Inserts multiple resumes using Weaviate batch API
// @Tags 			Resume
// @Accept json
// @Produce json
// @Param resumes body []object true "Array of resume objects"
// @Success 201 {object} map[string]string "Added resumes"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /upload/batch [post]
func (h *Handler) BatchUploadResume(ctx *fiber.Ctx) error {
	var req []map[string]any

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	status, err := h.service.BatchUploadResume(req)
	if err != nil {
		return ctx.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(status).JSON(fiber.Map{
		"Message": "Added resumes",
	})
}

func (h *Handler) VectorSearch(ctx *fiber.Ctx) error {
	query := ctx.Query("query")
	if query == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": "Query cann't be empty",
		})
	}

	status, data, err := h.service.VectorSearch(query)
	if err != nil {
		return ctx.Status(status).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return ctx.Status(status).JSON(fiber.Map{
		"Results": data,
	})
}
