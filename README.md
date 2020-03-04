<p align="center">
  <img src="https://raw.githubusercontent.com/wiki/globocom/huskyCI/images/huskyCI-logo.png" align="center" height="" />
  <!-- logo font: Anton -->
</p>

<p align="center">
  <a href="https://github.com/globocom/huskyCI/releases"><img src="https://img.shields.io/github/v/release/globocom/huskyCI"/></a>
  <a href="https://circleci.com/gh/globocom/huskyCI/tree/master"><img src="https://img.shields.io/circleci/build/github/globocom/huskyCI/master?token=415bfb6b5aa0dfce8d2129878a66326da9533150"/></a>
  <a href="https://github.com/rafaveira3/writing-and-presentations/blob/master/DEFCON-27-APP-SEC-VILLAGE-Rafael-Santos-huskyCI-Finding-security-flaws-in-CI-before-deploying-them.pdf"><img src="https://img.shields.io/badge/DEFCON%2027-AppSec%20Village-black"/></a>
<a href="https://github.com/rafaveira3/contributions/blob/master/huskyCI-BlackHat-Europe-2019.pdf"><img src="https://img.shields.io/badge/Black%20Hat%20Europe%202019-Arsenal-black"/></a>
</p>

## Introduction

huskyCI is an open source tool that orchestrates security tests and centralizes all results into a database for further analysis and metrics. It can perform static security analysis in Python ([Bandit][Bandit] and [Safety][Safety]), Ruby ([Brakeman][Brakeman]), JavaScript ([Npm Audit][NpmAudit] and [Yarn Audit][YarnAudit]), Golang ([Gosec][Gosec]), and Java ([SpotBugs][SpotBugs] plus [Find Sec Bugs][FindSec]). It can also audit repositories for secrets like AWS Secret Keys, Private SSH Keys, and many others using [GitLeaks][Gitleaks].

## How does it work?

Developers can set up a new stage into their CI pipelines to check for vulnerabilities:

<p align="center"><img src="huskyCI-stage.png"/></p>

If security issues are found in the code, the severity, the confidence, the file, the line, and many more useful information can be shown, as exemplified: 

```
[HUSKYCI][*] poc-python-bandit -> https://github.com/globocom/huskyCI.git
[HUSKYCI][*] huskyCI analysis started! yDS9tb9mdt4QnnyvOBp3eVAXE1nWpTRQ

[HUSKYCI][!] Language: Python
[HUSKYCI][!] Tool: Bandit
[HUSKYCI][!] Severity: MEDIUM
[HUSKYCI][!] Confidence: HIGH
[HUSKYCI][!] Details: Use of exec detected.
[HUSKYCI][!] File: ./main.py
[HUSKYCI][!] Line: 7
[HUSKYCI][!] Code:
6 
7 exec(command)
8 

[HUSKYCI][SUMMARY] Python -> huskyci/bandit:1.6.2
[HUSKYCI][SUMMARY] High: 0
[HUSKYCI][SUMMARY] Medium: 1
[HUSKYCI][SUMMARY] Low: 0
[HUSKYCI][SUMMARY] NoSecHusky: 0

[HUSKYCI][SUMMARY] Total
[HUSKYCI][SUMMARY] High: 0
[HUSKYCI][SUMMARY] Medium: 1
[HUSKYCI][SUMMARY] Low: 0
[HUSKYCI][SUMMARY] NoSecHusky: 0

[HUSKYCI][*] Some HIGH/MEDIUM issues were found :(
ERROR: Job failed: exit code 1
```

## Getting Started

You can try huskyCI by setting up a local environment using Docker Compose following [this guide](http://husky.ci/docs/development/set-up-environment).

## Documentation

All guides and the full documentation can be found in the [official documentation page](http://husky.ci/docs/quickstart/overview).

## Contributing

Read our [contributing guide](https://github.com/globocom/huskyCI/blob/master/CONTRIBUTING.md) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes to huskyCI.

## Communication

We have a few channels for contact, feel free to reach out to us at:

- [GitHub Issues](https://github.com/globocom/huskyCI/issues)
- [Gitter](https://gitter.im/globocom/huskyCI)
- [Twitter](https://twitter.com/huskyCI)

## Contributors

This project exists thanks to all the [contributors]((https://github.com/globocom/huskyCI/graphs/contributors)). You rock!   ❤️🚀

## License

huskyCI is licensed under the [BSD 3-Clause "New" or "Revised" License](https://github.com/globocom/huskyCI/blob/master/LICENSE.md).

[Bandit]: https://github.com/PyCQA/bandit
[Safety]: https://github.com/pyupio/safety
[Brakeman]: https://github.com/presidentbeef/brakeman
[Gosec]: https://github.com/securego/gosec
[NpmAudit]: https://docs.npmjs.com/cli/audit
[YarnAudit]: https://yarnpkg.com/lang/en/docs/cli/audit/
[Gitleaks]: https://github.com/zricethezav/gitleaks
[SpotBugs]: https://spotbugs.github.io
[FindSec]: https://find-sec-bugs.github.io