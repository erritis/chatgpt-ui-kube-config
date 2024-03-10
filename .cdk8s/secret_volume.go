package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
)

func NewSecretVolume(scope constructs.Construct, id string, name *string, value *string) cdk8splus26.Volume {

	secret := cdk8splus26.NewSecret(
		scope,
		jsii.String(id),
		&cdk8splus26.SecretProps{Type: jsii.String("Opaque")},
	)
	secret.AddStringData(name, value)

	volume := cdk8splus26.Volume_FromSecret(
		scope,
		name,
		secret,
		&cdk8splus26.SecretVolumeOptions{
			Name: name,
			Items: &map[string]*cdk8splus26.PathMapping{
				"name": {Path: name},
			},
		},
	)
	return volume
}
