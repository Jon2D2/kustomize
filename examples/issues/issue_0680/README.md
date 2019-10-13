# Feature Test for Issue 0680


This folder contains files describing how to address [Issue 0680](https://github.com/kubernetes-sigs/kustomize/issues/0680)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/first-solution
mkdir -p ${DEMO_HOME}/second-solution
mkdir -p ${DEMO_HOME}/second-solution/base
mkdir -p ${DEMO_HOME}/second-solution/overlays
mkdir -p ${DEMO_HOME}/second-solution/overlays/prod
mkdir -p ${DEMO_HOME}/third-solution
mkdir -p ${DEMO_HOME}/third-solution/base
mkdir -p ${DEMO_HOME}/third-solution/overlays
mkdir -p ${DEMO_HOME}/third-solution/overlays/production
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/first-solution/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- resources.yaml
configMapGenerator:
- name: demo-configmap-parameters
  envs:
  - params.env
generatorOptions:
  disableNameSuffixHash: true
vars:
- name: foo
  objref:
    kind: ConfigMap
    name: demo-configmap-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.foo
- name: bar
  objref:
    kind: ConfigMap
    name: demo-configmap-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.bar
configurations:
- params.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/second-solution/base/kustomization.yaml
configMapGenerator:
- name: my-configmap
  files:
  - settings.cfg
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/second-solution/overlays/prod/kustomization.yaml

resources:
- ../../base

configMapGenerator:
- name: my-configmap
  behavior: merge
  files:
  - settings.cfg
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/third-solution/base/kustomization.yaml
resources:
- values.yaml
  
configMapGenerator:
- name: myJavaServerProps
  files:
  - app.properties
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/third-solution/overlays/production/kustomization.yaml
resources:
- ../../base

patchesStrategicMerge: 
- ./values.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/first-solution/params.yaml
varReference:
- path: data/settings
  kind: ConfigMap
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/first-solution/resources.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: demo-configmap
data:   
   settings:  |-
      {
         foo: $(foo),
         bar: $(bar)
      }
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/third-solution/base/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: file1
spec:
  foo: default-foo-value
  bar: default-bar-value
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/third-solution/overlays/production/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: file1
spec:
  foo: production-foo-value
  bar: production-foo-value
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/first-solution/params.env
foo=foo_1
bar=bar_2
EOF
```


### Preparation Step Other1

<!-- @createOther1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/second-solution/base/settings.cfg
Content of the base settings.cfg
EOF
```


### Preparation Step Other2

<!-- @createOther2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/second-solution/overlays/prod/settings.cfg
Content of the overlay settings.cfg
EOF
```


### Preparation Step Other3

<!-- @createOther3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/third-solution/base/app.properties
foo=$(Values.file1.spec.foo)
bar=$(Values.file1.spec.bar)
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/first-solution/ -o ${DEMO_HOME}/actual/first-solution.yaml
kustomize build ${DEMO_HOME}/second-solution/overlays/prod/ -o ${DEMO_HOME}/actual/second-solution.yaml
kustomize build ${DEMO_HOME}/third-solution/overlays/production/ -o ${DEMO_HOME}/actual/third-solution.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/first-solution.yaml
apiVersion: v1
data:
  settings: |-
    {
       foo: foo_1,
       bar: bar_2
    }
kind: ConfigMap
metadata:
  name: demo-configmap
---
apiVersion: v1
data:
  bar: bar_2
  foo: foo_1
kind: ConfigMap
metadata:
  name: demo-configmap-parameters
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/second-solution.yaml
apiVersion: v1
data:
  settings.cfg: |
    Content of the overlay settings.cfg
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: my-configmap-8gkkcddg96
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/third-solution.yaml
apiVersion: v1
data:
  app.properties: |
    foo=production-foo-value
    bar=production-foo-value
kind: ConfigMap
metadata:
  name: myJavaServerProps-7m6dhd66bg
---
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: file1
spec:
  bar: production-foo-value
  foo: production-foo-value
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

