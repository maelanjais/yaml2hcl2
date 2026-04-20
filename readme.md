# YAML to HCL2 Converter

## Description
Ce projet est un utilitaire de ligne de commande écrit en Go, conçu pour transformer des configurations YAML complexes en fichiers HCL2 (HashiCorp Configuration Language). L'outil garantit un résultat déterministe et syntaxiquement valide, prêt pour une utilisation immédiate avec Terraform ou d'autres outils de l'écosystème HashiCorp.

## Spécifications Techniques

| Composant | Technologie | Rôle |
| :--- | :--- | :--- |
| Langage | Go 1.21+ | Moteur d'exécution principal |
| Analyseur YAML | gopkg.in/yaml.v3 | Décodage des fichiers sources |
| Générateur HCL | hclwrite (v2) | Construction de l'AST et formatage |
| Validation | hclparse | Validation syntaxique pré-enregistrement |

## Fonctionnalités Avancées

### 1. Robustesse et Auto-Validation
L'outil intègre une étape de validation interne systématique. Avant toute écriture sur disque, le HCL généré est analysé par le parseur officiel de HashiCorp. Si une anomalie est détectée, le processus est interrompu pour éviter la génération de fichiers corrompus.

### 2. Nettoyage des Identifiants (Sanitization)
Pour garantir la conformité avec la syntaxe HCL, les clés YAML sont automatiquement traitées :
- Remplacement des caractères spéciaux et espaces par des underscores (`_`).
- Préfixage des clés commençant par un chiffre.
- Gestion des valeurs `null` et des structures vides.

### 3. Sortie Déterministe
Toutes les clés sont triées par ordre alphabétique à chaque niveau d'imbrication, assurant que le même YAML produira toujours exactement le même fichier HCL (indispensable pour le versioning Git).

## Architecture du Projet

- `main.go` : Gestion du flux d'entrée/sortie et interface utilisateur.
- `internal/converter/` :
  - `converter.go` : Algorithme récursif de conversion et sanitization.
  - `converter_test.go` : Suite de 10 tests unitaires intensifs incluant la validation syntaxique.

## Utilisation

### Installation
```bash
go mod tidy
go build -o yaml2hcl