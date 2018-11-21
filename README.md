# husky-client

HuskyCI Client is an open source tool that performs security tests inside CI pipelines of multiple projects by sending requests to [HuskyCI][HuskyCI] and gathering its security tests results.

The main goal of this project is to help development teams improve the quality of their code by finding vulnerabilities as soon as possible.

## What is this client all about?

Well, actually it is just a binary built in Golang that performs the proper requests to HuskyCI API, waits security tests finish and inteprets the results by returning errors (if vulnerabilities are found) or not:

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

## Environment variables needed (use your own configuration):

```sh
echo 'export HUSKYCI_REPO_URL="https://github.com/tsuru/cst.git"' > .env
echo 'export HUSKYCI_REPO_BRANCH="master"' >> .env
echo 'export HUSKYCI_API="http://127.0.0.1:8888"' >> .env
```

[HuskyCI]: https://github.com/globocom/husky
