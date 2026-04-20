package converter

import (
	"fmt"
	"sort"
    
	"regexp"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v3"
)

func ToHCL2(yamlBytes []byte) ([]byte, error) {

	var yamlData map[string]interface{}

	err := yaml.Unmarshal(yamlBytes, &yamlData)
	if err != nil {
		return nil, fmt.Errorf("Impossible de parser le YAML: %w", err)
	}

	//Création du fichier HCL ( vide au debut)
	hclFile := hclwrite.NewEmptyFile()
	body := hclFile.Body()

	var keys []string
	for k := range yamlData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// ensuite on parcourt le fichier yaml donné en entré
	for _, key := range keys {
		value := yamlData[key] // On récupère la valeur correspondante
		safeKey := sanitizeKey(key)

		ctyVal, err := convertToCTY(value)
		if err != nil {
			comment := fmt.Sprintf("\n# Erreur de conversion pour la clé '%s': %v\n", key, err)
			body.AppendUnstructuredTokens(hclwrite.TokensForIdentifier(comment))
			continue
		}

		body.SetAttributeValue(safeKey, ctyVal)
	}
	resultBytes := hclFile.Bytes() // On génère les octets
	
	parser := hclparse.NewParser()
	_, diags := parser.ParseHCL(resultBytes, "internal_validation.hcl")
	
	// Si notre propre code a généré du HCL invalide, on bloque tout 
	if diags.HasErrors() {
		return nil, fmt.Errorf("échec critique: le HCL généré est syntaxiquement invalide. Détails: %s", diags.Error())
	}

	return hclFile.Bytes(), nil

}

func convertToCTY(val interface{}) (cty.Value, error) {
	if val == nil {
		return cty.NullVal((cty.String)), nil
	}

	switch v := val.(type) {

	case string:
		return cty.StringVal(v), nil
	case int:
		return cty.NumberIntVal(int64(v)), nil
	case float64:
		return cty.NumberFloatVal(float64(v)), nil
	case bool:
		return cty.BoolVal(v), nil

	case []interface{}:
		var ctyValues []cty.Value
		for _, item := range v {
			parsedItem, err := convertToCTY(item)
			if err != nil {
				return cty.NilVal, err
			}
			ctyValues = append(ctyValues, parsedItem)
		}

		// On utilise TupleVal au lieu de ListVal car un YAML peut avoir
		// une liste avec des types mixtes

		if len(ctyValues) == 0 {
			return cty.EmptyTupleVal, nil
		}

		return cty.TupleVal(ctyValues), nil

	case map[string]interface{}:
		objValues := make(map[string]cty.Value)

		for key, item := range v {
			parsedItem, err := convertToCTY(item)

			if err != nil {
				return cty.NilVal, err
			}
			objValues[key] = parsedItem

		}
		if len(objValues) == 0 {
			return cty.EmptyObjectVal, nil
		}
		return cty.ObjectVal(objValues), nil

	default:
		return cty.NilVal, fmt.Errorf("type Go pas supporté: %T", v)

	}

}

// sanitizeKey nettoie une chaîne pour en faire un identifiant HCL valide.
func sanitizeKey(key string) string {
	if key == "" {
		return "cle_vide"
	}
	
	// Remplace tout ce qui n'est pas une lettre, un chiffre, un tiret ou un underscore par "_"
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	safeKey := re.ReplaceAllString(key, "_")

	// Si la clé commence par un chiffre, on ajoute un "_" au début
	if safeKey[0] >= '0' && safeKey[0] <= '9' {
		safeKey = "_" + safeKey
	}

	return safeKey
}