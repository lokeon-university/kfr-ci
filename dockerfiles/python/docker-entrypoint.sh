#!/bin/bash
trap 'exit' ERR

KFR_CONFIG_PRESENT="false"
KFR_CONFIG_FILE=./.kfr.json
REQ_FILE=./requirements.txt
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}Checkout${NC}"
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

echo -e "${GREEN}Dependencies${NC}"
pip install -v -r "$REQ_FILE"
echo ""
echo -e "${GREEN}Build/Test${NC}"

LEGHT=$(jq -r '. | .steps | length' "$KFR_CONFIG_FILE")
if [ $KFR_CONFIG_PRESENT ] && [ "$LEGHT" -ne "0" ]; then
    jq -r '. | .steps[]' "$KFR_CONFIG_FILE" | bash
else 
    echo "steps cannot be empty"
    exit 4 #steps cannot be empty
fi

exec "$@"