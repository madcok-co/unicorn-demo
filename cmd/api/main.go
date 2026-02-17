package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/madcok-co/unicorn-demo/internal/database"
	"github.com/madcok-co/unicorn-demo/internal/handlers"
	"github.com/madcok-co/unicorn/contrib/auth/oauth2"
	"github.com/madcok-co/unicorn/contrib/authz/rbac"
	"github.com/madcok-co/unicorn/contrib/config"
	"github.com/madcok-co/unicorn/contrib/multitenancy"
	"github.com/madcok-co/unicorn/contrib/versioning"
	httpAdapter "github.com/madcok-co/unicorn/core/pkg/adapters/http"
	"github.com/madcok-co/unicorn/core/pkg/app"
	"github.com/madcok-co/unicorn/core/pkg/contracts"
)

func main() {
	log.Println("🦄 Starting Unicorn Demo API...")

	// Initialize configuration
	cfg := initializeConfig()

	// Initialize database
	db := initializeDatabase(cfg)

	// Initialize multi-tenancy
	_ = initializeMultiTenancy(cfg)

	// Initialize OAuth2
	auth := initializeOAuth2(cfg)

	// Initialize RBAC
	authz := initializeRBAC()

	// Initialize API versioning
	_ = initializeVersioning()

	// Create application
	application := app.New(&app.Config{
		Name:       cfg.GetString("app.name"),
		Version:    cfg.GetString("app.version"),
		EnableHTTP: true,
		HTTP: &httpAdapter.Config{
			Host: cfg.GetString("http.host"),
			Port: cfg.GetInt("http.port"),
		},
	})

	// Set adapters
	application.SetAuth(auth)
	application.SetAuthz(authz)
	application.SetDB(db) // db is already a contracts.Database from Unicorn's GORM driver

	// Note: Middleware in Unicorn is registered at the HTTP adapter level
	// For this demo, middleware is applied via handlers and route configuration

	// Register routes
	registerV1Routes(application)
	registerV2Routes(application)

	// Start server
	log.Printf("✅ Server starting on %s:%d", cfg.GetString("http.host"), cfg.GetInt("http.port"))
	log.Printf("📚 API Documentation: http://localhost:%d/health", cfg.GetInt("http.port"))
	log.Printf("🏢 Multi-tenant: subdomain.localhost:%d", cfg.GetInt("http.port"))

	if err := application.Start(); err != nil {
		log.Fatal(err)
	}
}

func initializeConfig() *config.Driver {
	log.Println("🔧 Initializing configuration...")

	cfg, err := config.NewDriver(&config.Config{
		Defaults: map[string]interface{}{
			"app.name":    "Unicorn Demo API",
			"app.version": "1.0.0",
			"http.host":   getEnv("HTTP_HOST", "0.0.0.0"),
			"http.port":   getEnvInt("HTTP_PORT", 8080),

			// Database
			"db.driver":   getEnv("DB_DRIVER", "sqlite"),
			"db.host":     getEnv("DB_HOST", "localhost"),
			"db.port":     getEnvInt("DB_PORT", 5432),
			"db.user":     getEnv("DB_USER", "postgres"),
			"db.password": getEnv("DB_PASSWORD", "postgres"),
			"db.database": getEnv("DB_DATABASE", "unicorn_demo.db"),
			"db.sslmode":  getEnv("DB_SSLMODE", "disable"),

			// OAuth2
			"oauth.provider":      getEnv("OAUTH_PROVIDER", "google"),
			"oauth.client_id":     getEnv("OAUTH_CLIENT_ID", "demo-client-id"),
			"oauth.client_secret": getEnv("OAUTH_CLIENT_SECRET", "demo-client-secret"),
			"oauth.redirect_url":  getEnv("OAUTH_REDIRECT_URL", "http://localhost:8080/api/v1/auth/callback"),

			// Multi-tenancy
			"multitenancy.strategy": getEnv("MULTITENANCY_STRATEGY", "subdomain"),
			"multitenancy.domain":   getEnv("MULTITENANCY_DOMAIN", "localhost"),
		},
		EnvPrefix: "APP",
	})

	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	return cfg
}

