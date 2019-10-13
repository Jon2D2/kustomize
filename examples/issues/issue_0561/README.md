# Feature Test for Issue 0561


This folder contains files describing how to address [Issue 0561](https://github.com/kubernetes-sigs/kustomize/issues/0561)

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
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
  - resources.yaml

namespace: test
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
  - base

configMapGenerator:
  - name: config
    namespace: test
    files:
      - kustomization.yaml

patchesStrategicMerge:
  - patch.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/resources.yaml
apiVersion: v1
kind: Pod
metadata:
  name: test
spec:
  containers:
    - name: test
      image: bash
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/patch.yaml
apiVersion: v1
kind: Pod
metadata:
  name: test
  namespace: test
spec:
  volumes:
    - name: config
      configMap:
        name: config
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_config-9mk7dchmcb.yaml
apiVersion: v1
data:
  kustomization.yaml: |
    resources:
      - base
    configMapGenerator:
      - name: config
        namespace: test
        files:
          - kustomization.yaml
    patchesStrategicMerge:
      - patch.yaml
kind: ConfigMap
metadata:
  name: config-9mk7dchmcb
  namespace: test
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_pod_test.yaml
apiVersion: v1
kind: Pod
metadata:
  name: test
  namespace: test
spec:
  containers:
  - image: bash
    name: test
  volumes:
  - configMap:
      name: config-9mk7dchmcb
    name: config
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

