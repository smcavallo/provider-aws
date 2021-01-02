module github.com/crossplane/provider-aws

go 1.15

require (
	github.com/aws/aws-sdk-go v1.37.4
	github.com/aws/aws-sdk-go-v2 v0.31.0
	github.com/aws/aws-sdk-go-v2/config v0.4.0
	github.com/aws/aws-sdk-go-v2/credentials v0.2.0
	github.com/aws/aws-sdk-go-v2/service/acm v0.31.0
	github.com/aws/aws-sdk-go-v2/service/acmpca v0.31.0
	github.com/aws/aws-sdk-go-v2/service/cloudformation v0.31.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v0.31.0
	github.com/aws/aws-sdk-go-v2/service/ecr v0.31.0
	github.com/aws/aws-sdk-go-v2/service/eks v0.31.0
	github.com/aws/aws-sdk-go-v2/service/elasticache v0.31.0
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v0.31.0
	github.com/aws/aws-sdk-go-v2/service/iam v0.31.0
	github.com/aws/aws-sdk-go-v2/service/rds v0.31.0
	github.com/aws/aws-sdk-go-v2/service/redshift v0.31.0
	github.com/aws/aws-sdk-go-v2/service/route53 v0.31.0
	github.com/aws/aws-sdk-go-v2/service/s3 v0.31.0
	github.com/aws/aws-sdk-go-v2/service/sns v0.31.0
	github.com/aws/aws-sdk-go-v2/service/sqs v0.31.0
	github.com/aws/aws-sdk-go-v2/service/sts v0.31.0
	github.com/aws/smithy-go v0.5.0
	github.com/crossplane/crossplane-runtime v0.13.0
	github.com/crossplane/crossplane-tools v0.0.0-20201201125637-9ddc70edfd0d
	github.com/evanphx/json-patch v4.9.0+incompatible
	github.com/go-ini/ini v1.46.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/go-cmp v0.5.4
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/copystructure v1.0.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/gomega v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/tools v0.0.0-20200916195026-c9a70fc28ce3 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	k8s.io/api v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
	sigs.k8s.io/controller-runtime v0.8.0
	sigs.k8s.io/controller-tools v0.4.0
)
