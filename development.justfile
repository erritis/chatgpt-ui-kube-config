_synth-development:
  cd ./.cdk8s && go run ./src -e=Development

_encrypt-development:
  werf helm secret values encrypt .cdk8s/secret-values.development.yaml -o .helm-development/secret-values.yaml
_decrypt-development:
  werf helm secret values decrypt .helm-development/secret-values.yaml -o .cdk8s/secret-values.development.yaml


_up-development *FLAGS:
  werf converge --config='werf.development.yaml' {{FLAGS}};
_down-development *FLAGS:
  werf dismiss --config='werf.development.yaml' {{FLAGS}};
