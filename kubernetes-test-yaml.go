package kubernetes

//this file has yaml s to test and show how to use interfaces
//Redis tested here as well

const TestdeploymentYAML = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: sname
  namespace: default
  labels:
    helm.sh/chart: default-springboot-3.6.0
    app.kubernetes.io/name: default-springboot
    app.kubernetes.io/instance: RELEASE-NAME
    app.kubernetes.io/managed-by: Helm
data:
  JAVA_OPTS: >
    -Xms3g -Xmx4g
  TOMCAT6_SECURITY: "no"
  SERVICE_NAME: "fee-tax"
  SPRING_CLOUD_CONSUL_ENABLED: "false"
  SHUTDOWN_WAIT: "10"
  HAZELCAST_CONSUL_ENABLED: "false"
  HAZELCAST_DISCOVERY_ENABLED: "true"
  HAZELCAST_NETWORK_CONFIG_JOIN_TCPIP_ENABLED: "false"
  HAZELCAST_NETWORK_CONFIG_JOIN_MULTICAST_ENABLED: "false"
  HAZELCAST_NETWORK_CONFIG_JOIN_KUBERNETES_ENABLED: "true"
  HAZELCAST_NETWORK_CONFIG_JOIN_KUBERNETES_NAMESPACE: "default"
  HAZELCAST_NETWORK_CONFIG_JOIN_KUBERNETES_SERVICE_NAME: "RELEASE-NAME"
  CONFIG_DEFAULT_RESOURCE_OVERRIDE: "config/config-default-k8s-dc-shared-services.yaml"
`



const TestdeploymentYAMLDep = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment-sk10
#  namespace: default
  namespace: vcase
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
`


