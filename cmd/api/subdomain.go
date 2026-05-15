package main

import (
	"errors"
	"regexp"
	"strings"
)

var (
	errInvalidSubdomain  = errors.New("subdomain must be 3 to 63 characters and contain only lowercase letters, numbers, and hyphens")
	errReservedSubdomain = errors.New("subdomain is reserved")
	subdomainPattern     = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{1,61}[a-z0-9])$`)
	reservedSubdomains   = map[string]struct{}{
		"admin":      {},
		"api":        {},
		"app":        {},
		"assets":     {},
		"cdn":        {},
		"dashboard":  {},
		"eazymarket": {},
		"help":       {},
		"mail":       {},
		"root":       {},
		"status":     {},
		"store":      {},
		"support":    {},
		"www":        {},
	}
)

func normalizeSubdomain(subdomain string) (string, error) {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	if !subdomainPattern.MatchString(subdomain) {
		return "", errInvalidSubdomain
	}

	if _, reserved := reservedSubdomains[subdomain]; reserved {
		return "", errReservedSubdomain
	}

	return subdomain, nil
}
