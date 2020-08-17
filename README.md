# integreatly-operator-cleanup-harness

This harness sets up the `cluster-service` image by with the required aws credentials and cluster id.
When running this image there are two args that  can be passed in.

- `cleanup` -> required to run being the cluster-service. Without this the image will exit.
- `dry-run` -> set the cluster-service into dry run mode 
