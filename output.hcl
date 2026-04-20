chaine_vide = ""
cluster_kubernetes = {
  haute_disponibilite = true
  nom                 = "k8s-prod"
  nœuds = {
    controle = {
      labels = {
        environnement = "production"
        role          = "master"
      }
      quantite      = 3
      type_instance = "t3.medium"
    }
    travailleurs = {
      labels = {
        gpu_accelere = false
        role         = "worker"
      }
      quantite_max  = 20
      quantite_min  = 5
      type_instance = "m5.large"
    }
  }
  version = "1.28.2"
}
liste_vide  = []
mode_strict = true
objet_vide  = {}
parametres_globaux = {
  multi_zones       = true
  region_principale = "eu-west-3"
  taux_erreurs_max  = 0.05
  timeout_secondes  = 30
}
politiques_securite = {
  regles_pare_feu = [{
    action = "allow"
    meta_donnees = {
      audite = true
      tags   = ["web", 80, true]
    }
    nom      = "autoriser-http"
    priorite = 100
    sources  = ["0.0.0.0/0"]
    }, {
    action       = "deny"
    meta_donnees = {}
    nom          = "bloquer-tout"
    priorite     = 999
    sources      = ["0.0.0.0/0"]
  }]
}
ports_ouverts = [80, 443, 22]
reseaux_virtuels = [{
  actif = true
  cidr  = "10.0.0.0/16"
  nom   = "vpc-frontend"
  sous_reseaux = [{
    cidr       = "10.0.1.0/24"
    nom        = "subnet-public-1"
    passerelle = "10.0.1.1"
    }, {
    cidr       = "10.0.2.0/24"
    nom        = "subnet-public-2"
    passerelle = "10.0.2.1"
  }]
  }, {
  actif = true
  cidr  = "10.1.0.0/16"
  nom   = "vpc-backend"
  sous_reseaux = [{
    cidr       = "10.1.1.0/24"
    nom        = "subnet-prive-1"
    passerelle = null
  }]
  }, {
  actif        = false
  cidr         = "192.168.0.0/16"
  nom          = "vpc-archive"
  sous_reseaux = []
}]
serveurs_dns    = ["1.1.1.1", "8.8.8.8"]
titre_projet    = "infrastructure-globale"
tuple_de_stress = ["chaine", 42, 3.14, false]
valeur_nulle    = null
version_config  = 4.2
