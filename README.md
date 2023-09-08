# Memnix Rest

Memnix Rest API

[![Go Report Card](https://goreportcard.com/badge/github.com/memnix/memnixrest)](https://goreportcard.com/report/github.com/memnix/memnixrest) [![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/memnix/memnix-rest/badges/quality-score.png?b=main)](https://scrutinizer-ci.com/g/memnix/memnix-rest/?branch=main) [![Build Status](https://scrutinizer-ci.com/g/memnix/memnix-rest/badges/build.png?b=main)](https://scrutinizer-ci.com/g/memnix/memnix-rest/build-status/main) ![GitHub](https://img.shields.io/github/license/Memnix/memnix-rest?style=flat-square) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/memnix/memnix-rest?style=flat-square) [![DeepSource](https://deepsource.io/gh/memnix/memnix-rest.svg/?label=active+issues&token=jIwgCj7nvzKqjNGXfdGIQwvJ)](https://deepsource.io/gh/memnix/memnix-rest/?ref=repository-badge) [![Maintainability](https://api.codeclimate.com/v1/badges/7d4403b2b97cf390983f/maintainability)](https://codeclimate.com/github/memnix/memnix-rest/maintainability) [![Codeac.io](https://static.codeac.io/badges/2-413846692.svg "Codeac.io")](https://app.codeac.io/github/memnix/memnix-rest)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest?ref=badge_shield)[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest.svg?type=shield&issueType=security)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest?ref=badge_shield)

---

## Memnix REST API

Memnix REST is the rest API that is used by user interfaces to create decks and cards and play on MemnixAPP. It's
using [Gofiber](https://github.com/gofiber/fiber) to
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


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest?ref=badge_large)