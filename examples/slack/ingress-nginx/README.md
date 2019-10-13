# Feature Test for Issue ingress-nginx

The original content has been build from the [ingress-nginx repo](https://github.com/kubernetes/ingress-nginx/)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/aws
mkdir -p ${DEMO_HOME}/aws/l4
mkdir -p ${DEMO_HOME}/aws/l7
mkdir -p ${DEMO_HOME}/aws/nlb
mkdir -p ${DEMO_HOME}/baremetal
mkdir -p ${DEMO_HOME}/cloud-generic
mkdir -p ${DEMO_HOME}/cluster-wide
mkdir -p ${DEMO_HOME}/grafana
mkdir -p ${DEMO_HOME}/minikube
mkdir -p ${DEMO_HOME}/prometheus
mkdir -p ${DEMO_HOME}/static
mkdir -p ${DEMO_HOME}/static/provider
mkdir -p ${DEMO_HOME}/static/provider/aws
mkdir -p ${DEMO_HOME}/static/provider/baremetal
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/aws/l4/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../cloud-generic
patchesStrategicMerge:
- service-l4.yaml
configMapGenerator:
- name: nginx-configuration
  behavior: merge
  literals:
  - use-proxy-protocol=true
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/aws/l7/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../cloud-generic
patchesStrategicMerge:
- service-l7.yaml
configMapGenerator:
- name: nginx-configuration
  behavior: merge
  literals:
  - use-proxy-protocol=false
  - use-forwarded-headers=true
  - proxy-real-ip-cidr=0.0.0.0/0 # restrict this to the IP addresses of ELB
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/aws/nlb/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../cloud-generic
patchesStrategicMerge:
- service-nlb.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/baremetal/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../cloud-generic
patchesStrategicMerge:
- service-nodeport.yaml
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cloud-generic/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: ingress-nginx
commonLabels:
  app.kubernetes.io/name: ingress-nginx
  app.kubernetes.io/part-of: ingress-nginx
resources:
- deployment.yaml
- role-binding.yaml
- role.yaml
- service-account.yaml
- service.yaml
images:
- name: quay.io/kubernetes-ingress-controller/nginx-ingress-controller
  newTag: 0.25.1
vars:
- fieldref:
    fieldPath: metadata.name
  name: NGINX_CONFIGMAP_NAME
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: nginx-configuration
- fieldref:
    fieldPath: metadata.name
  name: TCP_CONFIGMAP_NAME
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: tcp-services
- fieldref:
    fieldPath: metadata.name
  name: UDP_CONFIGMAP_NAME
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: udp-services
- fieldref:
    fieldPath: metadata.name
  name: SERVICE_NAME
  objref:
    apiVersion: v1
    kind: Service
    name: ingress-nginx
configMapGenerator:
- name: nginx-configuration
- name: tcp-services
- name: udp-services
generatorOptions:
  disableNameSuffixHash: true
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cluster-wide/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
commonLabels:
  app.kubernetes.io/name: ingress-nginx
  app.kubernetes.io/part-of: ingress-nginx
resources:
- cluster-role.yaml
- cluster-role-binding.yaml
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/grafana/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: ingress-nginx
commonLabels:
  app.kubernetes.io/name: grafana
  app.kubernetes.io/part-of: ingress-nginx
resources:
- deployment.yaml
- service.yaml
images:
- name: grafana/grafana
  newTag: 6.1.6
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/minikube/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: ingress-nginx
bases:
- ../baremetal
- ../cluster-wide
images:
- name: quay.io/kubernetes-ingress-controller/nginx-ingress-controller
  newName: ingress-controller/nginx-ingress-controller
  newTag: dev
EOF
```


### Preparation Step KustomizationFile8

<!-- @createKustomizationFile8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/prometheus/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: ingress-nginx
commonLabels:
  app.kubernetes.io/name: prometheus
  app.kubernetes.io/part-of: ingress-nginx
resources:
- role.yaml
- service-account.yaml
- role-binding.yaml
- deployment.yaml
- service.yaml
images:
- name: prom/prometheus
  newTag: v2.3.2
configMapGenerator:
- name: prometheus-configuration
  files:
  - prometheus.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/aws/l4/service-l4.yaml
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
  annotations:
    # Enable PROXY protocol
    service.beta.kubernetes.io/aws-load-balancer-proxy-protocol: "*"
    # Ensure the ELB idle timeout is less than nginx keep-alive timeout. By default,
    # NGINX keep-alive is set to 75s. If using WebSockets, the value will need to be
    # increased to '3600' to avoid any potential issues.
    service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: "60"
spec:
  externalTrafficPolicy: Cluster
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/aws/l7/service-l7.yaml
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
  annotations:
    # replace with the correct value of the generated certificate in the AWS console
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: "arn:aws:acm:us-west-2:XXXXXXXX:certificate/XXXXXX-XXXXXXX-XXXXXXX-XXXXXXXX"
    # the backend instances are HTTP
    service.beta.kubernetes.io/aws-load-balancer-backend-protocol: "http"
    # Map port 443
    service.beta.kubernetes.io/aws-load-balancer-ssl-ports: "https"
    # Ensure the ELB idle timeout is less than nginx keep-alive timeout. By default,
    # NGINX keep-alive is set to 75s. If using WebSockets, the value will need to be
    # increased to '3600' to avoid any potential issues.
    service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: "60"
spec:
  externalTrafficPolicy: Cluster
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/aws/nlb/service-nlb.yaml
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
  annotations:
    # by default the type is elb (classic load balancer).
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/baremetal/service-nodeport.yaml
apiVersion: v1
kind: Service
metadata:
  name: ingress-nginx
spec:
  type: NodePort
  ports:
    - name: http
      port: 80
      targetPort: 80
      protocol: TCP
    - name: https
      port: 443
      targetPort: 443
      protocol: TCP
  externalTrafficPolicy: Cluster
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cloud-generic/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-ingress-controller
spec:
  replicas: 1
  template:
    metadata:
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
    spec:
      serviceAccountName: nginx-ingress-serviceaccount
      containers:
        - name: nginx-ingress-controller
          image: quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.25.1
          args:
            - /nginx-ingress-controller
            - --configmap=$(POD_NAMESPACE)/$(NGINX_CONFIGMAP_NAME)
            - --tcp-services-configmap=$(POD_NAMESPACE)/$(TCP_CONFIGMAP_NAME)
            - --udp-services-configmap=$(POD_NAMESPACE)/$(UDP_CONFIGMAP_NAME)
            - --publish-service=$(POD_NAMESPACE)/$(SERVICE_NAME)
            - --annotations-prefix=nginx.ingress.kubernetes.io
          securityContext:
            allowPrivilegeEscalation: true
            capabilities:
              drop:
                - ALL
              add:
                - NET_BIND_SERVICE
            # www-data -> 33
            runAsUser: 33
          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
          ports:
            - name: http
              containerPort: 80
            - name: https
              containerPort: 443
          livenessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            initialDelaySeconds: 10
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 10
          readinessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 10
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cloud-generic/role-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: nginx-ingress-role-nisa-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: nginx-ingress-role
subjects:
  - kind: ServiceAccount
    name: nginx-ingress-serviceaccount
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cloud-generic/role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: nginx-ingress-role
rules:
  - apiGroups:
      - ""
    resources:
      - configmaps
      - pods
      - secrets
      - namespaces
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - configmaps
    resourceNames:
      # Defaults to "<election-id>-<ingress-class>"
      # Here: "<ingress-controller-leader>-<nginx>"
      # This has to be adapted if you change either parameter
      # when launching the nginx-ingress-controller.
      - "ingress-controller-leader-nginx"
    verbs:
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - configmaps
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - endpoints
    verbs:
      - get
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cloud-generic/service-account.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nginx-ingress-serviceaccount
EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cloud-generic/service.yaml
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
spec:
  externalTrafficPolicy: Local
  type: LoadBalancer
  ports:
    - name: http
      port: 80
      targetPort: http
    - name: https
      port: 443
      targetPort: https
EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cluster-wide/cluster-role-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: nginx-ingress-clusterrole-nisa-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: nginx-ingress-clusterrole
subjects:
  - kind: ServiceAccount
    name: nginx-ingress-serviceaccount
EOF
```


### Preparation Step Resource10

<!-- @createResource10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/cluster-wide/cluster-role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: nginx-ingress-clusterrole
rules:
  - apiGroups:
      - ""
    resources:
      - configmaps
      - endpoints
      - nodes
      - pods
      - secrets
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - services
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
  - apiGroups:
      - "extensions"
      - "networking.k8s.io"
    resources:
      - ingresses
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - "extensions"
      - "networking.k8s.io"
    resources:
      - ingresses/status
    verbs:
      - update

EOF
```


### Preparation Step Resource11

<!-- @createResource11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/grafana/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grafana
  namespace: ingress-nginx
spec:
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    spec:
      containers:
        - image: grafana/grafana
          name: grafana
          ports:
            - containerPort: 3000
              protocol: TCP
          resources:
            limits:
              cpu: 500m
              memory: 2500Mi
            requests:
              cpu: 100m
              memory: 100Mi
          volumeMounts:
            - mountPath: /var/lib/grafana
              name: data
      restartPolicy: Always
      volumes:
        - emptyDir: {}
          name: data
EOF
```


### Preparation Step Resource12

<!-- @createResource12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/grafana/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: grafana
spec:
  ports:
    - port: 3000
      protocol: TCP
      targetPort: 3000
  type: NodePort
EOF
```


### Preparation Step Resource13

<!-- @createResource13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/prometheus/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus-server
spec:
  replicas: 1
  template:
    spec:
      serviceAccountName: prometheus-server
      containers:
        - name: prometheus
          image: prom/prometheus
          args:
            - "--config.file=/etc/prometheus/prometheus.yaml"
            - "--storage.tsdb.path=/prometheus/"
          ports:
            - containerPort: 9090
          volumeMounts:
            - name: prometheus-config-volume
              mountPath: /etc/prometheus/
            - name: prometheus-storage-volume
              mountPath: /prometheus/
      volumes:
        - name: prometheus-config-volume
          configMap:
            name: prometheus-configuration
        - name: prometheus-storage-volume
          emptyDir: {}
EOF
```


### Preparation Step Resource14

<!-- @createResource14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/prometheus/prometheus.yaml
global:
  scrape_interval: 10s
scrape_configs:
- job_name: 'ingress-nginx-endpoints'
  kubernetes_sd_configs:
  - role: pod
    namespaces:
      names:
      - ingress-nginx
  relabel_configs:
  - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
    action: keep
    regex: true
  - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scheme]
    action: replace
    target_label: __scheme__
    regex: (https?)
  - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
    action: replace
    target_label: __metrics_path__
    regex: (.+)
  - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
    action: replace
    target_label: __address__
    regex: ([^:]+)(?::\d+)?;(\d+)
    replacement: $1:$2
  - source_labels: [__meta_kubernetes_service_name]
    regex: prometheus-server
    action: drop
