---
title: Submit a request
weight: 920
description: Solo customers can submit new requests and update existing ones by using the Support Portal, email (`support@solo.io`), or Slack. 
---
Solo customers can submit new requests and update existing ones by using the [Support Portal](#support-portal), email (`support@solo.io`), or [Slack](#slack).
## Use the Support Portal {#support-portal}

### Create a Support Portal account

1. Navigate to the [Support Portal](https://support.solo.io) and click **Sign up**.
2. On the next page, enter your full name and email address. You receive an email with instructions for how to create a password.
3. Follow the **Create password** link in that email. Check your Spam folder if you do not see it in your inbox.
4. You **MUST** add your phone number to your newly created profile. Failing to do so prevents your calls from being routed to a Solo Support Engineer when you call the Support Hotline to report Urgent incidents. 

### Submit a ticket

1. Navigate to the [Support Portal](https://support.solo.io). 
2. Log in, and ensure your phone number is added to your profile. This is an important step, as it allows our system to recognize your account when you call the Support Hotline for urgent incidents.
3. Click **Create a ticket**.
4. Fill out the form. Start by entering any email addresses that you want updates to be sent to while Solo Support works on the ticket. Add a subject related to the issue that you are reporting.
5. Add details for the issue in the **description** field. Refer to [Details to include in your support request](#ticket-details) to help Solo Support better understand your request and provide timely and accurate assistance. 
6. Enter the priority for this request. To learn more about priority levels and how to assign the right priority level to your ticket, see [Priority levels](#priority-levels). 
7. Optional: Enter a GitHub issue and upload any attachments. 
   {{< notice note >}}
   The ticketing system has an attachment file size limit of 50MB.
   {{< /notice >}}
8. Submit your request. A Solo Support Engineer will be in touch following the [Targeted times for initial response](#response-time).

To update an existing ticket, log in to your account, and click **View pending tickets**. You see a new page where you can select the ticket you want to update and add comments to.

## Details to include in your support request {#ticket-details}

### Environment
- Versions for Gloo Edge and Kubernetes
- Infrastructure provider, such as AWS or on-premises VMs
- Production, development, or other environment details

### Setup
- Installation method, such as Helm or `glooctl`
- Configuration used for installation, such as Helm values
- Details of the following Gloo Edge resources,
  - `Gateway`
  - `VirtualService`
  - `RouteTable`
  - `Upstream`
- If Gloo Edge is exposed by a Load Balancer, details of this configuration. If not, provide details about how Gloo Edge is configured to proxy client traffic.

<!--
- Management cluster details, such as whether it also runs workloads
- Number of workload clusters
- East-west gateway details and namespace, if you use Gloo Gateway with Gloo Mesh or your own Istio installation
- Federated trust details, such as self-signed certificates or certificates provided by a certificate management provider
- Snapshot of your Gloo management-agent relay connection
-->

### Issue
- Description of the issue
- If the issue is reproducible, steps to reproduce the issue
- Severity of the issue, such as urgent, high, normal, or low
- Impact of the issue, such as blocking an update, blocking a demo, data loss or the system is down
- Attach any relevant configuration/YAML files related to the issue, such as `Gateway`, `VirtualService`, `RouteTable` or `Upstream` resources

### Product-specific details

##### Control plane

  - Capture the output of the `glooctl check` command. Typically, the command output indicates any errors in the control plane components or associated resources. </br>
    An example output is shown below.
    ```
    Checking deployments... 1 Errors!
    Checking pods... 2 Errors!
    Checking upstreams... OK
    Checking upstream groups... OK
    Checking auth configs... OK
    Checking rate limit configs... OK
    Checking VirtualHostOptions... OK
    Checking RouteOptions... OK
    Checking secrets... OK
    Checking virtual services... OK
    Checking gateways... OK
    Checking proxies... Skipping due to an error in checking deployments
    Skipping due to an error in checking deployments
    Error: 5 errors occurred:
    * Deployment gloo in namespace gloo-system is not available! Message: Deployment does not have minimum availability.
    * Pod gloo-8ddc4ff4c-g4mnf in namespace gloo-system is not ready! Message: containers with unready status: [gloo]
    * Not all containers in pod gloo-8ddc4ff4c-g4mnf in namespace gloo-system are ready! Message: containers with unready status: [gloo]
    * proxy check was skipped due to an error in checking deployments
    * xds metrics check was skipped due to an error in checking deployment
    ```
  - Collect the logs from various control plane components, such as `gloo` by using the `debug` log level (if possible). To enable the `debug` log level, see [Debugging control plane]({{< versioned_link_path fromRoot="/operations/debugging_gloo#debug-control-plane" >}}).

##### Data plane

  - Capture the currently served xDS configuration with the `glooctl proxy served-config > served-config.yaml` command.
  - Get the configuration that is served in the `gateway-proxy` Envoy pod(s). For more information, see [Dumping Envoy configuration]({{< versioned_link_path fromRoot="/operations/debugging_gloo#dump-envoy-configuration" >}}).
  - Get the access log(s) for failed request from the `gateway-proxy` pod(s). If Access logging is not enabled, refer to [this guide]({{< versioned_link_path fromRoot="/guides/security/access_logging" >}}) to enable it.
  - If possible, collect the logs from the `gateway-proxy` Envoy pod(s) in `debug` log level for the failed request. For more information, see [Viewing Envoy logs]({{< versioned_link_path fromRoot="/operations/debugging_gloo#view-envoy-logs" >}}).

<!--
- **Gloo Mesh and Gloo Mesh Gateway**: 
  - Input and output snapshots
  - Control plane issues:
    - Logs, debug mode if possible, from the management server and relevant agent pods
  - Data plane issues:
    - Config_dump(s) from relevant Envoy(s)
    - Envoy access log(s) for failing requests from relevant Envoy(s)
    - If possible, debug logs from relevant Envoy(s) during a failing request


Istio:
- Installation method: 
  - Gloo IstioLifecycleManager and GatewayLifecycleManager resources
  - Helm
  - Kustomize
  - IstioOperator
-->

## Priority levels {#priority-levels}

{{% notice note %}}
For the latest list of priority levels and descriptions, see the [Technical Support Policy](https://legal.solo.io/#technical-support-policy).
{{% /notice %}}

|Priority level|Description|
|----------------|--|
|Urgent priority**|A problem that severely impacts your use of the software in a production environment (such as loss of production data or in which your production systems are not functioning). The situation halts your business operations, or your revenue or brand are impacted and no procedural workaround exists.|
|High priority|A problem where the production environment is operational but functionality is severely reduced. The situation is causing a high impact to portions of your business operations, or your revenue or brand are threatened and no procedural workaround exists.|
|Normal priority|A problem that involves partial, non-critical loss of use of the software in a production environment or development environment. For production environments, there is a medium-to-low impact on your business, but your business continues to function, including by using a procedural workaround. For development environments, where the situation is causing your project to no longer continue or migrate into production.|
|Low priority|A general usage question, reporting of a documentation error, or recommendation for a future product enhancement or modification. For production environments, there is low-to-no impact on your business or the performance or functionality of your system. For development environments, there is a medium-to-low impact on your business, but your business continues to function, including by using a procedural workaround.|

## Targeted times for initial response {#response-time}

{{% notice note %}}
For the latest list of targeted times for the initial response, see the [Technical Support Policy](https://legal.solo.io/#technical-support-policy).
{{% /notice %}}

When you contact Solo Support, you can choose a priority for your request:

|Priority Level| Standard Support Policy    |Enhanced Support Policy|
|--|----------------------------|--|
|Urgent**| 1 hour (24/7/365)          |15 minutes (24/7/365)|
|High|4 Business Hours Local Time |2 hours (24/7/365)|
|Normal|8 Business Hours Local Time|4 Business Hours Local Time|
|Low|24 Business Hours Local Time|12 Business Hours Local Time|

{{% notice warning %}}
**To report Urgent priority, production-related issues, you must contact Solo’s Support Hotline at `+1-601-476-5646`. 
{{% /notice %}} 

{{% notice note %}}
Solo Support reserves the right to adjust the priority you select if it does not align with the priorities documented above. 
{{% /notice %}}

## Join the solo.io Slack community {#slack}

Join the [solo.io Slack community](https://soloio.slack.com) to participate in conversations with other users and ask questions to the Solo team. 

Navigate through the channels to find topics relevant to you. We recommend `#gloo`<!--, `#gloo-enterprise`, `#service-mesh-hub` for specific product conversations-->, and `#general` for general discussion.

To migrate a discussion from Slack into an email-based support ticket, you can use the `/zendesk create_ticket` macro within the channel where the conversation started. Select **Support** as the assignee to ensure your request reaches the Solo Support team. Selecting other options will delay or prevent our team from seeing your ticket. 

If the macro hasn't been added to a channel, ask a Solo team member to assist you.

## Additional resources

Check out the following guides to start troubleshooting your environment and to collect the information to provide in your support request. 

- Downloadable version of our versioned [Technical Support Policy](https://legal.solo.io/#technical-support-policy)
- [Gloo Edge contribution guidelines]({{% versioned_link_path fromRoot="/contributing/" %}})
- [Gloo Edge community](https://github.com/solo-io/gloo/tree/main)



<!--
- [Communities of Practice repository on GitHub](https://github.com/solo-io/solo-cop)
- [Gloo Gateway troubleshooting guide](https://docs.solo.io/gloo-gateway/latest/troubleshooting/)
- [Gloo Mesh troubleshooting guide](https://docs.solo.io/gloo-mesh-enterprise/latest/troubleshooting/)
-->



