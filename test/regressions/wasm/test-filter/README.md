This compiles an example filter for envoy WASM.

# build filter
build with
```
bazel build :filter.wasm
```

Filter will be in:
```
./bazel-bin/filter.wasm
```

build and push  wasme image:
```
wasme build precompiled --tag webassemblyhub.io/yuval/header-test:v0.3 --config runtime-config.json ./bazel-bin/filter.wasm
wasme push webassemblyhub.io/yuval/header-test:v0.3
```
