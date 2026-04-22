# YAML2HCL2 - Système de Plugins de Configuration

## Description

Ce projet est un utilitaire modulaire écrit en Go, conçu pour analyser, convertir et évaluer des fichiers de configuration. Initialement conçu comme un simple convertisseur YAML vers HCL2, le projet a évolué vers une architecture extensible basée sur le système **`go-plugin` d'HashiCorp**.

Cette architecture permet à l'application principale de charger dynamiquement des sous-processus (plugins) capables de traiter différents formats de fichiers (YAML, HCL) via des appels RPC.

## Architecture

L'application repose sur le modèle de plugin d'HashiCorp (utilisé par Terraform, Vault, etc.) :
- **Hôte (`main.go`)** : Point d'entrée du programme. Il détecte l'extension du fichier cible, lance le processus binaire du plugin approprié en arrière-plan, et communique avec lui via gRPC/NetRPC.
- **Interface Partagée (`shared/interface.go`)** : Contient le contrat strict `Converter` que chaque plugin doit implémenter.
- **Plugins (`plugins/`)** : Exécutables compilés indépendamment qui effectuent le traitement lourd.

### Les Plugins Actuels

1. **Plugin YAML (`plugins/yaml/main.go` & `internal/yaml/`)** :
   - Prend en entrée un fichier YAML.
   - Construit l'AST (Abstract Syntax Tree) HCL via `hclwrite`.
   - Effectue une validation syntaxique avec `hclparse` pour garantir un fichier propre.
   - Génère un fichier de sortie `.hcl`.

2. **Plugin HCL (`plugins/hcl/main.go` & `internal/hcl/`)** :
   - Prend en entrée un fichier HCL.
   - Construit un `hcl.EvalContext` permettant l'évaluation d'expressions.
   - Supporte des **variables** (ex: variables d'environnement via l'objet `env`).
   - Supporte des **fonctions** (ex: `upper()`, `lower()`, `max()`, `min()`).
   - Évalue la configuration et la restitue sous forme de fichier `.json` pour validation.

---

## Guide d'Utilisation et de Test

### Prérequis
- **Go 1.21+** installé sur votre machine.
- `make` pour l'utilisation des commandes de compilation (inclus par défaut sur macOS/Linux).

### 1. Compilation du Projet

L'architecture go-plugin nécessite que l'application principale **et** chaque plugin soient compilés en tant que binaires distincts.
Pour tout compiler d'un seul coup, utilisez la commande `make` :

```bash
make
```

> **Ce que ça fait** :
> 1. Crée les dossiers nécessaires.
> 2. Compile le plugin YAML (`plugins/yaml/plugin`).
> 3. Compile le plugin HCL (`plugins/hcl/plugin`).
> 4. Compile l'exécutable hôte (`yaml2hcl2` à la racine).

### 2. Tester le Plugin YAML (YAML -> HCL)

Testez la conversion du fichier d'exemple inclus :

```bash
./yaml2hcl2 testSupfile.yml
```

**Résultat attendu** :
Le programme lancera le plugin YAML via RPC, qui générera un fichier nommé `output.hcl` contenant la représentation syntaxique valide de votre YAML. Vous verrez un message :
`Opération terminée ! Le résultat est disponible dans le fichier 'output.hcl'`

### 3. Tester le Plugin HCL (Évaluation HCL -> JSON)

Testez l'évaluation dynamique de variables et de fonctions :

```bash
./yaml2hcl2 test.hcl
```

> *Note : Si vous examinez `test.hcl`, vous verrez qu'il utilise des fonctions comme `upper()` et appelle des variables comme `env.USER`.*

**Résultat attendu** :
Le programme lancera le plugin HCL, qui évaluera le code et générera un fichier `output.json`.
Si vous ouvrez `output.json`, vous verrez que les fonctions (ex: "MY-SUPER-PROJECT") et les variables systèmes ont bien été résolues !

### 4. Nettoyage de l'Espace de Travail

Pour supprimer tous les binaires compilés et les fichiers générés par les tests, exécutez simplement :

```bash
make clean
```

---

## Pour aller plus loin

- **Ajouter de nouvelles fonctions** : Modifiez `internal/hcl/converter.go` dans la map `Functions` du `hcl.EvalContext` (vous pouvez importer des fonctions depuis `github.com/zclconf/go-cty/cty/function/stdlib`).
- **Ajouter de nouvelles variables** : Modifiez l'objet passé dans la map `Variables` du contexte.
- **Ajouter un nouveau plugin** : Créez un nouveau dossier sous `plugins/`, implémentez l'interface `shared.Converter`, compilez-le en tant que binaire, et ajoutez la condition dans `main.go`.