# AWS Lambda Plugin for Gloo


#### Description

The [AWS Lambda Plugin for Gloo](https://github.com/solo-io/gloo-plugins/tree/master/aws) brings the power of AWS Lambda
to the Envoy Proxy. This plugin allows HTTP requests routed by Envoy to become AWS Lambda invocations, automatically
signed by [AWS' API Signature signing process](https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html). 

This allows clients to make regular HTTP requests (including `GET` requests - even from a browser) through Gloo that
automatically route to AWS Lambda invocations. Clients are abstracted from the AWS API without having to specify
any AWS-specific headers or signature in their requests. `POST`s containing a JSON body will become the payload
for the Lambda function.

To jump right in, follow the AWS Lambda Getting Started Guide<!--(TODO)-->.


#### Upstream Spec Configuration

The **Upstream Type** for AWS Lambda upstreams is `aws`. 

The [upstream spec](../v1/upstream.md#v1.Upstream) for AWS Lambda Upstreams has the following structure:

```yaml
region: string
secret_ref: string
```
| Field | Type |  Description |
| ----- | ---- |  ----------- |
| region | string |  The [AWS Region](https://docs.aws.amazon.com/general/latest/gr/rande.html) this upstream points to. The credentials specified in the secret must be valid for this region. **required** |
| secret_ref | string |  The [AWS Credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) this upstream points to. The credentials specified in the secret must be valid for this region. **required** |

The `secret_ref` must be the name of a secret in Gloo's [secret storage backend](../introduction/concepts.md#secrets).



#### Secrets Configuration

The content of the secret must follow the format:

```yaml
access_key: <aws-access-key-id>
secret_key: <aws-secret-access-key>
```



#### Function Spec Configuration
The AWS Lambda Upstream Type supports [functions](../introduction/concepts.md#Functions) (and is in fact useless without them).
The [function spec](../v1/upstream.md#v1.Function) has the following stucture:


```yaml
function_name: string
qualifier: string
```

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| function_name | string |  | The [Lambda Function Name](https://docs.aws.amazon.com/lambda/latest/dg/API_Invoke.html#API_Invoke_RequestSyntax) used to invoke the function. **required** |
| qualifier | string |  | The [Lambda Function Qualifier (Version)](https://docs.aws.amazon.com/lambda/latest/dg/API_Invoke.html#API_Invoke_RequestSyntax) used to invoke the function. If empty, the `$LATEST` version will be used. |


#### Example Lambda Upstream

The following is an example of a valid Lambda [Upstream](../introduction/concepts.md#Upstreams):

```yaml
name: my-lambda-upstream
spec:
  region: "us-east-1"
  secret_ref: "my-aws-secrets"
type: aws
functions:
- name: func1
  spec:
    function_name: func1
- name: func2
  spec:
    function_name: func2
    qualifier: v1

```

#### Discovery

The Gloo Function Discovery Service<!--(TODO)--> will automatically discover functions for Lambda upstreams if it is running.
Simply create a Lambda upstream for Gloo to track, and the function discovery service will auomatically populate it with your
available lambda functions and keep it up to date with your AWS account. 
