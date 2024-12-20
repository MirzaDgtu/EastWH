package apiserver

import (
	"context"
	"eastwh/internal/model"
	"eastwh/internal/store"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not autenticated")
)

var (
	hmacSampleSecret = "8a046a6b436496d9c7af3e196a73ee9948677eb30b251706667ad59d6261bd78d2f6f501a6dea0118cfb3b0dcd62d6c9eb88142e2c24c2c686133a935cd65651"
)

type server struct {
	router *gin.Engine
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: gin.Default(),
		store:  store,
	}

	s.router.Use()

	s.router.Use(func(ctx *gin.Context) {
		fmt.Println("Запрос с сайта: ", ctx.Request.Header.Get("Origin"))
		ctx.Next()
	})
	confCors := cors.DefaultConfig()
	confCors.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	confCors.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	confCors.ExposeHeaders = []string{"Content-Length"}
	confCors.AllowCredentials = true
	confCors.MaxAge = 12 * time.Hour
	confCors.AllowOriginFunc = func(origin string) bool {
		return true
	}

	s.router.Use(cors.New(confCors))

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	apiGroup := s.router.Group("/api/v1")
	{

		apiGroup.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		userGroup := apiGroup.Group("/user") //, s.AuthMW
		{
			userGroup.POST("/logout/", s.Logout)
			userGroup.PUT("/", s.UpdateUser)
			userGroup.POST("/update/password/", s.UpdatePassword)
			userGroup.POST("/block/", s.BlockedUser)
			userGroup.GET("/profile/", s.GetUserProfile)
			userGroup.GET("/employees", s.GetEmployeeByUserID)

			userProjectsGroup := userGroup.Group("/projects")
			{
				userProjectsGroup.POST("", s.AddUserProjects)
				userProjectsGroup.GET("", s.GetUserProjects)
				userProjectsGroup.GET("/user/", s.GetUserProjectsByUserId)
				userProjectsGroup.GET("/project/", s.GetUserProjectsByProjectId)
			}

			userProjectGroup := userGroup.Group("/project")
			{
				userProjectGroup.GET("/", s.GetUserProjectById)
				userProjectGroup.PUT("/", s.UpdateUserProject)
				userProjectGroup.DELETE("/", s.DeleteUserProject)
				userProjectGroup.DELETE("/user/", s.DeleteProjectByUserID)
			}

			userRolesGroup := userGroup.Group("/roles")
			{
				userRolesGroup.POST("", s.AddUserRoles)
				userRolesGroup.GET("", s.GetUserRoles)
				userRolesGroup.GET("/user/", s.GetUserRolesByUserId)
				userRolesGroup.GET("/role/", s.GetUserRolesByRoleId)
			}

			userRoleGroup := userGroup.Group("/role")
			{
				userRoleGroup.GET("/", s.GetUserRoleById)
				userRoleGroup.PUT("/", s.UpdateUserRole)
				userRoleGroup.DELETE("/", s.DeleteUserRole)
			}

			userTeamsGroup := userGroup.Group("/teams")
			{
				userTeamsGroup.POST("", s.AddUserTeams)
				userTeamsGroup.GET("", s.GetUserTeams)
				userTeamsGroup.GET("/user/", s.GetUserTeamsByUserId)
				userTeamsGroup.GET("/team/", s.GetUserTeamsByTeamId)
				userTeamsGroup.DELETE("/", s.DeleteUserTeamByID)
			}

			userTeamGroup := userGroup.Group("team")
			{
				userTeamGroup.GET("/", s.GetUserTeamById)
				userTeamGroup.PUT("/", s.UpdateUserTeam)
				userTeamGroup.DELETE("/", s.DeleteUserTeam)
			}
		}

		usersGroup := apiGroup.Group("/users")
		{
			usersGroup.POST("", s.AddUser)
			usersGroup.POST("/login", s.Login)
			usersGroup.GET("", s.GetUsers)
			usersGroup.POST("/password/restore", s.RestoreUserPassword)
		}

		employeesGroup := apiGroup.Group("/employees")
		{
			employeesGroup.POST("", s.AddEmployee)
			employeesGroup.GET("", s.GetEmployees)
		}

		employeeGroup := apiGroup.Group("/employee")
		{
			employeeGroup.GET("/", s.GetEmployeeByID)
			employeeGroup.GET("/code/", s.GetEmployeeByCode)
			employeeGroup.PUT("/", s.UpdateEmployee)
			employeeGroup.DELETE("/", s.DeleteEmployee)

			employeeTeamsGroup := employeeGroup.Group("/teams")
			{
				employeeTeamsGroup.POST("", s.AddEmployeeTeams)
				employeeTeamsGroup.GET("", s.GetEmployeeTeams)
				employeeTeamsGroup.GET("/employee/", s.GetEmployeeTeamsByEmployeeId)
				employeeTeamsGroup.GET("/team/", s.GetEmployeeTeamsByTeamId)
			}
			employeeTeamGroup := employeeGroup.Group("/team")
			{
				employeeTeamGroup.GET("/", s.GetEmployeeTeamById)
				employeeTeamGroup.PUT("/", s.UpdateEmployeeTeam)
				employeeTeamGroup.DELETE("/id/", s.DeleteEmployeeTeamByID)
				employeeTeamGroup.DELETE("/", s.DeleteEmployeeTeam)

			}
		}

		orderGroup := apiGroup.Group("/order")
		{
			orderGroup.GET("/", s.GetOrderByID)
			orderGroup.GET("/uid/", s.GetOrderByUID)
			orderGroup.PUT("/collector/", s.UpdateOrderCollector)
			orderGroup.PUT("/check", s.UpdateOrderCheck)
		}

		ordersGroup := apiGroup.Group("/orders")
		{
			ordersGroup.GET("/user/", s.GetOrdersByUserId)
			ordersGroup.GET("/daterange/", s.GetOrdersByDateRange)
			ordersGroup.POST("/access/", s.GetOrdersByAccessUser)
			ordersGroup.POST("", s.AddOrders)
			ordersGroup.GET("", s.GetOrders)
			ordersGroup.POST("/assembly/", s.GetAssemblyOrders)
			orderGroup.POST("/check", s.GetOrdersChecked)
		}

		teamGroup := apiGroup.Group("/team")
		{
			teamGroup.GET("/", s.GetTeamByID)
			teamGroup.PUT("/", s.UpdateTeam)
			teamGroup.DELETE("/", s.DeleteTeam)
		}

		teamsGroup := apiGroup.Group("/teams")
		{
			teamsGroup.POST("", s.AddTeams)
			teamsGroup.GET("", s.GetTeams)
		}

		projectGroup := apiGroup.Group("/project")
		{
			projectGroup.GET("/", s.GetProjectById)
			projectGroup.DELETE("/", s.DeleteProject)
			projectGroup.PUT("/", s.UpdateProject)
		}

		projectsGroup := apiGroup.Group("/projects")
		{
			projectsGroup.POST("", s.AddProject)
			projectsGroup.GET("", s.GetProjects)

		}

		rolesGroup := apiGroup.Group("/roles")
		{
			rolesGroup.POST("", s.AddRoles)
			rolesGroup.GET("", s.GetRoles)
		}

		roleGroup := apiGroup.Group("/role")
		{
			roleGroup.GET("/", s.GetRoleByID)
			roleGroup.PUT("/", s.UpdateRole)
			roleGroup.DELETE("/", s.DeleteRole)
		}
	}
}

