# K8s proxy image swapper

This project redirects all public images to an (internal) proxy by dynamically
renaming those images. The proxy is configured through `PROXY_URL`.

The goal is to prevent errors from dockerhub usage rate limit by providing a
cache as well as doing some security analysis on the proxy side.
(Note that this project only patches images name, to redirect them on a proxy 
registry; it does not deploy the proxy by itself.
the official docker-registry can be configured as a proxy registry).

Patching is done automatically, there is no need to change anything on
your deployments.

Images patched by k8s-proxy-image-swapper (k8s-pisw) also have a label
"k8s-proxy-image-swapper: patched-image", to recognize them.

# Setup

Certificates are managed by Cert-manager. Be sure to install it on your cluster 
before using this software.

Afterwards you can use the helm chart in ./chart to deploy this software.
You need to have the image of the patcher built and pushed on a custom registry,
because no exception is currently made on k8s-proxy-image-swapper (self).
You need to provide the base64 of the CA for the cluster. Further details
are available in ./chart/values.yaml.

# Troubleshouting

In case of issue with image patching, to remove image-swapper, use the following command:

`kubectl delete mutatingwebhookconfiguration k8s-proxy-image-swapper-webhook`

This will delete the webhook and unblock your cluster. Upgrade the chart to reinstate
the webhook.

# Inner working

This software uses MutatingWebHook (from dynamic admission control in k8s)
to patch the `image` field in a pod (`containers` and `initContainers`) to
use a proxy registry (docker registry for instance).

See the unit tests in `mutate/mutate_test.go` for patching examples.
Note that this image must be stored in a registry different than the
Docker Hub. Otherwise you may have a chicken and egg problem.

The simple solution to unblock yourself when the proxy doesn't work for instance
is to simply delete the mutating webhook :
`kubectl delete MutatingWebHookConfiguration -n kube-system k8s-proxy-image-swapper-webhook`

# Contributing

Contributions implies licensing those contributions under the terms of LICENSE.

The GitHub repository https://github.com/Polyconseil/k8s-proxy-image-swapper
is the official repository.

To contribute you need a GitHub account.

## Opening issues

Please make sure there is no open issue on the topic.

## Submitting changes

Please format the commit messages according to semantic-release.

A good commit message includes relevant information about *why* a change
has been made (this might also be a good idea to put this kind of information
in comments), so that other developers can later understand why a change was made.

The project follows a semver versionning scheme.

# Architecture

- ./main.go contains the setup code and configuration code.
- ./mutate contains the code that patches the images. (if you ever modify this
code, please run the tests and add or modify the tests accordingly).
- ./chart contains the chart to deploy the software.
- ./tests contains some manifests to help test the software.


## Building
The Docker image can be build with the Dockerfile with :

```
docker build .
```

Or with Nix :

```
nix build .#packages.x86_64-linux.oci-k8s-proxy-image-swapper
```
