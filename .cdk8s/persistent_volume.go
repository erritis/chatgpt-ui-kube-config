package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
)

type TuplePersistent struct {
	persistent cdk8splus26.PersistentVolume
	volume     cdk8splus26.Volume
	claim      cdk8splus26.IPersistentVolumeClaim
}

func NewVolume(scope constructs.Construct, storageClassName string, id string, capacity cdk8s.Size) TuplePersistent {

	claim_id := fmt.Sprintf("%s-claim", id)

	claim_name := fmt.Sprintf("%s-%s", *scope.Node().Id(), claim_id)

	claim := cdk8splus26.PersistentVolumeClaim_FromClaimName(
		scope,
		jsii.String(claim_id),
		jsii.String(claim_name),
	)

	persistentProps := cdk8splus26.PersistentVolumeProps{
		VolumeMode: cdk8splus26.PersistentVolumeMode_FILE_SYSTEM,
		AccessModes: &[]cdk8splus26.PersistentVolumeAccessMode{
			cdk8splus26.PersistentVolumeAccessMode_READ_WRITE_ONCE,
			cdk8splus26.PersistentVolumeAccessMode_READ_ONLY_MANY,
		},
		ReclaimPolicy:    cdk8splus26.PersistentVolumeReclaimPolicy_RETAIN,
		Storage:          capacity,
		StorageClassName: jsii.String(storageClassName),
		Claim:            claim,
	}

	persistent := cdk8splus26.NewPersistentVolume(
		scope,
		jsii.String(id),
		&persistentProps,
	)

	volume := cdk8splus26.Volume_FromPersistentVolumeClaim(
		scope,
		jsii.String(fmt.Sprintf("%s-ref", claim_id)),
		claim,
		&cdk8splus26.PersistentVolumeClaimVolumeOptions{
			Name: jsii.String(claim_name),
		},
	)

	return TuplePersistent{
		persistent: persistent,
		volume:     volume,
		claim:      claim,
	}
}
