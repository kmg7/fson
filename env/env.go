package env

import "os"

func GetModeDebug() bool {
	return os.Getenv("FSON_DEBUG") != ""
}

func SetModeDebug() error {
	return os.Setenv("FSON_DEBUG", "A")
}