const TestdeploymentYAMLReadis1 = `
#NAME: my-release
#LAST DEPLOYED: Tue May 11 17:37:03 2021
#NAMESPACE: default
#STATUS: pending-install
#REVISION: 1
#TEST SUITE: None
#HOOKS:
#MANIFEST:
---
# Source: redis/templates/serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-release-redis
  namespace: "default"
  labels:
    app.kubernetes.io/name: redis
    helm.sh/chart: redis-14.1.1
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/managed-by: Helm
`
const TestdeploymentYAMLReadis2 = `
---
# Source: redis/templates/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-release-redis
  namespace: "default"
  labels:
    app.kubernetes.io/name: redis
    helm.sh/chart: redis-14.1.1
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/managed-by: Helm
type: Opaque
data:
  redis-password: "MWczNHZlRk1NUQ=="
`
const TestdeploymentYAMLReadis3 = `
---
# Source: redis/templates/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-release-redis-configuration
  namespace: "default"
  labels:
    app.kubernetes.io/name: redis
    helm.sh/chart: redis-14.1.1
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/managed-by: Helm
data:
  redis.conf: |-
    # User-supplied common configuration:
    # Enable AOF https://redis.io/topics/persistence#append-only-file
    appendonly yes
    # Disable RDB persistence, AOF persistence already enabled.
    save ""
    # End of common configuration
  master.conf: |-
    dir /data
    # User-supplied master configuration:
    rename-command FLUSHDB ""
    rename-command FLUSHALL ""
    # End of master configuration
  replica.conf: |-
    dir /data
    slave-read-only yes
    # User-supplied replica configuration:
    rename-command FLUSHDB ""
    rename-command FLUSHALL ""
    # End of replica configuration
---
`
const TestdeploymentYAMLReadis4 = `
# Source: redis/templates/health-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-release-redis-health
  namespace: "default"
  labels:
    app.kubernetes.io/name: redis
    helm.sh/chart: redis-14.1.1
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/managed-by: Helm
data:
  ping_readiness_local.sh: |-
    #!/bin/bash

    [[ -f $REDIS_PASSWORD_FILE ]] && export REDIS_PASSWORD="$(< "${REDIS_PASSWORD_FILE}")"
    export REDISCLI_AUTH="$REDIS_PASSWORD"
    response=$(
      timeout -s 3 $1 \
      redis-cli \
        -h localhost \
        -p $REDIS_PORT \
        ping
    )
    if [ "$response" != "PONG" ]; then
      echo "$response"
      exit 1
    fi
  ping_liveness_local.sh: |-
    #!/bin/bash

    [[ -f $REDIS_PASSWORD_FILE ]] && export REDIS_PASSWORD="$(< "${REDIS_PASSWORD_FILE}")"
    export REDISCLI_AUTH="$REDIS_PASSWORD"
    response=$(
      timeout -s 3 $1 \
      redis-cli \
        -h localhost \
        -p $REDIS_PORT \
        ping
    )
    if [ "$response" != "PONG" ] && [ "$response" != "LOADING Redis is loading the dataset in memory" ]; then
      echo "$response"
      exit 1
    fi
  ping_readiness_master.sh: |-
    #!/bin/bash

    [[ -f $REDIS_MASTER_PASSWORD_FILE ]] && export REDIS_MASTER_PASSWORD="$(< "${REDIS_MASTER_PASSWORD_FILE}")"
    export REDISCLI_AUTH="$REDIS_MASTER_PASSWORD"
    response=$(
      timeout -s 3 $1 \
      redis-cli \
        -h $REDIS_MASTER_HOST \
        -p $REDIS_MASTER_PORT_NUMBER \
        ping
    )
    if [ "$response" != "PONG" ]; then
      echo "$response"
      exit 1
    fi
  ping_liveness_master.sh: |-
    #!/bin/bash

    [[ -f $REDIS_MASTER_PASSWORD_FILE ]] && export REDIS_MASTER_PASSWORD="$(< "${REDIS_MASTER_PASSWORD_FILE}")"
    export REDISCLI_AUTH="$REDIS_MASTER_PASSWORD"
    response=$(
      timeout -s 3 $1 \
      redis-cli \
        -h $REDIS_MASTER_HOST \
        -p $REDIS_MASTER_PORT_NUMBER \
        ping
    )
    if [ "$response" != "PONG" ] && [ "$response" != "LOADING Redis is loading the dataset in memory" ]; then
      echo "$response"
      exit 1
    fi
  ping_readiness_local_and_master.sh: |-
    script_dir="$(dirname "$0")"
    exit_status=0
    "$script_dir/ping_readiness_local.sh" $1 || exit_status=$?
    "$script_dir/ping_readiness_master.sh" $1 || exit_status=$?
    exit $exit_status
  ping_liveness_local_and_master.sh: |-
    script_dir="$(dirname "$0")"
    exit_status=0
    "$script_dir/ping_liveness_local.sh" $1 || exit_status=$?
    "$script_dir/ping_liveness_master.sh" $1 || exit_status=$?
    exit $exit_status
---
`
const TestdeploymentYAMLReadis5 = `
# Source: redis/templates/scripts-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-release-redis-scripts
  namespace: "default"
  labels:
    app.kubernetes.io/name: redis
    helm.sh/chart: redis-14.1.1
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/managed-by: Helm
data:
  start-master.sh: |
    #!/bin/bash

    [[ -f $REDIS_PASSWORD_FILE ]] && export REDIS_PASSWORD="$(< "${REDIS_PASSWORD_FILE}")"
    if [[ ! -f /opt/bitnami/redis/etc/master.conf ]];then
        cp /opt/bitnami/redis/mounted-etc/master.conf /opt/bitnami/redis/etc/master.conf
    fi
    if [[ ! -f /opt/bitnami/redis/etc/redis.conf ]];then
        cp /opt/bitnami/redis/mounted-etc/redis.conf /opt/bitnami/redis/etc/redis.conf
    fi
    ARGS=("--port" "${REDIS_PORT}")
    ARGS+=("--requirepass" "${REDIS_PASSWORD}")
    ARGS+=("--masterauth" "${REDIS_PASSWORD}")
    ARGS+=("--include" "/opt/bitnami/redis/etc/redis.conf")
    ARGS+=("--include" "/opt/bitnami/redis/etc/master.conf")
    exec redis-server "${ARGS[@]}"
  start-replica.sh: |
    #!/bin/bash

    [[ -f $REDIS_PASSWORD_FILE ]] && export REDIS_PASSWORD="$(< "${REDIS_PASSWORD_FILE}")"
    [[ -f $REDIS_MASTER_PASSWORD_FILE ]] && export REDIS_MASTER_PASSWORD="$(< "${REDIS_MASTER_PASSWORD_FILE}")"
    if [[ ! -f /opt/bitnami/redis/etc/replica.conf ]];then
        cp /opt/bitnami/redis/mounted-etc/replica.conf /opt/bitnami/redis/etc/replica.conf
    fi
    if [[ ! -f /opt/bitnami/redis/etc/redis.conf ]];then
        cp /opt/bitnami/redis/mounted-etc/redis.conf /opt/bitnami/redis/etc/redis.conf
    fi
    ARGS=("--port" "${REDIS_PORT}")
    ARGS+=("--slaveof" "${REDIS_MASTER_HOST}" "${REDIS_MASTER_PORT_NUMBER}")
    ARGS+=("--requirepass" "${REDIS_PASSWORD}")
    ARGS+=("--masterauth" "${REDIS_MASTER_PASSWORD}")
    ARGS+=("--include" "/opt/bitnami/redis/etc/redis.conf")
    ARGS+=("--include" "/opt/bitnami/redis/etc/replica.conf")
    exec redis-server "${ARGS[@]}"
---
`
const TestdeploymentYAMLReadis6 = `
# Source: redis/templates/headless-svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: my-release-redis-headless
  namespace: "default"
  labels:
    app.kubernetes.io/name: redis
    helm.sh/chart: redis-14.1.1
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/managed-by: Helm
spec:
  type: ClusterIP
  clusterIP: None
  ports:
    - name: tcp-redis
      port: 6379
      protocol: TCP
      targetPort: redis
  selector:
    app.kubernetes.io/name: redis
    app.kubernetes.io/instance: my-release
---
`
const TestdeploymentYAMLReadis7 = `
# Source: redis/templates/master/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: my-release-redis-master
  namespace: "default"
  labels:
    app.kubernetes.io/name: redis
    helm.sh/chart: redis-14.1.1
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/component: master
spec:
  type: ClusterIP
  
  ports:
    - name: tcp-redis
      port: 6379
      protocol: TCP
      targetPort: redis
      nodePort: null
  selector:
    app.kubernetes.io/name: redis
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/component: master
---
`
const TestdeploymentYAMLReadis8 = `
# Source: redis/templates/replicas/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: my-release-redis-replicas
  namespace: "default"
  labels:
    app.kubernetes.io/name: redis
    helm.sh/chart: redis-14.1.1
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/component: replica
spec:
  type: ClusterIP
  ports:
    - name: tcp-redis
      port: 6379
      protocol: TCP
      targetPort: redis
      nodePort: null
  selector:
    app.kubernetes.io/name: redis
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/component: replica
---
`
const TestdeploymentYAMLReadis9 = `
# Source: redis/templates/master/statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: my-release-redis-master
  namespace: "default"
  labels:
    app.kubernetes.io/name: redis
    helm.sh/chart: redis-14.1.1
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/component: master
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: redis
      app.kubernetes.io/instance: my-release
      app.kubernetes.io/component: master
  serviceName: my-release-redis-headless
  updateStrategy:
    rollingUpdate: {}
    type: RollingUpdate
  template:
    metadata:
      labels:
        app.kubernetes.io/name: redis
        helm.sh/chart: redis-14.1.1
        app.kubernetes.io/instance: my-release
        app.kubernetes.io/managed-by: Helm
        app.kubernetes.io/component: master
      annotations:
        checksum/configmap: f114e2f92cda8990cd8ed501eb94225a6d99c6d1348be0649c719744fe9fe412
        checksum/health: 96635d1783462a68444484173f6608baaad6adb01eaa4a18aab5ee1da2e3a612
        checksum/scripts: af6456e4449bd4a8fb24cd0f481dffc85f4943b1abe78b3f2d1b33e3f0d04aa1
        checksum/secret: 7da05b585cf096e1e323ab38963644e798edd6f652dba9c4ef3255ea20fb5bd6
    spec:
      
      securityContext:
        fsGroup: 1001
      serviceAccountName: my-release-redis
      affinity:
        podAffinity:
          
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - podAffinityTerm:
                labelSelector:
                  matchLabels:
                    app.kubernetes.io/name: redis
                    app.kubernetes.io/instance: my-release
                    app.kubernetes.io/component: master
                namespaces:
                  - "default"
                topologyKey: kubernetes.io/hostname
              weight: 1
        nodeAffinity:
          
      terminationGracePeriodSeconds: 30
      containers:
        - name: redis
          image: docker.io/bitnami/redis:6.2.3-debian-10-r0
          imagePullPolicy: "IfNotPresent"
          securityContext:
            runAsUser: 1001
          command:
            - /bin/bash
          args:
            - -c
            - /opt/bitnami/scripts/start-scripts/start-master.sh
          env:
            - name: BITNAMI_DEBUG
              value: "false"
            - name: REDIS_REPLICATION_MODE
              value: master
            - name: ALLOW_EMPTY_PASSWORD
              value: "no"
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: my-release-redis
                  key: redis-password
            - name: REDIS_TLS_ENABLED
              value: "no"
            - name: REDIS_PORT
              value: "6379"
          ports:
            - name: redis
              containerPort: 6379
              protocol: TCP
          livenessProbe:
            initialDelaySeconds: 5
            periodSeconds: 5
            # One second longer than command timeout should prevent generation of zombie processes.
            timeoutSeconds: 6
            successThreshold: 1
            failureThreshold: 5
            exec:
              command:
                - sh
                - -c
                - /health/ping_liveness_local.sh 5
          readinessProbe:
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 2
            successThreshold: 1
            failureThreshold: 5
            exec:
              command:
                - sh
                - -c
                - /health/ping_readiness_local.sh 1
          resources:
            limits: {}
            requests: {}
          volumeMounts:
            - name: start-scripts
              mountPath: /opt/bitnami/scripts/start-scripts
            - name: health
              mountPath: /health
            - name: redis-data
              mountPath: /data
              subPath: 
            - name: config
              mountPath: /opt/bitnami/redis/mounted-etc
            - name: redis-tmp-conf
              mountPath: /opt/bitnami/redis/etc/
            - name: tmp
              mountPath: /tmp
      volumes:
        - name: start-scripts
          configMap:
            name: my-release-redis-scripts
            defaultMode: 0755
        - name: health
          configMap:
            name: my-release-redis-health
            defaultMode: 0755
        - name: config
          configMap:
            name: my-release-redis-configuration
        - name: redis-tmp-conf
          emptyDir: {}
        - name: tmp
          emptyDir: {}
        - name: redis-data
          emptyDir: {}
---
`
const TestdeploymentYAMLReadis11 = `
# Source: redis/templates/replicas/statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: my-release-redis-replicas
  namespace: "default"
  labels:
    app.kubernetes.io/name: redis
    helm.sh/chart: redis-14.1.1
    app.kubernetes.io/instance: my-release
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/component: replica
spec:
  replicas: 3
  selector:
    matchLabels:
      app.kubernetes.io/name: redis
      app.kubernetes.io/instance: my-release
      app.kubernetes.io/component: replica
  serviceName: my-release-redis-headless
  updateStrategy:
    rollingUpdate: {}
    type: RollingUpdate
  template:
    metadata:
      labels:
        app.kubernetes.io/name: redis
        helm.sh/chart: redis-14.1.1
        app.kubernetes.io/instance: my-release
        app.kubernetes.io/managed-by: Helm
        app.kubernetes.io/component: replica
      annotations:
        checksum/configmap: f114e2f92cda8990cd8ed501eb94225a6d99c6d1348be0649c719744fe9fe412
        checksum/health: 96635d1783462a68444484173f6608baaad6adb01eaa4a18aab5ee1da2e3a612
        checksum/scripts: af6456e4449bd4a8fb24cd0f481dffc85f4943b1abe78b3f2d1b33e3f0d04aa1
        checksum/secret: 0d5f63c1215d3f7d257bd3e3e937ff812df088aad174ef444bf8f7e22957c0a4
    spec:
      
      securityContext:
        fsGroup: 1001
      serviceAccountName: my-release-redis
      affinity:
        podAffinity:
          
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - podAffinityTerm:
                labelSelector:
                  matchLabels:
                    app.kubernetes.io/name: redis
                    app.kubernetes.io/instance: my-release
                    app.kubernetes.io/component: replica
                namespaces:
                  - "default"
                topologyKey: kubernetes.io/hostname
              weight: 1
        nodeAffinity:
          
      terminationGracePeriodSeconds: 30
      containers:
        - name: redis
          image: docker.io/bitnami/redis:6.2.3-debian-10-r0
          imagePullPolicy: "IfNotPresent"
          securityContext:
            runAsUser: 1001
          command:
            - /bin/bash
          args:
            - -c
            - /opt/bitnami/scripts/start-scripts/start-replica.sh
          env:
            - name: BITNAMI_DEBUG
              value: "false"
            - name: REDIS_REPLICATION_MODE
              value: slave
            - name: REDIS_MASTER_HOST
              value: my-release-redis-master-0.my-release-redis-headless.default.svc.cluster.local
            - name: REDIS_MASTER_PORT_NUMBER
              value: "6379"
            - name: ALLOW_EMPTY_PASSWORD
              value: "no"
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: my-release-redis
                  key: redis-password
            - name: REDIS_MASTER_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: my-release-redis
                  key: redis-password
            - name: REDIS_TLS_ENABLED
              value: "no"
            - name: REDIS_PORT
              value: "6379"
          ports:
            - name: redis
              containerPort: 6379
              protocol: TCP
          livenessProbe:
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 6
            successThreshold: 1
            failureThreshold: 5
            exec:
              command:
                - sh
                - -c
                - /health/ping_liveness_local_and_master.sh 5
          readinessProbe:
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 2
            successThreshold: 1
            failureThreshold: 5
            exec:
              command:
                - sh
                - -c
                - /health/ping_readiness_local_and_master.sh 1
          resources:
            limits: {}
            requests: {}
          volumeMounts:
            - name: start-scripts
              mountPath: /opt/bitnami/scripts/start-scripts
            - name: health
              mountPath: /health
            - name: redis-data
              mountPath: /data
              subPath: 
            - name: config
              mountPath: /opt/bitnami/redis/mounted-etc
            - name: redis-tmp-conf
              mountPath: /opt/bitnami/redis/etc
      volumes:
        - name: start-scripts
          configMap:
            name: my-release-redis-scripts
            defaultMode: 0755
        - name: health
          configMap:
            name: my-release-redis-health
            defaultMode: 0755
        - name: config
          configMap:
            name: my-release-redis-configuration
        - name: redis-tmp-conf
          emptyDir: {}
        - name: redis-data
          emptyDir: {}
#
#NOTES:
#** Please be patient while the chart is being deployed **
#
#Redis(TM) can be accessed on the following DNS names from within your cluster:
#
#    my-release-redis-master.default.svc.cluster.local for read/write operations (port 6379)
#    my-release-redis-replicas.default.svc.cluster.local for read-only operations (port 6379)
#
#
#
#To get your password run:
#
#    export REDIS_PASSWORD=$(kubectl get secret --namespace default my-release-redis -o jsonpath="{.data.redis-password}" | base64 --decode)
#
#To connect to your Redis(TM) server:
#
#1. Run a Redis(TM) pod that you can use as a client:
#
#   kubectl run --namespace default redis-client --restart='Never'  --env REDIS_PASSWORD=$REDIS_PASSWORD  --image docker.io/bitnami/redis:6.2.3-debian-10-r0 --command -- sleep infinity
#
#   Use the following command to attach to the pod:
#
#   kubectl exec --tty -i redis-client \
#   --namespace default -- bash
#
#2. Connect using the Redis(TM) CLI:
#   redis-cli -h my-release-redis-master -a $REDIS_PASSWORD
#   redis-cli -h my-release-redis-replicas -a $REDIS_PASSWORD
#
#To connect to your database from outside the cluster execute the following commands:
#
#    kubectl port-forward --namespace default svc/my-release-redis-master 6379:6379 &
#    redis-cli -h 127.0.0.1 -p 6379 -a $REDIS_PASSWORD
`

