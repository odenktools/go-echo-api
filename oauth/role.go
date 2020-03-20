package oauth

import (
	"errors"
	"go-echo-api/entity"
)

var (
	// ErrRoleNotFound ...
	ErrRoleNotFound = errors.New("Role not found")
)

// FindRoleByID looks up a role by ID and returns it
func (s *Service) FindRoleByID(id string) (*entity.OauthRole, error) {
	role := new(entity.OauthRole)
	if s.db.Where("id = ?", id).First(role).RecordNotFound() {
		return nil, ErrRoleNotFound
	}
	return role, nil
}
