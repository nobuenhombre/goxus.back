package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	v1middlewares "goxus/src/internal/app/goxus/api/server/router/v1/middlewares"
	userdomain "goxus/src/internal/app/goxus/domain/user"
	"goxus/src/internal/pkg/db/goxus"
)

// ---- DTOs ----

// UserResponse is the public representation of a user (no password).
type UserResponse struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CreateUserRequest is the request body for creating a user.
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=1"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateUserRequest is the request body for updating a user.
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required,min=1"`
	Email string `json:"email" binding:"required,email"`
}

// AssignRoleRequest is the request body for assigning a role.
type AssignRoleRequest struct {
	RoleSlug string `json:"role_slug" binding:"required"`
}

// ---- helpers ----

func userToResponse(u *goxus.User) UserResponse {
	resp := UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.EmailVerifiedAt.Valid {
		resp.EmailVerifiedAt = &u.EmailVerifiedAt.Time
	}
	return resp
}

// authContext builds a context with the authenticated actor from gin.Context.
func authContext(c *gin.Context) (context.Context, bool) {
	actorVal, exists := c.Get(v1middlewares.User)
	if !exists {
		return nil, false
	}
	actor, ok := actorVal.(*goxus.User)
	if !ok {
		return nil, false
	}
	return userdomain.WithActorID(c.Request.Context(), actor.ID), true
}

// ---- User HTTP handlers ----

// CreateUser handles POST /api/v1/entity/user/
func (h *HttpHandler) CreateUser(c *gin.Context) {
	ctx, ok := authContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Domain.CreateUser(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"version": "v1",
		"data":    userToResponse(user),
	})
}

// ListUsers handles GET /api/v1/entity/user/
func (h *HttpHandler) ListUsers(c *gin.Context) {
	ctx, ok := authContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	users, err := h.Domain.ListUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, userToResponse(u))
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"data":    responses,
	})
}

// GetUser handles GET /api/v1/entity/user/:id
func (h *HttpHandler) GetUser(c *gin.Context) {
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

	user, err := h.Domain.GetUser(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"data":    userToResponse(user),
	})
}

// UpdateUser handles PUT /api/v1/entity/user/:id
func (h *HttpHandler) UpdateUser(c *gin.Context) {
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

	var req UpdateUserRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Domain.UpdateUser(ctx, id, req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"data":    userToResponse(user),
	})
}

// DeleteUser handles DELETE /api/v1/entity/user/:id
func (h *HttpHandler) DeleteUser(c *gin.Context) {
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

	err = h.Domain.DeleteUser(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"message": "user deleted",
	})
}

// RestoreUser handles POST /api/v1/entity/user/:id/restore
func (h *HttpHandler) RestoreUser(c *gin.Context) {
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

	user, err := h.Domain.RestoreUser(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"data":    userToResponse(user),
		"message": "user restored",
	})
}

// GetUserRoles handles GET /api/v1/entity/user/:id/roles
func (h *HttpHandler) GetUserRoles(c *gin.Context) {
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

	roles, err := h.Domain.GetUserRoles(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"data":    roles,
	})
}

// AssignUserRole handles POST /api/v1/entity/user/:id/roles
func (h *HttpHandler) AssignUserRole(c *gin.Context) {
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

	var req AssignRoleRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Domain.AssignUserRole(ctx, id, req.RoleSlug)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"message": "role assigned",
	})
}

// RevokeUserRole handles DELETE /api/v1/entity/user/:id/roles/:slug
func (h *HttpHandler) RevokeUserRole(c *gin.Context) {
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

	roleSlug := c.Param("slug")
	if roleSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role slug is required"})
		return
	}

	err = h.Domain.RevokeUserRole(ctx, id, roleSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"message": "role revoked",
	})
}
