# Feature Test for Issue 1563


This folder contains files describing how to address [Issue 1563](https://github.com/kubernetes-sigs/kustomize/issues/1563)

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
mkdir -p ${DEMO_HOME}/staging
mkdir -p ${DEMO_HOME}/staging/base
mkdir -p ${DEMO_HOME}/staging/overlay
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
  - deployment.yaml

secretGenerator:
  - name: mysecret
    literals:
      - key=value

EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/staging/base/kustomization.yaml
resources:
  - ../../base
  - another_deployment.yaml

namePrefix: staging-

EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/staging/overlay/kustomization.yaml
resources:
  - ../base

namePrefix: overlay-

EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
spec:
  template:
    spec:
      containers:
        - name: test
          image: test
          env:
            - name: MY_SECRET_KEY
              valueFrom:
                secretKeyRef:
                  name: mysecret
                  key: key

EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/staging/base/another_deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: another-deployment
spec:
  template:
    spec:
      containers:
        - name: another-test
          image: another-test
          env:
            - name: MY_SECRET_KEY
              valueFrom:
                secretKeyRef:
                  name: mysecret
                  key: key

EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/staging/overlay -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_overlay-staging-another-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: overlay-staging-another-deployment
spec:
  template:
    spec:
      containers:
      - env:
        - name: MY_SECRET_KEY
          valueFrom:
            secretKeyRef:
              key: key
              name: overlay-staging-mysecret-7m4thk7c67
        image: another-test
        name: another-test
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_overlay-staging-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: overlay-staging-deployment
spec:
  template:
    spec:
      containers:
      - env:
        - name: MY_SECRET_KEY
          valueFrom:
            secretKeyRef:
              key: key
              name: overlay-staging-mysecret-7m4thk7c67
        image: test
        name: test
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_secret_overlay-staging-mysecret-7m4thk7c67.yaml
apiVersion: v1
data:
  key: dmFsdWU=
kind: Secret
metadata:
  name: overlay-staging-mysecret-7m4thk7c67
type: Opaque
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

