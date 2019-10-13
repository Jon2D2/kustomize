# Feature Test for Issue 0506


This folder contains files describing how to address [Issue 0506](https://github.com/kubernetes-sigs/kustomize/issues/0506)

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
mkdir -p ${DEMO_HOME}/app1
mkdir -p ${DEMO_HOME}/app2
mkdir -p ${DEMO_HOME}/component1
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app1/kustomization.yaml
resources:
- ../component1

namePrefix: app1-
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app2/kustomization.yaml
resources:
- ../component1

namePrefix: app2-
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/component1/kustomization.yaml
resources:
  - resources.yaml

vars:
- name: POD_NAME
  objref:
    apiVersion: v1
    kind: Pod
    name: component1
  fieldref:
    fieldpath: metadata.name
- name: IMAGE_NAME
  objref:
    apiVersion: v1
    kind: Pod
    name: component1
  fieldref:
    fieldpath: spec.containers[0].image
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
- app1
- app2
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/component1/resources.yaml
apiVersion: v1
kind: Pod
metadata:
  name: component1
spec:
  containers:
    - name: component1
      image: bash
      env:
        - name: POD_NAME
          value: $(POD_NAME)
        - name: IMAGE_NAME
          value: $(IMAGE_NAME)
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/app1 -o ${DEMO_HOME}/actual/app1.yaml
kustomize build ${DEMO_HOME}/app2 -o ${DEMO_HOME}/actual/app2.yaml
kustomize build ${DEMO_HOME} -o ${DEMO_HOME}/actual/both.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/app1.yaml
apiVersion: v1
kind: Pod
metadata:
  name: app1-component1
spec:
  containers:
  - env:
    - name: POD_NAME
      value: app1-component1
    - name: IMAGE_NAME
      value: bash
    image: bash
    name: component1
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/app2.yaml
apiVersion: v1
kind: Pod
metadata:
  name: app2-component1
spec:
  containers:
  - env:
    - name: POD_NAME
      value: app2-component1
    - name: IMAGE_NAME
      value: bash
    image: bash
    name: component1
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/both.yaml
apiVersion: v1
kind: Pod
metadata:
  name: app1-component1
spec:
  containers:
  - env:
    - name: POD_NAME
      value: app1-component1
    - name: IMAGE_NAME
      value: bash
    image: bash
    name: component1
---
apiVersion: v1
kind: Pod
metadata:
  name: app2-component1
spec:
  containers:
  - env:
    - name: POD_NAME
      value: app2-component1
    - name: IMAGE_NAME
      value: bash
    image: bash
    name: component1
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