EOF
```


### Preparation Step Resource15

<!-- @createResource15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/prometheus/role-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: prometheus-server
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: prometheus-server
subjects:
  - kind: ServiceAccount
    name: prometheus-server
EOF
```


### Preparation Step Resource16

<!-- @createResource16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/prometheus/role.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: prometheus-server
rules:
  - apiGroups: [""]
    resources:
      - services
      - endpoints
      - pods
    verbs: ["get", "list", "watch"]
EOF
```


### Preparation Step Resource17

<!-- @createResource17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/prometheus/service-account.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: prometheus-server
EOF
```


### Preparation Step Resource18

<!-- @createResource18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/prometheus/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: prometheus-server
spec:
  type: NodePort
  ports:
    - port: 9090
      targetPort: 9090
EOF
```


### Preparation Step Resource19

<!-- @createResource19 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/configmap.yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: nginx-configuration
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---
kind: ConfigMap
apiVersion: v1
metadata:
  name: tcp-services
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---
kind: ConfigMap
apiVersion: v1
metadata:
  name: udp-services
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---
EOF
```


### Preparation Step Resource20

<!-- @createResource20 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/mandatory.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---

kind: ConfigMap
apiVersion: v1
metadata:
  name: nginx-configuration
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---
kind: ConfigMap
apiVersion: v1
metadata:
  name: tcp-services
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---
kind: ConfigMap
apiVersion: v1
metadata:
  name: udp-services
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: nginx-ingress-clusterrole
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
rules:
  - apiGroups:
      - ""
    resources:
      - configmaps
      - endpoints
      - nodes
      - pods
      - secrets
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - services
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
  - apiGroups:
      - "extensions"
      - "networking.k8s.io"
    resources:
      - ingresses
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - "extensions"
      - "networking.k8s.io"
    resources:
      - ingresses/status
    verbs:
      - update

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: nginx-ingress-role
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
rules:
  - apiGroups:
      - ""
    resources:
      - configmaps
      - pods
      - secrets
      - namespaces
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - configmaps
    resourceNames:
      # Defaults to "<election-id>-<ingress-class>"
      # Here: "<ingress-controller-leader>-<nginx>"
      # This has to be adapted if you change either parameter
      # when launching the nginx-ingress-controller.
      - "ingress-controller-leader-nginx"
    verbs:
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - configmaps
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - endpoints
    verbs:
      - get

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: nginx-ingress-role-nisa-binding
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: nginx-ingress-role
subjects:
  - kind: ServiceAccount
    name: nginx-ingress-serviceaccount
    namespace: ingress-nginx

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: nginx-ingress-clusterrole-nisa-binding
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: nginx-ingress-clusterrole
subjects:
  - kind: ServiceAccount
    name: nginx-ingress-serviceaccount
    namespace: ingress-nginx

