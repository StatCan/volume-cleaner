# Kubernetes Volume Cleaner

<p align="center">
    <img src="./assets/volume_cleaner.jpg" alt="Volume Cleaner Logo" width="400"/>
</p>

A Kubernetes CronJob that automatically identifies and cleans up stale Persistent Volume Claims and Persistent Volumes (K8S) linked to an associated Azure disk.

# Contents

- [Requirements](#-requirements)
- [Documentation](#-documentation)
- [How to Contribute](#-how-to-contribute)
- [Code of Conduct](#code-of-conduct)
- [License](#-license)

# Requirements

# Documentation

## File Structure

This project follows the [golang-standard project layout](https://github.com/golang-standards/project-layout). Furthermore, file and directory names are written in `snake_case`

```
.github/
├─ ISSUE_TEMPLATE/
├─ workflows/
├─ PULL_REQUEST_TEMPLATE.md
assets/
cmd/
├─ controller/
├─ scheduler/
configs/
├─ controller/
├─ scheduler/
internal/
├─ kubernetes/
├─ structure/
├─ ...
manifests/
├─ controller/
├─ scheduler/
scripts/
```

Exception: This project does not use a `test` folder, `_test.go` files are keep in their package folders

# How to Contribute

## Code of Conduct

## Contributing Guide

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

## License

