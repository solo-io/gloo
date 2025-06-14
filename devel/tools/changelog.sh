#!/bin/bash

# Create a changelog file. This script automatically creates a file in the
# correct changelog directory. You can also pass the `-e` flag if you want to
# edit the file after creation. You can also use `-t` to specify which kind of
# changelog entry you would like to create (`-t FIX` or `-t DEPENDENCY_BUMP`
# for example). This script can be called multiple times, and it will
# automatically append additional changelog entries to the changelog file in
# question. So you could run `devel/tools/changelog.sh -t FIX` and then after
# that you could run `devel/tools/changelog.sh -t DEPENDENCY_BUMP` and both of
# those entries would appear in the file.

set -euo pipefail

CURRENT_VERSION=$(git describe --tags --always --abbrev=0)
OLD_VERSION=$(echo -n $CURRENT_VERSION | sed -E 's/.*[^0-9]([0-9]+)$/\1/')
NEW_VERSION=$((OLD_VERSION + 1))
NEXT_VERSION=$(echo -n $CURRENT_VERSION | sed -E "s/$OLD_VERSION$/$NEW_VERSION/")
BRANCH_NAME=$(git symbolic-ref -q HEAD | sed 's#^.*/##')
DESCRIPTION=${DESCRIPTION:=""}

CHANGELOG_DIR="changelog/$NEXT_VERSION"
mkdir -p "$CHANGELOG_DIR"

CHANGELOG_FILE="$CHANGELOG_DIR/$BRANCH_NAME.yaml"

##################################################
# OPTION FLAGS
OPT_EDIT=false # should we open the changelog file in $EDITOR?
OPT_CHANGELOG_TYPE="" # NEW_FEATURE, DEPENDENCY_BUMP, ENVOY_DEPENDENCY_BUMP, NON_USER_FACING

##################################################
# OPTION PARSING
while [[ $# -gt 0 ]]; do
  case $1 in
    -e|--edit)
      OPT_EDIT=true
      shift
      ;;
    --no-edit)
      OPT_EDIT=false
      shift
      ;;
    -t|--type)
      OPT_CHANGELOG_TYPE="$2"
      shift
      shift
      ;;
  esac
done
##################################################
# if the file does not exist, create it by adding the initial "changelog": line
test -f "${CHANGELOG_FILE}" || echo "changelog:" >> "${CHANGELOG_FILE}"

##################################################
# process the -t flag to add the templated changelog entries
case "$OPT_CHANGELOG_TYPE" in
  FIX)
    echo "creating FIX changelog"
    # make sure to append to the file so we don't delete existing entries!!
    cat << EOF >> "${CHANGELOG_FILE}"
  - type: FIX
    issueLink:
    resolvesIssue: false
    description: >-
      "${DESCRIPTION}"
EOF
    ;;
  NEW_FEATURE)
    echo "creating NEW_FEATURE changelog"
    # make sure to append to the file so we don't delete existing entries!!
    cat << EOF >> "${CHANGELOG_FILE}"
  - type: NEW_FEATURE
    issueLink:
    resolvesIssue: false
    description: >-
      "${DESCRIPTION}"
EOF
    ;;
  DEPENDENCY_BUMP)
    echo "creating DEPENDENCY_BUMP changelog"
    # make sure to append to the file so we don't delete existing entries!!
    cat << EOF >> "${CHANGELOG_FILE}"
  - type: DEPENDENCY_BUMP
    issueLink:
    resolvesIssue: false
    dependencyOwner: # solo-io
    dependencyRepo: # gloo
    dependencyTag: # v0.7.0
EOF
    ;;
  ENVOY_DEPENDENCY_BUMP)
    echo "creating ENVOY DEPENDENCY_BUMP changelog"
    # make sure to append to the file so we don't delete existing entries!!
    cat << EOF >> "${CHANGELOG_FILE}"
  - type: DEPENDENCY_BUMP
    issueLink:
    resolvesIssue: false
    dependencyOwner: solo-io
    dependencyRepo: envoy-gloo
    dependencyTag: # 1.33.0-patch3
EOF
    ;;
esac

##################################################
# Edit the file if `-e` was provided
if $OPT_EDIT ; then
  exec "$EDITOR" "${CHANGELOG_FILE}"
fi
