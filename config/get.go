package config

import "github.com/xackery/eqcp/tlog"

func Get(configConst int) string {
	switch configConst {
	case IsDiscoveredOnly:
		return config[IsDiscoveredOnly]
	default:
		tlog.Errorf("config.Get: unknown configConst: %v", configConst)
		return ""
	}
}

func IsEnabled(configConst int) bool {
	return Get(configConst) == "1"
}
