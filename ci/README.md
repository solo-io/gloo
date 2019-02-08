# Commands use to setup build environment

## One time init:

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

## Github key

ssh-keyscan -t rsa github.com > ./ci/github_known_hosts

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

# More info:
https://cloud.google.com/cloud-build/docs/securing-builds/use-encrypted-secrets-credentials
https://cloud.google.com/cloud-build/docs/access-private-github-repos