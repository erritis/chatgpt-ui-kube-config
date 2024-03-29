# README

ChatGPT UI configurations for deploying to Kubernetes.

To organize a set of commands into a project, use [Just](https://github.com/casey/just).

---

### Requirements:

To work with this project you need:

- Kubernetes (you can [Minikube](https://minikube.sigs.k8s.io/docs/start/))

- [helm](https://helm.sh/docs/intro/install/)

- [kompose](https://kompose.io/installation/)

- [Werf](https://werf.io/installation.html)

- Ingress Nginx Controller (or set up ingress in minikube like in [this article](https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/))

---

### Preset storage:

> **Warning**
>
> This is not required for minikube and kubernetes with default storage providers configured. Skip this step if this is your case.

Then you need to create a folder on the server:

> sudo mkdir /mnt/chatgptdb

The default folder will not be cleared. To remove it, go to the server and run:

> sudo rm -Rf /mnt/chatgptdb

---

### Deploy configurations:

Before deploying the project, set up configurations and secrets. In the [.origin](.origin) folder there is an example of configuration files with the **.example** extension, in order to use them copy the file, remove the **.example** extension and make changes if necessary.

To create configurations for Kubernetes, run:

> just werf-convert

After that, in the folder [.helm/templates](.helm/templates) there will be configuration files and secrets that do not fall into the git history.

---

### Application deployment:

To deploy, run:

> just werf-up

To uninstall an application, run:

> just werf-down

To completely remove, including configuration, secrets, and local-storage, run:

> just werf-clear

---

### Encrypt configuration files:

To encrypt the contents of the [.origin](.origin) folder, run:

> just werf-encrypt

To decrypt encrypted files back, run:

> just werf-decrypt

Note that you must have a **.werf_secret_key** encryption key file in your project folder (there are also [two other options](https://werf.io/documentation/v1.1/reference/deploy_process/working_with_secrets.html)).

To generate it run:

> werf helm secret generate-secret-key > .werf_secret_key
