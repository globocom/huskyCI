<p align="center">
  <img src="https://raw.githubusercontent.com/wiki/globocom/huskyCI/images/huskyCI-logo.png" align="center" height="" />
  <!-- logo font: Anton -->
</p>

<p align="center">
  <a href="https://github.com/globocom/huskyCI/releases"><img src="https://img.shields.io/github/v/release/globocom/huskyCI"/></a>
  <a href="https://github.com/rafaveira3/writing-and-presentations/blob/master/DEFCON-27-APP-SEC-VILLAGE-Rafael-Santos-huskyCI-Finding-security-flaws-in-CI-before-deploying-them.pdf"><img src="https://img.shields.io/badge/DEFCON%2027-AppSec%20Village-black"/></a>
<a href="https://github.com/rafaveira3/contributions/blob/master/huskyCI-BlackHat-Europe-2019.pdf"><img src="https://img.shields.io/badge/Black%20Hat%20Europe%202019-Arsenal-black"/></a>
<a href="https://defectdojo.readthedocs.io/en/latest/integrations.html#huskyci-report"><img src="https://img.shields.io/badge/DefectDojo-Compatible-brightgreen"/></a>
</p>

_This article can also be read in [Brazilian Portuguese](README-ptBR.md)._

## Introduction

huskyCI is an open source tool that orchestrates security tests and centralizes all results into a database for further analysis and metrics. It can perform static security analysis in Python ([Bandit][Bandit] and [Safety][Safety]), Ruby ([Brakeman][Brakeman]), JavaScript ([Npm Audit][NpmAudit] and [Yarn Audit][YarnAudit]), Golang ([Gosec][Gosec]), Java ([SpotBugs][SpotBugs] plus [Find Sec Bugs][FindSec]), and HCL ([TFSec][TFSec]). It can also audit repositories for secrets like AWS Secret Keys, Private SSH Keys, and many others using [GitLeaks][Gitleaks].

## How does it work?

Developers can set up a new stage into their CI pipelines to check for vulnerabilities:

<p align="center"><img src="huskyCI-stage.png"/></p>

If security issues are found in the code, the severity, the confidence, the file, the line, and many more useful information can be shown, as exemplified:

```
[HUSKYCI][*] poc-python-bandit -> https://github.com/globocom/huskyCI.git
[HUSKYCI][*] huskyCI analysis started! yDS9tb9mdt4QnnyvOBp3eVAXE1nWpTRQ

[HUSKYCI][!] Title: Use of exec detected.
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

[HUSKYCI][!] Title: Possible hardcoded password: 'password123!'
[HUSKYCI][!] Language: Python
[HUSKYCI][!] Tool: Bandit
[HUSKYCI][!] Severity: LOW
[HUSKYCI][!] Confidence: MEDIUM
[HUSKYCI][!] Details: Possible hardcoded password: 'password123!'
[HUSKYCI][!] File: ./main.py
[HUSKYCI][!] Line: 1
[HUSKYCI][!] Code:
1 secret = 'password123!'
2
3 password = 'thisisnotapassword' #nohusky
4

[HUSKYCI][SUMMARY] Python -> huskyci/bandit:1.6.2
[HUSKYCI][SUMMARY] High: 0
[HUSKYCI][SUMMARY] Medium: 1
[HUSKYCI][SUMMARY] Low: 1
[HUSKYCI][SUMMARY] NoSecHusky: 1

[HUSKYCI][SUMMARY] Total
[HUSKYCI][SUMMARY] High: 0
[HUSKYCI][SUMMARY] Medium: 1
[HUSKYCI][SUMMARY] Low: 1
[HUSKYCI][SUMMARY] NoSecHusky: 1

[HUSKYCI][*] The following securityTests were executed and no blocking vulnerabilities were found:
[HUSKYCI][*] [huskyci/gitleaks:2.1.0]
[HUSKYCI][*] Some HIGH/MEDIUM issues were found in these securityTests:
[HUSKYCI][*] [huskyci/bandit:1.6.2]
ERROR: Job failed: exit code 190
```

## Getting Started

