#!/bin/bash
set -e

# watch upstreams
output_directory=output/scale
mkdir -p $output_directory
PARENT_PATH=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
SCRIPTS_RELATIVE_PATH="/ci/eks"
SOLO_PROJECTS_PATH=${PARENT_PATH:0:(( ${#PARENT_PATH} - ${#SCRIPTS_RELATIVE_PATH} ))}

# setup output location
OUTPUT_LOCATION="${SOLO_PROJECTS_PATH}"/"${output_directory}/raw.csv"
POST_PROCESS_LOCATION="${SOLO_PROJECTS_PATH}"/"${output_directory}/postprocessed.csv"
echo "POST_PROCESS_LOCATION=${POST_PROCESS_LOCATION}" >> $GITHUB_ENV

echo Start upsteam watch
python "${SOLO_PROJECTS_PATH}"/ci/eks/scalescripts/watch-upstreams.py cluster-1 "$OUTPUT_LOCATION" "$ITERATIONS" &
P1=$!

# create templates
echo Start template applier
cd "${SOLO_PROJECTS_PATH}"/ci/eks/applier
go mod download
TEMPLATE_PATH_10="${SOLO_PROJECTS_PATH}"/ci/eks/templates/remotes_1_to_10.yaml
TEMPLATE_PATH_20="${SOLO_PROJECTS_PATH}"/ci/eks/templates/remotes_11_to_20.yaml
echo "$TEMPLATE_PATH"
kubectl config use-context mgmt-cluster
go run main.go apply -f "$TEMPLATE_PATH_10" --start 1 --iterations "$ITERATIONS" &
# this helps avoid a name collision when creating templates
NEXT_NAME_START=$((ITERATIONS * 3 + 1 ))
go run main.go apply -f "$TEMPLATE_PATH_20" --start $NEXT_NAME_START --iterations "$ITERATIONS" &


#wait for watch to end as that will mean all templates were created
wait $P1
echo completed watch
# post process data
python "${SOLO_PROJECTS_PATH}"/ci/eks/scalescripts/postprocess-data.py "$OUTPUT_LOCATION" "$POST_PROCESS_LOCATION"
