# Feature Test for Issue 1459


This folder contains files describing how to address [Issue 1459](https://github.com/kubernetes-sigs/kustomize/issues/1459)

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
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  common: label

configurations:
- kustomizeconfig/deployment.yaml

resources:
- deployment.yaml

EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: masterdata
    spec:
      affinity:
        podAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchLabels:
                  app: memcached
              topologyKey: kubernetes.io/hostname
      containers:
      - name: deployment
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig/deployment.yaml
commonLabels:
- path: spec/template/spec/affinity/podAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
  group: apps
  version: v1
  kind: Deployment
  skip: true 
  # behavior: drop
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    common: label
  name: deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      common: label
  template:
    metadata:
      labels:
        app: masterdata
        common: label
    spec:
      affinity:
        podAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchLabels:
                app: memcached
            topologyKey: kubernetes.io/hostname
      containers:
      - name: deployment
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

