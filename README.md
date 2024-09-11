# Altcha server

Serveur de génération de challenges altcha et de validation de la solution

# Utilisation

Lancer le serveur
```
altcha run
```

# Variables d'environement
| Nom                 | Description                                                                  | Valeur par défaut        | Requis |
|---------------------|------------------------------------------------------------------------------|--------------------------|--------|
| ALTCHA_PORT         | Port d'écoute du serveur                                                     | 3333                     | Non    |
| ALTCHA_HMAC_KEY     | Clé d'encodage des signatures                                                |                          | Oui    |
| ALTCHA_MAX_NUMBER   | Nombre d'itération maximum pour résoudre le challenge (défini la difficulté) | 1000000                  | Non    |
| ALTCHA_ALGORITHM    | Algorithme de hashage (valeurs possibles: SHA-1, SHA-256, SHA-512)           | SHA-256                  | Non    |
| ALTCHA_SALT         | Forcer le salt du challenge                                                  | *Généré automatiquement* | Non    |
| ALTCHA_EXPIRE       | Temps avant expiration du challenge (en secondes)                            | 600                      | Non    |
| ALTCHA_CHECK_EXPIRE | Vérifier si le challenge à expiré                                            | 1                        | Non    |