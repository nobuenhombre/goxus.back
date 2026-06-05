package domainapp

import (
	"context"

	"goxus/src/internal/app/goxus/cli"
	configapp "goxus/src/internal/app/goxus/config"
	userdomain "goxus/src/internal/app/goxus/domain/user"
	"goxus/src/internal/pkg/db/goxus"
	"goxus/src/internal/pkg/services/rbac"
)

// DomainService is the business-logic orchestrator.
// All domain operations are accessed through this interface;
// internal services (user, rbac, etc.) are encapsulated.
type DomainService interface {
	Run() error
	GetConfig() *configapp.Config

	// ---- User domain ----

	CreateUser(ctx context.Context, name, email, password string) (*goxus.User, error)
	// ListUsers returns users with pagination and total count.
	// limit=0 returns the default page size (50).
	// offset=0 starts from the beginning.
	ListUsers(ctx context.Context, limit, offset int) ([]*goxus.User, int64, error)
	GetUser(ctx context.Context, id int64) (*goxus.User, error)
	UpdateUser(ctx context.Context, id int64, name, email string) (*goxus.User, error)
	// UpdateUserPassword updates the password of a user.
	UpdateUserPassword(ctx context.Context, id int64, password string) error
	DeleteUser(ctx context.Context, id int64) error

	// RestoreUser restores a soft-deleted user.
	RestoreUser(ctx context.Context, id int64) (*goxus.User, error)

	Login(ctx context.Context, email, password string) (*goxus.User, *goxus.UsersToken, error)
	Logout(ctx context.Context, token string) error

	// ValidateToken validates a Bearer token and returns the associated user and token.
	ValidateToken(ctx context.Context, token string) (*goxus.User, *goxus.UsersToken, error)

	GetUserRoles(ctx context.Context, userID int64) ([]*goxus.RbacRole, error)
	AssignUserRole(ctx context.Context, userID int64, roleSlug string) error
	RevokeUserRole(ctx context.Context, userID int64, roleSlug string) error

	// DeleteExpiredTokens soft-deletes all tokens older than ttlDays days.
	// Internal system operation — used by cronjob, no RBAC.
	DeleteExpiredTokens(ctx context.Context, ttlDays int) error
}

// AppDomain is the concrete implementation of DomainService.
// It holds all internal services and orchestrates business flows.
type AppDomain struct {
	Cli    *cli.Config
	Config *configapp.Config
	Rbac   rbac.Service
	User   userdomain.Service
}

func New(cliConfig cli.Service, appConfig configapp.Service, rbacService rbac.Service, userService userdomain.Service) DomainService {
	return &AppDomain{
		Cli:    cliConfig.(*cli.Config),
		Config: appConfig.Get(),
		Rbac:   rbacService,
		User:   userService,
	}
}

func (d *AppDomain) Run() error {
	// Add your domain orchestration logic here
	return nil
}

func (d *AppDomain) GetConfig() *configapp.Config {
	return d.Config
}

// ---- User domain delegation ----

func (d *AppDomain) CreateUser(ctx context.Context, name, email, password string) (*goxus.User, error) {
	return d.User.Create(ctx, name, email, password)
}

func (d *AppDomain) ListUsers(ctx context.Context, limit, offset int) ([]*goxus.User, int64, error) {
	return d.User.List(ctx, limit, offset)
}

func (d *AppDomain) GetUser(ctx context.Context, id int64) (*goxus.User, error) {
	return d.User.GetByID(ctx, id)
}

func (d *AppDomain) UpdateUser(ctx context.Context, id int64, name, email string) (*goxus.User, error) {
	return d.User.Update(ctx, id, name, email)
}

func (d *AppDomain) UpdateUserPassword(ctx context.Context, id int64, password string) error {
	return d.User.UpdatePassword(ctx, id, password)
}

func (d *AppDomain) DeleteUser(ctx context.Context, id int64) error {
	return d.User.Delete(ctx, id)
}

func (d *AppDomain) RestoreUser(ctx context.Context, id int64) (*goxus.User, error) {
	return d.User.Restore(ctx, id)
}

func (d *AppDomain) Login(ctx context.Context, email, password string) (*goxus.User, *goxus.UsersToken, error) {
	return d.User.Login(ctx, email, password)
}

func (d *AppDomain) Logout(ctx context.Context, token string) error {
	return d.User.Logout(ctx, token)
}

func (d *AppDomain) ValidateToken(ctx context.Context, token string) (*goxus.User, *goxus.UsersToken, error) {
	return d.User.ValidateToken(ctx, token)
}

func (d *AppDomain) GetUserRoles(ctx context.Context, userID int64) ([]*goxus.RbacRole, error) {
	return d.User.GetRoles(ctx, userID)
}

func (d *AppDomain) AssignUserRole(ctx context.Context, userID int64, roleSlug string) error {
	return d.User.AssignRole(ctx, userID, roleSlug)
}

func (d *AppDomain) RevokeUserRole(ctx context.Context, userID int64, roleSlug string) error {
	return d.User.RevokeRole(ctx, userID, roleSlug)
}

func (d *AppDomain) DeleteExpiredTokens(ctx context.Context, ttlDays int) error {
	return d.User.DeleteExpiredTokens(ctx, ttlDays)
}
