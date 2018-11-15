# HuskyCI

[![CircleCI](https://circleci.com/gh/globocom/husky/tree/master.svg?style=svg&circle-token=415bfb6b5aa0dfce8d2129878a66326da9533150)](https://circleci.com/gh/globocom/husky/tree/master)

HuskyCI is an open source tool that performs security tests inside CI pipelines of multiple projects and centralizes all results into a database for further analysis and metrics. 

The main goal of this project is to help development teams improve the quality of their code by finding vulnerabilities as soon as possible.  

## How does it work?

Imagine that an organization has projects like `awesome-golang-project`, `awesome-python-project` and `awesome-ruby-project`. In each CI configuration file, the following example code may be included:

```
test-project:
  stage: HuskyCI
  script:
    - wget urlTo.huskyCI/huskyci-client
    - chmod +x huskyci-client
    - ./huskyci-client
```

By adding this simple stage on each project, requests will be made to HuskyCI API and it will start analyzing new code submitted via Pull Request using well-known open source static analysis tools, as shown in the example bellow: 

![architecture](images/arch-example-huskyCI.png)

## What is this HuskyCI Client all about?

Well, actually [HuskyCI Client][HuskyCI Client] is just a binary built in Golang that performs the proper requests to HuskyCI API, waits security tests finish and inteprets the results by returning errors (if vulnerabilities are found) or not: 

```
$ ./huskyci-client

[HUSKYCI][*] new-feature-branch -> https://url.to.repository/team/project.git
[HUSKYCI][*] HuskyCI analysis started! UUT3MoVnLio9r5syzhbOIZYdLqbx4EDT

[HUSKYCI][!] Severity: MEDIUM
[HUSKYCI][!] Confidence: HIGH
[HUSKYCI][!] Details: Potential file inclusion via variable
[HUSKYCI][!] File: /go/src/code/example/example/example.go
[HUSKYCI][!] Line: 76
[HUSKYCI][!] Code: os.Open(path)

[HUSKYCI][!] Severity: LOW
[HUSKYCI][!] Confidence: HIGH
[HUSKYCI][!] Details: Errors unhandled.
[HUSKYCI][!] File: /go/src/code/example2/example2/example2.go
[HUSKYCI][!] Line: 132
[HUSKYCI][!] Code: subdirs, _ := ioutil.ReadDir(p)

[HUSKYCI][X] :(
exit status 1
```

## Cool! So HuskyCI can check vulnerabilities in all languages ever?

Wow! Hold on! At this moment, HuskyCI can only perform static security analysis in Python ([Bandit][Bandit]), Ruby ([Brakeman][Brakeman]) and Golang ([Gosec][Gosec]). However, if you want to contribute to HuskyCI by adding other cool security tests, you should check this documentation right away! 

## Running locally
 
The easiest way to deploy HuskyCI is by using Docker Compose, thus, you should have installed [Docker][Docker Install] and [Docker Compose][Docker Compose Install] on your machine. After cloning the repository, just run this to provision your local environment:

```
make install
```

#### Starting a new analysis:

```
curl -s -H "Content-Type: application/json" -d '{"repositoryURL":"https://github.com/tsuru/cst.git","repositoryBranch":"master"}' http://localhost:8888/husky
```

```
{"RID":"8L85jTJgtuN7o7pRi3sUQ3R4KuCjRcP9","details":"Request received.","result":"ok"}
```

#### Checking analysis status:

```
curl -s localhost:8888/husky/1HQfkskK69LYvLV7rWY03xv03YWoD47T
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests to HuskyCI.


[Docker Install]:  https://docs.docker.com/install/
[Docker Compose Install]: https://docs.docker.com/compose/install/
[HuskyCI Client]: https://github.com/globocom/husky-client
[Bandit]: https://github.com/PyCQA/bandit
[Brakeman]: https://github.com/presidentbeef/brakeman
[Gosec]: https://github.com/securego/gosec
