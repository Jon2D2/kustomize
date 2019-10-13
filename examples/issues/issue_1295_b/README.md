# Feature Test for Issue 1295


This folder contains files describing how to address [Issue 1295](https://github.com/kubernetes-sigs/kustomize/issues/1295)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/front
mkdir -p ${DEMO_HOME}/monitoring
mkdir -p ${DEMO_HOME}/monitoring/base
mkdir -p ${DEMO_HOME}/monitoring/GCS
mkdir -p ${DEMO_HOME}/monitoring/S3
mkdir -p ${DEMO_HOME}/serving
mkdir -p ${DEMO_HOME}/serving/base
mkdir -p ${DEMO_HOME}/serving/GCS
mkdir -p ${DEMO_HOME}/serving/local
mkdir -p ${DEMO_HOME}/training
mkdir -p ${DEMO_HOME}/training/base
mkdir -p ${DEMO_HOME}/training/GCS
mkdir -p ${DEMO_HOME}/training/local
mkdir -p ${DEMO_HOME}/training/S3
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/front/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generatorOptions:
  disableNameSuffixHash: true

resources:
- deployment.yaml
- service.yaml

namespace: kubeflow
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/monitoring/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- deployment.yaml
- service.yaml

namespace: kubeflow

generatorOptions:
  disableNameSuffixHash: true

configurations:
- params.yaml

vars:
- fieldref:
    fieldPath: data.logDir
  name: logDir
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/monitoring/GCS/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generatorOptions:
  disableNameSuffixHash: true

configurations:
- params.yaml

vars:
- fieldref:
    fieldPath: data.GOOGLE_APPLICATION_CREDENTIALS
  name: GOOGLE_APPLICATION_CREDENTIALS
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
- fieldref:
    fieldPath: data.secretName
  name: secretName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
- fieldref:
    fieldPath: data.secretMountPath
  name: secretMountPath
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring

patchesJson6902:
- path: deployment_patch.yaml
  target:
    group: apps
    kind: Deployment
    name: tensorboard-tb
    namespace: kubeflow
    version: v1beta1

resources:
- ../base

configMapGenerator:
- literals:
  - name=mnist-gcs-dist
  - exportDir=gs://my-bucket/my-model/export
  - logDir=/tmp
  - secretName=user-gcp-sa
  - secretMountPath=/var/secrets
  - GOOGLE_APPLICATION_CREDENTIALS=/var/secrets/user-gcp-sa.json
  name: mnist-map-monitoring
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/monitoring/S3/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generatorOptions:
  disableNameSuffixHash: true

configurations:
- params.yaml

vars:
- fieldref:
    fieldPath: data.S3_ENDPOINT
  name: S3_ENDPOINT
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
- fieldref:
    fieldPath: data.AWS_ENDPOINT_URL
  name: AWS_ENDPOINT_URL
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
- fieldref:
    fieldPath: data.AWS_REGION
  name: AWS_REGION
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
- fieldref:
    fieldPath: data.BUCKET_NAME
  name: BUCKET_NAME
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
- fieldref:
    fieldPath: data.S3_USE_HTTPS
  name: S3_USE_HTTPS
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
- fieldref:
    fieldPath: data.S3_VERIFY_SSL
  name: S3_VERIFY_SSL
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
- fieldref:
    fieldPath: data.awsSecretName
  name: awsSecretName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
- fieldref:
    fieldPath: data.awsAccessKeyIDName
  name: awsAccessKeyIDName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring
- fieldref:
    fieldPath: data.awsSecretAccessKeyName
  name: awsSecretAccessKeyName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-monitoring

patchesJson6902:
- path: deployment_patch.yaml
  target:
    group: apps
    kind: Deployment
    name: tensorboard-tb
    namespace: kubeflow
    version: v1beta1

resources:
- ../base

configMapGenerator:
- literals:
  - logDir=/tmp
  - awsSecretName=aws-creds
  - awsAccessKeyIDName=awsAccessKeyID
  - awsSecretAccessKeyName=awsSecretAccessKey
  - S3_ENDPOINT=s3.us-west-2.amazonaws.com
  - AWS_ENDPOINT_URL=https://s3.us-west-2.amazonaws.com
  - AWS_REGION=us-west-2
  - BUCKET_NAME=mybucket
  - S3_USE_HTTPS=1
  - S3_VERIFY_SSL=1
  name: mnist-map-monitoring
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/serving/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- deployment.yaml
- mnist-deploy-config.yaml
- service.yaml

namespace: kubeflow

generatorOptions:
  disableNameSuffixHash: true

configurations:
- params.yaml

vars:
- fieldref:
    fieldPath: data.name
  name: svcName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-serving
- fieldref:
    fieldPath: data.modelBasePath
  name: modelBasePath
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-serving
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/serving/GCS/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generatorOptions:
  disableNameSuffixHash: true

patchesJson6902:
- path: deployment_patch.yaml
  target:
    group: apps
    kind: Deployment
    name: $(svcName)
    namespace: kubeflow
    version: v1

resources:
- ../base

configMapGenerator:
- literals:
  - name=mnist-gcs-dist
  - modelBasePath=//export-dir
  name: mnist-map-serving
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/serving/local/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generatorOptions:
  disableNameSuffixHash: true

vars:
- fieldref:
    fieldPath: data.pvcName
  name: pvcName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-serving
- fieldref:
    fieldPath: data.pvcMountPath
  name: pvcMountPath
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-serving

configurations:
- params.yaml

patchesJson6902:
- path: deployment_patch.yaml
  target:
    group: apps
    kind: Deployment
    name: $(svcName)
    namespace: kubeflow
    version: v1

resources:
- ../base

configMapGenerator:
- literals:
  - name=mnist-service-local
  - pvcName=local
  - pvcMountPath=/mnt
  - modelBasePath=/mnt/export
  name: mnist-map-serving
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- Chief.yaml

namespace: kubeflow

generatorOptions:
  disableNameSuffixHash: true

configurations:
- params.yaml

vars:
- fieldref:
    fieldPath: data.name
  name: trainingName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.modelDir
  name: modelDir
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.exportDir
  name: exportDir
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.trainSteps
  name: trainSteps
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.batchSize
  name: batchSize
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.learningRate
  name: learningRate
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
patchesStrategicMerge:
- Ps.yaml
- Worker.yaml
EOF
```


### Preparation Step KustomizationFile8

<!-- @createKustomizationFile8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/GCS/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generatorOptions:
  disableNameSuffixHash: true

configurations:
- params.yaml

# TBD (jinchihe) Need move the image to base file once.
# the issue addressed: kubernetes-sigs/kustomize/issues/1040
# TBD (jinchihe) Need to update the image once
# the issue addressed: kubeflow/testing/issues/373
images:
- name: training-image

vars:
- fieldref:
    fieldPath: data.GOOGLE_APPLICATION_CREDENTIALS
  name: GOOGLE_APPLICATION_CREDENTIALS
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.secretName
  name: secretName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.secretMountPath
  name: secretMountPath
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
      
patchesJson6902:
- path: Worker_patch.yaml
  target:
    group: kubeflow.org
    kind: TFJob
    name: $(trainingName)
    namespace: kubeflow
    version: v1beta2
- path: Ps_patch.yaml
  target:
    group: kubeflow.org
    kind: TFJob
    name: $(trainingName)
    namespace: kubeflow
    version: v1beta2
- path: Chief_patch.yaml
  target:
    group: kubeflow.org
    kind: TFJob
    name: $(trainingName)
    namespace: kubeflow
    version: v1beta2

resources:
- ../base

configMapGenerator:
- literals:
  - name=mnist-train-dist
  - trainSteps=200
  - batchSize=100
  - learningRate=0.01
  - modelDir=gs://my-bucket/my-model
  - exportDir=gs://my-bucket/my-model/export
  - secretName=user-gcp-sa
  - secretMountPath=/var/secrets
  - GOOGLE_APPLICATION_CREDENTIALS=/var/secrets/user-gcp-sa.json
  name: mnist-map-training
EOF
```


### Preparation Step KustomizationFile9

<!-- @createKustomizationFile9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/local/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generatorOptions:
  disableNameSuffixHash: true

configurations:
- params.yaml

# TBD (jinchihe) Need move the image to base file once.
# the issue addressed: kubernetes-sigs/kustomize/issues/1040
# TBD (jinchihe) Need to update the image once
# the issue addressed: kubeflow/testing/issues/373
images:
- name: training-image
  newName: gcr.io/kubeflow-examples/mnist/model
  newTag: v20190111-v0.2-148-g313770f

vars:
- fieldref:
    fieldPath: data.pvcName
  name: pvcName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.pvcMountPath
  name: pvcMountPath
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
      
patchesJson6902:
- path: Worker_patch.yaml
  target:
    group: kubeflow.org
    kind: TFJob
    name: $(trainingName)
    namespace: kubeflow
    version: v1beta2
- path: Ps_patch.yaml
  target:
    group: kubeflow.org
    kind: TFJob
    name: $(trainingName)
    namespace: kubeflow
    version: v1beta2
- path: Chief_patch.yaml
  target:
    group: kubeflow.org
    kind: TFJob
    name: $(trainingName)
    namespace: kubeflow
    version: v1beta2

resources:
- ../base

configMapGenerator:
- literals:
  - name=mnist-train-local
  - trainSteps=200
  - batchSize=100
  - learningRate=0.01
  - pvcName=local
  - pvcMountPath=/mnt
  - modelDir=/mnt
  - exportDir=/mnt/export
  name: mnist-map-training
EOF
```


### Preparation Step KustomizationFile10

<!-- @createKustomizationFile10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/S3/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generatorOptions:
  disableNameSuffixHash: true

configurations:
- params.yaml

# TBD (jinchihe) Need move the image to base file once.
# the issue addressed: kubernetes-sigs/kustomize/issues/1040
# TBD (jinchihe) Need to update the image once
# the issue addressed: kubeflow/testing/issues/373
images:
- name: training-image
  newName: gcr.io/kubeflow-examples/mnist/model
  newTag: v20190111-v0.2-148-g313770f

vars:
- fieldref:
    fieldPath: data.S3_ENDPOINT
  name: S3_ENDPOINT
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.AWS_ENDPOINT_URL
  name: AWS_ENDPOINT_URL
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.AWS_REGION
  name: AWS_REGION
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.BUCKET_NAME
  name: BUCKET_NAME
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.S3_USE_HTTPS
  name: S3_USE_HTTPS
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.S3_VERIFY_SSL
  name: S3_VERIFY_SSL
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.awsSecretName
  name: awsSecretName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.awsAccessKeyIDName
  name: awsAccessKeyIDName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
- fieldref:
    fieldPath: data.awsSecretAccessKeyName
  name: awsSecretAccessKeyName
  objref:
    apiVersion: v1
    kind: ConfigMap
    name: mnist-map-training
      
patchesJson6902:
- path: Worker_patch.yaml
  target:
    group: kubeflow.org
    kind: TFJob
    name: $(trainingName)
    namespace: kubeflow
    version: v1beta2
- path: Ps_patch.yaml
  target:
    group: kubeflow.org
    kind: TFJob
    name: $(trainingName)
    namespace: kubeflow
    version: v1beta2
- path: Chief_patch.yaml
  target:
    group: kubeflow.org
    kind: TFJob
    name: $(trainingName)
    namespace: kubeflow
    version: v1beta2

resources:
- ../base

configMapGenerator:
- literals:
  - name=mnist-train-dist
  - trainSteps=200
  - batchSize=100
  - learningRate=0.01
  - modelDir=s3://path
  - exportDir=s3://export
  - awsSecretName=aws-creds
  - awsAccessKeyIDName=awsAccessKeyID
  - awsSecretAccessKeyName=awsSecretAccessKey
  - S3_ENDPOINT=s3.us-west-2.amazonaws.com
  - AWS_ENDPOINT_URL=https://s3.us-west-2.amazonaws.com
  - AWS_REGION=us-west-2
  - BUCKET_NAME=mybucket
  - S3_USE_HTTPS=1
  - S3_VERIFY_SSL=1
  name: mnist-map-training

secretGenerator:
- literals:
  - awsAccessKeyID=xxxxx
  - awsSecretAccessKey=xxxxx
  name: aws-creds
  namespace: kubeflow
  type: Opaque
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/front/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: web-ui
  namespace: kubeflow
spec:
  replicas: 1
  selector:
    matchLabels:
      app: web-ui
  template:
    metadata:
      labels:
        app: web-ui
    spec:
      containers:
      - image: gcr.io/kubeflow-examples/mnist/web-ui:v20190112-v0.2-142-g3b38225
        name: web-ui
        ports:
        - containerPort: 5000
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/front/service.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    getambassador.io/config: |-
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: web-ui_mapping
      prefix: /kubeflow/mnist/
      rewrite: /
      service: web-ui.kubeflow
  name: web-ui
  namespace: kubeflow
spec:
  ports:
  - port: 80
    targetPort: 5000
  selector:
    app: web-ui
  type: ClusterIP
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/monitoring/base/deployment.yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: tensorboard-tb
  namespace: kubeflow
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: tensorboard
        tb-job: tensorboard
      name: tensorboard
      namespace: kubeflow
    spec:
      containers:
      - command:
        - /usr/local/bin/tensorboard
        - --logdir=$(logDir)
        - --port=80
        env:
        - name: logDir
          value: $(logDir)
        image: tensorflow/tensorflow:1.11.0
        name: tensorboard
        ports:
        - containerPort: 80
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/monitoring/base/params.yaml
varReference:
- path: spec/template/spec/containers/env/value
  kind: Deployment
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/monitoring/base/service.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    getambassador.io/config: |-
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tensorboard_mapping
      prefix: /kubeflow/tensorboard/mnist
      rewrite: /
      service: tensorboard-tb.kubeflow
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tensorboard_mapping_data
      prefix: /kubeflow/tensorboard/mnist/data/
      rewrite: /data/
      service: tensorboard-tb.kubeflow
  name: tensorboard-tb
  namespace: kubeflow
spec:
  ports:
  - name: http
    port: 80
    targetPort: 80
  selector:
    app: tensorboard
    tb-job: tensorboard
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/monitoring/GCS/deployment_patch.yaml
- op: add
  path: /spec/template/spec/containers/0/volumeMounts
  value:
    - mountPath: $(secretMountPath)
      name: user-gcp-sa
      readOnly: true
- op: add
  path: /spec/template/spec/volumes
  value:
    - name: user-gcp-sa
      secret:
        secretName: $(secretName)
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: GOOGLE_APPLICATION_CREDENTIALS
    value: $(GOOGLE_APPLICATION_CREDENTIALS)
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/monitoring/GCS/params.yaml
varReference:
- path: spec/template/spec/volumes/secret/secretName
  kind: Deployment
- path: spec/template/spec/containers/volumeMounts/mountPath
  kind: Deployment
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/monitoring/S3/deployment_patch.yaml
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: S3_ENDPOINT
    value: $(S3_ENDPOINT)
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: AWS_ENDPOINT_URL
    value: $(AWS_ENDPOINT_URL)
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: AWS_REGION
    value: $(AWS_REGION)
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: BUCKET_NAME
    value: $(BUCKET_NAME)
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: S3_USE_HTTPS
    value: $(S3_USE_HTTPS)
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: S3_VERIFY_SSL
    value: $(S3_VERIFY_SSL)
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: AWS_ACCESS_KEY_ID
    valueFrom:
      secretKeyRef:
        key: $(awsAccessKeyIDName)
        name: $(awsSecretName)
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: AWS_SECRET_ACCESS_KEY
    valueFrom:
      secretKeyRef:
        key: $(awsSecretAccessKeyName)
        name: $(awsSecretName)
EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/monitoring/S3/params.yaml
varReference:
- path: spec/template/spec/containers/env/valueFrom/secretKeyRef/name
  kind: Deployment
- path: spec/template/spec/containers/env/valueFrom/secretKeyRef/key
  kind: Deployment
EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/serving/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: mnist
  name: $(svcName)
  namespace: kubeflow
spec:
  template:
    metadata:
      labels:
        app: mnist
        version: v1
    spec:
      containers:
      - args:
        - --port=9000
        - --rest_api_port=8500
        - --model_name=mnist
        - --model_base_path=$(modelBasePath)
        - --monitoring_config_file=/var/config/monitoring_config.txt
        command:
        - /usr/bin/tensorflow_model_server
        env:
        - name: modelBasePath
          value: $(modelBasePath)        
        image: tensorflow/serving:1.11.1
        imagePullPolicy: IfNotPresent
        livenessProbe:
          initialDelaySeconds: 30
          periodSeconds: 30
          tcpSocket:
            port: 9000
        name: mnist
        ports:
        - containerPort: 9000
        - containerPort: 8500
        resources:
          limits:
            cpu: "4"
            memory: 4Gi
          requests:
            cpu: "1"
            memory: 1Gi
        volumeMounts:
        - mountPath: /var/config/
          name: config-volume
      volumes:
      - configMap:
          name: mnist-deploy-config
        name: config-volume
EOF
```


### Preparation Step Resource10

<!-- @createResource10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/serving/base/mnist-deploy-config.yaml
apiVersion: v1
data:
  monitoring_config.txt: |-
    prometheus_config: {
      enable: true,
      path: "/monitoring/prometheus/metrics"
    }
kind: ConfigMap
metadata:
  name: mnist-deploy-config
  namespace: kubeflow
EOF
```


### Preparation Step Resource11

<!-- @createResource11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/serving/base/params.yaml
varReference:
- path: spec/template/spec/containers/env/value
  kind: Deployment
- path: metadata/name
  kind: Service
- path: metadata/name
  kind: Deployment
EOF
```


### Preparation Step Resource12

<!-- @createResource12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/serving/base/service.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    getambassador.io/config: |-
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tfserving-predict-mapping-mnist
      prefix: /tfserving/models/mnist
      rewrite: /v1/models/mnist:predict
      method: POST
      service: mnist-service.kubeflow:8500
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tfserving-predict-mapping-mnist-get
      prefix: /tfserving/models/mnist
      rewrite: /v1/models/mnist
      method: GET
      service: mnist-service.kubeflow:8500
    prometheus.io/path: /monitoring/prometheus/metrics
    prometheus.io/port: "8500"
    prometheus.io/scrape: "true"
  labels:
    app: mnist
  name: $(svcName)
  namespace: kubeflow
spec:
  ports:
  - name: grpc-tf-serving
    port: 9000
    targetPort: 9000
  - name: http-tf-serving
    port: 8500
    targetPort: 8500
  selector:
    app: mnist
  type: ClusterIP
EOF
```


### Preparation Step Resource13

<!-- @createResource13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/serving/GCS/deployment_patch.yaml
- op: add
  path: /spec/template/spec/containers/0/volumeMounts/-
  value:
    mountPath: /secret/gcp-credentials
    name: user-gcp-sa
    readOnly: true
- op: add
  path: /spec/template/spec/volumes/-
  value:
    name: user-gcp-sa
    secret:
      secretName: user-gcp-sa
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: GOOGLE_APPLICATION_CREDENTIALS
    value: /secret/gcp-credentials/user-gcp-sa.json
EOF
```


### Preparation Step Resource14

<!-- @createResource14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/serving/local/deployment_patch.yaml
- op: add
  path: /spec/template/spec/containers/0/volumeMounts/-
  value:
    mountPath: $(pvcMountPath)
    name: local-storage

- op: add
  path: /spec/template/spec/volumes/-
  value:
    name: local-storage
    persistentVolumeClaim:
      claimName: $(pvcName)
EOF
```


### Preparation Step Resource15

<!-- @createResource15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/serving/local/params.yaml
varReference:
- path: spec/template/spec/volumes/persistentVolumeClaim/claimName
  kind: Deployment
- path: spec/template/spec/containers/volumeMounts/mountPath
  kind: TFJob
EOF
```


### Preparation Step Resource16

<!-- @createResource16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/base/Chief.yaml
apiVersion: kubeflow.org/v1beta2
kind: TFJob
metadata:
  name: $(trainingName)
  namespace: kubeflow
spec:
  tfReplicaSpecs:
    Chief:
      replicas: 1
      template:
        spec:
          containers:
          - name: tensorflow
            command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: $(modelDir)
            - name: exportDir
              value: $(exportDir)
            - name: trainSteps
              value: $(trainSteps)
            - name: batchSize
              value: $(batchSize)
            - name: learningRate
              value: $(learningRate)
            image: training-image
            workingDir: /opt
          restartPolicy: OnFailure
EOF
```


### Preparation Step Resource17

<!-- @createResource17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/base/params.yaml
varReference:
- path: metadata/name
  kind: TFJob
- path: spec/tfReplicaSpecs/Chief/template/spec/containers/env/value
  kind: TFJob
- path: spec/tfReplicaSpecs/Worker/template/spec/containers/env/value
  kind: TFJob
- path: spec/tfReplicaSpecs/Ps/template/spec/containers/env/value
  kind: TFJob
EOF
```


### Preparation Step Resource18

<!-- @createResource18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/base/Ps.yaml
apiVersion: kubeflow.org/v1beta2
kind: TFJob
metadata:
  name: $(trainingName)
  namespace: kubeflow
spec:
  tfReplicaSpecs:
    Ps:
      replicas: 1
      template:
        spec:
          containers:
          - name: tensorflow
            command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: $(modelDir)
            - name: exportDir
              value: $(exportDir)
            - name: trainSteps
              value: $(trainSteps)
            - name: batchSize
              value: $(batchSize)
            - name: learningRate
              value: $(learningRate)
            image: training-image
            workingDir: /opt
          restartPolicy: OnFailure
EOF
```


### Preparation Step Resource19

<!-- @createResource19 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/base/Worker.yaml
apiVersion: kubeflow.org/v1beta2
kind: TFJob
metadata:
  name: $(trainingName)
  namespace: kubeflow
spec:
  tfReplicaSpecs:
    Worker:
      replicas: 2
      template:
        spec:
          containers:
          - name: tensorflow
            command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)            
            env:
            - name: modelDir
              value: $(modelDir)
            - name: exportDir
              value: $(exportDir)
            - name: trainSteps
              value: $(trainSteps)
            - name: batchSize
              value: $(batchSize)
            - name: learningRate
              value: $(learningRate)
            image: training-image
            name: tensorflow
            workingDir: /opt
          restartPolicy: OnFailure
EOF
```


### Preparation Step Resource20

<!-- @createResource20 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/GCS/Chief_patch.yaml
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/volumeMounts
  value:
    - mountPath: $(secretMountPath)
      name: user-gcp-sa
      readOnly: true
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/volumes
  value:
    - name: user-gcp-sa
      secret:
        secretName: $(secretName)
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/env/-
  value:
    name: GOOGLE_APPLICATION_CREDENTIALS
    value: $(GOOGLE_APPLICATION_CREDENTIALS)
EOF
```


### Preparation Step Resource21

<!-- @createResource21 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/GCS/params.yaml
varReference:
- path: metadata/name
  kind: TFJob
- path: spec/tfReplicaSpecs/Chief/template/spec/volumes/secret/secretName
  kind: TFJob
- path: spec/tfReplicaSpecs/Chief/template/spec/containers/volumeMounts/mountPath
  kind: TFJob
- path: spec/tfReplicaSpecs/Worker/template/spec/volumes/secret/secretName
  kind: TFJob 
- path: spec/tfReplicaSpecs/Worker/template/spec/containers/volumeMounts/mountPath
  kind: TFJob
- path: spec/tfReplicaSpecs/Ps/template/spec/volumes/secret/secretName
  kind: TFJob
- path: spec/tfReplicaSpecs/Ps/template/spec/containers/volumeMounts/mountPath
  kind: TFJob
EOF
```


### Preparation Step Resource22

<!-- @createResource22 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/GCS/Ps_patch.yaml
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/volumeMounts
  value:
    - mountPath: $(secretMountPath)
      name: user-gcp-sa
      readOnly: true
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/volumes
  value:
    - name: user-gcp-sa
      secret:
        secretName: $(secretName)
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/env/-
  value:
    name: GOOGLE_APPLICATION_CREDENTIALS
    value: $(GOOGLE_APPLICATION_CREDENTIALS)
EOF
```


### Preparation Step Resource23

<!-- @createResource23 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/GCS/Worker_patch.yaml
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/volumeMounts
  value:
    - mountPath: $(secretMountPath)
      name: user-gcp-sa
      readOnly: true
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/volumes
  value:
    - name: user-gcp-sa
      secret:
        secretName: $(secretName)
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/env/-
  value:
    name: GOOGLE_APPLICATION_CREDENTIALS
    value: $(GOOGLE_APPLICATION_CREDENTIALS)
EOF
```


### Preparation Step Resource24

<!-- @createResource24 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/local/Chief_patch.yaml
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/volumeMounts
  value:
    - mountPath: $(pvcMountPath)
      name: local-storage
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/volumes
  value:
    - name: local-storage
      persistentVolumeClaim:
        claimName: $(pvcName)
EOF
```


### Preparation Step Resource25

<!-- @createResource25 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/local/params.yaml
varReference:
- path: metadata/name
  kind: TFJob
- path: spec/tfReplicaSpecs/Chief/template/spec/volumes/persistentVolumeClaim/claimName
  kind: TFJob
- path: spec/tfReplicaSpecs/Worker/template/spec/volumes/persistentVolumeClaim/claimName
  kind: TFJob
- path: spec/tfReplicaSpecs/Ps/template/spec/volumes/persistentVolumeClaim/claimName
  kind: TFJob
- path: spec/tfReplicaSpecs/Chief/template/spec/containers/volumeMounts/mountPath
  kind: TFJob
- path: spec/tfReplicaSpecs/Worker/template/spec/containers/volumeMounts/mountPath
  kind: TFJob
- path: spec/tfReplicaSpecs/Ps/template/spec/containers/volumeMounts/mountPath
  kind: TFJob
EOF
```


### Preparation Step Resource26

<!-- @createResource26 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/local/Ps_patch.yaml
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/volumeMounts
  value:
    - mountPath: $(pvcMountPath)
      name: local-storage
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/volumes
  value:
    - name: local-storage
      persistentVolumeClaim:
        claimName: $(pvcName)
EOF
```


### Preparation Step Resource27

<!-- @createResource27 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/local/Worker_patch.yaml
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/volumeMounts
  value:
    - mountPath: $(pvcMountPath)
      name: local-storage
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/volumes
  value:
    - name: local-storage
      persistentVolumeClaim:
        claimName: $(pvcName)
EOF
```


### Preparation Step Resource28

<!-- @createResource28 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/S3/Chief_patch.yaml
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/env/-
  value:
    name: S3_ENDPOINT
    value: $(S3_ENDPOINT)
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/env/-
  value:
    name: AWS_ENDPOINT_URL
    value: $(AWS_ENDPOINT_URL)
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/env/-
  value:
    name: AWS_REGION
    value: $(AWS_REGION)
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/env/-
  value:
    name: BUCKET_NAME
    value: $(BUCKET_NAME)
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/env/-
  value:
    name: S3_USE_HTTPS
    value: $(S3_USE_HTTPS)
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/env/-
  value:
    name: S3_VERIFY_SSL
    value: $(S3_VERIFY_SSL)
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/env/-
  value:
    name: AWS_ACCESS_KEY_ID
    valueFrom:
      secretKeyRef:
        key: $(awsAccessKeyIDName)
        name: $(awsSecretName)
- op: add
  path: /spec/tfReplicaSpecs/Chief/template/spec/containers/0/env/-
  value:
    name: AWS_SECRET_ACCESS_KEY
    valueFrom:
      secretKeyRef:
        key: $(awsSecretAccessKeyName)
        name: $(awsSecretName)
EOF
```


### Preparation Step Resource29

<!-- @createResource29 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/S3/params.yaml
varReference:
- path: metadata/name
  kind: TFJob
- path: metadata/name
  kind: TFJob
- path: spec/tfReplicaSpecs/Chief/template/spec/containers/env/valueFrom/secretKeyRef/name
  kind: TFJob
- path: spec/tfReplicaSpecs/Chief/template/spec/containers/env/valueFrom/secretKeyRef/key
  kind: TFJob
- path: spec/tfReplicaSpecs/Ps/template/spec/containers/env/valueFrom/secretKeyRef/name
  kind: TFJob
- path: spec/tfReplicaSpecs/Ps/template/spec/containers/env/valueFrom/secretKeyRef/key
  kind: TFJob
- path: spec/tfReplicaSpecs/Worker/template/spec/containers/env/valueFrom/secretKeyRef/name
  kind: TFJob
- path: spec/tfReplicaSpecs/Worker/template/spec/containers/env/valueFrom/secretKeyRef/key
  kind: TFJob
EOF
```


### Preparation Step Resource30

<!-- @createResource30 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/S3/Ps_patch.yaml
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/env/-
  value:
    name: S3_ENDPOINT
    value: $(S3_ENDPOINT)
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/env/-
  value:
    name: AWS_ENDPOINT_URL
    value: $(AWS_ENDPOINT_URL)
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/env/-
  value:
    name: AWS_REGION
    value: $(AWS_REGION)
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/env/-
  value:
    name: BUCKET_NAME
    value: $(BUCKET_NAME)
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/env/-
  value:
    name: S3_USE_HTTPS
    value: $(S3_USE_HTTPS)
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/env/-
  value:
    name: S3_VERIFY_SSL
    value: $(S3_VERIFY_SSL)
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/env/-
  value:
    name: AWS_ACCESS_KEY_ID
    valueFrom:
      secretKeyRef:
        key: $(awsAccessKeyIDName)
        name: $(awsSecretName)
- op: add
  path: /spec/tfReplicaSpecs/Ps/template/spec/containers/0/env/-
  value:
    name: AWS_SECRET_ACCESS_KEY
    valueFrom:
      secretKeyRef:
        key: $(awsSecretAccessKeyName)
        name: $(awsSecretName)
EOF
```


### Preparation Step Resource31

<!-- @createResource31 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/training/S3/Worker_patch.yaml
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/env/-
  value:
    name: S3_ENDPOINT
    value: $(S3_ENDPOINT)
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/env/-
  value:
    name: AWS_ENDPOINT_URL
    value: $(AWS_ENDPOINT_URL)
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/env/-
  value:
    name: AWS_REGION
    value: $(AWS_REGION)
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/env/-
  value:
    name: BUCKET_NAME
    value: $(BUCKET_NAME)
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/env/-
  value:
    name: S3_USE_HTTPS
    value: $(S3_USE_HTTPS)
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/env/-
  value:
    name: S3_VERIFY_SSL
    value: $(S3_VERIFY_SSL)
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/env/-
  value:
    name: AWS_ACCESS_KEY_ID
    valueFrom:
      secretKeyRef:
        key: $(awsAccessKeyIDName)
        name: $(awsSecretName)
- op: add
  path: /spec/tfReplicaSpecs/Worker/template/spec/containers/0/env/-
  value:
    name: AWS_SECRET_ACCESS_KEY
    valueFrom:
      secretKeyRef:
        key: $(awsSecretAccessKeyName)
        name: $(awsSecretName)
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir -p ${DEMO_HOME}/actual
mkdir -p ${DEMO_HOME}/actual/training/local
mkdir -p ${DEMO_HOME}/actual/training/GCS
mkdir -p ${DEMO_HOME}/actual/training/S3
mkdir -p ${DEMO_HOME}/actual/serving/local
mkdir -p ${DEMO_HOME}/actual/serving/GCS
mkdir -p ${DEMO_HOME}/actual/front
mkdir -p ${DEMO_HOME}/actual/monitoring/GCS
mkdir -p ${DEMO_HOME}/actual/monitoring/S3

kustomize build ${DEMO_HOME}/training/local -o ${DEMO_HOME}/actual/training/local
kustomize build ${DEMO_HOME}/training/GCS -o ${DEMO_HOME}/actual/training/GCS
kustomize build ${DEMO_HOME}/training/S3 -o ${DEMO_HOME}/actual/training/S3
kustomize build ${DEMO_HOME}/serving/local -o ${DEMO_HOME}/actual/serving/local
kustomize build ${DEMO_HOME}/serving/GCS -o ${DEMO_HOME}/actual/serving/GCS
kustomize build ${DEMO_HOME}/front -o ${DEMO_HOME}/actual/front
kustomize build ${DEMO_HOME}/monitoring/GCS -o ${DEMO_HOME}/actual/monitoring/GCS
kustomize build ${DEMO_HOME}/monitoring/S3 -o ${DEMO_HOME}/actual/monitoring/S3
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir -p ${DEMO_HOME}/expected
mkdir -p ${DEMO_HOME}/expected/training/local
mkdir -p ${DEMO_HOME}/expected/training/GCS
mkdir -p ${DEMO_HOME}/expected/training/S3
mkdir -p ${DEMO_HOME}/expected/serving/local
mkdir -p ${DEMO_HOME}/expected/serving/GCS
mkdir -p ${DEMO_HOME}/expected/front
mkdir -p ${DEMO_HOME}/expected/monitoring/GCS
mkdir -p ${DEMO_HOME}/expected/monitoring/S3
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/front/apps_v1beta2_deployment_web-ui.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: web-ui
  namespace: kubeflow
spec:
  replicas: 1
  selector:
    matchLabels:
      app: web-ui
  template:
    metadata:
      labels:
        app: web-ui
    spec:
      containers:
      - image: gcr.io/kubeflow-examples/mnist/web-ui:v20190112-v0.2-142-g3b38225
        name: web-ui
        ports:
        - containerPort: 5000
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/front/~g_v1_service_web-ui.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    getambassador.io/config: |-
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: web-ui_mapping
      prefix: /kubeflow/mnist/
      rewrite: /
      service: web-ui.kubeflow
  name: web-ui
  namespace: kubeflow
spec:
  ports:
  - port: 80
    targetPort: 5000
  selector:
    app: web-ui
  type: ClusterIP
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/monitoring/GCS/default_~g_v1_configmap_mnist-map-monitoring.yaml
apiVersion: v1
data:
  GOOGLE_APPLICATION_CREDENTIALS: /var/secrets/user-gcp-sa.json
  exportDir: gs://my-bucket/my-model/export
  logDir: /tmp
  name: mnist-gcs-dist
  secretMountPath: /var/secrets
  secretName: user-gcp-sa
kind: ConfigMap
metadata:
  name: mnist-map-monitoring
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/monitoring/GCS/kubeflow_apps_v1beta1_deployment_tensorboard-tb.yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: tensorboard-tb
  namespace: kubeflow
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: tensorboard
        tb-job: tensorboard
      name: tensorboard
      namespace: kubeflow
    spec:
      containers:
      - command:
        - /usr/local/bin/tensorboard
        - --logdir=/tmp
        - --port=80
        env:
        - name: logDir
          value: /tmp
        - name: GOOGLE_APPLICATION_CREDENTIALS
          value: /var/secrets/user-gcp-sa.json
        image: tensorflow/tensorflow:1.11.0
        name: tensorboard
        ports:
        - containerPort: 80
        volumeMounts:
        - mountPath: /var/secrets
          name: user-gcp-sa
          readOnly: true
      volumes:
      - name: user-gcp-sa
        secret:
          secretName: user-gcp-sa
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/monitoring/GCS/kubeflow_~g_v1_service_tensorboard-tb.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    getambassador.io/config: |-
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tensorboard_mapping
      prefix: /kubeflow/tensorboard/mnist
      rewrite: /
      service: tensorboard-tb.kubeflow
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tensorboard_mapping_data
      prefix: /kubeflow/tensorboard/mnist/data/
      rewrite: /data/
      service: tensorboard-tb.kubeflow
  name: tensorboard-tb
  namespace: kubeflow
spec:
  ports:
  - name: http
    port: 80
    targetPort: 80
  selector:
    app: tensorboard
    tb-job: tensorboard
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/monitoring/S3/default_~g_v1_configmap_mnist-map-monitoring.yaml
apiVersion: v1
data:
  AWS_ENDPOINT_URL: https://s3.us-west-2.amazonaws.com
  AWS_REGION: us-west-2
  BUCKET_NAME: mybucket
  S3_ENDPOINT: s3.us-west-2.amazonaws.com
  S3_USE_HTTPS: "1"
  S3_VERIFY_SSL: "1"
  awsAccessKeyIDName: awsAccessKeyID
  awsSecretAccessKeyName: awsSecretAccessKey
  awsSecretName: aws-creds
  logDir: /tmp
kind: ConfigMap
metadata:
  name: mnist-map-monitoring
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/monitoring/S3/kubeflow_apps_v1beta1_deployment_tensorboard-tb.yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: tensorboard-tb
  namespace: kubeflow
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: tensorboard
        tb-job: tensorboard
      name: tensorboard
      namespace: kubeflow
    spec:
      containers:
      - command:
        - /usr/local/bin/tensorboard
        - --logdir=/tmp
        - --port=80
        env:
        - name: logDir
          value: /tmp
        - name: S3_ENDPOINT
          value: s3.us-west-2.amazonaws.com
        - name: AWS_ENDPOINT_URL
          value: https://s3.us-west-2.amazonaws.com
        - name: AWS_REGION
          value: us-west-2
        - name: BUCKET_NAME
          value: mybucket
        - name: S3_USE_HTTPS
          value: "1"
        - name: S3_VERIFY_SSL
          value: "1"
        - name: AWS_ACCESS_KEY_ID
          valueFrom:
            secretKeyRef:
              key: awsAccessKeyID
              name: aws-creds
        - name: AWS_SECRET_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              key: awsSecretAccessKey
              name: aws-creds
        image: tensorflow/tensorflow:1.11.0
        name: tensorboard
        ports:
        - containerPort: 80
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/monitoring/S3/kubeflow_~g_v1_service_tensorboard-tb.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    getambassador.io/config: |-
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tensorboard_mapping
      prefix: /kubeflow/tensorboard/mnist
      rewrite: /
      service: tensorboard-tb.kubeflow
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tensorboard_mapping_data
      prefix: /kubeflow/tensorboard/mnist/data/
      rewrite: /data/
      service: tensorboard-tb.kubeflow
  name: tensorboard-tb
  namespace: kubeflow
spec:
  ports:
  - name: http
    port: 80
    targetPort: 80
  selector:
    app: tensorboard
    tb-job: tensorboard
EOF
```


### Verification Step Expected8

<!-- @createExpected8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/serving/GCS/default_~g_v1_configmap_mnist-map-serving.yaml
apiVersion: v1
data:
  modelBasePath: //export-dir
  name: mnist-gcs-dist
kind: ConfigMap
metadata:
  name: mnist-map-serving
EOF
```


### Verification Step Expected9

<!-- @createExpected9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/serving/GCS/kubeflow_apps_v1_deployment_mnist-gcs-dist.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: mnist
  name: mnist-gcs-dist
  namespace: kubeflow
spec:
  template:
    metadata:
      labels:
        app: mnist
        version: v1
    spec:
      containers:
      - args:
        - --port=9000
        - --rest_api_port=8500
        - --model_name=mnist
        - --model_base_path=//export-dir
        - --monitoring_config_file=/var/config/monitoring_config.txt
        command:
        - /usr/bin/tensorflow_model_server
        env:
        - name: modelBasePath
          value: //export-dir
        - name: GOOGLE_APPLICATION_CREDENTIALS
          value: /secret/gcp-credentials/user-gcp-sa.json
        image: tensorflow/serving:1.11.1
        imagePullPolicy: IfNotPresent
        livenessProbe:
          initialDelaySeconds: 30
          periodSeconds: 30
          tcpSocket:
            port: 9000
        name: mnist
        ports:
        - containerPort: 9000
        - containerPort: 8500
        resources:
          limits:
            cpu: "4"
            memory: 4Gi
          requests:
            cpu: "1"
            memory: 1Gi
        volumeMounts:
        - mountPath: /var/config/
          name: config-volume
        - mountPath: /secret/gcp-credentials
          name: user-gcp-sa
          readOnly: true
      volumes:
      - configMap:
          name: mnist-deploy-config
        name: config-volume
      - name: user-gcp-sa
        secret:
          secretName: user-gcp-sa
EOF
```


### Verification Step Expected10

<!-- @createExpected10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/serving/GCS/kubeflow_~g_v1_configmap_mnist-deploy-config.yaml
apiVersion: v1
data:
  monitoring_config.txt: |-
    prometheus_config: {
      enable: true,
      path: "/monitoring/prometheus/metrics"
    }
kind: ConfigMap
metadata:
  name: mnist-deploy-config
  namespace: kubeflow
EOF
```


### Verification Step Expected11

<!-- @createExpected11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/serving/GCS/kubeflow_~g_v1_service_mnist-gcs-dist.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    getambassador.io/config: |-
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tfserving-predict-mapping-mnist
      prefix: /tfserving/models/mnist
      rewrite: /v1/models/mnist:predict
      method: POST
      service: mnist-service.kubeflow:8500
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tfserving-predict-mapping-mnist-get
      prefix: /tfserving/models/mnist
      rewrite: /v1/models/mnist
      method: GET
      service: mnist-service.kubeflow:8500
    prometheus.io/path: /monitoring/prometheus/metrics
    prometheus.io/port: "8500"
    prometheus.io/scrape: "true"
  labels:
    app: mnist
  name: mnist-gcs-dist
  namespace: kubeflow
spec:
  ports:
  - name: grpc-tf-serving
    port: 9000
    targetPort: 9000
  - name: http-tf-serving
    port: 8500
    targetPort: 8500
  selector:
    app: mnist
  type: ClusterIP
EOF
```


### Verification Step Expected12

<!-- @createExpected12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/serving/local/default_~g_v1_configmap_mnist-map-serving.yaml
apiVersion: v1
data:
  modelBasePath: /mnt/export
  name: mnist-service-local
  pvcMountPath: /mnt
  pvcName: local
kind: ConfigMap
metadata:
  name: mnist-map-serving
EOF
```


### Verification Step Expected13

<!-- @createExpected13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/serving/local/kubeflow_apps_v1_deployment_mnist-service-local.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: mnist
  name: mnist-service-local
  namespace: kubeflow
spec:
  template:
    metadata:
      labels:
        app: mnist
        version: v1
    spec:
      containers:
      - args:
        - --port=9000
        - --rest_api_port=8500
        - --model_name=mnist
        - --model_base_path=/mnt/export
        - --monitoring_config_file=/var/config/monitoring_config.txt
        command:
        - /usr/bin/tensorflow_model_server
        env:
        - name: modelBasePath
          value: /mnt/export
        image: tensorflow/serving:1.11.1
        imagePullPolicy: IfNotPresent
        livenessProbe:
          initialDelaySeconds: 30
          periodSeconds: 30
          tcpSocket:
            port: 9000
        name: mnist
        ports:
        - containerPort: 9000
        - containerPort: 8500
        resources:
          limits:
            cpu: "4"
            memory: 4Gi
          requests:
            cpu: "1"
            memory: 1Gi
        volumeMounts:
        - mountPath: /var/config/
          name: config-volume
        - mountPath: /mnt
          name: local-storage
      volumes:
      - configMap:
          name: mnist-deploy-config
        name: config-volume
      - name: local-storage
        persistentVolumeClaim:
          claimName: local
EOF
```


### Verification Step Expected14

<!-- @createExpected14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/serving/local/kubeflow_~g_v1_configmap_mnist-deploy-config.yaml
apiVersion: v1
data:
  monitoring_config.txt: |-
    prometheus_config: {
      enable: true,
      path: "/monitoring/prometheus/metrics"
    }
kind: ConfigMap
metadata:
  name: mnist-deploy-config
  namespace: kubeflow
EOF
```


### Verification Step Expected15

<!-- @createExpected15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/serving/local/kubeflow_~g_v1_service_mnist-service-local.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    getambassador.io/config: |-
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tfserving-predict-mapping-mnist
      prefix: /tfserving/models/mnist
      rewrite: /v1/models/mnist:predict
      method: POST
      service: mnist-service.kubeflow:8500
      ---
      apiVersion: ambassador/v0
      kind:  Mapping
      name: tfserving-predict-mapping-mnist-get
      prefix: /tfserving/models/mnist
      rewrite: /v1/models/mnist
      method: GET
      service: mnist-service.kubeflow:8500
    prometheus.io/path: /monitoring/prometheus/metrics
    prometheus.io/port: "8500"
    prometheus.io/scrape: "true"
  labels:
    app: mnist
  name: mnist-service-local
  namespace: kubeflow
spec:
  ports:
  - name: grpc-tf-serving
    port: 9000
    targetPort: 9000
  - name: http-tf-serving
    port: 8500
    targetPort: 8500
  selector:
    app: mnist
  type: ClusterIP
EOF
```


### Verification Step Expected16

<!-- @createExpected16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/training/GCS/default_~g_v1_configmap_mnist-map-training.yaml
apiVersion: v1
data:
  GOOGLE_APPLICATION_CREDENTIALS: /var/secrets/user-gcp-sa.json
  batchSize: "100"
  exportDir: gs://my-bucket/my-model/export
  learningRate: "0.01"
  modelDir: gs://my-bucket/my-model
  name: mnist-train-dist
  secretMountPath: /var/secrets
  secretName: user-gcp-sa
  trainSteps: "200"
kind: ConfigMap
metadata:
  name: mnist-map-training
EOF
```


### Verification Step Expected17

<!-- @createExpected17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/training/GCS/kubeflow_kubeflow.org_v1beta2_tfjob_mnist-train-dist.yaml
apiVersion: kubeflow.org/v1beta2
kind: TFJob
metadata:
  name: mnist-train-dist
  namespace: kubeflow
spec:
  tfReplicaSpecs:
    Chief:
      replicas: 1
      template:
        spec:
          containers:
          - command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: gs://my-bucket/my-model
            - name: exportDir
              value: gs://my-bucket/my-model/export
            - name: trainSteps
              value: "200"
            - name: batchSize
              value: "100"
            - name: learningRate
              value: "0.01"
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /var/secrets/user-gcp-sa.json
            image: training-image
            name: tensorflow
            volumeMounts:
            - mountPath: /var/secrets
              name: user-gcp-sa
              readOnly: true
            workingDir: /opt
          restartPolicy: OnFailure
          volumes:
          - name: user-gcp-sa
            secret:
              secretName: user-gcp-sa
    Ps:
      replicas: 1
      template:
        spec:
          containers:
          - command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: gs://my-bucket/my-model
            - name: exportDir
              value: gs://my-bucket/my-model/export
            - name: trainSteps
              value: "200"
            - name: batchSize
              value: "100"
            - name: learningRate
              value: "0.01"
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /var/secrets/user-gcp-sa.json
            image: training-image
            name: tensorflow
            volumeMounts:
            - mountPath: /var/secrets
              name: user-gcp-sa
              readOnly: true
            workingDir: /opt
          restartPolicy: OnFailure
          volumes:
          - name: user-gcp-sa
            secret:
              secretName: user-gcp-sa
    Worker:
      replicas: 2
      template:
        spec:
          containers:
          - command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: gs://my-bucket/my-model
            - name: exportDir
              value: gs://my-bucket/my-model/export
            - name: trainSteps
              value: "200"
            - name: batchSize
              value: "100"
            - name: learningRate
              value: "0.01"
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /var/secrets/user-gcp-sa.json
            image: training-image
            name: tensorflow
            volumeMounts:
            - mountPath: /var/secrets
              name: user-gcp-sa
              readOnly: true
            workingDir: /opt
          restartPolicy: OnFailure
          volumes:
          - name: user-gcp-sa
            secret:
              secretName: user-gcp-sa
EOF
```


### Verification Step Expected18

<!-- @createExpected18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/training/local/default_~g_v1_configmap_mnist-map-training.yaml
apiVersion: v1
data:
  batchSize: "100"
  exportDir: /mnt/export
  learningRate: "0.01"
  modelDir: /mnt
  name: mnist-train-local
  pvcMountPath: /mnt
  pvcName: local
  trainSteps: "200"
kind: ConfigMap
metadata:
  name: mnist-map-training
EOF
```


### Verification Step Expected19

<!-- @createExpected19 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/training/local/kubeflow_kubeflow.org_v1beta2_tfjob_mnist-train-local.yaml
apiVersion: kubeflow.org/v1beta2
kind: TFJob
metadata:
  name: mnist-train-local
  namespace: kubeflow
spec:
  tfReplicaSpecs:
    Chief:
      replicas: 1
      template:
        spec:
          containers:
          - command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: /mnt
            - name: exportDir
              value: /mnt/export
            - name: trainSteps
              value: "200"
            - name: batchSize
              value: "100"
            - name: learningRate
              value: "0.01"
            image: gcr.io/kubeflow-examples/mnist/model:v20190111-v0.2-148-g313770f
            name: tensorflow
            volumeMounts:
            - mountPath: /mnt
              name: local-storage
            workingDir: /opt
          restartPolicy: OnFailure
          volumes:
          - name: local-storage
            persistentVolumeClaim:
              claimName: local
    Ps:
      replicas: 1
      template:
        spec:
          containers:
          - command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: /mnt
            - name: exportDir
              value: /mnt/export
            - name: trainSteps
              value: "200"
            - name: batchSize
              value: "100"
            - name: learningRate
              value: "0.01"
            image: gcr.io/kubeflow-examples/mnist/model:v20190111-v0.2-148-g313770f
            name: tensorflow
            volumeMounts:
            - mountPath: /mnt
              name: local-storage
            workingDir: /opt
          restartPolicy: OnFailure
          volumes:
          - name: local-storage
            persistentVolumeClaim:
              claimName: local
    Worker:
      replicas: 2
      template:
        spec:
          containers:
          - command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: /mnt
            - name: exportDir
              value: /mnt/export
            - name: trainSteps
              value: "200"
            - name: batchSize
              value: "100"
            - name: learningRate
              value: "0.01"
            image: gcr.io/kubeflow-examples/mnist/model:v20190111-v0.2-148-g313770f
            name: tensorflow
            volumeMounts:
            - mountPath: /mnt
              name: local-storage
            workingDir: /opt
          restartPolicy: OnFailure
          volumes:
          - name: local-storage
            persistentVolumeClaim:
              claimName: local
EOF
```


### Verification Step Expected20

<!-- @createExpected20 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/training/S3/default_~g_v1_configmap_mnist-map-training.yaml
apiVersion: v1
data:
  AWS_ENDPOINT_URL: https://s3.us-west-2.amazonaws.com
  AWS_REGION: us-west-2
  BUCKET_NAME: mybucket
  S3_ENDPOINT: s3.us-west-2.amazonaws.com
  S3_USE_HTTPS: "1"
  S3_VERIFY_SSL: "1"
  awsAccessKeyIDName: awsAccessKeyID
  awsSecretAccessKeyName: awsSecretAccessKey
  awsSecretName: aws-creds
  batchSize: "100"
  exportDir: s3://export
  learningRate: "0.01"
  modelDir: s3://path
  name: mnist-train-dist
  trainSteps: "200"
kind: ConfigMap
metadata:
  name: mnist-map-training
EOF
```


### Verification Step Expected21

<!-- @createExpected21 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/training/S3/kubeflow_~g_v1_secret_aws-creds.yaml
apiVersion: v1
data:
  awsAccessKeyID: eHh4eHg=
  awsSecretAccessKey: eHh4eHg=
kind: Secret
metadata:
  name: aws-creds
  namespace: kubeflow
type: Opaque
EOF
```


### Verification Step Expected22

<!-- @createExpected22 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/training/S3/kubeflow_kubeflow.org_v1beta2_tfjob_mnist-train-dist.yaml
apiVersion: kubeflow.org/v1beta2
kind: TFJob
metadata:
  name: mnist-train-dist
  namespace: kubeflow
spec:
  tfReplicaSpecs:
    Chief:
      replicas: 1
      template:
        spec:
          containers:
          - command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: s3://path
            - name: exportDir
              value: s3://export
            - name: trainSteps
              value: "200"
            - name: batchSize
              value: "100"
            - name: learningRate
              value: "0.01"
            - name: S3_ENDPOINT
              value: s3.us-west-2.amazonaws.com
            - name: AWS_ENDPOINT_URL
              value: https://s3.us-west-2.amazonaws.com
            - name: AWS_REGION
              value: us-west-2
            - name: BUCKET_NAME
              value: mybucket
            - name: S3_USE_HTTPS
              value: "1"
            - name: S3_VERIFY_SSL
              value: "1"
            - name: AWS_ACCESS_KEY_ID
              valueFrom:
                secretKeyRef:
                  key: awsAccessKeyID
                  name: aws-creds
            - name: AWS_SECRET_ACCESS_KEY
              valueFrom:
                secretKeyRef:
                  key: awsSecretAccessKey
                  name: aws-creds
            image: gcr.io/kubeflow-examples/mnist/model:v20190111-v0.2-148-g313770f
            name: tensorflow
            workingDir: /opt
          restartPolicy: OnFailure
    Ps:
      replicas: 1
      template:
        spec:
          containers:
          - command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: s3://path
            - name: exportDir
              value: s3://export
            - name: trainSteps
              value: "200"
            - name: batchSize
              value: "100"
            - name: learningRate
              value: "0.01"
            - name: S3_ENDPOINT
              value: s3.us-west-2.amazonaws.com
            - name: AWS_ENDPOINT_URL
              value: https://s3.us-west-2.amazonaws.com
            - name: AWS_REGION
              value: us-west-2
            - name: BUCKET_NAME
              value: mybucket
            - name: S3_USE_HTTPS
              value: "1"
            - name: S3_VERIFY_SSL
              value: "1"
            - name: AWS_ACCESS_KEY_ID
              valueFrom:
                secretKeyRef:
                  key: awsAccessKeyID
                  name: aws-creds
            - name: AWS_SECRET_ACCESS_KEY
              valueFrom:
                secretKeyRef:
                  key: awsSecretAccessKey
                  name: aws-creds
            image: gcr.io/kubeflow-examples/mnist/model:v20190111-v0.2-148-g313770f
            name: tensorflow
            workingDir: /opt
          restartPolicy: OnFailure
    Worker:
      replicas: 2
      template:
        spec:
          containers:
          - command:
            - /usr/bin/python
            - /opt/model.py
            - --tf-model-dir=$(modelDir)
            - --tf-export-dir=$(exportDir)
            - --tf-train-steps=$(trainSteps)
            - --tf-batch-size=$(batchSize)
            - --tf-learning-rate=$(learningRate)
            env:
            - name: modelDir
              value: s3://path
            - name: exportDir
              value: s3://export
            - name: trainSteps
              value: "200"
            - name: batchSize
              value: "100"
            - name: learningRate
              value: "0.01"
            - name: S3_ENDPOINT
              value: s3.us-west-2.amazonaws.com
            - name: AWS_ENDPOINT_URL
              value: https://s3.us-west-2.amazonaws.com
            - name: AWS_REGION
              value: us-west-2
            - name: BUCKET_NAME
              value: mybucket
            - name: S3_USE_HTTPS
              value: "1"
            - name: S3_VERIFY_SSL
              value: "1"
            - name: AWS_ACCESS_KEY_ID
              valueFrom:
                secretKeyRef:
                  key: awsAccessKeyID
                  name: aws-creds
            - name: AWS_SECRET_ACCESS_KEY
              valueFrom:
                secretKeyRef:
                  key: awsSecretAccessKey
                  name: aws-creds
            image: gcr.io/kubeflow-examples/mnist/model:v20190111-v0.2-148-g313770f
            name: tensorflow
            workingDir: /opt
          restartPolicy: OnFailure
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

