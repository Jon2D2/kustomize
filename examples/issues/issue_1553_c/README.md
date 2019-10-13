# Feature Test for Issue 1553


This folder contains files describing how to address [Issue 1553](https://github.com/kubernetes-sigs/kustomize/issues/1553)

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
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/production
mkdir -p ${DEMO_HOME}/overlays/staging
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  app: sentry

resources:
- ./deployment.yaml
- ./values.nodeenv.yaml
- ./values.sentryenv.yaml

secretGenerator:
- name: core
  type: Opaque
  envs:
  - secrets.txt
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base

commonLabels:
  env: $(Values.nodeenv.spec.env)

secretGenerator:
- name: core
  type: Opaque
  envs:
  - secrets.txt
  behavior: replace

patchesStrategicMerge:
- values.nodeenv.yaml
- values.sentryenv.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base

commonLabels:
  env: $(Values.nodeenv.spec.env)

secretGenerator:
- name: core
  type: Opaque
  envs:
  - secrets.txt
  behavior: replace

patchesStrategicMerge:
- values.nodeenv.yaml
- values.sentryenv.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: core
spec:
  template:
    spec:
      containers:
      - name: sentry
        imagePullPolicy: Always
        image: sentryimage:1.0
        readinessProbe: $(Values.sentryenv.spec.readinessProbe)
        env:
        - name: NODE_ENV
          value: $(Values.nodeenv.spec.env)
        - name: SENTRY_ENV
          value: $(Values.sentryenv.spec.env)
        - name: STRIPE_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: $(Secret.core.metadata.name)
              key: stripeSecretKey
        - name: MG_API_KEY
          valueFrom:
            secretKeyRef:
              name: $(Secret.core.metadata.name)
              key: mailgunAPIKey
        - name: PGPASSWORD
          valueFrom:
            secretKeyRef:
              name: $(Secret.core.metadata.name)
              key: postgreSQLPassword
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/values.nodeenv.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: nodeenv
spec:
  env: dev
  args:
    param1: defaultvalue
    param2: defaultvalue
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/values.sentryenv.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: sentryenv
spec:
  env: dev
  readinessProbe:
    exec:
      command:
      - /opt/sentryenv/bin/zkOK.sh
    initialDelaySeconds: 10
    timeoutSeconds: 2
    periodSeconds: 5
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/values.nodeenv.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: nodeenv
spec:
  env: production
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/values.sentryenv.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: sentryenv
spec:
  env: production
  readinessProbe:
    periodSeconds: 22
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/values.nodeenv.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: nodeenv
spec:
  env: staging
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/values.sentryenv.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: sentryenv
spec:
  env: staging
  readinessProbe:
    timeoutSeconds: 7
    periodSeconds: 8
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/secrets.txt
username=userNameDefaultValue
password=passwordDefaultValue
sentryDNS=sentryDNSDefaultValue
stripeSecretKey=stripeSecretKeyDefaultValue
mailgunAPIKey=mailgunAPIKeyDefaultValue
postgreSQLPassword=postgreSQLPasswordDefaultValue
EOF
```


### Preparation Step Other1

<!-- @createOther1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/secrets.txt
username=userNameStagingValue
password=passwordStagingValue
sentryDNS=sentryDNSStagingValue
stripeSecretKey=stripeSecretKeyStagingValue
mailgunAPIKey=mailgunAPIKeyStagingValue
postgreSQLPassword=postgreSQLPasswordStagingValue
EOF
```


### Preparation Step Other2

<!-- @createOther2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/secrets.txt
username=userNameStagingValue
password=passwordStagingValue
sentryDNS=sentryDNSStagingValue
stripeSecretKey=stripeSecretKeyStagingValue
mailgunAPIKey=mailgunAPIKeyStagingValue
postgreSQLPassword=postgreSQLPasswordStagingValue
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/staging -o ${DEMO_HOME}/actual/staging.yaml
kustomize build ${DEMO_HOME}/overlays/production -o ${DEMO_HOME}/actual/production.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production.yaml
apiVersion: v1
data:
  mailgunAPIKey: bWFpbGd1bkFQSUtleVN0YWdpbmdWYWx1ZQ==
  password: cGFzc3dvcmRTdGFnaW5nVmFsdWU=
  postgreSQLPassword: cG9zdGdyZVNRTFBhc3N3b3JkU3RhZ2luZ1ZhbHVl
  sentryDNS: c2VudHJ5RE5TU3RhZ2luZ1ZhbHVl
  stripeSecretKey: c3RyaXBlU2VjcmV0S2V5U3RhZ2luZ1ZhbHVl
  username: dXNlck5hbWVTdGFnaW5nVmFsdWU=
kind: Secret
metadata:
  annotations: {}
  labels:
    app: sentry
    env: production
  name: core-tk6tmthbgm
type: Opaque
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: sentry
    env: production
  name: core
spec:
  selector:
    matchLabels:
      app: sentry
      env: production
  template:
    metadata:
      labels:
        app: sentry
        env: production
    spec:
      containers:
      - env:
        - name: NODE_ENV
          value: production
        - name: SENTRY_ENV
          value: production
        - name: STRIPE_SECRET_KEY
          valueFrom:
            secretKeyRef:
              key: stripeSecretKey
              name: core-tk6tmthbgm
        - name: MG_API_KEY
          valueFrom:
            secretKeyRef:
              key: mailgunAPIKey
              name: core-tk6tmthbgm
        - name: PGPASSWORD
          valueFrom:
            secretKeyRef:
              key: postgreSQLPassword
              name: core-tk6tmthbgm
        image: sentryimage:1.0
        imagePullPolicy: Always
        name: sentry
        readinessProbe:
          exec:
            command:
            - /opt/sentryenv/bin/zkOK.sh
          initialDelaySeconds: 10
          periodSeconds: 22
          timeoutSeconds: 2
---
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  labels:
    app: sentry
    env: production
  name: nodeenv
spec:
  args:
    param1: defaultvalue
    param2: defaultvalue
  env: production
---
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  labels:
    app: sentry
    env: production
  name: sentryenv
spec:
  env: production
  readinessProbe:
    exec:
      command:
      - /opt/sentryenv/bin/zkOK.sh
    initialDelaySeconds: 10
    periodSeconds: 22
    timeoutSeconds: 2
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/staging.yaml
apiVersion: v1
data:
  mailgunAPIKey: bWFpbGd1bkFQSUtleVN0YWdpbmdWYWx1ZQ==
  password: cGFzc3dvcmRTdGFnaW5nVmFsdWU=
  postgreSQLPassword: cG9zdGdyZVNRTFBhc3N3b3JkU3RhZ2luZ1ZhbHVl
  sentryDNS: c2VudHJ5RE5TU3RhZ2luZ1ZhbHVl
  stripeSecretKey: c3RyaXBlU2VjcmV0S2V5U3RhZ2luZ1ZhbHVl
  username: dXNlck5hbWVTdGFnaW5nVmFsdWU=
kind: Secret
metadata:
  annotations: {}
  labels:
    app: sentry
    env: staging
  name: core-tk6tmthbgm
type: Opaque
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: sentry
    env: staging
  name: core
spec:
  selector:
    matchLabels:
      app: sentry
      env: staging
  template:
    metadata:
      labels:
        app: sentry
        env: staging
    spec:
      containers:
      - env:
        - name: NODE_ENV
          value: staging
        - name: SENTRY_ENV
          value: staging
        - name: STRIPE_SECRET_KEY
          valueFrom:
            secretKeyRef:
              key: stripeSecretKey
              name: core-tk6tmthbgm
        - name: MG_API_KEY
          valueFrom:
            secretKeyRef:
              key: mailgunAPIKey
              name: core-tk6tmthbgm
        - name: PGPASSWORD
          valueFrom:
            secretKeyRef:
              key: postgreSQLPassword
              name: core-tk6tmthbgm
        image: sentryimage:1.0
        imagePullPolicy: Always
        name: sentry
        readinessProbe:
          exec:
            command:
            - /opt/sentryenv/bin/zkOK.sh
          initialDelaySeconds: 10
          periodSeconds: 8
          timeoutSeconds: 7
---
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  labels:
    app: sentry
    env: staging
  name: nodeenv
spec:
  args:
    param1: defaultvalue
    param2: defaultvalue
  env: staging
---
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  labels:
    app: sentry
    env: staging
  name: sentryenv
spec:
  env: staging
  readinessProbe:
    exec:
      command:
      - /opt/sentryenv/bin/zkOK.sh
    initialDelaySeconds: 10
    periodSeconds: 8
    timeoutSeconds: 7
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

