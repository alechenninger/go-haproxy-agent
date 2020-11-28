Tried criteo agent, but didn't work:

```
$ go run ./
INFO[0000] spoe: listening on [::]:9000                 
2020/11/28 11:46:26 message received

ERRO[0029] spoe: error handling connection: disconnect error: a timeout occurred 
```

Docker for mac supports host network differently. Use host.docker.internal hostname to refer to host
network.

go mod tidy – removes unused modules, adds transitively used modules