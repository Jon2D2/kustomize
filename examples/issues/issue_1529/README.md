# Feature Test for Issue 1529


This folder contains files describing how to address [Issue 1529](https://github.com/kubernetes-sigs/kustomize/issues/1529)

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

resources:
- ./cm.yaml
- ./deployment.yaml

vars:
  - name: KUSTOMIZE_MYCONFIGMAP_VARIABLE
    objref:
      kind: ConfigMap
      name: myConfigMap
      apiVersion: v1
    fieldref:
      fieldpath: data[another.variable.with.dots.in.it]
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cm.yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: myConfigMap
data:
  a_variable: "a value"
  another.variable.with.dots.in.it: "another value"
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: nginx
        image: nginx
        env:
        - name: SOME_ENV
          value: $(KUSTOMIZE_MYCONFIGMAP_VARIABLE)
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME} -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_nginx.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 1
  template:
    spec:
      containers:
      - env:
        - name: SOME_ENV
          value: another value
        image: nginx
        name: nginx
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_myconfigmap.yaml
apiVersion: v1
data:
  a_variable: a value
  another.variable.with.dots.in.it: another value
kind: ConfigMap
metadata:
  name: myConfigMap
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

