# Feature Test for Issue 1248


This folder contains files describing how to address [Issue 1248](https://github.com/kubernetes-sigs/kustomize/issues/1248)

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
mkdir -p ${DEMO_HOME}/overlay1
mkdir -p ${DEMO_HOME}/overlay2
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
configMapGenerator:
- literals:
  - MY_ENV=foo
  name: example

vars:
- fieldref:
    fieldPath: data.MY_ENV
  name: MY_ENV
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: example

resources:
- deployment.yaml

EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
- ./overlay1
- ./overlay2
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay1/kustomization.yaml
resources:
- ../base

namespace: overlay1
nameSuffix: -overlay1

configMapGenerator:
- literals:
  - MY_ENV=bar
  name: example
  behavior: replace

EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay2/kustomization.yaml
resources:
- ../base

namespace: overlay2
nameSuffix: -overlay2

configMapGenerator:
- literals:
  - MY_ENV=baz
  name: example
  behavior: replace

EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dep
spec:
  template:
    spec:
      containers:
      - name: dep
        env:
        - name: SOME_ENV
          value: $(MY_ENV)

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
cat <<'EOF' >${DEMO_HOME}/expected/overlay1_apps_v1_deployment_dep-overlay1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dep-overlay1
  namespace: overlay1
spec:
  template:
    spec:
      containers:
      - env:
        - name: SOME_ENV
          value: bar
        name: dep
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/overlay1_~g_v1_configmap_example-overlay1-d6httb926d.yaml
apiVersion: v1
data:
  MY_ENV: bar
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: example-overlay1-d6httb926d
  namespace: overlay1
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/overlay2_apps_v1_deployment_dep-overlay2.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dep-overlay2
  namespace: overlay2
spec:
  template:
    spec:
      containers:
      - env:
        - name: SOME_ENV
          value: baz
        name: dep
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/overlay2_~g_v1_configmap_example-overlay2-fm7c465842.yaml
apiVersion: v1
data:
  MY_ENV: baz
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: example-overlay2-fm7c465842
  namespace: overlay2
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

