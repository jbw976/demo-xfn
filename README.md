# xfn-demo

## pre-reqs

* install CLI with https://docs.crossplane.io/v1.14/cli/#installing-the-cli
* golang 1.21+

## Initialize a new Function Project

Start by initializing a new Function project into a new directory called
`xfn-demo`:
```
crossplane beta xpkg init xfn-demo function-template-go -d xfn-demo
```

Run the function locally so it is ready to process input and serve responses:
```
cd xfn-demo
go run . --insecure --debug
```

## Workflow: Updating Composition and Testing

Use the `render` command to test the unedited Function and composition input:
```
cd example
crossplane beta render xr.yaml composition.yaml functions.yaml -r
```

We will see the function output the XR and a `Result` that includes the default
`Input`'s "Hello world" message:
```yaml
❯ crossplane beta render xr.yaml composition.yaml functions.yaml -r
---
apiVersion: example.crossplane.io/v1
kind: XR
metadata:
  name: example-xr
---
apiVersion: render.crossplane.io/v1beta1
kind: Result
message: I was run with input "Hello world"!
severity: SEVERITY_NORMAL
step: run-the-template
```

### Update the Composition and Test Again

Now we'll update just the Composition input in `composition.yaml` to specify a
new message, something like "Hello Kubecon Chicago":
```yaml
    input:
      apiVersion: template.fn.crossplane.io/v1beta1
      kind: Input
      example: "Hello Kubecon Chicago"
```

`render` the function again to see that the new `Input` from the Composition
affects the output `Result`:
```yaml
❯ crossplane beta render xr.yaml composition.yaml functions.yaml -r
---
apiVersion: example.crossplane.io/v1
kind: XR
metadata:
  name: example-xr
---
apiVersion: render.crossplane.io/v1beta1
kind: Result
message: I was run with input "Hello Kubecon Chicago"!
severity: SEVERITY_NORMAL
step: run-the-template
```

## Workflow: Updating Function and Testing

We've been editing and testing just the composition input, let's take this a
step further by editing the actual code of the function itself.

Edit the `fn.go` function code now to create a simple S3 bucket.

The full/final `fn.go` file is available in this repo for you to copy:
[`fn.go`](./fn.go)

After your `fn.go` is updated, make sure all go modules are available: 
```
go mod tidy
```

Run the function locally:
```
go run . --insecure --debug
```

`render` the new function code we wrote to test our changes.  We should see a S3
bucket in the output of desired resources:
```yaml
❯ crossplane beta render xr.yaml composition.yaml functions.yaml -r
---
apiVersion: example.crossplane.io/v1
kind: XR
metadata:
  name: example-xr
---
apiVersion: s3.aws.upbound.io/v1beta1
kind: Bucket
metadata:
  annotations:
    crossplane.io/composition-resource-name: xbuckets-test-bucket
    crossplane.io/external-name: test-bucket
  generateName: example-xr-
  labels:
    crossplane.io/composite: example-xr
  ownerReferences:
  - apiVersion: example.crossplane.io/v1
    blockOwnerDeletion: true
    controller: true
    kind: XR
    name: example-xr
    uid: ""
spec:
  forProvider:
    region: us-east-2
```
