package hcl

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// Evaluate convertit et évalue un fichier HCL vers un JSON
func Evaluate(input []byte) ([]byte, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(input, "input.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("Erreur de parsing HCL: %s", diags.Error())
	}

	// Configuration du contexte d'évaluation avec des variables et fonctions
	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"env": buildEnvVars(),
		},
		Functions: map[string]function.Function{
			"upper": stdlib.UpperFunc,
			"lower": stdlib.LowerFunc,
			"max":   stdlib.MaxFunc,
			"min":   stdlib.MinFunc,
		},
	}

	// On extrait uniquement les attributs à la racine
	attrs, diags := file.Body.JustAttributes()
	if diags.HasErrors() {
		return nil, fmt.Errorf("Erreur d'extraction des attributs HCL: %s", diags.Error())
	}

	// Évaluation des attributs et conversion en JSON
	result := make(map[string]interface{})
	for name, attr := range attrs {
		val, valDiags := attr.Expr.Value(ctx)
		if valDiags.HasErrors() {
			return nil, fmt.Errorf("Erreur d'évaluation de l'attribut %q: %s", name, valDiags.Error())
		}

		// Convertit la valeur évaluée en JSON
		jsonBytes, err := ctyjson.Marshal(val, val.Type())
		if err != nil {
			return nil, fmt.Errorf("Erreur de conversion JSON pour %q: %v", name, err)
		}

		var goVal interface{}
		if err := json.Unmarshal(jsonBytes, &goVal); err != nil {
			return nil, err
		}
		result[name] = goVal
	}

	// Retourne le JSON formaté (la configuration HCL évaluée)
	return json.MarshalIndent(result, "", "  ")
}

// buildEnvVars crée un objet cty contenant les variables d'environnement
func buildEnvVars() cty.Value {
	envMap := make(map[string]cty.Value)
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = cty.StringVal(parts[1])
		}
	}

	if len(envMap) == 0 {
		return cty.EmptyObjectVal
	}
	return cty.ObjectVal(envMap)
}
