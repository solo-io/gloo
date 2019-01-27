# running locally on minikube

to update containers locally, you'll want to do the following:

1) `export SOLO_PROJECTS_DIR=$GOPATH/src/github.com/solo-io/solo-projects`
2) run the `restart.sh` script in this directory.
3) run `kubectl get pods -n gloo-system` and check that all pods are running

This script will rebuild all container simultaneously and then create a 
new template with the tag to use all of the local images.
Very useful for locally testing with all changes.