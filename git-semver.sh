#! /bin/bash

# Converts git describe to valid SemVer2 format.
#
# This script acts as the single "source of truth" when specifying a version
# for locally-built helm charts, docker images, etc.
#
# ===== Helm Notes =====
# Valid SemVer2 formatted versions are required by Helm. If the version
# specified for a chart is already publicly released, helm will pull in the
# published chart rather than using the locally-built chart.
#
# ===== Docker Notes =====
# Docker image tags "must be valid ASCII and may contain lowercase and
# uppercase letters, digits, underscores, periods and dashes. A tag name may
# not start with a period or a dash and may contain a maximum of 128
# characters." Since '+' characters are not permitted in tags, we avoid
# using build metadata (https://semver.org/#spec-item-10) when formatting the
# version (with the -no-meta) to keep both helm and docker happy.

go install github.com/mdomke/git-semver/v6@v6.2.0
VERSION=$(git-semver -no-meta)
VERSION=${VERSION:1} # Remove leading 'v'
echo "$VERSION"