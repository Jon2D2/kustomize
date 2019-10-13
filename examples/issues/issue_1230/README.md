# Feature Test for Issue 1230


This folder contains files describing how to address [Issue 1230](https://github.com/kubernetes-sigs/kustomize/issues/1230)

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
mkdir -p ${DEMO_HOME}/overlay
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- cm.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/kustomization.yaml
resources:
- ../base
patchesStrategicMerge:
- cm.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/cm.yaml
apiVersion: v1
kind: ConfigMap

metadata:
  name: hello-world-map
  namespace: hello-world-ns

data: {}
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/cm.yaml
apiVersion: v1
kind: ConfigMap

metadata:
  name: hello-world-map

data:
  ingress_host: example.org
EOF
```

## Execution

The following build calls triggers an exception in kustomize 2.1.0.

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
# kustomize build ${DEMO_HOME}/overlay -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

