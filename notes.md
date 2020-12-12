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

Would like to be able to

* Provide a rego policy (e.g. allow { path !contains .. })
* Provide query for what should be set as haproxy vars (e.g. allow)
* Provide input from arbitrary spoe event args (e.g. path)

What to do with additional results? SPOP does not support arrays. Could maybe express in string
somehow? Or just only expect one result.

Binary argument format is a kind of catch all. Can pass complex data like headers (see req.hdrs_bin)
to agent.

Interfaces in go take some getting used to. Instead of looking for an appropriate interface, you
might simply define one that matches what you care about, since callers don't have to even know it
exists.

https://github.com/golang/go/wiki/CodeReviewComments
https://dave.cheney.net/resources-for-new-go-programmers
https://www.honeybadger.io/blog/golang-logging/
https://dave.cheney.net/2015/11/05/lets-talk-about-logging

## CLI

kubernetes uses combination of pflag and cobra – looks like moving towards cobra.
docker uses cobra
