# Feature Test for Issue allinone


This folder contains files describing how to address [Issue allinone](https://github.com/kubernetes-sigs/kustomize/issues/allinone)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/base
mkdir -p ${DEMO_HOME}/base/catalogues
mkdir -p ${DEMO_HOME}/base/crds
mkdir -p ${DEMO_HOME}/base/kustomizeconfig
mkdir -p ${DEMO_HOME}/base/mysql
mkdir -p ${DEMO_HOME}/base/wordpress
mkdir -p ${DEMO_HOME}/dev
mkdir -p ${DEMO_HOME}/production
mkdir -p ${DEMO_HOME}/production/common
mkdir -p ${DEMO_HOME}/production/site1
mkdir -p ${DEMO_HOME}/production/site2
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ./catalogues/versions.yaml
- ./catalogues/endpoints.yaml
- ./catalogues/common-addresses.yaml
- ./wordpress/deployment.yaml
- ./wordpress/service.yaml
- ./mysql/deployment.yaml
- ./mysql/secret.yaml
- ./mysql/service.yaml

configurations:
- kustomizeconfig/Chart.yaml
- kustomizeconfig/CommonAddresses.yaml
- kustomizeconfig/EndpointCatalogue.yaml
- kustomizeconfig/SoftwareVersions.yaml
- kustomizeconfig/Deployment.yaml
- kustomizeconfig/Service.yaml

vars:
- name: Service.wordpress.metadata.name
  objref:
    kind: Service
    name: wordpress
    apiVersion: v1
- name: Service.mysql.metadata.name
  objref:
    kind: Service
    name: mysql
    apiVersion: v1
- name: SoftwareVersions.software-versions.spec.images.wordpress.tag
  objref:
    apiVersion: my.group.org/v1alpha1
    kind: SoftwareVersions
    name: software-versions
  fieldref:
    fieldpath: spec.images.wordpress.tag
- name: SoftwareVersions.software-versions.spec.images.mysql.tag
  objref:
    apiVersion: my.group.org/v1alpha1
    kind: SoftwareVersions
    name: software-versions
  fieldref:
    fieldpath: spec.images.mysql.tag
# Demonstrate the ability to fetch specific index of an index
- name: CommonAddresses.common-addresses.spec.dns.upstream_servers[2]
  objref:
    apiVersion: my.group.org/v1alpha1
    kind: CommonAddresses
    name: common-addresses
  fieldref:
    fieldpath: spec.dns.upstream_servers[2]
# Demonstrate the ability to clone entire tree from on object to the other
- name: Deployment.wordpress.spec.template.spec.initContainers
  fieldref:
    fieldpath: spec.template.spec.initContainers
  objref:
    apiVersion: apps/v1beta2
    kind: Deployment
    name: wordpress
- name: EndpointCatalogue.endpoints.spec.wordpress.labels
  objref:
    apiVersion: my.group.org/v1alpha1
    kind: EndpointCatalogue
    name: endpoints
  fieldref:
    fieldpath: spec.wordpress.labels
- name: EndpointCatalogue.endpoints.spec.mysql.labels
  objref:
    apiVersion: my.group.org/v1alpha1
    kind: EndpointCatalogue
    name: endpoints
  fieldref:
    fieldpath: spec.mysql.labels
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/dev/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namePrefix: dev-

patchesStrategicMerge:
- passphrases.yaml
- versions.yaml
- endpoints.yaml
- common-addresses.yaml

resources:
- ../base
- devtools.yaml

vars:
- name: SoftwareVersions.software-versions.spec.charts.wordpress
  objref:
    apiVersion: my.group.org/v1alpha1
    kind: SoftwareVersions
    name: software-versions
  fieldref:
    fieldpath: spec.charts.wordpress
- name: SoftwareVersions.software-versions.spec.charts.mysql
  objref:
    apiVersion: my.group.org/v1alpha1
    kind: SoftwareVersions
    name: software-versions
  fieldref:
    fieldpath: spec.charts.mysql
- name: SoftwareVersions.software-versions.spec.images.wordpress
  objref:
    apiVersion: my.group.org/v1alpha1
    kind: SoftwareVersions
    name: software-versions
  fieldref:
    fieldpath: spec.images.wordpress
- name: SoftwareVersions.software-versions.spec.images.mysql
  objref:
    apiVersion: my.group.org/v1alpha1
    kind: SoftwareVersions
    name: software-versions
  fieldref:
    fieldpath: spec.images.mysql
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/common/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

patchesStrategicMerge:
- ./endpoints.yaml
- ./versions.yaml

resources:
- ../../base

vars:
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/site1/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

patchesStrategicMerge:
- passphrases.yaml
- common-addresses.yaml

resources:
- ../common

vars:
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/site2/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

patchesStrategicMerge:
- passphrases.yaml
- common-addresses.yaml

resources:
- ../common

vars:
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/catalogues/common-addresses.yaml
---
apiVersion: my.group.org/v1alpha1
kind: CommonAddresses
metadata:
  name: common-addresses
spec:
  dns:
    cluster_domain: cluster.local
    ingress_domain: global.my.group.org
    service_ip: 10.96.0.10
    upstream_servers:
    - 8.8.8.8
    - 9.9.9.9
    - 10.0.2.3
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/catalogues/endpoints.yaml
---
apiVersion: my.group.org/v1alpha1
kind: EndpointCatalogue
metadata:
  name: endpoints
spec:
  mysql:
    labels:
       app: mysql
  wordpress:
    labels:
       app: wordpress

EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/catalogues/versions.yaml
---
apiVersion: my.group.org/v1alpha1
kind: SoftwareVersions
metadata:
  name: software-versions
spec:
  charts:
    wordpress:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: wordpress
      type: git
    mysql:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: mysql
      type: git
  images:
    wordpress:
      registory: docker.io
      repository: wordpress
      tag: 4.8-apache
    mysql:
      registory: docker.io
      repository: mysql
      tag: '5.6'
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/crds/Chart.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: charts.my.group.org
spec:
  additionalPrinterColumns:
  - JSONPath: .status.actual_state
    description: State
    name: State
    type: string
  - JSONPath: .spec.target_state
    description: Target State
    name: Target State
    type: string
  - JSONPath: .status.satisfied
    description: Satisfied
    name: Satisfied
    type: boolean
  group: my.group.org
  names:
    kind: Chart
    plural: charts
    shortNames:
    - act
  scope: Namespaced
  subresources:
    status: {}
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
            source:
              description: provide a path to a ``git repo``, ``local dir``, or ``tarball
                url`` chart
              properties:
                location:
                  description: '``url`` or ``path`` to the chart''s parent directory'
                  type: string
                reference:
                  description: (optional) branch, commit, or reference in the repo
                    (``master`` if not specified)
                  type: string
                subpath:
                  description: (optional) relative path to target chart from parent
                    (``.`` if not specified)
                  type: string
                type:
                  description: 'source to build the chart: ``git``, ``local``, or
                    ``tar``'
                  type: string
              required:
              - location
              - subpath
              - type
              type: object
            values:
              description: override any default values in the charts
              properties:
                endpoints:
                  description: 'endpoints contains tbd. JEB: Would have been too consistent.
                    Different structures are used depending on the direction of the
                    wind. Endpoints *map[string]AVEndpoint `json:"endpoints,omitempty"`'
                  properties:
                    identity:
                      description: identity contains tbd
                      properties:
                        auth:
                          additionalProperties:
                            properties:
                              bind:
                                description: bind contains tbd
                                type: string
                              database:
                                description: database contains tbd
                                type: string
                              email:
                                description: email contains tbd
                                type: string
                              password:
                                description: password contains tbd
                                type: string
                              role:
                                description: role contains tbd
                                type: string
                              tls:
                                description: tls contains tbd
                                properties:
                                  ca:
                                    description: ca contains tbd
                                    type: string
                                  client:
                                    description: client contains tbd
                                    properties:
                                      ca:
                                        description: ca contains tbd
                                        type: string
                                    type: object
                                  crt:
                                    description: crt contains tbd
                                    type: string
                                  key:
                                    description: key contains tbd
                                    type: string
                                  peer:
                                    description: peer contains tbd
                                    properties:
                                      ca:
                                        description: ca contains tbd
                                        type: string
                                    type: object
                                type: object
                              tmpurlkey:
                                description: tmpurlkey contains tbd
                                type: string
                              username:
                                description: username contains tbd
                                type: string
                            type: object
                          description: auth contains tbd
                          type: object
                        hosts:
                          description: hosts contains tbd
                          properties:
                            default:
                              description: default contains tbd
                              type: string
                            discovery:
                              description: discovery contains tbd
                              type: string
                            public:
                              description: public contains tbd
                              type: string
                          type: object
                        name:
                          description: name contains tbd
                          type: string
                        namespace:
                          description: namespace contains tbd
                          type: string
                        path:
                          description: path contains tbd
                          properties:
                            default:
                              description: default contains tbd
                              type: string
                            discovery:
                              description: discovery contains tbd
                              type: string
                            public:
                              description: public contains tbd
                              type: string
                          type: object
                        port:
                          additionalProperties:
                            properties:
                              default:
                                description: default contains tbd
                                format: int64
                                type: integer
                              internal:
                                description: internal contains tbd
                                format: int64
                                type: integer
                              nodeport:
                                description: nodeport contains tbd
                                format: int64
                                type: integer
                              public:
                                description: public contains tbd
                                format: int64
                                type: integer
                            type: object
                          description: port contains tbd
                          type: object
                        type:
                          description: type contains tbd
                          type: string
                      type: object
                images:
                  description: images contains tbd
                  properties:
                    tags:
                      additionalProperties:
                        type: string
                      description: tags contains tbd
                      type: object
                  type: object
                labels:
                  additionalProperties:
                    type: string
                  type: object
                pod:
                  description: pod contains tbd
                  properties:
                    affinity:
                      description: affinity contains tbd
                      type: object
                    lifecycle:
                      description: lifecycle contains tbd
                      type: object
                    replicas:
                      additionalProperties:
                        format: int64
                        type: integer
                      description: replicas contains tbd
                      type: object
                    resources:
                      additionalProperties:
                        properties:
                          limits:
                            description: limits contains tbd
                            properties:
                              cpu:
                                description: cpu contains tbd
                                type: string
                              memory:
                                description: memory contains tbd
                                type: string
                            type: object
                          requests:
                            description: requests contains tbd
                            properties:
                              cpu:
                                description: cpu contains tbd
                                type: string
                              memory:
                                description: memory contains tbd
                                type: string
                            type: object
                        type: object
                      description: resources contains tbd
                      type: object
                  type: object
              type: object
          required:
          - source
          type: object
        status:
          properties:
            conditions:
              description: 'List of conditions and states related to the resource.'
              items:
                properties:
                  lastTransitionTime:
                    format: date-time
                    type: string
                  message:
                    type: string
                  reason:
                    type: string
                  resourceName:
                    type: string
                  resourceVersion:
                    format: int32
                    type: integer
                  status:
                    type: string
                  type:
                    type: string
                required:
                - type
                - status
                type: object
              type: array
            reason:
              description: Reason indicates the reason for any related failures.
              type: string
            satisfied:
              description: Satisfied indicates if the release's ActualState satisfies
                its target state
              type: boolean
          required:
          - satisfied
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/crds/CommonAddresses.yaml
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: commonaddressess.my.group.org
spec:
  additionalPrinterColumns:
  group: my.group.org
  version: v1alpha1
  names:
    kind: CommonAddresses
    plural: commonaddressess
    shortNames:
    - pcaddr
  scope: Namespaced
  subresources:
    status: {}
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
          type: object
          properties:
            dns:
              type: object
              properties:
                cluster_domain:
                  type: string
                service_ip:
                  type: string
                upstream_servers:
                  type: array
                  items:
                    type: string
            etcd:
              type: object
              properties:
                container_port:
                  type: number
                haproxy_port:
                  type: number
            node_ports:
              type: object
              properties:
                drydock_api:
                  type: number
                maas_api:
                  type: number
                maas_proxy:
                  type: number
                shipyard_api:
                  type: number
                airflow_web:
                  type: number
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
...
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/crds/EndpointCatalogue.yaml
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: endpointcatalogues.my.group.org
spec:
  additionalPrinterColumns:
  group: my.group.org
  version: v1alpha1
  names:
    kind: EndpointCatalogue
    plural: endpointcatalogues
    shortNames:
    - pendptcat
  scope: Namespaced
  subresources:
    status: {}
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
          type: object
          description: resources contains tbd
          additionalProperties:
            type: object
            properties:
              labels:
                additionalProperties:
                  type: string
                type: object
              pod:
                description: pod contains tbd
                properties:
                  affinity:
                    description: affinity contains tbd
                    type: object
                  lifecycle:
                    description: lifecycle contains tbd
                    type: object
                  replicas:
                    additionalProperties:
                      format: int64
                      type: integer
                    description: replicas contains tbd
                    type: object
                  resources:
                    additionalProperties:
                      properties:
                        limits:
                          description: limits contains tbd
                          properties:
                            cpu:
                              description: cpu contains tbd
                              type: string
                            memory:
                              description: memory contains tbd
                              type: string
                          type: object
                        requests:
                          description: requests contains tbd
                          properties:
                            cpu:
                              description: cpu contains tbd
                              type: string
                            memory:
                              description: memory contains tbd
                              type: string
                          type: object
                      type: object
                    description: resources contains tbd
                    type: object
                type: object
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
...
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/crds/SoftwareVersions.yaml
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: softwareversionss.my.group.org
spec:
  additionalPrinterColumns:
  group: my.group.org
  version: v1alpha1
  names:
    kind: SoftwareVersions
    plural: softwareversionss
    shortNames:
    - psoftver
  scope: Namespaced
  subresources:
    status: {}
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
            charts:
              type: object
              properties:
                wordpress:
                  type: object
                  properties:
                    type:
                      type: string
                    location:
                      type: string
                    subpath:
                      type: string
                    reference:
                      type: string
                mysql:
                  type: object
                  properties:
                    type:
                      type: string
                    location:
                      type: string
                    subpath:
                      type: string
                    reference:
                      type: string
            images:
              type: object
              properties:
                wordpress:
                  type: object
                  properties:
                    registry:
                      type: string
                    repository:
                      type: string
                    tag:
                      type: string
                mysql:
                  type: object
                  properties:
                    registry:
                      type: string
                    repository:
                      type: string
                    tag:
                      type: string
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
...
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/Chart.yaml
varReference:
- kind: Chart
  path: spec/values/endpoints/messaging/auth/user/password
- kind: Chart
  path: spec/source
- kind: Chart
  path: spec/values/images
- kind: Chart
  path: spec/values/labels
EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/CommonAddresses.yaml
varReference:

EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/Deployment.yaml
varReference:
- kind: Deployment
  path: spec/template/spec/containers/image
- kind: Deployment
  path: metadata/labels
- kind: Deployment
  path: spec/template/metadata/labels
- kind: Deployment
  path: spec/selector/matchLabels
- kind: Deployment
  path: spec/template/spec/initContainers
EOF
```


### Preparation Step Resource10

<!-- @createResource10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/EndpointCatalogue.yaml
varReference:

EOF
```


### Preparation Step Resource11

<!-- @createResource11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/Service.yaml
varReference:
- kind: Service
  path: metadata/labels
- kind: Service
  path: spec/selector
EOF
```


### Preparation Step Resource12

<!-- @createResource12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomizeconfig/SoftwareVersions.yaml
varReference:

EOF
```


### Preparation Step Resource13

<!-- @createResource13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/mysql/deployment.yaml
apiVersion: apps/v1beta2 # for versions before 1.9.0 use apps/v1beta2
kind: Deployment
metadata:
  name: mysql
  labels: $(EndpointCatalogue.endpoints.spec.mysql.labels)
spec:
  selector:
    matchLabels: $(EndpointCatalogue.endpoints.spec.mysql.labels)
  strategy:
    type: Recreate
  template:
    metadata:
      labels: $(EndpointCatalogue.endpoints.spec.mysql.labels)
    spec:
      initContainers: $(Deployment.wordpress.spec.template.spec.initContainers)
      containers:
      - image: mysql:$(SoftwareVersions.software-versions.spec.images.mysql.tag)
        name: mysql
        env:
        - name: MYSQL_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mysql-pass
              key: password
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - name: mysql-persistent-storage
          mountPath: /var/lib/mysql
      volumes:
      - name: mysql-persistent-storage
        emptyDir: {}
EOF
```


### Preparation Step Resource14

<!-- @createResource14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/mysql/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysql-pass
type: Opaque
data:
  # Default password is "admin".
  password: YWRtaW4=
EOF
```


### Preparation Step Resource15

<!-- @createResource15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/mysql/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql
  labels: $(EndpointCatalogue.endpoints.spec.mysql.labels)
spec:
  ports:
    - port: 3306
  selector: $(EndpointCatalogue.endpoints.spec.mysql.labels)
EOF
```


### Preparation Step Resource16

<!-- @createResource16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/wordpress/deployment.yaml
apiVersion: apps/v1beta2 # for versions before 1.9.0 use apps/v1beta2
kind: Deployment
metadata:
  name: wordpress
  labels: $(EndpointCatalogue.endpoints.spec.wordpress.labels)
spec:
  selector:
    matchLabels: $(EndpointCatalogue.endpoints.spec.wordpress.labels)
  strategy:
    type: Recreate
  template:
    metadata:
      # wordpress.labels tree will be inserted into metatda.labels.
      # foo: bar label will then be added
      labels:
        parent-inline: $(EndpointCatalogue.endpoints.spec.wordpress.labels)
        foo: bar
    spec:
      initContainers:
      - name: init1
        image: busybox
        command:
        - 'sh'
        - '-c'
        - 'echo $(Service.wordpress.metadata.name) && echo $(Service.mysql.metadata.name)'
      - name: init2
        image: busybox
        command:
        - 'sh'
        - '-c'
        - 'echo $(CommonAddresses.common-addresses.spec.dns.upstream_servers[2])'
      containers:
      - image: wordpress:$(SoftwareVersions.software-versions.spec.images.wordpress.tag)
        name: wordpress
        ports:
        - containerPort: 80
          name: wordpress
        volumeMounts:
        - name: wordpress-persistent-storage
          mountPath: /var/www/html
        env:
        - name: WORDPRESS_DB_HOST
          value: $(Service.mysql.metadata.name)
        - name: WORDPRESS_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mysql-pass
              key: password
      volumes:
      - name: wordpress-persistent-storage
        emptyDir: {}
EOF
```


### Preparation Step Resource17

<!-- @createResource17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/wordpress/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: wordpress
  labels: $(EndpointCatalogue.endpoints.spec.wordpress.labels)
spec:
  ports:
    - port: 80
  selector: $(EndpointCatalogue.endpoints.spec.wordpress.labels)
  type: LoadBalancer
EOF
```


### Preparation Step Resource18

<!-- @createResource18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/dev/common-addresses.yaml
---
apiVersion: my.group.org/v1alpha1
kind: CommonAddresses
metadata:
  name: common-addresses
spec:
  dns:
    cluster_domain: cluster.local
    ingress_domain: dev.my.group.org
    service_ip: 10.96.0.10
    upstream_servers:
    - 8.8.8.8
    - 8.8.4.4
    - 10.0.2.3
EOF
```


### Preparation Step Resource19

<!-- @createResource19 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/dev/devtools.yaml
---
apiVersion: my.group.org/v1alpha1
kind: Chart
metadata:
  name: wordpress
spec:
  source: $(SoftwareVersions.software-versions.spec.charts.wordpress)
  values:
    images: $(SoftwareVersions.software-versions.spec.images.wordpress)
    labels: $(EndpointCatalogue.endpoints.spec.wordpress.labels)
    pod:
      replicas:
        api: 1
---
apiVersion: my.group.org/v1alpha1
kind: Chart
metadata:
  name: mysql
spec:
  source: $(SoftwareVersions.software-versions.spec.charts.mysql)
  values:
    images: $(SoftwareVersions.software-versions.spec.images.mysql)
    labels: $(EndpointCatalogue.endpoints.spec.mysql.labels)
    pod:
      replicas:
        api: 1


EOF
```


### Preparation Step Resource20

<!-- @createResource20 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/dev/endpoints.yaml
---
apiVersion: my.group.org/v1alpha1
kind: EndpointCatalogue
metadata:
  name: endpoints
spec:
  mysql:
    labels:
       app: mysql
  wordpress:
    labels:
       app: wordpress

EOF
```


### Preparation Step Resource21

<!-- @createResource21 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/dev/passphrases.yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysql-pass
type: Opaque
data:
  # dev password is "devmysqlpasswd".
  password: ZGV2bXlzcWxwYXNzd2Q=
EOF
```


### Preparation Step Resource22

<!-- @createResource22 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/dev/versions.yaml
---
apiVersion: my.group.org/v1alpha1
kind: SoftwareVersions
metadata:
  name: software-versions
spec:
  charts:
    wordpress:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: wordpress
      type: git
    mysql:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: mysql
      type: git
  images:
    wordpress:
      registory: docker.io
      repository: wordpress
      tag: 5.2.1-apache
    mysql:
      registory: docker.io
      repository: mysql
      tag: '5.7'
EOF
```


### Preparation Step Resource23

<!-- @createResource23 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/common/endpoints.yaml
---
apiVersion: my.group.org/v1alpha1
kind: EndpointCatalogue
metadata:
  name: endpoints
spec:
  mysql:
    labels:
       app: mysql
  wordpress:
    labels:
       app: wordpress

EOF
```


### Preparation Step Resource24

<!-- @createResource24 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/common/versions.yaml
---
apiVersion: my.group.org/v1alpha1
kind: SoftwareVersions
metadata:
  name: software-versions
spec:
  charts:
    wordpress:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: wordpress
      type: git
    mysql:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: mysql
      type: git
  images:
    wordpress:
      registory: docker.io
      repository: wordpress
      tag: 4.8-apache
    mysql:
      registory: docker.io
      repository: mysql
      tag: '5.6'
EOF
```


### Preparation Step Resource25

<!-- @createResource25 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/site1/common-addresses.yaml
---
apiVersion: my.group.org/v1alpha1
kind: CommonAddresses
metadata:
  name: common-addresses
spec:
  dns:
    cluster_domain: cluster.local
    ingress_domain: site1.my.group.org
    service_ip: 10.96.0.10
    upstream_servers:
    - 8.8.8.8
    - 8.8.4.4
    - 10.0.2.3
EOF
```


### Preparation Step Resource26

<!-- @createResource26 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/site1/passphrases.yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysql-pass
type: Opaque
data:
  # site1 password is "site1mysqlpassword".
  password: c2l0ZTFteXNxbHBhc3N3b3JkCg==
EOF
```


### Preparation Step Resource27

<!-- @createResource27 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/site2/common-addresses.yaml
---
apiVersion: my.group.org/v1alpha1
kind: CommonAddresses
metadata:
  name: common-addresses
spec:
  dns:
    cluster_domain: cluster.local
    ingress_domain: site2.my.group.org
    service_ip: 10.96.0.10
    upstream_servers:
    - 8.8.8.8
    - 8.8.4.4
    - 10.0.2.3
EOF
```


### Preparation Step Resource28

<!-- @createResource28 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/production/site2/passphrases.yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysql-pass
type: Opaque
data:
  # site2 password is "site2mysqlpassword".
  password: c2l0ZTJteXNxbHBhc3N3b3JkCg==
EOF
```

## Execution

<!-- @build @dev @test -->
```bash
mkdir -p ${DEMO_HOME}/actual/dev
kustomize build ${DEMO_HOME}/dev -o ${DEMO_HOME}/actual/dev
```

<!-- @build @site1 @test -->
```bash
mkdir -p ${DEMO_HOME}/actual/site1
kustomize build ${DEMO_HOME}/production/site1 -o ${DEMO_HOME}/actual/site1
```

<!-- @build @site2 @test -->
```bash
mkdir -p ${DEMO_HOME}/actual/site2
kustomize build ${DEMO_HOME}/production/site2 -o ${DEMO_HOME}/actual/site2
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir -p ${DEMO_HOME}/expected
mkdir -p ${DEMO_HOME}/expected/dev
mkdir -p ${DEMO_HOME}/expected/site1
mkdir -p ${DEMO_HOME}/expected/site2
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev/apps_v1beta2_deployment_dev-mysql.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: mysql
  name: dev-mysql
spec:
  selector:
    matchLabels:
      app: mysql
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - env:
        - name: MYSQL_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: dev-mysql-pass
        image: mysql:5.7
        name: mysql
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - mountPath: /var/lib/mysql
          name: mysql-persistent-storage
      initContainers:
      - command:
        - sh
        - -c
        - echo dev-wordpress && echo dev-mysql
        image: busybox
        name: init1
      - command:
        - sh
        - -c
        - echo 10.0.2.3
        image: busybox
        name: init2
      volumes:
      - emptyDir: {}
        name: mysql-persistent-storage
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev/apps_v1beta2_deployment_dev-wordpress.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: wordpress
  name: dev-wordpress
spec:
  selector:
    matchLabels:
      app: wordpress
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: wordpress
        foo: bar
    spec:
      containers:
      - env:
        - name: WORDPRESS_DB_HOST
          value: dev-mysql
        - name: WORDPRESS_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: dev-mysql-pass
        image: wordpress:5.2.1-apache
        name: wordpress
        ports:
        - containerPort: 80
          name: wordpress
        volumeMounts:
        - mountPath: /var/www/html
          name: wordpress-persistent-storage
      initContainers:
      - command:
        - sh
        - -c
        - echo dev-wordpress && echo dev-mysql
        image: busybox
        name: init1
      - command:
        - sh
        - -c
        - echo 10.0.2.3
        image: busybox
        name: init2
      volumes:
      - emptyDir: {}
        name: wordpress-persistent-storage
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev/~g_v1_secret_dev-mysql-pass.yaml
apiVersion: v1
data:
  password: ZGV2bXlzcWxwYXNzd2Q=
kind: Secret
metadata:
  name: dev-mysql-pass
type: Opaque
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev/~g_v1_service_dev-mysql.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: mysql
  name: dev-mysql
spec:
  ports:
  - port: 3306
  selector:
    app: mysql
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev/~g_v1_service_dev-wordpress.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: wordpress
  name: dev-wordpress
spec:
  ports:
  - port: 80
  selector:
    app: wordpress
  type: LoadBalancer
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev/my.group.org_v1alpha1_chart_dev-mysql.yaml
apiVersion: my.group.org/v1alpha1
kind: Chart
metadata:
  name: dev-mysql
spec:
  source:
    location: https://github.com/helm/charts/blob/
    reference: latest
    subpath: mysql
    type: git
  values:
    images:
      registory: docker.io
      repository: mysql
      tag: "5.7"
    labels:
      app: mysql
    pod:
      replicas:
        api: 1
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev/my.group.org_v1alpha1_chart_dev-wordpress.yaml
apiVersion: my.group.org/v1alpha1
kind: Chart
metadata:
  name: dev-wordpress
spec:
  source:
    location: https://github.com/helm/charts/blob/
    reference: latest
    subpath: wordpress
    type: git
  values:
    images:
      registory: docker.io
      repository: wordpress
      tag: 5.2.1-apache
    labels:
      app: wordpress
    pod:
      replicas:
        api: 1
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev/my.group.org_v1alpha1_commonaddresses_dev-common-addresses.yaml
apiVersion: my.group.org/v1alpha1
kind: CommonAddresses
metadata:
  name: dev-common-addresses
spec:
  dns:
    cluster_domain: cluster.local
    ingress_domain: dev.my.group.org
    service_ip: 10.96.0.10
    upstream_servers:
    - 8.8.8.8
    - 8.8.4.4
    - 10.0.2.3
EOF
```


### Verification Step Expected8

<!-- @createExpected8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev/my.group.org_v1alpha1_endpointcatalogue_dev-endpoints.yaml
apiVersion: my.group.org/v1alpha1
kind: EndpointCatalogue
metadata:
  name: dev-endpoints
spec:
  mysql:
    labels:
      app: mysql
  wordpress:
    labels:
      app: wordpress
EOF
```


### Verification Step Expected9

<!-- @createExpected9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/dev/my.group.org_v1alpha1_softwareversions_dev-software-versions.yaml
apiVersion: my.group.org/v1alpha1
kind: SoftwareVersions
metadata:
  name: dev-software-versions
spec:
  charts:
    mysql:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: mysql
      type: git
    wordpress:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: wordpress
      type: git
  images:
    mysql:
      registory: docker.io
      repository: mysql
      tag: "5.7"
    wordpress:
      registory: docker.io
      repository: wordpress
      tag: 5.2.1-apache
EOF
```


### Verification Step Expected10

<!-- @createExpected10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site1/apps_v1beta2_deployment_mysql.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: mysql
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - env:
        - name: MYSQL_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: mysql-pass
        image: mysql:5.6
        name: mysql
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - mountPath: /var/lib/mysql
          name: mysql-persistent-storage
      initContainers:
      - command:
        - sh
        - -c
        - echo wordpress && echo mysql
        image: busybox
        name: init1
      - command:
        - sh
        - -c
        - echo 10.0.2.3
        image: busybox
        name: init2
      volumes:
      - emptyDir: {}
        name: mysql-persistent-storage
EOF
```


### Verification Step Expected11

<!-- @createExpected11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site1/apps_v1beta2_deployment_wordpress.yaml
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
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: wordpress
        foo: bar
    spec:
      containers:
      - env:
        - name: WORDPRESS_DB_HOST
          value: mysql
        - name: WORDPRESS_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: mysql-pass
        image: wordpress:4.8-apache
        name: wordpress
        ports:
        - containerPort: 80
          name: wordpress
        volumeMounts:
        - mountPath: /var/www/html
          name: wordpress-persistent-storage
      initContainers:
      - command:
        - sh
        - -c
        - echo wordpress && echo mysql
        image: busybox
        name: init1
      - command:
        - sh
        - -c
        - echo 10.0.2.3
        image: busybox
        name: init2
      volumes:
      - emptyDir: {}
        name: wordpress-persistent-storage
EOF
```


### Verification Step Expected12

<!-- @createExpected12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site1/~g_v1_secret_mysql-pass.yaml
apiVersion: v1
data:
  password: c2l0ZTFteXNxbHBhc3N3b3JkCg==
kind: Secret
metadata:
  name: mysql-pass
type: Opaque
EOF
```


### Verification Step Expected13

<!-- @createExpected13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site1/~g_v1_service_mysql.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: mysql
  name: mysql
spec:
  ports:
  - port: 3306
  selector:
    app: mysql
EOF
```


### Verification Step Expected14

<!-- @createExpected14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site1/~g_v1_service_wordpress.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: wordpress
  name: wordpress
spec:
  ports:
  - port: 80
  selector:
    app: wordpress
  type: LoadBalancer
EOF
```


### Verification Step Expected15

<!-- @createExpected15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site1/my.group.org_v1alpha1_commonaddresses_common-addresses.yaml
apiVersion: my.group.org/v1alpha1
kind: CommonAddresses
metadata:
  name: common-addresses
spec:
  dns:
    cluster_domain: cluster.local
    ingress_domain: site1.my.group.org
    service_ip: 10.96.0.10
    upstream_servers:
    - 8.8.8.8
    - 8.8.4.4
    - 10.0.2.3
EOF
```


### Verification Step Expected16

<!-- @createExpected16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site1/my.group.org_v1alpha1_endpointcatalogue_endpoints.yaml
apiVersion: my.group.org/v1alpha1
kind: EndpointCatalogue
metadata:
  name: endpoints
spec:
  mysql:
    labels:
      app: mysql
  wordpress:
    labels:
      app: wordpress
EOF
```


### Verification Step Expected17

<!-- @createExpected17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site1/my.group.org_v1alpha1_softwareversions_software-versions.yaml
apiVersion: my.group.org/v1alpha1
kind: SoftwareVersions
metadata:
  name: software-versions
spec:
  charts:
    mysql:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: mysql
      type: git
    wordpress:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: wordpress
      type: git
  images:
    mysql:
      registory: docker.io
      repository: mysql
      tag: "5.6"
    wordpress:
      registory: docker.io
      repository: wordpress
      tag: 4.8-apache
EOF
```


### Verification Step Expected18

<!-- @createExpected18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site2/apps_v1beta2_deployment_mysql.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: mysql
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - env:
        - name: MYSQL_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: mysql-pass
        image: mysql:5.6
        name: mysql
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - mountPath: /var/lib/mysql
          name: mysql-persistent-storage
      initContainers:
      - command:
        - sh
        - -c
        - echo wordpress && echo mysql
        image: busybox
        name: init1
      - command:
        - sh
        - -c
        - echo 10.0.2.3
        image: busybox
        name: init2
      volumes:
      - emptyDir: {}
        name: mysql-persistent-storage
EOF
```


### Verification Step Expected19

<!-- @createExpected19 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site2/apps_v1beta2_deployment_wordpress.yaml
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
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: wordpress
        foo: bar
    spec:
      containers:
      - env:
        - name: WORDPRESS_DB_HOST
          value: mysql
        - name: WORDPRESS_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: mysql-pass
        image: wordpress:4.8-apache
        name: wordpress
        ports:
        - containerPort: 80
          name: wordpress
        volumeMounts:
        - mountPath: /var/www/html
          name: wordpress-persistent-storage
      initContainers:
      - command:
        - sh
        - -c
        - echo wordpress && echo mysql
        image: busybox
        name: init1
      - command:
        - sh
        - -c
        - echo 10.0.2.3
        image: busybox
        name: init2
      volumes:
      - emptyDir: {}
        name: wordpress-persistent-storage
EOF
```


### Verification Step Expected20

<!-- @createExpected20 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site2/~g_v1_secret_mysql-pass.yaml
apiVersion: v1
data:
  password: c2l0ZTJteXNxbHBhc3N3b3JkCg==
kind: Secret
metadata:
  name: mysql-pass
type: Opaque
EOF
```


### Verification Step Expected21

<!-- @createExpected21 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site2/~g_v1_service_mysql.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: mysql
  name: mysql
spec:
  ports:
  - port: 3306
  selector:
    app: mysql
EOF
```


### Verification Step Expected22

<!-- @createExpected22 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site2/~g_v1_service_wordpress.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: wordpress
  name: wordpress
spec:
  ports:
  - port: 80
  selector:
    app: wordpress
  type: LoadBalancer
EOF
```


### Verification Step Expected23

<!-- @createExpected23 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site2/my.group.org_v1alpha1_commonaddresses_common-addresses.yaml
apiVersion: my.group.org/v1alpha1
kind: CommonAddresses
metadata:
  name: common-addresses
spec:
  dns:
    cluster_domain: cluster.local
    ingress_domain: site2.my.group.org
    service_ip: 10.96.0.10
    upstream_servers:
    - 8.8.8.8
    - 8.8.4.4
    - 10.0.2.3
EOF
```


### Verification Step Expected24

<!-- @createExpected24 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site2/my.group.org_v1alpha1_endpointcatalogue_endpoints.yaml
apiVersion: my.group.org/v1alpha1
kind: EndpointCatalogue
metadata:
  name: endpoints
spec:
  mysql:
    labels:
      app: mysql
  wordpress:
    labels:
      app: wordpress
EOF
```


### Verification Step Expected25

<!-- @createExpected25 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/site2/my.group.org_v1alpha1_softwareversions_software-versions.yaml
apiVersion: my.group.org/v1alpha1
kind: SoftwareVersions
metadata:
  name: software-versions
spec:
  charts:
    mysql:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: mysql
      type: git
    wordpress:
      location: https://github.com/helm/charts/blob/
      reference: latest
      subpath: wordpress
      type: git
  images:
    mysql:
      registory: docker.io
      repository: mysql
      tag: "5.6"
    wordpress:
      registory: docker.io
      repository: wordpress
      tag: 4.8-apache
EOF
```

<!-- @compareActualToExpected @dev @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual/dev $DEMO_HOME/expected/dev | wc -l); \
echo $?
```

<!-- @compareActualToExpected @site1 @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual/site1 $DEMO_HOME/expected/site1 | wc -l); \
echo $?
```

<!-- @compareActualToExpected @site2 @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual/site2 $DEMO_HOME/expected/site2 | wc -l); \
echo $?
```
