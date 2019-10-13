# Feature Test for Issue 1347


This folder contains files describing how to address [Issue 1347](https://github.com/kubernetes-sigs/kustomize/issues/1347)

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
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
  - ./base/demo.yaml
namespace: demo1
namePrefix: prefix-
nameSuffix: -suffix
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/demo.yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: storage
  labels:
    k8s-addon: storage-aws.addons.k8s.io
provisioner: kubernetes.io/aws-ebs
reclaimPolicy: Retain
allowVolumeExpansion: true
parameters:
  type: gp2
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: demo
  namespace: demo
spec:
  selector:
    matchLabels:
      app: demo
  replicas: 3
  updateStrategy:
    type: RollingUpdate
  podManagementPolicy: Parallel
  template:
    metadata:
      labels:
        app: demo
      annotations:
    spec:
      containers:
        - name: demo
          image: alpine:3.9
          volumeMounts:
            - name: data
              mountPath: /data
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes: [ "ReadWriteOnce" ]
        storageClassName: storage
        resources:
          requests:
            storage: 10Gi
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_statefulset_prefix-demo-suffix.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: prefix-demo-suffix
  namespace: demo1
spec:
  podManagementPolicy: Parallel
  replicas: 3
  selector:
    matchLabels:
      app: demo
  template:
    metadata:
      annotations: null
      labels:
        app: demo
    spec:
      containers:
      - image: alpine:3.9
        name: demo
        volumeMounts:
        - mountPath: /data
          name: data
  updateStrategy:
    type: RollingUpdate
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 10Gi
      storageClassName: prefix-storage-suffix
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/storage.k8s.io_v1_storageclass_prefix-storage-suffix.yaml
allowVolumeExpansion: true
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  labels:
    k8s-addon: storage-aws.addons.k8s.io
  name: prefix-storage-suffix
parameters:
  type: gp2
provisioner: kubernetes.io/aws-ebs
reclaimPolicy: Retain
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

