# Feature Test for Issue 1113


This folder contains files describing how to address [Issue 1113](https://github.com/kubernetes-sigs/kustomize/issues/1113)

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
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- service.yaml
- deployment.yaml
- values.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
apiVersion: apps/v1beta2 # for versions before 1.9.0 use apps/v1beta2
kind: MyDeployment
metadata:
  name: mysql
  labels:
    app: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  strategy:
    type: $(Values.file1.spec.strategy)
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - image: mysql:5.6
        name: mysql
        ports:
        - containerPort: $(Values.file1.spec.port)
          name: mysql
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/service.yaml
apiVersion: v1
kind: MyService
metadata:
  name: mysql
  labels:
    app: mysql
spec:
  ports:
    - port: $(Values.file1.spec.port)
  selector:
    app: mysql
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: file1
spec:
  port: 3306
  strategy: Recreate
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1beta2_mydeployment_mysql.yaml
apiVersion: apps/v1beta2
kind: MyDeployment
metadata:
  labels:
    app: mysql
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - image: mysql:5.6
        name: mysql
        ports:
        - containerPort: 3306
          name: mysql
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_myservice_mysql.yaml
apiVersion: v1
kind: MyService
metadata:
  labels:
    app: mysql
  name: mysql
spec:
  ports:
  - port: 3306
  selector:
    app: mysql
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kustomize.config.k8s.io_v1_values_file1.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: file1
spec:
  port: 3306
  strategy: Recreate
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

