# Feature Test for TestCase configmaps


This folder contains files for old test-case configmaps

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
mkdir -p ${DEMO_HOME}/base/myapp
mkdir -p ${DEMO_HOME}/base/myapp/mycomponent
mkdir -p ${DEMO_HOME}/base/myapp/mycomponent2
mkdir -p ${DEMO_HOME}/overlay
mkdir -p ${DEMO_HOME}/overlay/dev
mkdir -p ${DEMO_HOME}/overlay/dev/myapp
mkdir -p ${DEMO_HOME}/overlay/dev/myapp/mycomponent
mkdir -p ${DEMO_HOME}/overlay/dev/myapp/mycomponent2
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/myapp/mycomponent2/kustomization.yaml
namePrefix: p2-
configMapGenerator:
- name: com2
  behavior: create
  literals:
    - from=base
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/myapp/mycomponent/kustomization.yaml
namePrefix: p1-
configMapGenerator:
- name: com1
  behavior: create
  literals:
    - from=base
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/dev/kustomization.yaml
resources:
- myapp/mycomponent
- myapp/mycomponent2
configMapGenerator:
- name: com1
  behavior: merge
  literals:
    - foo=bar
    - baz=qux
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/dev/myapp/mycomponent2/kustomization.yaml
resources:
- ../../../../base/myapp/mycomponent2
configMapGenerator:
- name: com2
  behavior: merge
  literals:
    - from=overlay
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/dev/myapp/mycomponent/kustomization.yaml
resources:
- ../../../../base/myapp/mycomponent
configMapGenerator:
- name: com1
  behavior: merge
  literals:
    - from=overlay
EOF
```

## Execution

<!-- @build @test -->
```bash
kustomize build ${DEMO_HOME}/overlay/dev -o ${DEMO_HOME}/actual.yaml
```

## Verification


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected.yaml
apiVersion: v1
data:
  baz: qux
  foo: bar
  from: overlay
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: p1-com1-dhbbm922gd
---
apiVersion: v1
data:
  from: overlay
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: p2-com2-c4b8md75k9
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual.yaml $DEMO_HOME/expected.yaml | wc -l); \
echo $?
```

