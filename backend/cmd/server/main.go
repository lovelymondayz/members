package main

import (
	"log"
	"time"

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
	seedDummyData()

	middleware.SetJWTSecret(cfg.JWTSecret)

	userRepo := repository.NewUserRepository()
	storeRepo := repository.NewStoreRepository()
	memberRepo := repository.NewMemberRepository()
	invoiceRepo := repository.NewInvoiceRepository()
	paymentRepo := repository.NewPaymentRepository()

	authSvc := service.NewAuthService(userRepo, storeRepo, cfg)

	authHandler := handler.NewAuthHandler(authSvc)
	adminHandler := handler.NewAdminHandler(authSvc)
	storeHandler := handler.NewStoreHandler(storeRepo)
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

	stores := protected.Group("/stores")
	stores.Use(middleware.RoleRequired("super_admin"))
	stores.GET("", storeHandler.GetStores)
	stores.POST("", storeHandler.CreateStore)
	stores.PUT("/:id", storeHandler.UpdateStore)
	stores.DELETE("/:id", storeHandler.DeleteStore)

	admin := protected.Group("/admin")
	admin.Use(middleware.RoleRequired("super_admin"))
	admin.GET("/admins", adminHandler.GetAdmins)
	admin.POST("/admins", adminHandler.CreateAdmin)
	admin.PUT("/admins/:id", adminHandler.UpdateAdmin)
	admin.DELETE("/admins/:id", adminHandler.DeleteAdmin)
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

func seedDummyData() {
	// Skip if data already exists
	var count int64
	config.DB.Model(&models.User{}).Count(&count)
	if count > 3 {
		log.Println("Seed: data exists, skipping dummy seed")
		return
	}

	log.Println("Seed: creating dummy data...")

	// ─── Admin user + store ───
	hashedPass, _ := service.HashPassword("admin123")
	adminUser := &models.User{
		RoleID:   2,
		Email:    "admin@tokotest.id",
		Password: hashedPass,
		Name:     "Toko Test Admin",
	}
	config.DB.Create(adminUser)
	adminStore := &models.Store{
		AdminID:      adminUser.UserID,
		Name:         "Toko Test",
		CardColorHex: "#059669",
	}
	config.DB.Create(adminStore)

	// ─── Members for admin store ───
	members := []models.Member{
		{
			StoreID:    adminStore.StoreID,
			MemberCode: "M001",
			Tier:       "gold",
		},
		{
			StoreID:    adminStore.StoreID,
			MemberCode: "M002",
			Tier:       "silver",
		},
		{
			StoreID:    adminStore.StoreID,
			MemberCode: "M003",
			Tier:       "premium",
		},
	}
	for i := range members {
		config.DB.Create(&members[i])
	}

	// ─── Invoices for first two members ───
	now := time.Now()
	invoices := []models.Invoice{
		{
			StoreID:       adminStore.StoreID,
			MemberID:      members[0].MemberID,
			InvoiceNumber: "INV-001",
			Amount:        150000,
			Description:   "Pembelian paket Gold",
			Status:        models.InvoicePaid,
			DueDate:       &now,
		},
		{
			StoreID:       adminStore.StoreID,
			MemberID:      members[1].MemberID,
			InvoiceNumber: "INV-002",
			Amount:        85000,
			Description:   "Pembelian paket Silver",
			Status:        models.InvoiceSent,
			DueDate:       &now,
		},
		{
			StoreID:       adminStore.StoreID,
			MemberID:      members[0].MemberID,
			InvoiceNumber: "INV-003",
			Amount:        250000,
			Description:   "Top up saldo",
			Status:        models.InvoiceDraft,
			DueDate:       &now,
		},
	}
	for i := range invoices {
		config.DB.Create(&invoices[i])
	}

	// ─── Payment for first invoice ───
	config.DB.Create(&models.Payment{
		InvoiceID: invoices[0].InvoiceID,
		Amount:    150000,
		Method:    "cash",
		Reference: "PAY-001",
		Note:      "Lunas tunai",
		PaidAt:    now,
	})

	log.Println("Seed: dummy data created ✅")
	log.Println("  Super Admin: admin@gantengbanget.id / admin123")
	log.Println("  Store Admin: admin@tokotest.id  / admin123")
	log.Println("  Members: M001 (gold), M002 (silver), M003 (premium)")
	log.Println("  Invoices: INV-001 (paid), INV-002 (sent), INV-003 (draft)")
}
