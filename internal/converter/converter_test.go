package converter

import (
	"strings"
	"testing"


	"github.com/stretchr/testify/assert"
)

// 1. Test des types primitifs (Texte, Nombre, Booléen)
func TestTypesPrimitifs(t *testing.T) {
	inputYAML := `
chaine: "texte"
entier: 42
decimal: 3.14
booleen: true
`
	expectedHCL := `booleen = true
chaine  = "texte"
decimal = 3.14
entier  = 42
`

	resultBytes, err := ToHCL2([]byte(inputYAML))

	assert.NoError(t, err, "Il ne devrait pas y avoir d'erreur de conversion")
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)))
}

// 2. Test des listes et des tuples
func TestCollections(t *testing.T) {
	inputYAML := `
liste_nombres: [1, 2, 3]
tuple_mixte: ["un", 2, true]
`
	expectedHCL := `liste_nombres = [1, 2, 3]
tuple_mixte   = ["un", 2, true]
`

	resultBytes, err := ToHCL2([]byte(inputYAML))

	assert.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)))
}

// 3. Test des dictionnaires (Objets HCL)
func TestObjetsImbriques(t *testing.T) {
	inputYAML := `
serveur:
  port: 8080
  actif: true
`
	expectedHCL := `serveur = {
  actif = true
  port  = 8080
}
`

	resultBytes, err := ToHCL2([]byte(inputYAML))

	assert.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)))
}

// 4. Test des cas limites (Listes et objets vides)
func TestStructuresVides(t *testing.T) {
	inputYAML := `
liste_vide: []
objet_vide: {}
`
	expectedHCL := `liste_vide = []
objet_vide = {}
`

	resultBytes, err := ToHCL2([]byte(inputYAML))

	assert.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)))
}

// 5. Test d'une configuration complexe complète
func TestCombinaisonComplexe(t *testing.T) {
	inputYAML := `
app:
  nom: "mon-api"
  reseaux:
    - "interne"
    - "externe"
`
	expectedHCL := `app = {
  nom     = "mon-api"
  reseaux = ["interne", "externe"]
}
`

	resultBytes, err := ToHCL2([]byte(inputYAML))

	assert.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)))
}

// 6. Test de la gestion des erreurs (Un mauvais fichier YAML)
func TestErreurYAMLInvalide(t *testing.T) {
	// On simule un YAML cassé (fermeture d'une liste avec une accolade)
	inputYAML := `mauvais_yaml: [1, 2}`

	resultBytes, err := ToHCL2([]byte(inputYAML))

	// Cette fois-ci, on ASSERT qu'il y a bien une erreur
	assert.Error(t, err, "Le programme aurait dû planter car le YAML est invalide")
	// On vérifie que le résultat est bien vide en cas d'erreur
	assert.Nil(t, resultBytes)
}

// 7. Le "Boss Final" : Un YAML massif et extrêmement complexe
func TestYAMLComplexe(t *testing.T) {
	// Un YAML digne d'une vraie configuration d'infrastructure cloud complexe
	inputYAML := `
projet: "mega-infrastructure"
version: 3.5
actif: true

parametres:
  environnement: "production"
  tags: ["backend", "critique"]
  quota:
    cpu: 128
    ram_gb: 512

serveurs:
  - nom: "web-01"
    ip: "10.0.0.10"
    roles: ["nginx", "frontend"]
  - nom: "db-01"
    ip: "10.0.0.20"
    roles: ["postgres"]
    master: true

divers: ["texte", 42, false, {"sous_cle": "sous_valeur"}]
`

	// Le résultat HCL attendu, avec toutes les clés triées par ordre alphabétique
	// hclwrite gère automatiquement l'indentation et l'alignement des signes '='
	expectedHCL := `actif = true
divers = ["texte", 42, false, {
  sous_cle = "sous_valeur"
}]
parametres = {
  environnement = "production"
  quota = {
    cpu    = 128
    ram_gb = 512
  }
  tags = ["backend", "critique"]
}
projet = "mega-infrastructure"
serveurs = [{
  ip    = "10.0.0.10"
  nom   = "web-01"
  roles = ["nginx", "frontend"]
  }, {
  ip     = "10.0.0.20"
  master = true
  nom    = "db-01"
  roles  = ["postgres"]
}]
version = 3.5
`

	// Exécution
	resultBytes, err := ToHCL2([]byte(inputYAML))

	// Assertions
	assert.NoError(t, err, "Le programme n'aurait pas dû planter sur ce gros YAML")
	assert.Equal(t, strings.TrimSpace(expectedHCL), strings.TrimSpace(string(resultBytes)), "Le HCL généré ne correspond pas au format attendu")
}