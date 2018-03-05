<h1 align="center">
    <img src="Gloo-01.png" alt="Gloo Storage Client" width="200" height="242">
  <br>
</h1>


<h3 align="center">Gloo Storage Client</h3>
<BR>

Gloo is a function gateway built on top of the [Envoy Proxy](https://www.Envoyproxy.io). Gloo provides a unified entry point
for access to all services and serverless functions, translating from any interface spoken by a client to any interface
spoken by a backend. Gloo aggregates REST APIs and events calls from clients, "glueing" together services in-cluster, 
out of cluster, across clusters, along with any provider of serverless functions.

This Repo 
----

This repository contains the sources for the Gloo Storage Client. The Gloo Storage Client is the centerpiece for the Gloo
API. Both clients to Gloo (including the Gloo discovery services) and Gloo itself consume this library to interact with 
the storage layer, the universal source of truth in Gloo's world.

Developers of integrations should use this repository as a client library for their application (if it's written in Go).
Support for more languages is in our roadmap.  

For information about storage in Gloo and writing client integrations, see our [documentation](https://gloo.solo.io). 

Documentation
-----

Get started by reading our docs here: [https://gloo.solo.io/](https://gloo.solo.io/)

Community
-----
Join us on our slack channel: [https://slack.solo.io/](https://slack.solo.io/)
