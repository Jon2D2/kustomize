# Feature Test for Issue 0952


This folder contains files describing how to address [Issue 0952](https://github.com/kubernetes-sigs/kustomize/issues/0952)

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
- cronjob.yaml
- deployment.yaml
- values.yaml

configurations:
- kustomizeconfig.yaml

vars:
- name : Values.shared.spec.env
  objref:
    apiVersion: v1
    kind: Values
    name: shared
  fieldref:
    fieldpath: spec.env
EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig.yaml
varReference:
- kind: Deployment
  path: spec/template/spec/containers[]/env

- kind: CronJob
  path: spec/jobTemplate/spec/template/spec/containers[]/env
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cronjob.yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: wordpress-cron
  labels:
    app: wordpress
spec:
  schedule: "*/10 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - image: wordpress:4.8-apache
            name: wordpress
            command:
            - php
            args:
            - /path/to/wp-cron.php
            env: $(Values.shared.spec.env)
          restartPolicy: OnFailure
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: wordpress
  labels:
    app: wordpress
spec:
  selector:
    matchLabels:
      app: wordpress
  template:
    metadata:
      labels:
        app: wordpress
    spec:
      containers:
      - image: wordpress:4.8-apache
        name: wordpress
        ports:
        - containerPort: 80
          name: wordpress
        env: $(Values.shared.spec.env)
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: shared
spec:
  env:
  - name: WORDPRESS_DB_USER
    valueFrom:
      secretKeyRef:
        name: wordpress-db-auth
        key: user
  - name: WORDPRESS_DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: wordpress-db-auth
        key: password
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
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1beta2_deployment_wordpress.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: wordpress
  name: wordpress
spec:
  selector:
    matchLabels:
      app: wordpress
  template:
    metadata:
      labels:
        app: wordpress
    spec:
      containers:
      - env:
        - name: WORDPRESS_DB_USER
          valueFrom:
            secretKeyRef:
              key: user
              name: wordpress-db-auth
        - name: WORDPRESS_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: wordpress-db-auth
        image: wordpress:4.8-apache
        name: wordpress
        ports:
        - containerPort: 80
          name: wordpress
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/batch_v1beta1_cronjob_wordpress-cron.yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  labels:
    app: wordpress
  name: wordpress-cron
spec:
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - args:
            - /path/to/wp-cron.php
            command:
            - php
            env:
            - name: WORDPRESS_DB_USER
              valueFrom:
                secretKeyRef:
                  key: user
                  name: wordpress-db-auth
            - name: WORDPRESS_DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  key: password
                  name: wordpress-db-auth
            image: wordpress:4.8-apache
            name: wordpress
          restartPolicy: OnFailure
  schedule: '*/10 * * * *'
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kustomize.config.k8s.io_v1_values_shared.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: shared
spec:
  env:
  - name: WORDPRESS_DB_USER
    valueFrom:
      secretKeyRef:
        key: user
        name: wordpress-db-auth
  - name: WORDPRESS_DB_PASSWORD
    valueFrom:
      secretKeyRef:
        key: password
        name: wordpress-db-auth
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

