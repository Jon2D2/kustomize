# Feature Test for Issue 1343


This folder contains files describing how to address [Issue 1343](https://github.com/kubernetes-sigs/kustomize/issues/1343)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/environment-specific
mkdir -p ${DEMO_HOME}/overlays/environment-specific/patches
mkdir -p ${DEMO_HOME}/overlays/environment-specific/properties
mkdir -p ${DEMO_HOME}/overlays/simple
mkdir -p ${DEMO_HOME}/overlays/instance-specific
mkdir -p ${DEMO_HOME}/templates
mkdir -p ${DEMO_HOME}/templates/mysqld-exporter
mkdir -p ${DEMO_HOME}/templates/mysqld-exporter/resources
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/environment-specific/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

commonLabels:
  app.kubernetes.io/instance: review-myapp-mysqld-exporter
  app.mintel.com/pipeline-stage: review
  app.mintel.com/env: dev

namePrefix: review-

resources:
- ../instance-specific

patchesStrategicMerge:
- patches/remove-cloud-sql-proxy.yaml
- patches/add-vars-to-exporter-main.yaml

configMapGenerator:
- envs:
  - properties/cluster.properties
  name: cluster-properties
  namespace: mynamespace
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/simple/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../../templates/mysqld-exporter

patchesStrategicMerge:
- ./remove-cloud-sql-proxy.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/instance-specific/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

resources:
- ../../templates/mysqld-exporter

commonLabels:
  app.kubernetes.io/managed-by: pipeline
  app.kubernetes.io/name: myapp-mysqld-exporter
  app.kubernetes.io/owner: myowner
  app.kubernetes.io/part-of: myapp
  k8s-app: myapp
  name: myapp-mysqld-exporter

namePrefix: myapp-

namespace: mynamespace

images:
- name: prom/mysqld-exporter
  newTag: v0.11.0
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/templates/mysqld-exporter/kustomization.yaml
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

commonLabels:
  app.kubernetes.io/component: prometheus-exporter

configurations:
- commonlabels.yaml

resources:
- resources/mysqld-exporter-deployment.yaml
- resources/mysqld-exporter-service.yaml
- resources/mysqld-exporter-serviceMonitor.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/environment-specific/patches/add-vars-to-exporter-main.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-mysqld-exporter
  namespace: mynamespace
spec:
  template:
    spec:
      containers:
      - name: mysqld-exporter
        env:
        - name: DATA_SOURCE_NAME
          value: user:password@(host:3306)/
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/environment-specific/patches/remove-cloud-sql-proxy.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-mysqld-exporter
  namespace: mynamespace
spec:
  template:
    spec:
      containers:
      - $patch: delete
        name: cloud-sql-proxy
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/simple/remove-cloud-sql-proxy.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysqld-exporter
spec:
  template:
    spec:
      containers:
      - $patch: delete
        name: cloud-sql-proxy
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/templates/mysqld-exporter/resources/mysqld-exporter-service.yaml
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/templates/mysqld-exporter/resources/mysqld-exporter-deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysqld-exporter
spec:
  template:
    spec:
      containers:
      - name: mysqld-exporter
        image: prom/mysqld-exporter:v0.11.0
        command:
        - /bin/mysqld_exporter
        args:
        - --web.listen-address=:9104
        - --log.level=info
        - --log.format=logger:stderr?json=true
        livenessProbe:
          httpGet:
            path: /
            port: 9104
          initialDelaySeconds: 30
        readinessProbe:
          httpGet:
            path: /
            port: 9104
        env:
        - name: DATA_SOURCE_NAME
          value: $(DB_USER):$(DB_PASSWORD)@(127.0.0.1:3306)/
        ports:
        - name: metrics
          containerPort: 9104
        resources:
          requests:
            cpu: 30m
            memory: 50Mi
          limits:
            cpu: 100m
            memory: 200Mi
        securityContext:
          allowPrivilegeEscalation: false
          runAsUser: 65534
      - name: cloud-sql-proxy
        image: gcr.io/cloudsql-docker/gce-proxy
        imagePullPolicy: IfNotPresent
        command:
        - /cloud_sql_proxy
        args:
        - -instances=$(MASTER_INSTANCE)=tcp:127.0.0.1:3306
        - -credential_file=/secrets/cloudsql/google_credentials
        - -dir=/cloudsql
        ports:
        - name: mysql
          containerPort: 3306
        resources:
          requests:
            cpu: 30m
            memory: 50Mi
          limits:
            cpu: 100m
            memory: 200Mi
        securityContext:
          allowPrivilegeEscalation: false
          runAsUser: 65534
        volumeMounts:
        - name: cloudsql-settings
          mountPath: /secrets/cloudsql
          readOnly: true
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/templates/mysqld-exporter/resources/mysqld-exporter-serviceMonitor.yaml
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/templates/mysqld-exporter/commonlabels.yaml
EOF
```

### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/environment-specific/properties/cluster.properties
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/environment-specific -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/apps_v1_deployment_review-myapp-mysqld-exporter.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/component: prometheus-exporter
    app.kubernetes.io/instance: review-myapp-mysqld-exporter
    app.kubernetes.io/managed-by: pipeline
    app.kubernetes.io/name: myapp-mysqld-exporter
    app.kubernetes.io/owner: myowner
    app.kubernetes.io/part-of: myapp
    app.mintel.com/env: dev
    app.mintel.com/pipeline-stage: review
    k8s-app: myapp
    name: myapp-mysqld-exporter
  name: review-myapp-mysqld-exporter
  namespace: mynamespace
spec:
  selector:
    matchLabels:
      app.kubernetes.io/component: prometheus-exporter
      app.kubernetes.io/instance: review-myapp-mysqld-exporter
      app.kubernetes.io/managed-by: pipeline
      app.kubernetes.io/name: myapp-mysqld-exporter
      app.kubernetes.io/owner: myowner
      app.kubernetes.io/part-of: myapp
      app.mintel.com/env: dev
      app.mintel.com/pipeline-stage: review
      k8s-app: myapp
      name: myapp-mysqld-exporter
  template:
    metadata:
      labels:
        app.kubernetes.io/component: prometheus-exporter
        app.kubernetes.io/instance: review-myapp-mysqld-exporter
        app.kubernetes.io/managed-by: pipeline
        app.kubernetes.io/name: myapp-mysqld-exporter
        app.kubernetes.io/owner: myowner
        app.kubernetes.io/part-of: myapp
        app.mintel.com/env: dev
        app.mintel.com/pipeline-stage: review
        k8s-app: myapp
        name: myapp-mysqld-exporter
    spec:
      containers:
      - args:
        - --web.listen-address=:9104
        - --log.level=info
        - --log.format=logger:stderr?json=true
        command:
        - /bin/mysqld_exporter
        env:
        - name: DATA_SOURCE_NAME
          value: user:password@(host:3306)/
        image: prom/mysqld-exporter:v0.11.0
        livenessProbe:
          httpGet:
            path: /
            port: 9104
          initialDelaySeconds: 30
        name: mysqld-exporter
        ports:
        - containerPort: 9104
          name: metrics
        readinessProbe:
          httpGet:
            path: /
            port: 9104
        resources:
          limits:
            cpu: 100m
            memory: 200Mi
          requests:
            cpu: 30m
            memory: 50Mi
        securityContext:
          allowPrivilegeEscalation: false
          runAsUser: 65534
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_review-cluster-properties-c928988k6m.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/instance: review-myapp-mysqld-exporter
    app.mintel.com/env: dev
    app.mintel.com/pipeline-stage: review
  name: review-cluster-properties-c928988k6m
  namespace: mynamespace
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

