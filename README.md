
<h1 align="center">
    <img src="Gloo-01.png" alt="Gloo" width="200" height="242">
  <br>
  Gloo Kubernetes Upstream Discovery Service
</h1>


<h4 align="center"></h4>
<BR>

Gloo is a function gateway built on top of the [Envoy Proxy](https://www.Envoyproxy.io). Gloo provides a unified entry point
for access to all services and serverless functions, translating from any interface spoken by a client to any interface
spoken by a backend. Gloo aggregates REST APIs and events calls from clients, "glueing" together services in-cluster, 
out of cluster, across clusters, along with any provider of serverless functions.

This Repo 
----
Here you will find the source files for the Kubernetes Upstream Discovery Service. The k8s upstream discovery service is a
discovery service that runs outside of Gloo, discovering Kubernetes services and creating Gloo upstreams from them. The 
discovery service is an add-on for Gloo to facilitate easier routing and-self service for Gloo clients.

See our [documentation](https://gloo.solo.io) for more information about the Gloo discovery services. 

Documentation
-----

Get started by reading our docs here: [https://gloo.solo.io/](https://gloo.solo.io/)

Community
-----
Join us on our slack channel: [https://slack.solo.io/](https://slack.solo.io/)
