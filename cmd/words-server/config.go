package main

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Dictionaries []string `yaml:"dictionaries"`
	MinWordLen   int      `yaml:"min_word_len"`
	MaxWordLen   int      `yaml:"max_word_len"`
	Port         int      `yaml:"port"` // meh
}

func (config *Config) ReadConfigurationFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return config.readConfiguration(f)
}

func (config *Config) readConfiguration(r io.Reader) error {
	var err error

	decoder := yaml.NewDecoder(r)
	decoder.SetStrict(true)

	if err = decoder.Decode(config); err != nil {
		return err
	}

	if err = config.validate(); err != nil {
		return err
	}
	return config.instantiate()
}

func (config *Config) validate() error {
	if len(conf.Dictionaries) == 0 {
		return errors.New("No dictionaries configured")
	}
	return nil
}

func (config *Config) instantiate() error {
	if config.MinWordLen == 0 {
		config.MinWordLen = 4
	}
	if config.MaxWordLen == 0 {
		config.MaxWordLen = 9
	}
	if config.Port == 0 {
		config.Port = 8080
	}
	return nil
}
