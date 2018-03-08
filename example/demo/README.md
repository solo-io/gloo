Self Demo
==========
This document outlines a demo of gloo composing ('glooing') an application from a monolyth, a microservice and AWS Lambda.

# Setup the environment

## Install Kubernetes

```shell
minikube start --extra-config=apiserver.Authorization.Mode=RBAC --cpus 4 --memory 4096
kubectl create clusterrolebinding permissive-binding \
         --clusterrole=cluster-admin \
         --user=admin \
         --user=kubelet \
         --group=system:serviceaccounts
```

## Install Gloo
```shell
kubectl apply \
          -f https://raw.githubusercontent.com/solo-io/gloo-install/master/kube/install.yaml
```

## Get the url of the ingress
If you installed kubernetes using minikube as mentioned above, you can use this command:
```shell
GATEWAY_ADDR=$(kubectl get po -l gloo=ingress -n gloo-system -o 'jsonpath={.items[0].status.hostIP}'):$(kubectl get svc ingress -n gloo-system -o 'jsonpath={.spec.ports[?(@.name=="http")].nodePort}')

export GATEWAY_URL=http://$GATEWAY_ADDR
```

# Deploy
## Deploy the pet store monolith

Note - source code for this demo is here: https://github.com/solo-io/spring-petclinic, https://github.com/solo-io/petclinic-vet.

```shell
kubectl apply -f yamls/pet-clinic.yaml
```

Verify it is available as an upstream:

```shell
kubectl get upstreams -n gloo-system
```

you should see `default-petclinic-80` in the list.

## Add a route to it!

```shell
glooctl route create --path-prefix / --upstream default-petclinic-80
```

Now you should see pet clinict in your browser if you open the $GATEWAY_URL.

Notice that the vets page is missing the city column. To fix the bug, we will deploy a microservice that replaces the vets page in the monolith.

Deploy the microservice:

```shell
kubectl apply -f yamls/vets.yaml
```

Add a route to the new microservice from gloo:

```shell
glooctl route create --sort --path-exact /vets.html --upstream default-petclinic-vets-80
```

Now the vets page will contain a city column!

# Add some Cloud™

Let's expand the app functionality by displaying a contact form and saving the contact response to an AWS S3 bucket.

Note: the current Gloo-AWS integration sends the aws keys unencrypted from gloo to envoy. This means
that they may be sent in the clear over the local network, and appear in the envoy debug logs (which are off by default). we plan to fix that very soon.

In this section, we will:
1. Create an S3 bucket to receive form responses
2. Create a policy and a role to allow a lambda function to put objects in the bucket.
3. Create a lambda functions to display and process the form.
4. Configure gloo to route requests to the lambda function.

## Configure AWS

Create an S3 bucket, set the BUCKET variable to a new name.
```shell
BUCKET=name-of-bucket-to-add
aws s3api  create-bucket  --acl private --bucket $BUCKET
```

change the policy to match your bucket:

```shell
sed s/io.solo.petclinic/$BUCKET/ aws/policy-document_template.json > aws/policy-document.json

```

Create the needed policy and role that allows the lambda function to access the bucket:

```shell
POLICY_ARN=$(aws iam list-policies | jq -r '.Policies[] | select(.PolicyName == "gloo-contact-lambda-policy") | .Arn')
if [ -z "$POLICY_ARN" ]; then
    aws iam create-policy --policy-name gloo-contact-lambda-policy --policy-document file://aws/policy-document.json
    POLICY_ARN=$(aws iam list-policies | jq -r '.Policies[] | select(.PolicyName == "gloo-contact-lambda-policy") | .Arn')
fi

ROLE_ARN=$(aws iam list-roles | jq -r '.Roles[] | select(.RoleName == "gloo-contact-lambda-role") | .Arn')
if [ -z "$ROLE_ARN" ]; then
    aws iam create-role --role-name gloo-contact-lambda-role --assume-role-policy-document file://aws/assume-role-policy-document.json
    aws iam attach-role-policy --role-name gloo-contact-lambda-role --policy-arn $POLICY_ARN
    ROLE_ARN=$(aws iam list-roles | jq -r '.Roles[] | select(.RoleName == "gloo-contact-lambda-role") | .Arn')
fi

```

Create the lambda function:

```shell
aws lambda create-function \
--region us-east-1 \
--function-name processContact \
--zip-file fileb://aws/contact.zip \
--handler index.handler \
--runtime nodejs6.10 \
--role $ROLE_ARN \
--environment "Variables={BUCKET=$BUCKET}"
```

## Route to AWS Lambda from gloo

Upload the aws secret to gloo, so that gloo can call your function:

```shell
glooctl secret create aws --name aws-lambda-us-east-1
```

Create the aws upstream in gloo:

```shell
glooctl upstream create -f yamls/aws-upstream.yaml
```

Route the relevant paths to the function:

```shell
glooctl route create --sort --path-exact /contact --upstream aws-lambda-us-east-1 --function 'processContact:$LATEST'
glooctl route create --sort --path-exact /contact.html --upstream aws-lambda-us-east-1 --function 'processContact:$LATEST'
```

Go to the contact page. Notice that you see json in stead of HTML. To fix that, we will attach a 
transformation to the route:

```shell
glooctl route update --path-exact /contact.html --upstream aws-lambda-us-east-1 --function 'processContact:$LATEST' --extensions ./yamls/transformation.yaml
```

Now the contact form will be presented! Post some messages through the contact form.

To view the sent messages, have a look at the S3 bucket:

```shell
aws s3 ls $BUCKET
```

# Cleanup
## Cloud reosurces

```shell
aws lambda delete-function --function-name processContact
aws iam detach-role-policy --role-name gloo-contact-lambda-role --policy-arn $POLICY_ARN
aws iam    delete-role --role-name gloo-contact-lambda-role
aws iam    delete-policy --policy-arn "$POLICY_ARN"
```
Delete all the objects in the bucket. and then:
```
aws s3api  delete-bucket --bucket $BUCKET
```

## Minikube
```
minikube delete
```