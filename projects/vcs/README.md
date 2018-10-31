# VCS
Version control for Gloo resources

## WIP
- We are still working out the spec and implementation details for this feature
- The code in this directory is for research purposes at the moment
- Version control may become a feature of all Solo.io products so we aim for reuseability

## Notes on current code
- features two clients (kubernetes and filesystem) pointing at the same resources (virtualServices)


# What to persist?
## Solo-io CRDs
We may not want to put all of our CRDs in git
- gateways.gateway.solo.io                                      
  - yes?
- proxies.gloo.solo.io                                          
  - yes?
- resolvermaps.sqoop.solo.io                                    
  - yes
- schemas.sqoop.solo.io
  - no - this will be deprecated, merged with resolver maps
- settings.gloo.solo.io
  - yes
- upstreams.gloo.solo.io
  - no - these are discovered
  - what to do about manually entered upstreams?
    - create a new crd called something like ManualUpstream?
- virtualservices.gateway.solo.io
  - yes

## Other CRDs
- there are many crds (see `crd_list.txt` for some examples)
- do we care to persist any of them?
  - I (mitch) think we should only persist the crds that we own. We can make new solo.io crds to hold any built-in crd information that we want to manage. This will give us greater control and storage-agnostic

# Usage
- The current implementation is for exploration only. To run it:
```
go run cmd/main.go
```
- It will create a directory called `gloo/gloo-system/virtualService` and it will contain the first VS that exists in your deployment under the gloo-system namespace
