package acm

import (
	"errors"
	"strings"
    "context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"

	"github.com/crossplane/provider-aws/apis/acm/v1alpha1"
	awsclients "github.com/crossplane/provider-aws/pkg/clients"
)

// Client defines the CertificateManager operations
type Client interface {
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

// NewClient returns a new client using AWS credentials as JSON encoded data.
func NewClient(conf aws.Config) Client {
	return acm.New(conf)
}

// GenerateCreateCertificateInput from CertificateSpec
func GenerateCreateCertificateInput(name string, p *v1alpha1.CertificateParameters) *acm.RequestCertificateInput {

	m := &acm.RequestCertificateInput{
		DomainName:              aws.String(p.DomainName),
		CertificateAuthorityArn: p.CertificateAuthorityARN,
	}

	if p.CertificateTransparencyLoggingPreference != nil {
		m.Options = &types.CertificateOptions{CertificateTransparencyLoggingPreference: *p.CertificateTransparencyLoggingPreference}
	}

	if p.ValidationMethod != nil {
		m.ValidationMethod = *p.ValidationMethod
	}

	if len(p.DomainValidationOptions) != 0 {
		m.DomainValidationOptions = make([]*types.DomainValidationOption, len(p.DomainValidationOptions))
		for i, val := range p.DomainValidationOptions {
			m.DomainValidationOptions[i] = &types.DomainValidationOption{
				DomainName:       aws.String(val.DomainName),
				ValidationDomain: aws.String(val.ValidationDomain),
			}
		}
	}

	if p.SubjectAlternativeNames != nil {
		m.SubjectAlternativeNames = make([]*string, len(p.SubjectAlternativeNames))
		for i := range p.SubjectAlternativeNames {
			m.SubjectAlternativeNames[i] = p.SubjectAlternativeNames[i]
		}
	}

	m.Tags = make([]*types.Tag, len(p.Tags))
	for i, val := range p.Tags {
		m.Tags[i] = &types.Tag{
			Key:   aws.String(val.Key),
			Value: aws.String(val.Value),
		}
	}
	return m
}

// GenerateCertificateStatus is used to produce CertificateExternalStatus from acm.certificateStatus
func GenerateCertificateStatus(certificate types.CertificateDetail) v1alpha1.CertificateExternalStatus {
	return v1alpha1.CertificateExternalStatus{
		CertificateARN:     aws.StringValue(certificate.CertificateArn),
		RenewalEligibility: certificate.RenewalEligibility,
		Status:             certificate.Status,
		Type:               certificate.Type,
	}
}

// LateInitializeCertificate fills the empty fields in *v1beta1.CertificateParameters with
// the values seen in iam.Certificate.
func LateInitializeCertificate(in *v1alpha1.CertificateParameters, certificate *types.CertificateDetail) { // nolint:gocyclo
	if certificate == nil {
		return
	}

	in.DomainName = awsclients.LateInitializeString(in.DomainName, certificate.DomainName)

	if aws.StringValue(in.CertificateAuthorityARN) == "" && certificate.CertificateAuthorityArn != nil {
		in.CertificateAuthorityARN = certificate.CertificateAuthorityArn
	}

	if in.CertificateTransparencyLoggingPreference == nil && certificate.Options != nil {
		in.CertificateTransparencyLoggingPreference = &certificate.Options.CertificateTransparencyLoggingPreference
	}

	if in.ValidationMethod == nil && len(certificate.DomainValidationOptions) != 0 {
		in.ValidationMethod = &certificate.DomainValidationOptions[0].ValidationMethod
	}

	if len(in.SubjectAlternativeNames) == 0 && len(certificate.SubjectAlternativeNames) != 0 {
		in.SubjectAlternativeNames = make([]*string, len(certificate.SubjectAlternativeNames))
		for i := range certificate.SubjectAlternativeNames {
			in.SubjectAlternativeNames[i] = certificate.SubjectAlternativeNames[i]
		}
	}

	if len(in.DomainValidationOptions) == 0 && len(certificate.DomainValidationOptions) != 0 {
		in.DomainValidationOptions = make([]*v1alpha1.DomainValidationOption, len(certificate.DomainValidationOptions))
		for i, val := range certificate.DomainValidationOptions {
			in.DomainValidationOptions[i] = &v1alpha1.DomainValidationOption{
				DomainName:       aws.StringValue(val.DomainName),
				ValidationDomain: aws.StringValue(val.ValidationDomain),
			}
		}
	}

}

// IsCertificateUpToDate checks whether there is a change in any of the modifiable fields.
func IsCertificateUpToDate(p v1alpha1.CertificateParameters, cd types.CertificateDetail, tags []types.Tag) bool { // nolint:gocyclo

	if *p.CertificateTransparencyLoggingPreference != cd.Options.CertificateTransparencyLoggingPreference {
		return false
	}

	if len(p.Tags) != len(tags) {
		return false
	}

	pTags := make(map[string]string, len(p.Tags))
	for _, tag := range p.Tags {
		pTags[tag.Key] = tag.Value
	}
	for _, tag := range tags {
		val, ok := pTags[aws.StringValue(tag.Key)]
		if !ok || !strings.EqualFold(val, aws.StringValue(tag.Value)) {
			return false
		}
	}

	return !aws.BoolValue(p.RenewCertificate)
}

// IsErrorNotFound returns true if the error code indicates that the item was not found
func IsErrorNotFound(err error) bool {
	if err != nil {
		var rne *types.ResourceNotFoundException
		if errors.As(err, &rne) {
			return true
		}
	}
	return false
}
