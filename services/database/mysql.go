package database

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"app/config"
)

var _database *gorm.DB
var once sync.Once

func MySqlInitDatabase() {
	once.Do(func() {
		conf := config.GetConfig().Database

		//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]       //db_user:password@tcp(localhost:3306)/my_db
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", conf.UserName, conf.Password, conf.Host, conf.Port, conf.DatabaseName)
		db, err := gorm.Open("mysql", connectionString)

		if err != nil {
			panic("failed to connect database")
		}

		if err := db.DB().Ping(); err != nil {
			panic("failed to ping database")
		}

		db.DB().SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetimeMin))
		db.DB().SetMaxIdleConns(conf.MaxIdleConnections) //https://github.com/jinzhu/gorm/issues/246
		db.DB().SetMaxOpenConns(conf.MaxOpenConnections)
		db.LogMode(conf.Debug)

		_database = db
	})
}

type InTransaction func(tx *gorm.DB) error

func MySqlDoInTransaction(fn InTransaction) error {
	//Start Transaction and Check if database begin-transactions causes any errors
	tx := _database.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	//execute
	err := fn(tx)

	// if a panic occurred, rollback and re-panic
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	}

	// check if something went wrong, during execution of codes inside DoInTransaction block. If yes , rollback
	if err != nil {
		rollBackError := tx.Rollback().Error
		if rollBackError != nil {
			return rollBackError
		}
		return err
	} else { //Everything OK: Commit
		//if commit error occurred, return commit error, otherwise return the original error
		if err = tx.Commit().Error; err != nil {
			return err
		}
	}

	return nil
}

func MySqlClose() {
	if err := _database.Close(); err != nil {
		log.Errorln("error closing database", err.Error())
	} else {
		log.Debugln("database connection closed")
	}
}

const (
	OAuthAccessTokenTable      = "oauth_access_tokens"
	OAuthClientTable           = "oauth_clients"
	OAuthSessionsTable         = "oauth_sessions"
	OathAccessTokenScopesTable = "oauth_access_token_scopes"
	UserTable                  = "users"
	UserSocialProfile          = "user_social_profiles"
	DriverMetaTable            = "drivers_meta"
	UserMetaTable              = "users_meta"
	AuthOtpTable               = "authentication_otp"
	CountryTable               = "countries"
	CloneTable                 = "clone_app"
	PasswordResetTable         = "password_resets"
	OtpVerificationTable       = "otp_verification"
)

func Db() *gorm.DB {
	return _database
}

func OAuthClientTokenDb() *gorm.DB {
	return _database.Table(OAuthClientTable)
}
func OAuthAccessTokenDb() *gorm.DB {
	return _database.Table(OAuthAccessTokenTable)
}

func OAuthSessionsDb() *gorm.DB {
	return _database.Table(OAuthSessionsTable)
}

func OathAccessTokenScopesDb() *gorm.DB {
	return _database.Table(OathAccessTokenScopesTable)
}

func UserDb() *gorm.DB {
	return _database.Table(UserTable)
}

func UserSocialProfileDb() *gorm.DB {
	return _database.Table(UserSocialProfile)
}
func UserMeta() *gorm.DB {
	return _database.Table(UserMetaTable)
}

func DriverMetaDb() *gorm.DB {
	return _database.Table(DriverMetaTable)
}
func AuthOtpDb() *gorm.DB {
	return _database.Table(AuthOtpTable)
}
func CountryDb() *gorm.DB {
	return _database.Table(CountryTable)
}
func CloneApp() *gorm.DB {
	return _database.Table(CloneTable)
}

func PasswordReset() *gorm.DB {
	return _database.Table(PasswordResetTable)
}
func OtpVerification() *gorm.DB {
	return _database.Table(OtpVerificationTable)
}
func Ping() error {
	return _database.DB().Ping()
}
