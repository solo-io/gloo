
<h1 align="center">
    <img src="Gloo-01.png" alt="Gloo" width="200" height="242">
  <br>
</h1>


<h3 align="center">Gloo Ingress Controller</h3>
<BR>

Gloo is a function gateway built on top of the [Envoy Proxy](https://www.Envoyproxy.io). Gloo provides a unified entry point
for access to all services and serverless functions, translating from any interface spoken by a client to any interface
spoken by a backend. Gloo aggregates REST APIs and events calls from clients, "glueing" together services in-cluster, 
out of cluster, across clusters, along with any provider of serverless functions.

This Repo 
----
Here you will find the source files for Gloo's Kubernetes Ingress Controller. The Kubernetes Ingress Controller
extends the Gloo configuration language to support the Kubernetes Ingress language, enabling Gloo (and Envoy) to 
monitor and apply [Kubernetes Ingress rules](https://kubernetes.io/docs/concepts/services-networking/ingress/).

For more understanding about how the ingress controller works, head over to
our [documentation](https://gloo.solo.io).

Documentation
-----

Get started by reading our docs here: [https://gloo.solo.io/](https://gloo.solo.io/)

Community
-----
Join us on our slack channel: [https://slack.solo.io/](https://slack.solo.io/)
