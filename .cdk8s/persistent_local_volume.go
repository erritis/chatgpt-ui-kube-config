package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

func NewLocalVolume(
	scope constructs.Construct,
	storageClassName string,
	id string,
	capacity cdk8s.Size,
	folder string,
	nodes *[]string,
) TuplePersistent {

	dbData := NewVolume(scope, storageClassName, id, capacity)

	dbData.persistent.ApiObject().AddJsonPatch(
		cdk8s.JsonPatch_Add(
			jsii.String("/spec/local"),
			&map[string]string{"path": folder},
		),
	)

	dbData.persistent.ApiObject().AddJsonPatch(
		cdk8s.JsonPatch_Add(
			jsii.String("/spec/nodeAffinity"),
			&map[string]interface{}{
				"required": &map[string]interface{}{
					"nodeSelectorTerms": &[]interface{}{
						&map[string]interface{}{
							"matchExpressions": &[]interface{}{
								&map[string]interface{}{
									"key":      jsii.String("kubernetes.io/hostname"),
									"operator": jsii.String("In"),
									"values":   nodes,
								},
							},
						},
					},
				},
			},
		),
	)

	return dbData
}
