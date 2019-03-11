# huskyCI - Performing security tests inside your CI

<img src="https://raw.githubusercontent.com/wiki/globocom/huskyci/images/huskyCI-logo.png" align="center" height="" />
<!-- logo font: Anton -->

[![CircleCI](https://circleci.com/gh/globocom/huskyci/tree/master.svg?style=svg&circle-token=415bfb6b5aa0dfce8d2129878a66326da9533150)](https://circleci.com/gh/globocom/husky/tree/master)

huskyCI is an open source tool that performs security tests inside CI pipelines of multiple projects and centralizes all results into a database for further analysis and metrics.

The main goal of this project is to help development teams improve the quality of their code by finding vulnerabilities as soon as possible.

## How does it work?

Check out our [wiki as the huskyCI works](https://github.com/globocom/huskyci/wiki/how-works)!


## Cool! So huskyCI can check vulnerabilities in all languages ever?

Wow! Hold on! At this moment huskyCI can only perform static security analysis in Python ([Bandit][Bandit]), Ruby ([Brakeman][Brakeman]) and Golang ([Gosec][Gosec]). However, if you want to contribute to huskyCI by adding other cool security tests, you should check [this documentation](https://github.com/globocom/huskyci/wiki/how-add-new-security-tests) right away!

## Running locally

The easiest way to deploy huskyCI is by using Docker Compose, thus, you should have [Docker][Docker Install] and [Docker Compose][Docker Compose Install] installed on your machine. After cloning the repository, just run this to provision your local environment:

```
make install
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests to huskyCI.


## Documentation

Take a look at our [documentation](https://github.com/globocom/huskyci/wiki)!

## License

This project is licensed under the BSD 3-Clause "New" or "Revised" License - read [LICENSE.md](LICENSE.md) file for details.

[Docker Install]:  https://docs.docker.com/install/
[Docker Compose Install]: https://docs.docker.com/compose/install/
[huskyCI Client]: https://github.com/globocom/huskyci-client
[Bandit]: https://github.com/PyCQA/bandit
[Brakeman]: https://github.com/presidentbeef/brakeman
[Gosec]: https://github.com/securego/gosec
