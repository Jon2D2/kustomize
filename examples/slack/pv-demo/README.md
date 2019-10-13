# Feature Test for PersistentVolume slack thread


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
mkdir -p ${DEMO_HOME}/base
mkdir -p ${DEMO_HOME}/base/kustomizeconfig
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- pv.yaml
- values.yaml

secretGenerator:
- name: dumpster-blobfuse-credentials
  literals:
  - username=admin
  - password=somepw

vars:
- name: CONTAINER_NAME
  objref:
    apiVersion: kustomize.config.k8s.io/v1
    kind: Values
    name: file1
  fieldref:
    fieldpath: spec.containername

configurations:
- kustomizeconfig/pv.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/pv.yaml
nameReference:
- kind: Secret
  version: v1
  fieldSpecs:
  - path: spec/flexVolume/secretRef/name
    kind: PersistentVolume

varReference:
- path: spec/flexVolume/options/container
  kind: PersistentVolume
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/pv.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: dumpster-pv
spec:
  accessModes:
  - ReadWriteMany
  capacity:
    storage: 100Gi
  flexVolume:
    driver: azure/blobfuse
    options:
      container: $(CONTAINER_NAME)
      mountoptions: --file-cache-timeout-in-seconds=120
      tmppath: /tmp/blobfuse
    secretRef:
      name: dumpster-blobfuse-credentials
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: file1
spec:
  containername: pgdumps
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/base -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_persistentvolume_dumpster-pv.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: dumpster-pv
spec:
  accessModes:
  - ReadWriteMany
  capacity:
    storage: 100Gi
  flexVolume:
    driver: azure/blobfuse
    options:
      container: pgdumps
      mountoptions: --file-cache-timeout-in-seconds=120
      tmppath: /tmp/blobfuse
    secretRef:
      name: dumpster-blobfuse-credentials-b2khc622hm
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_secret_dumpster-blobfuse-credentials-b2khc622hm.yaml
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  name: dumpster-blobfuse-credentials-b2khc622hm
type: Opaque
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
  containername: pgdumps
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

