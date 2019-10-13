# Feature Test for Issue 1399


This folder contains files describing how to address [Issue 1399](https://github.com/kubernetes-sigs/kustomize/issues/1399)

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
mkdir -p ${DEMO_HOME}/base/app
mkdir -p ${DEMO_HOME}/base/app/base
mkdir -p ${DEMO_HOME}/base/app/overlays
mkdir -p ${DEMO_HOME}/base/app/overlays/web
mkdir -p ${DEMO_HOME}/base/app/overlays/worker
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/prod
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/app/base/kustomization.yaml
resources:
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/app/overlays/web/kustomization.yaml
resources:
- ../../base
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/app/overlays/worker/kustomization.yaml
resources:
- ../../base
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- ./app/overlays/worker
- ./app/overlays/web
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/prod/kustomization.yaml
resources:
- ../../base

patchesStrategicMerge:
- patch-worker.yaml
- patch-web.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/app/overlays/web/deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy-web
  labels:
    mylabel: web
spec:
  replicas: 1
  selector:
    matchLabels:
      mylabel: web
  template:
    metadata:
      labels:
        mylabel: web
    spec:
      containers:
      - name: container1
        image: web-image:v1.0
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/app/overlays/worker/deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy-worker
  labels:
    mylabel: worker
spec:
  replicas: 1
  selector:
    matchLabels:
      mylabel: worker
  template:
    metadata:
      labels:
        mylabel: worker
    spec:
      containers:
      - name: container1
        image: worker-image:v1.0
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/prod/patch-web.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy-worker
spec:
  replicas: 5
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/prod/patch-worker.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy-web
spec:
  replicas: 3
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/prod -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_deploy-web.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    mylabel: web
  name: deploy-web
spec:
  replicas: 3
  selector:
    matchLabels:
      mylabel: web
  template:
    metadata:
      labels:
        mylabel: web
    spec:
      containers:
      - image: web-image:v1.0
        name: container1
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_deploy-worker.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    mylabel: worker
  name: deploy-worker
spec:
  replicas: 5
  selector:
    matchLabels:
      mylabel: worker
  template:
    metadata:
      labels:
        mylabel: worker
    spec:
      containers:
      - image: worker-image:v1.0
        name: container1
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

