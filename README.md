# Lemon-Wallet

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

## Overview

This project provides functionalities to manage _users_ and _movements_ in a wallet.
It is entirely written in the Go language, with a package-oriented design and a mysql database.

## Endpoints

- `POST /users` : Registration of a user. Users with the same alias nor the same email are not allowed.
- `GET /users/:id` : Get a user. 
- `POST /movements` : Register a new movement.
- `GET /movements/search` : List all user movements with optional filters such as: limit, offset, type of movement and currency.

## How To Run This Project