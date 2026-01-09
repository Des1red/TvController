package internal

import (
	"fmt"
	"time"
	"tvctrl/internal/avtransport"
	"tvctrl/internal/cache"
	"tvctrl/internal/identity"
	"tvctrl/internal/ssdp"
	"tvctrl/logger"
)

func trySSDP(cfg *Config) bool {
	logger.Notify("Running SSDP discovery scan")
	devices, _ := ssdp.ListenNotify(3 * time.Second)

	if len(devices) == 0 {
		logger.Notify("No NOTIFY devices from SSDP listen, trying SSDP discover")
		devices, _ = ssdp.Discover(3 * time.Second)
	}

	if len(devices) == 0 {
		return false
	}

	tv, err := ssdp.FetchAndDetect(devices[0].Location)
	if err != nil {
		return false
	}

	if tv.IP != "" {
		cfg.TIP = tv.IP
	}
	if tv.Port != "" {
		cfg.TPort = tv.Port
	}
	if tv.ControlURL != "" {
		cfg.TPath = tv.ControlURL
	}
	if tv.Vendor != "" {
		cfg.TVVendor = tv.Vendor
	}

	caps, err := avtransport.EnrichCapabilities(
		tv.AVTransportSCPD,
		tv.ConnectionManagerCtrl,
		avtransport.Target{
			ControlURL: cfg.ControlURL(),
		},
	)

	info, err := identity.Enrich(
		cfg.BaseUrl(),
		3*time.Second,
	)

	update := cache.Device{
		ControlURL: cfg.ControlURL(),
		Vendor:     tv.Vendor,
	}

	if err == nil {
		update.Identity = map[string]any{
			"friendly_name": info.FriendlyName,
			"manufacturer":  info.Manufacturer,
			"model_name":    info.ModelName,
			"model_number":  info.ModelNumber,
			"udn":           info.UDN,
			"presentation":  info.Presentation,
		}
	}

	if err == nil && caps != nil {
		update.Actions = caps.Actions
		update.Media = caps.Media
	} else {
		logger.Notify("Capability enrichment failed: %v", err)
	}

	storeInCache(cfg, update)

	return true
}

func probeAVTransport(cfg *Config) (bool, error) {
	if cfg.TIP == "" {
		return false, fmt.Errorf("probe requires -Tip")
	}

	logger.Notify("Probing AVTransport directly : %s", cfg.TIP)

	target, err := avtransport.Probe(cfg.TIP, 8*time.Second, cfg.DeepSearch)
	if err != nil {
		return false, err
	}

	observedActions := avtransport.ValidateActions(*target)

	// update cfg so playback can continue
	cfg._CachedControlURL = target.ControlURL
	info, err := identity.Enrich(
		"http://"+cfg.TIP,
		3*time.Second,
	)
	update := cache.Device{
		ControlURL: target.ControlURL,
		Vendor:     cfg.TVVendor,
	}

	if err == nil {
		update.Identity = map[string]any{
			"friendly_name": info.FriendlyName,
			"manufacturer":  info.Manufacturer,
			"model_name":    info.ModelName,
			"model_number":  info.ModelNumber,
			"udn":           info.UDN,
			"presentation":  info.Presentation,
		}
	} else {
		logger.Notify("%v", err)
	}

	if len(observedActions) > 0 {
		update.Actions = observedActions
	}

	storeInCache(cfg, update)

	logger.Success("\n=== AVTransport Probe Summary ===")

	logger.Result(" IP        : %s", cfg.TIP)
	logger.Result(" ControlURL: %s", target.ControlURL)

	return true, nil
}
