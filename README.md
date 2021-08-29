# network health admission controller

This admission controller acts as a [MutatingAdmissionWebhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook) and adds [network health sidecars](https://github.com/damoon/network-health-sidecar) to pods.


## Installation

1. install the [dependencies](#Dependencies)
2. download and verify [setup.yaml](setup.yaml)
3. deploy admission controller `kubectl apply -f setup.yaml`

## Dependencies

- [Cert Manager](https://cert-manager.io/docs/installation/helm/#installing-with-helm) is used to [set up certificates](https://cert-manager.io/docs/concepts/ca-injector/) to validate the webhook against the kubernetes control plane.


## Usage

Create a namespace and add the label `network-health-sidecar/enabled: "true"`.

``` yaml
apiVersion: v1
kind: Namespace
metadata:
  name: example
  labels:
    network-health-sidecar/enabled: "true"
```

All pods created in this namespace start with an additional [network health sidecar](https://github.com/damoon/network-health-sidecar) container.


## local development

1. install [tilt](https://docs.tilt.dev/install.html), [helm](https://helm.sh/docs/intro/install/#from-script), [helmfile](https://github.com/roboll/helmfile#installation), [helm diff](https://github.com/databus23/helm-diff#using-helm-plugin-manager--23x), and [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
2. setup [kind with local registry](https://github.com/tilt-dev/kind-local#how-to-try-it)
3. deploy dependencies `helmfile sync`
4. start environment `tilt up`
