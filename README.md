# Kubernetes Volume Cleaner

<p align="center">
    <img src="./assets/volume_cleaner.jpg" alt="Volume Cleaner Logo" width="400"/>
</p>

A Kubernetes CronJob that automatically identifies and cleans up stale Persistent Volume Claims and Persistent Volumes (K8S) linked to an associated Azure disk.

# Contents

- [Requirements](#requirements)
- [Documentation](#documentation)
- [How to Contribute](#how-to-contribute)
- [Code of Conduct](#code-of-conduct)
- [License](#license)

# Requirements

# Documentation

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

# How to Contribute

See [CONTRIBUTING.md](CONTRIBUTING.md)

## Code of Conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

## License

Volume Cleaner is distributed under [AGPL-3.0.only](LICENSE.md).

______________________

