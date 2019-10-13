# Feature Test for Issue 1727


This folder contains files describing how to address [Issue 1727](https://github.com/kubernetes-sigs/kustomize/issues/1727)

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
mkdir -p ${DEMO_HOME}/base
mkdir -p ${DEMO_HOME}/base/kustomizeconfig
mkdir -p ${DEMO_HOME}/overlay
mkdir -p ${DEMO_HOME}/overlay/staging
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- deployment.yaml
- vertical_pod_autoscaler.yaml

configurations:
- kustomizeconfig/name_references.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/staging/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

nameprefix: stagingpfx-
namespace: stagingns

bases:
- ../../base

patchesStrategicMerge:
- vertical_pod_autoscaler.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: puppetserver
  labels:
    app: puppetserver
spec:
  selector:
    matchLabels:
      app: puppetserver
  replicas: 1
  template:
    metadata:
      labels:
        app: puppetserver
    spec:
      containers:
      - name: main
        image: hack4easy/kubesim_health-amd64:latest
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
        ports:
        - name: pupperserver
          containerPort: 8081
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/name_references.yaml
nameReference:
- kind: Deployment
  fieldSpecs:
  - path: spec/targetRef/name
    group: autoscaling.k8s.io
    kind: VerticalPodAutoscaler

- kind: ReplicationController
  fieldSpecs:
  - path: spec/targetRef/name
    group: autoscaling.k8s.io
    kind: VerticalPodAutoscaler

- kind: ReplicaSet
  fieldSpecs:
  - path: spec/targetRef/name
    group: autoscaling.k8s.io
    kind: VerticalPodAutoscaler

EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/vertical_pod_autoscaler.yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: puppetserver
spec:
  targetRef:
    apiVersion: apps/v1
    kind:       Deployment
    name:       puppetserver
  updatePolicy:
    updateMode: Auto
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/staging/vertical_pod_autoscaler.yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: puppetserver
spec:
  updatePolicy:
    updateMode: Off
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlay/staging -o ${DEMO_HOME}/actual/staging.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/staging.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: puppetserver
  name: stagingpfx-puppetserver
  namespace: stagingns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: puppetserver
  template:
    metadata:
      labels:
        app: puppetserver
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:latest
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        name: main
        ports:
        - containerPort: 8081
          name: pupperserver
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
---
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: stagingpfx-puppetserver
  namespace: stagingns
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: stagingpfx-puppetserver
  updatePolicy:
    updateMode: false
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

