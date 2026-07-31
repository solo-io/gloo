---
title: Listener Configuration
weight: 20
---

**Gateway** definitions set up the protocols and ports on which Gloo Gateway listens for traffic.  For example, by default Gloo Gateway will have a gateway configured for HTTP and HTTPS traffic. Gloo Gateway allows you to configure properties of your gateways with several plugins.

These guides show you how to apply these advanced listener configurations to refine your gateways' behavior.

---

## Overview

For demonstration purposes, let's edit the default gateways that are installed with `glooctl install gateway`. You can list and edit gateways with `kubectl`.

```bash
kubectl get gateway --all-namespaces
```

```bash
NAMESPACE     NAME                AGE
gloo-system   gateway-proxy       2d
gloo-system   gateway-proxy-ssl   2d
```

```bash
kubectl edit gateway -n gloo-system gateway-proxy
```

### Plugin summary

The listener plugin portion of the gateway Custom Resource (CR) is shown below.

{{< highlight yaml "hl_lines=7-12" >}}
apiVersion: gateway.solo.io/v1
kind: Gateway
metadata: # collapsed for brevity
spec:
  bindAddress: '::'
  bindPort: 8080
  httpGateway:
    options:
      grpcWeb:
        disable: true
      httpConnectionManagerSettings:
        via: reference-string
  useProxyProto: false
status: # collapsed for brevity
{{< /highlight >}}


### Verify your configuration

To verify that your configuration was accepted, inspect the proxy.

```bash
glooctl get proxy -n gloo-system gateway-proxy -o yaml
```

{{< highlight yaml "hl_lines=6-10" >}}
...
  spec:
    listeners:
    - bindAddress: '::'
      bindPort: 8080
      httpListener:
        options:
          grpcWeb:
            disable: true
          httpConnectionManagerSettings:
            via: reference-string
        virtualHosts:
        - domains:
          - '*'
          name: gloo-system.merged-*
        - domains:
          - solo.io
          name: gloo-system.myvs3
      name: listener-::-8080
      useProxyProto: false
    - bindAddress: '::'
      bindPort: 8443
      httpListener: {}
      name: listener-::-8443
      useProxyProto: false
  status: # collapsed for brevity
kind: List
metadata: # collapsed for brevity
{{< /highlight >}}

---

## IPv4-only environments {#ipv4-only}

By default, Gloo Gateway binds listeners to the IPv6 wildcard address `::`, which listens for traffic on all interfaces. In IPv4-only Kubernetes environments, this causes gateways to fail to bind and proxies to receive an invalid xDS fallback configuration. Use the `gatewayProxies.<NAME>.gatewaySettings.ipv4Only` Helm value to switch the Gateway spec's `bindAddress` to `0.0.0.0` instead.

1. Get the current Helm values for your installation and save them to a file.
   ```sh
   helm get values gloo -n gloo-system -o yaml > gloo-gateway.yaml
   open gloo-gateway.yaml
   ```

2. Add the following value to your Helm values file. Replace `<NAME>` with the name of your gateway proxy, such as `gatewayProxy`.
   ```yaml
   gatewayProxies:
     <NAME>:
       gatewaySettings:
         ipv4Only: true
   ```

3. Upgrade your Gloo Gateway installation with the updated values.
   ```sh
   helm upgrade -n gloo-system gloo gloo/gloo \
     -f gloo-gateway.yaml
   ```

---

## Next Steps

The guide above showed some basics on how to manipulate settings for the gateway listener in Gloo Gateway. The following guides provide additional detail around what you may want to change in the gateway CR and how to do so.

{{% children description="true" depth="1" %}}
