# huskyCI - Performing security tests inside your CI

<img src="https://raw.githubusercontent.com/wiki/globocom/huskyCI/images/huskyCI-logo.png" align="center" height="" />
<!-- logo font: Anton -->

[![Coverage Status](https://coveralls.io/repos/github/globocom/huskyCI/badge.svg?branch=master)](https://coveralls.io/github/globocom/huskyCI?branch=master) [![CircleCI](https://circleci.com/gh/globocom/huskyCI/tree/master.svg?style=svg&circle-token=415bfb6b5aa0dfce8d2129878a66326da9533150)](https://circleci.com/gh/globocom/huskyCI/tree/master) [![Join the chat at https://gitter.im/globocom/huskyCI](https://badges.gitter.im/globocom/huskyCI.svg)](https://gitter.im/globocom/huskyCI?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

huskyCI is an open-source tool that performs security tests inside CI pipelines of multiple projects and centralizes all results into a database for further analysis and metrics.

## How does it work?

The main goal of this project is to help development teams improve the quality of their code by finding vulnerabilities as soon as possible.

huskyCI can perform static security analysis in Python ([Bandit][Bandit] and [Safety][Safety]), Ruby ([Brakeman][Brakeman]), JavaScript ([Npm Audit][NpmAudit]) and Golang ([Gosec][Gosec]). You should check our [wiki](https://github.com/globocom/huskyCI/wiki/How-does-huskyCI-work%3F) to better understand how this tool could help securing your organization projects!

<p allign="center">
  <img src="huskyCI.gif" />
</p>

## Requirements

#### Docker and Docker-Compose

The easiest way to deploy huskyCI locally is by using [Docker][Docker Install] and [Docker Compose][Docker Compose Install], thus you should have them installed on your machine.

#### Golang

You must also have [Go](https://golang.org/doc/install) installed and huskyCI needs to be inside your [$GOPATH](https://github.com/golang/go/wiki/GOPATH) to run properly.

## Installing

After cloning this repository, simply run the command inside huskyCI's folder:

```
make install
```

## Running

After installing, an `.env` file with instructions to huskyCI should be generated:

```sh
$ cat .env
export HUSKYCI_CLIENT_REPO_URL="https://github.com/globocom/huskyCI.git"
export HUSKYCI_CLIENT_REPO_BRANCH="vulns-Golang"
export HUSKYCI_CLIENT_API_ADDR="http://localhost:8888"
export HUSKYCI_CLIENT_API_USE_HTTPS="false"
export HUSKYCI_CLIENT_TOKEN="{YOUR_TOKEN_HERE}"
```

You can change the repository and branch being analysed by modifying the contents of `HUSKYCI_CLIENT_REPO_URL` and `HUSKYCI_CLIENT_REPO_BRANCH`. Then simply source it through the command:

```sh
. .env
```

Mac OS:
```sh
make run-client
```

Linux:
```sh
make run-client-linux
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests to huskyCI.

## Documentation

You can find huskyCI documentation [here](https://github.com/globocom/huskyCI/wiki).

## License

This project is licensed under the BSD 3-Clause "New" or "Revised" License - read [LICENSE.md](LICENSE.md) file for details.

[Docker Install]:  https://docs.docker.com/install/
[Docker Compose Install]: https://docs.docker.com/compose/install/
[huskyCI Client]: https://github.com/globocom/huskyCI-client
[Bandit]: https://github.com/PyCQA/bandit
[Safety]: https://github.com/pyupio/safety
[Brakeman]: https://github.com/presidentbeef/brakeman
[Gosec]: https://github.com/securego/gosec
[RetireJS]: https://github.com/retirejs/retire.js
[NpmAudit]: https://docs.npmjs.com/cli/audit