---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-ingress-controller
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: ingress-nginx
      app.kubernetes.io/part-of: ingress-nginx
  template:
    metadata:
      labels:
        app.kubernetes.io/name: ingress-nginx
        app.kubernetes.io/part-of: ingress-nginx
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
    spec:
      serviceAccountName: nginx-ingress-serviceaccount
      containers:
        - name: nginx-ingress-controller
          image: quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.25.1
          args:
            - /nginx-ingress-controller
            - --configmap=$(POD_NAMESPACE)/nginx-configuration
            - --tcp-services-configmap=$(POD_NAMESPACE)/tcp-services
            - --udp-services-configmap=$(POD_NAMESPACE)/udp-services
            - --publish-service=$(POD_NAMESPACE)/ingress-nginx
            - --annotations-prefix=nginx.ingress.kubernetes.io
          securityContext:
            allowPrivilegeEscalation: true
            capabilities:
              drop:
                - ALL
              add:
                - NET_BIND_SERVICE
            # www-data -> 33
            runAsUser: 33
          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
          ports:
            - name: http
              containerPort: 80
            - name: https
              containerPort: 443
          livenessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            initialDelaySeconds: 10
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 10
          readinessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 10

---
EOF
```


### Preparation Step Resource21

<!-- @createResource21 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---

EOF
```


### Preparation Step Resource22

<!-- @createResource22 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/provider/aws/patch-configmap-l4.yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: nginx-configuration
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
data:
  use-proxy-protocol: "true"
EOF
```


### Preparation Step Resource23

<!-- @createResource23 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/provider/aws/patch-configmap-l7.yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: nginx-configuration
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
data:
  use-proxy-protocol: "false"
  use-forwarded-headers: "true"
  proxy-real-ip-cidr: "0.0.0.0/0" # restrict this to the IP addresses of ELB
---

EOF
```


### Preparation Step Resource24

<!-- @createResource24 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/provider/aws/service-l4.yaml
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  annotations:
    # Enable PROXY protocol
    service.beta.kubernetes.io/aws-load-balancer-proxy-protocol: "*"
    # Ensure the ELB idle timeout is less than nginx keep-alive timeout. By default,
    # NGINX keep-alive is set to 75s. If using WebSockets, the value will need to be
    # increased to '3600' to avoid any potential issues.
    service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: "60"
spec:
  type: LoadBalancer
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  ports:
    - name: http
      port: 80
      targetPort: http
    - name: https
      port: 443
      targetPort: https

---

EOF
```


### Preparation Step Resource25

<!-- @createResource25 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/provider/aws/service-l7.yaml
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  annotations:
    # replace with the correct value of the generated certificate in the AWS console
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: "arn:aws:acm:us-west-2:XXXXXXXX:certificate/XXXXXX-XXXXXXX-XXXXXXX-XXXXXXXX"
    # the backend instances are HTTP
    service.beta.kubernetes.io/aws-load-balancer-backend-protocol: "http"
    # Map port 443
    service.beta.kubernetes.io/aws-load-balancer-ssl-ports: "https"
    # Ensure the ELB idle timeout is less than nginx keep-alive timeout. By default,
    # NGINX keep-alive is set to 75s. If using WebSockets, the value will need to be
    # increased to '3600' to avoid any potential issues.
    service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: "60"
spec:
  type: LoadBalancer
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  ports:
    - name: http
      port: 80
      targetPort: http
    - name: https
      port: 443
      targetPort: http

---

EOF
```


### Preparation Step Resource26

<!-- @createResource26 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/provider/aws/service-nlb.yaml
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  annotations:
    # by default the type is elb (classic load balancer).
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
spec:
  # this setting is to make sure the source IP address is preserved.
  externalTrafficPolicy: Local
  type: LoadBalancer
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  ports:
    - name: http
      port: 80
      targetPort: http
    - name: https
      port: 443
      targetPort: https

---

