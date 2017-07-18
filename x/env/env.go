package env

import (
	"os"
	"strconv"
	"time"
)

// String inspects the environment variables specified in envvars.
// If all of these environment variables are empty, it returns defaultValue.
func String(defaultValue string, envvars ...string) string {
	for _, envvar := range envvars {
		if s := os.Getenv(envvar); s != "" {
			return s
		}
	}
	return defaultValue
}

// Int inspects the environment variables specified in envvars.
// If all of these environment variables are empty, it returns defaultValue.
func Int(defaultValue int, envvars ...string) int {
	for _, envvar := range envvars {
		if s := os.Getenv(envvar); s != "" {
			if i, err := strconv.Atoi(s); err == nil {
				return i
			}
		}
	}
	return defaultValue
}

// Int64 inspects the environment variables specified in envvars.
// If all of these environment variables are empty, it returns defaultValue.
func Int64(defaultValue int64, envvars ...string) int64 {
	for _, envvar := range envvars {
		if s := os.Getenv(envvar); s != "" {
			if i, err := strconv.ParseInt(s, 10, 64); err == nil {
				return i
			}
		}
	}
	return defaultValue
}

// Bool inspects the environment variables specified in envvars.
// If all of these environment variables are empty, it returns defaultValue.
func Bool(defaultValue bool, envvars ...string) bool {
	for _, envvar := range envvars {
		if s := os.Getenv(envvar); s != "" {
			if flag, err := strconv.ParseBool(s); err == nil {
				return flag
			}
		}
	}
	return defaultValue
}

// Duration inspects the environment variables specified in envvars.
// If all of these environment variables are empty, it returns defaultValue.
func Duration(defaultValue time.Duration, envvars ...string) time.Duration {
	for _, envvar := range envvars {
		if s := os.Getenv(envvar); s != "" {
			if d, err := time.ParseDuration(s); err == nil {
				return d
			}
		}
	}
	return defaultValue
}
