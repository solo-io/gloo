# Resources to manually test/demo LDAP plugin functionality

- `docker-run-dlap.sh`: idempotent script to run OpenLDAP via docker with the config in `/ldif`
- `setup-ldap.sh`: script to get the same OpenLDAP server running in Kubernetes:
    ```bash
    ./kube/setup-ldap.sh <optional_namespace_argument>
    ```
- `ldif`: contains [.ldif](https://en.wikipedia.org/wiki/LDAP_Data_Interchange_Format) files to initialize the LDAP server with