# integreatly-operator-cleanup-harness

This harness sets up the [cluster-service](https://github.com/integr8ly/cluster-service) image by with the required aws credentials and cluster id.
When running this image there are two args that  can be passed in.

- `cleanup` -> required to run being the cluster-service. Without this the image will exit.
- `dry-run` -> set the cluster-service into dry run mode 

## Building the image
To run the image on a cluster you must build and push the image to an image repository.
```
podman build -t <path/to/repository>/integreatly-operator-cleanup-harness .
```
```
podman push <path/to/repository>/integreatly-operator-cleanup-harness
```

## Running the image
Run the image with `cluster-admin` permissions.
Set the cluster role bindings.
```
kubectl create clusterrolebinding --user system:serviceaccount:kube-system:default namespace-cluster-admin --clusterrole cluster-admin
```
Deploy the image into the `kube-system` namespace with the `cleanup` arg.
`dry-run` is optional and will run teh cluster-service in a dry run mode.
```
oc run -n kube-system --restart=Never --image <path/to/repository>/integreatly-operator-cleanup-harness -- integreatly-operator-cleanup-harness cleanup
```

## On Cluster
On the cluster when the pod have been deployed in the `kube-system` name space.
Two pods will be created and should run to completion without restarting.
These pods are:

- integreatly-operator-cleanup-harness
- integreatly-operator-cluster-service