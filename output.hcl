serveur = {
  domaine = "api.mon-projet.com"
  port    = 8080
  tags = {
    equipe   = "backend"
    priorite = 1
  }
}
reseaux_autorises = ["192.168.1.0/24", "10.0.0.0/8"]
base_de_donnees = {
  connexions_max = 100
  ssl_actif      = true
  type           = "postgresql"
}
nom_du_projet = "super-api"
version       = 1.5
en_production = false
