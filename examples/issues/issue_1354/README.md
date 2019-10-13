# Feature Test for Issue 1354


This folder contains files describing how to address [Issue 1354](https://github.com/kubernetes-sigs/kustomize/issues/1354)

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
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ./deployment.yaml

patchesStrategicMerge:
- ./patch-add-label.yaml
- ./patch-delete-container.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy1
  labels:
    mylabel: myapp
spec:
  selector:
    matchLabels:
      mylabel: myapp
  template:
    metadata:
      labels:
        mylabel: myapp
    spec:
      containers:
      - name: container1
        image: image1:v1.0
      - name: container2
        image: image2:v1.0
      - name: container3
        image: image3:v1.0
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy2
  labels:
    mylabel: myapp
spec:
  selector:
    matchLabels:
      mylabel: myapp
  template:
    metadata:
      labels:
        mylabel: myapp
    spec:
      containers:
      - name: container1
        image: image1:v1.0
      - name: container2
        image: image2:v1.0
      - name: container3
        image: image3:v1.0

EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/patch-add-label.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy1
spec:
  template:
    spec:
      containers:
      - name: container1
        env:
        - name: SOME_NAME
          value: somevalue
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy2
spec:
  template:
    spec:
      containers:
      - name: container1
        env:
        - name: SOME_NAME
          value: somevalue
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/patch-delete-container.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy1
spec:
  template:
    spec:
      containers:
      - $patch: delete
        name: container2
      - $patch: delete
        name: container3
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy2
spec:
  template:
    spec:
      containers:
      - $patch: delete
        name: container2
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_deploy1.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    mylabel: myapp
  name: deploy1
spec:
  selector:
    matchLabels:
      mylabel: myapp
  template:
    metadata:
      labels:
        mylabel: myapp
    spec:
      containers:
      - env:
        - name: SOME_NAME
          value: somevalue
        image: image1:v1.0
        name: container1
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_deploy2.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    mylabel: myapp
  name: deploy2
spec:
  selector:
    matchLabels:
      mylabel: myapp
  template:
    metadata:
      labels:
        mylabel: myapp
    spec:
      containers:
      - env:
        - name: SOME_NAME
          value: somevalue
        image: image1:v1.0
        name: container1
      - image: image3:v1.0
        name: container3
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

