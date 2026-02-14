//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Service struct {
	URL    string `json:"url,omitempty" toml:"url"`
	Domain string `json:"domain,omitempty" toml:"domain"`
	Host   string `json:"host,omitempty" toml:"host"`
}

type Environment struct {
	Services map[string]Service `json:"service,omitempty" toml:"service"`
}

type Region struct {
	Environments map[string]Environment `json:"env,omitempty" toml:"env"`
	Services     map[string]Service     `json:"service,omitempty" toml:"service"`
}

type World struct {
	Regions map[string]Region `json:"region" toml:"region"`
}

func main() {
	var world World
	if _, err := toml.DecodeFile("hsdp.toml", &world); err != nil {
		fmt.Printf("Error decoding TOML: %v\n", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(world, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Add trailing newline to match usual JSON file formatting
	data = append(data, '\n')

	if err := os.WriteFile("hsdp.json", data, 0644); err != nil {
		fmt.Printf("Error writing JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("hsdp.json updated successfully")
}
