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

package s3

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/crossplane/crossplane-runtime/pkg/meta"

	"github.com/crossplane/provider-aws/apis/s3/v1beta1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
)

var (
	// CORSErrCode is the error code sent by AWS when the CORS configuration does not exist
	CORSErrCode = "NoSuchCORSConfiguration"
	// ReplicationErrCode is the error code sent by AWS when the replication config does not exist
	ReplicationErrCode = "ReplicationConfigurationNotFoundError"
	// LifecycleErrCode is the error code sent by AWS when the lifecycle config does not exist
	LifecycleErrCode = "NoSuchLifecycleConfiguration"
	// SSEErrCode is the error code sent by AWS when the SSE config does not exist
	SSEErrCode = "ServerSideEncryptionConfigurationNotFoundError"
	// TaggingErrCode is the error code sent by AWS when the tagging does not exist
	TaggingErrCode = "NoSuchTagSet"
	// WebsiteErrCode is the error code sent by AWS when the website config does not exist
	WebsiteErrCode = "NoSuchWebsiteConfiguration"
	// MethodNotAllowed is the error code sent by AWS when the request method for an object is not allowed
	MethodNotAllowed = "MethodNotAllowed"
	// UnsupportedArgument is the error code sent by AWS when the request fields contain an argument that is not supported
	UnsupportedArgument = "UnsupportedArgument"
)

