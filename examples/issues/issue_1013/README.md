# Feature Test for Issue 1013


This folder contains files describing how to address [Issue 1013](https://github.com/kubernetes-sigs/kustomize/issues/1013)

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
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

commonLabels:
  app.kubernetes.io/app: elasticsearch
  app: elasticsearch
  role: master

resources:
- resource.yaml

configurations:
- ./kustomizeconfig/labels.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig/labels.yaml
commonLabels:
- path: spec/template/spec/affinity/podAntiAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
  create: true
  group: apps
  kind: StatefulSet
  behavior: replace
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/resource.yaml
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: elasticsearch-data
spec:
  serviceName: elasticsearch-data
  replicas: 3
  updateStrategy:
    type: RollingUpdate
  selector:
    matchLabels:
      tier: logging-plane
  template:
    metadata:
      labels:
        tier: logging-plane
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                # matchLabels:
                matchExpressions:
                - key: app.kubernetes.io/app
                  operator: In
                  values:
                  - elasticsearch
                - key: role
                  operator: In
                  values:
                  - data
              topologyKey: kubernetes.io/hostname
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_statefulset_elasticsearch-data.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    app: elasticsearch
    app.kubernetes.io/app: elasticsearch
    role: master
  name: elasticsearch-data
spec:
  replicas: 3
  selector:
    matchLabels:
      app: elasticsearch
      app.kubernetes.io/app: elasticsearch
      role: master
      tier: logging-plane
  serviceName: elasticsearch-data
  template:
    metadata:
      labels:
        app: elasticsearch
        app.kubernetes.io/app: elasticsearch
        role: master
        tier: logging-plane
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app.kubernetes.io/app
                  operator: In
                  values:
                  - elasticsearch
                - key: role
                  operator: In
                  values:
                  - data
                matchLabels:
                  app: elasticsearch
                  app.kubernetes.io/app: elasticsearch
                  role: master
              topologyKey: kubernetes.io/hostname
            weight: 100
  updateStrategy:
    type: RollingUpdate
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

