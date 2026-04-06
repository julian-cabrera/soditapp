package main

import (
	"context"
	"database/sql"
	"io/fs"
	"os"

	"soditapp/backend/database"
	"soditapp/backend/repository"
	"soditapp/backend/service"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	db           *sql.DB
	migrationsFS fs.FS

	ZoneSvc          *service.ZoneService
	ClientSvc        *service.ClientService
	ProductSvc       *service.ProductService
	ClientProductSvc *service.ClientProductService
}

// NewApp creates a new App application struct.
// Services are allocated here (so Bind pointers are non-nil) and initialized in startup().
func NewApp(migrationsFS fs.FS) *App {
	return &App{
		migrationsFS:     migrationsFS,
		ZoneSvc:          &service.ZoneService{},
		ClientSvc:        &service.ClientService{},
		ProductSvc:       &service.ProductService{},
		ClientProductSvc: &service.ClientProductService{},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	connStr := os.Getenv("SODITAPP_DB_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	db, err := database.Connect(connStr)
	if err != nil {
		runtime.LogFatalf(ctx, "Failed to connect to database: %s", err.Error())
		return
	}
	a.db = db

	if err := database.RunMigrations(a.db, a.migrationsFS); err != nil {
		runtime.LogFatalf(ctx, "Failed to run migrations: %s", err.Error())
		return
	}

	// Create repositories (internal, not exposed to frontend)
	zoneRepo := repository.NewZoneRepository(db)
	clientRepo := repository.NewClientRepository(db)
	productRepo := repository.NewProductRepository(db)
	clientProductRepo := repository.NewClientProductRepository(db)

	// Initialize services with runtime dependencies
	a.ZoneSvc.Init(ctx, zoneRepo)
	a.ClientSvc.Init(ctx, clientRepo)
	a.ProductSvc.Init(ctx, productRepo)
	a.ClientProductSvc.Init(ctx, db, clientProductRepo, zoneRepo)
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}
