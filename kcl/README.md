# function-kcl demo

This demo shows the usage of
[`function-kcl`](https://github.com/crossplane-contrib/function-kcl) to compose
resources within a Crossplane composition and its function pipeline.

The [XRD](definition.yaml) offers a simple `Network` claim abstraction that
allows developers to self-service provision network infrastructure, such as
`VPC` and `InternetGateway`. The following configuration values are exposed on
this claim:

* `count`: The number of network objects to create.
* `includeGateway`: True to create an InternetGateway in addition to the VPC.
* `id`: ID of this Network that will be included in its labels for other objects
  to discover easily.

The [composition](./composition.yaml) uses [KCL](https://kcl-lang.io/) to
dynamically compose resources in accordance with the user provided input from
the XR.

* All input values from the XR are safely extracted, with defaults provided if
  values were not provided on the XR
* A labels dictionary is initialized from the `spec.id` value and is set on
  every created resource
* In a `for` loop, a `count` number of VPC resources are initialized
* If `includeGateway` is true, a `count` number of `InternetGateway` resources
  are initialized in another `for` loop
* The resulting VPC and gateway objects are returned in the `items` collection

## Usage

This demo provides two XRs that can be used to initiate the creation of network
resources.

Create a variable number of network resources that includes both `VPC` and
`InternetGateway` objects:

```console
crossplane beta render xr-with-gateway.yaml composition.yaml functions.yaml -r
```

Create network resources that only include `VPC` objects:

```console
crossplane beta render xr-without-gateway.yaml composition.yaml functions.yaml -r
```

## Example Output

If we create a `XNetwork` object with `count: 2` and `includeGateway: true`, we
should see the following output from `crossplane render`, which shows 2 `VPC`
created and 2 `InternetGateway` created, along with the results providing by
`function-kcl`:

```console
❯ crossplane beta render xr-with-gateway.yaml composition.yaml functions.yaml -r
---
apiVersion: demo-kcl.crossplane.io/v1alpha1
kind: XNetwork
metadata:
  name: xnetwork-kcl
status:
  conditions:
  - lastTransitionTime: "2024-01-01T00:00:00Z"
    message: 'Unready resources: resource-gateway-0, resource-gateway-1, resource-vpc-0,
      and 1 more'
    reason: Creating
    status: "False"
    type: Ready
---
apiVersion: ec2.aws.upbound.io/v1beta1
kind: InternetGateway
metadata:
  annotations:
    crossplane.io/composition-resource-name: resource-gateway-0
  generateName: xnetwork-kcl-
  labels:
    crossplane.io/composite: xnetwork-kcl
    networks.meta.fn.crossplane.io/network-id: xnetwork-kcl
  name: gateway-0
  ownerReferences:
  - apiVersion: demo-kcl.crossplane.io/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: XNetwork
    name: xnetwork-kcl
    uid: ""
spec:
  forProvider:
    region: eu-west-1
    vpcIdSelector:
      matchControllerRef: true
---
apiVersion: ec2.aws.upbound.io/v1beta1
kind: InternetGateway
metadata:
  annotations:
    crossplane.io/composition-resource-name: resource-gateway-1
  generateName: xnetwork-kcl-
  labels:
    crossplane.io/composite: xnetwork-kcl
    networks.meta.fn.crossplane.io/network-id: xnetwork-kcl
  name: gateway-1
  ownerReferences:
  - apiVersion: demo-kcl.crossplane.io/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: XNetwork
    name: xnetwork-kcl
    uid: ""
spec:
  forProvider:
    region: eu-west-1
    vpcIdSelector:
      matchControllerRef: true
---
apiVersion: ec2.aws.upbound.io/v1beta1
kind: VPC
metadata:
  annotations:
    crossplane.io/composition-resource-name: resource-vpc-0
  generateName: xnetwork-kcl-
  labels:
    crossplane.io/composite: xnetwork-kcl
    networks.meta.fn.crossplane.io/network-id: xnetwork-kcl
  name: vpc-0
  ownerReferences:
  - apiVersion: demo-kcl.crossplane.io/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: XNetwork
    name: xnetwork-kcl
    uid: ""
spec:
  forProvider:
    cidrBlock: 192.168.0.0/16
    enableDnsHostnames: true
    enableDnsSupport: true
    region: eu-west-1
---
apiVersion: ec2.aws.upbound.io/v1beta1
kind: VPC
metadata:
  annotations:
    crossplane.io/composition-resource-name: resource-vpc-1
  generateName: xnetwork-kcl-
  labels:
    crossplane.io/composite: xnetwork-kcl
    networks.meta.fn.crossplane.io/network-id: xnetwork-kcl
  name: vpc-1
  ownerReferences:
  - apiVersion: demo-kcl.crossplane.io/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: XNetwork
    name: xnetwork-kcl
    uid: ""
spec:
  forProvider:
    cidrBlock: 192.168.0.0/16
    enableDnsHostnames: true
    enableDnsSupport: true
    region: eu-west-1
---
apiVersion: render.crossplane.io/v1beta1
kind: Result
message: created resource "gateway-0:InternetGateway"
severity: SEVERITY_NORMAL
step: render-with-kcl
---
apiVersion: render.crossplane.io/v1beta1
kind: Result
message: created resource "gateway-1:InternetGateway"
severity: SEVERITY_NORMAL
step: render-with-kcl
---
apiVersion: render.crossplane.io/v1beta1
kind: Result
message: created resource "vpc-0:VPC"
severity: SEVERITY_NORMAL
step: render-with-kcl
---
apiVersion: render.crossplane.io/v1beta1
kind: Result
message: created resource "vpc-1:VPC"
severity: SEVERITY_NORMAL
step: render-with-kcl
```