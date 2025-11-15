package registries

import "os"

func getEnv() string {
	env := os.Getenv("GO_ENV")
	if os.Getenv("GO_ENV") == "" {
		env = "development"
	}
	return env
}
