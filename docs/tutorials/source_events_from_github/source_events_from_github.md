# Deploy Gloo
k apply -f https://raw.githubusercontent.com/solo-io/gloo-install/master/kube/install.yaml

# Deploy NATS
k apply -f deploy-nats.yaml

# Create a route for nats
glooctl route create --sort \
    --path-exact /github-webhooks \
    --upstream default-nats-streaming-4222 \
    --function github-webhooks

# deploy the image-pusher service
k apply -f  deploy-image-pusher.yaml

# deploy the 
k apply -f deploy-
