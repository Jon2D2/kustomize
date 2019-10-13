# Feature Test for Issue 1295


This folder contains files describing how to address [Issue 1295](https://github.com/kubernetes-sigs/kustomize/issues/1295)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/common
mkdir -p ${DEMO_HOME}/component1
mkdir -p ${DEMO_HOME}/myapp
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization


resources:
- Chief.yaml

configurations:
- kustomizeconfig.yaml

vars:
- fieldref:
    fieldPath: data.batchSize
  name: batchSize
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/component1/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generatorOptions:
  disableNameSuffixHash: true

configMapGenerator:
- literals:
  - name=mnist-train-local
  - batchSize=100
  name: mnist-map-training
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/myapp/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../common
- ../component1
EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/kustomizeconfig.yaml
varReference:
- path: spec/tfReplicaSpecs/Chief/template/spec/containers/env/value
  kind: TFJob
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/common/Chief.yaml
apiVersion: kubeflow.org/v1beta2
kind: TFJob
metadata:
  name: training-name
spec:
  tfReplicaSpecs:
    Chief:
      replicas: 1
      template:
        spec:
          containers:
          - name: tensorflow
            command:
            - /usr/bin/python
            - /opt/model.py
            env:
            - name: batchSize
              value: $(batchSize)
            image: training-image
            workingDir: /opt
          restartPolicy: OnFailure
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/myapp -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_mnist-map-training.yaml
apiVersion: v1
data:
  batchSize: "100"
  name: mnist-train-local
kind: ConfigMap
metadata:
  name: mnist-map-training
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kubeflow.org_v1beta2_tfjob_training-name.yaml
apiVersion: kubeflow.org/v1beta2
kind: TFJob
metadata:
  name: training-name
spec:
  tfReplicaSpecs:
    Chief:
      replicas: 1
      template:
        spec:
          containers:
          - command:
            - /usr/bin/python
            - /opt/model.py
            env:
            - name: batchSize
              value: "100"
            image: training-image
            name: tensorflow
            workingDir: /opt
          restartPolicy: OnFailure
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

