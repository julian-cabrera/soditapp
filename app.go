package main

import (
	"context"
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"

	"soditapp/backend/database"
	"soditapp/backend/repository"
	"soditapp/backend/service"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	db           *sql.DB
	connStr      string
	migrationsFS fs.FS
	seedsFS      fs.FS

	ZoneSvc          *service.ZoneService
	ClientSvc        *service.ClientService
	ProductSvc       *service.ProductService
	ClientProductSvc *service.ClientProductService
	BackupSvc        *service.BackupService
	ReportSvc        *service.ReportService
}

// NewApp creates a new App application struct.
// Services are allocated here (so Bind pointers are non-nil) and initialized in startup().
func NewApp(migrationsFS, seedsFS fs.FS) *App {
	return &App{
		migrationsFS:     migrationsFS,
		seedsFS:          seedsFS,
		ZoneSvc:          &service.ZoneService{},
		ClientSvc:        &service.ClientService{},
		ProductSvc:       &service.ProductService{},
		ClientProductSvc: &service.ClientProductService{},
		BackupSvc:        &service.BackupService{},
		ReportSvc:        &service.ReportService{},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	connStr := os.Getenv("SODITAPP_DB_URL")
	dbPath := ""
	if connStr == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			runtime.LogFatalf(ctx, "Failed to get config directory: %s", err.Error())
			return
		}
		dbDir := filepath.Join(configDir, "SoditApp")
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			runtime.LogFatalf(ctx, "Failed to create config directory: %s", err.Error())
			return
		}
		dbPath = filepath.Join(dbDir, "soditapp.db")
		connStr = "file:" + dbPath + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
	}
	a.connStr = connStr

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

	if err := database.RunSeeds(a.db, a.seedsFS); err != nil {
		runtime.LogFatalf(ctx, "Failed to run seeds: %s", err.Error())
		return
	}

	// Create repositories (internal, not exposed to frontend)
	zoneRepo := repository.NewZoneRepository(db)
	clientRepo := repository.NewClientRepository(db)
	productRepo := repository.NewProductRepository(db)
	clientProductRepo := repository.NewClientProductRepository(db)
	deliveryTxRepo := repository.NewDeliveryTransactionRepository(db)

	// Initialize services with runtime dependencies
	a.ZoneSvc.Init(ctx, zoneRepo)
	a.ClientSvc.Init(ctx, clientRepo)
	a.ProductSvc.Init(ctx, productRepo)
	a.ClientProductSvc.Init(ctx, db, clientProductRepo, zoneRepo, deliveryTxRepo)

	// BackupService needs a reopen callback for import
	reopenDB := func() error {
		if a.db != nil {
			a.db.Close()
		}
		newDB, err := database.Connect(a.connStr)
		if err != nil {
			return err
		}
		a.db = newDB
		// Reinitialize repos with the new connection
		zoneRepo = repository.NewZoneRepository(newDB)
		clientRepo = repository.NewClientRepository(newDB)
		productRepo = repository.NewProductRepository(newDB)
		clientProductRepo = repository.NewClientProductRepository(newDB)
		deliveryTxRepo = repository.NewDeliveryTransactionRepository(newDB)
		// Reinitialize all services with new repos
		a.ZoneSvc.Reinit(zoneRepo)
		a.ClientSvc.Reinit(clientRepo)
		a.ProductSvc.Reinit(productRepo)
		a.ClientProductSvc.Reinit(newDB, clientProductRepo, zoneRepo, deliveryTxRepo)
		a.ReportSvc.Reinit(deliveryTxRepo)
		return nil
	}

	a.BackupSvc.Init(ctx, db, dbPath, reopenDB)
	a.ReportSvc.Init(ctx, deliveryTxRepo)
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}
