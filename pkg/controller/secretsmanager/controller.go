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

package secretsmanager

import (
	"context"
	"reflect"
	"sort"

	"github.com/google/go-cmp/cmp"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssecretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	awssecretsmanagertypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	runtimev1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-aws/apis/secretsmanager/v1alpha1"
	awsclient "github.com/crossplane/provider-aws/pkg/clients"
	"github.com/crossplane/provider-aws/pkg/clients/secretsmanager"
)

const (
	errNotSecret            = "managed resource is not a secret custom resource"
	errCreateFailed         = "failed to create secret"
	errDeleteFailed         = "failed to delete secret"
	secretMarkedForDeletion = "secret is marked for deletion"
	errUpdateFailed         = "failed to update secret"
	errDescribeSecretFailed = "failed to describe secret"
	errGetSecretValueFailed = "failed to get secret value"
	errK8sSecretNotFound    = "failed to get Kubernetes secret"
	errKubeUpdateFailed     = "failed to update Secret custom resource"
	errCreateTags           = "failed to create tags for the secret"
	errRemoveTags           = "failed to remove tags for the secret"
)

// SetupSecret adds a controller that reconciles a Secret.
func SetupSecret(mgr ctrl.Manager, l logging.Logger) error {
	name := managed.ControllerName(v1alpha1.SecretGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.Secret{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.SecretGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), newClientFn: secretsmanager.NewSecretsmanagerClient}),
			managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
			managed.WithInitializers(managed.NewDefaultProviderConfig(mgr.GetClient()), managed.NewNameAsExternalName(mgr.GetClient()), &tagger{kube: mgr.GetClient()}),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

type connector struct {
	kube        client.Client
	newClientFn func(aws.Config) secretsmanager.Client
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return nil, errors.New(errNotSecret)
	}
	cfg, err := awsclient.GetConfig(ctx, c.kube, mg, cr.Spec.ForProvider.Region)
	if err != nil {
		return nil, err
	}
	return &external{c.newClientFn(*cfg), c.kube}, nil
}

type external struct {
	client secretsmanager.Client
	kube   client.Client
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) { // nolint:gocyclo

	cr, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotSecret)
	}

	// Trigger creation of secret if external name is not set
	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	// Check the existence of the secret
	req, err := e.client.DescribeSecret(ctx, &awssecretsmanager.DescribeSecretInput{SecretId: awsclient.String(meta.GetExternalName(cr))})
	if err != nil {
		return managed.ExternalObservation{}, awsclient.Wrap(resource.Ignore(secretsmanager.IsErrorNotFound, err), errDescribeSecretFailed)
	}

	// Check the existence of the value of the secret.
	rsp, err := e.client.GetSecretValue(ctx, &awssecretsmanager.GetSecretValueInput{SecretId: awsclient.String(meta.GetExternalName(cr))})
	// Ignore empty secret values and deleted secrets in case of error
	if err != nil && !secretsmanager.IsErrorNotFound(err) && req.DeletedDate == nil {
		return managed.ExternalObservation{}, awsclient.Wrap(err, errGetSecretValueFailed)
	}

	// Update Crossplane secret if Kubernetes and AWS secret are different
	current := cr.Spec.ForProvider.DeepCopy()
	secretsmanager.LateInitialize(&cr.Spec.ForProvider, req)
	if !reflect.DeepEqual(current, &cr.Spec.ForProvider) {
		if err := e.kube.Update(ctx, cr); err != nil {
			return managed.ExternalObservation{}, errors.Wrap(err, errKubeUpdateFailed)
		}
	}

	secretsmanager.UpdateObservation(&cr.Status.AtProvider, rsp, req)
	cr.Status.SetConditions(runtimev1.Available())

	var resourceUpToDate bool
	if req.DeletedDate == nil {
		secret, err := secretsmanager.GetSecretValue(ctx, e.kube, cr)
		if err != nil {
			return managed.ExternalObservation{}, errors.New(errK8sSecretNotFound)
		}
		resourceUpToDate = secretsmanager.IsUpToDate(cr, req, secret, rsp)
	}
	if req.DeletedDate != nil {
		resourceUpToDate = true
		cr.Status.SetConditions(runtimev1.Deleting().WithMessage(secretMarkedForDeletion))
	}

	return managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        resourceUpToDate,
		ResourceLateInitialized: !cmp.Equal(current, &cr.Spec.ForProvider),
	}, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotSecret)
	}

	cr.SetConditions(runtimev1.Creating())

	secret, err := secretsmanager.GetSecretValue(ctx, e.kube, cr)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	_, err = e.client.CreateSecret(ctx, secretsmanager.GenerateCreateSecretsmanagerInput(meta.GetExternalName(cr), &cr.Spec.ForProvider, secret))
	return managed.ExternalCreation{}, awsclient.Wrap(err, errCreateFailed)
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {

	cr, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotSecret)
	}

	req, err := e.client.DescribeSecret(ctx, &awssecretsmanager.DescribeSecretInput{SecretId: awsclient.String(meta.GetExternalName(cr))})
	if err != nil {
		return managed.ExternalUpdate{}, awsclient.Wrap(err, errDescribeSecretFailed)
	}

	err = e.updateTags(ctx, cr, req)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	secret, err := secretsmanager.GetSecretValue(ctx, e.kube, cr)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	_, err = e.client.UpdateSecret(ctx, secretsmanager.GenerateUpdateSecretInput(meta.GetExternalName(cr), cr.Spec.ForProvider, secret))
	if err != nil {
		return managed.ExternalUpdate{}, awsclient.Wrap(err, errUpdateFailed)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {

	cr, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return errors.New(errNotSecret)
	}

	cr.SetConditions(runtimev1.Deleting())

	// DeletionDate is set based on the return value of the DeleteSecretRequest call below
	if cr.Status.AtProvider.DeletionDate != nil && cr.Spec.ForProvider.RecoveryWindowInDays != nil {
		// only request a new deletion if the user has changed the recovery window of secret object
		oldDeletionDate := cr.Status.AtProvider.DeletionDate
		newDeletionDate := &metav1.Time{Time: cr.Status.AtProvider.DeletedDate.AddDate(0, 0, int(*cr.Spec.ForProvider.RecoveryWindowInDays))}
		if oldDeletionDate.Equal(newDeletionDate) {
			return nil
		}
	}

	rsp, err := e.client.DeleteSecret(ctx, secretsmanager.GenerateDeleteSecretInput(meta.GetExternalName(cr), cr.Spec.ForProvider))
	if err != nil {
		return awsclient.Wrap(resource.Ignore(secretsmanager.IsErrorNotFound, err), errDeleteFailed)
	}

	// DeletionDate is returned by the DeleteSecretRequest call
	cr.Status.AtProvider.DeletionDate = &metav1.Time{Time: *rsp.DeletionDate}

	return nil
}

