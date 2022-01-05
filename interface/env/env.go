package env

type Env interface {
	Init()
	LoadEnvFile(string)
}
