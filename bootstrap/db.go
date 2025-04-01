// bootstrap/database.go
package bootstrap

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func InitDB(env *Env) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		(*env).DBHost,
		env.DBPort,
		env.DBUser,
		env.DBName,
		env.DBPassword,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database", err)
	}

	return db
}
