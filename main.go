package main

import (
	"fmt"
	"log"
	"os"
	
	"yaml2hcl2/internal/converter"
)

func main(){
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <fichier.yaml>\n", os.Args[0])
	}
	
	inputFile := os.Args[1]
	outputFile := "output.hcl"

	yamlData, err := os.ReadFile(inputFile)
	if err != nil{
		log.Fatalf("Erreur lors de la lecture du fichier : %v", err)
	}

	hclData, err := converter.ToHCL2(yamlData)
	if err != nil {
		log.Fatalf("Erreur lors de la conversion du fichier : %v", err)
	}

	// ecriture du fichier de sortie 
	err = os.WriteFile(outputFile, hclData, 0644)
	if err != nil {
		log.Fatalf("Erreur lors de l'écriture du fichier %s : %v", outputFile, err)
	}	

	fmt.Println((string(hclData)))


}