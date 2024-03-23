_synth-production:
  cd ./.cdk8s && go run ./src -e=Production

_encrypt-production:
  werf helm secret values encrypt .cdk8s/secret-values.production.yaml -o .helm-production/secret-values.yaml
_decrypt-production:
  werf helm secret values decrypt .helm-production/secret-values.yaml -o .cdk8s/secret-values.production.yaml


_up-production *FLAGS:
  werf converge --config='werf.production.yaml' {{FLAGS}};
_down-production *FLAGS:
  werf dismiss --config='werf.production.yaml' {{FLAGS}};
