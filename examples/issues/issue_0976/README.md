# Feature Test for Issue 0976


This folder contains files describing how to address [Issue 0976](https://github.com/kubernetes-sigs/kustomize/issues/0976)

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
mkdir -p ${DEMO_HOME}/secret
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
- admission_configurations.yaml

namePrefix: test-

commonLabels:
  foo: bar

configurations:
  - kustomizeconfig.yaml

secretGenerator:
- name: webhook-server-cert
  files:
  - tls.crt=secret/tls.cert
  - tls.key=secret/tls.key
  type: "kubernetes.io/tls"

vars:
  - name: TLSCERT
    objref:
      kind: Secret
      version: v1
      name: webhook-server-cert
    fieldref:
      fieldpath: data[tls.crt]

EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig.yaml
varReference:
  - path: webhooks/clientConfig/caBundle
    kind: ValidatingWebhookConfiguration
  - path: webhooks/clientConfig/caBundle
    kind: MutatingWebhookConfiguration

EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/admission_configurations.yaml
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: mutating-webhook-configuration
webhooks:
  - name: mutating-create-update
    clientConfig:
      url: https://example.com
      caBundle: $(TLSCERT)
    failurePolicy: Fail
    rules:
      - apiGroups:
          - mygroup
        apiVersions:
          - v1alpha1
        operations:
          - CREATE
          - UPDATE
        resources:
          - myresource
    sideEffects: None
---
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: validating-webhook-configuration
webhooks:
  - name: validating-create-update
    clientConfig:
      url: https://example.com
      caBundle: $(TLSCERT)
    failurePolicy: Fail
    rules:
      - apiGroups:
          - mygroup
        apiVersions:
          - v1alpha1
        operations:
          - CREATE
          - UPDATE
        resources:
          - myresource
    sideEffects: None

EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/secret/tls.cert
-----BEGIN CERTIFICATE-----
Li4u
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
Li4u
-----END CERTIFICATE-----

EOF
```


### Preparation Step Other1

<!-- @createOther1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/secret/tls.key
-----BEGIN RSA PRIVATE KEY-----
Li4u
-----END RSA PRIVATE KEY-----

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
cat <<'EOF' >${DEMO_HOME}/expected/admissionregistration.k8s.io_v1beta1_mutatingwebhookconfiguration_test-mutating-webhook-configuration.yaml
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  labels:
    foo: bar
  name: test-mutating-webhook-configuration
webhooks:
- clientConfig:
    caBundle: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCkxpNHUKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQotLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTGk0dQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    url: https://example.com
  failurePolicy: Fail
  name: mutating-create-update
  rules:
  - apiGroups:
    - mygroup
    apiVersions:
    - v1alpha1
    operations:
    - CREATE
    - UPDATE
    resources:
    - myresource
  sideEffects: None
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/admissionregistration.k8s.io_v1beta1_validatingwebhookconfiguration_test-validating-webhook-configuration.yaml
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  labels:
    foo: bar
  name: test-validating-webhook-configuration
webhooks:
- clientConfig:
    caBundle: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCkxpNHUKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQotLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTGk0dQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    url: https://example.com
  failurePolicy: Fail
  name: validating-create-update
  rules:
  - apiGroups:
    - mygroup
    apiVersions:
    - v1alpha1
    operations:
    - CREATE
    - UPDATE
    resources:
    - myresource
  sideEffects: None
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_secret_test-webhook-server-cert-454t54c54k.yaml
apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCkxpNHUKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQotLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTGk0dQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  tls.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpMaTR1Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
kind: Secret
metadata:
  labels:
    foo: bar
  name: test-webhook-server-cert-454t54c54k
type: kubernetes.io/tls
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

