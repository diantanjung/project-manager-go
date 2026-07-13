package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	NodeEnv             string
	Port                int
	DatabaseURL         string
	FrontendURL         string
	JWTSecret           string
	JWTExpiresIn        time.Duration
	JWTRefreshSecret    string
	JWTRefreshExpiresIn time.Duration
	UploadDir           string
}

func Load() (Config, error) {
	if err := loadEnvFile(".env"); err != nil {
		return Config{}, err
	}

	port, err := strconv.Atoi(getenv("PORT", "3000"))
	if err != nil {
		return Config{}, fmt.Errorf("parsing PORT: %w", err)
	}

	accessTTL, err := parseDuration(getenv("JWT_EXPIRES_IN", "15m"))
	if err != nil {
		return Config{}, fmt.Errorf("parsing JWT_EXPIRES_IN: %w", err)
	}

	refreshTTL, err := parseDuration(getenv("JWT_REFRESH_EXPIRES_IN", "7d"))
	if err != nil {
		return Config{}, fmt.Errorf("parsing JWT_REFRESH_EXPIRES_IN: %w", err)
	}

	return Config{
		NodeEnv:             getenv("NODE_ENV", "development"),
		Port:                port,
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		FrontendURL:         getenv("FRONTEND_URL", "http://localhost:5173"),
		JWTSecret:           getenv("JWT_SECRET", "supersecret"),
		JWTExpiresIn:        accessTTL,
		JWTRefreshSecret:    getenv("JWT_REFRESH_SECRET", "superrefreshsecret"),
		JWTRefreshExpiresIn: refreshTTL,
		UploadDir:           getenv("UPLOAD_DIR", "uploads"),
	}, nil
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func parseDuration(value string) (time.Duration, error) {
	if len(value) < 2 {
		return 0, fmt.Errorf("invalid duration %q", value)
	}

	unit := value[len(value)-1]
	number, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0, err
	}

	switch unit {
	case 'd':
		return time.Duration(number) * 24 * time.Hour, nil
	case 'h':
		return time.Duration(number) * time.Hour, nil
	case 'm':
		return time.Duration(number) * time.Minute, nil
	case 's':
		return time.Duration(number) * time.Second, nil
	default:
		return 0, fmt.Errorf("unsupported duration unit %q", unit)
	}
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("opening env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" || os.Getenv(key) != "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("setting env %s: %w", key, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}
	return nil
}
