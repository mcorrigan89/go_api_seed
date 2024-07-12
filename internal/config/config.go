package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN     string
		Logging bool
	}
	Cors struct {
		TrustedOrigins []string
	}
	OAuth struct {
		Google struct {
			ClientID     string
			ClientSecret string
		}
	}
	Logging struct {
		NewRelicAppName    string
		NewRelicLicenseKey string
	}
}

func LoadConfig(cfg *Config) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Load ENV
	env := os.Getenv("ENV")
	if env == "" {
		cfg.Env = "local"
	} else {
		cfg.Env = env
	}

	// Load PORT
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("PORT not available in .env")
	}

	cfg.Port = port

	// Load DATABASE_URL
	postgres_url := os.Getenv("POSTGRES_URL")
	if postgres_url == "" {
		log.Fatalf("POSTGRES_URL not available in .env")
	}

	cfg.DB.DSN = postgres_url

	cfg.Cors.TrustedOrigins = []string{"http://localhost:3000"}

	// Load OAuth
	google_client_id := os.Getenv("GOOGLE_CLIENT_ID")
	if google_client_id == "" {
		log.Fatalf("GOOGLE_CLIENT_ID not available in .env")
	}

	google_client_secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if google_client_secret == "" {
		log.Fatalf("GOOGLE_CLIENT_SECRET not available in .env")
	}

	cfg.OAuth.Google.ClientID = google_client_id
	cfg.OAuth.Google.ClientSecret = google_client_secret

	new_relic_app_name := os.Getenv("NEW_RELIC_APP_NAME")
	if new_relic_app_name == "" {
		log.Fatalf("NEW_RELIC_APP_NAME not available in .env")
	}
	cfg.Logging.NewRelicAppName = new_relic_app_name

	new_relic_license_key := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if new_relic_license_key == "" {
		log.Fatalf("NEW_RELIC_LICENSE_KEY not available in .env")
	}
	cfg.Logging.NewRelicLicenseKey = new_relic_license_key

}
