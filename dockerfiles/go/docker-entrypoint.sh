#!/bin/bash
trap 'exit' ERR

source /etc/profile

echo "<h3>Starting the build</h3>"

echo "<h3>Adding SSH keys</h3>"
mkdir -p /root/.ssh/ && cp -R .ssh/* "$_"
chmod 600 /root/.ssh/* && \
    ssh-keyscan github.com > /root/.ssh/known_hosts 
echo ""

echo "<h3>Checkout source code<h/3>"
git clone $REPO_URL $REPO_NAME
cd $REPO_NAME
git checkout $REPO_BRANCH
cd .
echo

# comprobamos si existe kfr.yml
KFR_CONFIG_PRESENT = false
KFR_CONFIG_FILE = ./kfr.yml

if [ -r ./kfr.yml ]; then
    KFR_CONFIG_PRESENT = true
fi

echo "<h3>Dependencies</h3>"

go get -v ./...

echo "<h3>Setup</h3>"
if ! ($KFR_CONFIG_PRESENT ); then
    go build 
fi

exec "$@"