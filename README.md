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