func setCookie(ctx *gin.Context, token string) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Auth", token, 3600*24*100, "", "", false, true)
}

func (s *server) AuthMW(ctx *gin.Context) {
	// Получение токена из куки
	tokenStr, err := ctx.Cookie("Auth")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No auth token"})
		ctx.Abort()
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse JWT"})
		ctx.Abort()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "JWT Claims failed"})
		ctx.Abort()
	}

	if claims["ttl"].(float64) < float64(time.Now().Unix()) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token expired"})
		ctx.Abort()
	}

	user, err := s.store.User().ByID(uint((claims["userID"].(float64))))

	if user.ID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Could not find the user!"})
		ctx.Abort()
	}

	ctx.Set("user", user)

	ctx.Next()
}

func createAndSignJWT(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"ttl":    time.Now().Add(time.Hour * 24 * 100).Unix(),
	})

	return token.SignedString([]byte(hmacSampleSecret))
}

// User ...
func (s *server) AddUser(ctx *gin.Context) {
	var user model.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка создания пользователя",
			"error": err.Error()})
		return
	}

	user, err = s.store.User().Add(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка создания пользователя",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно создан",
		"user": user})
}

func (s *server) UpdateUser(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	var user model.User
	err = ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность обновляемых данных",
			"error": err.Error()})
		return
	}

	user.ID = uint(ID)

	user, err = s.store.User().Update(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления данных пользователя",
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Данные пользователя успешно обновлены",
		"user": user})
}

