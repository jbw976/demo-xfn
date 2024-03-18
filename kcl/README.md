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
* `region`: The AWS region to create the resources in

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

## Deploying this Composition Resources to AWS

Now that we have tested locally, we can create these resources in AWS.

In the next few steps, we'll:

* Install Crossplane
* Install the AWS Provider and set up Authentication
* Install the KCL Function
* Apply the CompositeResourceDefinition (XRD) and Composition
* Create your Claim which will trigger the Composition to create the Resources

### Installing Crossplane

Install Crossplane via the [Getting Started](https://docs.crossplane.io/v1.15/getting-started/provider-aws/) guide. For this example a minimum Crossplane version of 1.14.x is required.

### Installing the EC2 Provider

Once Crossplane is installed, install the AWS Provider from the manifest at <https://marketplace.upbound.io/providers/upbound/provider-aws-ec2/v1.2.0>. A copy of the manifest is in this directory.

The major Cloud providers (AWS, Azure, GCP) are broken up in to "family" providers in order to reduce CRDs installed on the cluster, which is why there is a separate provider for the EC2 related resources.

```shell
kubectl apply -f provider.yaml
```

Verify that your Provider was installed correctly. You should see that the `upbound-provider-family-aws` provider was automatically installed. This provider supplies the common `ProviderConfig` CRD.

```shell
❯ kubectl get provider.pkg
NAME                      INSTALLED   HEALTHY   PACKAGE                                          AGE
upbound-provider-aws-ec2      True        True      xpkg.upbound.io/upbound/provider-aws-ec2:v1.2.0      5h17m
upbound-provider-family-aws   True        True      xpkg.upbound.io/upbound/provider-family-aws:v1.2.0   5h16m
```

## Creating a Secret to Authenticate to AWS

In this example, we are going to create a Kubernetes secret that contains AWS credentials to allow Crossplane to provision AWS Resources.

The provider will look for `[default]` credentials in the secret with the following format. Save your credentials into the file `aws-credentials.txt`.

```ini
[default]
aws_access_key_id = ...
aws_secret_access_key = ...
```

Next create the secret. Note the secret `name`, `namespace`. In this example, the secret key will be `creds`.

```shell
kubectl create secret \
generic aws-secret \
-n crossplane-system \
--from-file=creds=./aws-credentials.txt

```

### Create the ProviderConfig

`ProviderConfigs` tell Crossplane Providers how to authenticate to remote APIs. There can be multiple ProviderConfigs on a Crossplane Cluster. The `default` ProviderConfig is used by default.

Your ProviderConfig should match the `name`, `namespace`, and `key` of the secret you created:

```yaml
apiVersion: aws.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: aws-secret
      key: creds
```

### Installing the Functions

Functions are also installed via a Crossplane package. The KCL function can be found at <https://marketplace.upbound.io/functions/crossplane-contrib/function-kcl/v0.2.0>. We will also be using the [auto-ready](https://marketplace.upbound.io/functions/crossplane-contrib/function-auto-ready/v0.2.1) function to ensure resource ready status is propagated to the Composite Resource.

```yaml
apiVersion: pkg.crossplane.io/v1beta1
kind: Function
metadata:
  name: function-kcl
spec:
  package: xpkg.upbound.io/crossplane-contrib/function-kcl:v0.2.0
```

```console
❯ kubectl apply -f functions.yaml
function.pkg.crossplane.io/function-kcl created
function.pkg.crossplane.io/function-auto-ready created
```

Ensure the KCL functions are healthy:

```console
❯ kubectl get -f functions.yaml 
NAME                  INSTALLED   HEALTHY   PACKAGE                                                         AGE
function-kcl          True        True      xpkg.upbound.io/crossplane-contrib/function-kcl:v0.2.0          4h47m
function-auto-ready   True        True      xpkg.upbound.io/crossplane-contrib/function-auto-ready:v0.2.1   3h56m                                                 AGE

```

### Creating Resources

Now we can create resources in AWS. There are two example manifests

```console
❯ kubectl apply -f xr-with-gateway.yaml 
xnetwork.demo-kcl.crossplane.io/xnetwork-kcl created
```

We can use the `crossplane beta trace` command to get the status of all the resources in the composition. It may take a few minutes for all the resources to become available.

```console
❯ crossplane beta trace xnetwork.demo-kcl.crossplane.io/xnetwork-kcl
NAME                           SYNCED   READY   STATUS
XNetwork/xnetwork-kcl          True     True    Available
├─ InternetGateway/gateway-0   True     True    Available
├─ InternetGateway/gateway-1   True     True    Available
├─ VPC/vpc-0                   True     True    Available
└─ VPC/vpc-1                   True     True    Available
```

### Cleanup

To remove the resources from AWS, delete any `XNetworks` that you created. For example:

```console
❯ kubectl delete -f xr-with-gateway.yaml 
xnetwork.demo-kcl.crossplane.io "xnetwork-kcl" deleted
```

