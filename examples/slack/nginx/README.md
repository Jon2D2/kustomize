# Feature Test for UseCase nginx


This folder contains files matching a slack thread.
Original files have been pulled from [here](https://github.com/neith00/kustomize-demo)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/base
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/dev
mkdir -p ${DEMO_HOME}/overlays/prod
mkdir -p ${DEMO_HOME}/overlays/prod/base
mkdir -p ${DEMO_HOME}/overlays/prod/cluster1
mkdir -p ${DEMO_HOME}/overlays/prod/cluster2
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- nginx-deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/dev/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namePrefix: dev-

commonLabels:
  variant: dev

commonAnnotations:
  note: manifests for dev environement

patchesStrategicMerge:
- replicas_count.yaml

resources:
- ../../base
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/prod/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namePrefix: prod-

commonLabels:
  variant: prod

commonAnnotations:
  note: manifests for prod environement

resources:
- ../../../base

patchesStrategicMerge:
- replicas_count.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/prod/cluster1/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../base

patchesStrategicMerge:
- cpu_limit.yaml
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/prod/cluster2/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../base

patchesStrategicMerge:
- cpu_limit.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/nginx-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx:1.16.1
        name: nginx
        resources:
          limits:
            cpu: 100m
            memory: 128Mi
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/dev/replicas_count.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 3
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/prod/base/replicas_count.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 6
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/prod/cluster1/cpu_limit.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
       containers:
       - name: nginx
         resources:
           limits:
             cpu: 1234m
             memory: 128Mi
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/prod/cluster2/cpu_limit.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
       containers:
       - name: nginx
         resources:
           limits:
             cpu: 432m
             memory: 128Mi
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/dev -o ${DEMO_HOME}/actual/dev.yaml
kustomize build ${DEMO_HOME}/overlays/prod/cluster1 -o ${DEMO_HOME}/actual/cluster1.yaml
kustomize build ${DEMO_HOME}/overlays/prod/cluster2 -o ${DEMO_HOME}/actual/cluster2.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cluster1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: manifests for prod environement
  labels:
    app: nginx
    variant: prod
  name: prod-nginx
spec:
  replicas: 6
  selector:
    matchLabels:
      app: nginx
      variant: prod
  template:
    metadata:
      annotations:
        note: manifests for prod environement
      labels:
        app: nginx
        variant: prod
    spec:
      containers:
      - image: nginx:1.16.1
        name: nginx
        resources:
          limits:
            cpu: 1234m
            memory: 128Mi
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cluster2.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: manifests for prod environement
  labels:
    app: nginx
    variant: prod
  name: prod-nginx
spec:
  replicas: 6
  selector:
    matchLabels:
      app: nginx
      variant: prod
  template:
    metadata:
      annotations:
        note: manifests for prod environement
      labels:
        app: nginx
        variant: prod
    spec:
      containers:
      - image: nginx:1.16.1
        name: nginx
        resources:
          limits:
            cpu: 432m
            memory: 128Mi
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: manifests for dev environement
  labels:
    app: nginx
    variant: dev
  name: dev-nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
      variant: dev
  template:
    metadata:
      annotations:
        note: manifests for dev environement
      labels:
        app: nginx
        variant: dev
    spec:
      containers:
      - image: nginx:1.16.1
        name: nginx
        resources:
          limits:
            cpu: 100m
            memory: 128Mi
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

