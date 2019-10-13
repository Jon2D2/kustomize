# Feature Test for Issue 1592


This folder contains files describing how to address [Issue 1592](https://github.com/kubernetes-sigs/kustomize/issues/1592)

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
  app: centraldashboard

configMapGenerator:
- envs:
  - params.env
  name: parameters

vars:
- fieldref:
    fieldPath: data.registry
  name: registry
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: parameters
- fieldref:
    fieldPath: data.centraldashboardImageName
  name: centraldashboardImageName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: parameters
- fieldref:
    fieldPath: data.centraldashboardImageTag
  name: centraldashboardImageTag
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: parameters

images:
- name: defaultCentralDashboardImageName
  newName: $(registry)/$(centraldashboardImageName)
  newTag: $(centraldashboardImageTag)

resources:
- deployment.yaml

configurations:
- kustomizeconfig.yaml

EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig.yaml
varReference:
- kind: Deployment
  path: spec/template/spec/containers[]/image
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
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
        image: defaultCentralDashboardImageName:latest
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/params.env
registry=gcr.io
centraldashboardImageName=kubeflow-images-public/centraldashboard
centraldashboardImageTag=v20190823-v0.6.0-rc.0-69-gcb7dab59
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
    app: centraldashboard
  name: deployment
spec:
  selector:
    matchLabels:
      app: centraldashboard
  template:
    metadata:
      labels:
        app: centraldashboard
    spec:
      containers:
      - image: gcr.io/kubeflow-images-public/centraldashboard:v20190823-v0.6.0-rc.0-69-gcb7dab59
        name: main
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_parameters-2bfg9t98bh.yaml
apiVersion: v1
data:
  centraldashboardImageName: kubeflow-images-public/centraldashboard
  centraldashboardImageTag: v20190823-v0.6.0-rc.0-69-gcb7dab59
  registry: gcr.io
kind: ConfigMap
metadata:
  labels:
    app: centraldashboard
  name: parameters-2bfg9t98bh
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

