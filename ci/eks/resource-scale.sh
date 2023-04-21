#!/bin/bash
set -e
# get env data
output_directory=output/scale
mkdir -p $output_directory
PARENT_PATH=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
SCRIPTS_RELATIVE_PATH="/ci/eks"
SOLO_PROJECTS_PATH=${PARENT_PATH:0:(( ${#PARENT_PATH} - ${#SCRIPTS_RELATIVE_PATH} ))}

# setup output location
CREATION_TIMES_OUTPUT=${SOLO_PROJECTS_PATH}"/creation_times.csv"
DELETION_TIMES_OUTPUT=${SOLO_PROJECTS_PATH}"/deletion_times.csv"
echo "CREATION_TIMES_OUTPUT=${CREATION_TIMES_OUTPUT}" >> $GITHUB_ENV
echo "DELETION_TIMES_OUTPUT=${DELETION_TIMES_OUTPUT}" >> $GITHUB_ENV

# run template data
echo Start template applier
cd "${SOLO_PROJECTS_PATH}"/ci/eks/applier
go mod download
TEMPLATE_PATH_RESOURCE_SCALE="${SOLO_PROJECTS_PATH}""${SCRIPTS_RELATIVE_PATH}"/templates/resource_scale.yaml


# run tests at given scale
echo "$TEMPLATE_PATH"
kubectl config use-context mgmt-cluster
go run main.go apply -f "$TEMPLATE_PATH_RESOURCE_SCALE" --start 1 --iterations "$ITERATIONS" --async true --workers 20


# run creation
cd "${SOLO_PROJECTS_PATH}"/ci/eks/scalescripts

#set to correct folder
go run main.go scaleup -f "$CREATION_TIMES_OUTPUT" -r 11 --test-request true

# run deletion
go run main.go scaledown -f "$DELETION_TIMES_OUTPUT" -r 10