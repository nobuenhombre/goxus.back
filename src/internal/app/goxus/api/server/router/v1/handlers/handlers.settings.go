package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	settingsdomain "goxus/src/internal/app/goxus/domain/settings"
	"goxus/src/internal/pkg/db/goxus"
)

// UpdateUserSettingRequest is the request body for updating a user setting.
type UpdateUserSettingRequest struct {
	Value goxus.JSON `json:"value" binding:"required"`
}

// GetSettingsDefinitions handles GET /api/v1/entity/settings
// Returns all setting definitions enriched with type and group info.
func (h *HttpHandler) GetSettingsDefinitions(c *gin.Context) {
	ctx, ok := authContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	definitions, err := h.Domain.GetSettingsDefinitions(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if definitions == nil {
		definitions = []*settingsdomain.SettingsDefinition{}
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"data":    definitions,
	})
}

// GetUserSettings handles GET /api/v1/entity/user/:id/settings
// Returns the current settings for the specified user.
func (h *HttpHandler) GetUserSettings(c *gin.Context) {
	ctx, ok := authContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	settings, err := h.Domain.GetUserSettings(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"data":    settings,
	})
}

// UpsertUserSetting handles PUT /api/v1/entity/user/:id/settings/:settings_id
// Creates or updates a user-specific setting value.
func (h *HttpHandler) UpsertUserSetting(c *gin.Context) {
	ctx, ok := authContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	settingsID, err := strconv.ParseInt(c.Param("settings_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid settings id"})
		return
	}

	var req UpdateUserSettingRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Domain.UpsertUserSetting(ctx, userID, settingsID, req.Value)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"message": "setting updated",
	})
}
