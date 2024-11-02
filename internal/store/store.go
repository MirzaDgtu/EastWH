package store

type Store interface {
	Employee() EmployeeRepository
	Order() OrderRepository
	Project() ProjectRepository
	Role() RoleRepository
	Team() TeamRepository
	User() UserRepository
	UserRole() UserRoleRepository
	UserTeam() UserTeamRepository
	UserProject() UserProjectRepository
}
