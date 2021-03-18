package acm

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/google/go-cmp/cmp"

	"github.com/crossplane/provider-aws/apis/acm/v1alpha1"
	aws "github.com/crossplane/provider-aws/pkg/clients"
)

var (
	domainName              = "somedomain"
	certificateArn          = "somearn"
	certificateAuthorityArn = "someauthorityarn"
)

func TestGenerateCreateCertificateInput(t *testing.T) {
	certificateTransparencyLoggingPreference := types.CertificateTransparencyLoggingPreferenceDisabled
	validationMethod := types.ValidationMethodDns
	cases := map[string]struct {
		in  v1alpha1.CertificateParameters
		out acm.RequestCertificateInput
	}{
		"FilledInput": {
			in: v1alpha1.CertificateParameters{
				DomainName:                               domainName,
				CertificateAuthorityARN:                  aws.String(certificateAuthorityArn),
				CertificateTransparencyLoggingPreference: &certificateTransparencyLoggingPreference,
				ValidationMethod:                         &validationMethod,
				Tags: []v1alpha1.Tag{{
					Key:   "key1",
					Value: "value1",
				}},
			},
			out: acm.RequestCertificateInput{
				DomainName:              aws.String(domainName),
				CertificateAuthorityArn: aws.String(certificateAuthorityArn),
				Options:                 &types.CertificateOptions{CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled},
				ValidationMethod:        types.ValidationMethodDns,
				Tags: []*types.Tag{{
					Key:   aws.String("key1"),
					Value: aws.String("value1"),
				}},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r := GenerateCreateCertificateInput(name, &tc.in)

			if diff := cmp.Diff(r, &tc.out); diff != "" {
				t.Errorf("GenerateCreateCertificateInput(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestLateInitializeCertificate(t *testing.T) {
	certificateTransparencyLoggingPreference := types.CertificateTransparencyLoggingPreferenceDisabled
	type args struct {
		spec *v1alpha1.CertificateParameters
		in   *types.CertificateDetail
	}
	cases := map[string]struct {
		args args
		want *v1alpha1.CertificateParameters
	}{
		"AllFilledNoDiff": {
			args: args{
				spec: &v1alpha1.CertificateParameters{
					DomainName:                               domainName,
					CertificateAuthorityARN:                  aws.String(certificateAuthorityArn),
					CertificateTransparencyLoggingPreference: &certificateTransparencyLoggingPreference,
				},
				in: &types.CertificateDetail{
					DomainName:              aws.String(domainName),
					CertificateAuthorityArn: aws.String(certificateAuthorityArn),
					Options:                 &types.CertificateOptions{CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled},
				},
			},
			want: &v1alpha1.CertificateParameters{
				DomainName:                               domainName,
				CertificateAuthorityARN:                  aws.String(certificateAuthorityArn),
				CertificateTransparencyLoggingPreference: &certificateTransparencyLoggingPreference,
			},
		},
		"AllFilledExternalDiff": {
			args: args{
				spec: &v1alpha1.CertificateParameters{
					DomainName:                               domainName,
					CertificateAuthorityARN:                  aws.String(certificateAuthorityArn),
					CertificateTransparencyLoggingPreference: &certificateTransparencyLoggingPreference,
				},
				in: &types.CertificateDetail{
					DomainName:              aws.String(domainName),
					CertificateAuthorityArn: aws.String(certificateAuthorityArn),
					Options:                 &types.CertificateOptions{CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled},
				},
			},
			want: &v1alpha1.CertificateParameters{
				DomainName:                               domainName,
				CertificateAuthorityARN:                  aws.String(certificateAuthorityArn),
				CertificateTransparencyLoggingPreference: &certificateTransparencyLoggingPreference,
			},
		},
		"PartialFilled": {
			args: args{
				spec: &v1alpha1.CertificateParameters{
					DomainName:              domainName,
					CertificateAuthorityARN: aws.String(certificateAuthorityArn),
				},
				in: &types.CertificateDetail{
					DomainName:              aws.String(domainName),
					CertificateAuthorityArn: aws.String(certificateAuthorityArn),
					Options:                 &types.CertificateOptions{CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled},
				},
			},
			want: &v1alpha1.CertificateParameters{
				DomainName:                               domainName,
				CertificateAuthorityARN:                  aws.String(certificateAuthorityArn),
				CertificateTransparencyLoggingPreference: &certificateTransparencyLoggingPreference,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			LateInitializeCertificate(tc.args.spec, tc.args.in)
			if diff := cmp.Diff(tc.args.spec, tc.want); diff != "" {
				t.Errorf("LateInitializeCertificate(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateCertificateStatus(t *testing.T) {
	cases := map[string]struct {
		in  types.CertificateDetail
		out v1alpha1.CertificateExternalStatus
	}{
		"AllFilled": {
			in: types.CertificateDetail{
				CertificateArn:     aws.String(certificateArn),
				RenewalEligibility: types.RenewalEligibilityEligible,
			},
			out: v1alpha1.CertificateExternalStatus{
				CertificateARN:     certificateArn,
				RenewalEligibility: types.RenewalEligibilityEligible,
			},
		},
		"NoRoleId": {
			in: types.CertificateDetail{
				CertificateArn:     nil,
				RenewalEligibility: types.RenewalEligibilityEligible,
			},
			out: v1alpha1.CertificateExternalStatus{
				CertificateARN:     "",
				RenewalEligibility: types.RenewalEligibilityEligible,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r := GenerateCertificateStatus(tc.in)
			if diff := cmp.Diff(r, tc.out); diff != "" {
				t.Errorf("GenerateCertificateStatus(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestIsCertificateUpToDate(t *testing.T) {
	certificateTransparencyLoggingPreference := types.CertificateTransparencyLoggingPreferenceDisabled
	type args struct {
		p    v1alpha1.CertificateParameters
		cd   types.CertificateDetail
		tags []types.Tag
	}

	cases := map[string]struct {
		args args
		want bool
	}{
		"SameFields": {
			args: args{
				cd: types.CertificateDetail{
					Options: &types.CertificateOptions{CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled},
				},
				p: v1alpha1.CertificateParameters{
					CertificateTransparencyLoggingPreference: &certificateTransparencyLoggingPreference,
					RenewCertificate:                         aws.Bool(false),
					Tags: []v1alpha1.Tag{{
						Key:   "key1",
						Value: "value1",
					}},
				},
				tags: []types.Tag{{
					Key:   aws.String("key1"),
					Value: aws.String("value1"),
				}},
			},
			want: true,
		},
		"DifferentFields": {
			args: args{
				cd: types.CertificateDetail{
					Options: &types.CertificateOptions{CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceEnabled},
				},
				p: v1alpha1.CertificateParameters{
					CertificateTransparencyLoggingPreference: &certificateTransparencyLoggingPreference,
					RenewCertificate:                         aws.Bool(false),
					Tags: []v1alpha1.Tag{{
						Key:   "key1",
						Value: "value1",
					}},
				},
				tags: []types.Tag{{
					Key:   aws.String("key1"),
					Value: aws.String("value1"),
				}},
			},
			want: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := IsCertificateUpToDate(tc.args.p, tc.args.cd, tc.args.tags)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}
