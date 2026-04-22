package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-plugin"

	"yaml2hcl2/shared"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <fichier_config>\n", os.Args[0])
	}

	inputFile := os.Args[1]
	ext := strings.ToLower(filepath.Ext(inputFile))

	// Déterminer le plugin à utiliser selon l'extension
	var pluginName string
	var outputFile string

	if ext == ".yaml" || ext == ".yml" {
		pluginName = "yaml"
		outputFile = "output.hcl"
	} else if ext == ".hcl" {
		pluginName = "hcl"
		outputFile = "output.json"
	} else {
		log.Fatalf("Extension non supportée: %s (utilisez .yaml ou .hcl)", ext)
	}

	pluginPath := fmt.Sprintf("./plugins/%s/plugin", pluginName)

	// Configuration du client du plugin
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
		},
	})
	defer client.Kill()

	// Connexion RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatalf("Erreur de client RPC: %s", err)
	}

	// Demande de l'implémentation du plugin "converter"
	raw, err := rpcClient.Dispense("converter")
	if err != nil {
		log.Fatalf("Erreur de dispense de plugin: %s", err)
	}

	converter := raw.(shared.Converter)

	// 1. Lecture du fichier
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier : %v", err)
	}

	// 2. Conversion/Évaluation via le plugin
	resultData, err := converter.Convert(data)
	if err != nil {
		log.Fatalf("Erreur lors de l'exécution du plugin : %v", err)
	}

	// 3. Écriture du fichier de sortie
	err = os.WriteFile(outputFile, resultData, 0644)
	if err != nil {
		log.Fatalf("Erreur lors de l'écriture du fichier %s : %v", outputFile, err)
	}

	fmt.Printf("Opération terminée ! Le résultat est disponible dans le fichier '%s'\n", outputFile)
}
