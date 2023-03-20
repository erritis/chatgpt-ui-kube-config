old_repo := "registry.argiago.ru"

werf-set-repo repo:
  old_repo={{old_repo}} \
  && repo={{repo}} \
  && sed -i "s/$(echo $old_repo | sed -e 's/[\/&]/\\&/g')/$(echo $repo | sed -e 's/[\/&]/\\&/g')/g" .justfile;

werf-convert:
  kompose convert -f docker-compose.yml -o ./.helm/templates;
  rm ./.helm/templates/*-persistentvolumeclaim.yaml;
  find ./.helm/templates -type f -exec sed -i "s/'{{{{ \(.*\) }}'/{{{{ \1 }}/g" {} +;
  find ./.helm/templates -type f -exec sed -i "s/\.values/\.Values/g" {} +;

werf-encrypt:
  werf helm secret values encrypt .origin/secret-values.yaml -o .helm/secret-values.yaml
  bash -c 'for filename in .origin/config/*; do name=${filename##*/}; werf helm secret file encrypt ".origin/config/$name" -o ".helm/config/$name"; done;';
  bash -c 'for filename in .origin/secret/*; do name=${filename##*/}; werf helm secret file encrypt ".origin/secret/$name" -o ".helm/secret/$name"; done;';
  rm ./.helm/*/*.example;
werf-decrypt:
  werf helm secret values decrypt .helm/secret-values.yaml -o .origin/secret-values.yaml
  bash -c 'for filename in .helm/config/*; do name=${filename##*/}; werf helm secret file encrypt ".helm/config/$name" -o ".origin/config/$name"; done;';
  bash -c 'for filename in .helm/secret/*; do name=${filename##*/}; werf helm secret file encrypt ".helm/secret/$name" -o ".origin/secret/$name"; done;';

werf-up-storage:
  kubectl apply -f local-storage.yaml;
  kubectl apply -f chatgptdb-pv-0.yaml;
werf-down-storage:
  kubectl delete -f chatgptdb-pv-0.yaml;
  kubectl delete -f local-storage.yaml;

werf-up-conf:
  kubectl create namespace chatgpt-ui &>/dev/null || exit 0;
  kubectl config set-context --current --namespace=chatgpt-ui;
  kubectl apply -Rf './.helm/templates/*-configmap.yaml';
  kubectl apply -Rf './.helm/templates/*-secret.yaml';
werf-down-conf:
  kubectl delete -Rf '/.helm/templates/*-configmap.yaml';
  kubectl delete -Rf './.helm/templates/*-secret.yaml';

werf-up *FLAGS:
  werf converge --repo {{old_repo}}/chatgpt-discord-bot {{FLAGS}};
werf-down *FLAGS:
  werf dismiss --repo {{old_repo}}/chatgpt-discord-bot {{FLAGS}};
  
werf-clear *FLAGS:
  werf dismiss --repo {{old_repo}}/chatgpt-discord-bot {{FLAGS}};
  kubectl delete namespace chatgpt-ui;
  kubectl delete -f chatgptdb-pv-0.yaml;
  kubectl delete -f local-storage.yaml;
  