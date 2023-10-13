# FIPS Support

## What is FIPS?
FIPS is a set of rules that outline the basic security needs of cryptographic modules used in computer and telecommunication systems. Compliance with these rules is mandatory for certain industries (ie. healthcare, finance) that utilize cryptographic modules to protect sensitive data. The publications and documents associated with FIPS are issued by the [National Institute of Standards and Technology (NIST)](https://www.nist.gov/).

The most recent publication of FIPS, [FIPS 140-2](https://en.wikipedia.org/wiki/FIPS_140-2), is a current set of cryptographic standards that applications may need to adhere to.

## Why do customers use it?
As explained above, FIPS standards are mandatory for certain industries. In order to be compliant, applications must use FIPS-compliant cryptographic modules. 

One example is PCI compliance, which is a set of policies and procedures developed to protect credit, debit and cash card transactions and prevent the misuse of cardholders' personal information.

## How is it supported in Go?
In Golang, native cryptography is not FIPS friendly. Per a [golang issue for fips support](https://github.com/golang/go/issues/11658#issuecomment-120448974), "Go's crypto is not FIPS 140 validated and I'm afraid that there is no possibility of that happening in the future either."

Since [go 1.19](https://github.com/golang/go/issues/51940) BoringSSL based crypto is part of the main branch. This means that in go 1.19 and up, pass `GOEXPERIMENT=boringcrypto` to the go tool during build time.

## How is it configured in Gloo Edge?
Since Gloo Edge is built in Go, we can use the `GOEXPERIMENT=boringcrypto` flag to build a FIPS compliant version of modules.

FIPS support was introduced in the following work:
- [Initial support of FIPS](https://github.com/solo-io/solo-projects/issues/2420)
- [Update runtime to use go 1.20](https://github.com/solo-io/solo-projects/pull/4586)
- [Fix FIPS and support discovery](https://github.com/solo-io/solo-projects/pull/5368)

## How is it validated in Gloo Edge?
We validate FIPS compliance by running the following command on our images during the build pipeline:
```
goversion -crypto [BINARY]] | grep "(boring crypto)" > /dev/null
```

This command checks if the binary was built with the `GOEXPERIMENT=boringcrypto` flag. If the binary was built with the flag, the command will return a 0 exit code. If the binary was not built with the flag, the command will return a non-zero exit code. For ease of use, we expose this as a target in our Makefile:
```
validate-boring-crypto
```

## Which components have FIPS variants built?
Gloo Edge Enterprise supports a FIPS variant of the data-plane, since those components are responsible for handling sensitive data.

We build a FIPS variant of the following images:
- gateway-proxy
- gloo
- discovery
- ext-auth
- ext-auth-plugins
- rate-limit

## Which components DO NOT have FIPS variants built?
- observability (Enterprise)
- caching (Enterprise)
- all Gloo Federation components (Enterprise)
- certgen (OSS)
- sds (OSS)
- ingress (OSS)
- access-logger (OSS)
- kubectl (OSS)
