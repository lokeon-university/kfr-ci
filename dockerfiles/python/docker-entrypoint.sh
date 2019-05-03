#!/bin/bash
trap 'exit' ERR

echo "<h3>Starting the build</h3>"

echo "<h3>Checkout source code<h/3>"
git clone $REPO_URL $REPO_NAME
cd $REPO_NAME
git checkout $REPO_BRANCH
cd .
echo ""

# comprobamos si existe kfr.json
KFR_CONFIG_PRESENT=false
KFR_CONFIG_FILE=./.kfr.json

REQ_PRESENT=true
REQ_FILE=./requirements.txt

if [ -r ./.kfr.json ]; then
    echo "Found .kfr.json file." 1>&2
    KFR_CONFIG_PRESENT=true
fi

echo "<h3>Dependencies</h3>"

if ![ -r ./requirements.txt]; then
    echo "ERROR. FALTA 'requirements.txt' " 1>&2
    exit 1
fi

echo "<h3>Build</h3>"

if ($KFR_CONFIG_PRESENT && $REQ_PRESENT); then
    cat .kfr-ci.json | jq -r '. | .build[]' | bash
fi

exec "$@"