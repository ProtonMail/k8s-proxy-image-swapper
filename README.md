# K8s proxy image swapper

This software uses MutatingWebHook (from dynamic admission control in k8s)
to patch the `image` field in a pod (`containers` and `initContainers`) to
use a proxy registry (docker registry for instance).

See the unit tests in `mutate/mutate_test.go` for patching examples.
Note that the image must be stored in a registry different than the
Docker Hub. Otherwise you may have a chicken and egg problem.

The simple solution to unblock yourself when the proxy doesn't work for instance
is to simply delete the mutating webhook :
`kubectl delete MutatingWebHookConfiguration -n kube-system k8s-proxy-image-swapper-webhook`
