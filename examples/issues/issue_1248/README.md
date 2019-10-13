# Feature Test for Issue 1248


This folder contains files describing how to address [Issue 1248](https://github.com/kubernetes-sigs/kustomize/issues/1248)

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
mkdir -p ${DEMO_HOME}/folder1
mkdir -p ${DEMO_HOME}/folder1/app-installations
mkdir -p ${DEMO_HOME}/folder1/app-installations/subfolder
mkdir -p ${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite
mkdir -p ${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib
mkdir -p ${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s
mkdir -p ${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s-kubedb-mysql
mkdir -p ${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite
mkdir -p ${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite-with-init-script
mkdir -p ${DEMO_HOME}/folder2
mkdir -p ${DEMO_HOME}/folder2/app-installations
mkdir -p ${DEMO_HOME}/folder2/app-installations/subfolder
mkdir -p ${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite
mkdir -p ${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib
mkdir -p ${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s
mkdir -p ${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s-kubedb-mysql
mkdir -p ${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite
mkdir -p ${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite-with-init-script
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s-kubedb-mysql/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- ../kubedb-mysql-etherpad-lite-with-init-script
- ../etherpad-lite-k8s
patchesStrategicMerge:
- configmap.yaml
- deployment.yaml
images:
- name: etherpad/etherpad
  # This is required until 1.8 comes out to be able to use env vars in settings.json
  newTag: latest
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- configmap.yaml
- deployment.yaml
- service.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- etherpad-mysql.yaml
vars:
- name: MYSQL_SERVICE1
  objref:
    apiVersion: kubedb.com/v1alpha1
    kind: MySQL
    name: etherpad-mysql
    # namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
  fieldref:
    fieldpath: metadata.name
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite-with-init-script/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- ../kubedb-mysql-etherpad-lite
- etherpad-mysql-init-configmap.yaml
patchesStrategicMerge:
- etherpad-mysql-with-init-script.yaml
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
  - ./etherpad-lite/lib/etherpad-lite-k8s-kubedb-mysql
commonLabels:
  k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
kind: Kustomization
namePrefix: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
resources:
  - namespace.yaml
  - app-installations/subfolder
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s-kubedb-mysql/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- ../kubedb-mysql-etherpad-lite-with-init-script
- ../etherpad-lite-k8s
patchesStrategicMerge:
- configmap.yaml
- deployment.yaml
images:
- name: etherpad/etherpad
  # This is required until 1.8 comes out to be able to use env vars in settings.json
  newTag: latest
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- configmap.yaml
- deployment.yaml
- service.yaml
EOF
```


### Preparation Step KustomizationFile8

<!-- @createKustomizationFile8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- etherpad-mysql.yaml
vars:
- name: MYSQL_SERVICE2
  objref:
    apiVersion: kubedb.com/v1alpha1
    kind: MySQL
    name: etherpad-mysql
    # namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
  fieldref:
    fieldpath: metadata.name
EOF
```


### Preparation Step KustomizationFile9

<!-- @createKustomizationFile9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite-with-init-script/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- ../kubedb-mysql-etherpad-lite
- etherpad-mysql-init-configmap.yaml
patchesStrategicMerge:
- etherpad-mysql-with-init-script.yaml
EOF
```


### Preparation Step KustomizationFile10

<!-- @createKustomizationFile10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
  - ./etherpad-lite/lib/etherpad-lite-k8s-kubedb-mysql
commonLabels:
  k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
kind: Kustomization
namePrefix: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ym
EOF
```


### Preparation Step KustomizationFile11

<!-- @createKustomizationFile11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
resources:
  - namespace.yaml
  - app-installations/subfolder

EOF
```


### Preparation Step KustomizationFile12

<!-- @createKustomizationFile12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - folder2
  - folder1

EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: etherpad
data:
  settings.json: |
    {
      "skinName":"colibris",
      "title":"Etherpad on Kubernetes"
    }
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: etherpad
spec:
  replicas: 1
  selector:
    matchLabels:
      app: etherpad
  template:
    metadata:
      labels:
        app: etherpad
    spec:
      containers:
      - name: etherpad
        image: etherpad/etherpad:1.7.5
        ports:
        - containerPort: 9001
          name: web
        volumeMounts:
        - name: "config"
          mountPath: "/opt/etherpad/settings.json"
          subPath: "settings.json"
      volumes:
      - name: config
        configMap:
          name: etherpad
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s-kubedb-mysql/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: etherpad
data:
  settings.json: |
    {
      "skinName":"colibris",
      "title":"Etherpad on Kubernetes w/ MySQL",
      "dbType": "${ETHERPAD_DB_TYPE:mysql}",
      "dbSettings": {
        "database": "${ETHERPAD_DB_DATABASE}",
        "host": "${ETHERPAD_DB_HOST}",
        "password": "${ETHERPAD_DB_PASSWORD}",
        "user": "${ETHERPAD_DB_USER}"
      }
    }
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s-kubedb-mysql/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: etherpad
spec:
  template:
    spec:
      containers:
      - name: etherpad
        env:
        - name: ETHERPAD_DB_TYPE
          value: mysql
        - name: ETHERPAD_DB_HOST
          value: $(MYSQL_SERVICE1)
        - name: ETHERPAD_DB_DATABASE
          value: etherpad_lite_db
        - name: ETHERPAD_DB_USER
          valueFrom:
            secretKeyRef:
              name: etherpad-mysql-auth
              key: username
        - name: ETHERPAD_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: etherpad-mysql-auth
              key: password
        volumeMounts:
        - name: "config"
          mountPath: "/opt/etherpad-lite/settings.json"
          subPath: "settings.json"
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: etherpad
spec:
  selector:
    app: etherpad
  ports:
  - name: web
    port: 80
    targetPort: web
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite/etherpad-mysql.yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: etherpad-mysql
spec:
  version: "5.7.25"
  storageType: Durable
  terminationPolicy: WipeOut
  storage:
    storageClassName: "default"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite-with-init-script/etherpad-mysql-init-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: etherpad-mysql-init
data:
  init.sql: |
    create database `etherpad_lite_db`;
    use `etherpad_lite_db`;

    CREATE TABLE `store` (
      `key` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
      `value` longtext COLLATE utf8mb4_bin NOT NULL,
      PRIMARY KEY (`key`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite-with-init-script/etherpad-mysql-with-init-script.yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: etherpad-mysql
spec:
  init:
    scriptSource:
      configMap:
        name: etherpad-mysql-init
EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder1/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: etherpad
data:
  settings.json: |
    {
      "skinName":"colibris",
      "title":"Etherpad on Kubernetes"
    }
EOF
```


### Preparation Step Resource10

<!-- @createResource10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: etherpad
spec:
  replicas: 1
  selector:
    matchLabels:
      app: etherpad
  template:
    metadata:
      labels:
        app: etherpad
    spec:
      containers:
      - name: etherpad
        image: etherpad/etherpad:1.7.5
        ports:
        - containerPort: 9001
          name: web
        volumeMounts:
        - name: "config"
          mountPath: "/opt/etherpad/settings.json"
          subPath: "settings.json"
      volumes:
      - name: config
        configMap:
          name: etherpad
EOF
```


### Preparation Step Resource11

<!-- @createResource11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s-kubedb-mysql/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: etherpad
data:
  settings.json: |
    {
      "skinName":"colibris",
      "title":"Etherpad on Kubernetes w/ MySQL",
      "dbType": "${ETHERPAD_DB_TYPE:mysql}",
      "dbSettings": {
        "database": "${ETHERPAD_DB_DATABASE}",
        "host": "${ETHERPAD_DB_HOST}",
        "password": "${ETHERPAD_DB_PASSWORD}",
        "user": "${ETHERPAD_DB_USER}"
      }
    }
EOF
```


### Preparation Step Resource12

<!-- @createResource12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s-kubedb-mysql/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: etherpad
spec:
  template:
    spec:
      containers:
      - name: etherpad
        env:
        - name: ETHERPAD_DB_TYPE
          value: mysql
        - name: ETHERPAD_DB_HOST
          value: $(MYSQL_SERVICE2)
        - name: ETHERPAD_DB_DATABASE
          value: etherpad_lite_db
        - name: ETHERPAD_DB_USER
          valueFrom:
            secretKeyRef:
              name: etherpad-mysql-auth
              key: username
        - name: ETHERPAD_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: etherpad-mysql-auth
              key: password
        volumeMounts:
        - name: "config"
          mountPath: "/opt/etherpad-lite/settings.json"
          subPath: "settings.json"
EOF
```


### Preparation Step Resource13

<!-- @createResource13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/etherpad-lite-k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: etherpad
spec:
  selector:
    app: etherpad
  ports:
  - name: web
    port: 80
    targetPort: web
EOF
```


### Preparation Step Resource14

<!-- @createResource14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite/etherpad-mysql.yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: etherpad-mysql
spec:
  version: "5.7.25"
  storageType: Durable
  terminationPolicy: WipeOut
  storage:
    storageClassName: "default"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
EOF
```


### Preparation Step Resource15

<!-- @createResource15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite-with-init-script/etherpad-mysql-init-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: etherpad-mysql-init
data:
  init.sql: |
    create database `etherpad_lite_db`;
    use `etherpad_lite_db`;

    CREATE TABLE `store` (
      `key` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
      `value` longtext COLLATE utf8mb4_bin NOT NULL,
      PRIMARY KEY (`key`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
EOF
```


### Preparation Step Resource16

<!-- @createResource16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/app-installations/subfolder/etherpad-lite/lib/kubedb-mysql-etherpad-lite-with-init-script/etherpad-mysql-with-init-script.yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: etherpad-mysql
spec:
  init:
    scriptSource:
      configMap:
        name: etherpad-mysql-init
EOF
```


### Preparation Step Resource17

<!-- @createResource17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/folder2/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/folder1 -o ${DEMO_HOME}/actual/folder1.yaml
kustomize build ${DEMO_HOME}/folder2 -o ${DEMO_HOME}/actual/folder2.yaml
kustomize build ${DEMO_HOME} -o ${DEMO_HOME}/actual/both.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/both.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
---
apiVersion: v1
kind: Namespace
metadata:
  name: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
---
apiVersion: v1
data:
  settings.json: |
    {
      "skinName":"colibris",
      "title":"Etherpad on Kubernetes w/ MySQL",
      "dbType": "${ETHERPAD_DB_TYPE:mysql}",
      "dbSettings": {
        "database": "${ETHERPAD_DB_DATABASE}",
        "host": "${ETHERPAD_DB_HOST}",
        "password": "${ETHERPAD_DB_PASSWORD}",
        "user": "${ETHERPAD_DB_USER}"
      }
    }
kind: ConfigMap
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad
  namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
---
apiVersion: v1
data:
  init.sql: |
    create database `etherpad_lite_db`;
    use `etherpad_lite_db`;

    CREATE TABLE `store` (
      `key` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
      `value` longtext COLLATE utf8mb4_bin NOT NULL,
      PRIMARY KEY (`key`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
kind: ConfigMap
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad-mysql-init
  namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
---
apiVersion: v1
data:
  settings.json: |
    {
      "skinName":"colibris",
      "title":"Etherpad on Kubernetes w/ MySQL",
      "dbType": "${ETHERPAD_DB_TYPE:mysql}",
      "dbSettings": {
        "database": "${ETHERPAD_DB_DATABASE}",
        "host": "${ETHERPAD_DB_HOST}",
        "password": "${ETHERPAD_DB_PASSWORD}",
        "user": "${ETHERPAD_DB_USER}"
      }
    }
kind: ConfigMap
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad
  namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
---
apiVersion: v1
data:
  init.sql: |
    create database `etherpad_lite_db`;
    use `etherpad_lite_db`;

    CREATE TABLE `store` (
      `key` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
      `value` longtext COLLATE utf8mb4_bin NOT NULL,
      PRIMARY KEY (`key`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
kind: ConfigMap
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad-mysql-init
  namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
---
apiVersion: v1
kind: Service
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad
  namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
spec:
  ports:
  - name: web
    port: 80
    targetPort: web
  selector:
    app: etherpad
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
---
apiVersion: v1
kind: Service
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad
  namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
spec:
  ports:
  - name: web
    port: 80
    targetPort: web
  selector:
    app: etherpad
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad
  namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
spec:
  replicas: 1
  selector:
    matchLabels:
      app: etherpad
      k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  template:
    metadata:
      labels:
        app: etherpad
        k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
    spec:
      containers:
      - env:
        - name: ETHERPAD_DB_TYPE
          value: mysql
        - name: ETHERPAD_DB_HOST
          value: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad-mysql
        - name: ETHERPAD_DB_DATABASE
          value: etherpad_lite_db
        - name: ETHERPAD_DB_USER
          valueFrom:
            secretKeyRef:
              key: username
              name: etherpad-mysql-auth
        - name: ETHERPAD_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: etherpad-mysql-auth
        image: etherpad/etherpad:latest
        name: etherpad
        ports:
        - containerPort: 9001
          name: web
        volumeMounts:
        - mountPath: /opt/etherpad-lite/settings.json
          name: config
          subPath: settings.json
        - mountPath: /opt/etherpad/settings.json
          name: config
          subPath: settings.json
      volumes:
      - configMap:
          name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad
        name: config
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad
  namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
spec:
  replicas: 1
  selector:
    matchLabels:
      app: etherpad
      k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  template:
    metadata:
      labels:
        app: etherpad
        k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
    spec:
      containers:
      - env:
        - name: ETHERPAD_DB_TYPE
          value: mysql
        - name: ETHERPAD_DB_HOST
          value: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad-mysql
        - name: ETHERPAD_DB_DATABASE
          value: etherpad_lite_db
        - name: ETHERPAD_DB_USER
          valueFrom:
            secretKeyRef:
              key: username
              name: etherpad-mysql-auth
        - name: ETHERPAD_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: etherpad-mysql-auth
        image: etherpad/etherpad:latest
        name: etherpad
        ports:
        - containerPort: 9001
          name: web
        volumeMounts:
        - mountPath: /opt/etherpad-lite/settings.json
          name: config
          subPath: settings.json
        - mountPath: /opt/etherpad/settings.json
          name: config
          subPath: settings.json
      volumes:
      - configMap:
          name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad
        name: config
---
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad-mysql
  namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
spec:
  init:
    scriptSource:
      configMap:
        name: etherpad-mysql-init
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
    storageClassName: default
  storageType: Durable
  terminationPolicy: WipeOut
  version: 5.7.25
---
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad-mysql
  namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
spec:
  init:
    scriptSource:
      configMap:
        name: etherpad-mysql-init
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
    storageClassName: default
  storageType: Durable
  terminationPolicy: WipeOut
  version: 5.7.25
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/folder1.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
---
apiVersion: v1
data:
  settings.json: |
    {
      "skinName":"colibris",
      "title":"Etherpad on Kubernetes w/ MySQL",
      "dbType": "${ETHERPAD_DB_TYPE:mysql}",
      "dbSettings": {
        "database": "${ETHERPAD_DB_DATABASE}",
        "host": "${ETHERPAD_DB_HOST}",
        "password": "${ETHERPAD_DB_PASSWORD}",
        "user": "${ETHERPAD_DB_USER}"
      }
    }
kind: ConfigMap
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad
  namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
---
apiVersion: v1
data:
  init.sql: |
    create database `etherpad_lite_db`;
    use `etherpad_lite_db`;

    CREATE TABLE `store` (
      `key` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
      `value` longtext COLLATE utf8mb4_bin NOT NULL,
      PRIMARY KEY (`key`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
kind: ConfigMap
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad-mysql-init
  namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
---
apiVersion: v1
kind: Service
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad
  namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
spec:
  ports:
  - name: web
    port: 80
    targetPort: web
  selector:
    app: etherpad
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad
  namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
spec:
  replicas: 1
  selector:
    matchLabels:
      app: etherpad
      k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  template:
    metadata:
      labels:
        app: etherpad
        k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
    spec:
      containers:
      - env:
        - name: ETHERPAD_DB_TYPE
          value: mysql
        - name: ETHERPAD_DB_HOST
          value: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad-mysql
        - name: ETHERPAD_DB_DATABASE
          value: etherpad_lite_db
        - name: ETHERPAD_DB_USER
          valueFrom:
            secretKeyRef:
              key: username
              name: etherpad-mysql-auth
        - name: ETHERPAD_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: etherpad-mysql-auth
        image: etherpad/etherpad:latest
        name: etherpad
        ports:
        - containerPort: 9001
          name: web
        volumeMounts:
        - mountPath: /opt/etherpad-lite/settings.json
          name: config
          subPath: settings.json
        - mountPath: /opt/etherpad/settings.json
          name: config
          subPath: settings.json
      volumes:
      - configMap:
          name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad
        name: config
---
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: add961a2-b5c7-4ccd-b3c7-66f7c03c9c6e
  name: ai-zv58kz2nbox64fkrqptr94nurvqoxz88o52etherpad-mysql
  namespace: cstmr-72n16kk2an86kc4855ujzk9a9plo293274l
spec:
  init:
    scriptSource:
      configMap:
        name: etherpad-mysql-init
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
    storageClassName: default
  storageType: Durable
  terminationPolicy: WipeOut
  version: 5.7.25
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/folder2.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
---
apiVersion: v1
data:
  settings.json: |
    {
      "skinName":"colibris",
      "title":"Etherpad on Kubernetes w/ MySQL",
      "dbType": "${ETHERPAD_DB_TYPE:mysql}",
      "dbSettings": {
        "database": "${ETHERPAD_DB_DATABASE}",
        "host": "${ETHERPAD_DB_HOST}",
        "password": "${ETHERPAD_DB_PASSWORD}",
        "user": "${ETHERPAD_DB_USER}"
      }
    }
kind: ConfigMap
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad
  namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
---
apiVersion: v1
data:
  init.sql: |
    create database `etherpad_lite_db`;
    use `etherpad_lite_db`;

    CREATE TABLE `store` (
      `key` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
      `value` longtext COLLATE utf8mb4_bin NOT NULL,
      PRIMARY KEY (`key`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
kind: ConfigMap
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad-mysql-init
  namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
---
apiVersion: v1
kind: Service
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad
  namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
spec:
  ports:
  - name: web
    port: 80
    targetPort: web
  selector:
    app: etherpad
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad
  namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
spec:
  replicas: 1
  selector:
    matchLabels:
      app: etherpad
      k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  template:
    metadata:
      labels:
        app: etherpad
        k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
    spec:
      containers:
      - env:
        - name: ETHERPAD_DB_TYPE
          value: mysql
        - name: ETHERPAD_DB_HOST
          value: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad-mysql
        - name: ETHERPAD_DB_DATABASE
          value: etherpad_lite_db
        - name: ETHERPAD_DB_USER
          valueFrom:
            secretKeyRef:
              key: username
              name: etherpad-mysql-auth
        - name: ETHERPAD_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: etherpad-mysql-auth
        image: etherpad/etherpad:latest
        name: etherpad
        ports:
        - containerPort: 9001
          name: web
        volumeMounts:
        - mountPath: /opt/etherpad-lite/settings.json
          name: config
          subPath: settings.json
        - mountPath: /opt/etherpad/settings.json
          name: config
          subPath: settings.json
      volumes:
      - configMap:
          name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad
        name: config
---
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  labels:
    k8s.permanent.cloud/appInstallation.id: 45170c85-ec8b-4008-9d57-4524aa16f93f
  name: ai-w613mmojuo0qqir4pvc1l4rsr96mm6110ymetherpad-mysql
  namespace: cstmr-zvyjvn35b81rkfkr87fznrpanw5op9x5yo0
spec:
  init:
    scriptSource:
      configMap:
        name: etherpad-mysql-init
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
    storageClassName: default
  storageType: Durable
  terminationPolicy: WipeOut
  version: 5.7.25
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

