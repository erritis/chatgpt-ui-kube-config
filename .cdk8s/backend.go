package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
)

type TupleBackend struct {
	deployment cdk8splus26.Deployment
	service    cdk8splus26.Service
}

func NewBackend(
	scope constructs.Construct,
	id string,
	image *string,
	port *float64,
	containerPort *float64,
	network string,
	variables *map[*string]*string,
) TupleBackend {

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

	deployment := cdk8splus26.NewDeployment(
		scope,
		jsii.String("deployment"),
		&cdk8splus26.DeploymentProps{
			Replicas: jsii.Number(1),
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

	deployment.AttachContainer(container)

	deployment.Metadata().AddLabel(jsii.String("io.service"), jsii.String(id))

	service := deployment.ExposeViaService(&cdk8splus26.DeploymentExposeViaServiceOptions{
		Name:        jsii.String(fmt.Sprintf("%s-service", id)),
		ServiceType: cdk8splus26.ServiceType_CLUSTER_IP,
		Ports: &[]*cdk8splus26.ServicePort{
			{
				Port:       port,
				TargetPort: containerPort,
			},
		},
	})

	return TupleBackend{
		deployment: deployment,
		service:    service,
	}
}
