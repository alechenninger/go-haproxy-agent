package example

default allow = false                               # unless otherwise defined, allow is false
default role = "anonymous"

allow = true {                                      # allow is true if...
    request_allowed
}

role {
    input.headers.x-user
}

# permissions := {
#     "anonymous": ["read"],
#     "authenticated": ["read", "write"]
# }

request_allowed {
    input.method == "GET"
    # permissions[]
}