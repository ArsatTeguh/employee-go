// routers/router.go
package routers

import (
	"backend/controller"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(server *gin.Engine, databases *gorm.DB) *gin.Engine {

	user := controller.NewUserService(databases)
	employee := controller.NewServiceEmployee(databases)
	position := controller.PositionContoroller{DB: databases}
	project := controller.ProjectController{DB: databases}
	attedance := controller.AttedanceController{DB: databases}
	payroll := controller.PayrollController{DB: databases}
	leave := controller.LeaveController{DB: databases}
	task := controller.TaskController{DB: databases}
	wallet := controller.WalletContoroller{DB: databases}

	// Public API
	r := server.Group("api/v1")
	r.POST("/auth/register", user.Register)
	r.POST("/auth/login", user.Login)
	r.GET("/auth/token", user.RefreshToken)

	// Private API
	protected := server.Group("api/v1")
	protected.Use(middlewares.JwtAuthMiddleware())

	protected.GET("/auth/user-one", user.GetOneUser)
	protected.GET("/auth/user", user.GetAll)
	protected.PATCH("/auth/user/:id", user.Update)
	protected.POST("/auth/logout", user.Logout)

	protected.GET("/employees", employee.GetAllEmployee)
	protected.GET("/employee-popup", employee.PopupEmployee)
	protected.GET("/employee-one", employee.GetOneEmployee)
	protected.POST("/employee", employee.SaveEmployee)
	protected.PATCH("/upload", employee.UploadProfile)
	protected.PATCH("/employee", employee.Update)
	protected.DELETE("/employee-delete/:id", employee.Delete)
	protected.GET("/profile/:id", employee.GetProfile)

	protected.GET("/positions", position.GetAllPosition)
	protected.GET("/position/:id", position.GetOnePosition)
	protected.POST("/positionAllById", position.GetAllPositionById)
	protected.POST("/position", position.SavePosition)
	protected.GET("/position-project", position.GetByProject)
	// protected.PATCH("/position/:id", position.Update)
	protected.PATCH("/positions", position.SyncPositions)

	protected.GET("/projects", project.GetAllProject)
	protected.GET("/project-popup", project.PopupProject)
	protected.GET("/project/:id", project.GetOne)
	protected.POST("/project", project.Saved)
	protected.PATCH("/project/:id", project.Update)
	protected.DELETE("/project/:id", project.Delete)
	protected.GET("/project-master", project.ProjectMaster)

	protected.GET("/attedances", attedance.GetAll)
	protected.GET("/attedance/:id", attedance.GetOne)
	protected.POST("/attedance", attedance.Created)
	protected.PATCH("/attedance", attedance.Update)
	protected.POST("/chekout", attedance.EmployeeCheckout)

	protected.GET("/payrolls", payroll.GetAll)
	protected.POST("/payroll", payroll.Payroll)
	protected.POST("/payroll-email", payroll.EmailPayslip)
	protected.POST("/payroll-excel", payroll.ExportExcelHandler)

	protected.GET("/leaves", leave.GetAll)
	protected.GET("/leaves-employee/:id", leave.GetAllByEmployee)
	protected.GET("/leave-employee", leave.GetOneByEmployee)
	protected.POST("/leave", leave.Created)
	protected.PATCH("/leave/:id", leave.Approve)
	protected.DELETE("/leave/:id", leave.Delete)

	protected.PATCH("/task/:id", task.Update)
	protected.GET("/task/:projectId", task.GetOne)
	protected.POST("/task", task.SaveTask)
	protected.DELETE("/task/:id", task.Delete)

	protected.GET("/wallet/:id", wallet.GetOneWallet)
	protected.PATCH("/wallet/:id", wallet.UpdatedWallet)

	return server
}
