# husky-client
Husky client to be downloaded and executed inside a CI.

#### Environment variables required (use your own configuration):

```sh
echo 'export HUSKYCI_REPO_URL="git@github.com:globocom/husky-client.git"' > .env
echo 'export HUSKYCI_REPO_BRANCH="master"' >> .env
echo 'export HUSKYCI_API="http://127.0.0.1:4444"' >> .env
```

