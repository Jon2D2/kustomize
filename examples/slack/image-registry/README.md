# Feature Test for Issue image-registry

This demonstrate how to use a mix of variable
and image transformer to change the registry
used by your deployments.

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/
mkdir -p ${DEMO_HOME}/kustomizeconfig
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
- deployment.yaml
- values.yaml

images:
- name: docker.io/nginx
  newName: from.the.kustomization.yaml/nginx
- name: docker.io/busybox
  newName: $(Values.my-values.spec.registry)/busybox

# The following sections are only necessary if you
# don't have the automatic variables declaration
# and varRef creation: https://github.com/kubernetes-sigs/kustomize/pull/1208.
# Uncomment out the vars and configurations section
# if you don't have the PR
#
# vars:
# - name: Values.my-values.spec.registry
#   objref:
#     apiVersion: kustomize.config.k8s.io/v1
#     kind: Values
#     name: my-values
#   fieldref:
#     fieldpath: spec.registry
# 
# configurations:
# - kustomizeconfig/varreference.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: dep1
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: dep1
  template:
    metadata:
      labels:
        app: dep1
    spec:
      serviceAccountName: dep1
      initContainers:
        - name: init
          image: $(Values.my-values.spec.registry)/withoutimagetransformer:latest
      containers:
        - name: nginx
          image: docker.io/nginx:latest
          env:
            - name: NFS_PATH
              value: /var/nfs
        - name: busybox
          image: docker.io/busybox:latest
          env:
            - name: PROVISIONER_NAME
              value: fuseim.pri/ifs

EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig/varreference.yaml
varReference:
- path: spec/template/spec/containers[]/image
  kind: Deployment
- path: spec/template/spec/initContainers[]/image
  kind: Deployment
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: my-values
spec:
  registry: from.the.values.io 
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build $DEMO_HOME -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_dep1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dep1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dep1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: dep1
    spec:
      containers:
      - env:
        - name: NFS_PATH
          value: /var/nfs
        image: from.the.kustomization.yaml/nginx:latest
        name: nginx
      - env:
        - name: PROVISIONER_NAME
          value: fuseim.pri/ifs
        image: from.the.values.io/busybox:latest
        name: busybox
      initContainers:
      - image: from.the.values.io/withoutimagetransformer:latest
        name: init
      serviceAccountName: dep1
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kustomize.config.k8s.io_v1_values_my-values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: my-values
spec:
  registry: from.the.values.io
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

