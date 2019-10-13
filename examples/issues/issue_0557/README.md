# Feature Test for Issue 0557


This folder contains files describing how to address [Issue 0557](https://github.com/kubernetes-sigs/kustomize/issues/0557)

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
mkdir -p ${DEMO_HOME}/conflict
mkdir -p ${DEMO_HOME}/conflict/overlay_a
mkdir -p ${DEMO_HOME}/conflict/overlay_b
mkdir -p ${DEMO_HOME}/works
mkdir -p ${DEMO_HOME}/works/overlay_1
mkdir -p ${DEMO_HOME}/works/overlay_2
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- resources.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/conflict/kustomization.yaml
resources:
- ./overlay_a
- ./overlay_b
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/conflict/overlay_a/kustomization.yaml
namePrefix: extra-

bases:
- ../../works/overlay_1
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/conflict/overlay_b/kustomization.yaml
namePrefix: extra-

bases:
- ../../works/overlay_2
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/works/kustomization.yaml
# These bases work fine
resources:
- ./overlay_1
- ./overlay_2
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/works/overlay_1/kustomization.yaml
namePrefix: overlay1-

bases:
- ../../base

# patches:
#   - patch.yaml
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/works/overlay_2/kustomization.yaml
namePrefix: overlay2-

bases:
- ../../base

# patches:
#   - patch.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/resources.yaml
---
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: data
spec:
  resources:
    requests:
      storage: 1Gi
  accessModes:
    - ReadWriteOnce
---
apiVersion: apps/v1
kind: Pod
metadata:
  name: pod
spec:
  containers:
    - name: nginx
      image: nginx
  volumes:
    - name: data
      persistentVolumeClaim:
        claimName: data
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir -p ${DEMO_HOME}/actual
mkdir -p ${DEMO_HOME}/actual/works
mkdir -p ${DEMO_HOME}/actual/conflict
kustomize build ${DEMO_HOME}/works -o ${DEMO_HOME}/actual/works
kustomize build ${DEMO_HOME}/conflict -o ${DEMO_HOME}/actual/conflict
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir -p ${DEMO_HOME}/expected
mkdir -p ${DEMO_HOME}/expected/works
mkdir -p ${DEMO_HOME}/expected/conflict
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/conflict/apps_v1_pod_extra-overlay1-pod.yaml
apiVersion: apps/v1
kind: Pod
metadata:
  name: extra-overlay1-pod
spec:
  containers:
  - image: nginx
    name: nginx
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: extra-overlay1-data
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/conflict/apps_v1_pod_extra-overlay2-pod.yaml
apiVersion: apps/v1
kind: Pod
metadata:
  name: extra-overlay2-pod
spec:
  containers:
  - image: nginx
    name: nginx
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: extra-overlay2-data
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/conflict/~g_v1_persistentvolumeclaim_extra-overlay1-data.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: extra-overlay1-data
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/conflict/~g_v1_persistentvolumeclaim_extra-overlay2-data.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: extra-overlay2-data
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/works/apps_v1_pod_overlay1-pod.yaml
apiVersion: apps/v1
kind: Pod
metadata:
  name: overlay1-pod
spec:
  containers:
  - image: nginx
    name: nginx
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: overlay1-data
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/works/apps_v1_pod_overlay2-pod.yaml
apiVersion: apps/v1
kind: Pod
metadata:
  name: overlay2-pod
spec:
  containers:
  - image: nginx
    name: nginx
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: overlay2-data
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/works/~g_v1_persistentvolumeclaim_overlay1-data.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: overlay1-data
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/works/~g_v1_persistentvolumeclaim_overlay2-data.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: overlay2-data
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

