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

package bucket

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	awss3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/crossplane/provider-aws/apis/s3/v1beta1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
	"github.com/crossplane/provider-aws/pkg/clients/s3"
)

const (
	lifecycleGetFailed    = "cannot get Bucket lifecycle configuration"
	lifecyclePutFailed    = "cannot put Bucket lifecycle configuration"
	lifecycleDeleteFailed = "cannot delete Bucket lifecycle configuration"
)

// LifecycleConfigurationClient is the client for API methods and reconciling the LifecycleConfiguration
type LifecycleConfigurationClient struct {
	client s3.BucketClient
}

// LateInitialize does nothing because LifecycleConfiguration might have been be
// deleted by the user.
func (*LifecycleConfigurationClient) LateInitialize(_ context.Context, _ *v1beta1.Bucket) error {
	return nil
}

// NewLifecycleConfigurationClient creates the client for Accelerate Configuration
func NewLifecycleConfigurationClient(client s3.BucketClient) *LifecycleConfigurationClient {
	return &LifecycleConfigurationClient{client: client}
}

// Observe checks if the resource exists and if it matches the local configuration
func (in *LifecycleConfigurationClient) Observe(ctx context.Context, bucket *v1beta1.Bucket) (ResourceStatus, error) {
	response, err := in.client.GetBucketLifecycleConfiguration(ctx, &awss3.GetBucketLifecycleConfigurationInput{Bucket: awsclient.String(meta.GetExternalName(bucket))})
	if bucket.Spec.ForProvider.LifecycleConfiguration == nil && s3.LifecycleConfigurationNotFound(err) {
		return Updated, nil
	}
	if resource.Ignore(s3.LifecycleConfigurationNotFound, err) != nil {
		return NeedsUpdate, awsclient.Wrap(err, lifecycleGetFailed)
	}
	var local []v1beta1.LifecycleRule
	if bucket.Spec.ForProvider.LifecycleConfiguration != nil {
		local = bucket.Spec.ForProvider.LifecycleConfiguration.Rules
	}
	var external []awss3types.LifecycleRule
	if response != nil {
		external = response.Rules
	}
	sortFilterTags(external)
	switch {
	case len(external) != 0 && len(local) == 0:
		return NeedsDeletion, nil
	// NOTE(muvaf): We ignore ID because it might have been auto-assigned by AWS
	// and we don't have late-init for this subresource. Besides, a change in ID
	// is almost never expected.
	case cmp.Equal(external, GenerateLifecycleRules(local),
		cmpopts.IgnoreFields(awss3types.LifecycleRule{}, "ID")):
		return Updated, nil
	default:
		return NeedsUpdate, nil
	}
}

// GenerateLifecycleConfiguration creates the PutBucketLifecycleConfigurationInput for the AWS SDK
func GenerateLifecycleConfiguration(name string, config *v1beta1.BucketLifecycleConfiguration) *awss3.PutBucketLifecycleConfigurationInput {
	if config == nil {
		return nil
	}
	return &awss3.PutBucketLifecycleConfigurationInput{
		Bucket:                 awsclient.String(name),
		LifecycleConfiguration: &awss3types.BucketLifecycleConfiguration{Rules: GenerateLifecycleRules(config.Rules)},
	}
}

