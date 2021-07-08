package config

import (
	"os"

	configo "github.com/jxsl13/simple-configo"
)

type Config interface {
	configo.Config
	PostParse() error
	Close() error
}

func convertConfig(cs []Config) []configo.Config {
	c := make([]configo.Config, 0, len(cs))
	for _, v := range cs {
		c = append(c, v)
	}
	return c
}

func parse(cs ...Config) error {
	err := configo.ParseEnvFileOrEnv("./.env", convertConfig(cs)...)
	if err != nil {
		err = configo.ParseEnvFileOrEnv(envFileKey, convertConfig(cs)...)
		if err != nil {
			return err
		}
	} else {
		// in case we did use the ./.env file, set the fenvironment variable to that
		// value
		os.Setenv(envFileKey, "./.env")
	}
	for _, c := range cs {
		err = c.PostParse()
		if err != nil {
			return err
		}
	}
	return nil
}

func unparse(cs ...Config) error {
	err := configo.UnparseEnvFile(envFileKey, convertConfig(cs)...)
	if err != nil {
		return err
	}
	for _, c := range cs {
		err = c.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