EOF
```


### Preparation Step Resource27

<!-- @createResource27 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/provider/baremetal/service-nodeport.yaml
apiVersion: v1
kind: Service
metadata:
  name: ingress-nginx
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
spec:
  type: NodePort
  ports:
    - name: http
      port: 80
      targetPort: 80
      protocol: TCP
    - name: https
      port: 443
      targetPort: 443
      protocol: TCP
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---

EOF
```


### Preparation Step Resource28

<!-- @createResource28 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/provider/cloud-generic.yaml
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
spec:
  externalTrafficPolicy: Local
  type: LoadBalancer
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  ports:
    - name: http
      port: 80
      targetPort: http
    - name: https
      port: 443
      targetPort: https

---

EOF
```


### Preparation Step Resource29

<!-- @createResource29 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: nginx-ingress-clusterrole
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
rules:
  - apiGroups:
      - ""
    resources:
      - configmaps
      - endpoints
      - nodes
      - pods
      - secrets
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - services
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
  - apiGroups:
      - "extensions"
      - "networking.k8s.io"
    resources:
      - ingresses
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - "extensions"
      - "networking.k8s.io"
    resources:
      - ingresses/status
    verbs:
      - update

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: nginx-ingress-role
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
rules:
  - apiGroups:
      - ""
    resources:
      - configmaps
      - pods
      - secrets
      - namespaces
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - configmaps
    resourceNames:
      # Defaults to "<election-id>-<ingress-class>"
      # Here: "<ingress-controller-leader>-<nginx>"
      # This has to be adapted if you change either parameter
      # when launching the nginx-ingress-controller.
      - "ingress-controller-leader-nginx"
    verbs:
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - configmaps
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - endpoints
    verbs:
      - get

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: nginx-ingress-role-nisa-binding
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: nginx-ingress-role
subjects:
  - kind: ServiceAccount
    name: nginx-ingress-serviceaccount
    namespace: ingress-nginx

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: nginx-ingress-clusterrole-nisa-binding
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: nginx-ingress-clusterrole
subjects:
  - kind: ServiceAccount
    name: nginx-ingress-serviceaccount
    namespace: ingress-nginx

---

EOF
```


### Preparation Step Resource30

<!-- @createResource30 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/static/with-rbac.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-ingress-controller
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: ingress-nginx
      app.kubernetes.io/part-of: ingress-nginx
  template:
    metadata:
      labels:
        app.kubernetes.io/name: ingress-nginx
        app.kubernetes.io/part-of: ingress-nginx
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
    spec:
      serviceAccountName: nginx-ingress-serviceaccount
      containers:
        - name: nginx-ingress-controller
          image: quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.25.1
          args:
            - /nginx-ingress-controller
            - --configmap=$(POD_NAMESPACE)/nginx-configuration
            - --tcp-services-configmap=$(POD_NAMESPACE)/tcp-services
            - --udp-services-configmap=$(POD_NAMESPACE)/udp-services
            - --publish-service=$(POD_NAMESPACE)/ingress-nginx
            - --annotations-prefix=nginx.ingress.kubernetes.io
          securityContext:
            allowPrivilegeEscalation: true
            capabilities:
              drop:
                - ALL
              add:
                - NET_BIND_SERVICE
            # www-data -> 33
            runAsUser: 33
          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
          ports:
            - name: http
              containerPort: 80
            - name: https
              containerPort: 443
          livenessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            initialDelaySeconds: 10
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 10
          readinessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 10

---

EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir -p ${DEMO_HOME}/actual
mkdir -p ${DEMO_HOME}/actual/aws/nlb
mkdir -p ${DEMO_HOME}/actual/aws/l4
mkdir -p ${DEMO_HOME}/actual/aws/l7
mkdir -p ${DEMO_HOME}/actual/cluster-wide
mkdir -p ${DEMO_HOME}/actual/grafana
mkdir -p ${DEMO_HOME}/actual/minikube
mkdir -p ${DEMO_HOME}/actual/baremetal
mkdir -p ${DEMO_HOME}/actual/cloud-generic
mkdir -p ${DEMO_HOME}/actual/prometheus
kustomize build ${DEMO_HOME}/aws/nlb -o ${DEMO_HOME}/actual/aws/nlb
kustomize build ${DEMO_HOME}/aws/l4 -o ${DEMO_HOME}/actual/aws/l4
kustomize build ${DEMO_HOME}/aws/l7 -o ${DEMO_HOME}/actual/aws/l7
kustomize build ${DEMO_HOME}/cluster-wide -o ${DEMO_HOME}/actual/cluster-wide
kustomize build ${DEMO_HOME}/grafana -o ${DEMO_HOME}/actual/grafana
kustomize build ${DEMO_HOME}/minikube -o ${DEMO_HOME}/actual/minikube
kustomize build ${DEMO_HOME}/baremetal -o ${DEMO_HOME}/actual/baremetal
kustomize build ${DEMO_HOME}/cloud-generic -o ${DEMO_HOME}/actual/cloud-generic
kustomize build ${DEMO_HOME}/prometheus -o ${DEMO_HOME}/actual/prometheus
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir -p ${DEMO_HOME}/expected
mkdir -p ${DEMO_HOME}/expected/aws/nlb
mkdir -p ${DEMO_HOME}/expected/aws/l4
mkdir -p ${DEMO_HOME}/expected/aws/l7
mkdir -p ${DEMO_HOME}/expected/cluster-wide
mkdir -p ${DEMO_HOME}/expected/grafana
mkdir -p ${DEMO_HOME}/expected/minikube
mkdir -p ${DEMO_HOME}/expected/baremetal
mkdir -p ${DEMO_HOME}/expected/cloud-generic
mkdir -p ${DEMO_HOME}/expected/prometheus
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l4/apps_v1_deployment_nginx-ingress-controller.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-controller
  namespace: ingress-nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: ingress-nginx
      app.kubernetes.io/part-of: ingress-nginx
  template:
    metadata:
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
      labels:
        app.kubernetes.io/name: ingress-nginx
        app.kubernetes.io/part-of: ingress-nginx
    spec:
      containers:
      - args:
        - /nginx-ingress-controller
        - --configmap=$(POD_NAMESPACE)/nginx-configuration
        - --tcp-services-configmap=$(POD_NAMESPACE)/tcp-services
        - --udp-services-configmap=$(POD_NAMESPACE)/udp-services
        - --publish-service=$(POD_NAMESPACE)/ingress-nginx
        - --annotations-prefix=nginx.ingress.kubernetes.io
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.25.1
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        name: nginx-ingress-controller
        ports:
        - containerPort: 80
          name: http
        - containerPort: 443
          name: https
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        securityContext:
          allowPrivilegeEscalation: true
          capabilities:
            add:
            - NET_BIND_SERVICE
            drop:
            - ALL
          runAsUser: 33
      serviceAccountName: nginx-ingress-serviceaccount
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l4/~g_v1_configmap_nginx-configuration.yaml
apiVersion: v1
data:
  use-proxy-protocol: "true"
