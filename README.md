# xfn-demo

## pre-reqs

* install CLI with https://docs.crossplane.io/latest/cli/#installing-the-cli
* golang 1.21+

## Initialize a new Function Project

Start by initializing a new Function project into a new directory called
`xfn-demo`:
```
crossplane beta xpkg init xfn-demo function-template-go -d xfn-demo
```

The `init` command will show us some notes about next steps for our functions
project and also offer to run a helpful initialization script to customize the
project further for us.

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
new message, something like "Hello Rejekts Paris":
```yaml
    input:
      apiVersion: template.fn.crossplane.io/v1beta1
      kind: Input
      example: "Hello Rejekts Paris"
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
message: I was run with input "Hello Rejekts Paris"!
severity: SEVERITY_NORMAL
step: run-the-template
```

## Workflow: Updating Function, Testing, and Validation

We've been editing and testing just the composition input, let's take this a
step further by editing the actual code of the function itself.

Edit the `fn.go` function code now to create a simple S3 bucket.

The full/final `fn.go` file is available in this repo for you to copy:
[`fn.go`](./fn.go). The file can be directly copied into your repo with:
```console
 curl -s https://raw.githubusercontent.com/jbw976/demo-xfn/main/fn.go > fn.go
 ```

After your `fn.go` is updated, make sure all go modules are available: 
```
go mod tidy
```

Run the function locally:
```
go run . --insecure --debug
```

### Test with `render`

`render` the new function code we wrote to test our changes.  We should see a S3
bucket in the output of desired resources:
```yaml
❯ crossplane beta render xr.yaml composition.yaml functions.yaml -x
---
apiVersion: example.crossplane.io/v1
kind: XR
metadata:
  name: example-xr
spec: {}
status:
  conditions:
  - lastTransitionTime: "2024-01-01T00:00:00Z"
    message: 'Unready resources: xbuckets-test-bucket'
    reason: Creating
    status: "False"
    type: Ready
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

### Validate Output

In addition to testing our function and composition locally with `render`, we
can also verify the resources it generated are compliant with their defined
schemas using `crossplane validate`. We tell `validate` the providers and
resources we are using, so it can automatically download their schemas.

Which providers we are using is captured in an `extensions.yaml` file and is
available in this repo for you to copy: [`extensions.yaml`](./extensions.yaml).
The file can be directly copied into your repo with:
```console
 curl -s https://raw.githubusercontent.com/jbw976/demo-xfn/main/extensions.yaml > extensions.yaml
 ```

Then we can run `validate` on the piped output of `render` to verify our
generated resources are well formed:
```console
crossplane beta render xr.yaml composition.yaml functions.yaml -x | crossplane beta validate extensions.yaml -
```

The full output of this will look like this:
```console
❯ crossplane beta render xr.yaml composition.yaml functions.yaml -x | crossplane beta validate extensions.yaml -
[✓] example.crossplane.io/v1, Kind=XR, example-xr validated successfully
[✓] s3.aws.upbound.io/v1beta1, Kind=Bucket, xbuckets-test-bucket validated successfully
Total 2 resources: 0 missing schemas, 2 success cases, 0 failure cases
```