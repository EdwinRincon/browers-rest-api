package handler

import (
	"net/http"
	"strconv"

	"github.com/EdwinRincon/browersfc-api/api/constants"
	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/api/service"
	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/gin-gonic/gin"
)

type ClassificationHandler struct {
	ClassificationService service.ClassificationService
}

func NewClassificationHandler(classificationService service.ClassificationService) *ClassificationHandler {
	return &ClassificationHandler{
		ClassificationService: classificationService,
	}
}

func (h *ClassificationHandler) CreateClassification(c *gin.Context) {
	var classification model.Classifications
	if err := c.ShouldBindJSON(&classification); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err := h.ClassificationService.CreateClassification(ctx, &classification)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to create classification", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusCreated, classification, "Classification created successfully")
}

func (h *ClassificationHandler) GetClassificationByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid ID", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	classification, err := h.ClassificationService.GetClassificationByID(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusNotFound, constants.ErrClassificationNotFound.Error(), ""), false)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, classification, "Classification retrieved successfully")
}

func (h *ClassificationHandler) ListClassifications(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid page number", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	classifications, err := h.ClassificationService.ListClassifications(ctx, page)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to list classifications", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, classifications, "Classifications listed successfully")
}

func (h *ClassificationHandler) UpdateClassification(c *gin.Context) {
	var classification model.Classifications
	if err := c.ShouldBindJSON(&classification); err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, constants.ErrInvalidInput, err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err := h.ClassificationService.UpdateClassification(ctx, &classification)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to update classification", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, classification, "Classification updated successfully")
}

func (h *ClassificationHandler) DeleteClassification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid ID", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	err = h.ClassificationService.DeleteClassification(ctx, id)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to delete classification", err.Error()), true)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, nil, "Classification deleted successfully")
}

func (h *ClassificationHandler) GetClassificationBySeason(c *gin.Context) {
	seasonID, err := strconv.ParseUint(c.Param("seasonID"), 10, 8)
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusBadRequest, "Invalid season ID", err.Error()), true)
		return
	}

	ctx := c.Request.Context()
	classifications, err := h.ClassificationService.GetClassificationBySeason(ctx, uint8(seasonID))
	if err != nil {
		helper.HandleError(c, helper.NewAppError(http.StatusInternalServerError, "Failed to retrieve classification by season", err.Error()), false)
		return
	}

	helper.HandleSuccess(c, http.StatusOK, classifications, "Classification by season retrieved successfully")
}
