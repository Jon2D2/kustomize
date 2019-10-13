# Feature Test for Issue 0915


This folder contains files describing how to address [Issue 0915](https://github.com/kubernetes-sigs/kustomize/issues/0915)

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
mkdir -p ${DEMO_HOME}/common
mkdir -p ${DEMO_HOME}/overlay
mkdir -p ${DEMO_HOME}/parent
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/kustomization.yaml
images:     
- name: cr.agilicus.com/group/image
  newName: cr.agilicus.com/group/image
  newTag: 1.0.0
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/kustomization.yaml
namespace: X

resources:
- ../parent
- ../common
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/parent/kustomization.yaml
namespace: X

resources:
- dep.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/parent/dep.yaml
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: X
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: X
          image: cr.agilicus.com/group/image

EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlay -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1beta1_deployment_x.yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: X
  namespace: X
spec:
  replicas: 1
  template:
    spec:
      containers:
      - image: cr.agilicus.com/group/image
        name: X
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