type tagger struct {
	kube client.Client
}

// TODO(knappek): split this out as it is used in several controllers
func (t *tagger) Initialize(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return errors.New(errNotSecret)
	}
	tagMap := map[string]string{}
	for _, tags := range cr.Spec.ForProvider.Tags {
		tagMap[tags.Key] = tags.Value
	}
	for k, v := range resource.GetExternalTags(mg) {
		tagMap[k] = v
	}
	cr.Spec.ForProvider.Tags = make([]v1alpha1.Tag, len(tagMap))
	i := 0
	for k, v := range tagMap {
		cr.Spec.ForProvider.Tags[i] = v1alpha1.Tag{Key: k, Value: v}
		i++
	}
	sort.Slice(cr.Spec.ForProvider.Tags, func(i, j int) bool {
		return cr.Spec.ForProvider.Tags[i].Key < cr.Spec.ForProvider.Tags[j].Key
	})
	return errors.Wrap(t.kube.Update(ctx, cr), errKubeUpdateFailed)
}

func (e *external) updateTags(ctx context.Context, secret *v1alpha1.Secret, req *awssecretsmanager.DescribeSecretOutput) error {
	add, remove := DiffTags(secret.Spec.ForProvider.Tags, req.Tags)
	if len(remove) != 0 {
		if _, err := e.client.UntagResource(ctx, &awssecretsmanager.UntagResourceInput{
			SecretId: awsclient.String(meta.GetExternalName(secret)),
			TagKeys:  remove,
		}); err != nil {
			return awsclient.Wrap(err, errRemoveTags)
		}
	}
	if len(add) != 0 {
		if _, err := e.client.TagResource(ctx, &awssecretsmanager.TagResourceInput{
			SecretId: awsclient.String(meta.GetExternalName(secret)),
			Tags:     add,
		}); err != nil {
			return awsclient.Wrap(err, errCreateTags)
		}

	}
	return nil
}

// GenerateSecretTags generates a tag array with type that secretsmanager client expects.
func GenerateSecretTags(tags []v1alpha1.Tag) []awssecretsmanagertypes.Tag {
	res := make([]awssecretsmanagertypes.Tag, len(tags))
	for i, t := range tags {
		res[i] = awssecretsmanagertypes.Tag{Key: aws.String(t.Key), Value: aws.String(t.Value)}
	}
	return res
}

// LocalTagsToMap converts []awssecretsmanagertypes.Tag to map
func LocalTagsToMap(tags []v1alpha1.Tag) map[string]string {
	result := make(map[string]string)
	for i := range tags {
		result[tags[i].Key] = tags[i].Value
	}
	return result
}

// RemoteTagsToMap converts []awssecretsmanagertypes.Tag to map
func RemoteTagsToMap(secretsmanagerTags []awssecretsmanagertypes.Tag) map[string]string {
	result := make(map[string]string)
	for i := range secretsmanagerTags {
		result[aws.ToString(secretsmanagerTags[i].Key)] = aws.ToString(secretsmanagerTags[i].Value)
	}
	return result
}

// DiffTags returns tags that should be added or removed.
func DiffTags(spec []v1alpha1.Tag, current []awssecretsmanagertypes.Tag) (addTags []awssecretsmanagertypes.Tag, remove []string) {
	local := LocalTagsToMap(spec)
	remote := RemoteTagsToMap(current)
	add := make(map[string]string, len(local))
	remove = []string{}
	for k, v := range local {
		add[k] = v
	}
	for k, v := range remote {
		switch val, ok := local[k]; {
		case ok && val != v:
			remove = append(remove, k)
		case !ok:
			remove = append(remove, k)
			delete(add, k)
		default:
			delete(add, k)
		}
	}
	addTags = []awssecretsmanagertypes.Tag{}
	for key, value := range add {
		value := value
		key := key
		addTags = append(addTags, awssecretsmanagertypes.Tag{Key: &key, Value: &value})
	}
	return
}
