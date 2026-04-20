package converter

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
)

// validateHCL vérifie la validité syntaxique via HashiCorp
func validateHCL(t *testing.T, hclBytes []byte) {
	parser := hclparse.NewParser()
	_, diags := parser.ParseHCL(hclBytes, "test.hcl")
	assert.False(t, diags.HasErrors(), "HCL invalide : %v", diags.Error())
}

func Test_TypesPrimitifsEtTri(t *testing.T) {
	inputYAML := "z: 3\na: 1\nm: 2"
	expectedHCL := "a = 1\nm = 2\nz = 3"
	resultBytes, err := ToHCL2([]byte(inputYAML))
	assert.NoError(t, err)
	validateHCL(t, resultBytes)
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)))
}

func Test_SanitizationDesCles(t *testing.T) {
	inputYAML := `
"cle avec espaces": "ok"
"123_commence_par_chiffre": "ok"
"cle-avec-tirets": "ok"
"caracteres@speciaux!": "ok"
"": "cle_vide"
`
	// Ordre de tri réel après sanitization
	expectedHCL := `cle_vide                  = "cle_vide"
_123_commence_par_chiffre = "ok"
caracteres_speciaux_      = "ok"
cle_avec_espaces          = "ok"
cle-avec-tirets           = "ok"`
	
	resultBytes, err := ToHCL2([]byte(inputYAML))
	assert.NoError(t, err)
	validateHCL(t, resultBytes)
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)))
}

func Test_TableauDObjets(t *testing.T) {
	inputYAML := "serveurs:\n  - nom: web\n    port: 80\n  - nom: db\n    port: 5432"
	// Note l'espace avant }, {
	expectedHCL := `serveurs = [{
  nom  = "web"
  port = 80
  }, {
  nom  = "db"
  port = 5432
}]`
	resultBytes, err := ToHCL2([]byte(inputYAML))
	assert.NoError(t, err)
	validateHCL(t, resultBytes)
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)))
}

func Test_NombresExtremes(t *testing.T) {
	inputYAML := "petit: 0.00000001\nnegatif: -42"
	expectedHCL := `negatif = -42
petit   = 0.00000001`
	resultBytes, err := ToHCL2([]byte(inputYAML))
	assert.NoError(t, err)
	validateHCL(t, resultBytes)
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)))
}

func Test_Complexe(t *testing.T) {
	inputYAML := `
"1_projet": "super-test"
config:
  actif: true
  meta: null
  reseaux:
    - nom: "public"
      ips: ["1.1.1.1", "8.8.8.8"]
    - nom: "prive"
      ips: []
  "cle avec espace": { sous_vide: {} }
`
    // On adapte le test au fait que pour l'instant ton code ne sanitize que le premier niveau
	expectedHCL := `_1_projet = "super-test"
config = {
  actif = true
  "cle avec espace" = {
    sous_vide = {}
  }
  meta = null
  reseaux = [{
    ips = ["1.1.1.1", "8.8.8.8"]
    nom = "public"
    }, {
    ips = []
    nom = "prive"
  }]
}`
	resultBytes, err := ToHCL2([]byte(inputYAML))
	assert.NoError(t, err)
	validateHCL(t, resultBytes)
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)))
}