func (s *server) Login(ctx *gin.Context) {
	type request struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var req request
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность введенных данных",
			"error": err.Error()})
		return
	}

	user, err := s.store.User().Login(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка авторизации",
			"error": errIncorrectEmailOrPassword})
		return
	}

	tokenString, err := createAndSignJWT(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"messge": "Ошибка создания JWT токена",
			"error": err.Error()})
		return
	}

	err = s.store.User().UpdateToken(user.ID, tokenString)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления JWT токена",
			"error": err.Error()})
		return
	}

	user.Token = tokenString
	setCookie(ctx, tokenString)
	ctx.JSON(http.StatusOK, gin.H{"message": "Вы успешно авторизованы",
		"user": user})
}

func (s *server) RestoreUserPassword(ctx *gin.Context) {
	email := ctx.Query("email")

	if utf8.RuneCountInString(email) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "Email не должен быть пустым"})
		return
	}

	password, err := s.store.User().Restore(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка установки временного пароля",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"password": password})
}

func (s *server) GetUsers(ctx *gin.Context) {
	users, err := s.store.User().All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка пользователей",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (s *server) Logout(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	err = s.store.User().Logout(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка выхода из аккаунта",
			"error": err.Error()})
		return
	}

	ctx.SetCookie("Auth", "deleted", 0, "", "", false, false)

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Вы успешно вышли из аккаунта"})
}

func (s *server) UpdatePassword(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	type request struct {
		Password string
	}
	var req request
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность введенного пароля",
			"error": err.Error()})
		return
	}

	if utf8.RuneCountInString(req.Password) <= 6 {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "Длина пароля должна быть не меньше 6"})
		return
	}

	fmt.Println(ID, " ", req.Password)

	err = s.store.User().ChangePassword(uint(ID), req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка изменения пароля пользователя",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Пароль успешно изменен"})
}

func (s *server) BlockedUser(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	type request struct {
		Blocked bool `json:"blocked"`
	}
	var req request

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность значения блокировки",
			"error": err.Error()})
		return
	}

	err = s.store.User().BlockedUser(uint(ID), req.Blocked)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка блокировки пользователя",
			"error": err.Error()})
		return
	}

	var msg string
	if req.Blocked {
		msg = "Пользователь " + pID + " заблокирован"
	} else {
		msg = "Пользователь " + pID + " разблокирован"
	}
	ctx.JSON(http.StatusOK, msg)
}

func (s *server) GetUserProfile(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	user, err := s.store.User().Profile(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения профиля пользователя",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (s *server) GetEmployeeByUserID(ctx *gin.Context) {
	pID := ctx.Query("user_id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}
	employee, err := s.store.User().EmployeeByUserID(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения сотрудников",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, employee)
}

// Employee...
func (s *server) AddEmployee(ctx *gin.Context) {
	var employees []model.Employee

	err := ctx.ShouldBindJSON(&employees)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка чтения данных",
			"error":   err.Error(),
		})
		return
	}

	// Создаем каналы для результатов и ошибок
	results := make(chan model.Employee, len(employees))
	errors := make(chan error, len(employees))

	// WaitGroup для отслеживания завершения всех горутин
	var wg sync.WaitGroup

	// Запускаем горутину для каждого сотрудника
	for _, emp := range employees {
		wg.Add(1)
		go func(emp model.Employee) {
			defer wg.Done()
			employee, err := s.store.Employee().Add(emp)
			if err != nil {
				errors <- err
				return
			}
			results <- employee
		}(emp)
	}

	// Горутина для закрытия каналов после завершения всех операций
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Собираем результаты и ошибки
	var addedEmployees []model.Employee
	var errorMessages []string

	// Читаем из каналов, пока они не закрыты
	for {
		select {
		case employee, ok := <-results:
			if !ok {
				results = nil
				continue
			}
			addedEmployees = append(addedEmployees, employee)

		case err, ok := <-errors:
			if !ok {
				errors = nil
				continue
			}
			errorMessages = append(errorMessages, err.Error())

		// Если оба канала закрыты, завершаем цикл
		default:
			if results == nil && errors == nil {
				goto Done
			}
		}
	}

Done:
	// Формируем ответ
	response := gin.H{
		"added_employees": addedEmployees,
	}

	if len(errorMessages) > 0 {
		response["errors"] = errorMessages
		ctx.JSON(http.StatusMultiStatus, response)
		return
	}

	ctx.JSON(http.StatusCreated, addedEmployees)
}

