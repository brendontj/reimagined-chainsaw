# Reimagined Chainsaw Project

This project is a real time conversation between different users. Also, it's possible for a user query a stock price in real time.

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Features](#features)

## Requirements

- Go 1.21 or later
- Docker
- Additional dependencies as specified in `go.mod`

## Installation

1. Clone this repository:

   ```sh
   git clone git@github.com:brendontj/reimagined-chainsaw.git
   cd reimagined-chainsaw
   ```

2. Install dependencies:

   ```sh
   go mod download
   ```

3. Start docker container with rabbitmq image: 
    ```sh
    docker run -d --hostname my-rabbit --name some-rabbit -p 5672:5672 rabbitmq:3
    ```

4. Run the bot:
    
   ```sh
   go run stock_prices_bot/main.go
   ```

5. Run the project:

   ```sh
   go run main.go
   ```
   
6. Access the main page:
   - localhost:8000

## Features

- Multi channels
- Sign up User
- Sign in User
- Messages throughout each specific channel
- Query of stock closing prices
- (not implemented) Show only the last 50 messages per chat.
- (not implemented) Unit tests.
- (not implemented) Persist data in a database and not in a memory storage.)



