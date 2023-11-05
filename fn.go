package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/upbound/provider-aws/apis/s3/v1beta1"

	"github.com/crossplane/function-sdk-go/errors"
	"github.com/crossplane/function-sdk-go/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/resource/composed"
	"github.com/crossplane/function-sdk-go/response"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction observes an example composite resource (XR). It simple adds one
// S3 bucket to the desired state.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.log.Info("Running Function", "tag", req.GetMeta().GetTag())
	rsp := response.To(req, response.DefaultTTL)

	// create a single test S3 bucket
	_ = v1beta1.AddToScheme(composed.Scheme)
	name := "test-bucket"
	b := &v1beta1.Bucket{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"crossplane.io/external-name": name,
			},
		},
		Spec: v1beta1.BucketSpec{
			ForProvider: v1beta1.BucketParameters{
				Region: ptr.To[string]("us-east-2"),
			},
		},
	}

	// read the desired composed resources so we can update them with our bucket
	desired, err := request.GetDesiredComposedResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get desired resources from %T", req))
		return rsp, nil
	}

	// add our bucket to the desired composed resources
	cd, _ := composed.From(b)
	desired[resource.Name("xbuckets-"+name)] = &resource.DesiredComposed{Resource: cd}
	if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composed resources in %T", rsp))
		return rsp, nil
	}

	return rsp, nil
}
