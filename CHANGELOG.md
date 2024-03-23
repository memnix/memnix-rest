## v2.0.0-a7 (2024-03-15)

### 📌➕⬇️ ➖⬆️  Dependencies

- upgrade deps and remove fiber

## v2.0.0-a6 (2024-03-15)

### ⚡️ Performance

- **middlewares**: better config for performances

## v2.0.0-a5 (2024-03-11)

### 🐛🚑️ Fixes

- **docker**: fix path

## v2.0.0-a4 (2024-03-11)

### 🔒️ Security

- **csp**: insane csp

## v2.0.0-a3 (2024-03-10)

### ♻️  Refactorings

- **config**: remove old config files

### fix

- **deps**: update all non-major dependencies
- **deps**: update all non-major dependencies
- **deps**: update module github.com/gofiber/contrib/otelfiber to v2
- **deps**: update module github.com/gofiber/contrib/otelfiber to v2
- **deps**: update golang.org/x/exp digest to 814bf88
- **deps**: update golang.org/x/exp digest to 814bf88

### 💚👷 CI & Build

- fix small issuers

## v2.0.0-a2 (2024-03-10)

### 🔒️ Security

- **password**: verify password strengh >>> ⏰ 3h

### 🎨🏗️ Style & Architecture

- **pgx**: improve pgx config

## v2.0.0-a1 (2024-03-10)

### 💥 Boom

- switch to sqlc >>> ⏰ 3h

### ♻️  Refactorings

- **views**: better ui

### BREAKING CHANGE

- remove gorm for sqlc

### ✅🤡🧪 Tests

- **views**: update tests
- **tests**: remove tests that need the db

### 💄🚸 UI & UIX

- **flashmessages**: auto remove after a timer

### 💚👷 CI & Build

- fix weird stuff

### 📌➕⬇️ ➖⬆️  Dependencies

- **deps**: update deps

### 🔐🚧📈✏️ 💩👽️🍻💬🥚🌱🚩🥅🩺 Others

- **sqlc**: try sqlc

### 🔧🔨📦️ Configuration, Scripts, Packages

- **devbox**: add atlas

### 🗃️ Database

- add atlas migration

### 🚨 Linting

- fix packages

## v2.0.0-a0 (2024-02-24)

### ✨ Features

- **auth**: register page
- **views**: image optimization
- **pkg**: image optimization
- **auth**: login error
- **auth**: login template
- **handlers**: create page handler struct >>> ⏰ 45m
- **v2**: tailwind >>> ⏰ 2h

### 🐛🚑️ Fixes

- remove config file

### 🔒️ Security

- **pre-commit**: add ggshield hook

### ♻️  Refactorings

- **views**: improved flashmessages
- **config**: better config management >>> ⏰ 1h30
- **api**: moved http/views

### BREAKING CHANGE

- Echo replaces Fiber
- change api to htmx handlers

### ⚗️  Experiments

- **v2**: jwt middleware >>> ⏰ 2h

### ✅🤡🧪 Tests

- **views**: test register component
- **infrastructures**: add tests
- **pkg**: add tests to utils
- **images**: test for images conversion
- **pkg**: crypto tests
- **domain**: add tests
- **domain**: card tests
- **pkg**: add tests
- **views**: add tests

### 👔 logic

- **setup**: setup v2 with infra and htmx >>> ⏰ 2h

### 💄🚸 UI & UIX

- **tailwind**: fonts

### 💚👷 CI & Build

- **devbox**: auto install tools
- **hadolint**: wget -q >>> ⏰ 1m
- **docker**: docker image building >>> ⏰ 10m

### 🔐🚧📈✏️ 💩👽️🍻💬🥚🌱🚩🥅🩺 Others

- **app**: htmx implementation >>> ⏰ 1h

### 🔥⚰️  Clean up

- **json**: remove useless benchmark

### 🙈 Ignore

- remove config files

## v1.2.2 (2024-02-20)

### 💚👷 CI & Build

- **devbox**: better commit lint >>> ⏰ 20m
