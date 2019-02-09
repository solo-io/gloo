

# Usage
## Install a new glooe distribution from scratch on Kubernetes
```bash
# Setup your repo
make init
make pin-repos
dep ensure -v # you may need to repeat this and make pin-repos once or twice
make allprojects

# at this point you should have gloo built to you ./_output/ directory
# make the manifest
VERSION="1.10.0" make manifest # note that there is no "v" in the version

# install
# prep: create a secret with you docker credentials
./_output/glooctl install kube -f ./install/manifest/glooe-distribution.yaml
# NOTE: glooe-distribution.yaml is the same as glooe-release.yaml except that "distribution" uses an IfNotPresent pull policy
```

# Additional Notes
- Shared projects across Solo.io.
- This repo contains the git history for Gloo and Solo-Kit. 
