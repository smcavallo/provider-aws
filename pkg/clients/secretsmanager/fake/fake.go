/*
Copyright 2020 The Crossplane Authors.

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

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// MockClient is a fake implementation of secretsmanager.Client.
type MockClient struct {
	MockDescribeSecret func(ctx context.Context, input *secretsmanager.DescribeSecretInput, opts []func(*secretsmanager.Options)) (*secretsmanager.DescribeSecretOutput, error)
	MockGetSecretValue func(ctx context.Context, input *secretsmanager.GetSecretValueInput, opts []func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
	MockCreateSecret   func(ctx context.Context, input *secretsmanager.CreateSecretInput, opts []func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error)
	MockDeleteSecret   func(ctx context.Context, input *secretsmanager.DeleteSecretInput, opts []func(*secretsmanager.Options)) (*secretsmanager.DeleteSecretOutput, error)
	MockUpdateSecret   func(ctx context.Context, input *secretsmanager.UpdateSecretInput, opts []func(*secretsmanager.Options)) (*secretsmanager.UpdateSecretOutput, error)
	MockTagResource    func(ctx context.Context, input *secretsmanager.TagResourceInput, opts []func(*secretsmanager.Options)) (*secretsmanager.TagResourceOutput, error)
	MockUntagResource  func(ctx context.Context, input *secretsmanager.UntagResourceInput, opts []func(*secretsmanager.Options)) (*secretsmanager.UntagResourceOutput, error)
}

// DescribeSecret calls the underlying MockDescribeSecret method.
func (c *MockClient) DescribeSecret(ctx context.Context, input *secretsmanager.DescribeSecretInput, opts ...func(*secretsmanager.Options)) (*secretsmanager.DescribeSecretOutput, error) {
	return c.MockDescribeSecret(ctx, input, opts)
}

// GetSecretValue calls the underlying MockGetSecretValue method.
func (c *MockClient) GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, opts ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return c.MockGetSecretValue(ctx, input, opts)
}

// CreateSecret calls the underlying MockCreateSecret method.
func (c *MockClient) CreateSecret(ctx context.Context, input *secretsmanager.CreateSecretInput, opts ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error) {
	return c.MockCreateSecret(ctx, input, opts)
}

// DeleteSecret calls the underlying MockDeleteSecret method.
func (c *MockClient) DeleteSecret(ctx context.Context, input *secretsmanager.DeleteSecretInput, opts ...func(*secretsmanager.Options)) (*secretsmanager.DeleteSecretOutput, error) {
	return c.MockDeleteSecret(ctx, input, opts)
}

// UpdateSecret calls the underlying MockUpdateSecret method.
func (c *MockClient) UpdateSecret(ctx context.Context, input *secretsmanager.UpdateSecretInput, opts ...func(*secretsmanager.Options)) (*secretsmanager.UpdateSecretOutput, error) {
	return c.MockUpdateSecret(ctx, input, opts)
}

// TagResource calls the underlying MockTagResource method.
func (c *MockClient) TagResource(ctx context.Context, input *secretsmanager.TagResourceInput, opts ...func(*secretsmanager.Options)) (*secretsmanager.TagResourceOutput, error) {
	return c.MockTagResource(ctx, input, opts)
}

// UntagResource calls the underlying UntagResource method.
func (c *MockClient) UntagResource(ctx context.Context, input *secretsmanager.UntagResourceInput, opts ...func(*secretsmanager.Options)) (*secretsmanager.UntagResourceOutput, error) {
	return c.MockUntagResource(ctx, input, opts)
}