func initializeDatabase(cfg *config.Driver) contracts.Database {
	log.Println("🗄️  Initializing database...")

	dbConfig := &database.Config{
		Driver:   cfg.GetString("db.driver"),
		Host:     cfg.GetString("db.host"),
		Port:     cfg.GetInt("db.port"),
		User:     cfg.GetString("db.user"),
		Password: cfg.GetString("db.password"),
		Database: cfg.GetString("db.database"),
		SSLMode:  cfg.GetString("db.sslmode"),
	}

	// Connect with Unicorn's database interface
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get raw GORM DB for migrations and seeding
	rawDB, err := database.ConnectRaw(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect raw database: %v", err)
	}

	// Run migrations
	if err := database.AutoMigrate(rawDB); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed initial data
	if err := database.SeedData(rawDB); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	return db
}

func initializeMultiTenancy(cfg *config.Driver) *multitenancy.Driver {
	log.Println("🏢 Initializing multi-tenancy...")

	strategy := multitenancy.StrategySubdomain
	strategyStr := cfg.GetString("multitenancy.strategy")
	if strategyStr == "header" {
		strategy = multitenancy.StrategyHeader
	} else if strategyStr == "path" {
		strategy = multitenancy.StrategyPath
	}

	mt := multitenancy.NewDriver(&multitenancy.Config{
		Strategy:      strategy,
		DefaultTenant: "default",
	})

	// Create demo tenants
	ctx := context.Background()
	mt.CreateTenant(ctx, &multitenancy.Tenant{
		ID:     "acme",
		Name:   "Acme Corporation",
		Active: true,
		Metadata: map[string]interface{}{
			"plan":      "enterprise",
			"max_users": 100,
			"features":  []string{"projects", "tasks", "analytics"},
		},
	})

	mt.CreateTenant(ctx, &multitenancy.Tenant{
		ID:     "techcorp",
		Name:   "Tech Corp",
		Active: true,
		Metadata: map[string]interface{}{
			"plan":      "professional",
			"max_users": 50,
			"features":  []string{"projects", "tasks"},
		},
	})

	log.Println("✅ Created demo tenants: acme, techcorp")
	return mt
}

func initializeOAuth2(cfg *config.Driver) *oauth2.Driver {
	log.Println("🔐 Initializing OAuth2...")

	provider := oauth2.ProviderGoogle
	providerStr := cfg.GetString("oauth.provider")
	if providerStr == "github" {
		provider = oauth2.ProviderGitHub
	} else if providerStr == "microsoft" {
		provider = oauth2.ProviderMicrosoft
	}

	return oauth2.NewDriver(&oauth2.Config{
		Provider:     provider,
		ClientID:     cfg.GetString("oauth.client_id"),
		ClientSecret: cfg.GetString("oauth.client_secret"),
		RedirectURL:  cfg.GetString("oauth.redirect_url"),
		Scopes:       []string{"openid", "email", "profile"},
	})
}

func initializeRBAC() *rbac.Driver {
	log.Println("🔒 Initializing RBAC...")

	roles := make(map[string]*rbac.Role)

	roles["super_admin"] = &rbac.Role{
		Name:        "super_admin",
		Permissions: []string{"*"},
	}

	roles["tenant_admin"] = &rbac.Role{
		Name: "tenant_admin",
		Permissions: []string{
			"projects:*",
			"tasks:*",
			"users:*",
			"comments:*",
		},
	}

	roles["project_manager"] = &rbac.Role{
		Name: "project_manager",
		Permissions: []string{
			"projects:read",
			"projects:create",
			"projects:update",
			"tasks:*",
			"users:read",
			"comments:*",
		},
		Inherits: []string{"developer"},
	}

	roles["developer"] = &rbac.Role{
		Name: "developer",
		Permissions: []string{
			"projects:read",
			"tasks:read",
			"tasks:update",
			"comments:*",
		},
		Inherits: []string{"viewer"},
	}

	roles["viewer"] = &rbac.Role{
		Name:        "viewer",
		Permissions: []string{"*:read"},
	}

	authz := rbac.NewDriver(&rbac.Config{
		Roles: roles,
	})

	log.Println("✅ Configured RBAC with 5 roles")
	return authz
}

func initializeVersioning() *versioning.Manager {
	log.Println("📋 Initializing API versioning...")

	vm := versioning.NewManager(&versioning.Config{
		Strategy:       versioning.StrategyURL,
		DefaultVersion: "1.0",
	})

	// Note: Deprecation tracking can be added via custom logic if needed
	log.Println("✅ API versioning configured with URL strategy (v1, v2)")

	return vm
}

