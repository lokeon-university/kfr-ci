#!/bin/bash
trap 'exit' ERR

echo "<h3>Starting the build</h3>"

echo "<h3>Checkout source code<h/3>"
git clone $REPO_URL $REPO_NAME
cd $REPO_NAME
git checkout $REPO_BRANCH
echo

# comprobamos si existe kfr-ci.json
KFR_CONFIG_PRESENT=false
KFR_CONFIG_FILE=./.kfr-ci.json

if [ -r $KFR_CONFIG_FILE ]; then
    KFR_CONFIG_PRESENT=true
fi

echo "<h3>Dependencies</h3>"

cat .kfr-ci.json | jq -r '. | .submodules != null'

if $KFR_CONFIG_PRESENT; then
    source <(cat .kfr-ci.json | jq -r '. | .env[]')
fi

echo "<h3>Build</h3>"

if $KFR_CONFIG_PRESENT; then
    cat .kfr-ci.json | jq -r '. | .build[]' | bash
fi

exec "$@"
