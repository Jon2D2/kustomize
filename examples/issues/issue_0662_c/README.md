# Feature Test for Issue 0662


This folder contains files describing how to address [Issue 0662](https://github.com/kubernetes-sigs/kustomize/issues/0662)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/base
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- ./deploy.yaml
namePrefix: lalala-

secretGenerator:
- name: xxx
  literals:
  - password=123456
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
- ./base

secretGenerator:
- name: xxx
  behavior: merge
  literals:
  - password=12345699999

configmapGenerator:
- name: yyy
  literals:
  - password=123456
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deploy.yaml
kind: Deployment
metadata:
  name: foobar
spec:
  template:
    spec:
      containers:
      - name: foobar
        image: busybox
        envFrom:
        - secretRef:
            name: xxx
        - configMapRef:
            name: yyy
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_yyy-bkfmbb8t66.yaml
apiVersion: v1
data:
  password: "123456"
kind: ConfigMap
metadata:
  name: yyy-bkfmbb8t66
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_secret_lalala-xxx-ctd2ckb7h7.yaml
apiVersion: v1
data:
  password: MTIzNDU2OTk5OTk=
kind: Secret
metadata:
  annotations: {}
  labels: {}
  name: lalala-xxx-ctd2ckb7h7
type: Opaque
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_~v_deployment_lalala-foobar.yaml
kind: Deployment
metadata:
  name: lalala-foobar
spec:
  template:
    spec:
      containers:
      - envFrom:
        - secretRef:
            name: lalala-xxx-ctd2ckb7h7
        - configMapRef:
            name: yyy-bkfmbb8t66
        image: busybox
        name: foobar
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

