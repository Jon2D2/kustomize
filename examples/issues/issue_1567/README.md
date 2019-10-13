# Feature Test for Issue 1567


This folder contains files describing how to address [Issue 1567](https://github.com/kubernetes-sigs/kustomize/issues/1567)

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
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

configurations:
- commonlabels.yaml

commonLabels:
  app.kubernetes.io/part-of: elasticsearch

resources:
- statefulset.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/commonlabels.yaml
commonLabels:
- path: spec/selector/matchLabels
  kind: StatefulSet
  behavior: remove
- path: spec/selector/matchLabels
  create: false
  group: apps
  version: v1beta1
  kind: StatefulSet
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/statefulset.yaml
---
# Do not change this apiVersion to let it deploy on
# Kubernetes 1.16. It is used to showcase handling of
# multiple apiVersion for a same kind
apiVersion: apps/v1beta1
kind: StatefulSet
metadata:
  name: elasticsearch-data
spec:
  serviceName: elasticsearch-data
  replicas: 3
  updateStrategy:
    type: RollingUpdate
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1beta1_statefulset_elasticsearch-data.yaml
apiVersion: apps/v1beta1
kind: StatefulSet
metadata:
  labels:
    app.kubernetes.io/part-of: elasticsearch
  name: elasticsearch-data
spec:
  replicas: 3
  serviceName: elasticsearch-data
  template:
    metadata:
      labels:
        app.kubernetes.io/part-of: elasticsearch
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

