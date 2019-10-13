# Feature Test for Issue 1501


This folder contains files describing how to address [Issue 1501](https://github.com/kubernetes-sigs/kustomize/issues/1501)

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
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
configMapGenerator:
- name: a-config-map
  envs:
  - a.properties
- name: b-config-map
  envs:
  - b.properties

# generatorOptions:
#   disableNameSuffixHash: true

patchesStrategicMerge:
- patch-a.yaml
- patch-b.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/patch-a.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: a-config-map
  labels:
    apps: a
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/patch-b.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: b-config-map
  labels:
    apps: b
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/a.properties
a=0
EOF
```


### Preparation Step Other1

<!-- @createOther1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/b.properties
b=1
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_a-config-map-26kgmbk2md.yaml
apiVersion: v1
data:
  a: "0"
kind: ConfigMap
metadata:
  labels:
    apps: a
  name: a-config-map-26kgmbk2md
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_b-config-map-gh56k7ggb5.yaml
apiVersion: v1
data:
  b: "1"
kind: ConfigMap
metadata:
  labels:
    apps: b
  name: b-config-map-gh56k7ggb5
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

