package converter

import (
	"fmt"
	"sort"

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

		ctyVal, err := convertToCTY(value)
		if err != nil {
			comment := fmt.Sprintf("\n# Erreur de conversion pour la clé '%s': %v\n", key, err)
			body.AppendUnstructuredTokens(hclwrite.TokensForIdentifier(comment))
			continue
		}

		body.SetAttributeValue(key, ctyVal)
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