You can try huskyCI by setting up a local environment using Docker Compose following [this guide](https://huskyci.opensource.globo.com/docs/development/set-up-environment).

## Documentation

All guides and the full documentation can be found in the [official documentation page](https://huskyci.opensource.globo.com/docs/quickstart/overview).

## Contributing

Read our [contributing guide](https://github.com/globocom/huskyCI/blob/master/CONTRIBUTING.md) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes to huskyCI.

## Troubleshooting

If you encounter a problem, please reach out to us at [GitHub Issues](https://github.com/globocom/huskyCI/issues) or [Gitter](https://gitter.im/globocom/huskyCI). And if you have found a solution for a commom problem, please add the problem and the solution Hher

## Communication

We have a few channels for contact, feel free to reach out to us at:

- [GitHub Issues](https://github.com/globocom/huskyCI/issues)
- [Gitter](https://gitter.im/globocom/huskyCI)
- [Twitter](https://twitter.com/huskyCI)

## Contributors

<!-- CONTRIBUTORS_START -->
<table><tr><td align="center"><a href="https://github.com/rafaveira3"><img src="https://avatars.githubusercontent.com/u/8943477?v=4" width="100" style="border-radius: 50%;"><br>rafaveira3</a></td><td align="center"><a href="https://github.com/Krlier"><img src="https://avatars.githubusercontent.com/u/40367872?v=4" width="100" style="border-radius: 50%;"><br>Krlier</a></td><td align="center"><a href="https://github.com/carlosljr"><img src="https://avatars.githubusercontent.com/u/25513224?v=4" width="100" style="border-radius: 50%;"><br>carlosljr</a></td><td align="center"><a href="https://github.com/spimpaov"><img src="https://avatars.githubusercontent.com/u/22274988?v=4" width="100" style="border-radius: 50%;"><br>spimpaov</a></td><td align="center"><a href="https://github.com/joserenatosilva"><img src="https://avatars.githubusercontent.com/u/11424945?v=4" width="100" style="border-radius: 50%;"><br>joserenatosilva</a></td></tr><tr><td align="center"><a href="https://github.com/janiltonmaciel"><img src="https://avatars.githubusercontent.com/u/987588?v=4" width="100" style="border-radius: 50%;"><br>janiltonmaciel</a></td><td align="center"><a href="https://github.com/gabriel-cantergiani"><img src="https://avatars.githubusercontent.com/u/27586618?v=4" width="100" style="border-radius: 50%;"><br>gabriel-cantergiani</a></td><td align="center"><a href="https://github.com/marcelomagina"><img src="https://avatars.githubusercontent.com/u/12450277?v=4" width="100" style="border-radius: 50%;"><br>marcelomagina</a></td><td align="center"><a href="https://github.com/nettoclaudio"><img src="https://avatars.githubusercontent.com/u/7503687?v=4" width="100" style="border-radius: 50%;"><br>nettoclaudio</a></td><td align="center"><a href="https://github.com/edersonbrilhante"><img src="https://avatars.githubusercontent.com/u/1094995?v=4" width="100" style="border-radius: 50%;"><br>edersonbrilhante</a></td></tr><tr><td align="center"><a href="https://github.com/GabhenDM"><img src="https://avatars.githubusercontent.com/u/38007503?v=4" width="100" style="border-radius: 50%;"><br>GabhenDM</a></td><td align="center"><a href="https://github.com/mdjunior"><img src="https://avatars.githubusercontent.com/u/3290669?v=4" width="100" style="border-radius: 50%;"><br>mdjunior</a></td><td align="center"><a href="https://github.com/gustavocovas"><img src="https://avatars.githubusercontent.com/u/11429002?v=4" width="100" style="border-radius: 50%;"><br>gustavocovas</a></td><td align="center"><a href="https://github.com/rodrigo-brito"><img src="https://avatars.githubusercontent.com/u/7620947?v=4" width="100" style="border-radius: 50%;"><br>rodrigo-brito</a></td><td align="center"><a href="https://github.com/vitoriario"><img src="https://avatars.githubusercontent.com/u/17754098?v=4" width="100" style="border-radius: 50%;"><br>vitoriario</a></td></tr><tr><td align="center"><a href="https://github.com/fguisso"><img src="https://avatars.githubusercontent.com/u/5755568?v=4" width="100" style="border-radius: 50%;"><br>fguisso</a></td><td align="center"><a href="https://github.com/abzcoding"><img src="https://avatars.githubusercontent.com/u/10992695?v=4" width="100" style="border-radius: 50%;"><br>abzcoding</a></td><td align="center"><a href="https://github.com/lzakharov"><img src="https://avatars.githubusercontent.com/u/26368218?v=4" width="100" style="border-radius: 50%;"><br>lzakharov</a></td><td align="center"><a href="https://github.com/arnaudlanna"><img src="https://avatars.githubusercontent.com/u/11250299?v=4" width="100" style="border-radius: 50%;"><br>arnaudlanna</a></td><td align="center"><a href="https://github.com/brenol"><img src="https://avatars.githubusercontent.com/u/7342697?v=4" width="100" style="border-radius: 50%;"><br>brenol</a></td></tr><tr><td align="center"><a href="https://github.com/chinchila"><img src="https://avatars.githubusercontent.com/u/3947315?v=4" width="100" style="border-radius: 50%;"><br>chinchila</a></td><td align="center"><a href="https://github.com/LuanaGP"><img src="https://avatars.githubusercontent.com/u/34948516?v=4" width="100" style="border-radius: 50%;"><br>LuanaGP</a></td><td align="center"><a href="https://github.com/Lucbm99"><img src="https://avatars.githubusercontent.com/u/45500959?v=4" width="100" style="border-radius: 50%;"><br>Lucbm99</a></td><td align="center"><a href="https://github.com/mportela"><img src="https://avatars.githubusercontent.com/u/1304652?v=4" width="100" style="border-radius: 50%;"><br>mportela</a></td><td align="center"><a href="https://github.com/marcelometal"><img src="https://avatars.githubusercontent.com/u/665903?v=4" width="100" style="border-radius: 50%;"><br>marcelometal</a></td></tr><tr><td align="center"><a href="https://github.com/aranhams"><img src="https://avatars.githubusercontent.com/u/18319426?v=4" width="100" style="border-radius: 50%;"><br>aranhams</a></td><td align="center"><a href="https://github.com/ragoso"><img src="https://avatars.githubusercontent.com/u/9319775?v=4" width="100" style="border-radius: 50%;"><br>ragoso</a></td><td align="center"><a href="https://github.com/rafaelsq"><img src="https://avatars.githubusercontent.com/u/1598854?v=4" width="100" style="border-radius: 50%;"><br>rafaelsq</a></td><td align="center"><a href="https://github.com/rafaelrubbioli"><img src="https://avatars.githubusercontent.com/u/15738138?v=4" width="100" style="border-radius: 50%;"><br>rafaelrubbioli</a></td><td align="center"><a href="https://github.com/renatoaquino"><img src="https://avatars.githubusercontent.com/u/516453?v=4" width="100" style="border-radius: 50%;"><br>renatoaquino</a></td></tr></table>
<!-- CONTRIBUTORS_END -->

This project exists thanks to all the [contributors](<(https://github.com/globocom/huskyCI/graphs/contributors)>). You rock! ❤️🚀

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
[TFSec]: https://github.com/liamg/tfsec
