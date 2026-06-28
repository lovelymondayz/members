package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lovelymondayz/members/backend/src/config"
	"github.com/lovelymondayz/members/backend/src/handler"
	"github.com/lovelymondayz/members/backend/src/middleware"
	"github.com/lovelymondayz/members/backend/src/models"
	"github.com/lovelymondayz/members/backend/src/repository"
	"github.com/lovelymondayz/members/backend/src/service"
)

func main() {
	cfg := config.Load()
	config.Connect(cfg)

	config.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Store{},
		&models.Member{},
		&models.Invoice{},
		&models.Payment{},
	)
	seedRoles()

	middleware.SetJWTSecret(cfg.JWTSecret)

	userRepo := repository.NewUserRepository()
	storeRepo := repository.NewStoreRepository()
	memberRepo := repository.NewMemberRepository()
	invoiceRepo := repository.NewInvoiceRepository()
	paymentRepo := repository.NewPaymentRepository()

	authSvc := service.NewAuthService(userRepo, storeRepo, cfg)

	authHandler := handler.NewAuthHandler(authSvc)
	adminHandler := handler.NewAdminHandler(authSvc)
	memberHandler := handler.NewMemberHandler(memberRepo)
	invoiceHandler := handler.NewInvoiceHandler(invoiceRepo, paymentRepo)

	router := gin.Default()

	allowedOrigins := map[string]bool{
		"https://members.arjism.com": true,
		"http://localhost:3003":      true,
		"http://localhost:5173":      true,
	}
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.GET("/health", handler.HealthCheck)

	public := router.Group("/api/auth")
	public.POST("/register", authHandler.Register)
	public.POST("/login", authHandler.Login)
	public.GET("/google", authHandler.GoogleLogin)
	public.GET("/google/callback", authHandler.GoogleCallback)

	protected := router.Group("/api")
	protected.Use(middleware.AuthRequired())
	protected.GET("/auth/me", authHandler.Me)
	protected.POST("/auth/logout", authHandler.Logout)

	members := protected.Group("/members")
	members.Use(middleware.RoleRequired("admin", "super_admin"))
	members.GET("", memberHandler.GetMembers)
	members.POST("", memberHandler.CreateMember)
	members.GET("/:id", memberHandler.GetMember)
	members.PUT("/:id", memberHandler.UpdateMember)
	members.DELETE("/:id", memberHandler.DeleteMember)
	members.GET("/:id/card", memberHandler.GetMemberCard)

	invoices := protected.Group("/invoices")
	invoices.Use(middleware.RoleRequired("admin", "super_admin"))
	invoices.GET("", invoiceHandler.GetInvoices)
	invoices.POST("", invoiceHandler.CreateInvoice)
	invoices.GET("/:id", invoiceHandler.GetInvoice)
	invoices.PUT("/:id", invoiceHandler.UpdateInvoice)
	invoices.DELETE("/:id", invoiceHandler.DeleteInvoice)
	invoices.POST("/:id/pay", invoiceHandler.RecordPayment)
	invoices.GET("/member", invoiceHandler.GetMemberInvoices)

	admin := protected.Group("/admin")
	admin.Use(middleware.RoleRequired("super_admin"))
	admin.GET("/admins", adminHandler.GetAdmins)
	admin.POST("/admins", adminHandler.CreateAdmin)
	admin.GET("/dashboard", adminHandler.DashboardStats)

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func seedRoles() {
	roles := []models.Role{
		{RoleID: 1, Name: "super_admin"},
		{RoleID: 2, Name: "admin"},
		{RoleID: 3, Name: "member"},
	}
	for _, role := range roles {
		config.DB.FirstOrCreate(&role, models.Role{RoleID: role.RoleID})
	}
}
