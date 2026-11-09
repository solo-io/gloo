#!/usr/bin/env bash
# Waits for a Cloud Build in project solo-public to reach a terminal state.
# Usage: wait-for-build.sh <build-id>
# Exits 0 on SUCCESS, 1 otherwise. Appends a summary to $GITHUB_STEP_SUMMARY when set.
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 <build-id>" >&2
  exit 2
fi
build_id="$1"

echo "Cloud Build ${build_id}. Cancelling this run does not stop the build. To stop it: gcloud builds cancel ${build_id} --project solo-public"

describe() {
  local i out
  for i in 1 2 3 4 5; do
    if out="$(gcloud builds describe "$1" --project solo-public --format="${2:-value(status)}")"; then
      printf '%s\n' "${out}"
      return 0
    fi
    echo "gcloud builds describe failed (attempt ${i}/5)" >&2
    if [[ "${i}" -lt 5 ]]; then sleep 10; fi
  done
  return 1
}

if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
  {
    echo "Cloud Build \`${build_id}\`: $(describe "${build_id}" 'value(logUrl)')"
    echo
    echo "Cancelling this run does not stop the build. To stop it: \`gcloud builds cancel ${build_id} --project solo-public\`"
  } >> "${GITHUB_STEP_SUMMARY}"
fi

while true; do
  status="$(describe "${build_id}")"
  echo "status: ${status}"
  case "${status}" in
    SUCCESS|FAILURE|INTERNAL_ERROR|TIMEOUT|CANCELLED|EXPIRED) break ;;
  esac
  sleep 30
done
echo "Logs: $(describe "${build_id}" 'value(logUrl)')"
if [[ "${status}" != "SUCCESS" ]]; then
  echo "::error::Cloud Build ${build_id} ended ${status}"
  exit 1
fi
