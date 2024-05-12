# Simple Core Bank
Simple core bank: send and withdraw transaction

Features:
- One user, multiple accounts
- Send/Withdraw money
- Get list of accounts per user
- Transaction histories per account

### API Documentation
You can look to docs directory or open the documentation here

### Tech-stack:
- Supertokens
- Golang(GIN)
- PostgreSQL 16 (GORM)
- Docker
 
### Prerequisites
You need to have:
- [go `1.22`](https://go.dev/doc/install)
- [air](https://github.com/cosmtrek/air)
- [docker cli](https://docs.docker.com/get-docker/)
- [GNU Make](https://www.gnu.org/software/make/) (optional), if you used Linux it's already installed
- good internet connection

## Installation
1. Clone the repository 
2. Open cloned folder in your IDE or text editor
3. Open terminal in the current working directory
4. Copy .env.example to .env (do both account and payment)
5. You need to install air
6. Run make up or `docker compose -f compose.dev.yml up -d`
```sh
$ git clone https://github.com/tryoasnafi/be-assignment
$ cd be-assignment/account
$ cp -p .env.example .env
$ go install
$ cd ../payment
$ cp -p .env.example .env
$ go install
$ go install github.com/cosmtrek/air@latest
```
