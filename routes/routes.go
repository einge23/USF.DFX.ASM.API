package routes

import (
	"gin-api/controllers"
	"gin-api/middleware"

	"github.com/gin-gonic/gin"
)

//routes here are sort of cascaded / nested. example: on line 28 we see reservePrinter, this is nested inside
//of /printers, which itself is nested inside of /api. So the route is localhost:3000/api/printers/reservePrinter

//defines all routes and groups for the router
func SetupRouter(r *gin.Engine) {
	api := r.Group("/api") //master group
	{
		auth := api.Group("/auth") //authorization routes: logs in / gives tokens
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/refreshToken", controllers.RefreshToken)
		}
		protected := api.Group("") //protected: only available to users with a valid token
		protected.Use(middleware.AuthMiddleware())
		{
			printers := protected.Group("/printers") //user-level printer routes
			{
				printers.GET("/getPrinters", controllers.GetPrinters)
				printers.PUT("/reservePrinter", controllers.ReservePrinter)
			}
			users := protected.Group("/users") //user-level user routes
			{
				users.GET("/reservations/:userID",
					controllers.GetUserReservations,
				)
				users.GET("/activeReservations/:userID",
					controllers.GetActiveUserReservations,
				)
				users.PUT("/cancelActiveReservation",
					controllers.CancelActiveReservation,
				)
			}
			settings := protected.Group("/settings") //user-level settings routes
			{
				settings.GET("/getSettings", controllers.GetSettings)
			}
			reservations := protected.Group("/reservations") //user-level reservations routes
			{
				reservations.GET("/getActiveReservations", controllers.GetActiveReservations)
			}

			//Admin routes
			admin := protected.Group("/admin") //admin: only available to admin users
			admin.Use(middleware.AdminPermission())
			{
				users := admin.Group("/users") //admin-level user routes
				{
					users.POST("/create", controllers.CreateUser)
					users.POST("/getUser", controllers.GetUserById)
					users.PUT("/setTrained/:userID", controllers.SetUserTrained)
					users.PUT("/setExecutiveAccess/:userID", controllers.SetUserExecutiveAccess)
					users.PUT("/addWeeklyMinutes/:userID", controllers.AddUserWeeklyMinutes)
					users.PUT("/setBanTime/:userID", controllers.SetUserBanTime)
				}
				printers := admin.Group("/printers") //admin-level printers routes
				{
					printers.POST("/create", controllers.AddPrinter)
					printers.PUT("/setExecutive/:printerID", controllers.SetPrinterExecutive)
					printers.PUT("/update/:printerID", controllers.UpdatePrinter)
				}
				settings := admin.Group("/settings") //admin-level settings routes
				{
					settings.PUT("/setSettings", controllers.SetSettings)
				}
			}

		}
	}
}
