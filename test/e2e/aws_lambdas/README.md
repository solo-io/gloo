# AWS Lambda functions used by `test/e2e/aws_test.go`

These are the Lambda functions deployed in the test AWS account that
`aws_test.go` exercises via `LogicalName`/`LambdaFunctionName`. Source lives
here since Envoy/Gloo only knows the function name — the code is not part of
the repo build.

## uppercase

Used by `testProxy`, `testLambdaWithVirtualService`, and the STS credential
variants. No unwrap/transform involved — the raw request body (a JSON
string, e.g. `"solo.io"`) is the `event`. Returns it uppercased.

```python
import json

def lambda_handler(event, context):
    print(event)
    status_code = event.get('statusCode', 200) if isinstance(event, dict) else 200
    return {
        'statusCode': status_code,
        'body': event.upper()
    }
```

## contact-form

Used by `testProxyWithResponseTransform`. Response transformation on the
route unwraps this into an HTML body; test asserts the response contains the
`<meta ... text/html ...>` tag.

```python
import json

def lambda_handler(event, context):
    html_string = """<meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
<html>
<head>
    <title>Test</title>
</head>
<body>
    <h1>Test</h1>
</body>
</html>
    """
    status_code = event.get('statusCode', 200) if isinstance(event, dict) else 200
    return {
        'statusCode': status_code,
        'body': html_string
    }
```

## dumpContext

Used by `testProxyWithRequestTransform` and `testProxyWithRequestAndResponseTransforms`.
Request transformation wraps the incoming HTTP request (headers, body,
method, path, query string, ...) into an API-Gateway-style `event`; this
lambda just dumps that `event` back out as JSON so the test can assert on
its shape.

```python
import json

def lambda_handler(event, context):
    status_code = event.get('statusCode', 200) if isinstance(event, dict) else 200
    return {
        'statusCode': status_code,
        'body': json.dumps(event)
    }
```

## echo

Used by `testProxyWithUnwrapAsApiGateway` and (via `RequestTransformation`)
by `testLambdaTransformations`. With `unwrapAsApiGateway` set, Envoy maps
whatever this lambda returns (`statusCode`/`headers`/`multiValueHeaders`/`body`)
onto the real HTTP response, so it just needs to echo the event back
unchanged — the test controls status/headers/body by putting them directly
in the request payload.

```python
def lambda_handler(event, context):
    # unwrapAsApiGateway pass-through: caller sends an API-Gateway-shaped
    # payload (statusCode/headers/multiValueHeaders/body/...) as the request
    # body; Envoy's unwrap logic maps whatever this lambda returns back onto
    # the real HTTP response, so just echo the event back unchanged.
    return event
```

## non-string-headers-test

Used by `testProxyWithUnwrapAsApiGatewayNonStringHeaderResponse`. Regression
test for a case where Envoy 500'd on non-string values in
`multiValueHeaders`. Returns a `multiValueHeaders.foo` list mixing `None`,
a string, and a number; test expects header `Foo: null,bar,123` and body
`"test body"`.

```python
import json

def lambda_handler(event, context):
    status_code = event.get('statusCode', 200) if isinstance(event, dict) else 200
    return {
        'statusCode': status_code,
        'body': json.dumps('test body'),
        'multiValueHeaders': {
            'foo': [
                None,
                "bar",
                123
            ]
        }
    }
```

## malformed-headers-test

Used by `testProxyWithUnwrapAsApiGatewayMalformedHeaderResponse`. Regression
test for malformed `multiValueHeaders` (an object instead of an array of
values), which previously 500'd. Test expects Envoy to fall back to a 200
with empty body (`{}`) and only default headers (`Date`, `Server`,
`Content-Length`).

```python
import json

def lambda_handler(event, context):
    status_code = event.get('statusCode', 200) if isinstance(event, dict) else 200
    return {
        'statusCode': status_code,
        'body': {},
        # malformed on purpose: multiValueHeaders values must be arrays;
        # 'foo' is an object here, which Envoy's unwrap should reject,
        # falling back to a 200 with no body/headers.
        'multiValueHeaders': {
            'foo': {
                None,
                "bar",
                123
            }
        }
    }
```

## resource-based-cross-account-hello

Used by `testProxyWithCrossAccountLambda`, invoked via a resource-based
policy from a separate AWS account (`AwsAccountId: "986112284769"` on the
upstream). No unwrap involved, so the `statusCode` field is inert — the
test only asserts the response body contains `"Hello from Lambda!"`
(JSON-encoded).

```python
import json

def lambda_handler(event, context):
    status_code = event.get('statusCode', 201) if isinstance(event, dict) else 201
    return {
        'statusCode': status_code,
        'body': json.dumps('Hello from Lambda!')
    }
```
