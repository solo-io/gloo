Storage Client
----

This repository contains the sources for the Gloo Storage Client. The Gloo Storage Client is the centerpiece for the Gloo
API. Both clients to Gloo (including the Gloo discovery services) and Gloo itself consume this library to interact with 
the storage layer, the universal source of truth in Gloo's world.

Developers of integrations should use this repository as a client library for their application (if it's written in Go).
Support for more languages is in our roadmap.  

For information about storage in Gloo and writing client integrations, see our [documentation](https://gloo.solo.io). 
