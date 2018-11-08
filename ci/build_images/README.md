# Commands use to setup build environment

# build container:

```
cd dep && docker build -t soloio/dep . && docker push soloio/dep
```

# One time init:

```
KEYRING=build
KEYNAME=build-key
SERVICE_ACCOUNT=825641009090

gcloud kms keyrings create ${KEYRING} \
  --location=global

gcloud kms keys create ${KEYNAME} \
  --location=global \
  --keyring=${KEYRING} \
  --purpose=encryption

gcloud kms keys add-iam-policy-binding \
    ${KEYNAME} --location=global --keyring=${KEYRING} \
    --member=serviceAccount:${SERVICE_ACCOUNT}@cloudbuild.gserviceaccount.com \
    --role=roles/cloudkms.cryptoKeyEncrypterDecrypter
```

## Encrypt secrets:
Get the solobot private key and use this to encrypt:

```
gcloud kms encrypt \
  --plaintext-file=${HOME}/Documents/solo/bot/id_rsa \
  --ciphertext-file=./ci/id_rsa.enc \
  --location=global \
  --keyring=${KEYRING} \
  --key=${KEYNAME}
```

# Test dep container locally

copy the unencrypted id_rsa to the ci folder. and run:

```
cp ${HOME}/Documents/solo/bot/id_rsa ./ci/id_rsa
docker run --rm -t -i -v $PWD:/workspace -w /workspace -e PROJECT_ROOT=github.com/solo-io/solo-kit \
    --entrypoint=/bin/bash soloio/dep
rm ./ci/id_rsa
```

Run a build with:

```
gcloud builds submit --config ./cloudbuild.yaml .
```