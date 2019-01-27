#!/usr/bin/env bash
cd $SOLO_PROJECTS_DIR/hack/openshift
minishift delete
minishift start
minishift addon apply admin-user

# get the oc command
eval $(minishift oc-env)

oc new-project gloo-system

# external script for creds
# ask a teammate for the docker password and put this script in your ~/scripts/secret/ dir
# format is:
# oc create secret docker-registry solo-bot-docker-hub --docker-server=https://index.docker.io/v1/ --docker-username=soloiobot --docker-password=<password> --docker-email=bot@solo.io
~/scripts/secret/docker_credential.sh

# need to be an admin to use cluster role bindings
oc login -u system:admin

./apply.sh

./setup_new_minishift.sh

# after the pods are ready:
# kubectl port-forward deploy/apiserver-ui 8080 8082

# TODO - add way to re-hydrate state (maybe gitops will do this for us)
