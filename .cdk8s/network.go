package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
)

func NewNetworkPolicy(
	scope constructs.Construct,
	id string,
	network string,
) cdk8splus26.NetworkPolicy {

	selector := cdk8splus26.Pods_Select(
		scope,
		jsii.String("selector"),
		&cdk8splus26.PodsSelectOptions{
			Labels: &map[string]*string{
				network: jsii.String("true"),
			},
		},
	)

	networkPolicy := cdk8splus26.NewNetworkPolicy(
		scope,
		jsii.String(id),
		&cdk8splus26.NetworkPolicyProps{
			Ingress: &cdk8splus26.NetworkPolicyTraffic{
				Rules: &[]*cdk8splus26.NetworkPolicyRule{
					{
						Peer: selector,
					},
				},
			},
			Selector: selector,
		},
	)

	return networkPolicy
}