func (s *server) GetEmployees(ctx *gin.Context) {
	employees, err := s.store.Employee().All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка сотрудников",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, employees)
}

func (s *server) GetEmployeeByID(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	employee, err := s.store.Employee().ByID(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения сотрудника по ID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, employee)
}

func (s *server) GetEmployeeByCode(ctx *gin.Context) {
	pCode := ctx.Query("code")

	employee, err := s.store.Employee().ByCode(pCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения сотрудника по Code",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, employee)
}

func (s *server) UpdateEmployee(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	var employee model.Employee
	err = ctx.ShouldBindJSON(&employee)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность обновляемых данных",
			"error": err.Error()})
		return
	}

	employee.ID = uint(ID)
	employee, err = s.store.Employee().Update(employee)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления данных сотрудника",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Данные сотрудника успешно обновлены",
		"employee": employee})
}

func (s *server) DeleteEmployee(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	err = s.store.Employee().Delete(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления сотрудника по ID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Данные сотрудника успешно удалены"})
}

// Order...
func (s *server) AddOrders(ctx *gin.Context) {
	var orders []model.Order

	err := ctx.ShouldBindJSON(&orders)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка проверки данных заказа",
			"error":   err.Error(),
		})
		return
	}

	// Создаем контекст с таймаутом
	ctxTimeout, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	// Каналы для результатов и ошибок
	results := make(chan model.Order, len(orders))
	errors := make(chan struct {
		order model.Order
		err   error
	}, len(orders))

	// WaitGroup для отслеживания завершения всех горутин
	var wg sync.WaitGroup

	// Запускаем горутину для каждого заказа
	for _, order := range orders {
		wg.Add(1)
		go func(order model.Order) {
			defer wg.Done()

			// Проверяем контекст перед обработкой
			select {
			case <-ctxTimeout.Done():
				errors <- struct {
					order model.Order
					err   error
				}{order, fmt.Errorf("timeout processing order")}
				return
			default:
			}

			createdOrder, err := s.store.Order().Add(order)
			if err != nil {
				errors <- struct {
					order model.Order
					err   error
				}{order, err}
				return
			}
			results <- createdOrder
		}(order)
	}

	// Горутина для закрытия каналов после завершения всех операций
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Собираем результаты
	var addedOrders []model.Order
	var failedOrders []struct {
		Order   model.Order `json:"order"`
		Message string      `json:"error"`
	}

	// Читаем из каналов до их закрытия
	for {
		select {
		case <-ctxTimeout.Done():
			ctx.JSON(http.StatusGatewayTimeout, gin.H{
				"message":       "Превышено время ожидания при создании заказов",
				"added_orders":  addedOrders,
				"failed_orders": failedOrders,
				"error":         "timeout",
			})
			return

		case order, ok := <-results:
			if !ok {
				results = nil
				continue
			}
			addedOrders = append(addedOrders, order)

		case errorData, ok := <-errors:
			if !ok {
				errors = nil
				continue
			}
			failedOrders = append(failedOrders, struct {
				Order   model.Order `json:"order"`
				Message string      `json:"error"`
			}{
				Order:   errorData.order,
				Message: errorData.err.Error(),
			})

		default:
			if results == nil && errors == nil {
				goto Done
			}
		}
	}

Done:
	// Формируем итоговый ответ
	response := gin.H{
		"message":      "Обработка заказов завершена",
		"added_orders": addedOrders,
	}

	if len(failedOrders) > 0 {
		response["failed_orders"] = failedOrders
		ctx.JSON(http.StatusMultiStatus, response)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (s *server) GetOrders(ctx *gin.Context) {
	orders, err := s.store.Order().All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка заказов",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (s *server) GetOrderByID(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	order, err := s.store.Order().ByID(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения заказа по ID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (s *server) GetOrderByUID(ctx *gin.Context) {
	pID := ctx.Query("order_uid")
	OrderUID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность OrderUID",
			"error": err.Error()})
		return
	}

	order, err := s.store.Order().ByOrderUID(uint(OrderUID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения заказа по OrderUID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (s *server) GetOrdersByUserId(ctx *gin.Context) {
	pID := ctx.Query("user_id")
	UserID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность UserUID",
			"error": err.Error()})
		return
	}

	order, err := s.store.Order().ByUserID(uint(UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения заказа по UserUID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (s *server) UpdateOrderCheck(ctx *gin.Context) {
	type request struct {
		OrderUID uint `json:"order_uid"`
		UserID   uint `json:"user_id"`
		Check    bool `json:"check"`
	}
	var reqs []request
	if err := ctx.ShouldBindJSON(&reqs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Проверьте корректность передаваемых данных",
			"error":   err.Error(),
		})
		return
	}

	for _, req := range reqs {
		err := s.store.Order().SetCheck(req.OrderUID, req.UserID, req.Check)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления информации",
				"error": err.Error()})
			continue
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Обновление выполнено успешно"})
}

