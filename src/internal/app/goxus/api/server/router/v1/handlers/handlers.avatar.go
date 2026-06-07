package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	userdomain "goxus/src/internal/app/goxus/domain/user"
)

// ---- Avatar HTTP handlers ----

// GetUserAvatar handles GET /api/v1/entity/user/:id/avatar
func (h *HttpHandler) GetUserAvatar(c *gin.Context) {
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

	data, contentType, err := h.Domain.GetAvatar(ctx, id)
	if err != nil {
		if errors.Is(err, userdomain.ErrAvatarNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "avatar not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, contentType, data)
}

// UploadUserAvatar handles POST /api/v1/entity/user/:id/avatar
func (h *HttpHandler) UploadUserAvatar(c *gin.Context) {
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

	// Read uploaded file
	file, _, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar file is required"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read uploaded file"})
		return
	}

	err = h.Domain.UploadAvatar(ctx, id, data)
	if err != nil {
		if errors.Is(err, userdomain.ErrInvalidImageSize) || errors.Is(err, userdomain.ErrInvalidImageFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"message": "avatar uploaded",
	})
}

// DeleteUserAvatar handles DELETE /api/v1/entity/user/:id/avatar
func (h *HttpHandler) DeleteUserAvatar(c *gin.Context) {
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

	err = h.Domain.DeleteAvatar(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"message": "avatar deleted",
	})
}
