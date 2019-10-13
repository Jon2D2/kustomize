# Feature Test for TestCase variable-ref


This folder contains files for old test-case variable-ref

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
mkdir -p ${DEMO_HOME}/in/overlay
mkdir -p ${DEMO_HOME}/in/package
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/in/overlay/kustomization.yaml
namePrefix: dev-
resources:
- ../package

EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/in/package/kustomization.yaml
namePrefix: base-
resources:
 - cockroachdb-statefulset-secure.yaml
 - cronjob.yaml
configMapGenerator:
- name: test-config-map
  literals:
    - foo=bar
    - baz=qux
vars:
 - name: CDB_PUBLIC_SVC
   objref:
        kind: Service
        name: cockroachdb-public
        apiVersion: v1
   fieldref:
        fieldpath: metadata.name
 - name: CDB_STATEFULSET_NAME
   objref:
        kind: StatefulSet
        name: cockroachdb
        apiVersion: apps/v1beta1
   fieldref:
        fieldpath: metadata.name
 - name: CDB_STATEFULSET_SVC
   objref:
        kind: Service
        name: cockroachdb
        apiVersion: v1
   fieldref:
        fieldpath: metadata.name

 - name: TEST_CONFIG_MAP
   objref:
        kind: ConfigMap
        name: test-config-map
        apiVersion: v1
   fieldref:
        fieldpath: metadata.name
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/in/package/cockroachdb-statefulset-secure.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cockroachdb
  labels:
    app: cockroachdb
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: cockroachdb
  labels:
    app: cockroachdb
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - get
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: cockroachdb
  labels:
    app: cockroachdb
rules:
- apiGroups:
  - certificates.k8s.io
  resources:
  - certificatesigningrequests
  verbs:
  - create
  - get
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: cockroachdb
  labels:
    app: cockroachdb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cockroachdb
subjects:
- kind: ServiceAccount
  name: cockroachdb
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cockroachdb
  labels:
    app: cockroachdb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cockroachdb
subjects:
- kind: ServiceAccount
  name: cockroachdb
  namespace: default
---
apiVersion: v1
kind: Service
metadata:
  # This service is meant to be used by clients of the database. It exposes a ClusterIP that will
  # automatically load balance connections to the different database pods.
  name: cockroachdb-public
  labels:
    app: cockroachdb
spec:
  ports:
  # The main port, served by gRPC, serves Postgres-flavor SQL, internode
  # traffic and the cli.
  - port: 26257
    targetPort: 26257
    name: grpc
  # The secondary port serves the UI as well as health and debug endpoints.
  - port: 8080
    targetPort: 8080
    name: http
  selector:
    app: cockroachdb
---
apiVersion: v1
kind: Service
metadata:
  # This service only exists to create DNS entries for each pod in the stateful
  # set such that they can resolve each other's IP addresses. It does not
  # create a load-balanced ClusterIP and should not be used directly by clients
  # in most circumstances.
  name: cockroachdb
  labels:
    app: cockroachdb
  annotations:
    # This is needed to make the peer-finder work properly and to help avoid
    # edge cases where instance 0 comes up after losing its data and needs to
    # decide whether it should create a new cluster or try to join an existing
    # one. If it creates a new cluster when it should have joined an existing
    # one, we'd end up with two separate clusters listening at the same service
    # endpoint, which would be very bad.
    service.alpha.kubernetes.io/tolerate-unready-endpoints: "true"
    # Enable automatic monitoring of all instances when Prometheus is running in the cluster.
    prometheus.io/scrape: "true"
    prometheus.io/path: "_status/vars"
    prometheus.io/port: "8080"
spec:
  ports:
  - port: 26257
    targetPort: 26257
    name: grpc
  - port: 8080
    targetPort: 8080
    name: http
  clusterIP: None
  selector:
    app: cockroachdb
---
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: cockroachdb-budget
  labels:
    app: cockroachdb
spec:
  selector:
    matchLabels:
      app: cockroachdb
  maxUnavailable: 1
---
apiVersion: apps/v1beta1
kind: StatefulSet
metadata:
  name: cockroachdb