func (s *server) UpdateOrderCollector(ctx *gin.Context) {
	type request struct {
		OrderUID   uint `json:"order_uid"`
		UserID     uint `json:"user_id"`
		EmployeeID uint `json:"employee_id"`
	}

	var reqs []request
	if err := ctx.ShouldBindJSON(&reqs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Проверьте корректность передаваемых данных",
			"error":   err.Error(),
		})
		return
	}

	errors := make(chan error, len(reqs))
	var wg sync.WaitGroup

	for _, req := range reqs {
		wg.Add(1)
		go func(req request) {
			defer wg.Done()
			if err := s.store.Order().SetCollector(req.OrderUID, req.UserID, req.EmployeeID); err != nil {
				errors <- err
			}
		}(req)
	}

	// Закрываем канал после завершения всех горутин
	go func() {
		wg.Wait()
		close(errors)
	}()

	// Собираем ошибки
	var errorMessages []string
	for err := range errors {
		errorMessages = append(errorMessages, err.Error())
	}

	if len(errorMessages) > 0 {
		ctx.JSON(http.StatusMultiStatus, gin.H{"errors": errorMessages})
		return
	}

	// Возвращаем успешный статус, если ошибок нет
	ctx.JSON(http.StatusOK, gin.H{"message": "Обновление выполнено успешно"})
}

func (s *server) GetOrdersByDateRange(ctx *gin.Context) {

	type request struct {
		DtStart  string `json:"dt_start"`
		DtFinish string `json:"dt_finish"`
	}

	var req request
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	findedOrders, err := s.store.Order().ByDateRange(req.DtStart, req.DtFinish)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, findedOrders)
}

func (s *server) GetAssemblyOrders(ctx *gin.Context) {
	type request struct {
		StartDT  string `json:"start_dt"`
		FinishDT string `json:"finish_dt"`
	}

	var req request
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность передаваемых данных",
			"error": err.Error()})
		return
	}

	assemblyOrders, err := s.store.Order().AssemblyOrder(req.StartDT, req.FinishDT)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка собранных заказов",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Собранные заказы успешно получены",
		"orders": assemblyOrders})
}

func (s *server) GetOrdersByAccessUser(ctx *gin.Context) {
	pUserID := ctx.Query("user_id")
	UserID, err := strconv.Atoi(pUserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность user_id",
			"error": err.Error()})
		return
	}

	type request struct {
		StartDT  string `json:"start_dt"`
		FinishDT string `json:"finish_dt"`
	}

	var req request
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте коррекность передаваемых данных",
			"error": err.Error()})
		return
	}

	orders, err := s.store.Order().ByAccessUser(uint(UserID), req.StartDT, req.FinishDT)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка заказов",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (s *server) GetOrdersChecked(ctx *gin.Context) {
	type request struct {
		StartDT  string `json:"start_dt"`
		FinishDT string `json:"finish_dt"`
		Check    bool   `json:"check"`
	}

	var req request
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность передаваемых данных",
			"error": err.Error()})
		return
	}

	ChekedOrders, err := s.store.Order().CheckedList(req.StartDT, req.FinishDT, req.Check)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка собранных заказов",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Собранные заказы успешно получены",
		"orders": ChekedOrders})
}

// Teams
func (s *server) AddTeams(ctx *gin.Context) {
	var teams []model.Team

	err := ctx.ShouldBindJSON(&teams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка проверки данных команды.",
			"error": err.Error()})
		return
	}

	var addedTeams []model.Team
	for _, team := range teams {
		team, err = s.store.Team().Add(team)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления команды " + team.Name,
				"error": err.Error()})
			continue
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"message": "Команда - " + team.Name + " успешно добавлена"})
			addedTeams = append(addedTeams, team)
		}
	}

	ctx.JSON(http.StatusCreated, addedTeams)
}

func (s *server) GetTeams(ctx *gin.Context) {
	teams, err := s.store.Team().All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка команд",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, teams)
}

func (s *server) GetTeamByID(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	team, err := s.store.Team().ByID(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения заказа по ID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, team)
}

