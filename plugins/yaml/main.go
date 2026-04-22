package main

import (
	"github.com/hashicorp/go-plugin"

	"yaml2hcl2/internal/yaml"
	"yaml2hcl2/shared"
)

// YAMLConverter est l'implémentation du plugin pour YAML
type YAMLConverter struct{}

func (YAMLConverter) Convert(input []byte) ([]byte, error) {
	// Utilise la logique existante de conversion
	return yaml.ToHCL2(input)
}

func main() {
	// Sert le plugin via go-plugin
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"converter": &shared.ConverterPlugin{Impl: &YAMLConverter{}},
		},
	})
}
