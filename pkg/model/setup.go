package model

import (
	"encoding/json"
	"fmt"
	l "log"
	"os"
	"strings"
	"time"

	"github.com/misoboy/cdule/pkg"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB gorm DB
var DB *gorm.DB

// CduleRepos repositories
var CduleRepos *Repositories

// Repositories struct
type Repositories struct {
	CduleRepository CduleRepository
	DB              *gorm.DB
}

// ConnectDataBase to create a database connection
func ConnectDataBase(param []string) (*pkg.CduleConfig, error) {
	cduleConfig, err := readConfig(param)
	if nil != err {
		log.Error(err)
		panic("Failed to read config!")
	}
	printConfig(cduleConfig)
	var db *gorm.DB
	if cduleConfig.Cdule.Type == string(pkg.DATABASE) {
		if cduleConfig.Database.Type == "postgres" {
			db = postgresConn(cduleConfig.Database.Url)
		} else if cduleConfig.Database.Type == "mysql" {
			db = mysqlConn(cduleConfig.Database.Url)
		}
	} else if cduleConfig.Cdule.Type == string(pkg.MEMORY) {
		db = sqliteConn(cduleConfig.Database.Url)
	}

	logLevel := logger.Error
	if len(param) > 2 && param[2] != "errorLogType" {
		logLevel = logger.Info
	}
	// Set LogLevel to `logger.Silent` to stop logging sqls
	sqlLogger := logger.New(
		l.New(os.Stdout, "\r\n", l.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	db.Logger = sqlLogger
	Migrate(db)
	DB = db

	// Initialise CduleRepositories
	CduleRepos = &Repositories{
		CduleRepository: NewCduleRepository(db),
		DB:              db,
	}
	return cduleConfig, err
}

func postgresConn(dbDSN string) (db *gorm.DB) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbDSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Errorf("Error Connecting Postgressql %s, %s", dbDSN, err.Error())
		panic("Failed to connect to database! " + dbDSN)
	}
	return db
}

func mysqlConn(dbDSN string) (db *gorm.DB) {
	// splitting DSN to only use the string after mysql://
	splitDSN := strings.Split(dbDSN, "mysql://")
	db, err := gorm.Open(mysql.Open(splitDSN[1]), &gorm.Config{})
	if err != nil {
		log.Errorf("Error Connecting MySQL %s, %s", dbDSN, err.Error())
		panic("Failed to connect to database! " + dbDSN)
	}
	return db
}

func sqliteConn(dbDSN string) (db *gorm.DB) {
	//db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	//db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})

	// If you would use file based as mentioned above db
	db, err := gorm.Open(sqlite.Open(dbDSN), &gorm.Config{})
	if err != nil {
		log.Error(dbDSN)
		panic("Failed to connect to database! " + dbDSN)
	}
	return db
}

// Migrate database schema
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Job{})
	db.AutoMigrate(&JobHistory{})
	db.AutoMigrate(&Schedule{})
	db.AutoMigrate(&Worker{})
}

func printConfig(config *pkg.CduleConfig) {
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("Configuration %s\n", string(configJSON))
}
