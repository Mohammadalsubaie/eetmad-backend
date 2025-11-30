package database

import (
"fmt"
"log"
"os"

"github.com/joho/godotenv"
"gorm.io/driver/postgres"
"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
err := godotenv.Load()
if err != nil {
log.Println("No .env file found, using system environment")
}

dsn := os.Getenv("DATABASE_URL")
if dsn == "" {
dsn = "postgresql://eetmad_user:wPtzFX5Z7s2aNlQkxn1%2FerpSJtisLG9AopNMwMju2bc=@localhost:5432/eetmad_prod?sslmode=disable"
}

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
log.Fatal("[error] failed to initialize database, got error ", err)
}

fmt.Println("DB connected")
DB = db
}