func (s *server) UpdateTeam(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	var team model.Team
	err = ctx.ShouldBindJSON(&team)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность передаваемых данных",
			"error": err.Error()})
		return
	}

	team.ID = uint(ID)

	team, err = s.store.Team().Update(team)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления информации о команде",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Информация о команде успешно обновлена",
		"team": team})
}

func (s *server) DeleteTeam(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	err = s.store.Team().Delete(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления команды",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Команда успешно удалена"})
}

// Project

func (s *server) AddProject(ctx *gin.Context) {
	var project []model.Project

	err := ctx.ShouldBindJSON(&project)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка чтения данных",
			"error": err.Error()})
		return
	}

	var addedProject []model.Project
	for _, req := range project {
		project, err := s.store.Project().Add(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления проекта",
				"error": err.Error()})
			continue
		} else {
			addedProject = append(addedProject, project)
		}
	}
	ctx.JSON(http.StatusCreated, addedProject)
}

func (s *server) GetProjects(ctx *gin.Context) {
	projects, err := s.store.Project().All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка проектов",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}

func (s *server) GetProjectById(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	project, err := s.store.Project().ByID(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения проекта по ID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, project)
}

func (s *server) DeleteProject(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	err = s.store.Project().Delete(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления проекта по ID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Данные проекта успешно удалены"})
}

func (s *server) DeleteProjectByUserID(ctx *gin.Context) {
	pUserID := ctx.Query("user_id")
	UserID, err := strconv.Atoi(pUserID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность user_id",
			"error": err.Error()})
		return
	}

	type request struct {
		ProjectID uint `json:"project_id"`
	}

	var req request
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка получения project_id",
			"error": err.Error()})
		return
	}

	err = s.store.UserProject().DeleteUserProject(uint(UserID), req.ProjectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления проекта пользователя",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Проект успешно удален"})
}

func (s *server) UpdateProject(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	var project model.Project
	err = ctx.ShouldBindJSON(&project)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность обновляемых данных",
			"error": err.Error()})
		return
	}

	project.ID = uint(ID)

	project, err = s.store.Project().Update(project)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления данных проекта",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Данные проекта успешно обновлены",
		"project": project})
}

// UserProjects ...

func (s *server) AddUserProjects(ctx *gin.Context) {
	var userProjects []model.UserProject

	err := ctx.ShouldBindJSON(&userProjects)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка чтения данных",
			"error": err.Error()})
		return
	}

	var addedUP []model.UserProject
	for _, req := range userProjects {
		userProjects, err := s.store.UserProject().Add(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления проекта пользователя",
				"error": err.Error()})
			continue
		} else {
			addedUP = append(addedUP, userProjects)
		}
	}
	ctx.JSON(http.StatusCreated, addedUP)
}

func (s *server) GetUserProjects(ctx *gin.Context) {
	up, err := s.store.UserProject().All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка проектов пользователей",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, up)
}

func (s *server) GetUserProjectsByUserId(ctx *gin.Context) {
	pUserID := ctx.Query("user_id")
	UserID, err := strconv.Atoi(pUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность UserId",
			"error": err.Error()})
		return
	}

	UserProject, err := s.store.UserProject().ByUserID(uint(UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка проектов пользователей по UserId",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, UserProject)
}

func (s *server) GetUserProjectsByProjectId(ctx *gin.Context) {
	pProjectID := ctx.Query("project_id")
	ProjectID, err := strconv.Atoi(pProjectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ProjectId",
			"error": err.Error()})
		return
	}

	UserProject, err := s.store.UserProject().ByProjectID(uint(ProjectID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка проектов пользователей по ProjectId",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, UserProject)
}

func (s *server) GetUserProjectById(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	userProject, err := s.store.UserProject().ByID(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения проекта пользователя по ID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, userProject)
}

func (s *server) UpdateUserProject(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	var userProject model.UserProject
	err = ctx.ShouldBindJSON(&userProject)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность передаваемых данных",
			"error": err.Error()})
		return
	}

	userProject.ID = uint(ID)

	userProject, err = s.store.UserProject().Update(userProject)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления информации о проекте пользователя",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Информация о проекте пользователя успешно обновлена",
		"user_project": userProject})
}

func (s *server) DeleteUserProject(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	err = s.store.UserProject().Delete(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления проекта пользователя",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Проект пользователя успешно удален"})
}

// Roles

