# Contributing

([Français](#comment-contribuer))

## How to Contribute

When contributing, post comments and discuss changes you wish to make via Issues.

Feel free to propose changes by creating Pull Requests. If you don't have write access, editing a file will create a Fork of this project for you to save your proposed changes to. Submitting a change to a file will write it to a new Branch in your Fork, so you can send a Pull Request.

If this is your first time contributing on GitHub, don't worry! Let us know if you have any questions.

### Security

**Do not post any security issues on the public repository!** See [SECURITY.md](SECURITY.md)

______________________

## Comment contribuer

Lorsque vous contribuez, veuillez également publier des commentaires et discuter des modifications que vous souhaitez apporter par l'entremise des enjeux (Issues).

N'hésitez pas à proposer des modifications en créant des demandes de tirage (Pull Requests). Si vous n'avez pas accès au mode de rédaction, la modification d'un fichier créera une copie (Fork) de ce projet afin que vous puissiez enregistrer les modifications que vous proposez. Le fait de proposer une modification à un fichier l'écrira dans une nouvelle branche dans votre copie (Fork), de sorte que vous puissiez envoyer une demande de tirage (Pull Request).

Si c'est la première fois que vous contribuez à GitHub, ne vous en faites pas! Faites-nous part de vos questions.

### Sécurité

**Ne publiez aucun problème de sécurité sur le dépôt publique!** Voir [SECURITY.md](SECURITY.md)








---
## Contributing Guide

## File Structure

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

### Git Branching Strategy

- Volume Cleaner follows a [Github Flow Branching Strategy](https://www.gitkraken.com/learn/git/best-practices/git-branch-strategy#github-flow-branch-strategy)

### Branch Naming Convention

Branch names follow the following guidelines when naming branches

- Lowercase and Hyphen-separated ([Kebab Case](https://developer.mozilla.org/en-US/docs/Glossary/Kebab_case)) - for example `feature/btis-100-new-login` or `bugfix/btis-50-header-styling`

- Alphanumeric Characters - Only use Alphanumeric Characters (a-z,A-Z,0-9) and hyphens. Avoid punctuation, spaces, underscores.

- No Continuous or Trailing Hyphens - for example `feature/btis-120--new--login` or `feature/btis-30-new-login-` 


For Branch prefix names, the prefix used in branch names is flexible and does not follow hard guidelines. However, please be reasonable in your prefix names. Here are some common prefixes.

- `feature/` - For new features or enhancements (e.g., `feature/btis-1000-user-authentication`)

- `bugfix/` - For fixing bugs (e.g., `bugfix/btis-20-login-error`)

- `docs/` - For documentation changes (e.g., `docs/btis-10-api-guide`)

- `ci/` - For changes to CI/CD pipelines or config (e.g., `ci/btis-30-github-actions`)

- `chore/` - For routine tasks like dependency updates or refactoring (e.g., `chore/btis-180-update-deps`)

Here are some sample branch names:

- `docs/btis-1001-contribution-guide`

- `ci/btis-1002-gh-actions-ci-jobs`

- `feature/btis-1000-go-graph-client`

### Commit Messages

- Commit Messages are to follow the structure of [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification)

### Issue Structure

Bug Reports are to be structured with the following headings

- Describe the bug - _A clear and concise description of what the bug is_ 

- Environment Info - _Details on what environment used when the bug was encountered_

- Files - _A list of relevant files for this issue_

- Steps to reproduce - _Steps to reproduce the behaviour_

- Expected Behaviour - _A clear and concise description of what you expected to happen_

- Screenshots - _Screenshots of the problem_

- Additional Context - _Other details_ 

Feature Requests are to be structured with the following headings

- Is your feature request related to a problem? Please link issue ticket

- Describe the solution you'd like - _A clear and concise description of what you want to happen_

- Describe alternatives you've considered - _A clear and concise description of any alternative solutions or features you've considered_

- Additional Context - _Other details_

### Pull Request Structure

Pull Requests (PRs) are to be structured with the following headings

- Title - _Concise Summary of the change_

- Proposed Changes/Description - _Details on what the PR does, Why the change is needed, and any additional context about this PR_

- Type of Change - _labels for the change, based on the PR template in the project (i.e., bugfix, new feature, breaking change, documentation update)_

- Testing - _Steps to reproduce the issue and verify the change/fix_

- Screenshots (if applicable) - _Include before/after screenshots or links_

- Related Issue/Ticket - _Link to the related issue(s) on Github or Jira_
