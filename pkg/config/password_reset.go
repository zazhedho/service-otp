package config

import (
	"strings"
	"time"

	"service-otp/utils"
)

type PasswordResetConfig struct {
	TTL         time.Duration
	Cooldown    time.Duration
	RateWindow  time.Duration
	RateLimit   int
	Secret      string
	URLTemplate string
}

func LoadPasswordResetConfig() PasswordResetConfig {
	ttl := time.Duration(utils.GetEnv("RESET_TTL_SECONDS", 900).(int)) * time.Second
	if v := strings.TrimSpace(utils.GetEnv("RESET_TTL", "").(string)); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			ttl = d
		}
	}

	cooldown := time.Duration(utils.GetEnv("RESET_COOLDOWN_SECONDS", 60).(int)) * time.Second
	if v := strings.TrimSpace(utils.GetEnv("RESET_COOLDOWN", "").(string)); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cooldown = d
		}
	}

	rateWindow := time.Duration(utils.GetEnv("RESET_RATE_WINDOW_SECONDS", int(ttl.Seconds())).(int)) * time.Second
	if v := strings.TrimSpace(utils.GetEnv("RESET_RATE_WINDOW", "").(string)); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			rateWindow = d
		}
	}

	rateLimit := utils.GetEnv("RESET_RATE_LIMIT", 5).(int)

	secret := strings.TrimSpace(utils.GetEnv("RESET_SECRET", "reset-secret").(string))
	urlTemplate := strings.TrimSpace(utils.GetEnv("RESET_URL_TEMPLATE", "").(string))
	if urlTemplate == "" {
		urlTemplate = strings.TrimSpace(utils.GetEnv("RESET_URL", "").(string))
	}

	return PasswordResetConfig{
		TTL:         ttl,
		Cooldown:    cooldown,
		RateWindow:  rateWindow,
		RateLimit:   rateLimit,
		Secret:      secret,
		URLTemplate: urlTemplate,
	}
}
