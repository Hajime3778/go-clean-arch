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

// Init: 環境の初期化
func (e *Env) Init() {
	// 本番環境 or 検証環境
	env := os.Getenv("ENVIRONMENT")
	if env == PRODUCTION || env == STAGING {
		return
	}

	// ローカル環境
	e.LoadEnvFile(".env")
}

// Load: 環境変数の読み込み
func (e *Env) LoadEnvFile(path string) {
	// ローカル環境
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Can't Load Environment: %v", err)
	}
}
