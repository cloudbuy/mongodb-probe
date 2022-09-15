# mongodb-probe

Small app replacing the functionality of the liveness/readiness probes in the Bitnami mongodb chart.

## Building

```
go build
```

## Usage

```
$ ./mongodb-probe liveness
2022/09/15 12:17:17 attempting to connect
2022/09/15 12:17:17 running ping
2022/09/15 12:17:17 ping successful
```

```
$ ./mongodb-probe readiness
2022/09/15 12:17:12 attempting to connect
2022/09/15 12:17:12 running hello in admin database
2022/09/15 12:17:12 hello successful
```

## Using with the Bitnami mongodb chart

To use this probe in place of the existing shell + mongosh probes, the following addition to the
mongodb values can be used:

```yaml
extraVolumes:
- name: mongodb-probe
  emptyDir: {}
extraVolumeMounts:
- name: mongodb-probe
  mountPath: /opt/probes
initContainers:
- name: fetch-probe
  image: busybox:latest
  command:
  - /bin/sh
  args:
  - -ec
  - |
    /bin/wget -O /opt/probes/mongodb-probe https://github.com/cloudbuy/mongodb-probe/releases/download/v0.0.1/mongodb-probe
    chmod +x /opt/probes/mongodb-probe
  volumeMounts:
  - name: mongodb-probe
    mountPath: /opt/probes
  
livenessProbe:
  enabled: false
customLivenessProbe: |
  exec:
    command:
      - /opt/probes/mongodb-probe
      - liveness
  failureThreshold: 6
  initialDelaySeconds: 30
  periodSeconds: 20
  timeoutSeconds: 10
  successThreshold: 1
readinessProbe:
  enabled: false
customReadinessProbe: |
  exec:
    command:
      - /opt/probes/mongodb-probe
      - readiness
  failureThreshold: 6
  initialDelaySeconds: 30
  periodSeconds: 20
  timeoutSeconds: 10
  successThreshold: 1
customStartupProbe: |
  exec:
    command:
      - /opt/probes/mongodb-probe
      - startup
  failureThreshold: 6
  initialDelaySeconds: 30
  periodSeconds: 20
  timeoutSeconds: 10
  successThreshold: 1
```