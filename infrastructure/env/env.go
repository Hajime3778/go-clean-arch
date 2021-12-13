package env

import (
	"log"
	"os"

	"github.com/Hajime3778/go-clean-arch/interface/env"
	"github.com/joho/godotenv"
)

const (
	PRODUCTION = "PRODUCTION"
	STAGING    = "STAGING"
)

type Env struct{}

func NewEnv() env.Env {
	return &Env{}
}

// Load: 環境変数の読み込み
func (e *Env) Load() {
	// 本番環境 or 検証環境
	env := os.Getenv("ENVIRONMENT")
	if env == PRODUCTION || env == STAGING {
		return
	}

	// ローカル環境
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Can't Load Environment: %v", err)
	}
}
