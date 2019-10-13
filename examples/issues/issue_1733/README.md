# Feature Test for Issue 1733


This folder contains files describing how to address [Issue 1733](https://github.com/kubernetes-sigs/kustomize/issues/1733)
Original code is found [here](https://github.com/xmlking/micro-starter-kit/tree/develop/deploy)

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
mkdir -p ${DEMO_HOME}/bases
mkdir -p ${DEMO_HOME}/bases/etcd
mkdir -p ${DEMO_HOME}/bases/micros
mkdir -p ${DEMO_HOME}/bases/micros/account-srv
mkdir -p ${DEMO_HOME}/bases/micros/account-srv/config
mkdir -p ${DEMO_HOME}/bases/micros/emailer-srv
mkdir -p ${DEMO_HOME}/bases/micros/emailer-srv/config
mkdir -p ${DEMO_HOME}/bases/micros/gateway
mkdir -p ${DEMO_HOME}/bases/micros/greeter-srv
mkdir -p ${DEMO_HOME}/bases/micros/greeter-srv/config
mkdir -p ${DEMO_HOME}/bases/micros/proxy
mkdir -p ${DEMO_HOME}/bases/nats
mkdir -p ${DEMO_HOME}/bases/postgres
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/e2e
mkdir -p ${DEMO_HOME}/overlays/e2e/patches
mkdir -p ${DEMO_HOME}/overlays/production
mkdir -p ${DEMO_HOME}/overlays/production/patches
mkdir -p ${DEMO_HOME}/overlays/production/resources
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/etcd/kustomization.yaml
commonLabels:
  app.kubernetes.io/name: etcd
  app.kubernetes.io/instance: etcd-abcxzy
  app.kubernetes.io/component: infra
  app.kubernetes.io/part-of: micro-starter-kit
  app.kubernetes.io/managed-by: kustomize
commonAnnotations:
  org: acmeCorporation

resources:
  - deployment.yaml

vars:
  - name: ETCD_SRV_ENDPOINT
    objref:
      kind: EtcdCluster
      name: etcd-cluster
      apiVersion: etcd.database.coreos.com/v1beta2
    fieldref:
      fieldpath: metadata.name
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/account-srv/kustomization.yaml
commonLabels:
  app.kubernetes.io/name: account-srv
  app.kubernetes.io/instance: account-srv-abcxzy
  app.kubernetes.io/component: microservice

namePrefix: account

resources:
  - deployment.yaml
  - service.yaml

configMapGenerator:
  - name: config
    files:
      - config/config.yaml
  - name: env-vars
    literals:
      - MICRO_SERVER_NAME=accountsrv
      # - MICRO_SERVER_ADVERTISE="$(ACCOUNT_SRV_ENDPOINT):8080"
      - DATABASE_HOST=$(DATABASE_ENDPOINT)

vars:
  - name: ACCOUNT_SRV_ENDPOINT
    objref:
      kind: Service
      name: srv
      apiVersion: v1
    fieldref:
      fieldpath: metadata.name
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/emailer-srv/kustomization.yaml
commonLabels:
  app.kubernetes.io/name: emailer-srv
  app.kubernetes.io/instance: emailer-srv-abcxzy
  app.kubernetes.io/component: microservice

namePrefix: emailer

resources:
  - deployment.yaml
  - service.yaml

configMapGenerator:
  - name: config
    files:
      - config/config.yaml
  - name: env-vars
    literals:
      - MICRO_SERVER_NAME=emailersrv
      # - MICRO_SERVER_ADVERTISE="$(EMAILER_SRV_ENDPOINT):8080"

vars:
  - name: EMAILER_SRV_ENDPOINT
    objref:
      kind: Service
      name: srv
      apiVersion: v1
    fieldref:
      fieldpath: metadata.name
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/gateway/kustomization.yaml
commonLabels:
  app.kubernetes.io/name: gateway
  app.kubernetes.io/instance: gateway-abcxzy
  app.kubernetes.io/component: microservice
  app.kubernetes.io/part-of: micro-starter-kit
  app.kubernetes.io/managed-by: kustomize
commonAnnotations:
  org: acmeCorporation

namePrefix: gateway

resources:
  - deployment.yaml
  - service.yaml

configMapGenerator:
  - name: env-vars
    literals:
      - MICRO_SERVER_NAME=gatewaysrv
      # - MICRO_SERVER_ADVERTISE="$(GATEWAY_SRV_ENDPOINT):8080"
      - MICRO_API_NAMESPACE=""
      - MICRO_API_HANDLER=rpc
      - MICRO_API_ENABLE_RPC="true"
      - MICRO_LOG_LEVEL=debug
      - CORS_ALLOWED_HEADERS="Authorization,Content-Type"
      - CORS_ALLOWED_ORIGINS="*"
      - CORS_ALLOWED_METHODS="POST,GET"

vars:
  - name: GATEWAY_SRV_ENDPOINT
    objref:
      kind: Service
      name: srv
      apiVersion: v1
    fieldref:
      fieldpath: metadata.name
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/greeter-srv/kustomization.yaml
commonLabels:
  app.kubernetes.io/name: greeter-srv
  app.kubernetes.io/instance: greeter-srv-abcxzy
  app.kubernetes.io/component: microservice

namePrefix: greeter

resources:
  - deployment.yaml
  - service.yaml

configMapGenerator:
  - name: config
    files:
      - config/config.yaml
  - name: env-vars
    literals:
      - MICRO_SERVER_NAME=greetersrv
      # - MICRO_SERVER_ADVERTISE="$(GREETER_SRV_ENDPOINT):8080"

vars:
  - name: GREETER_SRV_ENDPOINT
    objref:
      kind: Service
      name: srv
      apiVersion: v1
    fieldref:
      fieldpath: metadata.name
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/kustomization.yaml
namespace: default
commonLabels:
  app.kubernetes.io/part-of: micro-starter-kit
  app.kubernetes.io/managed-by: kustomize
commonAnnotations:
  org: acmeCorporation

bases:
  - gateway
  - proxy
  - account-srv
  - emailer-srv
  - greeter-srv

configurations:
  - kconfig.yaml

configMapGenerator:
  # - name: env-vars
  - name: env-vars-common
    # behavior: merge
    literals:
      - MICRO_SERVER_ADDRESS=0.0.0.0:8080
      - MICRO_BROKER_ADDRESS=0.0.0.0:10001
      - APP_ENV=development
      - CONFIG_DIR=/config
      - CONFIG_FILE=config.yaml
      - MICRO_LOG_LEVEL=debug
      - MICRO_CLIENT_RETRIES=3
      - MICRO_CLIENT_REQUEST_TIMEOUT=5s

secretGenerator:
  - name: secrets
    literals:
      - DATABASE_PASSWORD=fake
  # - name: app-tls
  #   files:
  #     - secret/tls.cert
  #     - secret/tls.key
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/proxy/kustomization.yaml
commonLabels:
  app.kubernetes.io/name: proxy
  app.kubernetes.io/instance: proxy-abcxzy
  app.kubernetes.io/component: microservice
  app.kubernetes.io/part-of: micro-starter-kit
  app.kubernetes.io/managed-by: kustomize
commonAnnotations:
  org: acmeCorporation

namePrefix: proxy

resources:
  - deployment.yaml
  - service.yaml

configMapGenerator:
  - name: env-vars
    literals:
      - MICRO_SERVER_NAME=proxysrv
      # - MICRO_SERVER_ADVERTISE="$(PROXY_SRV_ENDPOINT):8888"
      - MICRO_PROXY_PROTOCOL=grpc
      # - MICRO_PROXY_ADDRESS=0.0.0.0:8081

vars:
  - name: PROXY_SRV_ENDPOINT
    objref:
      kind: Service
      name: srv
      apiVersion: v1
    fieldref:
      fieldpath: metadata.name
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/nats/kustomization.yaml
commonLabels:
  app.kubernetes.io/name: nats
  app.kubernetes.io/instance: nats-abcxzy
  app.kubernetes.io/component: infra
  app.kubernetes.io/part-of: micro-starter-kit
  app.kubernetes.io/managed-by: kustomize
commonAnnotations:
  org: acmeCorporation

resources:
  - nats.yaml

vars:
  - name: NATS_SRV_ENDPOINT
    objref:
      kind: Service
      name: nats
      apiVersion: v1
    fieldref:
      fieldpath: metadata.name
EOF
```


### Preparation Step KustomizationFile8

<!-- @createKustomizationFile8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/postgres/kustomization.yaml
commonLabels:
  app.kubernetes.io/name: postgres
  app.kubernetes.io/instance: postgres-abcxzy
  app.kubernetes.io/component: database

resources:
  - postgres.yaml
  - service.yaml

secretGenerator:
  - name: postgres-secrets
    literals:
      - postgres-password=postgres123

# labels for generated secrets at this level
generatorOptions:
  labels:
    app.kubernetes.io/name: postgres-secrets
    app.kubernetes.io/instance: postgres-secrets-abcxzy
    app.kubernetes.io/component: secrets

vars:
  - name: DATABASE_ENDPOINT
    objref:
      kind: Service
      name: postgres
      apiVersion: v1
    fieldref:
      fieldpath: metadata.name
EOF
```


### Preparation Step KustomizationFile9

<!-- @createKustomizationFile9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/e2e/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: $(NS)

commonLabels:
  environment: integration
  app.kubernetes.io/version: v1
commonAnnotations:
  note: Hello, I am integration!

bases:
  - ../../bases/micros

# enable mage-pull-policy and resource_limit as needed
patches:
  - path: patches/image-pull-policy-if-not-present.yaml
    target:
      kind: Deployment
      labelSelector: app.kubernetes.io/component=microservice

configMapGenerator:
  # - name: env-vars
  - name: env-vars-common
    behavior: merge
    literals:
      - APP_ENV=integration
      - LOG_LEVEL=debug
      - LOG_FORMAT=text
      - MICRO_LOG_LEVEL=debug
      # profile specific variables
      - MICRO_REGISTER_TTL="60"
      - MICRO_REGISTER_INTERVAL="30"
      # - MICRO_SELECTOR=static # static/memory still not working with gateway & proxy
      # - MICRO_REGISTRY=memory
      # following endpoint overwrites (in config.yaml) should be enabled only when  MICRO_SELECTOR=static is used.
      # - ACCOUNTSRV_ENDPOINT=$(ACCOUNT_SRV_ENDPOINT)
      # - GREETERSRV_ENDPOINT=$(GREETER_SRV_ENDPOINT)
      # - EMAILERSRV_ENDPOINT=$(EMAILER_SRV_ENDPOINT)
      # - GATEWAYSRV_ENDPOINT=$(GATEWAY_SRV_ENDPOINT)
      # - PROXYSRV_ENDPOINT=$(PROXY_SRV_ENDPOINT)

secretGenerator:
  - name: secrets
    behavior: replace
    literals:
      - DATABASE_PASSWORD=e2e-real-pass

replicas:
  - name: srv
    count: 1

images:
  - name: etcd
    newTag: 3.3.15
  - name: postgres
    newTag: 11.5-alpine
  - name: micro/micro
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/micro
    newTag: v1.15.1
  - name: xmlking/account-srv
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/account-srv
    newTag: $(IMAGE_VERSION)
  - name: xmlking/emailer-srv
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/emailer-srv
    newTag: $(IMAGE_VERSION)
  - name: xmlking/greeter-srv
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/greeter-srv
    newTag: $(IMAGE_VERSION)
EOF
```


### Preparation Step KustomizationFile10

<!-- @createKustomizationFile10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: $(NS)
namePrefix: prod-
nameSuffix: -v1

commonLabels:
  environment: production
  app.kubernetes.io/version: v1
commonAnnotations:
  note: Hello, I am production!

bases:
  - ../../bases/micros
  - ../../bases/etcd

# enable mage-pull-policy and resource_limit as needed
patches:
  - path: patches/health-sidecar.yaml
    target:
      kind: Deployment
      labelSelector: app.kubernetes.io/component=microservice
  - path: patches/health-sidecar-only-for-proxy.yaml
    target:
      kind: Deployment
      labelSelector: app.kubernetes.io/name=proxy
  # - path: patches/resource_limit.yaml
  #   target:
  #     kind: Deployment
  #     labelSelector: app.kubernetes.io/component=microservice
  - path: patches/image-pull-policy-if-not-present.yaml
    target:
      kind: Deployment
      labelSelector: app.kubernetes.io/component=microservice

configMapGenerator:
  # - name: env-vars
  - name: env-vars-common
    behavior: merge
    literals:
      - APP_ENV=production
      - MICRO_LOG_LEVEL=info
      # profile specific variables
      - MICRO_REGISTRY=etcd
      - MICRO_REGISTRY_ADDRESS="$(ETCD_SRV_ENDPOINT)-client"
      - MICRO_REGISTER_TTL="60"
      - MICRO_REGISTER_INTERVAL="30"

secretGenerator:
  - name: secrets
    behavior: replace
    literals:
      - DATABASE_PASSWORD=prod-real-pass

replicas:
  - name: srv
    count: 1

images:
  - name: etcd
    newTag: 3.3.15
  - name: postgres
    newTag: 11.5-alpine
  - name: micro/micro
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/micro
    newTag: v1.15.1
  - name: xmlking/account-srv
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/account-srv
    newTag: $(IMAGE_VERSION)
  - name: xmlking/emailer-srv
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/emailer-srv
    newTag: $(IMAGE_VERSION)
  - name: xmlking/greeter-srv
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/greeter-srv
    newTag: $(IMAGE_VERSION)
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/etcd/deployment.yaml
apiVersion: "etcd.database.coreos.com/v1beta2"
kind: EtcdCluster
metadata:
  name: etcd-cluster
  ## Adding this annotation make this cluster managed by clusterwide operators
  ## namespaced operators ignore it
  # annotations:
  #   etcd.database.coreos.com/scope: clusterwide
spec:
  size: 3
  version: "3.4.3"
  repository: "quay.io/coreos/etcd"
  pod:
    labels:
      application: micro
    busyboxImage: "busybox:1.28.0-glibc"
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/account-srv/config/config.yaml
environment: production
name: accountsrv
version: v0.1.0
log:
  level: info
  format: json
database:
  dialect: sqlite3
  host: 0.0.0.0
  port: 9999
  Username: sumo
  Password: demo
  database: ":memory:"
  logging: true
observability:
  metrics:
    address: prometheus:8125
    flushInterval: 1000000000
  tracing:
    address: jaeger:6831
    flushInterval: 5000000000
greetersrv:
  endpoint: greetersrv
  version: v0.1.0
emailersrv:
  endpoint: emailersrv
  version: v0.1.0
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/account-srv/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: srv
spec:
  template:
    spec:
      containers:
        - name: srv
          image: xmlking/account-srv:latest
          imagePullPolicy: Always
          ports:
            - name: grpc-port
              containerPort: 8080
              protocol: TCP
          volumeMounts:
            - name: config
              mountPath: /config/config.yaml
              subPath: config.yaml
              readOnly: true
          envFrom:
            - configMapRef:
                name: env-vars
            - configMapRef:
                name: env-vars-common
            - secretRef:
                name: secrets
      volumes:
        - name: config
          configMap:
            name: config
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/account-srv/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: srv
spec:
  type: ClusterIP
  ports:
    - name: grpc-port
      port: 8080
      protocol: TCP
      targetPort: grpc-port
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/emailer-srv/config/config.yaml
environment: production
name: emailersrv
version: v0.1.0
log:
  level: info
  format: json
observability:
  metrics:
    address: prometheus:8125
    flushInterval: 1000000000
  tracing:
    address: jaeger:6831
    flushInterval: 5000000000
email:
  username: yourGmailUsername
  password: yourGmailAppPassword
  emailServer: smtp.gmail.com
  port: 587
  from: xmlking@gmail.com
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/emailer-srv/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: srv
spec:
  template:
    spec:
      containers:
        - name: srv
          image: xmlking/emailer-srv:latest
          imagePullPolicy: Always
          ports:
            - name: grpc-port
              containerPort: 8080
              protocol: TCP
          volumeMounts:
            - name: config
              mountPath: /config/config.yaml
              subPath: config.yaml
              readOnly: true
          envFrom:
            - configMapRef:
                name: env-vars
            - configMapRef:
                name: env-vars-common
            - secretRef:
                name: secrets
      volumes:
        - name: config
          configMap:
            name: config
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/emailer-srv/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: srv
spec:
  type: ClusterIP
  ports:
    - name: grpc-port
      port: 8080
      protocol: TCP
      targetPort: grpc-port
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/gateway/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: srv
spec:
  template:
    # metadata:
    #   annotations:
    #     sidecar.istio.io/inject: "false"
    spec:
      containers:
        - name: srv
          image: micro/micro:latest
          imagePullPolicy: Always
          args:
            - "api"
            - "--handler=rpc"
            - "--enable_rpc=true"
            - "--address=0.0.0.0:8090"
          ports:
            - name: http-gateway
              containerPort: 8090
              protocol: TCP
          envFrom:
            - configMapRef:
                name: env-vars
            - configMapRef:
                name: env-vars-common
          env:
            - name: MICRO_API_ADDRESS
              value: 0.0.0.0:8090
EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/gateway/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: srv
spec:
  type: LoadBalancer
  ports:
    - name: http-gateway
      port: 8080
      protocol: TCP
      targetPort: http-gateway
EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/greeter-srv/config/config.yaml
environment: production
name: greetersrv
version: v0.1.0
log:
  level: info
  format: json
observability:
  metrics:
    address: prometheus:8125
    flushInterval: 1000000000
  tracing:
    address: jaeger:6831
    flushInterval: 5000000000
EOF
```


### Preparation Step Resource10

<!-- @createResource10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/greeter-srv/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: srv
spec:
  template:
    spec:
      containers:
        - name: srv
          image: xmlking/greeter-srv:latest
          imagePullPolicy: Always
          ports:
            - name: grpc-port
              containerPort: 8080
              protocol: TCP
          volumeMounts:
            - name: config
              mountPath: /config/config.yaml
              subPath: config.yaml
              readOnly: true
          envFrom:
            - configMapRef:
                name: env-vars
            - configMapRef:
                name: env-vars-common
            - secretRef:
                name: secrets
      volumes:
        - name: config
          configMap:
            name: config
EOF
```


### Preparation Step Resource11

<!-- @createResource11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/greeter-srv/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: srv
spec:
  type: ClusterIP
  ports:
    - name: grpc-port
      port: 8080
      protocol: TCP
      targetPort: grpc-port
EOF
```


### Preparation Step Resource12

<!-- @createResource12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/kconfig.yaml
varReference:
  - path: data
    kind: ConfigMap
EOF
```


### Preparation Step Resource13

<!-- @createResource13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/proxy/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: srv
spec:
  template:
    # metadata:
    #   annotations:
    #     sidecar.istio.io/inject: "false"
    spec:
      containers:
        - name: srv
          image: micro/micro:latest
          imagePullPolicy: Always
          args:
            - "proxy"
          ports:
            - name: grpc-proxy
              containerPort: 8081
              protocol: TCP
          envFrom:
            - configMapRef:
                name: env-vars
            - configMapRef:
                name: env-vars-common
EOF
```


### Preparation Step Resource14

<!-- @createResource14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/micros/proxy/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: srv
spec:
  type: LoadBalancer
  ports:
    - name: grpc-proxy
      port: 8888
      protocol: TCP
      targetPort: grpc-proxy
EOF
```


### Preparation Step Resource15

<!-- @createResource15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/nats/nats.yaml
## Install NATS Operator
#
# kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/00-prereqs.yaml
# kubectl apply -f https://github.com/nats-io/nats-operator/releases/latest/download/10-deployment.yaml
#
apiVersion: nats.io/v1alpha2
kind: NatsCluster
metadata:
  name: nats-cluster
spec:
  size: 3
  version: "2.1.0"
EOF
```


### Preparation Step Resource16

<!-- @createResource16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/postgres/postgres.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  labels:
    app: postgres
spec:
  replicas: 1
  serviceName: postgres-internal
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      terminationGracePeriodSeconds: 0
      containers:
        - name: postgres
          image: postgres:11.5-alpine
          imagePullPolicy: Always
          ports:
            - name: tcp-pg
              containerPort: 5432
              protocol: TCP
          env:
            - name: POSTGRES_DB
              value: postgres
            - name: POSTGRES_USER
              value: postgres
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secrets
                  key: postgres-password
          livenessProbe:
            exec:
              command: ["pg_isready", "-U", "$(POSTGRES_USER)"]
            initialDelaySeconds: 3
            timeoutSeconds: 2
          readinessProbe:
            exec:
              command: ["pg_isready", "-U", "$(POSTGRES_USER)"]
            initialDelaySeconds: 3
            timeoutSeconds: 2
          volumeMounts:
            - name: database-storage
              mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
    - metadata:
        name: postgres-storage
        labels:
          app: postgres
      spec:
        accessModes: ["ReadWriteOnce"]
        # storageClassName: <custom storage class>
        resources:
          requests:
            storage: 1Gi
EOF
```


### Preparation Step Resource17

<!-- @createResource17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/postgres/service-headless.yaml
---
apiVersion: v1
kind: Service
metadata:
  name: postgres-headless
  labels:
    app: postgres
spec:
  type: ClusterIP
  clusterIP: None
  ports:
    - name: tcp-pg
      port: 5432
      targetPort: tcp-pg
  selector:
    app: postgres
EOF
```


### Preparation Step Resource18

<!-- @createResource18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/postgres/service.yaml
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  labels:
    app: postgres
spec:
  type: NodePort
  ports:
    - name: tcp-pg
      port: 5432
      targetPort: tcp-pg
      nodePort: 31432
  selector:
    app: postgres
    role: master
EOF
```


### Preparation Step Resource19

<!-- @createResource19 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/bases/postgres/statefulset.yaml
---
apiVersion: apps/v1beta2
kind: StatefulSet
metadata:
  name: postgres
  labels:
    app: postgres
spec:
  replicas: 1
  serviceName: postgres-headless
  updateStrategy:
    type: RollingUpdate
  selector:
    matchLabels:
      app: postgres
      role: master
  template:
    metadata:
      name: postgres
      labels:
        app: postgres
        role: master
    spec:
      securityContext:
        fsGroup: 1001
      initContainers:
        - name: init-chmod-data
          image: docker.io/bitnami/minideb:latest
          imagePullPolicy: "Always"
          resources:
            requests:
              cpu: 250m
              memory: 256Mi

          command:
            - sh
            - -c
            - |
              mkdir -p /bitnami/postgresql/data
              chmod 700 /bitnami/postgresql/data
              find /bitnami/postgresql -mindepth 1 -maxdepth 1 -not -name ".snapshot" -not -name "lost+found" | \
                xargs chown -R 1001:1001
          securityContext:
            runAsUser: 0
          volumeMounts:
            - name: data
              mountPath: /bitnami/postgresql
              subPath:
      containers:
        - name: postgres
          image: docker.io/bitnami/postgresql:10.7.0
          imagePullPolicy: "IfNotPresent"
          resources:
            requests:
              cpu: 250m
              memory: 256Mi

          securityContext:
            runAsUser: 1001
          env:
            - name: BITNAMI_DEBUG
              value: "false"
            - name: POSTGRESQL_PORT_NUMBER
              value: "5432"
            - name: POSTGRESQL_VOLUME_DIR
              value: "/bitnami/postgresql"
            - name: PGDATA
              value: "/bitnami/postgresql/data"
            - name: POSTGRES_USER
              value: "postgres"
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres
                  key: postgres-password
          ports:
            - name: postgres
              containerPort: 5432
          livenessProbe:
            exec:
              command:
                - sh
                - -c
                - exec pg_isready -U "postgres" -h 127.0.0.1 -p 5432
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            successThreshold: 1
            failureThreshold: 6
          readinessProbe:
            exec:
              command:
                - sh
                - -c
                - |
                  pg_isready -U "postgres" -h 127.0.0.1 -p 5432
                  [ -f /opt/bitnami/postgresql/tmp/.initialized ]
            initialDelaySeconds: 5
            periodSeconds: 10
            timeoutSeconds: 5
            successThreshold: 1
            failureThreshold: 6
          volumeMounts:
            - name: data
              mountPath: /bitnami/postgresql
              subPath:
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes:
          - "ReadWriteOnce"
        resources:
          requests:
            storage: "8Gi"
EOF
```


### Preparation Step Resource20

<!-- @createResource20 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/e2e/kustomization-static.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: $(NS)

commonLabels:
  environment: integration
  app.kubernetes.io/version: v1
commonAnnotations:
  note: Hello, I am integration!

bases:
  - ../../bases/micros

# enable mage-pull-policy and resource_limit as needed
patches:
  - path: patches/image-pull-policy-if-not-present.yaml
    target:
      kind: Deployment
      labelSelector: app.kubernetes.io/component=microservice

configMapGenerator:
  # - name: env-vars
  - name: env-vars-common
    behavior: merge
    literals:
      - APP_ENV=integration
      - LOG_LEVEL=debug
      - LOG_FORMAT=text
      - MICRO_LOG_LEVEL=debug
      # profile specific variables
      - MICRO_REGISTER_TTL="60"
      - MICRO_REGISTER_INTERVAL="30"
      # static/memory still not working with gateway & proxy. publish to emailer not working too...
      - MICRO_SELECTOR=static
      - MICRO_REGISTRY=memory
      # following endpoint overwrites (in config.yaml) should be enabled only when  MICRO_SELECTOR=static is used.
      - ACCOUNTSRV_ENDPOINT=$(ACCOUNT_SRV_ENDPOINT)
      - GREETERSRV_ENDPOINT=$(GREETER_SRV_ENDPOINT)
      - EMAILERSRV_ENDPOINT=$(EMAILER_SRV_ENDPOINT)
      - GATEWAYSRV_ENDPOINT=$(GATEWAY_SRV_ENDPOINT)
      - PROXYSRV_ENDPOINT=$(PROXY_SRV_ENDPOINT)

secretGenerator:
  - name: secrets
    behavior: replace
    literals:
      - DATABASE_PASSWORD=e2e-real-pass

replicas:
  - name: srv
    count: 1

images:
  - name: etcd
    newTag: 3.3.15
  - name: postgres
    newTag: 11.5-alpine
  - name: micro/micro
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/micro
    newTag: v1.15.1
  - name: xmlking/account-srv
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/account-srv
    newTag: $(IMAGE_VERSION)
  - name: xmlking/emailer-srv
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/emailer-srv
    newTag: $(IMAGE_VERSION)
  - name: xmlking/greeter-srv
    newName: docker.pkg.github.com/xmlking/micro-starter-kit/greeter-srv
    newTag: $(IMAGE_VERSION)
EOF
```


### Preparation Step Resource21

<!-- @createResource21 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/e2e/patches/image-pull-policy-if-not-present.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: srv
spec:
  template:
    spec:
      containers:
        - name: srv
          imagePullPolicy: IfNotPresent
EOF
```


### Preparation Step Resource22

<!-- @createResource22 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/patches/health-sidecar-only-for-proxy.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: health-sidecar
spec:
  template:
    spec:
      containers:
        - name: health
          image: micro/micro:latest
          args:
            - "health"
            - "--address=:8088"
            - "--check_service=$(MICRO_SERVER_NAME)"
            - "--check_address=0.0.0.0:8081"
          envFrom:
            - configMapRef:
                name: env-vars
          livenessProbe:
            httpGet:
              path: /health
              port: 8088
            initialDelaySeconds: 30
            periodSeconds: 10
EOF
```


### Preparation Step Resource23

<!-- @createResource23 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/patches/health-sidecar.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: health-sidecar
spec:
  template:
    spec:
      containers:
        - name: health
          image: micro/micro:latest
          args:
            - "health"
            - "--address=:8088"
            - "--check_service=$(MICRO_SERVER_NAME)"
            - "--check_address=0.0.0.0:8080"
          envFrom:
            - configMapRef:
                name: env-vars
          livenessProbe:
            httpGet:
              path: /health
              port: 8088
            initialDelaySeconds: 30
            periodSeconds: 10
EOF
```


### Preparation Step Resource24

<!-- @createResource24 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/patches/image-pull-policy-if-not-present.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: srv
spec:
  template:
    spec:
      containers:
        - name: srv
          imagePullPolicy: IfNotPresent
EOF
```


### Preparation Step Resource25

<!-- @createResource25 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/patches/resource_limit.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: srv
spec:
  template:
    spec:
      containers:
        - name: srv
          resources:
            requests:
              cpu: 250m
              memory: 1G
            limits:
              cpu: 500m
              memory: 2G
EOF
```


### Preparation Step Resource26

<!-- @createResource26 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/resources/hpa.yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: account-srv-hpa
spec:
  minReplicas: 1
  maxReplicas: 5
  metrics:
    - resource:
        name: cpu
        target:
          averageUtilization: 70
          type: Utilization
      type: Resource
    - resource:
        name: memory
        target:
          averageUtilization: 70
          type: Utilization
      type: Resource
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: account-srv
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/e2e -o ${DEMO_HOME}/actual/e2e.yaml
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
cat <<'EOF' >${DEMO_HOME}/expected/e2e.yaml
apiVersion: v1
data:
  config.yaml: |
    environment: production
    name: accountsrv
    version: v0.1.0
    log:
      level: info
      format: json
    database:
      dialect: sqlite3
      host: 0.0.0.0
      port: 9999
      Username: sumo
      Password: demo
      database: ":memory:"
      logging: true
    observability:
      metrics:
        address: prometheus:8125
        flushInterval: 1000000000
      tracing:
        address: jaeger:6831
        flushInterval: 5000000000
    greetersrv:
      endpoint: greetersrv
      version: v0.1.0
    emailersrv:
      endpoint: emailersrv
      version: v0.1.0
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: account-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: account-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: accountconfig-6cccg5m4d7
  namespace: $(NS)
---
apiVersion: v1
data:
  DATABASE_HOST: $(DATABASE_ENDPOINT)
  MICRO_SERVER_NAME: accountsrv
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: account-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: account-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: accountenv-vars-f8dmc5gt6c
  namespace: $(NS)
---
apiVersion: v1
data:
  config.yaml: |
    environment: production
    name: emailersrv
    version: v0.1.0
    log:
      level: info
      format: json
    observability:
      metrics:
        address: prometheus:8125
        flushInterval: 1000000000
      tracing:
        address: jaeger:6831
        flushInterval: 5000000000
    email:
      username: yourGmailUsername
      password: yourGmailAppPassword
      emailServer: smtp.gmail.com
      port: 587
      from: xmlking@gmail.com
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: emailer-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: emailer-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: emailerconfig-2dhkfff2b2
  namespace: $(NS)
---
apiVersion: v1
data:
  MICRO_SERVER_NAME: emailersrv
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: emailer-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: emailer-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: emailerenv-vars-ch9m5b9577
  namespace: $(NS)
---
apiVersion: v1
data:
  APP_ENV: integration
  CONFIG_DIR: /config
  CONFIG_FILE: config.yaml
  LOG_FORMAT: text
  LOG_LEVEL: debug
  MICRO_BROKER_ADDRESS: 0.0.0.0:10001
  MICRO_CLIENT_REQUEST_TIMEOUT: 5s
  MICRO_CLIENT_RETRIES: "3"
  MICRO_LOG_LEVEL: debug
  MICRO_REGISTER_INTERVAL: "30"
  MICRO_REGISTER_TTL: "60"
  MICRO_SERVER_ADDRESS: 0.0.0.0:8080
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: env-vars-common-bgfk58fk86
  namespace: $(NS)
---
apiVersion: v1
data:
  CORS_ALLOWED_HEADERS: Authorization,Content-Type
  CORS_ALLOWED_METHODS: POST,GET
  CORS_ALLOWED_ORIGINS: '*'
  MICRO_API_ENABLE_RPC: "true"
  MICRO_API_HANDLER: rpc
  MICRO_API_NAMESPACE: ""
  MICRO_LOG_LEVEL: debug
  MICRO_SERVER_NAME: gatewaysrv
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: gateway-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: gateway
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: gatewayenv-vars-th5fd48962
  namespace: $(NS)
---
apiVersion: v1
data:
  config.yaml: |
    environment: production
    name: greetersrv
    version: v0.1.0
    log:
      level: info
      format: json
    observability:
      metrics:
        address: prometheus:8125
        flushInterval: 1000000000
      tracing:
        address: jaeger:6831
        flushInterval: 5000000000
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: greeter-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: greeter-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: greeterconfig-4h67chd28f
  namespace: $(NS)
---
apiVersion: v1
data:
  MICRO_SERVER_NAME: greetersrv
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: greeter-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: greeter-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: greeterenv-vars-6h9mfk5bk6
  namespace: $(NS)
---
apiVersion: v1
data:
  MICRO_PROXY_PROTOCOL: grpc
  MICRO_SERVER_NAME: proxysrv
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: proxy-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: proxy
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: proxyenv-vars-t88hd4g4f4
  namespace: $(NS)
---
apiVersion: v1
data:
  DATABASE_PASSWORD: ZTJlLXJlYWwtcGFzcw==
kind: Secret
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: secrets-tbd677b8ff
  namespace: $(NS)
type: Opaque
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: account-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: account-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: accountsrv
  namespace: $(NS)
spec:
  ports:
  - name: grpc-port
    port: 8080
    protocol: TCP
    targetPort: grpc-port
  selector:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: account-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: account-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: emailer-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: emailer-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: emailersrv
  namespace: $(NS)
spec:
  ports:
  - name: grpc-port
    port: 8080
    protocol: TCP
    targetPort: grpc-port
  selector:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: emailer-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: emailer-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: gateway-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: gateway
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: gatewaysrv
  namespace: $(NS)
spec:
  ports:
  - name: http-gateway
    port: 8080
    protocol: TCP
    targetPort: http-gateway
  selector:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: gateway-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: gateway
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  type: LoadBalancer
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: greeter-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: greeter-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: greetersrv
  namespace: $(NS)
spec:
  ports:
  - name: grpc-port
    port: 8080
    protocol: TCP
    targetPort: grpc-port
  selector:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: greeter-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: greeter-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: proxy-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: proxy
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: proxysrv
  namespace: $(NS)
spec:
  ports:
  - name: grpc-proxy
    port: 8888
    protocol: TCP
    targetPort: grpc-proxy
  selector:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: proxy-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: proxy
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: account-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: account-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: accountsrv
  namespace: $(NS)
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/component: microservice
      app.kubernetes.io/instance: account-srv-abcxzy
      app.kubernetes.io/managed-by: kustomize
      app.kubernetes.io/name: account-srv
      app.kubernetes.io/part-of: micro-starter-kit
      app.kubernetes.io/version: v1
      environment: integration
  template:
    metadata:
      annotations:
        note: Hello, I am integration!
        org: acmeCorporation
      labels:
        app.kubernetes.io/component: microservice
        app.kubernetes.io/instance: account-srv-abcxzy
        app.kubernetes.io/managed-by: kustomize
        app.kubernetes.io/name: account-srv
        app.kubernetes.io/part-of: micro-starter-kit
        app.kubernetes.io/version: v1
        environment: integration
    spec:
      containers:
      - envFrom:
        - configMapRef:
            name: accountenv-vars-f8dmc5gt6c
        - configMapRef:
            name: env-vars-common-bgfk58fk86
        - secretRef:
            name: secrets-tbd677b8ff
        image: docker.pkg.github.com/xmlking/micro-starter-kit/account-srv:$(IMAGE_VERSION)
        imagePullPolicy: IfNotPresent
        name: srv
        ports:
        - containerPort: 8080
          name: grpc-port
          protocol: TCP
        volumeMounts:
        - mountPath: /config/config.yaml
          name: config
          readOnly: true
          subPath: config.yaml
      volumes:
      - configMap:
          name: accountconfig-6cccg5m4d7
        name: config
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: emailer-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: emailer-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: emailersrv
  namespace: $(NS)
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/component: microservice
      app.kubernetes.io/instance: emailer-srv-abcxzy
      app.kubernetes.io/managed-by: kustomize
      app.kubernetes.io/name: emailer-srv
      app.kubernetes.io/part-of: micro-starter-kit
      app.kubernetes.io/version: v1
      environment: integration
  template:
    metadata:
      annotations:
        note: Hello, I am integration!
        org: acmeCorporation
      labels:
        app.kubernetes.io/component: microservice
        app.kubernetes.io/instance: emailer-srv-abcxzy
        app.kubernetes.io/managed-by: kustomize
        app.kubernetes.io/name: emailer-srv
        app.kubernetes.io/part-of: micro-starter-kit
        app.kubernetes.io/version: v1
        environment: integration
    spec:
      containers:
      - envFrom:
        - configMapRef:
            name: emailerenv-vars-ch9m5b9577
        - configMapRef:
            name: env-vars-common-bgfk58fk86
        - secretRef:
            name: secrets-tbd677b8ff
        image: docker.pkg.github.com/xmlking/micro-starter-kit/emailer-srv:$(IMAGE_VERSION)
        imagePullPolicy: IfNotPresent
        name: srv
        ports:
        - containerPort: 8080
          name: grpc-port
          protocol: TCP
        volumeMounts:
        - mountPath: /config/config.yaml
          name: config
          readOnly: true
          subPath: config.yaml
      volumes:
      - configMap:
          name: emailerconfig-2dhkfff2b2
        name: config
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: gateway-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: gateway
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: gatewaysrv
  namespace: $(NS)
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/component: microservice
      app.kubernetes.io/instance: gateway-abcxzy
      app.kubernetes.io/managed-by: kustomize
      app.kubernetes.io/name: gateway
      app.kubernetes.io/part-of: micro-starter-kit
      app.kubernetes.io/version: v1
      environment: integration
  template:
    metadata:
      annotations:
        note: Hello, I am integration!
        org: acmeCorporation
      labels:
        app.kubernetes.io/component: microservice
        app.kubernetes.io/instance: gateway-abcxzy
        app.kubernetes.io/managed-by: kustomize
        app.kubernetes.io/name: gateway
        app.kubernetes.io/part-of: micro-starter-kit
        app.kubernetes.io/version: v1
        environment: integration
    spec:
      containers:
      - args:
        - api
        - --handler=rpc
        - --enable_rpc=true
        - --address=0.0.0.0:8090
        env:
        - name: MICRO_API_ADDRESS
          value: 0.0.0.0:8090
        envFrom:
        - configMapRef:
            name: gatewayenv-vars-th5fd48962
        - configMapRef:
            name: env-vars-common-bgfk58fk86
        image: docker.pkg.github.com/xmlking/micro-starter-kit/micro:v1.15.1
        imagePullPolicy: IfNotPresent
        name: srv
        ports:
        - containerPort: 8090
          name: http-gateway
          protocol: TCP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: greeter-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: greeter-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: greetersrv
  namespace: $(NS)
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/component: microservice
      app.kubernetes.io/instance: greeter-srv-abcxzy
      app.kubernetes.io/managed-by: kustomize
      app.kubernetes.io/name: greeter-srv
      app.kubernetes.io/part-of: micro-starter-kit
      app.kubernetes.io/version: v1
      environment: integration
  template:
    metadata:
      annotations:
        note: Hello, I am integration!
        org: acmeCorporation
      labels:
        app.kubernetes.io/component: microservice
        app.kubernetes.io/instance: greeter-srv-abcxzy
        app.kubernetes.io/managed-by: kustomize
        app.kubernetes.io/name: greeter-srv
        app.kubernetes.io/part-of: micro-starter-kit
        app.kubernetes.io/version: v1
        environment: integration
    spec:
      containers:
      - envFrom:
        - configMapRef:
            name: greeterenv-vars-6h9mfk5bk6
        - configMapRef:
            name: env-vars-common-bgfk58fk86
        - secretRef:
            name: secrets-tbd677b8ff
        image: docker.pkg.github.com/xmlking/micro-starter-kit/greeter-srv:$(IMAGE_VERSION)
        imagePullPolicy: IfNotPresent
        name: srv
        ports:
        - containerPort: 8080
          name: grpc-port
          protocol: TCP
        volumeMounts:
        - mountPath: /config/config.yaml
          name: config
          readOnly: true
          subPath: config.yaml
      volumes:
      - configMap:
          name: greeterconfig-4h67chd28f
        name: config
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am integration!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: proxy-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: proxy
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: integration
  name: proxysrv
  namespace: $(NS)
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/component: microservice
      app.kubernetes.io/instance: proxy-abcxzy
      app.kubernetes.io/managed-by: kustomize
      app.kubernetes.io/name: proxy
      app.kubernetes.io/part-of: micro-starter-kit
      app.kubernetes.io/version: v1
      environment: integration
  template:
    metadata:
      annotations:
        note: Hello, I am integration!
        org: acmeCorporation
      labels:
        app.kubernetes.io/component: microservice
        app.kubernetes.io/instance: proxy-abcxzy
        app.kubernetes.io/managed-by: kustomize
        app.kubernetes.io/name: proxy
        app.kubernetes.io/part-of: micro-starter-kit
        app.kubernetes.io/version: v1
        environment: integration
    spec:
      containers:
      - args:
        - proxy
        envFrom:
        - configMapRef:
            name: proxyenv-vars-t88hd4g4f4
        - configMapRef:
            name: env-vars-common-bgfk58fk86
        image: docker.pkg.github.com/xmlking/micro-starter-kit/micro:v1.15.1
        imagePullPolicy: IfNotPresent
        name: srv
        ports:
        - containerPort: 8081
          name: grpc-proxy
          protocol: TCP
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production.yaml
apiVersion: v1
data:
  config.yaml: |
    environment: production
    name: accountsrv
    version: v0.1.0
    log:
      level: info
      format: json
    database:
      dialect: sqlite3
      host: 0.0.0.0
      port: 9999
      Username: sumo
      Password: demo
      database: ":memory:"
      logging: true
    observability:
      metrics:
        address: prometheus:8125
        flushInterval: 1000000000
      tracing:
        address: jaeger:6831
        flushInterval: 5000000000
    greetersrv:
      endpoint: greetersrv
      version: v0.1.0
    emailersrv:
      endpoint: emailersrv
      version: v0.1.0
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: account-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: account-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-accountconfig-v1-kkgm2mdk6d
  namespace: $(NS)
---
apiVersion: v1
data:
  DATABASE_HOST: $(DATABASE_ENDPOINT)
  MICRO_SERVER_NAME: accountsrv
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: account-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: account-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-accountenv-vars-v1-bkm8c88cth
  namespace: $(NS)
---
apiVersion: v1
data:
  config.yaml: |
    environment: production
    name: emailersrv
    version: v0.1.0
    log:
      level: info
      format: json
    observability:
      metrics:
        address: prometheus:8125
        flushInterval: 1000000000
      tracing:
        address: jaeger:6831
        flushInterval: 5000000000
    email:
      username: yourGmailUsername
      password: yourGmailAppPassword
      emailServer: smtp.gmail.com
      port: 587
      from: xmlking@gmail.com
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: emailer-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: emailer-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-emailerconfig-v1-9bm4dcdm69
  namespace: $(NS)
---
apiVersion: v1
data:
  MICRO_SERVER_NAME: emailersrv
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: emailer-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: emailer-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-emailerenv-vars-v1-7785t5275f
  namespace: $(NS)
---
apiVersion: v1
data:
  APP_ENV: production
  CONFIG_DIR: /config
  CONFIG_FILE: config.yaml
  MICRO_BROKER_ADDRESS: 0.0.0.0:10001
  MICRO_CLIENT_REQUEST_TIMEOUT: 5s
  MICRO_CLIENT_RETRIES: "3"
  MICRO_LOG_LEVEL: info
  MICRO_REGISTER_INTERVAL: "30"
  MICRO_REGISTER_TTL: "60"
  MICRO_REGISTRY: etcd
  MICRO_REGISTRY_ADDRESS: prod-etcd-cluster-v1-client
  MICRO_SERVER_ADDRESS: 0.0.0.0:8080
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-env-vars-common-v1-t4df5m59fk
  namespace: $(NS)
---
apiVersion: v1
data:
  CORS_ALLOWED_HEADERS: Authorization,Content-Type
  CORS_ALLOWED_METHODS: POST,GET
  CORS_ALLOWED_ORIGINS: '*'
  MICRO_API_ENABLE_RPC: "true"
  MICRO_API_HANDLER: rpc
  MICRO_API_NAMESPACE: ""
  MICRO_LOG_LEVEL: debug
  MICRO_SERVER_NAME: gatewaysrv
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: gateway-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: gateway
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-gatewayenv-vars-v1-ft7gft5265
  namespace: $(NS)
---
apiVersion: v1
data:
  config.yaml: |
    environment: production
    name: greetersrv
    version: v0.1.0
    log:
      level: info
      format: json
    observability:
      metrics:
        address: prometheus:8125
        flushInterval: 1000000000
      tracing:
        address: jaeger:6831
        flushInterval: 5000000000
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: greeter-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: greeter-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-greeterconfig-v1-9khmk64kh6
  namespace: $(NS)
---
apiVersion: v1
data:
  MICRO_SERVER_NAME: greetersrv
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: greeter-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: greeter-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-greeterenv-vars-v1-b852cg88tg
  namespace: $(NS)
---
apiVersion: v1
data:
  MICRO_PROXY_PROTOCOL: grpc
  MICRO_SERVER_NAME: proxysrv
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: proxy-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: proxy
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-proxyenv-vars-v1-6f5k5tk78d
  namespace: $(NS)
---
apiVersion: v1
data:
  DATABASE_PASSWORD: cHJvZC1yZWFsLXBhc3M=
kind: Secret
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-secrets-v1-t45t559g92
  namespace: $(NS)
type: Opaque
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: account-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: account-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-accountsrv-v1
  namespace: $(NS)
spec:
  ports:
  - name: grpc-port
    port: 8080
    protocol: TCP
    targetPort: grpc-port
  selector:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: account-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: account-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: emailer-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: emailer-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-emailersrv-v1
  namespace: $(NS)
spec:
  ports:
  - name: grpc-port
    port: 8080
    protocol: TCP
    targetPort: grpc-port
  selector:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: emailer-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: emailer-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: gateway-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: gateway
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-gatewaysrv-v1
  namespace: $(NS)
spec:
  ports:
  - name: http-gateway
    port: 8080
    protocol: TCP
    targetPort: http-gateway
  selector:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: gateway-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: gateway
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  type: LoadBalancer
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: greeter-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: greeter-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-greetersrv-v1
  namespace: $(NS)
spec:
  ports:
  - name: grpc-port
    port: 8080
    protocol: TCP
    targetPort: grpc-port
  selector:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: greeter-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: greeter-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: proxy-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: proxy
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-proxysrv-v1
  namespace: $(NS)
spec:
  ports:
  - name: grpc-proxy
    port: 8888
    protocol: TCP
    targetPort: grpc-proxy
  selector:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: proxy-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: proxy
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: account-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: account-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-accountsrv-v1
  namespace: $(NS)
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/component: microservice
      app.kubernetes.io/instance: account-srv-abcxzy
      app.kubernetes.io/managed-by: kustomize
      app.kubernetes.io/name: account-srv
      app.kubernetes.io/part-of: micro-starter-kit
      app.kubernetes.io/version: v1
      environment: production
  template:
    metadata:
      annotations:
        note: Hello, I am production!
        org: acmeCorporation
      labels:
        app.kubernetes.io/component: microservice
        app.kubernetes.io/instance: account-srv-abcxzy
        app.kubernetes.io/managed-by: kustomize
        app.kubernetes.io/name: account-srv
        app.kubernetes.io/part-of: micro-starter-kit
        app.kubernetes.io/version: v1
        environment: production
    spec:
      containers:
      - args:
        - health
        - --address=:8088
        - --check_service=$(MICRO_SERVER_NAME)
        - --check_address=0.0.0.0:8080
        envFrom:
        - configMapRef:
            name: prod-accountenv-vars-v1-bkm8c88cth
        image: docker.pkg.github.com/xmlking/micro-starter-kit/micro:v1.15.1
        livenessProbe:
          httpGet:
            path: /health
            port: 8088
          initialDelaySeconds: 30
          periodSeconds: 10
        name: health
      - envFrom:
        - configMapRef:
            name: prod-accountenv-vars-v1-bkm8c88cth
        - configMapRef:
            name: prod-env-vars-common-v1-t4df5m59fk
        - secretRef:
            name: prod-secrets-v1-t45t559g92
        image: docker.pkg.github.com/xmlking/micro-starter-kit/account-srv:$(IMAGE_VERSION)
        imagePullPolicy: IfNotPresent
        name: srv
        ports:
        - containerPort: 8080
          name: grpc-port
          protocol: TCP
        volumeMounts:
        - mountPath: /config/config.yaml
          name: config
          readOnly: true
          subPath: config.yaml
      volumes:
      - configMap:
          name: prod-accountconfig-v1-kkgm2mdk6d
        name: config
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: emailer-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: emailer-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-emailersrv-v1
  namespace: $(NS)
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/component: microservice
      app.kubernetes.io/instance: emailer-srv-abcxzy
      app.kubernetes.io/managed-by: kustomize
      app.kubernetes.io/name: emailer-srv
      app.kubernetes.io/part-of: micro-starter-kit
      app.kubernetes.io/version: v1
      environment: production
  template:
    metadata:
      annotations:
        note: Hello, I am production!
        org: acmeCorporation
      labels:
        app.kubernetes.io/component: microservice
        app.kubernetes.io/instance: emailer-srv-abcxzy
        app.kubernetes.io/managed-by: kustomize
        app.kubernetes.io/name: emailer-srv
        app.kubernetes.io/part-of: micro-starter-kit
        app.kubernetes.io/version: v1
        environment: production
    spec:
      containers:
      - args:
        - health
        - --address=:8088
        - --check_service=$(MICRO_SERVER_NAME)
        - --check_address=0.0.0.0:8080
        envFrom:
        - configMapRef:
            name: prod-emailerenv-vars-v1-7785t5275f
        image: docker.pkg.github.com/xmlking/micro-starter-kit/micro:v1.15.1
        livenessProbe:
          httpGet:
            path: /health
            port: 8088
          initialDelaySeconds: 30
          periodSeconds: 10
        name: health
      - envFrom:
        - configMapRef:
            name: prod-emailerenv-vars-v1-7785t5275f
        - configMapRef:
            name: prod-env-vars-common-v1-t4df5m59fk
        - secretRef:
            name: prod-secrets-v1-t45t559g92
        image: docker.pkg.github.com/xmlking/micro-starter-kit/emailer-srv:$(IMAGE_VERSION)
        imagePullPolicy: IfNotPresent
        name: srv
        ports:
        - containerPort: 8080
          name: grpc-port
          protocol: TCP
        volumeMounts:
        - mountPath: /config/config.yaml
          name: config
          readOnly: true
          subPath: config.yaml
      volumes:
      - configMap:
          name: prod-emailerconfig-v1-9bm4dcdm69
        name: config
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: gateway-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: gateway
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-gatewaysrv-v1
  namespace: $(NS)
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/component: microservice
      app.kubernetes.io/instance: gateway-abcxzy
      app.kubernetes.io/managed-by: kustomize
      app.kubernetes.io/name: gateway
      app.kubernetes.io/part-of: micro-starter-kit
      app.kubernetes.io/version: v1
      environment: production
  template:
    metadata:
      annotations:
        note: Hello, I am production!
        org: acmeCorporation
      labels:
        app.kubernetes.io/component: microservice
        app.kubernetes.io/instance: gateway-abcxzy
        app.kubernetes.io/managed-by: kustomize
        app.kubernetes.io/name: gateway
        app.kubernetes.io/part-of: micro-starter-kit
        app.kubernetes.io/version: v1
        environment: production
    spec:
      containers:
      - args:
        - health
        - --address=:8088
        - --check_service=$(MICRO_SERVER_NAME)
        - --check_address=0.0.0.0:8080
        envFrom:
        - configMapRef:
            name: prod-gatewayenv-vars-v1-ft7gft5265
        image: docker.pkg.github.com/xmlking/micro-starter-kit/micro:v1.15.1
        livenessProbe:
          httpGet:
            path: /health
            port: 8088
          initialDelaySeconds: 30
          periodSeconds: 10
        name: health
      - args:
        - api
        - --handler=rpc
        - --enable_rpc=true
        - --address=0.0.0.0:8090
        env:
        - name: MICRO_API_ADDRESS
          value: 0.0.0.0:8090
        envFrom:
        - configMapRef:
            name: prod-gatewayenv-vars-v1-ft7gft5265
        - configMapRef:
            name: prod-env-vars-common-v1-t4df5m59fk
        image: docker.pkg.github.com/xmlking/micro-starter-kit/micro:v1.15.1
        imagePullPolicy: IfNotPresent
        name: srv
        ports:
        - containerPort: 8090
          name: http-gateway
          protocol: TCP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: greeter-srv-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: greeter-srv
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-greetersrv-v1
  namespace: $(NS)
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/component: microservice
      app.kubernetes.io/instance: greeter-srv-abcxzy
      app.kubernetes.io/managed-by: kustomize
      app.kubernetes.io/name: greeter-srv
      app.kubernetes.io/part-of: micro-starter-kit
      app.kubernetes.io/version: v1
      environment: production
  template:
    metadata:
      annotations:
        note: Hello, I am production!
        org: acmeCorporation
      labels:
        app.kubernetes.io/component: microservice
        app.kubernetes.io/instance: greeter-srv-abcxzy
        app.kubernetes.io/managed-by: kustomize
        app.kubernetes.io/name: greeter-srv
        app.kubernetes.io/part-of: micro-starter-kit
        app.kubernetes.io/version: v1
        environment: production
    spec:
      containers:
      - args:
        - health
        - --address=:8088
        - --check_service=$(MICRO_SERVER_NAME)
        - --check_address=0.0.0.0:8080
        envFrom:
        - configMapRef:
            name: prod-greeterenv-vars-v1-b852cg88tg
        image: docker.pkg.github.com/xmlking/micro-starter-kit/micro:v1.15.1
        livenessProbe:
          httpGet:
            path: /health
            port: 8088
          initialDelaySeconds: 30
          periodSeconds: 10
        name: health
      - envFrom:
        - configMapRef:
            name: prod-greeterenv-vars-v1-b852cg88tg
        - configMapRef:
            name: prod-env-vars-common-v1-t4df5m59fk
        - secretRef:
            name: prod-secrets-v1-t45t559g92
        image: docker.pkg.github.com/xmlking/micro-starter-kit/greeter-srv:$(IMAGE_VERSION)
        imagePullPolicy: IfNotPresent
        name: srv
        ports:
        - containerPort: 8080
          name: grpc-port
          protocol: TCP
        volumeMounts:
        - mountPath: /config/config.yaml
          name: config
          readOnly: true
          subPath: config.yaml
      volumes:
      - configMap:
          name: prod-greeterconfig-v1-9khmk64kh6
        name: config
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: microservice
    app.kubernetes.io/instance: proxy-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: proxy
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-proxysrv-v1
  namespace: $(NS)
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/component: microservice
      app.kubernetes.io/instance: proxy-abcxzy
      app.kubernetes.io/managed-by: kustomize
      app.kubernetes.io/name: proxy
      app.kubernetes.io/part-of: micro-starter-kit
      app.kubernetes.io/version: v1
      environment: production
  template:
    metadata:
      annotations:
        note: Hello, I am production!
        org: acmeCorporation
      labels:
        app.kubernetes.io/component: microservice
        app.kubernetes.io/instance: proxy-abcxzy
        app.kubernetes.io/managed-by: kustomize
        app.kubernetes.io/name: proxy
        app.kubernetes.io/part-of: micro-starter-kit
        app.kubernetes.io/version: v1
        environment: production
    spec:
      containers:
      - args:
        - health
        - --address=:8088
        - --check_service=$(MICRO_SERVER_NAME)
        - --check_address=0.0.0.0:8081
        envFrom:
        - configMapRef:
            name: prod-proxyenv-vars-v1-6f5k5tk78d
        image: docker.pkg.github.com/xmlking/micro-starter-kit/micro:v1.15.1
        livenessProbe:
          httpGet:
            path: /health
            port: 8088
          initialDelaySeconds: 30
          periodSeconds: 10
        name: health
      - args:
        - proxy
        envFrom:
        - configMapRef:
            name: prod-proxyenv-vars-v1-6f5k5tk78d
        - configMapRef:
            name: prod-env-vars-common-v1-t4df5m59fk
        image: docker.pkg.github.com/xmlking/micro-starter-kit/micro:v1.15.1
        imagePullPolicy: IfNotPresent
        name: srv
        ports:
        - containerPort: 8081
          name: grpc-proxy
          protocol: TCP
---
apiVersion: etcd.database.coreos.com/v1beta2
kind: EtcdCluster
metadata:
  annotations:
    note: Hello, I am production!
    org: acmeCorporation
  labels:
    app.kubernetes.io/component: infra
    app.kubernetes.io/instance: etcd-abcxzy
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: etcd
    app.kubernetes.io/part-of: micro-starter-kit
    app.kubernetes.io/version: v1
    environment: production
  name: prod-etcd-cluster-v1
  namespace: $(NS)
spec:
  pod:
    busyboxImage: busybox:1.28.0-glibc
    labels:
      application: micro
  repository: quay.io/coreos/etcd
  size: 3
  version: 3.4.3
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