func (s *server) AddRoles(ctx *gin.Context) {
	var roles []model.Role

	err := ctx.ShouldBindJSON(&roles)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте передаваемые данные",
			"error": err.Error()})
		return
	}

	var addedRoles []model.Role
	for _, role := range roles {
		role, err = s.store.Role().Add(role)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления роли - " + role.Name,
				"error": err.Error()})
			continue
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"message": "Роль " + role.Name + " успешно создана"})
			addedRoles = append(addedRoles, role)
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Создание ролей успешно завершено",
		"roles": addedRoles})
}

func (s *server) GetRoles(ctx *gin.Context) {
	roles, err := s.store.Role().All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка ролей",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (s *server) GetRoleByID(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	role, err := s.store.Role().ByID(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения роли по ID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, role)
}

func (s *server) UpdateRole(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	var role model.Role
	err = ctx.ShouldBindJSON(&role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность изменяемых данных",
			"error": err.Error()})
		return
	}

	role.ID = uint(ID)
	role, err = s.store.Role().Update(role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления данных роли",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Данные роли успешно обновлены",
		"role": role})
}

func (s *server) DeleteRole(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	err = s.store.Role().Delete(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления роли",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Роль успешно удалена"})
}

// UserRoles ...

func (s *server) AddUserRoles(ctx *gin.Context) {
	var userRoles []model.UserRole

	err := ctx.ShouldBindJSON(&userRoles)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка чтения данных",
			"error": err.Error()})
		return
	}

	var addedUP []model.UserRole
	for _, req := range userRoles {
		userRoles, err := s.store.UserRole().Add(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления роли пользователя",

				"error": err.Error()})
			continue
		} else {
			addedUP = append(addedUP, userRoles)
		}
	}
	ctx.JSON(http.StatusCreated, addedUP)
}

func (s *server) GetUserRoles(ctx *gin.Context) {
	up, err := s.store.UserRole().All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка ролей пользователей",

			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, up)
}

func (s *server) GetUserRolesByUserId(ctx *gin.Context) {
	pUserID := ctx.Query("user_id")
	UserID, err := strconv.Atoi(pUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность UserId",
			"error": err.Error()})
		return
	}

	UserRole, err := s.store.UserRole().ByUserID(uint(UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка проектов ролей по UserId",

			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, UserRole)
}

func (s *server) GetUserRolesByRoleId(ctx *gin.Context) {
	pRoleID := ctx.Query("role_id")
	RoleID, err := strconv.Atoi(pRoleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность RoleId",
			"error": err.Error()})
		return
	}

	UserRole, err := s.store.UserRole().ByRoleID(uint(RoleID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка ролей пользователей по RoleId",

			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, UserRole)
}

func (s *server) GetUserRoleById(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	userRole, err := s.store.UserRole().ByID(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения роли пользователя по ID",

			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, userRole)
}

func (s *server) UpdateUserRole(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	var userRole model.UserRole
	err = ctx.ShouldBindJSON(&userRole)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность передаваемых данных",

			"error": err.Error()})
		return
	}

	userRole.ID = uint(ID)

	userRole, err = s.store.UserRole().Update(userRole)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления информации о роли пользователя",

			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Информация о роли пользователя успешно обновлена",
		"user_role": userRole})
}

func (s *server) DeleteUserRole(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	err = s.store.UserRole().Delete(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления роли пользователя",

			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Роль пользователя успешно удалена"})
}

// UserTeams ...

func (s *server) AddUserTeams(ctx *gin.Context) {
	var userTeams []model.UserTeam

	err := ctx.ShouldBindJSON(&userTeams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка чтения данных",
			"error": err.Error()})
		return
	}

	var addedUP []model.UserTeam
	for _, req := range userTeams {
		userTeams, err := s.store.UserTeam().Add(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления команды пользователя",

				"error": err.Error()})
			continue
		} else {
			addedUP = append(addedUP, userTeams)
		}
	}
	ctx.JSON(http.StatusCreated, addedUP)
}

func (s *server) GetUserTeams(ctx *gin.Context) {
	up, err := s.store.UserTeam().All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка ролей пользователей",

			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, up)
}

func (s *server) GetUserTeamsByUserId(ctx *gin.Context) {
	pUserID := ctx.Query("user_id")
	UserID, err := strconv.Atoi(pUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность UserId",
			"error": err.Error()})
		return
	}

	UserTeam, err := s.store.UserTeam().ByUserID(uint(UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка команд ролей по UserId",

			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, UserTeam)
}

