# Feature Test for Issue 1404


This folder contains files describing how to address [Issue 1404](https://github.com/kubernetes-sigs/kustomize/issues/1404)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/test-kus
mkdir -p ${DEMO_HOME}/test-kus/base
mkdir -p ${DEMO_HOME}/test-kus/overlays
mkdir -p ${DEMO_HOME}/test-kus/overlays/production
mkdir -p ${DEMO_HOME}/test-kus/overlays/production/blue
mkdir -p ${DEMO_HOME}/test-kus/overlays/production/green
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: test
resources:
  - test-kus/overlays
images:
  - name: leaf/esp-service
    newName: 11111.dkr.ecr.us-west-2.amazonaws.com/leaf/esp-service
    newTag: latest
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/test-kus/base/kustomization.yaml
resources:
    - deployment.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/test-kus/overlays/kustomization.yaml
resources:
  - production/green
  - production/blue
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/test-kus/overlays/production/blue/kustomization.yaml
commonAnnotations:
    note: This is the production blue environment
nameSuffix: -blue
resources:
  - ../../../base
patchesStrategicMerge:
  - deployment.yaml
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/test-kus/overlays/production/green/kustomization.yaml
commonAnnotations:
    note: This is the production green environment
nameSuffix: -green
resources:
  - ../../../base
patchesStrategicMerge:
  - deployment.yaml
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/test-kus/overlays/production/kustomization.yaml
resources:
  - green
  - blue
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/test-kus/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: esp
spec:
  template:
    spec:
      containers:
      - name: esp
        imagePullPolicy: Always
        image: leaf/esp-service
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/test-kus/overlays/production/blue/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: esp
spec:
  template:
    spec:
      containers:
      - name: esp
        image: leaf/esp-service:latest
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/test-kus/overlays/production/green/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: esp
spec:
  template:
    spec:
      containers:
      - name: esp
        image: leaf/esp-service:latest
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_esp-blue.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: This is the production blue environment
  name: esp-blue
  namespace: test
spec:
  template:
    metadata:
      annotations:
        note: This is the production blue environment
    spec:
      containers:
      - image: 11111.dkr.ecr.us-west-2.amazonaws.com/leaf/esp-service:latest
        imagePullPolicy: Always
        name: esp
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_esp-green.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: This is the production green environment
  name: esp-green
  namespace: test
spec:
  template:
    metadata:
      annotations:
        note: This is the production green environment
    spec:
      containers:
      - image: 11111.dkr.ecr.us-west-2.amazonaws.com/leaf/esp-service:latest
        imagePullPolicy: Always
        name: esp
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

