# Feature Test for Issue 0662


This folder contains files describing how to address [Issue 0662](https://github.com/kubernetes-sigs/kustomize/issues/0662)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/example
mkdir -p ${DEMO_HOME}/example/base
mkdir -p ${DEMO_HOME}/example/base/otherbase
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/example/base/kustomization.yaml
resources:
- otherbase
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/example/base/otherbase/kustomization.yaml
resources:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/example/kustomization.yaml
resources:
- base

secretGenerator:
- name: example
  literals:
  - HELLO=world
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/example/base/otherbase/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example
spec:
  template:
    spec:
      containers:
      - name: example
        image: example
        envFrom:
        - secretRef:
            name: example
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/example -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_example.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example
spec:
  template:
    spec:
      containers:
      - envFrom:
        - secretRef:
            name: example-5cm72hhtk8
        image: example
        name: example
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_secret_example-5cm72hhtk8.yaml
apiVersion: v1
data:
  HELLO: d29ybGQ=
kind: Secret
metadata:
  name: example-5cm72hhtk8
type: Opaque
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

