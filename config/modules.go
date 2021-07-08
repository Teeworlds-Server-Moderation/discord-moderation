package config

import (
	"fmt"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/parsers"
	"github.com/jxsl13/simple-configo/unparsers"
)

type moduleConfig struct {
	enabledDiscordLog   bool
	enabledVPNDetection bool
}

func (m *moduleConfig) PostParse() error {
	return nil
}

func (m *moduleConfig) Close() error {
	return nil
}

func (m *moduleConfig) Name() string {
	return "modules"
}

func (m *moduleConfig) Options() configo.Options {
	return configo.Options{
		{
			Key:             "ENABLE_DISCORD_LOGGING",
			Description:     "Whether to enable logging to discord channels",
			DefaultValue:    "true",
			ParseFunction:   parsers.Bool(&m.enabledDiscordLog),
			UnparseFunction: unparsers.Bool(&m.enabledDiscordLog),
		},
		{
			Key:             "ENABLE_VPN_DETECTION",
			Description:     "Whether to enable logging to discord channels",
			DefaultValue:    "false",
			ParseFunction:   parsers.Bool(&m.enabledVPNDetection),
			UnparseFunction: unparsers.Bool(&m.enabledVPNDetection),
		},
	}
}

func (m *moduleConfig) ErrIfDiscordLoggingDisabled() error {
	if !m.enabledDiscordLog {
		return fmt.Errorf("the discord logging module is enabled")
	}
	return nil
}

func (m *moduleConfig) ErrIfVPNDetectionDisabled() error {
	if !m.enabledVPNDetection {
		return fmt.Errorf("the vpn detection module is disabled")
	}
	return nil
}
