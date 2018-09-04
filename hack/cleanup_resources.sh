for ns in default gloo-system; do
    for resource in upstream proxy gateway virtualservice; do kubectl delete -n $ns $resource --all; done;
    for secret in  my-precious some-secret ssl-secret; do kubectl delete secret -n $ns $secret; done;
done