# Feature Test for Issue 1357


This folder contains files describing how to address [Issue 1357](https://github.com/kubernetes-sigs/kustomize/issues/1357)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
- pod.yaml
images:
- name: busybox
  newName: alpine
  digest: sha256:2222222222222222
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
spec:
  containers:
  - name: myapp-container
    image: busybox:1.29.0
    command: ['sh', '-c', 'echo The app is running! && sleep 3600']
  initContainers:
  - name: init-mydb
    image: busybox@sha256:1111111111111111
    command: ['sh', '-c', 'until nslookup mydb; do echo waiting for mydb; sleep 2; done;']
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_pod_myapp-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: myapp
  name: myapp-pod
spec:
  containers:
  - command:
    - sh
    - -c
    - echo The app is running! && sleep 3600
    image: alpine@sha256:2222222222222222
    name: myapp-container
  initContainers:
  - command:
    - sh
    - -c
    - until nslookup mydb; do echo waiting for mydb; sleep 2; done;
    image: alpine@sha256:2222222222222222
    name: init-mydb
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