spec:
  serviceName: "cockroachdb"
  replicas: 3
  template:
    metadata:
      labels:
        app: cockroachdb
    spec:
      serviceAccountName: cockroachdb
      # Init containers are run only once in the lifetime of a pod, before
      # it's started up for the first time. It has to exit successfully
      # before the pod's main containers are allowed to start.
      initContainers:
      # The init-certs container sends a certificate signing request to the
      # kubernetes cluster.
      # You can see pending requests using: kubectl get csr
      # CSRs can be approved using:         kubectl certificate approve <csr name>
      #
      # All addresses used to contact a node must be specified in the --addresses arg.
      #
      # In addition to the node certificate and key, the init-certs entrypoint will symlink
      # the cluster CA to the certs directory.
      - name: init-certs
        image: cockroachdb/cockroach-k8s-request-cert:0.2
        imagePullPolicy: IfNotPresent
        command:
        - "/bin/ash"
        - "-ecx"
        - "/request-cert -namespace=${POD_NAMESPACE} -certs-dir=/cockroach-certs -type=node -addresses=localhost,127.0.0.1,${POD_IP},$(hostname -f),$(hostname -f|cut -f 1-2 -d '.'),$(CDB_PUBLIC_SVC) -symlink-ca-from=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        volumeMounts:
        - name: certs
          mountPath: /cockroach-certs

      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - cockroachdb
              topologyKey: kubernetes.io/hostname
      containers:
      - name: cockroachdb
        image: cockroachdb/cockroach:v1.1.5
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 26257
          name: grpc
        - containerPort: 8080
          name: http
        volumeMounts:
        - name: datadir
          mountPath: /cockroach/cockroach-data
        - name: certs
          mountPath: /cockroach/cockroach-certs
        command:
          - "/bin/bash"
          - "-ecx"
          # The use of qualified `hostname -f` is crucial:
          # Other nodes aren't able to look up the unqualified hostname.
          # Once 2.0 is out, we should be able to switch from --host to --advertise-host to make port-forwarding work to the main port.
          - "exec /cockroach/cockroach start --logtostderr --certs-dir /cockroach/cockroach-certs --host $(hostname -f) --http-host 0.0.0.0 --join $(CDB_STATEFULSET_NAME)-0.$(CDB_STATEFULSET_SVC),$(CDB_STATEFULSET_NAME)-1.$(CDB_STATEFULSET_SVC),$(CDB_STATEFULSET_NAME)-2.$(CDB_STATEFULSET_SVC) --cache 25% --max-sql-memory 25%"
      # No pre-stop hook is required, a SIGTERM plus some time is all that's
      # needed for graceful shutdown of a node.
      terminationGracePeriodSeconds: 60
      volumes:
      - name: datadir
        persistentVolumeClaim:
          claimName: datadir
      - name: certs
        emptyDir: {}
  updateStrategy:
    type: RollingUpdate
  volumeClaimTemplates:
  - metadata:
      name: datadir
    spec:
      accessModes:
        - "ReadWriteOnce"
      resources:
        requests:
          storage: 1Gi
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/in/package/cronjob.yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: cronjob-example
spec:
  schedule: "*/1 * * * *"
  concurrencyPolicy: Forbid
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: cronjob-example
            image: cockroachdb/cockroach:v1.1.5
            command:
            - echo
            - "$(CDB_STATEFULSET_NAME)"
            - "$(TEST_CONFIG_MAP)"
            env:
              - name: CDB_PUBLIC_SVC
                value: "$(CDB_PUBLIC_SVC)"
