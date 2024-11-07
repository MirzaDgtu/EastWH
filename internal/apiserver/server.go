package apiserver

import (
	"eastwh/internal/model"
	"eastwh/internal/store"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
			userGroup.PUT("/update/", s.UpdateUser)
			userGroup.POST("/update/password/", s.UpdatePassword)
			userGroup.POST("/block/", s.BlockedUser)
			userGroup.GET("/profile/", s.GetUserProfile)
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
			employeeGroup.PUT("/update/", s.UpdateEmployee)
			employeeGroup.DELETE("/delete/", s.DeleteEmployee)
		}

		orderGroup := apiGroup.Group("/order")
		{
			orderGroup.GET("/", s.GetOrderByID)
			orderGroup.GET("/uid/", s.GetOrderByUID)
			orderGroup.PUT("/update/collector", s.UpdateOrderCollector)
		}

		ordersGroup := apiGroup.Group("/orders")
		{
			ordersGroup.GET("/user/", s.GetOrdersByUserId)
			ordersGroup.GET("/daterange/", s.GerOrderByDateRange)
			ordersGroup.POST("", s.AddOrders)
			ordersGroup.GET("", s.GetOrders)
		}

		teamGroup := apiGroup.Group("/team")
		{
			teamGroup.GET("/", s.GetTeamByID)
			teamGroup.PUT("/update/", s.UpdateTeam)
			teamGroup.DELETE("/delete/", s.DeleteTeam)
		}

		teamsGroup := apiGroup.Group("/teams")
		{
			teamsGroup.POST("", s.AddTeams)
			teamsGroup.GET("", s.GetTeams)
		}

		projectGroup := apiGroup.Group("/project")
		{
			projectGroup.GET("/id/", s.GetProjectById)
			projectGroup.DELETE("/delete/", s.DeleteProject)
			projectGroup.PUT("/update/", s.UpdateProject)
		}
		projectsGroup := apiGroup.Group("/projects")
		{
			projectsGroup.POST("", s.AddProject)
			projectsGroup.GET("", s.GetProjects)

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

// Employee...
func (s *server) AddEmployee(ctx *gin.Context) {
	var employees []model.Employee

	err := ctx.ShouldBindJSON(&employees)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка чтения данных",
			"error": err.Error()})
		return
	}

	var addedEmployees []model.Employee
	for _, req := range employees {
		employee, err := s.store.Employee().Add(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления сотрудника",
				"error": err.Error()})
			continue
		} else {
			addedEmployees = append(addedEmployees, employee)
		}
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
	var order model.Order

	err := ctx.ShouldBindJSON(&order)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка проверки данных заказа",
			"error": err.Error()})
		return
	}

	order, err = s.store.Order().Add(order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка создания заказа",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Заказ успешно создан",
		"order": order})
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

func (s *server) UpdateOrderCollector(ctx *gin.Context) {
	type request struct {
		OrderUID    uint `json:"order_uid"`
		KeeperID    uint `json:"keeper_id"`
		CollectorID uint `json:"collector_id"`
	}

	var req request

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Проверьте корректность передаваемых данных",
			"error": err.Error(),
		})
		return
	}

	err = s.store.Order().SetCollector(req.OrderUID, req.KeeperID, req.CollectorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления данных заказа",
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Данные заказа успешно обновлены"})
}

func (s *server) GerOrderByDateRange(ctx *gin.Context) {

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

