# Contributing

([Français](#comment-contribuer))

## How to Contribute

When contributing, post comments and discuss changes you wish to make via Issues.

Issues follow the template structures indicated in `.github/ISSUE_TEMPLATES/`

Feel free to propose changes by creating Pull Requests. If you don't have write access, editing a file will create a Fork of this project for you to save your proposed changes to. Submitting a change to a file will write it to a new Branch in your Fork, so you can send a Pull Request.

Pull Requests follow the template structure indicated in `.github/PULL_REQUEST_TEMPLATE.md`

If this is your first time contributing on GitHub, don't worry! Let us know if you have any questions.

### File structure

This project follows the [golang-standard project layout](https://github.com/golang-standards/project-layout). Furthermore, file and directory names are written in `snake_case`

```
volume-cleaner/
├─ .github/
│  ├─ configs/
│  ├─ ISSUE_TEMPLATES/
│  ├─ workflows/
│  ├─ PULL_REQUEST_TEMPLATE.md
├─ assets/
├─ cmd/
│  ├─ controller/
│  ├─ scheduler/
├─ docker/
│  ├─ controller/
│  ├─ scheduler/
├─ internal/
│  ├─ kubernetes/
│  ├─ structure/
│  ├─ utils/
│  ├─ ...
├─ manifests/
│  ├─ controller/
│  ├─ scheduler/
├─ scripts/
│  ├─ controller/
│  ├─ scheduler/
├─ testing/
```

### Testing

This project does not use a test folder, `_test.go` files are kept in their package folders, unit testing is performed using Golang's standard testing library along with the popular testing toolkit [Testify](https://github.com/stretchr/testify). Integration testing is performed in Github actions workflows defined within the `.github/workflows/` directory. Testing can be performed locally using the following:

```
// Unit and Integration Testing
go test -v -race ./...

// End-to-End Testing
act --rm
```

### Practices 

#### Branching

- Volume Cleaner follows a [Github Flow Branching Strategy](https://www.gitkraken.com/learn/git/best-practices/git-branch-strategy#github-flow-branch-strategy)

**Branch names follow the following guidelines when naming branches:**

- Branch names are structured like so `<type>/<issue number>-<detail>` 

- Lowercase and Hyphen-separated ([Kebab Case](https://developer.mozilla.org/en-US/docs/Glossary/Kebab_case)) - for example `feature/100-new-login` or `bugfix/50-header-styling`

- Alphanumeric Characters - Only use Alphanumeric Characters (a-z,A-Z,0-9) and hyphens. Avoid punctuation, spaces, underscores.

- No Continuous or Trailing Hyphens - for example `feature/120--new--login` or `feature/30-new-login-` 


For Branch prefix names, the prefix used in branch names is flexible and does not follow hard guidelines. However, please be reasonable in your prefix names. Here are some common prefixes.

- `feature/` : For new features or enhancements (e.g., `feature/1000-user-authentication`)

- `bugfix/` : For fixing bugs (e.g., `bugfix/20-login-error`)

- `docs/` : For documentation changes (e.g., `docs/10-api-guide`)

- `ci/` : For changes to CI/CD pipelines or config (e.g., `ci/30-github-actions`)

- `chore/` : For routine tasks like dependency updates or refactoring (e.g., `chore/180-update-deps`)

### Commit Messages

- Commit Messages are to follow the structure of [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification)

### Security

**Do not post any security issues on the public repository!** See [SECURITY.md](SECURITY.md)

______________________

## Comment contribuer

Lorsque vous contribuez, veuillez également publier des commentaires et discuter des modifications que vous souhaitez apporter par l'entremise des enjeux (Issues).

Les Issues suivent les modèles indiqués dans `.github/ISSUE_TEMPLATES/`

N'hésitez pas à proposer des modifications en créant des demandes de tirage (Pull Requests). Si vous n'avez pas accès au mode de rédaction, la modification d'un fichier créera une copie (Fork) de ce projet afin que vous puissiez enregistrer les modifications que vous proposez. Le fait de proposer une modification à un fichier l'écrira dans une nouvelle branche dans votre copie (Fork), de sorte que vous puissiez envoyer une demande de tirage (Pull Request).

Les Pull Requests suivent le modèle indiqué dans `.github/PULL_REQUEST_TEMPLATE.md`

Si c'est la première fois que vous contribuez à GitHub, ne vous en faites pas! Faites-nous part de vos questions.

### Structure des fichiers

Ce projet suit le [golang-standard project layout](https://github.com/golang-standards/project-layout). De plus, les noms de fichiers et de répertoires sont en `snake_case`

```
volume-cleaner/
├─ .github/
│  ├─ configs/
│  ├─ ISSUE_TEMPLATES/
│  ├─ workflows/
│  ├─ PULL_REQUEST_TEMPLATE.md
├─ assets/
├─ cmd/
│  ├─ controller/
│  ├─ scheduler/
├─ docker/
│  ├─ controller/
│  ├─ scheduler/
├─ internal/
│  ├─ kubernetes/
│  ├─ structure/
│  ├─ utils/
│  ├─ ...
├─ manifests/
│  ├─ controller/
│  ├─ scheduler/
├─ scripts/
│  ├─ controller/
│  ├─ scheduler/
├─ testing/
```

### Tests

Ce projet n'utilise pas de dossier de test ; les fichiers `_test.go` restent dans leur dossier de package. Les tests unitaires sont réalisés avec la bibliothèque de tests standard de Golang ainsi que le populaire toolkit [Testify](https://github.com/stretchr/testify). Les tests d'intégration sont exécutés via les workflows GitHub Actions définis dans le répertoire `.github/workflows/`. Les tests peuvent être lancés localement avec :

```

// Unit and Integration Testing
go test -v -race ./...

// End-to-End Testing
act --rm

```

### Bonnes Pratiques

#### Gestion des branches

- Volume Cleaner suit une [stratégie de branching GitHub Flow](https://www.gitkraken.com/learn/git/best-practices/git-branch-strategy#github-flow-branch-strategy).

**Les noms de branches suivent ces règles :**

- Structure `<type>/<numéro-de-issue>-<détail>`

- Minuscules et séparés par des tirets ([Kebab Case](https://developer.mozilla.org/en-US/docs/Glossary/Kebab_case)), par exemple : `feature/100-new-login` ou `bugfix/50-header-styling`.

- Caractères alphanumériques uniquement (a-z, A-Z, 0-9) et des tirets. Évitez la ponctuation, les espaces ou les underscores.

- Pas de tirets doubles ou en fin de nom, par exemple `feature/120--new--login` ou `feature/30-new-login-` sont à proscrire.

Pour le préfixe des branches, soyez raisonnable : bien qu'aucune règle stricte n'impose un standard, voici des préfixes courants :

- `feature/` : pour les nouvelles fonctionnalités (ex. `feature/1000-user-authentication`)

- `bugfix/` : pour corriger des bugs (ex. `bugfix/20-login-error`)

- `docs/` : pour les modifications de documentation (ex. `docs/10-api-guide`)

- `ci/` : pour les changements CI/CD ou config (ex. `ci/30-github-actions`)

- `chore/` : pour les tâches de routine (ex. `chore/180-update-deps`)

### Messages de commit

- Les messages de commit doivent suivre la convention [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification).

### Sécurité

**Ne publiez aucun problème de sécurité sur le dépôt publique!** Voir [SECURITY.md](SECURITY.md)









