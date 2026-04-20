# YAML to HCL2 Converter

## Description
Ce projet est un utilitaire écrit en langage Go permettant la conversion de fichiers de configuration du format YAML vers le format HCL2 (HashiCorp Configuration Language). Il a été conçu pour faciliter l'interopérabilité entre les systèmes de gestion de configuration et les outils de l'écosystème HashiCorp tels que Terraform.

## Spécifications Techniques

| Composant | Technologie | Rôle |
| :--- | :--- | :--- |
| Langage | Go 1.21+ | Moteur d'exécution principal |
| Analyseur YAML | gopkg.in/yaml.v3 | Parsing des fichiers sources |
| Générateur HCL | hclwrite (v2) | Construction de l'AST et formatage |
| Typage | zclconf/go-cty | Gestion du système de types HCL2 |

## Architecture du Projet

Le projet respecte une structure modulaire pour garantir la maintenabilité :

* `main.go` : Point d'entrée de l'application et gestion des entrées/sorties fichiers.
* `internal/converter/` : Logique métier isolée.
  * `converter.go` : Algorithme de conversion récursif et tri déterministe.
  * `converter_test.go` : Suite de tests unitaires couvrant les cas nominaux et d'erreurs.

## Utilisation

### Installation
```bash
go mod tidy
go build -o yaml2hcl
```

### Exécution
```bash
./yaml2hcl <source.yaml>
```
Le résultat est généré dans un fichier nommé `output.hcl` à la racine du répertoire.

## Validation de la Qualité
La robustesse de la conversion est garantie par une suite de tests unitaires validant :
1. La conversion des types primitifs (string, number, bool).
2. La gestion des collections complexes (listes mixtes, objets imbriqués).
3. Le respect de l'ordre alphabétique des clés pour le déterminisme.
4. La résilience face aux fichiers sources corrompus.