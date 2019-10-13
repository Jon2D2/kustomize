# Feature Test for Issue 1381


This folder contains files describing how to address [Issue 1381](https://github.com/kubernetes-sigs/kustomize/issues/1381)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/combined
mkdir -p ${DEMO_HOME}/persistentvolume
mkdir -p ${DEMO_HOME}/storageclass
mkdir -p ${DEMO_HOME}/validatingwebhook
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/combined/kustomization.yaml
resources:
- ../validatingwebhook
- ../persistentvolume
- ../storageclass
- clusterrole.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/persistentvolume/kustomization.yaml
nameprefix: persistentvolumes-

configurations:
- kustomizeconfig.yaml

resources:
- persistentvolume.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/storageclass/kustomization.yaml
nameprefix: storageclass-

configurations:
- kustomizeconfig.yaml

resources:
- storageclass.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/validatingwebhook/kustomization.yaml
nameprefix: validatingwebhook-

configurations:
- kustomizeconfig.yaml

resources:
- validatingwebhook.yaml
EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/persistentvolume/kustomizeconfig.yaml
nameReference:
- kind: PersistentVolume
  version: v1
  fieldSpecs:
  - path: rules/resourceNames
    kind: ClusterRole
EOF
```


### Preparation Step KustomizeConfig1

<!-- @createKustomizeConfig1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/storageclass/kustomizeconfig.yaml
nameReference:
- kind: StorageClass
  version: v1
  group: storage.k8s.io
  fieldSpecs:
  - path: rules/resourceNames
    kind: ClusterRole
EOF
```


### Preparation Step KustomizeConfig2

<!-- @createKustomizeConfig2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/validatingwebhook/kustomizeconfig.yaml
nameReference:
- kind: ValidatingWebhookConfiguration
  version: v1beta1
  group: admissionregistration.k8s.io
  fieldSpecs:
  - path: rules/resourceNames
    kind: ClusterRole
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/combined/clusterrole.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: orgname
rules:
  resourceNames:
  - orgname
  resources:
  - persistentvolumes
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/persistentvolume/persistentvolume.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: orgname
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/resources.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: orgname
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: orgname
rules:
  resourceNames:
  - orgname
  resources:
  - persistentvolumes
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: orgname
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/storageclass/storageclass.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: orgname
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/validatingwebhook/validatingwebhook.yaml
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: orgname
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/combined -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/admissionregistration.k8s.io_v1beta1_validatingwebhookconfiguration_validatingwebhook-orgname.yaml
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: validatingwebhook-orgname
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_persistentvolume_persistentvolumes-orgname.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: persistentvolumes-orgname
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/rbac.authorization.k8s.io_v1_clusterrole_orgname.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: orgname
rules:
  resourceNames:
  - validatingwebhook-orgname
  resources:
  - persistentvolumes
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/storage.k8s.io_v1_storageclass_storageclass-orgname.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: storageclass-orgname
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

