# Memnix Rest

Memnix Rest API

[![Go Report Card](https://goreportcard.com/badge/github.com/memnix/memnixrest)](https://goreportcard.com/report/github.com/memnix/memnixrest) [![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/memnix/memnix-rest/badges/quality-score.png?b=main)](https://scrutinizer-ci.com/g/memnix/memnix-rest/?branch=main) [![Build Status](https://scrutinizer-ci.com/g/memnix/memnix-rest/badges/build.png?b=main)](https://scrutinizer-ci.com/g/memnix/memnix-rest/build-status/main) ![GitHub](https://img.shields.io/github/license/Memnix/memnix-rest?style=flat-square) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/memnix/memnix-rest?style=flat-square) [![DeepSource](https://deepsource.io/gh/memnix/memnix-rest.svg/?label=active+issues&token=jIwgCj7nvzKqjNGXfdGIQwvJ)](https://deepsource.io/gh/memnix/memnix-rest/?ref=repository-badge) [![Maintainability](https://api.codeclimate.com/v1/badges/7d4403b2b97cf390983f/maintainability)](https://codeclimate.com/github/memnix/memnix-rest/maintainability) [![Codeac.io](https://static.codeac.io/badges/2-413846692.svg "Codeac.io")](https://app.codeac.io/github/memnix/memnix-rest)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest?ref=badge_shield)[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest.svg?type=shield&issueType=security)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest?ref=badge_shield)

---

## Memnix REST API

Memnix REST is the rest API that is used by user interfaces to create decks and cards and play on MemnixAPP. It's
using [Gofiber](https://github.com/gofiber/fiber) to
handle requests and [Gorm](https://github.com/go-gorm/gorm) as a layer for Postgres.

## Status

I wasn't happy with the state of the project, so I decided to rewrite it from scratch. 
The new version is in the `main` branch. The old version is in the `legacy` branch.

For now, the project is in a very early stage, and it's not ready for production. I haven't yet decided on the architecture, features, and so on. 

At the moment, I'm trying things out and experimenting with different approaches. 

Here are some of the things I'm considering:

### Web Framework

In the previous version, I used Gofiber, but now I'm considering using Echo or Chi to follow the idiomatic Go way of doing things. 
If you have any suggestions, please let me know. 

### Database

In the past, I used Gorm, but I'm considering using SQLBoiler or SQLC. I'm also considering using a NoSQL database like MongoDB for some parts of the application.

### Caching

I'm considering using Redis for caching and eventually for other things like pub/sub. 
I'm using Ristretto for some parts of the application, so it's even faster. 

### Authentication 

I'll keep using JWT for authentication, but I've added support for OAuth2. 

## Contributing

I'm actively looking for contributors. If you're interested in contributing, please let me know.

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmemnix%2Fmemnix-rest?ref=badge_large)