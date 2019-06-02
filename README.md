# ldap

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![Gopher Share](https://github.com/go-stuff/images/blob/master/GOPHER_SHARE_640x320.png)

Using [go-ldap.v3](https://github.com/go-ldap/ldap) to authenticate with LDAP and return the username and groups associated witht that user. An error is returned if authentication fails.

## Packages Imported

- Basic LDAP [github.com/go-ldap/ldap](https://github.com/go-ldap/ldap)

## Installation

The recommended way to get started using [github.com/go-stuff/ldap](https://github.com/go-stuff/ldap) is by using 'go get' to install the dependency in your project.

```go
go get "github.com/go-stuff/ldap"
```

## Usage

```go
import (
    "github.com/go-stuff/ldap"
)
```

## Example

This is an example of how it would be implemented. Of course the constants could be environment variables or in a configuration file, etc... this is just an example.

```go
package main

import (
    "fmt"

    "github.com/go-stuff/ldap"
)

const (
    LDAP_SERVER        string = "192.168.1.100"
    LDAP_PORT          int    = 636
    LDAP_BIND_DN       string = "cn=admin,dc=go-stuff,dc=ca"
    LDAP_BIND_PASS     string = "password"
    LDAP_USER_BASE_DN  string = "ou=people,dc=go-stuff,dc=ca"
    LDAP_GROUP_BASE_DN string = "ou=group,dc=go-stuff,dc=ca"
)

func main() {
    username, groups, err := ldap.Auth(
        LDAP_SERVER,
        LDAP_PORT,
        LDAP_BIND_DN,
        LDAP_BIND_PASS,
        LDAP_USER_BASE_DN,
        LDAP_GROUP_BASE_DN,
        "web-user",
        "password",
    )
    fmt.Printf("Username: %s\n", username)
    if err != nil {
        fmt.Println(err.Error())
    }
    for _, v := range groups {
        fmt.Printf("    Group: %s\n", v)
    }
```

### Example Output

```bash
Username: web-user
    Group: domain users
    Group: group-user
    Group: group-random1
    Group: group-random3
```

## License

[MIT License](LICENSE)