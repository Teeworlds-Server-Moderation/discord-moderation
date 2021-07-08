package config

import (
	"fmt"
	"log"
)

var (
	moduleCfg      *moduleConfig
	brokerCfg      *brokerConfig
	discordCfg     *discordConfig
	detectVPNCfg   *detectVPNConfig
	envFileKey              = "ENV_FILE"
	enabledModules []Config = make([]Config, 0)
)

func Broker() *brokerConfig {
	return brokerCfg
}

func Discord() *discordConfig {
	return discordCfg
}

func DetectVPN() *detectVPNConfig {
	return detectVPNCfg
}

func Modules() *moduleConfig {
	return moduleCfg
}

func EnabledModules() []string {
	names := make([]string, 0)
	for _, module := range enabledModules {
		names = append(names, module.Name())
	}
	return names
}

func init() {
	moduleCfg = &moduleConfig{}
	err := parse(moduleCfg)
	if err != nil {
		log.Fatalln(err)
	}

	brokerCfg = &brokerConfig{}
	enabledModules = append(enabledModules, brokerCfg)

	if moduleCfg.enabledDiscordLog {
		discordCfg = newDiscordConfig()
		enabledModules = append(enabledModules, discordCfg)
	}

	if moduleCfg.enabledVPNDetection {
		detectVPNCfg = &detectVPNConfig{}
		enabledModules = append(enabledModules, detectVPNCfg)
	}

	err = parse(enabledModules...)
	if err != nil {
		log.Fatalln(err)
	}
	enabledModules = append(enabledModules, moduleCfg)
}

// If you want to save any changed back to your config file, call this method
// at the end of oyur main function.
func Close() error {
	fmt.Println()
	log.Println("Saving configuration...")
	return unparse(enabledModules...)
}