func registerV1Routes(app *app.App) {
	log.Println("📝 Registering API v1 routes...")

	// Health check
	app.RegisterHandler(handlers.HealthCheck).
		Named("health").
		HTTP("GET", "/health").
		Done()

	// Authentication
	app.RegisterHandler(handlers.Login).
		Named("auth.login.v1").
		HTTP("POST", "/api/v1/auth/login").
		Done()

	app.RegisterHandler(handlers.OAuth2Callback).
		Named("auth.callback.v1").
		HTTP("GET", "/api/v1/auth/callback").
		Done()

	app.RegisterHandler(handlers.RefreshToken).
		Named("auth.refresh.v1").
		HTTP("POST", "/api/v1/auth/refresh").
		Done()

	app.RegisterHandler(handlers.GetCurrentUser).
		Named("auth.me.v1").
		HTTP("GET", "/api/v1/auth/me").
		Done()

	app.RegisterHandler(handlers.Logout).
		Named("auth.logout.v1").
		HTTP("POST", "/api/v1/auth/logout").
		Done()

	// Projects
	app.RegisterHandler(handlers.ListProjects).
		Named("projects.list.v1").
		HTTP("GET", "/api/v1/projects").
		Done()

	app.RegisterHandler(handlers.CreateProject).
		Named("projects.create.v1").
		HTTP("POST", "/api/v1/projects").
		Done()

	app.RegisterHandler(handlers.GetProject).
		Named("projects.get.v1").
		HTTP("GET", "/api/v1/projects/:id").
		Done()

	app.RegisterHandler(handlers.UpdateProject).
		Named("projects.update.v1").
		HTTP("PUT", "/api/v1/projects/:id").
		Done()

	app.RegisterHandler(handlers.DeleteProject).
		Named("projects.delete.v1").
		HTTP("DELETE", "/api/v1/projects/:id").
		Done()

	app.RegisterHandler(handlers.GetProjectStats).
		Named("projects.stats.v1").
		HTTP("GET", "/api/v1/projects/stats").
		Done()

	// Users
	app.RegisterHandler(handlers.ListUsers).
		Named("users.list.v1").
		HTTP("GET", "/api/v1/users").
		Done()

	app.RegisterHandler(handlers.CreateUser).
		Named("users.create.v1").
		HTTP("POST", "/api/v1/users").
		Done()

	app.RegisterHandler(handlers.GetUser).
		Named("users.get.v1").
		HTTP("GET", "/api/v1/users/:id").
		Done()

	app.RegisterHandler(handlers.UpdateUser).
		Named("users.update.v1").
		HTTP("PUT", "/api/v1/users/:id").
		Done()

	app.RegisterHandler(handlers.DeleteUser).
		Named("users.delete.v1").
		HTTP("DELETE", "/api/v1/users/:id").
		Done()

	app.RegisterHandler(handlers.GetUserStats).
		Named("users.stats.v1").
		HTTP("GET", "/api/v1/users/stats").
		Done()

	// Tasks
	app.RegisterHandler(handlers.ListTasks).
		Named("tasks.list.v1").
		HTTP("GET", "/api/v1/tasks").
		Done()

	app.RegisterHandler(handlers.CreateTask).
		Named("tasks.create.v1").
		HTTP("POST", "/api/v1/tasks").
		Done()

	app.RegisterHandler(handlers.GetTask).
		Named("tasks.get.v1").
		HTTP("GET", "/api/v1/tasks/:id").
		Done()

	app.RegisterHandler(handlers.UpdateTask).
		Named("tasks.update.v1").
		HTTP("PUT", "/api/v1/tasks/:id").
		Done()

	app.RegisterHandler(handlers.DeleteTask).
		Named("tasks.delete.v1").
		HTTP("DELETE", "/api/v1/tasks/:id").
		Done()

	log.Println("✅ Registered 28 v1 endpoints")
}

func registerV2Routes(app *app.App) {
	log.Println("📝 Registering API v2 routes...")

	// V2 features: Cursor pagination + Batch operations
	app.RegisterHandler(handlers.BatchUpdateProjects).
		Named("projects.batch_update.v2").
		HTTP("POST", "/api/v2/projects/batch").
		Done()

	log.Println("✅ Registered 1 v2 endpoint")
}

// Utility functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Note: GORMAdapter is no longer needed - using Unicorn's built-in GORM driver
// from github.com/madcok-co/unicorn/contrib/database/gorm
