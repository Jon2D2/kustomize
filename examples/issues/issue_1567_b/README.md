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

transformers:
- transformer.yaml

resources:
- statefulset.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
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


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/transformer.yaml
apiVersion: builtin
kind: LabelTransformer
metadata:
  name: labeltransformer
labels:
  app.kubernetes.io/part-of: elasticsearch
fieldSpecs:
- path: metadata/labels
  create: true
- path: spec/selector
  create: false
  version: v1
  kind: Service
- path: spec/selector
  create: true
  version: v1
  kind: ReplicationController
- path: spec/template/metadata/labels
  create: true
  version: v1
  kind: ReplicationController
- path: spec/selector/matchLabels
  create: true
  kind: Deployment
- path: spec/template/metadata/labels
  create: true
  kind: Deployment
- path: spec/template/spec/affinity/podAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
  create: false
  group: apps
  kind: Deployment
- path: spec/template/spec/affinity/podAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
  create: false
  group: apps
  kind: Deployment
- path: spec/template/spec/affinity/podAntiAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
  create: false
  group: apps
  kind: Deployment
- path: spec/template/spec/affinity/podAntiAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
  create: false
  group: apps
  kind: Deployment
- path: spec/selector/matchLabels
  create: true
  kind: ReplicaSet
- path: spec/template/metadata/labels
  create: true
  kind: ReplicaSet
- path: spec/selector/matchLabels
  create: true
  kind: DaemonSet
- path: spec/template/metadata/labels
  create: true
  kind: DaemonSet
- path: spec/selector/matchLabels
  create: false
  group: apps
  version: v1beta1
  kind: StatefulSet
- path: spec/template/metadata/labels
  create: true
  group: apps
  version: v1beta1
  kind: StatefulSet
- path: spec/template/spec/affinity/podAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
  create: false
  group: apps
  version: v1beta1
  kind: StatefulSet
- path: spec/template/spec/affinity/podAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
  create: false
  group: apps
  version: v1beta1
  kind: StatefulSet
- path: spec/template/spec/affinity/podAntiAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
  create: false
  group: apps
  version: v1beta1
  kind: StatefulSet
- path: spec/template/spec/affinity/podAntiAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
  create: false
  group: apps
  version: v1beta1
  kind: StatefulSet
- path: spec/volumeClaimTemplates[]/metadata/labels
  create: true
  group: apps
  version: v1beta1
  kind: StatefulSet
- path: spec/selector/matchLabels
  create: false
  group: batch
  kind: Job
- path: spec/template/metadata/labels
  create: true
  group: batch
  kind: Job
- path: spec/jobTemplate/spec/selector/matchLabels
  create: false
  group: batch
  kind: CronJob
- path: spec/jobTemplate/metadata/labels
  create: true
  group: batch
  kind: CronJob
- path: spec/jobTemplate/spec/template/metadata/labels
  create: true
  group: batch
  kind: CronJob
- path: spec/selector/matchLabels
  create: false
  group: policy
  kind: PodDisruptionBudget
- path: spec/podSelector/matchLabels
  create: false
  group: networking.k8s.io
  kind: NetworkPolicy
- path: spec/ingress/from/podSelector/matchLabels
  create: false
  group: networking.k8s.io
  kind: NetworkPolicy
- path: spec/egress/to/podSelector/matchLabels
  create: false
  group: networking.k8s.io
  kind: NetworkPolicy
- path: metadata/labels
  create: true
  kind: MyCRD
  skip: true
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

