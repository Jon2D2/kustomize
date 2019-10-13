# Feature Test for Issue pv-example


This example here was used in a Slack thread.
Original code is available [here](https://github.com/bobbyrward/kustomize-example)

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
mkdir -p ${DEMO_HOME}/prod
mkdir -p ${DEMO_HOME}/qa
mkdir -p ${DEMO_HOME}/stage
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
  - pv.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/prod/kustomization.yaml
bases:
  - ../base
namePrefix: prod-
patchesStrategicMerge:
  - change-volume-handle.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/qa/kustomization.yaml
bases:
  - ../base
namePrefix: qa-
patchesStrategicMerge:
  - change-volume-handle.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/stage/kustomization.yaml
bases:
  - ../base
namePrefix: stage-
patchesStrategicMerge:
  - change-volume-handle.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/pv.yaml
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mms-app-efs
  namespace: mms-app
spec:
  capacity:
    storage: 10Gi
  volumeMode: Filesystem
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  storageClassName: efs
  csi:
    driver: efs.csi.aws.com
    volumeHandle: fs-999aaa999
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/prod/change-volume-handle.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mms-app-efs
  namespace: mms-app
spec:
  cis:
    volumeHandle: fs-prodprodprod
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/qa/change-volume-handle.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mms-app-efs
  namespace: mms-app
spec:
  cis:
    volumeHandle: fs-qaqaqa
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/stage/change-volume-handle.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mms-app-efs
  namespace: mms-app
spec:
  cis:
    volumeHandle: fs-stagestagestage
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/qa -o ${DEMO_HOME}/actual/qa.yaml
kustomize build ${DEMO_HOME}/stage -o ${DEMO_HOME}/actual/stage.yaml
kustomize build ${DEMO_HOME}/prod -o ${DEMO_HOME}/actual/prod.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prod.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: prod-mms-app-efs
  namespace: mms-app
spec:
  accessModes:
  - ReadWriteMany
  capacity:
    storage: 10Gi
  cis:
    volumeHandle: fs-prodprodprod
  csi:
    driver: efs.csi.aws.com
    volumeHandle: fs-999aaa999
  persistentVolumeReclaimPolicy: Retain
  storageClassName: efs
  volumeMode: Filesystem
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/qa.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: qa-mms-app-efs
  namespace: mms-app
spec:
  accessModes:
  - ReadWriteMany
  capacity:
    storage: 10Gi
  cis:
    volumeHandle: fs-qaqaqa
  csi:
    driver: efs.csi.aws.com
    volumeHandle: fs-999aaa999
  persistentVolumeReclaimPolicy: Retain
  storageClassName: efs
  volumeMode: Filesystem
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/stage.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: stage-mms-app-efs
  namespace: mms-app
spec:
  accessModes:
  - ReadWriteMany
  capacity:
    storage: 10Gi
  cis:
    volumeHandle: fs-stagestagestage
  csi:
    driver: efs.csi.aws.com
    volumeHandle: fs-999aaa999
  persistentVolumeReclaimPolicy: Retain
  storageClassName: efs
  volumeMode: Filesystem
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

