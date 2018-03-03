
<h1 align="center">
    <img src="https://i.imgur.com/D1tw77U.png" alt="squash" width="200" height="242">
  <br>
  The Function Gateway
</h1>


<h4 align="center"></h4>
<BR>

Mock API Gateway for building Integrations &amp; Plugins

See the [example module](module/example) for an example of how to write modules.

By default, the example module will be loaded on run, with the [example module config](module/example/example_config.yml) used as the gateway config.

Note: Kubernetes support doesn't exist in this mock. It's designed to run directly on localhost. Therefore `SecretsToWatch` method on your module will not be called, and `secrets` map argument to `Translate` will always be nil (for now).

### Run
```bash
make run
```

### Add a module:

1. call `module.Register(yourModule)` in your module's `init()` function 
2. Add an underscore import to [module/install/modules.go](module/install/modules.go)

### Enable verbose logging
```bash
export DEBUG=1
```
