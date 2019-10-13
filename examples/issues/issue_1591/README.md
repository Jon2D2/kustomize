# Feature Test for Issue 1591


This folder contains files describing how to address [Issue 1591](https://github.com/kubernetes-sigs/kustomize/issues/1591)

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

commonLabels:
  app: api

configurations:
- commonlabels.yaml

resources:
- application.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/application.yaml
apiVersion: app.k8s.io/v1beta1
kind: Application
metadata:
  name: app-example
spec:
  descriptor:
    type: simple
    version: v1.2.3
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/commonlabels.yaml
commonLabels:
- path: spec/selector/matchLabels
  create: true
  group: app.k8s.io
  version: v1beta1
  kind: Application
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
cat <<'EOF' >${DEMO_HOME}/expected/app.k8s.io_v1beta1_application_app-example.yaml
apiVersion: app.k8s.io/v1beta1
kind: Application
metadata:
  labels:
    app: api
  name: app-example
spec:
  descriptor:
    type: simple
    version: v1.2.3
  selector:
    matchLabels:
      app: api
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