kind: ConfigMap
metadata:
  annotations: {}
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-configuration
  namespace: ingress-nginx
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l4/~g_v1_configmap_tcp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: tcp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l4/~g_v1_configmap_udp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: udp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l4/~g_v1_serviceaccount_nginx-ingress-serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l4/~g_v1_service_ingress-nginx.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: "60"
    service.beta.kubernetes.io/aws-load-balancer-proxy-protocol: '*'
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: ingress-nginx
  namespace: ingress-nginx
spec:
  externalTrafficPolicy: Cluster
  ports:
  - name: http
    port: 80
    targetPort: http
  - name: https
    port: 443
    targetPort: https
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  type: LoadBalancer
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l4/rbac.authorization.k8s.io_v1beta1_rolebinding_nginx-ingress-role-nisa-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role-nisa-binding
  namespace: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: nginx-ingress-role
subjects:
- kind: ServiceAccount
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l4/rbac.authorization.k8s.io_v1beta1_role_nginx-ingress-role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role
  namespace: ingress-nginx
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - pods
  - secrets
  - namespaces
  verbs:
  - get
- apiGroups:
  - ""
  resourceNames:
  - ingress-controller-leader-nginx
  resources:
  - configmaps
  verbs:
  - get
  - update
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - create
- apiGroups:
  - ""
  resources:
  - endpoints
  verbs:
  - get
EOF
```


### Verification Step Expected8

<!-- @createExpected8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l7/apps_v1_deployment_nginx-ingress-controller.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-controller
  namespace: ingress-nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: ingress-nginx
      app.kubernetes.io/part-of: ingress-nginx
  template:
    metadata:
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
      labels:
        app.kubernetes.io/name: ingress-nginx
        app.kubernetes.io/part-of: ingress-nginx
    spec:
      containers:
      - args:
        - /nginx-ingress-controller
        - --configmap=$(POD_NAMESPACE)/nginx-configuration
        - --tcp-services-configmap=$(POD_NAMESPACE)/tcp-services
        - --udp-services-configmap=$(POD_NAMESPACE)/udp-services
        - --publish-service=$(POD_NAMESPACE)/ingress-nginx
        - --annotations-prefix=nginx.ingress.kubernetes.io
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.25.1
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        name: nginx-ingress-controller
        ports:
        - containerPort: 80
          name: http
        - containerPort: 443
          name: https
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        securityContext:
          allowPrivilegeEscalation: true
          capabilities:
            add:
            - NET_BIND_SERVICE
            drop:
            - ALL
          runAsUser: 33
      serviceAccountName: nginx-ingress-serviceaccount
EOF
```


### Verification Step Expected9

<!-- @createExpected9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l7/~g_v1_configmap_nginx-configuration.yaml
apiVersion: v1
data:
  proxy-real-ip-cidr: 0.0.0.0/0
  use-forwarded-headers: "true"
  use-proxy-protocol: "false"
kind: ConfigMap
metadata:
  annotations: {}
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-configuration
  namespace: ingress-nginx
EOF
```


### Verification Step Expected10

<!-- @createExpected10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l7/~g_v1_configmap_tcp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: tcp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected11

<!-- @createExpected11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l7/~g_v1_configmap_udp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: udp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected12

<!-- @createExpected12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l7/~g_v1_serviceaccount_nginx-ingress-serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected13

<!-- @createExpected13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l7/~g_v1_service_ingress-nginx.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-backend-protocol: http
    service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: "60"
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: arn:aws:acm:us-west-2:XXXXXXXX:certificate/XXXXXX-XXXXXXX-XXXXXXX-XXXXXXXX
    service.beta.kubernetes.io/aws-load-balancer-ssl-ports: https
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: ingress-nginx
  namespace: ingress-nginx
spec:
  externalTrafficPolicy: Cluster
  ports:
  - name: http
    port: 80
    targetPort: http
  - name: https
    port: 443
    targetPort: https
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  type: LoadBalancer
EOF
```


### Verification Step Expected14

<!-- @createExpected14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l7/rbac.authorization.k8s.io_v1beta1_rolebinding_nginx-ingress-role-nisa-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role-nisa-binding
  namespace: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: nginx-ingress-role
subjects:
- kind: ServiceAccount
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected15

<!-- @createExpected15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/l7/rbac.authorization.k8s.io_v1beta1_role_nginx-ingress-role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role
  namespace: ingress-nginx
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - pods
  - secrets
  - namespaces
  verbs:
  - get
- apiGroups:
  - ""
  resourceNames:
  - ingress-controller-leader-nginx
  resources:
  - configmaps
  verbs:
  - get
  - update
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - create
- apiGroups:
  - ""
  resources:
  - endpoints
  verbs:
  - get
