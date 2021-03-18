/*
Copyright 2019 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"

	clientset "github.com/crossplane/provider-aws/pkg/clients/acm"
)

// this ensures that the mock implements the client interface
var _ clientset.Client = (*MockCertificateClient)(nil)

// MockCertificateClient is a type that implements all the methods for Certificate Client interface
type MockCertificateClient interface {
	// GetCertificateRequest(*acm.GetCertificateInput) acm.GetCertificateRequest
	DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(options *acm.Options)) (*acm.DescribeCertificateOutput, error)
	RequestCertificate(ctx context.Context, params *acm.RequestCertificateInput, optFns ...func(options *acm.Options)) (*acm.RequestCertificateOutput, error)
	DeleteCertificate(ctx context.Context, params *acm.DeleteCertificateInput, optFns ...func(options *acm.Options)) (*acm.DeleteCertificateOutput, error)
	UpdateCertificateOptions(ctx context.Context, params *acm.UpdateCertificateOptionsInput, optFns ...func(options *acm.Options)) (*acm.UpdateCertificateOptionsOutput, error)
	ListTagsForCertificate(ctx context.Context, params *acm.ListTagsForCertificateInput, optFns ...func(options *acm.Options)) (*acm.ListTagsForCertificateOutput, error)
	AddTagsToCertificate(ctx context.Context, params *acm.AddTagsToCertificateInput, optFns ...func(options *acm.Options)) (*acm.AddTagsToCertificateOutput, error)
	RenewCertificate(ctx context.Context, params *acm.RenewCertificateInput, optFns ...func(options *acm.Options)) (*acm.RenewCertificateOutput, error)
	RemoveTagsFromCertificate(ctx context.Context, params *acm.RemoveTagsFromCertificateInput, optFns ...func(options *acm.Options)) (*acm.RemoveTagsFromCertificateOutput, error)
}

type mockDescribeCertificateAPI func(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(options *acm.Options)) (*acm.DescribeCertificateOutput, error)

func (m mockDescribeCertificateAPI) DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(options *acm.Options)) (*acm.DescribeCertificateOutput, error){
	return m(ctx, params, optFns...)
}

type mockRequestCertificateAPI func(ctx context.Context, params *acm.RequestCertificateInput, optFns ...func(options *acm.Options)) (*acm.RequestCertificateOutput, error)

func (m mockRequestCertificateAPI) RequestCertificate(ctx context.Context, params *acm.RequestCertificateInput, optFns ...func(options *acm.Options)) (*acm.RequestCertificateOutput, error){
	return m(ctx, params, optFns...)
}

type mockDeleteCertificateAPI func(ctx context.Context, params *acm.DeleteCertificateInput, optFns ...func(options *acm.Options)) (*acm.DeleteCertificateOutput, error)

func (m mockDeleteCertificateAPI) DeleteCertificateCertificate(ctx context.Context, params *acm.DeleteCertificateInput, optFns ...func(options *acm.Options)) (*acm.DeleteCertificateOutput, error){
	return m(ctx, params, optFns...)
}

type mockUpdateCertificateOptionsAPI func(ctx context.Context, params *acm.UpdateCertificateOptionsInput, optFns ...func(options *acm.Options)) (*acm.UpdateCertificateOptionsOutput, error)

func (m mockUpdateCertificateOptionsAPI) UpdateCertificateOptions(ctx context.Context, params *acm.UpdateCertificateOptionsInput, optFns ...func(options *acm.Options)) (*acm.UpdateCertificateOptionsOutput, error){
	return m(ctx, params, optFns...)
}