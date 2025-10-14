package internal

import "os"

var (
	dB_URL string
)

func LoadEnv() {
	dB_URL = os.Getenv("DB_URL")
}
