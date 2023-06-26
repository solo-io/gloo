#!/bin/bash

################################################################
# This script checks to see if a changelog file has been added
# with any "skipCI-*" fields set to true.
#
# This script supports the following fields:
#   skipCI-kube-tests
#
# For each field that we support, the script will output the value
# into the $GITHUB_OUTPUT variable, which then can be consumed
# by other steps in our CI pipeline
#
# This script is based on the Gloo OSS: /ci/github-actions/process-skip-directives/script.sh
# It does not contain support for skipping docs builds, since those do not exist in this project
################################################################

set -ex

skipKubeTestsDirective="skipCI-kube-tests:true"
shouldSkipKubeTests=false

githubBaseRef=$1

if [[ $(git diff origin/$githubBaseRef HEAD --name-only | grep "changelog/" | wc -l) = "1" ]]; then
    echo "exactly one changelog added since main"
    changelogFileName=$(git diff origin/main HEAD --name-only | grep "changelog/")
    echo "changelog file name == $changelogFileName"
    if [[ $(cat $changelogFileName | grep $skipKubeTestsDirective) ]]; then
        shouldSkipKubeTests=true
    fi
else
    echo "no changelog found (or more than one changelog found) - not skipping CI"
fi

echo "skip-kube-tests=${shouldSkipKubeTests}" >> $GITHUB_OUTPUT