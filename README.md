# huskyCI - Performing security tests inside your CI

<img src="https://raw.githubusercontent.com/wiki/globocom/huskyCI/images/huskyCI-logo.png" align="center" height="" />
<!-- logo font: Anton -->

[![CircleCI](https://circleci.com/gh/globocom/huskyCI/tree/master.svg?style=svg&circle-token=415bfb6b5aa0dfce8d2129878a66326da9533150)](https://circleci.com/gh/globocom/huskyCI/tree/master) [![Join the chat at https://gitter.im/globocom/huskyCI](https://badges.gitter.im/globocom/huskyCI.svg)](https://gitter.im/globocom/huskyCI?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

huskyCI is an open source tool that performs security tests inside CI pipelines of multiple projects and centralizes all results into a database for further analysis and metrics.


## How does it work?

The main goal of this project is to help development teams improve the quality of their code by finding vulnerabilities as soon as possible.

huskyCI can perform static security analysis in Python ([Bandit][Bandit] and [Safety][Safety]), Ruby ([Brakeman][Brakeman]), JavaScript ([RetireJS][RetireJS]) and Golang ([Gosec][Gosec]). You should check our [wiki](https://github.com/globocom/huskyCI/wiki/How-does-huskyCI-work%3F) to better understand how this tool could help securing your organization projects!


## Installing

The easiest way to deploy huskyCI locally is by using Docker Compose, thus, you should have [Docker][Docker Install] and [Docker Compose][Docker Compose Install] installed on your machine. After cloning this repository, run this:

```
make install
```

## Running

After installing, a `.env` file will be generated which is needed to run huskyCI-client: 

```sh
. .env 
```
```sh
make run-client
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