package config

import (
	"errors"
	"fmt"
	"strconv"
)

const defaultHTTPAddr = "127.0.0.1:8080"

// ResolveHTTPAddr keeps the safe loopback default for ordinary local hosts
// while honouring platforms that expose the container's public listen port via
// PORT (for example Railway and Cloud Run). An explicitly configured address
// always wins.
//
// Hosts do not necessarily expand variable references inside another
// environment variable, so setting MCP_HTTP_ADDR=$PORT is not portable. When
// PORT is present we consume it directly and bind on all interfaces.
func ResolveHTTPAddr(configured string, explicitlySet bool, lookupEnv func(string) (string, bool)) (string, error) {
	if explicitlySet {
		return configured, nil
	}
	if configured == "" {
		configured = defaultHTTPAddr
	}
	if lookupEnv == nil {
		return configured, nil
	}

	if port, ok := lookupEnv("PORT"); ok && port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return "", fmt.Errorf("invalid PORT %q: expected an integer from 1 to 65535", port)
		}
		return ":" + port, nil
	}

	// Cloud Run promises PORT. If K_SERVICE is present but PORT is not, fail
	// loudly instead of silently binding only to loopback and looking healthy
	// while remaining unreachable.
	if _, cloudRun := lookupEnv("K_SERVICE"); cloudRun {
		return "", errors.New("Cloud Run K_SERVICE is set but PORT is missing") //nolint:staticcheck // Cloud Run and its env names are proper nouns.
	}

	return configured, nil
}