// GenerateLifecycleRules creates the list of LifecycleRules for the AWS SDK
func GenerateLifecycleRules(in []v1beta1.LifecycleRule) []awss3types.LifecycleRule { // nolint:gocyclo
	// NOTE(muvaf): prealloc is disabled due to AWS requiring nil instead
	// of 0-length for empty slices.
	var result []awss3types.LifecycleRule // nolint:prealloc
	for _, local := range in {
		rule := awss3types.LifecycleRule{
			ID:     local.ID,
			Status: awss3types.ExpirationStatus(local.Status),
		}
		if local.AbortIncompleteMultipartUpload != nil {
			rule.AbortIncompleteMultipartUpload = &awss3types.AbortIncompleteMultipartUpload{
				DaysAfterInitiation: local.AbortIncompleteMultipartUpload.DaysAfterInitiation,
			}
		}
		if local.Expiration != nil {
			rule.Expiration = &awss3types.LifecycleExpiration{
				Days:                      aws.ToInt32(local.Expiration.Days),
				ExpiredObjectDeleteMarker: aws.ToBool(local.Expiration.ExpiredObjectDeleteMarker),
			}
			if local.Expiration.Date != nil {
				rule.Expiration.Date = &local.Expiration.Date.Time
			}
		}
		if local.NoncurrentVersionExpiration != nil {
			rule.NoncurrentVersionExpiration = &awss3types.NoncurrentVersionExpiration{NoncurrentDays: aws.ToInt32(local.NoncurrentVersionExpiration.NoncurrentDays)}
		}
		if local.NoncurrentVersionTransitions != nil {
			rule.NoncurrentVersionTransitions = make([]awss3types.NoncurrentVersionTransition, len(local.NoncurrentVersionTransitions))
			for tIndex, transition := range local.NoncurrentVersionTransitions {
				rule.NoncurrentVersionTransitions[tIndex] = awss3types.NoncurrentVersionTransition{
					NoncurrentDays: aws.ToInt32(transition.NoncurrentDays),
					StorageClass:   awss3types.TransitionStorageClass(transition.StorageClass),
				}
			}
		}
		if local.Transitions != nil {
			rule.Transitions = make([]awss3types.Transition, len(local.Transitions))
			for tIndex, transition := range local.Transitions {
				rule.Transitions[tIndex] = awss3types.Transition{
					Days:         aws.ToInt32(transition.Days),
					StorageClass: awss3types.TransitionStorageClass(transition.StorageClass),
				}
				if transition.Date != nil {
					rule.Transitions[tIndex].Date = &transition.Date.Time
				}
			}
		}
		if local.Filter != nil {
			if local.Filter.Prefix != nil {
				rule.Filter = &awss3types.LifecycleRuleFilterMemberPrefix{Value: *local.Filter.Prefix}
			}
			if local.Filter.Tag != nil {
				rule.Filter = &awss3types.LifecycleRuleFilterMemberTag{Value: awss3types.Tag{Key: awsclient.String(local.Filter.Tag.Key), Value: awsclient.String(local.Filter.Tag.Value)}}
			}
			if local.Filter.And != nil {
				andOperator := awss3types.LifecycleRuleAndOperator{
					Prefix: local.Filter.And.Prefix,
				}
				if local.Filter.And.Tags != nil {
					andOperator.Tags = s3.SortS3TagSet(s3.CopyTags(local.Filter.And.Tags))
				}
				rule.Filter = &awss3types.LifecycleRuleFilterMemberAnd{Value: andOperator}
			}
		}
		result = append(result, rule)
	}
	return result
}

// CreateOrUpdate sends a request to have resource created on AWS
func (in *LifecycleConfigurationClient) CreateOrUpdate(ctx context.Context, bucket *v1beta1.Bucket) error {
	if bucket.Spec.ForProvider.LifecycleConfiguration == nil {
		return nil
	}
	input := GenerateLifecycleConfiguration(meta.GetExternalName(bucket), bucket.Spec.ForProvider.LifecycleConfiguration)
	_, err := in.client.PutBucketLifecycleConfiguration(ctx, input)
	return awsclient.Wrap(err, lifecyclePutFailed)

}

// Delete creates the request to delete the resource on AWS or set it to the default value.
func (in *LifecycleConfigurationClient) Delete(ctx context.Context, bucket *v1beta1.Bucket) error {
	_, err := in.client.DeleteBucketLifecycle(ctx,
		&awss3.DeleteBucketLifecycleInput{
			Bucket: awsclient.String(meta.GetExternalName(bucket)),
		},
	)
	return awsclient.Wrap(err, lifecycleDeleteFailed)
}

func sortFilterTags(rules []awss3types.LifecycleRule) {
	for i := range rules {
		andOperator, ok := rules[i].Filter.(*awss3types.LifecycleRuleFilterMemberAnd)
		if ok {
			andOperator.Value.Tags = s3.SortS3TagSet(andOperator.Value.Tags)
		}
	}
}
