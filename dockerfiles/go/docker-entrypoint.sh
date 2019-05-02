#!/bin/bash
trap 'exit' ERR

source /etc/profile

echo "<h3>Starting the build</h3>"

#echo "<h3>Adding SSH keys</h3>"
#mkdir -p /root/.ssh/ && cp -R .ssh/* "$_"
#chmod 600 /root/.ssh/* &&
#    ssh-keyscan github.com >/root/.ssh/known_hosts
#echo

echo "<h3>Checkout source code<h/3>"
git clone $REPO_URL $REPO_NAME
cd $REPO_NAME
git checkout $REPO_BRANCH
echo

# comprobamos si existe kfr.json
KFR_CONFIG_PRESENT=false
KFR_CONFIG_FILE=./kfr.json

if [ -r ./kfr.json ]; then
    KFR_CONFIG_PRESENT=true
fi

echo "<h3>Dependencies</h3>"

if $KFR_CONFIG_PRESENT; then
    source <(cat .kfr-ci.json | jq -r '. | .env[]')
fi

go get -v ./...

echo "<h3>Build</h3>"
if $KFR_CONFIG_PRESENT; then
    cat .kfr-ci.json | jq -r '. | .build[]' | bash
fi

if !($KFR_CONFIG_PRESENT); then
    go build
fi

exec "$@"
