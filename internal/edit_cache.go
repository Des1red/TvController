package internal

import (
	"fmt"
	"sort"
	"strings"
	"tvctrl/internal/cache"
	"tvctrl/logger"
)

func storeInCache(cfg *Config, update cache.Device) {
	if !cfg.UseCache || cfg.SelectCache != -1 {
		return
	}
	store, _ := cache.Load()
	dev := store[cfg.TIP]

	alreadyStored := dev.ControlURL != ""

	if !cfg.AutoCache && !alreadyStored {
		if !confirm("Store this AVTransport endpoint in cache?") {
			return
		}
	} else {
		logger.Notify("=============== CACHE ===============")
		logger.Notify("%s already stored in cache", dev.ControlURL)
		logger.Notify("=============== CACHE ===============")
	}

	// --- merge ControlURL ---
	if dev.ControlURL == "" && update.ControlURL != "" {
		dev.ControlURL = update.ControlURL
	}

	// --- merge Vendor ---
	if dev.Vendor == "" && update.Vendor != "" {
		dev.Vendor = update.Vendor
	}

	// --- merge Identity ---
	if update.Identity != nil {
		dev.Identity = update.Identity
	}

	// --- merge Actions ---
	if update.Actions != nil {
		if dev.Actions == nil {
			dev.Actions = map[string]bool{}
		}
		for k, v := range update.Actions {
			dev.Actions[k] = v
		}
	}

	// --- merge Media ---
	if update.Media != nil {
		dev.Media = update.Media
	}

	store[cfg.TIP] = dev
	_ = cache.Save(store)
}

func LoadCachedTV(cfg *Config) {
	ip, dev, ok := selectFromCache(cfg.SelectCache)
	if !ok {
		logger.Fatal("Invalid cache index: %d", cfg.SelectCache)
	}

	cfg.TIP = ip
	cfg.TVVendor = dev.Vendor
	cfg._CachedControlURL = dev.ControlURL

	logger.Success(
		"Using cached device [%d]: %s",
		cfg.SelectCache,
		dev.ControlURL,
	)
}

func selectFromCache(index int) (string, cache.Device, bool) {
	store, _ := cache.Load()
	keys := sortedCache(store)

	if index < 0 || index >= len(keys) {
		return "", cache.Device{}, false
	}

	ip := keys[index]
	return ip, store[ip], true
}

func sortedCache(store cache.Store) []string {
	keys := make([]string, 0, len(store))
	for ip := range store {
		keys = append(keys, ip)
	}
	sort.Strings(keys)
	return keys
}

func HandleCacheCommands(cfg Config) bool {
	if cfg.ListCache {
		handleListCache()
		return true
	}

	if cfg.ForgetCache != "" {
		handleForgetCache(cfg)
		return true
	}

	return false
}

func handleForgetCache(cfg Config) {
	store, err := cache.Load()
	if err != nil {
		logger.Fatal("%v", err)
	}

	if len(store) == 0 {
		logger.Notify("Cache is empty.")
		return
	}

	switch cfg.ForgetCache {

	case "all":
		if !confirm("Delete ALL cached devices?") {
			return
		}
		_ = cache.Save(cache.Store{})
		logger.Success("Cache cleared.")
		return

	case "interactive":
		logger.Info("\nCached devices:")
		for ip, dev := range store {
			fmt.Printf(" %s → %s\n", ip, dev.ControlURL)
		}

		logger.Notify("\nEnter IP to delete: ")
		var ip string
		fmt.Scanln(&ip)

		if _, ok := store[ip]; !ok {
			logger.Notify("IP not found in cache.")
			return
		}

		if !confirm("Delete this entry?") {
			return
		}

		delete(store, ip)
		_ = cache.Save(store)
		logger.Success("Deleted %s", ip)
		return

	default:
		if _, ok := store[cfg.ForgetCache]; !ok {
			logger.Notify("IP not found in cache.")
			return
		}

		if !confirm("Delete cached entry for " + cfg.ForgetCache + "?") {
			return
		}

		delete(store, cfg.ForgetCache)
		_ = cache.Save(store)
		logger.Success("Deleted %s", cfg.ForgetCache)
	}
}

func handleListCache() {
	store, err := cache.Load()
	if err != nil {
		logger.Fatal("Error: %v", err)
	}

	if len(store) == 0 {
		logger.Info("Cache is empty.")
		return
	}

	logger.Info("\n\nCached AVTransport devices:\n\n")
	fmt.Printf(" %-15s %-10s %s\n", "IP", "Vendor", "ControlURL")
	fmt.Println(strings.Repeat("-", 60))

	keys := sortedCache(store)

	for i, ip := range keys {
		dev := store[ip]
		fmt.Printf("[%d] %-15s %-10s %s\n", i, ip, dev.Vendor, dev.ControlURL)
	}

}
