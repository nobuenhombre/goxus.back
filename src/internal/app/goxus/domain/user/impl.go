package userdomain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/nobuenhombre/suikat/pkg/ge"

	"goxus/src/internal/pkg/db/goxus"
	"goxus/src/internal/pkg/services/rbac"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyTaken  = errors.New("email already taken")
	ErrAccessDenied       = errors.New("access denied")
	ErrTokenNotFound      = errors.New("token not found")
	ErrCannotDeleteSelf   = errors.New("cannot delete yourself")
	ErrUserAlreadyDeleted = errors.New("user already deleted")
	ErrUserNotDeleted     = errors.New("user not deleted")
)

// impl is the concrete implementation of Service with pure business logic.
// Authorization is handled by the authorizedService decorator.
type impl struct {
	repo    *goxus.DbGoxusRepo
	rbacSvc rbac.Service
}

// New creates a new user domain service with pure business logic.
func New(dbRepo *goxus.DbGoxusRepo, rbacSvc rbac.Service) Service {
	return &impl{
		repo:    dbRepo,
		rbacSvc: rbacSvc,
	}
}

// Create creates a new user.
func (s *impl) Create(_ context.Context, name, email, password string) (*goxus.User, error) {
	// check if email already taken (including soft-deleted — prevent restore conflicts)
	existing, err := s.repo.User.GetUserByEmail(email)
	if err == nil && existing != nil {
		return nil, ge.Pin(fmt.Errorf("email '%s': %w", email, ErrEmailAlreadyTaken))
	}

	now := time.Now()
	user := &goxus.User{
		Name:      name,
		Email:     email,
		Password:  password, // TODO: hash password before storing
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.repo.User.Save(user)
	if err != nil {
		return nil, ge.Pin(err)
	}

	return user, nil
}

// List returns all users.
func (s *impl) List(_ context.Context) ([]*goxus.User, error) {
	users, err := s.repo.User.GetAll()
	if err != nil {
		return nil, ge.Pin(err)
	}
	return users, nil
}

// GetByID returns a single user by ID.
func (s *impl) GetByID(_ context.Context, id int64) (*goxus.User, error) {
	user, err := s.repo.User.GetUserByID(id)
	if err != nil {
		return nil, ge.Pin(fmt.Errorf("user id '%d': %w", id, ErrUserNotFound))
	}
	return user, nil
}

// Update updates user name and email.
func (s *impl) Update(_ context.Context, id int64, name, email string) (*goxus.User, error) {
	user, err := s.repo.User.GetUserByID(id)
	if err != nil {
		return nil, ge.Pin(fmt.Errorf("user id '%d': %w", id, ErrUserNotFound))
	}

	// if email changed, check it's not taken (including soft-deleted — prevent restore conflicts)
	if email != user.Email {
		existing, err := s.repo.User.GetUserByEmail(email)
		if err == nil && existing != nil {
			return nil, ge.Pin(fmt.Errorf("email '%s': %w", email, ErrEmailAlreadyTaken))
		}
	}

	user.Name = name
	user.Email = email
	user.UpdatedAt = time.Now()

	err = s.repo.User.Save(user)
	if err != nil {
		return nil, ge.Pin(err)
	}

	return user, nil
}

// Delete soft-deletes a user by setting deleted_at.
func (s *impl) Delete(_ context.Context, id int64) error {
	user, err := s.repo.User.GetUserByID(id)
	if err != nil {
		return ge.Pin(fmt.Errorf("user id '%d': %w", id, ErrUserNotFound))
	}

	if user.DeletedAt.Valid {
		return ge.Pin(fmt.Errorf("user id '%d': %w", id, ErrUserAlreadyDeleted))
	}

	now := time.Now()
	user.DeletedAt = pq.NullTime{Time: now, Valid: true}
	user.UpdatedAt = now

	err = s.repo.User.Save(user)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

// Restore restores a soft-deleted user by clearing deleted_at.
func (s *impl) Restore(_ context.Context, id int64) (*goxus.User, error) {
	user, err := s.repo.User.GetUserByID(id)
	if err != nil {
		return nil, ge.Pin(fmt.Errorf("user id '%d': %w", id, ErrUserNotFound))
	}

	if !user.DeletedAt.Valid {
		return nil, ge.Pin(fmt.Errorf("user id '%d': %w", id, ErrUserNotDeleted))
	}

	now := time.Now()
	user.DeletedAt = pq.NullTime{Valid: false}
	user.UpdatedAt = now

	err = s.repo.User.Save(user)
	if err != nil {
		return nil, ge.Pin(err)
	}

	return user, nil
}

// GetRoles returns all roles assigned to a user.
func (s *impl) GetRoles(_ context.Context, userID int64) ([]*goxus.RbacRole, error) {
	// verify user exists
	_, err := s.repo.User.GetUserByID(userID)
	if err != nil {
		return nil, ge.Pin(fmt.Errorf("user id '%d': %w", userID, ErrUserNotFound))
	}

	roles, err := s.rbacSvc.GetUserRoles(userID)
	if err != nil {
		return nil, ge.Pin(err)
	}

	return roles, nil
}

// AssignRole assigns a role to a user.
func (s *impl) AssignRole(_ context.Context, userID int64, roleSlug string) error {
	// verify user exists
	_, err := s.repo.User.GetUserByID(userID)
	if err != nil {
		return ge.Pin(fmt.Errorf("user id '%d': %w", userID, ErrUserNotFound))
	}

	return s.rbacSvc.AssignRoleToUser(userID, roleSlug)
}

// RevokeRole removes a role from a user.
func (s *impl) RevokeRole(_ context.Context, userID int64, roleSlug string) error {
	// verify user exists
	_, err := s.repo.User.GetUserByID(userID)
	if err != nil {
		return ge.Pin(fmt.Errorf("user id '%d': %w", userID, ErrUserNotFound))
	}

	return s.rbacSvc.RevokeUserRole(userID, roleSlug)
}

// DeleteExpiredTokens soft-deletes all tokens older than ttlDays days.
func (s *impl) DeleteExpiredTokens(_ context.Context, ttlDays int) error {
	return s.repo.UsersToken.DeleteExpiredTokens(ttlDays)
}

// generateToken creates a UUID v4 token.
func generateToken() string {
	return uuid.New().String()
}

// Login authenticates a user by email and password, and returns a token.
func (s *impl) Login(_ context.Context, email, password string) (*goxus.User, *goxus.UsersToken, error) {
	// find active user by email (not soft-deleted)
	user, err := s.repo.User.GetActiveUserByEmail(email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ge.Pin(ErrUserNotFound)
		}
		return nil, nil, ge.Pin(err)
	}

	// verify password
	if user.Password != password {
		return nil, nil, ge.Pin(ErrAccessDenied)
	}

	// create token
	now := time.Now()
	token := &goxus.UsersToken{
		Token:     generateToken(),
		UserID:    user.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.repo.UsersToken.Save(token)
	if err != nil {
		return nil, nil, ge.Pin(err)
	}

	return user, token, nil
}

// ValidateToken validates a Bearer token: looks up the token, checks it
// is not soft-deleted, looks up the associated user, checks the user is
// not soft-deleted, updates last_used_at, and returns both.
func (s *impl) ValidateToken(_ context.Context, token string) (*goxus.User, *goxus.UsersToken, error) {
	ut, err := s.repo.UsersToken.GetUsersTokenByToken(token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ge.Pin(ErrTokenNotFound)
		}
		return nil, nil, ge.Pin(err)
	}

	if ut.DeletedAt.Valid {
		return nil, nil, ge.Pin(ErrTokenNotFound)
	}

	user, err := s.repo.User.GetUserByID(ut.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ge.Pin(ErrUserNotFound)
		}
		return nil, nil, ge.Pin(err)
	}

	if user.DeletedAt.Valid {
		return nil, nil, ge.Pin(ErrUserNotFound)
	}

	// Update last_used_at
	ut.LastUsedAt = pq.NullTime{Time: time.Now(), Valid: true}
	err = s.repo.UsersToken.Save(ut)
	if err != nil {
		return nil, nil, ge.Pin(err)
	}

	return user, ut, nil
}

// Logout invalidates a token by setting its deleted_at.
func (s *impl) Logout(_ context.Context, token string) error {
	ut, err := s.repo.UsersToken.GetUsersTokenByToken(token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ge.Pin(ErrTokenNotFound)
		}
		return ge.Pin(err)
	}

	// soft-delete the token
	now := time.Now()
	ut.DeletedAt = pq.NullTime{Time: now, Valid: true}
	ut.UpdatedAt = now

	err = s.repo.UsersToken.Save(ut)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}
