package main

import (
	"github.com/hashicorp/go-plugin"

	"yaml2hcl2/internal/hcl"
	"yaml2hcl2/shared"
)

// HCLConverter est l'implémentation du plugin pour HCL
type HCLConverter struct{}

func (HCLConverter) Convert(input []byte) ([]byte, error) {
	return hcl.Evaluate(input)
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"converter": &shared.ConverterPlugin{Impl: &HCLConverter{}},
		},
	})
}