EOF
```


### Verification Step Expected16

<!-- @createExpected16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/nlb/apps_v1_deployment_nginx-ingress-controller.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-controller
  namespace: ingress-nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: ingress-nginx
      app.kubernetes.io/part-of: ingress-nginx
  template:
    metadata:
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
      labels:
        app.kubernetes.io/name: ingress-nginx
        app.kubernetes.io/part-of: ingress-nginx
    spec:
      containers:
      - args:
        - /nginx-ingress-controller
        - --configmap=$(POD_NAMESPACE)/nginx-configuration
        - --tcp-services-configmap=$(POD_NAMESPACE)/tcp-services
        - --udp-services-configmap=$(POD_NAMESPACE)/udp-services
        - --publish-service=$(POD_NAMESPACE)/ingress-nginx
        - --annotations-prefix=nginx.ingress.kubernetes.io
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.25.1
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        name: nginx-ingress-controller
        ports:
        - containerPort: 80
          name: http
        - containerPort: 443
          name: https
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        securityContext:
          allowPrivilegeEscalation: true
          capabilities:
            add:
            - NET_BIND_SERVICE
            drop:
            - ALL
          runAsUser: 33
      serviceAccountName: nginx-ingress-serviceaccount
EOF
```


### Verification Step Expected17

<!-- @createExpected17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/nlb/~g_v1_configmap_nginx-configuration.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-configuration
  namespace: ingress-nginx
EOF
```


### Verification Step Expected18

<!-- @createExpected18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/nlb/~g_v1_configmap_tcp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: tcp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected19

<!-- @createExpected19 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/nlb/~g_v1_configmap_udp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: udp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected20

<!-- @createExpected20 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/nlb/~g_v1_serviceaccount_nginx-ingress-serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected21

<!-- @createExpected21 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/nlb/~g_v1_service_ingress-nginx.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: ingress-nginx
  namespace: ingress-nginx
spec:
  externalTrafficPolicy: Local
  ports:
  - name: http
    port: 80
    targetPort: http
  - name: https
    port: 443
    targetPort: https
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  type: LoadBalancer
EOF
```


### Verification Step Expected22

<!-- @createExpected22 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/nlb/rbac.authorization.k8s.io_v1beta1_rolebinding_nginx-ingress-role-nisa-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role-nisa-binding
  namespace: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: nginx-ingress-role
subjects:
- kind: ServiceAccount
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected23

<!-- @createExpected23 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/aws/nlb/rbac.authorization.k8s.io_v1beta1_role_nginx-ingress-role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role
  namespace: ingress-nginx
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - pods
  - secrets
  - namespaces
  verbs:
  - get
- apiGroups:
  - ""
  resourceNames:
  - ingress-controller-leader-nginx
  resources:
  - configmaps
  verbs:
  - get
  - update
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - create
- apiGroups:
  - ""
  resources:
  - endpoints
  verbs:
  - get
EOF
```


### Verification Step Expected24

<!-- @createExpected24 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/baremetal/apps_v1_deployment_nginx-ingress-controller.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-controller
  namespace: ingress-nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: ingress-nginx
      app.kubernetes.io/part-of: ingress-nginx
  template:
    metadata:
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
      labels:
        app.kubernetes.io/name: ingress-nginx
        app.kubernetes.io/part-of: ingress-nginx
    spec:
      containers:
      - args:
        - /nginx-ingress-controller
        - --configmap=$(POD_NAMESPACE)/nginx-configuration
        - --tcp-services-configmap=$(POD_NAMESPACE)/tcp-services
        - --udp-services-configmap=$(POD_NAMESPACE)/udp-services
        - --publish-service=$(POD_NAMESPACE)/ingress-nginx
        - --annotations-prefix=nginx.ingress.kubernetes.io
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.25.1
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        name: nginx-ingress-controller
        ports:
        - containerPort: 80
          name: http
        - containerPort: 443
          name: https
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        securityContext:
          allowPrivilegeEscalation: true
          capabilities:
            add:
            - NET_BIND_SERVICE
            drop:
            - ALL
          runAsUser: 33
      serviceAccountName: nginx-ingress-serviceaccount
EOF
```


### Verification Step Expected25

<!-- @createExpected25 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/baremetal/~g_v1_configmap_nginx-configuration.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-configuration
  namespace: ingress-nginx
EOF
```


### Verification Step Expected26

<!-- @createExpected26 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/baremetal/~g_v1_configmap_tcp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: tcp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected27

<!-- @createExpected27 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/baremetal/~g_v1_configmap_udp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: udp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected28

<!-- @createExpected28 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/baremetal/~g_v1_serviceaccount_nginx-ingress-serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected29

<!-- @createExpected29 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/baremetal/~g_v1_service_ingress-nginx.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: ingress-nginx
  namespace: ingress-nginx
spec:
  externalTrafficPolicy: Cluster
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 80
  - name: https
    port: 443
    protocol: TCP
    targetPort: 443
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  type: NodePort
EOF
```


### Verification Step Expected30

<!-- @createExpected30 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/baremetal/rbac.authorization.k8s.io_v1beta1_rolebinding_nginx-ingress-role-nisa-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role-nisa-binding
  namespace: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: nginx-ingress-role
subjects:
- kind: ServiceAccount
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected31

<!-- @createExpected31 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/baremetal/rbac.authorization.k8s.io_v1beta1_role_nginx-ingress-role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role
  namespace: ingress-nginx
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - pods
  - secrets
  - namespaces
  verbs:
  - get
- apiGroups:
  - ""
  resourceNames:
  - ingress-controller-leader-nginx
  resources:
  - configmaps
  verbs:
  - get
  - update
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - create
- apiGroups:
  - ""
  resources:
  - endpoints
  verbs:
  - get
