package utils

import "sync"

var (
	tokenBlacklist = make(map[string]struct{})
	blacklistMu    sync.RWMutex
)

func BlacklistToken(token string) {
	blacklistMu.Lock()
	defer blacklistMu.Unlock()
	tokenBlacklist[token] = struct{}{}
}

func IsTokenBlacklisted(token string) bool {
	blacklistMu.RLock()
	defer blacklistMu.RUnlock()
	_, exists := tokenBlacklist[token]
	return exists
}
