#!/bin/bash

################################################################
# This script checks to see if a changelog file has been added
# with any "skipCI-*" fields set to true.
#
# This script supports the following fields:
#   skipCI-kube-tests
#   skipCI-storybook-tests
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

skipStorybookTestsDirective="skipCI-storybook-tests:true"
shouldSkipStorybookTests=false

githubBaseRef=$1
# If `githubBaseRef` is not present, it means that this script is not running as part of a PR (probably running on a push to main or an LTS branch).
# In that case we ignore the skip directives since we need to run CI
if [ ! -z "$githubBaseRef" ]; then
    # If there is no changelog found, the grep command fails and in turn the entire script exits since the error on exit flag has been set
    # To avoid that, we are using `|| true` to ensure that even if there is no changelog, it doesn't exit
    changelog=$(git diff origin/$githubBaseRef HEAD --name-only | grep "changelog/" || true)
    # An empty string is also one line in bash. Hence adding the first check
    if [ ! -z "$changelog" ] && [[ $(echo $changelog | wc -l | tr -d ' ') = "1" ]]; then
        echo "exactly one changelog added since main"
        echo "changelog file name == $changelog"
        if [[ $(cat $changelog | grep $skipKubeTestsDirective) ]]; then
            shouldSkipKubeTests=true
        fi
        if [[ $(cat $changelog | grep $skipStorybookTestsDirective) ]]; then
            shouldSkipStorybookTests=true
        fi
    else
        echo "no changelog found (or more than one changelog found) - not skipping CI"
    fi
fi

echo "skip-kube-tests=${shouldSkipKubeTests}" >> $GITHUB_OUTPUT
echo "skip-storybook-tests=${shouldSkipStorybookTests}" >> $GITHUB_OUTPUT