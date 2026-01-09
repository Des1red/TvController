package avtransport

import (
	"errors"
	"fmt"
	"time"
)

var probePorts = []string{
	"9197", // Samsung AVTransport
	"7678", // Samsung AllShare
	"8187",
	"9119",
	"8080",
}

var defaultList = []string{
	"/dmr/upnp/control/AVTransport1",
	"/upnp/control/AVTransport",
	"/MediaRenderer/AVTransport/Control",
	"/AVTransport/control",
}

var bigList = []string{
	// --- Common / standard ---
	"/upnp/control/AVTransport",
	"/AVTransport/control",
	"/MediaRenderer/AVTransport/Control",
	"/dmr/upnp/control/AVTransport1",

	// --- Samsung ---
	"/smp_7_/AVTransport",
	"/smp_9_/AVTransport",
	"/upnp/control/AVTransport1",

	// --- Sony / Bravia ---
	"/sony/AVTransport",
	"/upnp/control/AVTransport/1",

	// --- LG webOS ---
	"/upnp/control/AVTransport",
	"/upnp/control/avtransport",

	// --- Generic DMR variants ---
	"/renderer/control/AVTransport",
	"/device/AVTransport/control",
	"/control/AVTransport",

	// --- Case / path quirks ---
	"/AVTransport/Control",
	"/avtransport/control",
}

func Probe(ip string, timeout time.Duration, list bool) (*Target, error) {
	var endpoints []string
	if list {
		endpoints = bigList
	} else {
		endpoints = defaultList
	}
	deadline := time.Now().Add(timeout)

	for _, port := range probePorts {
		for _, path := range endpoints {
			if time.Now().After(deadline) {
				return nil, errors.New("AVTransport probe timed out")
			}

			controlURL := fmt.Sprintf("http://%s:%s%s", ip, port, path)

			ok := soapProbe(controlURL, "")
			if ok {
				return &Target{
					ControlURL: controlURL,
				}, nil
			}
		}
	}

	return nil, errors.New("no AVTransport endpoint found")
}
