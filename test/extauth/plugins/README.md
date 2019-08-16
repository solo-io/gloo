# Sample ExtAuth plugins

These plugins are used to test the loader in the `projects/extauth/pkg/plugins` package. The tests should succeed
when running with `ginkgo -r`, but if you are trying to debug with `dlv` (this includes Goland IDE debug
configurations, which are backed by `dlv`), the `plugin.Open` function will likely fail.

In order to be able to debug the loader tests locally, you need to disable compiler optimizations and inlining
when building the plugins via the additional `gcflags` option:

```bash
go build -buildmode=plugin -gcflags="all=-N -l" -o TestPlugin.so interface.go
```

This is because `dlv` builds everything with `-gcflags="all=-N -l"` and the plugins must be compiled in the
same way. See [this comment](https://github.com/go-delve/delve/issues/865#issuecomment-480766102) for reference.

Conversely, building the plugin this way will cause `ginkgo -r` to fail on the aforementioned `plugin.Open`
function. If for some reason you want/need to run the loader test for the unoptimized plugin build with the
ginkgo CLI, you have to pass the same `gcflags` to ginkgo:

```bash
ginkgo -gcflags="all=-N -l"
```
