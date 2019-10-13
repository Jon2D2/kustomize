# Feature Test for Issue 0800


This folder contains files describing how to address [Issue 0800](https://github.com/kubernetes-sigs/kustomize/issues/0800)

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
mkdir -p ${DEMO_HOME}/base1
mkdir -p ${DEMO_HOME}/base1/base
mkdir -p ${DEMO_HOME}/base1/overlay
mkdir -p ${DEMO_HOME}/base2
mkdir -p ${DEMO_HOME}/base2/base
mkdir -p ${DEMO_HOME}/base2/overlay
mkdir -p ${DEMO_HOME}/postgres
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base1/base/kustomization.yaml
bases:
  - ../../postgres

namePrefix: base1-
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base1/overlay/kustomization.yaml
bases:
  - ../base

namePrefix: overlay-
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base2/base/kustomization.yaml
bases:
  - ../../postgres

namePrefix: base2-
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base2/overlay/kustomization.yaml
bases:
  - ../base

namePrefix: overlay-
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
bases:
  - base1/overlay
  - base2/overlay
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/postgres/kustomization.yaml
resources:
  - resources.yaml

commonLabels:
  component: postgres
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/postgres/resources.yaml
---
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: postgres
spec:
  storageClassName: do-block-storage
  resources:
    requests:
      storage: 1Gi
  accessModes:
    - ReadWriteOnce
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
spec:
  selector:
    matchLabels: {}
  strategy:
    type: Recreate
  template:
    spec:
      containers:
        - name: postgres
          image: postgres
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - mountPath: /var/lib/postgresql
              name: data
          ports:
            - name: postgres
              containerPort: 5432
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: postgres
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_overlay-base1-postgres.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    component: postgres
  name: overlay-base1-postgres
spec:
  selector:
    matchLabels:
      component: postgres
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        component: postgres
    spec:
      containers:
      - image: postgres
        imagePullPolicy: IfNotPresent
        name: postgres
        ports:
        - containerPort: 5432
          name: postgres
        volumeMounts:
        - mountPath: /var/lib/postgresql
          name: data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: overlay-base1-postgres
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_overlay-base2-postgres.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    component: postgres
  name: overlay-base2-postgres
spec:
  selector:
    matchLabels:
      component: postgres
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        component: postgres
    spec:
      containers:
      - image: postgres
        imagePullPolicy: IfNotPresent
        name: postgres
        ports:
        - containerPort: 5432
          name: postgres
        volumeMounts:
        - mountPath: /var/lib/postgresql
          name: data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: overlay-base2-postgres
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_persistentvolumeclaim_overlay-base1-postgres.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  labels:
    component: postgres
  name: overlay-base1-postgres
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: do-block-storage
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_persistentvolumeclaim_overlay-base2-postgres.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  labels:
    component: postgres
  name: overlay-base2-postgres
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: do-block-storage
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

