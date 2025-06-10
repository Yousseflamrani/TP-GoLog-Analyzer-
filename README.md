# GoLog Analyzer

Un outil en ligne de commande (CLI) écrit en Go pour analyser des fichiers de logs provenant de plusieurs sources, en parallèle, et générer un rapport d’analyse.

## Installation

1. **Clone le dépôt :**
   ```bash
   git clone https://github.com/votre-utilisateur/loganalyzer.git
   cd loganalyzer
   ```
2. **Installe les dépendances :**
   ```bash
   go mod download
   ```
3. **Compile et installe (optionnel) :**
   ```bash
   go install
   ```

## Utilisation

### Analyser des logs

```bash
go run main.go analyze --config config.json
```
ou
```bash
go run main.go analyze -c config.json -o rapport.json
```

Le fichier `config.json` doit contenir la liste des logs à analyser (voir l’exemple ci-dessous).

### Exemple de fichier `config.json`

```json
[
  {
    "id": "web-server-1",
    "path": "/var/log/nginx/access.log",
    "type": "nginx-access"
  }
]
```

### Bonus

- **Exporter le rapport dans un fichier JSON**  
  ```bash
  go run main.go analyze -c config.json -o rapport.json
  ```
- **Ajouter un log au fichier de configuration**  
  ```bash
  go run main.go add-log --id "db-server-1" --path "/var/log/db.log" --type "postgres" --file config.json
  ```
- **Filtrer les résultats par statut (OK/FAILED)**  
  ```bash
  go run main.go analyze -c config.json --status FAILED
  ```

## Structure du projet

```
loganalyzer/
├── cmd/
│   ├── root.go
│   └── analyze.go
├── internal/
│   ├── config/
│   ├── analyzer/
│   └── reporter/
├── go.mod
└── README.md
```

## Membres du groupe

- [Jong Hoa CHONG]
- [Youssef ALAOUI]

---

**Projet réalisé dans le cadre du TP GoLog Analyzer.**

Sources
