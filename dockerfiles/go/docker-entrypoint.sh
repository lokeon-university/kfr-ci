#!/bin/bash
trap 'exit' ERR

KFR_CONFIG_PRESENT="false"
KFR_CONFIG_FILE=./.kfr.json

echo "<h3>Checkout<h/3>"
git clone --progress "$REPO_URL" "$REPO_NAME"
cd "$REPO_NAME" || exit
git checkout --progress "$REPO_BRANCH"
cd . || exit
echo ""

if [ -r "$KFR_CONFIG_FILE" ]; then
    KFR_CONFIG_PRESENT="true"
else 
    echo "config file not found"
    exit 2 #file $KFR_CONFIG_FILE not found
fi

echo "<h3>Dependencies</h3>"
go get -v ./..
echo ""
echo "<h3>Build/Test</h3>"

LEGHT=$(jq -r '. | .steps | length' "$KFR_CONFIG_FILE")
if [ $KFR_CONFIG_PRESENT ] && [ "$LEGHT" -ne "0" ]; then
    jq -r '. | .steps[]' "$KFR_CONFIG_FILE" | bash
else 
    echo "steps cannot be empty"
    exit 4 #steps cannot be empty
fi

exec "$@"
