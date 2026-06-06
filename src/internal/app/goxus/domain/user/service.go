// Package userdomain provides the domain service for user management.
// It orchestrates user CRUD operations and integrates with the RBAC service
// for role management.
package userdomain

import (
	"context"

	"goxus/src/internal/pkg/db/goxus"
)

// contextKey is a private type to avoid context key collisions.
type contextKey string

const (
	// CtxKeyActorID is the context key for the authenticated user's ID.
	CtxKeyActorID contextKey = "actor_id"
)

// WithActorID returns a context with the actor user ID set.
func WithActorID(ctx context.Context, actorID int64) context.Context {
	return context.WithValue(ctx, CtxKeyActorID, actorID)
}

// ActorIDFromContext extracts the actor user ID from the context.
// Returns false if not set.
func ActorIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(CtxKeyActorID).(int64)
	return id, ok
}

// Service defines the user domain operations.
// The context.Context carries the authenticated user's identity
// for authorization in the decorator layer.
type Service interface {
	// Create creates a new user with the given name, email, and password.
	// Requires user_add permission.
	Create(ctx context.Context, name, email, password string) (*goxus.User, error)

	// List returns users with pagination and total count.
	// limit=0 returns the default page size (50).
	// offset=0 starts from the beginning.
	// Returns []*goxus.UserWithRole which includes aggregated role names in UserWithRole.Roles.
	List(ctx context.Context, limit, offset int) ([]*goxus.UserWithRole, int64, error)

	// GetByID returns a single user by ID.
	// Requires user_view permission (except when actorID == id).
	GetByID(ctx context.Context, id int64) (*goxus.User, error)

	// Update updates the name and email of a user.
	// Requires user_edit permission.
	Update(ctx context.Context, id int64, name, email string) (*goxus.User, error)

	// UpdatePassword updates the password of a user.
	// Requires user_edit permission.
	UpdatePassword(ctx context.Context, id int64, password string) error

	// Delete soft-deletes a user by ID (sets deleted_at).
	// Requires user_delete permission.
	Delete(ctx context.Context, id int64) error

	// Restore restores a soft-deleted user by clearing deleted_at.
	// Requires user_delete permission (same as delete — inverse operation).
	Restore(ctx context.Context, id int64) (*goxus.User, error)

	// Login authenticates a user by email and password, and returns a token.
	Login(ctx context.Context, email, password string) (*goxus.User, *goxus.UsersToken, error)

	// Logout invalidates a token by setting its deleted_at.
	Logout(ctx context.Context, token string) error

	// ValidateToken validates a Bearer token: looks up the token, checks it
	// is not soft-deleted, looks up the associated user, checks the user is
	// not soft-deleted, updates last_used_at, and returns both.
	ValidateToken(ctx context.Context, token string) (*goxus.User, *goxus.UsersToken, error)

	// GetRoles returns all roles assigned to a user.
	// Requires user_role_view permission.
	GetRoles(ctx context.Context, userID int64) ([]*goxus.RbacRole, error)

	// AssignRole assigns a role to a user.
	// Requires user_role_add permission.
	AssignRole(ctx context.Context, userID int64, roleSlug string) error

	// RevokeRole removes a role from a user.
	// Requires user_role_delete permission.
	RevokeRole(ctx context.Context, userID int64, roleSlug string) error

	// DeleteExpiredTokens soft-deletes all tokens older than ttlDays days.
	// Internal system operation — no RBAC check.
	DeleteExpiredTokens(ctx context.Context, ttlDays int) error
}
