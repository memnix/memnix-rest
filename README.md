# Memnix Rest

Memnix Rest API

[![Go Report Card](https://goreportcard.com/badge/github.com/memnix/memnix-rest)](https://goreportcard.com/report/github.com/memnix/memnix-rest) [![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/memnix/memnix-rest/badges/quality-score.png?b=main)](https://scrutinizer-ci.com/g/memnix/memnix-rest/?branch=main) [![Build Status](https://scrutinizer-ci.com/g/memnix/memnix-rest/badges/build.png?b=main)](https://scrutinizer-ci.com/g/memnix/memnix-rest/build-status/main) ![GitHub](https://img.shields.io/github/license/Memnix/memnix-rest?style=flat-square) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/memnix/memnix-rest?style=flat-square)

---

## Memnix REST API

Memnix REST is the rest API that is used by user interfaces to create decks and cards and play on MemnixAPP. It's using [Gofiber](https://github.com/gofiber/fiber) to 
handle requests and [Gorm](https://github.com/go-gorm/gorm) as a layer for Postgres.

## Env

Create a **.env** file

```env
# This is a sample config file

DB_USER="user"
DB_PASSWORD="password"
DB_PORT ="port"
DB_DB = "mydb"
DB_HOST="localhost"
```

## Version v0.1.0-beta4
