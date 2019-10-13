# Feature Test for Issue 1251


This folder contains files describing how to address [Issue 1251](https://github.com/kubernetes-sigs/kustomize/issues/1251)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/diamond
mkdir -p ${DEMO_HOME}/diamond/base
mkdir -p ${DEMO_HOME}/diamond/overlays
mkdir -p ${DEMO_HOME}/diamond/overlays/aggregate
mkdir -p ${DEMO_HOME}/diamond/overlays/overlay1
mkdir -p ${DEMO_HOME}/diamond/overlays/overlay2
mkdir -p ${DEMO_HOME}/diamond-with-patches
mkdir -p ${DEMO_HOME}/diamond-with-patches/base
mkdir -p ${DEMO_HOME}/diamond-with-patches/overlays
mkdir -p ${DEMO_HOME}/diamond-with-patches/overlays/aggregate
mkdir -p ${DEMO_HOME}/diamond-with-patches/overlays/overlay1
mkdir -p ${DEMO_HOME}/diamond-with-patches/overlays/overlay1/patches
mkdir -p ${DEMO_HOME}/diamond-with-patches/overlays/overlay2
mkdir -p ${DEMO_HOME}/diamond-with-patches/overlays/overlay2/patches
mkdir -p ${DEMO_HOME}/mixin
mkdir -p ${DEMO_HOME}/mixin/base
mkdir -p ${DEMO_HOME}/mixin/overlays
mkdir -p ${DEMO_HOME}/mixin/overlays/aggregate
mkdir -p ${DEMO_HOME}/mixin/overlays/overlay1
mkdir -p ${DEMO_HOME}/mixin/overlays/overlay1/patches
mkdir -p ${DEMO_HOME}/mixin/overlays/overlay2
mkdir -p ${DEMO_HOME}/mixin/overlays/overlay2/patches
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond/base/kustomization.yaml
resources:
  - deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond/overlays/aggregate/kustomization.yaml
resources:
  - ../overlay1
  - ../overlay2
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond/overlays/overlay1/kustomization.yaml
# nameprefix: p1-
resources:
  - ../../base
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond/overlays/overlay2/kustomization.yaml
# nameprefix: p2-
resources:
  - ../../base
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond-with-patches/base/kustomization.yaml
resources:
  - deployment.yaml
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond-with-patches/overlays/aggregate/kustomization.yaml
resources:
  - ../overlay1
  - ../overlay2
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond-with-patches/overlays/overlay1/kustomization.yaml
resources:
  - ../../base

patchesStrategicMerge:
  - patches/patch.yaml
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond-with-patches/overlays/overlay2/kustomization.yaml
resources:
  - ../../base

patchesStrategicMerge:
  - patches/patch.yaml
EOF
```


### Preparation Step KustomizationFile8

<!-- @createKustomizationFile8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/mixin/base/kustomization.yaml
resources:
  - deployment.yaml
EOF
```


### Preparation Step KustomizationFile9

<!-- @createKustomizationFile9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/mixin/overlays/aggregate/kustomization.yaml
resources:
  - ../../base
  - ../overlay1
  - ../overlay2
EOF
```


### Preparation Step KustomizationFile10

<!-- @createKustomizationFile10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/mixin/overlays/overlay1/kustomization.yaml
resources:
  - ../../base

patchesStrategicMerge:
  - patches/patch.yaml
EOF
```


### Preparation Step KustomizationFile11

<!-- @createKustomizationFile11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/mixin/overlays/overlay2/kustomization.yaml
resources:
  - ../../base

patchesStrategicMerge:
  - patches/patch.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
        image: my-image
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond-with-patches/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
        image: my-image
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond-with-patches/overlays/overlay1/patches/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  replicas: 2
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/diamond-with-patches/overlays/overlay2/patches/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    metadata:
      labels:
        app: my-app
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/mixin/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
        image: my-image
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/mixin/overlays/overlay1/patches/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  replicas: 2
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/mixin/overlays/overlay2/patches/patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    metadata:
      labels:
        app: my-app
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir -p ${DEMO_HOME}/actual
mkdir -p ${DEMO_HOME}/actual/diamond
kustomize build ${DEMO_HOME}/diamond/overlays/aggregate -o ${DEMO_HOME}/actual/diamond
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir -p ${DEMO_HOME}/expected
mkdir -p ${DEMO_HOME}/expected/diamond
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/diamond/apps_v1_deployment_my-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - image: my-image
        name: my-deployment
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

