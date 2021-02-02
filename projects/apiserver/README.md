# gloo-fed-apiserver
Gloo Federation Apiserver

The apiserver is used to power both the gloo-fed UI
and the glooctl fed extension's 'get' commands.

There are four 'packages' that encompass all Gloo Edge resources:
1. gloo.solo.io
2. gateway.solo.io
3. enterprise.gloo.solo.io
4. ratelimit.api.solo.io

The definitions for these resources must be skv2-compatible, or
they must have a Spec and Status in the form of:
```
message UpstreamSpec {
    ...
}

message UpstreamStatus {
    ...
}
```

These are defined in the solo-apis repository.

Gloo Federation then defines wrappers around these resources in the
following packages:
1. fed.gloo.solo.io
2. fed.gateway.solo.io
3. fed.enterprise.gloo.solo.io
4. fed.ratelimit.solo.io

Note: ratelimit.api.solo.io and fed.ratelimit.solo.io use v1alpha1 instead of v1.

Finally, there are a few custom types defined in Gloo Federation.

1. fed.solo.io
    - FailoverScheme
    - GlooInstance

The apiserver is responsible for making all of these resources
available to the UI.