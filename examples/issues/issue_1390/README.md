# Feature Test for Issue 1390


This folder contains files describing how to address [Issue 1390](https://github.com/kubernetes-sigs/kustomize/issues/1390)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/autovar
mkdir -p ${DEMO_HOME}/manual
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/autovar/kustomization.yaml
namePrefix: pfx-

resources:
- deployment.yaml
- secret.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/manual/kustomization.yaml
namePrefix: pfx-

resources:
- deployment.yaml
- secret.yaml

configurations:
- kustomizeconfig.yaml

vars:
- name : STORAGE_SECRET
  objref:
    kind: Secret
    name: storage-secret
    apiVersion: v1
  fieldref:
    fieldpath: metadata.name
EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/manual/kustomizeconfig.yaml
varReference:
- kind: Deployment
  path: spec/template/spec/volumes[]/azureFile/secretName
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/autovar/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
          - name: app
            volumeMounts:
            - name: appMnt
              mountPath: /data
      volumes:
      - name: appMnt
        azureFile:
          secretName: $(Secret.storage-secret.metadata.name)
          shareName: appMnt-share
          readOnly: true
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/autovar/secret.yaml
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  name: storage-secret
type: Opaque
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/manual/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
          - name: app
            volumeMounts:
            - name: appMnt
              mountPath: /data
      volumes:
      - name: appMnt
        azureFile:
          secretName: $(STORAGE_SECRET)
          shareName: appMnt-share
          readOnly: true
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/manual/secret.yaml
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  name: storage-secret
type: Opaque
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir -p ${DEMO_HOME}/actual
mkdir -p ${DEMO_HOME}/actual/manual
mkdir -p ${DEMO_HOME}/actual/autovar
kustomize build ${DEMO_HOME}/manual -o ${DEMO_HOME}/actual/manual
kustomize build ${DEMO_HOME}/autovar -o ${DEMO_HOME}/actual/autovar
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir -p ${DEMO_HOME}/expected
mkdir -p ${DEMO_HOME}/expected/manual
mkdir -p ${DEMO_HOME}/expected/autovar
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/autovar/apps_v1_deployment_pfx-app.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pfx-app
spec:
  template:
    spec:
      containers:
      - name: app
        volumeMounts:
        - mountPath: /data
          name: appMnt
      volumes:
      - azureFile:
          readOnly: true
          secretName: pfx-storage-secret
          shareName: appMnt-share
        name: appMnt
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/autovar/~g_v1_secret_pfx-storage-secret.yaml
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  name: pfx-storage-secret
type: Opaque
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/manual/apps_v1_deployment_pfx-app.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pfx-app
spec:
  template:
    spec:
      containers:
      - name: app
        volumeMounts:
        - mountPath: /data
          name: appMnt
      volumes:
      - azureFile:
          readOnly: true
          secretName: pfx-storage-secret
          shareName: appMnt-share
        name: appMnt
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/manual/~g_v1_secret_pfx-storage-secret.yaml
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  name: pfx-storage-secret
type: Opaque
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

