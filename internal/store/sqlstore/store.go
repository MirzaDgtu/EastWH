package sqlstore

import (
	"eastwh/internal/store"

	"gorm.io/gorm"
)

type Store struct {
	db                     *gorm.DB
	userRepository         *UserRepository
	userTeamRepository     *UserTeamRepository
	userRoleRepository     *UserRoleRepository
	userProjectRepository  *UserProjectRepository
	teamRepository         *TeamRepository
	orderRepository        *OrderRepository
	projectRepository      *ProjectRepository
	employeeRepository     *EmployeeRepository
	roleRepository         *RoleRepository
	employeeTeamRepository *EmployeeTeamRepository
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

func (s *Store) Team() store.TeamRepository {
	if s.teamRepository != nil {
		return s.teamRepository
	}

	s.teamRepository = &TeamRepository{
		store: s,
	}

	return s.teamRepository
}

func (s *Store) Order() store.OrderRepository {
	if s.orderRepository != nil {
		return s.orderRepository
	}

	s.orderRepository = &OrderRepository{
		store: s,
	}

	return s.orderRepository
}

func (s *Store) Project() store.ProjectRepository {
	if s.projectRepository != nil {
		return s.projectRepository
	}

	s.projectRepository = &ProjectRepository{
		store: s,
	}

	return s.projectRepository
}

func (s *Store) Role() store.RoleRepository {
	if s.roleRepository != nil {
		return s.roleRepository
	}

	s.roleRepository = &RoleRepository{
		store: s,
	}

	return s.roleRepository
}

func (s *Store) Employee() store.EmployeeRepository {
	if s.employeeRepository != nil {
		return s.employeeRepository
	}

	s.employeeRepository = &EmployeeRepository{
		store: s,
	}

	return s.employeeRepository
}

func (s *Store) UserTeam() store.UserTeamRepository {
	if s.userTeamRepository != nil {
		return s.userTeamRepository
	}

	s.userTeamRepository = &UserTeamRepository{
		store: s,
	}

	return s.userTeamRepository
}

func (s *Store) UserRole() store.UserRoleRepository {
	if s.userRoleRepository != nil {
		return s.userRoleRepository
	}

	s.userRoleRepository = &UserRoleRepository{
		store: s,
	}

	return s.userRoleRepository
}

func (s *Store) UserProject() store.UserProjectRepository {
	if s.userProjectRepository != nil {
		return s.userProjectRepository
	}

	s.userProjectRepository = &UserProjectRepository{
		store: s,
	}

	return s.userProjectRepository
}

func (s *Store) EmployeeTeam() store.EmployeeTeamRepository {
	if s.employeeTeamRepository != nil {
		return s.employeeTeamRepository
	}

	s.employeeTeamRepository = &EmployeeTeamRepository{
		store: s,
	}

	return s.employeeTeamRepository
}
