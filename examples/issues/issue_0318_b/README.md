# Feature Test for Issue 0318


This folder contains files describing how to address [Issue 0318](https://github.com/kubernetes-sigs/kustomize/issues/0318)

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
mkdir -p ${DEMO_HOME}/overlay
mkdir -p ${DEMO_HOME}/overlay/production
mkdir -p ${DEMO_HOME}/overlay/staging
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
commonAnnotations:
  service.beta.kubernetes.io/aws-load-balancer-ssl-cert: $(Values.public.spec.AWS_LOAD_BALANCER_SSL_CERT)
  service.beta.kubernetes.io/aws-load-balancer-extra-security-group: $(Values.public.spec.AWS_LOAD_BALANCER_EXTRA_SECURITY_GROUP)

resources:
- values.yaml
- resources.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/production/kustomization.yaml

resources:
- ../../base

patchesStrategicMerge:
- ./values.yaml

# current kustomize version calls for adding
# transformers here. It does not work because
# the filter applied here interferes with
# variables resolution. Such a filtering has
# to be apply after the variables are resolved,
#
# transformers:
# - kubectlapplyordertransformer.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/staging/kustomization.yaml

resources:
- ../../base

patchesStrategicMerge:
- ./values.yaml

# current kustomize version calls for adding
# transformers here. It does not work because
# the filter applied here interferes with
# variables resolution. Such a filtering has
# to be apply after the variables are resolved,
#
# transformers:
# - kubectlapplyordertransformer.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/resources.yaml
---
apiVersion: v1
kind: Service
metadata:
  name: ingress-validation-webhook
  namespace: ingress-nginx
  annotations:
      external-dns.alpha.kubernetes.io/hostname: $(Values.public.spec.EXTERNAL_DNS)
spec:
  ports:
  - name: admission
    port: 443
    protocol: TCP
    targetPort: 8080
  selector:
    app.kubernetes.io/name: ingress-nginx
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
      failure-domain.beta.kubernetes.io/zone: $(Values.public.spec.IAM_ZONE)
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
          image: ingress-nginx:latest
          args:
            - /nginx-ingress-controller
            - --configmap=$(POD_NAMESPACE)/nginx-configuration
            - --tcp-services-configmap=$(POD_NAMESPACE)/tcp-services
            - --udp-services-configmap=$(POD_NAMESPACE)/udp-services
            - --publish-service=$(POD_NAMESPACE)/ingress-nginx
            - --annotations-prefix=nginx.ingress.kubernetes.io
            - --validating-webhook=:8080
            - --validating-webhook-certificate=/usr/local/certificates/certificate.pem
            - --validating-webhook-key=/usr/local/certificates/key.pem
          volumeMounts:
          - name: webhook-cert
            mountPath: "/usr/local/certificates/"
            readOnly: true
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
            - name: webhook
              containerPort: 8080
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
      volumes:
      - name: webhook-cert
        secret:
          secretName: $(Values.public.spec.SECRET_NAME)
        awsElasiticBlockStore:
          volumeID: $(Values.public.spec.VOLUME_ID)
---
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: public
spec:
  AWS_LOAD_BALANCER_EXTRA_SECURITY_GROUP: aws.base.load.balancer.extra.security.group
  AWS_LOAD_BALANCER_SSL_CERT: aws.base.load.balancer.ssl.cert
  EXTERNAL_DNS: external.base.dns
  IAM_ZONE: iam-zone.base.example.com
  SECRET_NAME: base.secret-name
  VOLUME_ID: base-volumneid

EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/production/kubectlapplyordertransformer.yaml
apiVersion: builtin
kind: KindFilterTransformer
metadata:
  name: kindfiltertransformer
excludes:
- Values
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/production/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: public
spec:
  SECRET_NAME: production.secret-name
  VOLUME_ID: production-volumneid
  EXTERNAL_DNS: external.production.dns

EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/staging/kubectlapplyordertransformer.yaml
apiVersion: builtin
kind: KindFilterTransformer
metadata:
  name: kindfiltertransformer
excludes:
- Values
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/staging/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: public
spec:
  AWS_LOAD_BALANCER_SSL_CERT: aws.staging.load.balancer.ssl.cert
  AWS_LOAD_BALANCER_EXTRA_SECURITY_GROUP: aws.staging.load.balancer.extra.security.group
  IAM_ZONE: iam-zone.staging.example.com

EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlay/staging --reorder=kubectlapply -o ${DEMO_HOME}/actual/staging.yaml
kustomize build ${DEMO_HOME}/overlay/production --reorder=kubectlapply -o ${DEMO_HOME}/actual/production.yaml
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
kind: Service
metadata:
  annotations:
    external-dns.alpha.kubernetes.io/hostname: external.production.dns
    service.beta.kubernetes.io/aws-load-balancer-extra-security-group: aws.base.load.balancer.extra.security.group
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: aws.base.load.balancer.ssl.cert
  name: ingress-validation-webhook
  namespace: ingress-nginx
spec:
  ports:
  - name: admission
    port: 443
    protocol: TCP
    targetPort: 8080
  selector:
    app.kubernetes.io/name: ingress-nginx
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-extra-security-group: aws.base.load.balancer.extra.security.group
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: aws.base.load.balancer.ssl.cert
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
      failure-domain.beta.kubernetes.io/zone: iam-zone.base.example.com
  template:
    metadata:
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
        service.beta.kubernetes.io/aws-load-balancer-extra-security-group: aws.base.load.balancer.extra.security.group
        service.beta.kubernetes.io/aws-load-balancer-ssl-cert: aws.base.load.balancer.ssl.cert
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
        - --validating-webhook=:8080
        - --validating-webhook-certificate=/usr/local/certificates/certificate.pem
        - --validating-webhook-key=/usr/local/certificates/key.pem
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: ingress-nginx:latest
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
        - containerPort: 8080
          name: webhook
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
        volumeMounts:
        - mountPath: /usr/local/certificates/
          name: webhook-cert
          readOnly: true
      serviceAccountName: nginx-ingress-serviceaccount
      volumes:
      - awsElasiticBlockStore:
          volumeID: production-volumneid
        name: webhook-cert
        secret:
          secretName: production.secret-name
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/staging.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    external-dns.alpha.kubernetes.io/hostname: external.base.dns
    service.beta.kubernetes.io/aws-load-balancer-extra-security-group: aws.staging.load.balancer.extra.security.group
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: aws.staging.load.balancer.ssl.cert
  name: ingress-validation-webhook
  namespace: ingress-nginx
spec:
  ports:
  - name: admission
    port: 443
    protocol: TCP
    targetPort: 8080
  selector:
    app.kubernetes.io/name: ingress-nginx
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-extra-security-group: aws.staging.load.balancer.extra.security.group
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: aws.staging.load.balancer.ssl.cert
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
      failure-domain.beta.kubernetes.io/zone: iam-zone.staging.example.com
  template:
    metadata:
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
        service.beta.kubernetes.io/aws-load-balancer-extra-security-group: aws.staging.load.balancer.extra.security.group
        service.beta.kubernetes.io/aws-load-balancer-ssl-cert: aws.staging.load.balancer.ssl.cert
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
        - --validating-webhook=:8080
        - --validating-webhook-certificate=/usr/local/certificates/certificate.pem
        - --validating-webhook-key=/usr/local/certificates/key.pem
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: ingress-nginx:latest
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
        - containerPort: 8080
          name: webhook
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
        volumeMounts:
        - mountPath: /usr/local/certificates/
          name: webhook-cert
          readOnly: true
      serviceAccountName: nginx-ingress-serviceaccount
      volumes:
      - awsElasiticBlockStore:
          volumeID: base-volumneid
        name: webhook-cert
        secret:
          secretName: base.secret-name
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

