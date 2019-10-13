# Feature Test for Issue 0817


This folder contains files describing how to address [Issue 0817](https://github.com/kubernetes-sigs/kustomize/issues/0817)

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
- service.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: redis
spec:
  selector:
    app: redis
  ports:
  - name: server
    port: 6379
    protocol: TCP
    targetPort: redis
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
  app.kubernetes.io/part-of: argocd
fieldSpecs:
- path: metadata/labels
  create: true
# - path: spec/selector
#   create: false
#   version: v1
#   kind: Service
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
  create: true
  group: apps
  kind: StatefulSet
- path: spec/template/metadata/labels
  create: true
  group: apps
  kind: StatefulSet
- path: spec/template/spec/affinity/podAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
  create: false
  group: apps
  kind: StatefulSet
- path: spec/template/spec/affinity/podAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
  create: false
  group: apps
  kind: StatefulSet
- path: spec/template/spec/affinity/podAntiAffinity/preferredDuringSchedulingIgnoredDuringExecution/podAffinityTerm/labelSelector/matchLabels
  create: false
  group: apps
  kind: StatefulSet
- path: spec/template/spec/affinity/podAntiAffinity/requiredDuringSchedulingIgnoredDuringExecution/labelSelector/matchLabels
  create: false
  group: apps
  kind: StatefulSet
- path: spec/volumeClaimTemplates[]/metadata/labels
  create: true
  group: apps
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
kustomize build ${DEMO_HOME} -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_redis.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/part-of: argocd
  name: redis
spec:
  ports:
  - name: server
    port: 6379
    protocol: TCP
    targetPort: redis
  selector:
    app: redis
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