EOF
```


### Verification Step Expected32

<!-- @createExpected32 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cloud-generic/apps_v1_deployment_nginx-ingress-controller.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-controller
  namespace: ingress-nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: ingress-nginx
      app.kubernetes.io/part-of: ingress-nginx
  template:
    metadata:
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
      labels:
        app.kubernetes.io/name: ingress-nginx
        app.kubernetes.io/part-of: ingress-nginx
    spec:
      containers:
      - args:
        - /nginx-ingress-controller
        - --configmap=$(POD_NAMESPACE)/nginx-configuration
        - --tcp-services-configmap=$(POD_NAMESPACE)/tcp-services
        - --udp-services-configmap=$(POD_NAMESPACE)/udp-services
        - --publish-service=$(POD_NAMESPACE)/ingress-nginx
        - --annotations-prefix=nginx.ingress.kubernetes.io
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.25.1
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        name: nginx-ingress-controller
        ports:
        - containerPort: 80
          name: http
        - containerPort: 443
          name: https
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        securityContext:
          allowPrivilegeEscalation: true
          capabilities:
            add:
            - NET_BIND_SERVICE
            drop:
            - ALL
          runAsUser: 33
      serviceAccountName: nginx-ingress-serviceaccount
EOF
```


### Verification Step Expected33

<!-- @createExpected33 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cloud-generic/~g_v1_configmap_nginx-configuration.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-configuration
  namespace: ingress-nginx
EOF
```


### Verification Step Expected34

<!-- @createExpected34 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cloud-generic/~g_v1_configmap_tcp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: tcp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected35

<!-- @createExpected35 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cloud-generic/~g_v1_configmap_udp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: udp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected36

<!-- @createExpected36 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cloud-generic/~g_v1_serviceaccount_nginx-ingress-serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected37

<!-- @createExpected37 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cloud-generic/~g_v1_service_ingress-nginx.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: ingress-nginx
  namespace: ingress-nginx
spec:
  externalTrafficPolicy: Local
  ports:
  - name: http
    port: 80
    targetPort: http
  - name: https
    port: 443
    targetPort: https
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  type: LoadBalancer
EOF
```


### Verification Step Expected38

<!-- @createExpected38 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cloud-generic/rbac.authorization.k8s.io_v1beta1_rolebinding_nginx-ingress-role-nisa-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role-nisa-binding
  namespace: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: nginx-ingress-role
subjects:
- kind: ServiceAccount
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected39

<!-- @createExpected39 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cloud-generic/rbac.authorization.k8s.io_v1beta1_role_nginx-ingress-role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role
  namespace: ingress-nginx
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - pods
  - secrets
  - namespaces
  verbs:
  - get
- apiGroups:
  - ""
  resourceNames:
  - ingress-controller-leader-nginx
  resources:
  - configmaps
  verbs:
  - get
  - update
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - create
- apiGroups:
  - ""
  resources:
  - endpoints
  verbs:
  - get
EOF
```


### Verification Step Expected40

<!-- @createExpected40 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cluster-wide/rbac.authorization.k8s.io_v1beta1_clusterrolebinding_nginx-ingress-clusterrole-nisa-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-clusterrole-nisa-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: nginx-ingress-clusterrole
subjects:
- kind: ServiceAccount
  name: nginx-ingress-serviceaccount
EOF
```


### Verification Step Expected41

<!-- @createExpected41 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cluster-wide/rbac.authorization.k8s.io_v1beta1_clusterrole_nginx-ingress-clusterrole.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-clusterrole
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - endpoints
  - nodes
  - pods
  - secrets
  verbs:
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - get
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
- apiGroups:
  - extensions
  - networking.k8s.io
  resources:
  - ingresses
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - extensions
  - networking.k8s.io
  resources:
  - ingresses/status
  verbs:
  - update
EOF
```


### Verification Step Expected42

<!-- @createExpected42 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/grafana/apps_v1_deployment_grafana.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/name: grafana
    app.kubernetes.io/part-of: ingress-nginx
  name: grafana
  namespace: ingress-nginx
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: grafana
      app.kubernetes.io/part-of: ingress-nginx
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app.kubernetes.io/name: grafana
        app.kubernetes.io/part-of: ingress-nginx
    spec:
      containers:
      - image: grafana/grafana:6.1.6
        name: grafana
        ports:
        - containerPort: 3000
          protocol: TCP
        resources:
          limits:
            cpu: 500m
            memory: 2500Mi
          requests:
            cpu: 100m
            memory: 100Mi
        volumeMounts:
        - mountPath: /var/lib/grafana
          name: data
      restartPolicy: Always
      volumes:
      - emptyDir: {}
        name: data
EOF
```


### Verification Step Expected43

<!-- @createExpected43 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/grafana/~g_v1_service_grafana.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/name: grafana
    app.kubernetes.io/part-of: ingress-nginx
  name: grafana
  namespace: ingress-nginx
spec:
  ports:
  - port: 3000
    protocol: TCP
    targetPort: 3000
  selector:
    app.kubernetes.io/name: grafana
    app.kubernetes.io/part-of: ingress-nginx
  type: NodePort
EOF
```


### Verification Step Expected44

<!-- @createExpected44 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/minikube/apps_v1_deployment_nginx-ingress-controller.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-controller
  namespace: ingress-nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: ingress-nginx
      app.kubernetes.io/part-of: ingress-nginx
  template:
    metadata:
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
      labels:
        app.kubernetes.io/name: ingress-nginx
        app.kubernetes.io/part-of: ingress-nginx
    spec:
      containers:
      - args:
        - /nginx-ingress-controller
        - --configmap=$(POD_NAMESPACE)/nginx-configuration
        - --tcp-services-configmap=$(POD_NAMESPACE)/tcp-services
        - --udp-services-configmap=$(POD_NAMESPACE)/udp-services
        - --publish-service=$(POD_NAMESPACE)/ingress-nginx
        - --annotations-prefix=nginx.ingress.kubernetes.io
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: ingress-controller/nginx-ingress-controller:dev
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        name: nginx-ingress-controller
        ports:
        - containerPort: 80
          name: http
        - containerPort: 443
          name: https
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        securityContext:
          allowPrivilegeEscalation: true
          capabilities:
            add:
            - NET_BIND_SERVICE
            drop:
            - ALL
          runAsUser: 33
      serviceAccountName: nginx-ingress-serviceaccount
EOF
```


### Verification Step Expected45

<!-- @createExpected45 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/minikube/~g_v1_configmap_nginx-configuration.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-configuration
  namespace: ingress-nginx
