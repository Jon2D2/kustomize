# Feature Test for Issue 1540_b


This folder contains files describing how to address [Issue 1540_b](https://github.com/kubernetes-sigs/kustomize/issues/1540_b)

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
mkdir -p ${DEMO_HOME}/kustomizeconfig
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
- deployment.yaml
- values.yaml

vars:
- name: Values.my-values.spec.someip
  objref:
    apiVersion: kustomize.config.k8s.io/v1
    kind: Values
    name: my-values
  fieldref:
    fieldpath: spec.someip

configurations:
- kustomizeconfig/varreference.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: nfs-client-provisioner
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: nfs-client-provisioner
  template:
    metadata:
      labels:
        app: nfs-client-provisioner
    spec:
      serviceAccountName: nfs-client-provisioner
      containers:
        - name: nfs-client-provisioner
          image: quay.io/external_storage/nfs-client-provisioner:latest
          volumeMounts:
            - name: nfs-client-root
              mountPath: /persistentvolumes
          env:
            - name: PROVISIONER_NAME
              value: fuseim.pri/ifs
            - name: NFS_SERVER
              value: $(Values.my-values.spec.someip)
            - name: NFS_PATH
              value: /var/nfs
      volumes:
      - name: nfs-client-root
        nfs:
          server: $(Values.my-values.spec.someip)
          path: /var/nfs
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig/varreference.yaml
varReference:
- path: spec/template/spec/volumes[]/nfs
  kind: Deployment
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: my-values
spec:
  someip: 127.0.0.1
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_nfs-client-provisioner.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nfs-client-provisioner
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nfs-client-provisioner
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: nfs-client-provisioner
    spec:
      containers:
      - env:
        - name: PROVISIONER_NAME
          value: fuseim.pri/ifs
        - name: NFS_SERVER
          value: 127.0.0.1
        - name: NFS_PATH
          value: /var/nfs
        image: quay.io/external_storage/nfs-client-provisioner:latest
        name: nfs-client-provisioner
        volumeMounts:
        - mountPath: /persistentvolumes
          name: nfs-client-root
      serviceAccountName: nfs-client-provisioner
      volumes:
      - name: nfs-client-root
        nfs:
          path: /var/nfs
          server: 127.0.0.1
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kustomize.config.k8s.io_v1_values_my-values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: my-values
spec:
  someip: 127.0.0.1
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