func (s *server) GetUserTeamsByTeamId(ctx *gin.Context) {
	pTeamID := ctx.Query("team_id")
	TeamID, err := strconv.Atoi(pTeamID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность TeamId",
			"error": err.Error()})
		return
	}

	UserRole, err := s.store.UserTeam().ByTeamID(uint(TeamID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка команд пользователей по TeamId",

			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, UserRole)
}

func (s *server) GetUserTeamById(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	userRole, err := s.store.UserTeam().ByID(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения команды пользователя по ID",

			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, userRole)
}

func (s *server) UpdateUserTeam(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	var userTeam model.UserTeam
	err = ctx.ShouldBindJSON(&userTeam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность передаваемых данных",

			"error": err.Error()})
		return
	}

	userTeam.ID = uint(ID)

	userTeam, err = s.store.UserTeam().Update(userTeam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления информации о команде пользователя",

			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Информация о команде пользователя успешно обновлена",

		"user_role": userTeam})
}

func (s *server) DeleteUserTeamByID(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	err = s.store.UserTeam().Delete(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления команды пользователя",

			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Команда пользователя успешно удалена"})
}

func (s *server) DeleteUserTeam(ctx *gin.Context) {
	var user_team model.UserTeam
	err := ctx.ShouldBindJSON(&user_team)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность передаваемых данных",
			"error": err.Error()})
		return
	}

	err = s.store.UserTeam().DeleteUserTeam(user_team.TeamID, user_team.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления команды пользователя",

			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Команда пользователя успешно удалена"})
}

// EmployeeTeams ...

func (s *server) AddEmployeeTeams(ctx *gin.Context) {
	var employeeTeams []model.EmployeeTeam
	err := ctx.ShouldBindJSON(&employeeTeams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка чтения данных",
			"error": err.Error()})
		return
	}

	var addedET []model.EmployeeTeam
	for _, req := range employeeTeams {
		employeeTeams, err := s.store.EmployeeTeam().Add(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления команды пользователя",

				"error": err.Error()})
			continue
		} else {
			addedET = append(addedET, employeeTeams)
		}
	}
	ctx.JSON(http.StatusOK, addedET)
}

func (s *server) GetEmployeeTeams(ctx *gin.Context) {
	et, err := s.store.EmployeeTeam().All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка ролей пользователей",

			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, et)
}

func (s *server) GetEmployeeTeamsByEmployeeId(ctx *gin.Context) {
	pEmployeeID := ctx.Query("employee_id")
	EmployeeID, err := strconv.Atoi(pEmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность EmployeeID",
			"error": err.Error()})
		return
	}

	EmployeeTeam, err := s.store.EmployeeTeam().ByEmployeeID(uint(EmployeeID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка команд по EmployeeID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, EmployeeTeam)
}

func (s *server) GetEmployeeTeamsByTeamId(ctx *gin.Context) {
	pTeamID := ctx.Query("employee_id")
	TeamID, err := strconv.Atoi(pTeamID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность TeamID",
			"error": err.Error()})
		return
	}

	EmployeeTeam, err := s.store.EmployeeTeam().ByEmployeeID(uint(TeamID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения списка команд ролей по TeamID",
			"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, EmployeeTeam)
}

func (s *server) GetEmployeeTeamById(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Провеьте корректность ID",
			"error": err.Error()})
		return
	}

	et, err := s.store.Employee().ByID(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения команды пользователя по ID",

			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, et)
}

func (s *server) UpdateEmployeeTeam(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	var et model.EmployeeTeam
	err = ctx.ShouldBindJSON(&et)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность передаваемых данных",

			"error": err.Error()})
		return
	}

	et.ID = uint(ID)
	et, err = s.store.EmployeeTeam().Update(et)

	ctx.JSON(http.StatusOK, et)
}

func (s *server) DeleteEmployeeTeamByID(ctx *gin.Context) {
	pID := ctx.Query("id")
	ID, err := strconv.Atoi(pID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Проверьте корректность ID",
			"error": err.Error()})
		return
	}

	err = s.store.EmployeeTeam().Delete(uint(ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления команды пользователя",

			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Запись успешно удалена"})
}

func (s *server) DeleteEmployeeTeam(ctx *gin.Context) {
	var et model.EmployeeTeam
	err := ctx.ShouldBindJSON(&et)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность передаваемых данных",
			"error": err.Error()})
		return
	}

	err = s.store.EmployeeTeam().DeleteEmployeeTeam(et.EmployeeID, et.TeamID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления команды сотрудника",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Запись успешно удалена"})
}
