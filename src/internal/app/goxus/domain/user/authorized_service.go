package userdomain

import (
	"context"

	"github.com/nobuenhombre/suikat/pkg/ge"

	"goxus/src/internal/pkg/db/goxus"
	"goxus/src/internal/pkg/services/rbac"
	"goxus/src/internal/pkg/services/rbac/permission"
)

// authorizedService is a decorator around Service that adds RBAC permission checks.
// It implements the same Service interface and delegates business logic
// to the inner service after verifying authorization.
type authorizedService struct {
	inner   Service
	rbacSvc rbac.Service
}

// NewAuthorized creates a new authorized decorator around a user domain service.
// The inner service performs pure business logic; this layer adds permission checks.
func NewAuthorized(inner Service, rbacSvc rbac.Service) Service {
	return &authorizedService{
		inner:   inner,
		rbacSvc: rbacSvc,
	}
}

// Create creates a new user. Requires user_add permission.
func (a *authorizedService) Create(ctx context.Context, name, email, password string) (*goxus.User, error) {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return nil, ge.Pin(ErrAccessDenied)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserAdd)
	if err != nil {
		return nil, ge.Pin(err)
	}
	if !ok {
		return nil, ge.Pin(ErrAccessDenied)
	}

	return a.inner.Create(ctx, name, email, password)
}

// List returns users with pagination. Requires user_view permission.
func (a *authorizedService) List(ctx context.Context, limit, offset int) ([]*goxus.UserWithRole, int64, error) {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return nil, 0, ge.Pin(ErrAccessDenied)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserView)
	if err != nil {
		return nil, 0, ge.Pin(err)
	}
	if !ok {
		return nil, 0, ge.Pin(ErrAccessDenied)
	}

	return a.inner.List(ctx, limit, offset)
}

// GetByID returns a single user by ID.
// Requires user_view permission (except self-read: actorID == id).
func (a *authorizedService) GetByID(ctx context.Context, id int64) (*goxus.User, error) {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return nil, ge.Pin(ErrAccessDenied)
	}

	if actorID != id {
		ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserView)
		if err != nil {
			return nil, ge.Pin(err)
		}
		if !ok {
			return nil, ge.Pin(ErrAccessDenied)
		}
	}

	return a.inner.GetByID(ctx, id)
}

// Update updates user name and email. Requires user_edit permission.
func (a *authorizedService) Update(ctx context.Context, id int64, name, email string) (*goxus.User, error) {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return nil, ge.Pin(ErrAccessDenied)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserEdit)
	if err != nil {
		return nil, ge.Pin(err)
	}
	if !ok {
		return nil, ge.Pin(ErrAccessDenied)
	}

	return a.inner.Update(ctx, id, name, email)
}

// UpdatePassword updates the password of a user. Requires user_edit permission.
func (a *authorizedService) UpdatePassword(ctx context.Context, id int64, password string) error {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserEdit)
	if err != nil {
		return ge.Pin(err)
	}
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	return a.inner.UpdatePassword(ctx, id, password)
}

// Delete deletes a user by ID. Requires user_delete permission.
// A user cannot delete themselves.
func (a *authorizedService) Delete(ctx context.Context, id int64) error {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	// User cannot delete themselves
	if actorID == id {
		return ge.Pin(ErrCannotDeleteSelf)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserDelete)
	if err != nil {
		return ge.Pin(err)
	}
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	return a.inner.Delete(ctx, id)
}

// Restore restores a soft-deleted user. Requires user_delete permission (inverse of delete).
func (a *authorizedService) Restore(ctx context.Context, id int64) (*goxus.User, error) {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return nil, ge.Pin(ErrAccessDenied)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserDelete)
	if err != nil {
		return nil, ge.Pin(err)
	}
	if !ok {
		return nil, ge.Pin(ErrAccessDenied)
	}

	return a.inner.Restore(ctx, id)
}

// Login authenticates a user — no permission check (public endpoint).
func (a *authorizedService) Login(ctx context.Context, email, password string) (*goxus.User, *goxus.UsersToken, error) {
	return a.inner.Login(ctx, email, password)
}

// Logout invalidates a token — no permission check (token identifies itself).
func (a *authorizedService) Logout(ctx context.Context, token string) error {
	return a.inner.Logout(ctx, token)
}

// ValidateToken validates a Bearer token — no permission check (public entry point).
func (a *authorizedService) ValidateToken(ctx context.Context, token string) (*goxus.User, *goxus.UsersToken, error) {
	return a.inner.ValidateToken(ctx, token)
}

// GetRoles returns all roles assigned to a user. Requires user_role_view permission.
func (a *authorizedService) GetRoles(ctx context.Context, userID int64) ([]*goxus.RbacRole, error) {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return nil, ge.Pin(ErrAccessDenied)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserRoleView)
	if err != nil {
		return nil, ge.Pin(err)
	}
	if !ok {
		return nil, ge.Pin(ErrAccessDenied)
	}

	return a.inner.GetRoles(ctx, userID)
}

// AssignRole assigns a role to a user. Requires user_role_add permission.
func (a *authorizedService) AssignRole(ctx context.Context, userID int64, roleSlug string) error {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserRoleAdd)
	if err != nil {
		return ge.Pin(err)
	}
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	return a.inner.AssignRole(ctx, userID, roleSlug)
}

// RevokeRole removes a role from a user. Requires user_role_delete permission.
func (a *authorizedService) RevokeRole(ctx context.Context, userID int64, roleSlug string) error {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserRoleDelete)
	if err != nil {
		return ge.Pin(err)
	}
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	return a.inner.RevokeRole(ctx, userID, roleSlug)
}

// DeleteExpiredTokens soft-deletes tokens older than ttlDays — system operation, no RBAC.
func (a *authorizedService) DeleteExpiredTokens(ctx context.Context, ttlDays int) error {
	return a.inner.DeleteExpiredTokens(ctx, ttlDays)
}

// GetAvatar returns a user's avatar image.
// Requires user_view permission (except self-read: actorID == userID).
func (a *authorizedService) GetAvatar(ctx context.Context, userID int64) ([]byte, string, error) {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return nil, "", ge.Pin(ErrAccessDenied)
	}

	if actorID != userID {
		ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserView)
		if err != nil {
			return nil, "", ge.Pin(err)
		}
		if !ok {
			return nil, "", ge.Pin(ErrAccessDenied)
		}
	}

	return a.inner.GetAvatar(ctx, userID)
}

// UploadAvatar saves an avatar image for a user. Requires user_edit permission.
func (a *authorizedService) UploadAvatar(ctx context.Context, userID int64, data []byte) error {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserEdit)
	if err != nil {
		return ge.Pin(err)
	}
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	return a.inner.UploadAvatar(ctx, userID, data)
}

// DeleteAvatar removes a user's custom avatar. Requires user_edit permission.
func (a *authorizedService) DeleteAvatar(ctx context.Context, userID int64) error {
	actorID, ok := ActorIDFromContext(ctx)
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	ok, err := a.rbacSvc.CheckUserPermission(actorID, permission.UserEdit)
	if err != nil {
		return ge.Pin(err)
	}
	if !ok {
		return ge.Pin(ErrAccessDenied)
	}

	return a.inner.DeleteAvatar(ctx, userID)
}
