# Kustomize Regression Test based on kfserving

This folder is only used for kustomize regression testing.
The original files are located [here](https://github.com/kubeflow/kfserving/tree/master/config)

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
mkdir -p ${DEMO_HOME}/default
mkdir -p ${DEMO_HOME}/default/configmap
mkdir -p ${DEMO_HOME}/default/crds
mkdir -p ${DEMO_HOME}/default/manager
mkdir -p ${DEMO_HOME}/default/rbac
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/development
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/crds/kustomization.yaml
# Labels to add to all resources and selectors.
#commonLabels:
#  someName: someValue

# Each entry in this list must resolve to an existing
# resource definition in YAML.  These are the resource
# files that kustomize reads, modifies and emits as a
# YAML string, with resources separated by document
# markers ("---").
resources:
- serving_v1alpha2_kfservice.yaml

patches:
  # Knative uses VolatileTime in place of metav1.Time to exclude this from creating equality.Semantic difference,
  # however when generating crd last transition time is tagged as object thus getting validation error while updating
  # status object(status.conditions.lastTransitionTime in body must be of type object: "string").
  # Controller-runtime 2.0 supports +kubebuilder:validation:Type=string but we are still on 1.9.0 due to knative dep.
- crd_status_condition_patch.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/kustomization.yaml
# Adds namespace to all resources.
namespace: kfserving-system

# Labels to add to all resources and selectors.
#commonLabels:
#  someName: someValue

# Each entry in this list must resolve to an existing
# resource definition in YAML.  These are the resource
# files that kustomize reads, modifies and emits as a
# YAML string, with resources separated by document
# markers ("---").
resources:
- crds/serving_v1alpha2_kfservice.yaml
- configmap/kfservice.yaml
- rbac/rbac_role.yaml
- rbac/rbac_role_binding.yaml
- manager/manager.yaml
- manager/service.yaml
  # Comment the following 3 lines if you want to disable
  # the auth proxy (https://github.com/brancz/kube-rbac-proxy)
  # which protects your /metrics endpoint.
- rbac/auth_proxy_service.yaml
- rbac/auth_proxy_role.yaml
- rbac/auth_proxy_role_binding.yaml

patchesStrategicMerge:
- manager_image_patch.yaml
  # Protect the /metrics endpoint by putting it behind auth.
  # Only one of manager_auth_proxy_patch.yaml and
  # manager_prometheus_metrics_patch.yaml should be enabled.
- manager_auth_proxy_patch.yaml
  # If you want your controller-manager to expose the /metrics
  # endpoint w/o any authn/z, uncomment the following line and
  # comment manager_auth_proxy_patch.yaml.
  # Only one of manager_auth_proxy_patch.yaml and
  # manager_prometheus_metrics_patch.yaml should be enabled.
#- manager_prometheus_metrics_patch.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/development/kustomization.yaml
bases:
  - ../../default

patchesStrategicMerge:
  - manager_image_patch.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/configmap/kfservice.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: kfservice-config
  namespace: kfserving-system
data:
  frameworks: |-
    {
        "tensorflow": {
            "image": "tensorflow/serving"
        },
        "sklearn": {
            "image": "gcr.io/kfserving/sklearnserver"
        },
        "xgboost": {
            "image": "gcr.io/kfserving/xgbserver"
        },
        "pytorch": {
            "image": "gcr.io/kfserving/pytorchserver"
        },
        "tensorrt": {
            "image": "nvcr.io/nvidia/tensorrtserver"
        }
    }
  modelInitializer: |-
    {
        "image" : "gcr.io/kfserving/model-initializer:latest"
    }
  credentials: |-
    {
       "gcs": {
           "gcsCredentialFileName": "gcloud-application-credentials.json"
       },
       "s3": {
           "s3AccessKeyIDName": "awsAccessKeyID",
           "s3SecretAccessKeyName": "awsSecretAccessKey"
       }
    }
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/crds/crd_status_condition_patch.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: kfservices.serving.kubeflow.org
spec:
  validation:
    openAPIV3Schema:
      properties:
        status:
          properties:
            url:
              type: string
            conditions:
              items:
                properties:
                  lastTransitionTime:
                    type: string
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/crds/serving_v1alpha2_kfservice.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: kfservices.serving.kubeflow.org
spec:
  additionalPrinterColumns:
  - JSONPath: .status.conditions[?(@.type=='Ready')].status
    name: Ready
    type: string
  - JSONPath: .status.url
    name: URL
    type: string
  - JSONPath: .status.default.traffic
    name: Default Traffic
    type: integer
  - JSONPath: .status.canary.traffic
    name: Canary Traffic
    type: integer
  - JSONPath: .metadata.creationTimestamp
    name: Age
    type: date
  group: serving.kubeflow.org
  names:
    kind: KFService
    plural: kfservices
    shortNames:
    - kfservice
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            canary:
              description: Canary defines an alternate endpoints to route a percentage
                of traffic.
              properties:
                explainer:
                  description: Explainer defines the model explanation service spec
                    explainer service calls to transformer or predictor service
                  properties:
                    alibi:
                      description: The following fields follow a "1-of" semantic.
                        Users must specify exactly one spec.
                      properties:
                        config:
                          description: Inline custom parameter settings for explainer
                          type: object
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest Alibi Version.
                          type: string
                        storageUri:
                          description: The location of a trained explanation model
                          type: string
                        type:
                          description: The type of Alibi explainer
                          type: string
                      required:
                      - type
                      type: object
                    custom:
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
                predictor:
                  description: Predictor defines the model serving spec +required
                  properties:
                    custom:
                      description: The following fields follow a "1-of" semantic.
                        Users must specify exactly one spec.
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    pytorch:
                      properties:
                        modelClassName:
                          description: Defaults PyTorch model class name to 'PyTorchModel'
                          type: string
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest PyTorch Version
                          type: string
                      required:
                      - modelUri
                      type: object
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                    sklearn:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest SKLearn Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    tensorflow:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest TF Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    tensorrt:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest TensorRT Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    xgboost:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest XGBoost Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                  type: object
                transformer:
                  description: Transformer defines the transformer service spec for
                    pre/post processing transformer service calls to predictor service
                  properties:
                    custom:
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
              required:
              - predictor
              type: object
            canaryTrafficPercent:
              description: CanaryTrafficPercent defines the percentage of traffic
                going to canary KFService endpoints
              format: int64
              type: integer
            default:
              description: Default defines default KFService endpoints +required
              properties:
                explainer:
                  description: Explainer defines the model explanation service spec
                    explainer service calls to transformer or predictor service
                  properties:
                    alibi:
                      description: The following fields follow a "1-of" semantic.
                        Users must specify exactly one spec.
                      properties:
                        config:
                          description: Inline custom parameter settings for explainer
                          type: object
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest Alibi Version.
                          type: string
                        storageUri:
                          description: The location of a trained explanation model
                          type: string
                        type:
                          description: The type of Alibi explainer
                          type: string
                      required:
                      - type
                      type: object
                    custom:
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
                predictor:
                  description: Predictor defines the model serving spec +required
                  properties:
                    custom:
                      description: The following fields follow a "1-of" semantic.
                        Users must specify exactly one spec.
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    pytorch:
                      properties:
                        modelClassName:
                          description: Defaults PyTorch model class name to 'PyTorchModel'
                          type: string
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest PyTorch Version
                          type: string
                      required:
                      - modelUri
                      type: object
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                    sklearn:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest SKLearn Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    tensorflow:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest TF Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    tensorrt:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest TensorRT Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    xgboost:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest XGBoost Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                  type: object
                transformer:
                  description: Transformer defines the transformer service spec for
                    pre/post processing transformer service calls to predictor service
                  properties:
                    custom:
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
              required:
              - predictor
              type: object
          required:
          - default
          type: object
        status:
          properties:
            canary:
              properties:
                name:
                  type: string
                replicas:
                  format: int64
                  type: integer
                traffic:
                  format: int64
                  type: integer
              type: object
            conditions:
              description: Conditions the latest available observations of a resource's
                current state. +patchMergeKey=type +patchStrategy=merge
              items:
                properties:
                  lastTransitionTime:
                    description: LastTransitionTime is the last time the condition
                      transitioned from one status to another. We use VolatileTime
                      in place of metav1.Time to exclude this from creating equality.Semantic
                      differences (all other things held constant).
                    type: string
                  message:
                    description: A human readable message indicating details about
                      the transition.
                    type: string
                  reason:
                    description: The reason for the condition's last transition.
                    type: string
                  severity:
                    description: Severity with which to treat failures of this type
                      of condition. When this is not specified, it defaults to Error.
                    type: string
                  status:
                    description: Status of the condition, one of True, False, Unknown.
                      +required
                    type: string
                  type:
                    description: Type of condition. +required
                    type: string
                required:
                - type
                - status
                type: object
              type: array
            default:
              properties:
                name:
                  type: string
                replicas:
                  format: int64
                  type: integer
                traffic:
                  format: int64
                  type: integer
              type: object
            observedGeneration:
              description: ObservedGeneration is the 'Generation' of the Service that
                was last processed by the controller.
              format: int64
              type: integer
            url:
              type: string
          type: object
  version: v1alpha2
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/manager_auth_proxy_patch.yaml
# This patch inject a sidecar container which is a HTTP proxy for the controller manager,
# it performs RBAC authorization against the Kubernetes API using SubjectAccessReviews.
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: kfserving-controller-manager
  namespace: kfserving-system
spec:
  template:
    spec:
      containers:
      - name: kube-rbac-proxy
        image: gcr.io/kubebuilder/kube-rbac-proxy:v0.4.0
        args:
        - "--secure-listen-address=0.0.0.0:8443"
        - "--upstream=http://127.0.0.1:8080/"
        - "--logtostderr=true"
        - "--v=10"
        ports:
        - containerPort: 8443
          name: https
      - name: manager
        args:
        - "--metrics-addr=127.0.0.1:8080"
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/manager_image_patch.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: kfserving-controller-manager
  namespace: kfserving-system
spec:
  template:
    spec:
      containers:
      # Change the value of image field below to your controller image URL
      - image: gcr.io/kfserving/kfserving-controller:latest
        name: manager
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/manager/manager.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: kfserving-controller-manager
  namespace: kfserving-system
  labels:
    control-plane: kfserving-controller-manager
    controller-tools.k8s.io: "1.0"
spec:
  selector:
    matchLabels:
      control-plane: kfserving-controller-manager
      controller-tools.k8s.io: "1.0"
  serviceName: controller-manager-service
  template:
    metadata:
      labels:
        control-plane: kfserving-controller-manager
        controller-tools.k8s.io: "1.0"
    spec:
      containers:
      - command:
        - /manager
        image: github.com/kubeflow/kfserving/cmd/manager
        imagePullPolicy: Always
        name: manager
        env:
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: SECRET_NAME
            value: kfserving-webhook-server-secret
        resources:
          limits:
            cpu: 100m
            memory: 300Mi
          requests:
            cpu: 100m
            memory: 200Mi
        ports:
        - containerPort: 9876
          name: webhook-server
          protocol: TCP
        volumeMounts:
        - mountPath: /tmp/cert
          name: cert
          readOnly: true
      terminationGracePeriodSeconds: 10
      volumes:
      - name: cert
        secret:
          defaultMode: 420
          secretName: kfserving-webhook-server-secret
---
apiVersion: v1
kind: Secret
metadata:
  name: kfserving-webhook-server-secret
  namespace: kfserving-system
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/manager_prometheus_metrics_patch.yaml
# This patch enables Prometheus scraping for the manager pod.
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: kfserving-controller-manager
  namespace: kfserving-system
spec:
  template:
    metadata:
      annotations:
        prometheus.io/scrape: 'true'
    spec:
      containers:
      # Expose the prometheus metrics on default port
      - name: manager
        ports:
        - containerPort: 8080
          name: metrics
          protocol: TCP
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/manager/service.yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: kfserving-controller-manager
    controller-tools.k8s.io: "1.0"
    istio-injection: disabled
  name: kfserving-system
---
apiVersion: v1
kind: Service
metadata:
  name: kfserving-controller-manager-service
  namespace: kfserving-system
  labels:
    control-plane: kfserving-controller-manager
    controller-tools.k8s.io: "1.0"
spec:
  selector:
    control-plane: kfserving-controller-manager
    controller-tools.k8s.io: "1.0"
  ports:
  - port: 443
EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/rbac/auth_proxy_role_binding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kfserving-proxy-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: proxy-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: system
EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/rbac/auth_proxy_role.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kfserving-proxy-role
rules:
- apiGroups: ["authentication.k8s.io"]
  resources:
  - tokenreviews
  verbs: ["create"]
- apiGroups: ["authorization.k8s.io"]
  resources:
  - subjectaccessreviews
  verbs: ["create"]
EOF
```


### Preparation Step Resource10

<!-- @createResource10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/rbac/auth_proxy_service.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    prometheus.io/port: "8443"
    prometheus.io/scheme: https
    prometheus.io/scrape: "true"
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
  name: kfserving-controller-manager-metrics-service
  namespace: kfserving-system
spec:
  ports:
  - name: https
    port: 8443
    targetPort: https
  selector:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
EOF
```


### Preparation Step Resource11

<!-- @createResource11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/rbac/rbac_role_binding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  creationTimestamp: null
  name: manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: manager-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: system
EOF
```


### Preparation Step Resource12

<!-- @createResource12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/default/rbac/rbac_role.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: manager-role
rules:
- apiGroups:
  - serving.knative.dev
  resources:
  - configurations
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - serving.knative.dev
  resources:
  - configurations/status
  verbs:
  - get
  - update
  - patch
- apiGroups:
  - serving.knative.dev
  resources:
  - routes
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - serving.knative.dev
  resources:
  - routes/status
  verbs:
  - get
  - update
  - patch
- apiGroups:
  - serving.kubeflow.org
  resources:
  - kfservices
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - serving.kubeflow.org
  resources:
  - kfservices/status
  verbs:
  - get
  - update
  - patch
- apiGroups:
  - ""
  resources:
  - serviceaccounts
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - admissionregistration.k8s.io
  resources:
  - mutatingwebhookconfigurations
  - validatingwebhookconfigurations
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
EOF
```


### Preparation Step Resource13

<!-- @createResource13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/development/manager_image_patch.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: kfserving-controller-manager
  namespace: kfserving-system
spec:
  template:
    spec:
      containers:
      # Change the value of image field below to your controller image URL
      - image: gcr.io/kfserving/kfserving-controller:v0.1.0
        name: manager
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/development -o ${DEMO_HOME}/actual/result.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/result.yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: kfserving-controller-manager
    controller-tools.k8s.io: "1.0"
    istio-injection: disabled
  name: kfserving-system
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: kfservices.serving.kubeflow.org
spec:
  additionalPrinterColumns:
  - JSONPath: .status.conditions[?(@.type=='Ready')].status
    name: Ready
    type: string
  - JSONPath: .status.url
    name: URL
    type: string
  - JSONPath: .status.default.traffic
    name: Default Traffic
    type: integer
  - JSONPath: .status.canary.traffic
    name: Canary Traffic
    type: integer
  - JSONPath: .metadata.creationTimestamp
    name: Age
    type: date
  group: serving.kubeflow.org
  names:
    kind: KFService
    plural: kfservices
    shortNames:
    - kfservice
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            canary:
              description: Canary defines an alternate endpoints to route a percentage
                of traffic.
              properties:
                explainer:
                  description: Explainer defines the model explanation service spec
                    explainer service calls to transformer or predictor service
                  properties:
                    alibi:
                      description: The following fields follow a "1-of" semantic.
                        Users must specify exactly one spec.
                      properties:
                        config:
                          description: Inline custom parameter settings for explainer
                          type: object
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest Alibi Version.
                          type: string
                        storageUri:
                          description: The location of a trained explanation model
                          type: string
                        type:
                          description: The type of Alibi explainer
                          type: string
                      required:
                      - type
                      type: object
                    custom:
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
                predictor:
                  description: Predictor defines the model serving spec +required
                  properties:
                    custom:
                      description: The following fields follow a "1-of" semantic.
                        Users must specify exactly one spec.
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    pytorch:
                      properties:
                        modelClassName:
                          description: Defaults PyTorch model class name to 'PyTorchModel'
                          type: string
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest PyTorch Version
                          type: string
                      required:
                      - modelUri
                      type: object
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                    sklearn:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest SKLearn Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    tensorflow:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest TF Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    tensorrt:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest TensorRT Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    xgboost:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest XGBoost Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                  type: object
                transformer:
                  description: Transformer defines the transformer service spec for
                    pre/post processing transformer service calls to predictor service
                  properties:
                    custom:
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
              required:
              - predictor
              type: object
            canaryTrafficPercent:
              description: CanaryTrafficPercent defines the percentage of traffic
                going to canary KFService endpoints
              format: int64
              type: integer
            default:
              description: Default defines default KFService endpoints +required
              properties:
                explainer:
                  description: Explainer defines the model explanation service spec
                    explainer service calls to transformer or predictor service
                  properties:
                    alibi:
                      description: The following fields follow a "1-of" semantic.
                        Users must specify exactly one spec.
                      properties:
                        config:
                          description: Inline custom parameter settings for explainer
                          type: object
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest Alibi Version.
                          type: string
                        storageUri:
                          description: The location of a trained explanation model
                          type: string
                        type:
                          description: The type of Alibi explainer
                          type: string
                      required:
                      - type
                      type: object
                    custom:
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
                predictor:
                  description: Predictor defines the model serving spec +required
                  properties:
                    custom:
                      description: The following fields follow a "1-of" semantic.
                        Users must specify exactly one spec.
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    pytorch:
                      properties:
                        modelClassName:
                          description: Defaults PyTorch model class name to 'PyTorchModel'
                          type: string
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest PyTorch Version
                          type: string
                      required:
                      - modelUri
                      type: object
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                    sklearn:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest SKLearn Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    tensorflow:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest TF Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    tensorrt:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest TensorRT Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                    xgboost:
                      properties:
                        modelUri:
                          type: string
                        resources:
                          description: Defaults to requests and limits of 1CPU, 2Gb
                            MEM.
                          type: object
                        runtimeVersion:
                          description: Defaults to latest XGBoost Version.
                          type: string
                      required:
                      - modelUri
                      type: object
                  type: object
                transformer:
                  description: Transformer defines the transformer service spec for
                    pre/post processing transformer service calls to predictor service
                  properties:
                    custom:
                      properties:
                        container:
                          type: object
                      required:
                      - container
                      type: object
                    maxReplicas:
                      description: This is the up bound for autoscaler to scale to
                      format: int64
                      type: integer
                    minReplicas:
                      description: Minimum number of replicas, pods won't scale down
                        to 0 in case of no traffic
                      format: int64
                      type: integer
                    serviceAccountName:
                      description: ServiceAccountName is the name of the ServiceAccount
                        to use to run the service
                      type: string
                  type: object
              required:
              - predictor
              type: object
          required:
          - default
          type: object
        status:
          properties:
            canary:
              properties:
                name:
                  type: string
                replicas:
                  format: int64
                  type: integer
                traffic:
                  format: int64
                  type: integer
              type: object
            conditions:
              description: Conditions the latest available observations of a resource's
                current state. +patchMergeKey=type +patchStrategy=merge
              items:
                properties:
                  lastTransitionTime:
                    description: LastTransitionTime is the last time the condition
                      transitioned from one status to another. We use VolatileTime
                      in place of metav1.Time to exclude this from creating equality.Semantic
                      differences (all other things held constant).
                    type: string
                  message:
                    description: A human readable message indicating details about
                      the transition.
                    type: string
                  reason:
                    description: The reason for the condition's last transition.
                    type: string
                  severity:
                    description: Severity with which to treat failures of this type
                      of condition. When this is not specified, it defaults to Error.
                    type: string
                  status:
                    description: Status of the condition, one of True, False, Unknown.
                      +required
                    type: string
                  type:
                    description: Type of condition. +required
                    type: string
                required:
                - type
                - status
                type: object
              type: array
            default:
              properties:
                name:
                  type: string
                replicas:
                  format: int64
                  type: integer
                traffic:
                  format: int64
                  type: integer
              type: object
            observedGeneration:
              description: ObservedGeneration is the 'Generation' of the Service that
                was last processed by the controller.
              format: int64
              type: integer
            url:
              type: string
          type: object
  version: v1alpha2
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kfserving-proxy-role
rules:
- apiGroups:
  - authentication.k8s.io
  resources:
  - tokenreviews
  verbs:
  - create
- apiGroups:
  - authorization.k8s.io
  resources:
  - subjectaccessreviews
  verbs:
  - create
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: manager-role
rules:
- apiGroups:
  - serving.knative.dev
  resources:
  - configurations
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - serving.knative.dev
  resources:
  - configurations/status
  verbs:
  - get
  - update
  - patch
- apiGroups:
  - serving.knative.dev
  resources:
  - routes
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - serving.knative.dev
  resources:
  - routes/status
  verbs:
  - get
  - update
  - patch
- apiGroups:
  - serving.kubeflow.org
  resources:
  - kfservices
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - serving.kubeflow.org
  resources:
  - kfservices/status
  verbs:
  - get
  - update
  - patch
- apiGroups:
  - ""
  resources:
  - serviceaccounts
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - admissionregistration.k8s.io
  resources:
  - mutatingwebhookconfigurations
  - validatingwebhookconfigurations
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kfserving-proxy-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: proxy-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: kfserving-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  creationTimestamp: null
  name: manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: manager-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: kfserving-system
---
apiVersion: v1
data:
  credentials: |-
    {
       "gcs": {
           "gcsCredentialFileName": "gcloud-application-credentials.json"
       },
       "s3": {
           "s3AccessKeyIDName": "awsAccessKeyID",
           "s3SecretAccessKeyName": "awsSecretAccessKey"
       }
    }
  frameworks: |-
    {
        "tensorflow": {
            "image": "tensorflow/serving"
        },
        "sklearn": {
            "image": "gcr.io/kfserving/sklearnserver"
        },
        "xgboost": {
            "image": "gcr.io/kfserving/xgbserver"
        },
        "pytorch": {
            "image": "gcr.io/kfserving/pytorchserver"
        },
        "tensorrt": {
            "image": "nvcr.io/nvidia/tensorrtserver"
        }
    }
  modelInitializer: |-
    {
        "image" : "gcr.io/kfserving/model-initializer:latest"
    }
kind: ConfigMap
metadata:
  name: kfservice-config
  namespace: kfserving-system
---
apiVersion: v1
kind: Secret
metadata:
  name: kfserving-webhook-server-secret
  namespace: kfserving-system
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    prometheus.io/port: "8443"
    prometheus.io/scheme: https
    prometheus.io/scrape: "true"
  labels:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
  name: kfserving-controller-manager-metrics-service
  namespace: kfserving-system
spec:
  ports:
  - name: https
    port: 8443
    targetPort: https
  selector:
    control-plane: controller-manager
    controller-tools.k8s.io: "1.0"
---
apiVersion: v1
kind: Service
metadata:
  labels:
    control-plane: kfserving-controller-manager
    controller-tools.k8s.io: "1.0"
  name: kfserving-controller-manager-service
  namespace: kfserving-system
spec:
  ports:
  - port: 443
  selector:
    control-plane: kfserving-controller-manager
    controller-tools.k8s.io: "1.0"
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    control-plane: kfserving-controller-manager
    controller-tools.k8s.io: "1.0"
  name: kfserving-controller-manager
  namespace: kfserving-system
spec:
  selector:
    matchLabels:
      control-plane: kfserving-controller-manager
      controller-tools.k8s.io: "1.0"
  serviceName: controller-manager-service
  template:
    metadata:
      labels:
        control-plane: kfserving-controller-manager
        controller-tools.k8s.io: "1.0"
    spec:
      containers:
      - args:
        - --secure-listen-address=0.0.0.0:8443
        - --upstream=http://127.0.0.1:8080/
        - --logtostderr=true
        - --v=10
        image: gcr.io/kubebuilder/kube-rbac-proxy:v0.4.0
        name: kube-rbac-proxy
        ports:
        - containerPort: 8443
          name: https
      - args:
        - --metrics-addr=127.0.0.1:8080
        command:
        - /manager
        env:
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: SECRET_NAME
          value: kfserving-webhook-server-secret
        image: gcr.io/kfserving/kfserving-controller:v0.1.0
        imagePullPolicy: Always
        name: manager
        ports:
        - containerPort: 9876
          name: webhook-server
          protocol: TCP
        resources:
          limits:
            cpu: 100m
            memory: 300Mi
          requests:
            cpu: 100m
            memory: 200Mi
        volumeMounts:
        - mountPath: /tmp/cert
          name: cert
          readOnly: true
      terminationGracePeriodSeconds: 10
      volumes:
      - name: cert
        secret:
          defaultMode: 420
          secretName: kfserving-webhook-server-secret
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

