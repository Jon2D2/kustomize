# Feature Test for Issue 1816


This folder contains files describing how to address [Issue 1816](https://github.com/kubernetes-sigs/kustomize/issues/1816)

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
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

configurations:
- commonlabels.yaml

commonLabels:
  app: nginx

resources:
- services.yaml
- deployments.yaml
- cronjob.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/commonlabels.yaml
commonLabels:
- path: spec/selector
  version: v1
  kind: Service
  behavior: remove
- path: spec/jobTemplate/spec/selector/matchLabels
  group: batch
  kind: CronJob
  behavior: remove
- path: spec/jobTemplate/metadata/labels
  group: batch
  kind: CronJob
  behavior: remove
- path: spec/jobTemplate/spec/template/metadata/labels
  group: batch
  kind: CronJob
  behavior: remove
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cronjob.yaml
---
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: callout
spec:
  schedule: "*/1 * * * *"
  concurrencyPolicy: Forbid
  successfulJobsHistoryLimit: 1
  failedJobsHistoryLimit: 1
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: callout
            image: buildpack-deps:curl
            args:
            - /bin/sh
            - -ec
            - curl http://singleton
          restartPolicy: Never

EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployments.yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: main
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: singleton
  labels:
    special: singleton
spec:
  replicas: 1
  selector:
    matchLabels:
      special: singleton
  template:
    metadata:
      labels:
        special: singleton
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/services.yaml
kind: Service
apiVersion: v1
metadata:
  name: allpods
spec:
  selector:
    app: nginx
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
---
kind: Service
apiVersion: v1
metadata:
  name: singleton
spec:
  selector:
    special: singleton
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1beta1_deployment_main.yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: main
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx:1.7.9
        name: nginx
        ports:
        - containerPort: 80
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1beta1_deployment_singleton.yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  labels:
    app: nginx
    special: singleton
  name: singleton
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
      special: singleton
  template:
    metadata:
      labels:
        app: nginx
        special: singleton
    spec:
      containers:
      - image: nginx:1.7.9
        name: nginx
        ports:
        - containerPort: 80
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/batch_v1beta1_cronjob_callout.yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  labels:
    app: nginx
  name: callout
spec:
  concurrencyPolicy: Forbid
  failedJobsHistoryLimit: 1
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - args:
            - /bin/sh
            - -ec
            - curl http://singleton
            image: buildpack-deps:curl
            name: callout
          restartPolicy: Never
  schedule: '*/1 * * * *'
  successfulJobsHistoryLimit: 1
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_allpods.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: nginx
  name: allpods
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: nginx
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_singleton.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: nginx
  name: singleton
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    special: singleton
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

