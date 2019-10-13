# Feature Test for Issue 1268


This folder contains files describing how to address [Issue 1268](https://github.com/kubernetes-sigs/kustomize/issues/1268)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- rollout.yaml
- deployment.yaml
- values.yaml

configurations:
- kustomizeconfig.yaml

vars:
  - name: DEPLOYMENT_COLOR
    objref:
      apiVersion: kustomize.config.k8s.io/v1
      kind: Values
      name: file1
    fieldref:
      fieldpath: spec.deploymentColor
EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig.yaml
varReference:
- path: spec/template/spec/containers[]/env[]/value
  kind: Rollout
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: helloworld
  name: helloworld
spec:
  replicas: 2
  selector:
    matchLabels:
      app: helloworld
  strategy:
    blueGreen: $(Values.file1.spec.blueGreen)
  template:
    metadata:
      labels:
        app: helloworld
    spec:
      containers:
      - env:
        - name: DEPLOY_VERSION
          value: $(DEPLOYMENT_COLOR)
        image: someplage/hellopy
        imagePullPolicy: Always
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/rollout.yaml
---
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  labels:
    app: helloworld
  name: helloworld
spec:
  replicas: 2
  selector:
    matchLabels:
      app: helloworld
  strategy:
    blueGreen: $(Values.file1.spec.blueGreen)
  template:
    metadata:
      labels:
        app: helloworld
    spec:
      containers:
      - env:
        - name: DEPLOY_VERSION
          value: $(DEPLOYMENT_COLOR)
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: file1
spec:
  deploymentColor: orange
  blueGreen:
     activeService: helloworld-svc-active
     autoPromotionEnabled: false
     previewService: helloworld-svc-preview
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_helloworld.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: helloworld
  name: helloworld
spec:
  replicas: 2
  selector:
    matchLabels:
      app: helloworld
  strategy:
    blueGreen:
      activeService: helloworld-svc-active
      autoPromotionEnabled: false
      previewService: helloworld-svc-preview
  template:
    metadata:
      labels:
        app: helloworld
    spec:
      containers:
      - env:
        - name: DEPLOY_VERSION
          value: orange
        image: someplage/hellopy
        imagePullPolicy: Always
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/argoproj.io_v1alpha1_rollout_helloworld.yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  labels:
    app: helloworld
  name: helloworld
spec:
  replicas: 2
  selector:
    matchLabels:
      app: helloworld
  strategy:
    blueGreen:
      activeService: helloworld-svc-active
      autoPromotionEnabled: false
      previewService: helloworld-svc-preview
  template:
    metadata:
      labels:
        app: helloworld
    spec:
      containers:
      - env:
        - name: DEPLOY_VERSION
          value: orange
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kustomize.config.k8s.io_v1_values_file1.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: file1
spec:
  blueGreen:
    activeService: helloworld-svc-active
    autoPromotionEnabled: false
    previewService: helloworld-svc-preview
  deploymentColor: orange
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

