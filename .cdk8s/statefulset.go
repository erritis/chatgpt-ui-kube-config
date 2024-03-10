package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
)

func NewStatefulSet(
	scope constructs.Construct,
	id string,
	image *string,
	port *float64,
	containerPort *float64,
	network string,
	variables *map[*string]*string,
	volumes *map[*string]*cdk8splus26.Volume,
) cdk8splus26.StatefulSet {

	container := cdk8splus26.NewContainer(&cdk8splus26.ContainerProps{
		Name:       jsii.String(id),
		Image:      image,
		PortNumber: containerPort,
		Resources: &cdk8splus26.ContainerResources{
			Cpu:              nil,
			EphemeralStorage: nil,
			Memory:           nil,
		},
		SecurityContext: &cdk8splus26.ContainerSecurityContextProps{
			ReadOnlyRootFilesystem: jsii.Bool(false),
			EnsureNonRoot:          jsii.Bool(false),
		},
	})

	for k, v := range *variables {
		container.Env().AddVariable(k, cdk8splus26.EnvValue_FromValue(v))
	}

	for path, volume := range *volumes {
		var storage cdk8splus26.IStorage = *volume
		container.Mount(path, storage, nil)
	}

	statefulset := cdk8splus26.NewStatefulSet(
		scope,
		jsii.String("statefulset"),
		&cdk8splus26.StatefulSetProps{
			Replicas: jsii.Number(1),
			Service: cdk8splus26.NewService(
				scope,
				jsii.String("service"),
				&cdk8splus26.ServiceProps{
					Type: cdk8splus26.ServiceType_CLUSTER_IP,
					Ports: &[]*cdk8splus26.ServicePort{
						{
							Port:       port,
							TargetPort: containerPort,
						},
					},
				},
			),
			SecurityContext: &cdk8splus26.PodSecurityContextProps{
				EnsureNonRoot: jsii.Bool(false),
			},
			PodMetadata: &cdk8s.ApiObjectMetadata{
				Labels: &map[string]*string{
					network: jsii.String("true"),
				},
			},
		},
	)

	statefulset.AttachContainer(container)

	statefulset.Metadata().AddLabel(jsii.String("io.service"), jsii.String(id))
	statefulset.Service().Metadata().AddLabel(jsii.String("io.service"), jsii.String(id))

	for _, volume := range *volumes {
		statefulset.AddVolume(*volume)
	}

	return statefulset
}
