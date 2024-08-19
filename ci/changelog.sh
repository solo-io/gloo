#!/usr/bin/env bash

set -e

CURRENT_VERSION=$(git describe --tags --always --abbrev=0)
OLD_VERSION=$(echo -n $CURRENT_VERSION | sed -E 's/.*[^0-9]([0-9]+)$/\1/')
NEW_VERSION=$((OLD_VERSION + 1))
NEXT_VERSION=$(echo -n $CURRENT_VERSION | sed -E "s/$OLD_VERSION$/$NEW_VERSION/")
BRANCH_NAME=$(git symbolic-ref -q HEAD | sed 's#^.*/##')

export DESCRIPTION=${DESCRIPTION:=""}
export CHANGELOG_TYPE=${CHANGELOG_TYPE:-NON_USER_FACING}

CI_FILTER=${CI_FILTER:-} # Used to only run certain subset of tests

DOCUMENT='
changelog:
- type: ${CHANGELOG_TYPE}
  issueLink: ${ISSUE_LINK}
  description: >
    "${DESCRIPTION}"
  resolvesIssue: false
'

if [[ -z $ISSUE_LINK ]]; then
  DOCUMENT=$(envsubst <<< "$DOCUMENT" | yq -P 'del(.[][0].issueLink)')
fi

case $(echo $CI_FILTER| awk '{print tolower($0)}') in
  skip)
    DOCUMENT=$(envsubst <<< "$DOCUMENT" | yq -P '. | .[][0] += {"skipCI": true}')
  ;;
  ui)
    DOCUMENT=$(envsubst <<< "$DOCUMENT" | yq -P '. | .[][0] += {"onlyUITests": true}')
  ;;
  insights)
    DOCUMENT=$(envsubst <<< "$DOCUMENT" | yq -P '. | .[][0] += {"onlyInsightsTests": true}')
  ;;
esac

CHANGELOG_FILE=changelog/$NEXT_VERSION
if [[ $FILE != "" ]]; then
  CHANGELOG_FILE=$CHANGELOG_FILE/$FILE.yaml
else
  CHANGELOG_FILE=$CHANGELOG_FILE/$BRANCH_NAME.yaml
fi

mkdir -p "changelog/$NEXT_VERSION"
envsubst <<< "$DOCUMENT" | yq -P | sed '/^$/d' | tee $CHANGELOG_FILE &>/dev/null

echo Created $CHANGELOG_FILE
