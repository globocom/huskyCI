# Husky: Security CI

Husky will make security tests inside a CI.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

```
brew install vagrant
```

## Installing

#### Cloning Husky's repository:

```
cd $GOPATH && cd src/ && mkdir github.com && cd github.com && mkdir globocom && cd globocom
```

```
git clone https://github.com/globocom/husky.git && cd husky
```

#### Starting up VMs:

```
vagrant up vm2-db
```

```
vagrant up vm3-docker
```

#### Setting up environment variables (user your own configuration):

```
echo 'export DOCKER_HOST="192.168.50.6:2376"' > .env
echo 'export MONGO_HOST="192.168.50.5"' >> .ev
echo 'export MONGO_NAME="huskyDB"' >> .env
echo 'export MONGO_USER="husky"' >> .env
echo 'export MONGO_PASS="superENVPassword"' >> .env
```

```
source .env
```

#### Inserting a new MongoDB user:

```
vagrant ssh vm2-db
```

```
sudo su && mongo
```

```
use huskyDB
```

```
db.createUser({user:"husky", pwd:"superENVPassword", roles: ["readWrite"]})
```

#### Inserting securityTests examples:

```
go run server.go
```

```
curl -H "Content-Type: application/json" -d '{"name":"brakeman", "image": "huskyci/brakeman", "cmd": ["ifconfig", "whoami", "ls"], "language": "ruby", "default":true}' http://localhost:9999/securitytest
```

```
curl -H "Content-Type: application/json" -d '{"name":"enry", "image": "huskyci/enry", "cmd": ["ifconfig", "whoami", "ls"], "language": "generic", "default":true}' http://localhost:9999/securitytest
```

```
curl -H "Content-Type: application/json" -d '{"name":"gas", "image": "huskyci/gas", "cmd": ["ifconfig", "whoami", "ls"], "language": "golang", "default":true}' http://localhost:9999/securitytest
```

## Architecture draft

![architecture](images/architecture-draft.png)

## MongoDB draft

![db](images/mongoBD-draft.png)

