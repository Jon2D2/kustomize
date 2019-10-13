# Feature Test for Issue 1251_b


This folder contains files describing how to address [Issue 1251_b](https://github.com/kubernetes-sigs/kustomize/issues/1251_b)

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
mkdir -p ${DEMO_HOME}/base/backend
mkdir -p ${DEMO_HOME}/base/frontend
mkdir -p ${DEMO_HOME}/development
mkdir -p ${DEMO_HOME}/development/backend
mkdir -p ${DEMO_HOME}/development/frontend
mkdir -p ${DEMO_HOME}/production
mkdir -p ${DEMO_HOME}/production/backend
mkdir -p ${DEMO_HOME}/production/frontend
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/backend/kustomization.yaml
resources:
  - deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/frontend/kustomization.yaml
resources:
  - deployment.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ./backend/deployment.yaml
  - ./frontend/deployment.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/development/backend/kustomization.yaml
resources:
- ../../base/backend

patchesStrategicMerge:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/development/frontend/kustomization.yaml
resources:
- ../../base/frontend

patchesStrategicMerge:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/development/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ./backend
- ./frontend

commonLabels:
  env: dev
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/backend/kustomization.yaml
resources:
  - ../../base

patchesStrategicMerge:
  - deployment.yaml
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/frontend/kustomization.yaml
# resources:
#  - ../../base

patchesStrategicMerge:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile8

<!-- @createKustomizationFile8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../base

patchesStrategicMerge:
- ./backend/deployment.yaml
- ./frontend/deployment.yaml

commonLabels:
  env: production
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/backend/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-backend
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: my-backend
          image: my-backend-image
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/frontend/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-frontend
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: my-frontend
          image: my-frontend-image
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/development/backend/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-backend
spec:
  replicas: 2
  template:
    spec:
      containers:
        - name: my-backend
          image: my-backend-image
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/development/frontend/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-frontend
spec:
  replicas: 2
  template:
    spec:
      containers:
        - name: my-frontend
          image: my-frontend-image
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/backend/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-backend
spec:
  replicas: 3
  template:
    spec:
      containers:
        - name: my-backend
          image: my-backend-image
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/frontend/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-frontend
spec:
  replicas: 3
  template:
    spec:
      containers:
        - name: my-frontend
          image: my-frontend-image
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir -p ${DEMO_HOME}/actual
mkdir -p ${DEMO_HOME}/actual/development
mkdir -p ${DEMO_HOME}/actual/production
kustomize build ${DEMO_HOME}/development -o ${DEMO_HOME}/actual/development
kustomize build ${DEMO_HOME}/production -o ${DEMO_HOME}/actual/production
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir -p ${DEMO_HOME}/expected
mkdir -p ${DEMO_HOME}/expected/development
mkdir -p ${DEMO_HOME}/expected/production
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/development/apps_v1_deployment_my-backend.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    env: dev
  name: my-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      env: dev
  template:
    metadata:
      labels:
        env: dev
    spec:
      containers:
      - image: my-backend-image
        name: my-backend
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/development/apps_v1_deployment_my-frontend.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    env: dev
  name: my-frontend
spec:
  replicas: 2
  selector:
    matchLabels:
      env: dev
  template:
    metadata:
      labels:
        env: dev
    spec:
      containers:
      - image: my-frontend-image
        name: my-frontend
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production/apps_v1_deployment_my-backend.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    env: production
  name: my-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      env: production
  template:
    metadata:
      labels:
        env: production
    spec:
      containers:
      - image: my-backend-image
        name: my-backend
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production/apps_v1_deployment_my-frontend.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    env: production
  name: my-frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      env: production
  template:
    metadata:
      labels:
        env: production
    spec:
      containers:
      - image: my-frontend-image
        name: my-frontend
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual/development $DEMO_HOME/expected/development | wc -l); \
echo $?
```

```bash
test 0 == \
$(diff -r $DEMO_HOME/actual/production $DEMO_HOME/expected/production | wc -l); \
echo $?
```

