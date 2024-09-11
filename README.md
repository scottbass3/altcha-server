# Altcha server

Serveur de génération de challenges altcha et de validation de la solution

## Utilisation

### Depuis le binaire
```sh
$ ALTCHA_HMAC_KEY="CLÉ HMAC" bin/altcha run
```

### Depuis l'image docker
```sh
$ docker run -e ALTCHA_HMAC_KEY="CLÉ HMAC" reg.cadoles.com/cadoles/altcha
```

### Depuis les sources
```sh
$ ALTCHA_HMAC_KEY="CLÉ HMAC" go run ./cmd/altcha run
```

### Autres commandes
Générer un challenge
```sh
$ ALTCHA_HMAC_KEY="CLÉ HMAC" bin/altcha generate
```

Résoudre un challenge
```sh
$ ALTCHA_HMAC_KEY="CLÉ HMAC" bin/atlcha solve [CHALLENGE] [SALT]
```

Vérifier une solution
```sh
$ ALTCHA_HMAC_KEY="CLÉ HMAC" bin/altcha verify [CHALLENGE] [SALT] [SIGNATURE] [SOLUTION]
```

## Variables d'environement
| Nom                 | Description                                                                  | Valeur par défaut        | Requis |
|---------------------|------------------------------------------------------------------------------|--------------------------|--------|
| ALTCHA_PORT         | Port d'écoute du serveur                                                     | 3333                     | Non    |
| ALTCHA_HMAC_KEY     | Clé d'encodage des signatures                                                |                          | Oui    |
| ALTCHA_MAX_NUMBER   | Nombre d'itération maximum pour résoudre le challenge (défini la difficulté) | 1000000                  | Non    |
| ALTCHA_ALGORITHM    | Algorithme de hashage (valeurs possibles: SHA-1, SHA-256, SHA-512)           | SHA-256                  | Non    |
| ALTCHA_SALT         | Forcer le salt du challenge                                                  | *Généré automatiquement* | Non    |
| ALTCHA_EXPIRE       | Temps avant expiration du challenge (en secondes)                            | 600                      | Non    |
| ALTCHA_CHECK_EXPIRE | Vérifier si le challenge à expiré                                            | 1                        | Non    |

## Construire le binaire
```sh
$ make build
```

## Construire l'image docker
```sh
$ make build-image
```

## Publier l'image docker
```sh
$ make release-image
```