EOF
```

## Execution

<!-- @build @test -->
```bash
kustomize build ${DEMO_HOME}/in/overlay -o ${DEMO_HOME}/actual.yaml
```

## Verification


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: cockroachdb
  name: dev-base-cockroachdb
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  labels:
    app: cockroachdb
  name: dev-base-cockroachdb
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - get
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  labels:
    app: cockroachdb
  name: dev-base-cockroachdb
rules:
- apiGroups:
  - certificates.k8s.io
  resources:
  - certificatesigningrequests
  verbs:
  - create
  - get
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  labels:
    app: cockroachdb
  name: dev-base-cockroachdb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: dev-base-cockroachdb
subjects:
- kind: ServiceAccount
  name: dev-base-cockroachdb
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  labels:
    app: cockroachdb
  name: dev-base-cockroachdb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: dev-base-cockroachdb
subjects:
- kind: ServiceAccount
  name: dev-base-cockroachdb
  namespace: default
---
apiVersion: v1
data:
  baz: qux
  foo: bar
kind: ConfigMap
metadata:
  name: dev-base-test-config-map-b2g2dmd64b
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    prometheus.io/path: _status/vars
    prometheus.io/port: "8080"
    prometheus.io/scrape: "true"
    service.alpha.kubernetes.io/tolerate-unready-endpoints: "true"
  labels:
    app: cockroachdb
  name: dev-base-cockroachdb
spec:
  clusterIP: None
  ports:
  - name: grpc
    port: 26257
    targetPort: 26257
  - name: http
    port: 8080
    targetPort: 8080
  selector:
    app: cockroachdb
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: cockroachdb
  name: dev-base-cockroachdb-public
spec:
  ports:
  - name: grpc
    port: 26257
    targetPort: 26257
  - name: http
    port: 8080
    targetPort: 8080
  selector:
    app: cockroachdb
---
apiVersion: apps/v1beta1
kind: StatefulSet
metadata:
  name: dev-base-cockroachdb
spec:
  replicas: 3
  serviceName: dev-base-cockroachdb
  template:
    metadata:
      labels:
        app: cockroachdb
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - cockroachdb
              topologyKey: kubernetes.io/hostname
            weight: 100
      containers:
      - command:
        - /bin/bash
        - -ecx
        - exec /cockroach/cockroach start --logtostderr --certs-dir /cockroach/cockroach-certs
          --host $(hostname -f) --http-host 0.0.0.0 --join dev-base-cockroachdb-0.dev-base-cockroachdb,dev-base-cockroachdb-1.dev-base-cockroachdb,dev-base-cockroachdb-2.dev-base-cockroachdb
          --cache 25% --max-sql-memory 25%
        image: cockroachdb/cockroach:v1.1.5
        imagePullPolicy: IfNotPresent
        name: cockroachdb
        ports:
        - containerPort: 26257
          name: grpc
        - containerPort: 8080
          name: http
        volumeMounts:
        - mountPath: /cockroach/cockroach-data
          name: datadir
        - mountPath: /cockroach/cockroach-certs
          name: certs
      initContainers:
      - command:
        - /bin/ash
        - -ecx
        - /request-cert -namespace=${POD_NAMESPACE} -certs-dir=/cockroach-certs -type=node
          -addresses=localhost,127.0.0.1,${POD_IP},$(hostname -f),$(hostname -f|cut
          -f 1-2 -d '.'),dev-base-cockroachdb-public -symlink-ca-from=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: cockroachdb/cockroach-k8s-request-cert:0.2
        imagePullPolicy: IfNotPresent
        name: init-certs
        volumeMounts:
        - mountPath: /cockroach-certs
          name: certs
      serviceAccountName: dev-base-cockroachdb
      terminationGracePeriodSeconds: 60
      volumes:
      - name: datadir
        persistentVolumeClaim:
          claimName: datadir
      - emptyDir: {}
        name: certs
  updateStrategy:
    type: RollingUpdate
  volumeClaimTemplates:
  - metadata:
      name: datadir
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
---
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: dev-base-cronjob-example
spec:
  concurrencyPolicy: Forbid
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - command:
            - echo
            - dev-base-cockroachdb
            - dev-base-test-config-map-b2g2dmd64b
            env:
            - name: CDB_PUBLIC_SVC
              value: dev-base-cockroachdb-public
            image: cockroachdb/cockroach:v1.1.5
            name: cronjob-example
  schedule: '*/1 * * * *'
---
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  labels:
    app: cockroachdb
  name: dev-base-cockroachdb-budget
spec:
  maxUnavailable: 1
  selector:
    matchLabels:
      app: cockroachdb
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual.yaml $DEMO_HOME/expected.yaml | wc -l); \
echo $?
```

