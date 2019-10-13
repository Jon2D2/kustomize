# Feature Test for Issue 1190


This folder contains files describing how to address [Issue 1190](https://github.com/kubernetes-sigs/kustomize/issues/1190)

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
mkdir -p ${DEMO_HOME}/bar
mkdir -p ${DEMO_HOME}/foo
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bar/kustomization.yaml
resources:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/foo/kustomization.yaml
resources:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
  - foo
  - bar
namePrefix: inlining-example-
vars:
  - name: CUSTOM_TEMPLATE
    objref:
      kind: Deployment
      name: foo
      apiVersion: apps/v1beta2
    fieldref:
      fieldpath: spec.template
configurations:
  - transformer.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bar/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: bar
  labels:
    app: bar
spec:
  selector:
    matchLabels:
      app: bar
  template: $(CUSTOM_TEMPLATE)
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/foo/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: foo
  labels:
    app: foo
spec:
  selector:
    matchLabels:
      app: foo
  template:
    metadata:
      labels:
        app: foo
    spec:
      containers:
      - image: alpine
        name: foo
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/transformer.yaml
varReference:
- kind: Deployment
  path: spec/template
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1beta2_deployment_inlining-example-bar.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: bar
  name: inlining-example-bar
spec:
  selector:
    matchLabels:
      app: bar
  template:
    metadata:
      labels:
        app: foo
    spec:
      containers:
      - image: alpine
        name: foo
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1beta2_deployment_inlining-example-foo.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: foo
  name: inlining-example-foo
spec:
  selector:
    matchLabels:
      app: foo
  template:
    metadata:
      labels:
        app: foo
    spec:
      containers:
      - image: alpine
        name: foo
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

