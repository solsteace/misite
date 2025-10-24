package internal

import (
	"os"
	"path"
)

var (
	dB_URL    string
	dB_SCHEMA string

	aLPINE_URL string
	hTMX_URL   string
)

func LoadEnv() {
	dB_URL = os.Getenv("DB_URL")
	dB_SCHEMA = os.Getenv("DB_SCHEMA")

	// Use offline scripts instead of fetching it from CDN
	aLPINE_URL = os.Getenv("ALPINEJS_URL")
	hTMX_URL = os.Getenv("HTMX_URL")
	if root := os.Getenv("LOCAL_SCRIPT_URL"); root != "" {
		aLPINE_URL = path.Join(root, "alpinejs")
		hTMX_URL = path.Join(root, "htmx")
	}
}
