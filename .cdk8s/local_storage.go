package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2/k8s"
)

func NewLocalStorage(scope constructs.Construct, id string) k8s.KubeStorageClass {

	annotations := make(map[string]*string)

	annotations["storageclass.kubernetes.io/is-default-class"] = jsii.String("true")

	metadata := k8s.ObjectMeta{
		Name:        jsii.String(id),
		Annotations: &annotations,
	}

	storageProps := k8s.KubeStorageClassProps{
		Provisioner:       jsii.String("kubernetes.io/no-provisioner"),
		Metadata:          &metadata,
		VolumeBindingMode: jsii.String("WaitForFirstConsumer"),
	}

	storage := k8s.NewKubeStorageClass(scope, jsii.String(id), &storageProps)

	storage.AddJsonPatch(cdk8s.JsonPatch_Replace(jsii.String("/metadata/namespace"), new(*string)))

	return storage
}
