#!/usr/bin/env bash

#############################################################################
# This script can be used to quickly spin up a docker image running OpenLDAP.
# Custom config is provided through the ENVs and the files in ./ldif which
# are loaded by the server on startup.
#############################################################################
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Stop and remove any existing openldap containers; suppress output
(docker stop openldap && docker rm openldap) >/dev/null 2>/dev/null || true

# Start new openldap instance
docker run --name openldap -p 389:389 -p 636:636 --detach \
  --env LDAP_ORGANISATION="Solo.io" \
  --env LDAP_DOMAIN="solo.io" \
  --env LDAP_ADMIN_PASSWORD="solopwd" \
  --volume "$DIR/ldif:/container/service/slapd/assets/config/bootstrap/ldif/custom" \
  osixia/openldap:1.2.5 --copy-service --loglevel debug