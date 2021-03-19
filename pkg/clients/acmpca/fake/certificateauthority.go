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

	"github.com/aws/aws-sdk-go-v2/service/acmpca"

	clientset "github.com/crossplane/provider-aws/pkg/clients/acmpca"
)

// this ensures that the mock implements the client interface
var _ clientset.Client = (*MockCertificateAuthorityClient)(nil)

// MockCertificateAuthorityClient is a type that implements all the methods for Certificate Authority Client interface
type MockCertificateAuthorityClient interface {
	CreateCertificateAuthority(ctx context.Context, params *acmpca.CreateCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.CreateCertificateAuthorityOutput, error)
	DeleteCertificateAuthority(ctx context.Context, params *acmpca.DeleteCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.DeleteCertificateAuthorityOutput, error)
	UpdateCertificateAuthority(ctx context.Context, params *acmpca.UpdateCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.UpdateCertificateAuthorityOutput, error)
	DescribeCertificateAuthority(ctx context.Context, params *acmpca.DescribeCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.DescribeCertificateAuthorityOutput, error)
	ListTags(ctx context.Context, params *acmpca.ListTagsInput, optFns ...func(options *acmpca.Options)) (*acmpca.ListTagsOutput, error)
	TagCertificateAuthority(ctx context.Context, params *acmpca.TagCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.TagCertificateAuthorityOutput, error)
	UntagCertificateAuthority(ctx context.Context, params *acmpca.UntagCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.UntagCertificateAuthorityOutput, error)
}

type mockCreateCertificateAuthorityAPI func(ctx context.Context, params *acmpca.CreateCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.CreateCertificateAuthorityOutput, error)

func (m mockCreateCertificateAuthorityAPI) CreateCertificateAuthority(ctx context.Context, params *acmpca.CreateCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.CreateCertificateAuthorityOutput, error){
	return m(ctx, params, optFns...)
}

type mockDeleteCertificateAuthorityAPI func(ctx context.Context, params *acmpca.DeleteCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.DeleteCertificateAuthorityOutput, error)

func (m mockDeleteCertificateAuthorityAPI) DeleteCertificateAuthority(ctx context.Context, params *acmpca.DeleteCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.DeleteCertificateAuthorityOutput, error){
	return m(ctx, params, optFns...)
}
type mockUpdateCertificateAuthorityAPI func(ctx context.Context, params *acmpca.UpdateCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.UpdateCertificateAuthorityOutput, error)

func (m mockUpdateCertificateAuthorityAPI) UpdateCertificateAuthority(ctx context.Context, params *acmpca.UpdateCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.UpdateCertificateAuthorityOutput, error){
	return m(ctx, params, optFns...)
}
type mockDescribeCertificateAuthorityAPI func(ctx context.Context, params *acmpca.DescribeCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.DescribeCertificateAuthorityOutput, error)

func (m mockDescribeCertificateAuthorityAPI) DescribeCertificateAuthority(ctx context.Context, params *acmpca.DescribeCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.DescribeCertificateAuthorityOutput, error){
	return m(ctx, params, optFns...)
}

type mockListTagsAPI func(ctx context.Context, params *acmpca.ListTagsInput, optFns ...func(options *acmpca.Options)) (*acmpca.ListTagsOutput, error)

func (m mockListTagsAPI) ListTags(ctx context.Context, params *acmpca.ListTagsInput, optFns ...func(options *acmpca.Options)) (*acmpca.ListTagsOutput, error){
	return m(ctx, params, optFns...)
}

type mockTagCertificateAuthorityAPI func(ctx context.Context, params *acmpca.TagCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.TagCertificateAuthorityOutput, error)

func (m mockTagCertificateAuthorityAPI) TagCertificateAuthority(ctx context.Context, params *acmpca.TagCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.TagCertificateAuthorityOutput, error){
	return m(ctx, params, optFns...)
}

type mockUntagCertificateAuthorityAPI func(ctx context.Context, params *acmpca.UntagCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.UntagCertificateAuthorityOutput, error)

func (m mockUntagCertificateAuthorityAPI) UntagCertificateAuthority(ctx context.Context, params *acmpca.UntagCertificateAuthorityInput, optFns ...func(options *acmpca.Options)) (*acmpca.UntagCertificateAuthorityOutput, error){
	return m(ctx, params, optFns...)
}