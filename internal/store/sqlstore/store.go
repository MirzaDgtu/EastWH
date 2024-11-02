package sqlstore

import (
	"eastwh/internal/store"

	"gorm.io/gorm"
)

type Store struct {
	db                    *gorm.DB
	userRepository        *UserRepository
	userTeamRepository    *UserTeamRepository
	userRoleRepository    *UserRoleRepository
	userProjectRepository *UserProjectRepository
	teamRepository        *TeamRepository
	orderRepository       *OrderRepository
	projectRepository     *ProjectRepository
	employeeRepository    *EmployeeRepository
	roleRepository        *RoleRepository
}

func New(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

