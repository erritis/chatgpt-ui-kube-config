synth:
  cd ./.cdk8s && cdk8s synth
werf-encrypt:
  werf helm secret values encrypt generator/secret-values.yaml -o .helm/secret-values.yaml
werf-decrypt:
  werf helm secret values decrypt .helm/secret-values.yaml -o generator/secret-values.yaml

werf-up *FLAGS:
  werf converge {{FLAGS}};
werf-down *FLAGS:
  werf dismiss {{FLAGS}};

werf-clear *FLAGS:
  werf dismiss {{FLAGS}};
  kubectl delete namespace chatgpt-ui;
