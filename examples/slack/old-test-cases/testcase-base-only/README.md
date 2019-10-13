# Feature Test for TestCase base-only


This folder contains files for old test-case base-only

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/in
mkdir -p ${DEMO_HOME}/in/resources
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/in/kustomization.yaml
namePrefix: team-foo-
commonLabels:
  app: mynginx
  org: example.com
  team: foo
commonAnnotations:
  note: This is a test annotation
resources:
  - resources/deployment.yaml
  - resources/networkpolicy.yaml
  - resources/service.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/in/resources/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/in/resources/networkpolicy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: nginx
spec:
  podSelector:
    matchExpressions:
      - {key: app, operator: In, values: [test]}
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: nginx
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/in/resources/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  ports:
    - port: 80
  selector:
    app: nginx
EOF
```

## Execution

<!-- @build @test -->
```bash
kustomize build ${DEMO_HOME}/in -o ${DEMO_HOME}/actual.yaml
```

## Verification


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    org: example.com
    team: foo
  name: team-foo-nginx
spec:
  ports:
  - port: 80
  selector:
    app: mynginx
    org: example.com
    team: foo
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    org: example.com
    team: foo
  name: team-foo-nginx
spec:
  selector:
    matchLabels:
      app: mynginx
      org: example.com
      team: foo
  template:
    metadata:
      annotations:
        note: This is a test annotation
      labels:
        app: mynginx
        org: example.com
        team: foo
    spec:
      containers:
      - image: nginx
        name: nginx
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    org: example.com
    team: foo
  name: team-foo-nginx
spec:
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: mynginx
          org: example.com
          team: foo
  podSelector:
    matchExpressions:
    - key: app
      operator: In
      values:
      - test
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual.yaml $DEMO_HOME/expected.yaml | wc -l); \
echo $?
```

