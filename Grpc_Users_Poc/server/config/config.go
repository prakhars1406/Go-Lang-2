package config

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	Configuration Config
)

// Represents database server and credentials
type Config struct {
	ADDRESS string
}

// Read and parse the configuration file
func (c *Config) Read(filepath string) {
	fmt.Println("This is in config")
	if _, err := toml.DecodeFile(filepath, &c); err != nil {
		log.Fatal(err)
	}
	Configuration = *c
}
