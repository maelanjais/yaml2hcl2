_123_commence_par_chiffre = true
api_endpoint              = null
architecture = {
  couche_1 = {
    couche_2 = {
      couche_3 = {
        couche_4 = {
          couche_5 = {
            identifiant = "profondeur-maximale"
            valide      = true
          }
        }
      }
    }
  }
}
cas_limites = {
  chaine_vide = ""
  liste_vide  = []
  objet_vide  = {}
  zeros       = 0
}
cle_avec_espaces = "valeur"
clusters_regionaux = [{
  metadata = {
    cost_center = 1024
    env         = "prod"
  }
  noeuds = 12
  nom    = "europe-west"
  zones  = ["a", "b", "c"]
  }, {
  metadata = {
    cost_center = 512
    env         = "dev"
  }
  noeuds = 4
  nom    = "us-east"
  zones  = ["a"]
}]
cl__accentu_e = "test"
deprecated    = false
donnees_mixtes = ["chaine", 100, 45.5, false, {
  sous_objet = "valeur"
}, [1, 2, 3]]
identifiant-avec-tirets = 10
infrastructure_id       = "infra-prod-2024"
max_retries             = 5
politiques_securite = {
  entrants = [{
    action   = "ALLOW"
    priorite = 100
    regles = {
      ports      = [80, 443]
      protocoles = ["TCP", "UDP"]
    }
    }, {
    action   = "DENY"
    priorite = 200
    regles = {
      ports      = [0, 65535]
      protocoles = ["*"]
    }
  }]
}
scaling_factor = 1.25