// BucketClient is the interface for Client for making S3 Bucket requests.
type BucketClient interface {
	HeadBucket(ctx context.Context, input *s3.HeadBucketInput, opts ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
	CreateBucket(ctx context.Context, input *s3.CreateBucketInput, opts ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	DeleteBucket(ctx context.Context, input *s3.DeleteBucketInput, opts ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)

	PutBucketEncryption(ctx context.Context, input *s3.PutBucketEncryptionInput, opts ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error)
	GetBucketEncryption(ctx context.Context, input *s3.GetBucketEncryptionInput, opts ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error)
	DeleteBucketEncryption(ctx context.Context, input *s3.DeleteBucketEncryptionInput, opts ...func(*s3.Options)) (*s3.DeleteBucketEncryptionOutput, error)

	PutBucketVersioning(ctx context.Context, input *s3.PutBucketVersioningInput, opts ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error)
	GetBucketVersioning(ctx context.Context, input *s3.GetBucketVersioningInput, opts ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)

	PutBucketAccelerateConfiguration(ctx context.Context, input *s3.PutBucketAccelerateConfigurationInput, opts ...func(*s3.Options)) (*s3.PutBucketAccelerateConfigurationOutput, error)
	GetBucketAccelerateConfiguration(ctx context.Context, input *s3.GetBucketAccelerateConfigurationInput, opts ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error)

	PutBucketCors(ctx context.Context, input *s3.PutBucketCorsInput, opts ...func(*s3.Options)) (*s3.PutBucketCorsOutput, error)
	GetBucketCors(ctx context.Context, input *s3.GetBucketCorsInput, opts ...func(*s3.Options)) (*s3.GetBucketCorsOutput, error)
	DeleteBucketCors(ctx context.Context, input *s3.DeleteBucketCorsInput, opts ...func(*s3.Options)) (*s3.DeleteBucketCorsOutput, error)

	PutBucketWebsite(ctx context.Context, input *s3.PutBucketWebsiteInput, opts ...func(*s3.Options)) (*s3.PutBucketWebsiteOutput, error)
	GetBucketWebsite(ctx context.Context, input *s3.GetBucketWebsiteInput, opts ...func(*s3.Options)) (*s3.GetBucketWebsiteOutput, error)
	DeleteBucketWebsite(ctx context.Context, input *s3.DeleteBucketWebsiteInput, opts ...func(*s3.Options)) (*s3.DeleteBucketWebsiteOutput, error)

	PutBucketLogging(ctx context.Context, input *s3.PutBucketLoggingInput, opts ...func(*s3.Options)) (*s3.PutBucketLoggingOutput, error)
	GetBucketLogging(ctx context.Context, input *s3.GetBucketLoggingInput, opts ...func(*s3.Options)) (*s3.GetBucketLoggingOutput, error)

	PutBucketReplication(ctx context.Context, input *s3.PutBucketReplicationInput, opts ...func(*s3.Options)) (*s3.PutBucketReplicationOutput, error)
	GetBucketReplication(ctx context.Context, input *s3.GetBucketReplicationInput, opts ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error)
	DeleteBucketReplication(ctx context.Context, input *s3.DeleteBucketReplicationInput, opts ...func(*s3.Options)) (*s3.DeleteBucketReplicationOutput, error)

	PutBucketRequestPayment(ctx context.Context, input *s3.PutBucketRequestPaymentInput, opts ...func(*s3.Options)) (*s3.PutBucketRequestPaymentOutput, error)
	GetBucketRequestPayment(ctx context.Context, input *s3.GetBucketRequestPaymentInput, opts ...func(*s3.Options)) (*s3.GetBucketRequestPaymentOutput, error)

	PutBucketTagging(ctx context.Context, input *s3.PutBucketTaggingInput, opts ...func(*s3.Options)) (*s3.PutBucketTaggingOutput, error)
	GetBucketTagging(ctx context.Context, input *s3.GetBucketTaggingInput, opts ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
	DeleteBucketTagging(ctx context.Context, input *s3.DeleteBucketTaggingInput, opts ...func(*s3.Options)) (*s3.DeleteBucketTaggingOutput, error)

	PutBucketAnalyticsConfiguration(ctx context.Context, input *s3.PutBucketAnalyticsConfigurationInput, opts ...func(*s3.Options)) (*s3.PutBucketAnalyticsConfigurationOutput, error)
	GetBucketAnalyticsConfiguration(ctx context.Context, input *s3.GetBucketAnalyticsConfigurationInput, opts ...func(*s3.Options)) (*s3.GetBucketAnalyticsConfigurationOutput, error)

	PutBucketLifecycleConfiguration(ctx context.Context, input *s3.PutBucketLifecycleConfigurationInput, opts ...func(*s3.Options)) (*s3.PutBucketLifecycleConfigurationOutput, error)
	GetBucketLifecycleConfiguration(ctx context.Context, input *s3.GetBucketLifecycleConfigurationInput, opts ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error)
	DeleteBucketLifecycle(ctx context.Context, input *s3.DeleteBucketLifecycleInput, opts ...func(*s3.Options)) (*s3.DeleteBucketLifecycleOutput, error)

	PutBucketNotificationConfiguration(ctx context.Context, input *s3.PutBucketNotificationConfigurationInput, opts ...func(*s3.Options)) (*s3.PutBucketNotificationConfigurationOutput, error)
	GetBucketNotificationConfiguration(ctx context.Context, input *s3.GetBucketNotificationConfigurationInput, opts ...func(*s3.Options)) (*s3.GetBucketNotificationConfigurationOutput, error)

	GetBucketAcl(ctx context.Context, input *s3.GetBucketAclInput, opts ...func(*s3.Options)) (*s3.GetBucketAclOutput, error) //nolint
	PutBucketAcl(ctx context.Context, input *s3.PutBucketAclInput, opts ...func(*s3.Options)) (*s3.PutBucketAclOutput, error) //nolint
}

// NewClient returns a new client using AWS credentials as JSON encoded data.
func NewClient(cfg aws.Config) BucketClient {
	return s3.NewFromConfig(cfg)
}

// IsNotFound helper function to test for NotFound error
func IsNotFound(err error) bool {
	var nsb *s3types.NoSuchBucket
	return errors.As(err, &nsb)
}

// IsAlreadyExists helper function to test for ErrCodeBucketAlreadyOwnedByYou error
func IsAlreadyExists(err error) bool {
	var nsb *s3types.BucketAlreadyOwnedByYou
	return errors.As(err, &nsb)
}

// GenerateCreateBucketInput creates the input for CreateBucket S3 Client request
func GenerateCreateBucketInput(name string, s v1beta1.BucketParameters) *s3.CreateBucketInput {
	cbi := &s3.CreateBucketInput{
		ACL:                        s3types.BucketCannedACL(aws.ToString(s.ACL)),
		Bucket:                     aws.String(name),
		GrantFullControl:           s.GrantFullControl,
		GrantRead:                  s.GrantRead,
		GrantReadACP:               s.GrantReadACP,
		GrantWrite:                 s.GrantWrite,
		GrantWriteACP:              s.GrantWriteACP,
		ObjectLockEnabledForBucket: aws.ToBool(s.ObjectLockEnabledForBucket),
	}
	if s.LocationConstraint != "us-east-1" {
		cbi.CreateBucketConfiguration = &s3types.CreateBucketConfiguration{LocationConstraint: s3types.BucketLocationConstraint(s.LocationConstraint)}
	}
	return cbi
}

// GenerateBucketObservation generates the ARN string for the external status
func GenerateBucketObservation(name string) v1beta1.BucketExternalStatus {
	return v1beta1.BucketExternalStatus{
		ARN: fmt.Sprintf("arn:aws:s3:::%s", name),
	}
}

// CORSConfigurationNotFound is parses the aws Error and validates if the cors configuration does not exist
func CORSConfigurationNotFound(err error) bool {
	if awsErr, ok := err.(smithy.APIError); ok {
		if awsErr.ErrorCode() == CORSErrCode {
			return true
		}
	}
	return false
}

// ReplicationConfigurationNotFound is parses the aws Error and validates if the replication configuration does not exist
func ReplicationConfigurationNotFound(err error) bool {
	if awsErr, ok := err.(smithy.APIError); ok {
		if awsErr.ErrorCode() == ReplicationErrCode {
			return true
		}
	}
	return false
}

// LifecycleConfigurationNotFound is parses the aws Error and validates if the lifecycle configuration does not exist
func LifecycleConfigurationNotFound(err error) bool {
	if awsErr, ok := err.(smithy.APIError); ok {
		if awsErr.ErrorCode() == LifecycleErrCode {
			return true
		}
	}
	return false
}

// SSEConfigurationNotFound is parses the aws Error and validates if the SSE configuration does not exist
func SSEConfigurationNotFound(err error) bool {
	if awsErr, ok := err.(smithy.APIError); ok {
		if awsErr.ErrorCode() == SSEErrCode {
			return true
		}
	}
	return false
}

// TaggingNotFound is parses the aws Error and validates if the tagging configuration does not exist
func TaggingNotFound(err error) bool {
	if awsErr, ok := err.(smithy.APIError); ok {
		if awsErr.ErrorCode() == TaggingErrCode {
			return true
		}
	}
	return false
}

// WebsiteConfigurationNotFound is parses the aws Error and validates if the website configuration does not exist
func WebsiteConfigurationNotFound(err error) bool {
	if awsErr, ok := err.(smithy.APIError); ok {
		if awsErr.ErrorCode() == WebsiteErrCode {
			return true
		}
	}
	return false
}

// MethodNotSupported is parses the aws Error and validates if the method is allowed for a request
func MethodNotSupported(err error) bool {
	if awsErr, ok := err.(smithy.APIError); ok {
		if awsErr.ErrorCode() == MethodNotAllowed {
			return true
		}
	}
	return false
}

// ArgumentNotSupported is parses the aws Error and validates if parameters are now allowed for a request
func ArgumentNotSupported(err error) bool {
	if awsErr, ok := err.(smithy.APIError); ok {
		if awsErr.ErrorCode() == UnsupportedArgument {
			return true
		}
	}
	return false
}

// UpdateBucketACL creates the ACLInput, sends the request to put an ACL based on the bucket
func UpdateBucketACL(ctx context.Context, client BucketClient, bucket *v1beta1.Bucket) error {
	config := &s3.PutBucketAclInput{
		ACL:              s3types.BucketCannedACL(aws.ToString(bucket.Spec.ForProvider.ACL)),
		Bucket:           aws.String(meta.GetExternalName(bucket)),
		GrantFullControl: bucket.Spec.ForProvider.GrantFullControl,
		GrantRead:        bucket.Spec.ForProvider.GrantRead,
		GrantReadACP:     bucket.Spec.ForProvider.GrantReadACP,
		GrantWrite:       bucket.Spec.ForProvider.GrantWrite,
		GrantWriteACP:    bucket.Spec.ForProvider.GrantWriteACP,
	}
	_, err := client.PutBucketAcl(ctx, config)
	return err
}

// CopyTags converts a list of local v1beta.Tags to S3 Tags
func CopyTags(tags []v1beta1.Tag) []s3types.Tag {
	out := make([]s3types.Tag, 0)
	for _, one := range tags {
		out = append(out, s3types.Tag{Key: aws.String(one.Key), Value: aws.String(one.Value)})
	}
	return out
}

// CopyAWSTags converts a list of external s3.Tags to local Tags
func CopyAWSTags(tags []s3types.Tag) []v1beta1.Tag {
	out := make([]v1beta1.Tag, len(tags))
	for i, one := range tags {
		out[i] = v1beta1.Tag{Key: awsclient.StringValue(one.Key), Value: awsclient.StringValue(one.Value)}
	}
	return out
}

// SortS3TagSet stable sorts an external s3 tag list by the key and value.
func SortS3TagSet(tags []s3types.Tag) []s3types.Tag {
	outTags := make([]s3types.Tag, len(tags))
	copy(outTags, tags)
	sort.SliceStable(outTags, func(i, j int) bool {
		return aws.ToString(outTags[i].Key) < aws.ToString(outTags[j].Key)
	})
	return outTags
}
