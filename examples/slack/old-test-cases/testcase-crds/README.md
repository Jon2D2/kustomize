# Feature Test for TestCase crds


This folder contains files for old test-case crds

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/crd
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/crd/kustomization.yaml
crds:
- mycrd.json

resources:
- secret.yaml
- mykind.yaml
- bee.yaml

namePrefix: test-EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/crd/bee.yaml
apiVersion: v1beta1
kind: Bee
metadata:
  name: bee
spec:
  action: fly
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/crd/mykind.yaml
apiVersion: jingfang.example.com/v1beta1
kind: MyKind
metadata:
  name: mykind
spec:
  secretRef:
    name: crdsecret
  beeRef:
    name: bee
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/crd/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: crdsecret
data:
  PATH: YmJiYmJiYmIK
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/crd/mycrd.json
EOF
```

## Execution

<!-- @build @test -->
```bash
# kustomize build ${DEMO_HOME}/crd -o ${DEMO_HOME}/actual.yaml
```

## Verification


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected.yaml
apiVersion: v1
data:
  PATH: YmJiYmJiYmIK
kind: Secret
metadata:
  name: test-crdsecret
---
apiVersion: jingfang.example.com/v1beta1
kind: MyKind
metadata:
  name: test-mykind
spec:
  beeRef:
    name: test-bee
  secretRef:
    name: test-crdsecret
---
apiVersion: v1beta1
kind: Bee
metadata:
  name: test-bee
spec:
  action: fly
EOF
```


<!-- @compareActualToExpected @test -->
```bash
```

