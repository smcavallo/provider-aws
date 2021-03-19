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
var _ clientset.CAPermissionClient = (*MockCertificateAuthorityPermissionClient)(nil)

// MockCertificateAuthorityPermissionClient is a type that implements all the methods for Certificate Authority Permission Client interface
type MockCertificateAuthorityPermissionClient interface {

	CreatePermission(ctx context.Context, params *acmpca.CreatePermissionInput, optFns ...func(options *acmpca.Options)) (*acmpca.CreatePermissionOutput, error)
	DeletePermission(ctx context.Context, params *acmpca.DeletePermissionInput, optFns ...func(options *acmpca.Options)) (*acmpca.DeletePermissionOutput, error)
	ListPermissions(ctx context.Context, params *acmpca.ListPermissionsInput, optFns ...func(options *acmpca.Options)) (*acmpca.ListPermissionsOutput, error)

}

type mockCreatePermissionAPI func(ctx context.Context, params *acmpca.CreatePermissionInput, optFns ...func(options *acmpca.Options)) (*acmpca.CreatePermissionOutput, error)

func (m mockCreatePermissionAPI) CreatePermission(ctx context.Context, params *acmpca.CreatePermissionInput, optFns ...func(options *acmpca.Options)) (*acmpca.CreatePermissionOutput, error){
	return m(ctx, params, optFns...)
}

type mockDeletePermissionAPI func(ctx context.Context, params *acmpca.DeletePermissionInput, optFns ...func(options *acmpca.Options)) (*acmpca.DeletePermissionOutput, error)

func (m mockDeletePermissionAPI) DeletePermission(ctx context.Context, params *acmpca.DeletePermissionInput, optFns ...func(options *acmpca.Options)) (*acmpca.DeletePermissionOutput, error){
	return m(ctx, params, optFns...)
}

type mockListPermissionsAPI func(ctx context.Context, params *acmpca.ListPermissionsInput, optFns ...func(options *acmpca.Options)) (*acmpca.ListPermissionsOutput, error)

func (m mockListPermissionsAPI) ListPermissions(ctx context.Context, params *acmpca.ListPermissionsInput, optFns ...func(options *acmpca.Options)) (*acmpca.ListPermissionsOutput, error){
	return m(ctx, params, optFns...)
}