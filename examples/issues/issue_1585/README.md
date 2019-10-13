# Feature Test for Issue 1585


This folder contains files describing how to address [Issue 1585](https://github.com/kubernetes-sigs/kustomize/issues/1585)

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
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- configmap.yaml
- deployment.yaml

commonLabels:
  app.kubernetes.io\/instance: $(ConfigMap.cm.data.myvar)
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm
data:
  myvar: var-value-in-config-map
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
spec:
  template:
    spec:
      containers:
      - name: main
        image: myapp
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io\/instance: var-value-in-config-map
  name: deployment
spec:
  selector:
    matchLabels:
      app.kubernetes.io\/instance: var-value-in-config-map
  template:
    metadata:
      labels:
        app.kubernetes.io\/instance: var-value-in-config-map
    spec:
      containers:
      - image: myapp
        name: main
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_cm.yaml
apiVersion: v1
data:
  myvar: var-value-in-config-map
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io\/instance: var-value-in-config-map
  name: cm
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

