# ssh-test

## Prerequisites

Author used `go version go1.12.13 darwin/amd64`.
You must have `docker`, `docker-compose`, `openssl` installed on your computer.

### Reading

How to [generate certificate](https://gist.github.com/cecilemuller/9492b848eb8fe46d462abeb26656c4f8#how-to-create-an-https-certificate-for-localhost-domains).
How to [set up custom server](https://github.com/go-swagger/go-swagger/blob/master/docs/tutorial/custom-server.md)

I have decided to create the game that will not be a stupid one.
Hence, this I thought I could use Markov process.

I found [this article](https://ambareeshr.wordpress.com/2019/01/07/markov-decision-process-for-tic-tac-toe/) and
[this repo](https://github.com/revanurambareesh/mdp-tic-tac-toe).

## How to start
```.bash
# generate code
make generate

# install dependencies
make deps

# generate certificate
make generate_certificate

# start server and postgres
docker-compose up
```

## How to play

Load [this file](./ssh.postman_collection.json) (it is a collection) to Postman and play.

## How to run tests

```bash
make unit_test
```