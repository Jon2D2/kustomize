# Feature Test for Issue 1480


This folder contains files describing how to address [Issue 1480](https://github.com/kubernetes-sigs/kustomize/issues/1480)

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
mkdir -p ${DEMO_HOME}/1-shared-base
mkdir -p ${DEMO_HOME}/2-app-overlay
mkdir -p ${DEMO_HOME}/3-dev-overlay
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/1-shared-base/kustomization.yaml
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/2-app-overlay/kustomization.yaml
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: myapp-

resources:
- ../1-shared-base
- serviceaccount.yaml

patchesStrategicMerge:
- add-serviceaccount-to-deployment.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/3-dev-overlay/kustomization.yaml
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

namePrefix: dev-

resources:
  - ../2-app-overlay
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/1-shared-base/deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 2
  strategy:
    rollingUpdate:
      maxSurge: 15%
      maxUnavailable: 10%
    type: RollingUpdate
  template:
    metadata:
      labels:
        logformat: json
        tier: frontend
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              topologyKey: kubernetes.io/hostname
            weight: 100
      containers:
      - args: []
        command: []
        image: container-image
        livenessProbe:
          httpGet:
            path: /healthz
            port: http
          initialDelaySeconds: 90
        name: main
        ports:
        - containerPort: 8000
          name: http
        readinessProbe:
          httpGet:
            path: /readiness
            port: http
          initialDelaySeconds: 30
        resources:
          limits:
            cpu: 200m
            memory: 256Mi
          requests:
            cpu: 100m
            memory: 128Mi
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/2-app-overlay/add-serviceaccount-to-deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      serviceAccountName: app
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/2-app-overlay/serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/3-dev-overlay -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_dev-myapp-app.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dev-myapp-app
spec:
  replicas: 2
  strategy:
    rollingUpdate:
      maxSurge: 15%
      maxUnavailable: 10%
    type: RollingUpdate
  template:
    metadata:
      labels:
        logformat: json
        tier: frontend
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              topologyKey: kubernetes.io/hostname
            weight: 100
      containers:
      - args: []
        command: []
        image: container-image
        livenessProbe:
          httpGet:
            path: /healthz
            port: http
          initialDelaySeconds: 90
        name: main
        ports:
        - containerPort: 8000
          name: http
        readinessProbe:
          httpGet:
            path: /readiness
            port: http
          initialDelaySeconds: 30
        resources:
          limits:
            cpu: 200m
            memory: 256Mi
          requests:
            cpu: 100m
            memory: 128Mi
      serviceAccountName: dev-myapp-app
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_serviceaccount_dev-myapp-app.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: dev-myapp-app
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