EOF
```


### Verification Step Expected46

<!-- @createExpected46 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/minikube/~g_v1_configmap_tcp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: tcp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected47

<!-- @createExpected47 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/minikube/~g_v1_configmap_udp-services.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: udp-services
  namespace: ingress-nginx
EOF
```


### Verification Step Expected48

<!-- @createExpected48 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/minikube/~g_v1_serviceaccount_nginx-ingress-serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected49

<!-- @createExpected49 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/minikube/~g_v1_service_ingress-nginx.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: ingress-nginx
  namespace: ingress-nginx
spec:
  externalTrafficPolicy: Cluster
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 80
  - name: https
    port: 443
    protocol: TCP
    targetPort: 443
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  type: NodePort
EOF
```


### Verification Step Expected50

<!-- @createExpected50 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/minikube/rbac.authorization.k8s.io_v1beta1_clusterrolebinding_nginx-ingress-clusterrole-nisa-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-clusterrole-nisa-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: nginx-ingress-clusterrole
subjects:
- kind: ServiceAccount
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected51

<!-- @createExpected51 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/minikube/rbac.authorization.k8s.io_v1beta1_clusterrole_nginx-ingress-clusterrole.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-clusterrole
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - endpoints
  - nodes
  - pods
  - secrets
  verbs:
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - get
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
- apiGroups:
  - extensions
  - networking.k8s.io
  resources:
  - ingresses
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - extensions
  - networking.k8s.io
  resources:
  - ingresses/status
  verbs:
  - update
EOF
```


### Verification Step Expected52

<!-- @createExpected52 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/minikube/rbac.authorization.k8s.io_v1beta1_rolebinding_nginx-ingress-role-nisa-binding.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role-nisa-binding
  namespace: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: nginx-ingress-role
subjects:
- kind: ServiceAccount
  name: nginx-ingress-serviceaccount
  namespace: ingress-nginx
EOF
```


### Verification Step Expected53

<!-- @createExpected53 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/minikube/rbac.authorization.k8s.io_v1beta1_role_nginx-ingress-role.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: nginx-ingress-role
  namespace: ingress-nginx
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - pods
  - secrets
  - namespaces
  verbs:
  - get
- apiGroups:
  - ""
  resourceNames:
  - ingress-controller-leader-nginx
  resources:
  - configmaps
  verbs:
  - get
  - update
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - create
- apiGroups:
  - ""
  resources:
  - endpoints
  verbs:
  - get
EOF
```


### Verification Step Expected54

<!-- @createExpected54 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prometheus/apps_v1_deployment_prometheus-server.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/name: prometheus
    app.kubernetes.io/part-of: ingress-nginx
  name: prometheus-server
  namespace: ingress-nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: prometheus
      app.kubernetes.io/part-of: ingress-nginx
  template:
    metadata:
      labels:
        app.kubernetes.io/name: prometheus
        app.kubernetes.io/part-of: ingress-nginx
    spec:
      containers:
      - args:
        - --config.file=/etc/prometheus/prometheus.yaml
        - --storage.tsdb.path=/prometheus/
        image: prom/prometheus:v2.3.2
        name: prometheus
        ports:
        - containerPort: 9090
        volumeMounts:
        - mountPath: /etc/prometheus/
          name: prometheus-config-volume
        - mountPath: /prometheus/
          name: prometheus-storage-volume
      serviceAccountName: prometheus-server
      volumes:
      - configMap:
          name: prometheus-configuration-bc6bcg7b65
        name: prometheus-config-volume
      - emptyDir: {}
        name: prometheus-storage-volume
EOF
```


### Verification Step Expected55

<!-- @createExpected55 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prometheus/~g_v1_configmap_prometheus-configuration-bc6bcg7b65.yaml
apiVersion: v1
data:
  prometheus.yaml: |
    global:
      scrape_interval: 10s
    scrape_configs:
    - job_name: 'ingress-nginx-endpoints'
      kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
          - ingress-nginx
      relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scheme]
        action: replace
        target_label: __scheme__
        regex: (https?)
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
      - source_labels: [__meta_kubernetes_service_name]
        regex: prometheus-server
        action: drop
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/name: prometheus
    app.kubernetes.io/part-of: ingress-nginx
  name: prometheus-configuration-bc6bcg7b65
  namespace: ingress-nginx
EOF
```


### Verification Step Expected56

<!-- @createExpected56 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prometheus/~g_v1_serviceaccount_prometheus-server.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/name: prometheus
    app.kubernetes.io/part-of: ingress-nginx
  name: prometheus-server
  namespace: ingress-nginx
EOF
```


### Verification Step Expected57

<!-- @createExpected57 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prometheus/~g_v1_service_prometheus-server.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/name: prometheus
    app.kubernetes.io/part-of: ingress-nginx
  name: prometheus-server
  namespace: ingress-nginx
spec:
  ports:
  - port: 9090
    targetPort: 9090
  selector:
    app.kubernetes.io/name: prometheus
    app.kubernetes.io/part-of: ingress-nginx
  type: NodePort
EOF
```


### Verification Step Expected58

<!-- @createExpected58 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prometheus/rbac.authorization.k8s.io_v1_rolebinding_prometheus-server.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  labels:
    app.kubernetes.io/name: prometheus
    app.kubernetes.io/part-of: ingress-nginx
  name: prometheus-server
  namespace: ingress-nginx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: prometheus-server
subjects:
- kind: ServiceAccount
  name: prometheus-server
  namespace: ingress-nginx
EOF
```


### Verification Step Expected59

<!-- @createExpected59 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/prometheus/rbac.authorization.k8s.io_v1_role_prometheus-server.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  labels:
    app.kubernetes.io/name: prometheus
    app.kubernetes.io/part-of: ingress-nginx
  name: prometheus-server
  namespace: ingress-nginx
rules:
- apiGroups:
  - ""
  resources:
  - services
  - endpoints
  - pods
  verbs:
  - get
  - list
  - watch
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

