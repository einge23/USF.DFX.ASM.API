package routes

import (
	"gin-api/controllers"
	"gin-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/refreshToken", controllers.RefreshToken)
		}
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			printers := protected.Group("/printers")
			{
				printers.GET("/getPrinters", controllers.GetPrinters)
				printers.PUT("/reservePrinter", controllers.ReservePrinter)
			}
			users := protected.Group("/users")
			{
				users.GET("/reservations/:userID",
					middleware.UserOwnershipPermission(),
					controllers.GetUserReservations,
				)
				users.GET("/activeReservations/:userID",
					middleware.UserOwnershipPermission(),
					controllers.GetActiveUserReservations,
				)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminPermission())
			{
				users := admin.Group("/users")
				{
					users.POST("/create", controllers.CreateUser)
					users.POST("/getUser", controllers.GetUserById)
					users.PUT("/setTrained/:userID", controllers.SetUserTrained)
					users.PUT("/setExecutiveAccess/:userID", controllers.SetUserExecutiveAccess)
					users.PUT("/addWeeklyMinutes/:userID", controllers.AddUserWeeklyMinutes)
					users.PUT("/setBanTime/:userID", controllers.SetUserBanTime)
				}
				printers := admin.Group("/printers")
				{
					printers.PUT("/setExecutive/:printerID", controllers.SetPrinterExecutive)
				}
			}
		}
	}
}
