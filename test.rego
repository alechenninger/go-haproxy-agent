package example

default allow = false                               # unless otherwise defined, allow is false

allow = true {                                      # allow is true if...
    request_allowed
}

user := input.header["User"][0]
roles := data.users[user]

request_allowed {
    input.method == "GET"
    roles[_] == "read"
}

request_allowed {
    not input.method == "GET"
    roles[_] == "write"
}