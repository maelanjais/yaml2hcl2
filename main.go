package main

import (
	"fmt"
	"log"
	"os"
	

	"yaml2hcl2/internal/converter"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <fichier.yaml>\n", os.Args[0])
	}

	inputFile := os.Args[1]
	outputFile := "output.hcl"

	// 1. Lecture du fichier YAML
	yamlData, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier : %v", err)
	}

	// 2. Conversion et Auto-Validation (via ton package internal)
	hclData, err := converter.ToHCL2(yamlData)
	if err != nil {
		log.Fatalf(" Erreur lors de la conversion : %v", err)
	}

	// Si on arrive ici, c'est que l'auto-validation de converter.go a réussi !
	fmt.Println("Validation syntaxique HCL : Succès")

	// 3. Écriture du fichier de sortie
	err = os.WriteFile(outputFile, hclData, 0644)
	if err != nil {
		log.Fatalf("Erreur lors de l'écriture du fichier %s : %v", outputFile, err)
	}

	// 4. Confirmation finale claire et nette
	fmt.Printf("Conversion terminée ! Le résultat est disponible dans le fichier '%s'\n", outputFile)
}
