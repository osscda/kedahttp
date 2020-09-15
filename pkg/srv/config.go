package srv

import "os"

func EnvOr(envName, otherwise string) string {
	fromEnv := os.Getenv(envName)
	if fromEnv == "" {
		return otherwise
	}
	return fromEnv
}
