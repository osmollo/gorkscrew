# gorkscrew versions

## 1.0.16

- Improve dependencies installation
- Improve documentation
- Compiled with [GO 1.18.6](https://go.dev/dl)

## 1.0.15

- New workflow in CI for detect vulnerabilities and coding errors

## 1.0.14

- Translate README to english
- Code licensed under GPL3
- Licensed under GPL3
- New argument for enable logging (`--log`). Default: no logging
- Proxy timeout up to 5 secs

## 1.0.13

- New parameter `--version` that shows:
  - Gorkscrew version
  - Go compiler version

## 1.0.10

- CI uses GO version of `release.json`

## 1.0.9

- Manual update of go compiler to 1.15.1 in Client

## 1.0.8

- Deleted gitlab CI file. This repository will be the main git remote

## 1.0.6

- Functional CI:
  - Test for every commit in a PR
  - Create release and create assets for go binary and gorkscrew md5sum

## 1.0.4

- Send `git-shell-command` and `userauth` to squid for logging

## 1.0.1

- Send `repository` to squid for logging

## 1.0.0

- First functional release
