# Feature Test for Issue 1442


This folder contains files describing how to address [Issue 1442](https://github.com/kubernetes-sigs/kustomize/issues/1442)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/a
mkdir -p ${DEMO_HOME}/b
mkdir -p ${DEMO_HOME}/c
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/a/kustomization.yaml
configMapGenerator:
- name: cm
  literals:
  - FOO=BAR
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/b/kustomization.yaml
resources:
- ../a
namePrefix: renamed-
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/c/kustomization.yaml
resources:
- ../b
configMapGenerator:
- name: cm
  behavior: merge
  literals:
  - HELLO=WORLD
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/c -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_renamed-cm-27kg5g5b5c.yaml
apiVersion: v1
data:
  FOO: BAR
  HELLO: WORLD
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: renamed-cm-27kg5g5b5c
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

