# ldap

[![Go Report Card](https://goreportcard.com/badge/github.com/go-stuff/ldap)](https://goreportcard.com/report/github.com/go-stuff/ldap)
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
The reason there are so many variables is to allow for connecting to multiple environments, it has been tested against OpenLDAP and Active Directory, there are some minor differences in objectClass and attributes.

```go
package main

import (
    "fmt"

    "github.com/go-stuff/ldap"
)

// OpenLDAP
const (
    LDAP_SERVER             string = "192.168.1.100"
    LDAP_PORT               int    = 636
    LDAP_BIND_DN            string = "cn=admin,dc=go-stuff,dc=ca"
    LDAP_BIND_PASS          string = "password"
    LDAP_USER_BASE_DN       string = "ou=people,dc=go-stuff,dc=ca"
    LDAP_USER_OBJECT_CLASS  string = "person"
    LDAP_USER_SEARCH_ATTR   string = "uid"
    LDAP_GROUP_BASE_DN      string = "ou=group,dc=go-stuff,dc=ca"
    LDAP_GROUP_OBJECT_CLASS string = "posixGroup"
    LDAP_GROUP_SEARCH_ATTR  string = "memberUid"
    LDAP_GROUP_SEARCH_FULL  bool   = "false"
    LDAP_AUTH_ATTR          string = "uid"
)

// Active Dreictory
// const (
//     LDAP_SERVER             string = "LDAPSSL"
//     LDAP_PORT               int    = 636
//     LDAP_BIND_DN            string = "CN=admin,OU=Users,DC=go-stuff,DC=ca"
//     LDAP_BIND_PASS          string = "password"
//     LDAP_USER_BASE_DN       string = "OU=Users,DC=go-stuff,DC=ca"
//     LDAP_USER_OBJECT_CLASS  string = "person"
//     LDAP_USER_SEARCH_ATTR   string = "CN"
//     LDAP_GROUP_BASE_DN      string = "OU=Groups,DC=go-stuff,DC=ca"
//     LDAP_GROUP_OBJECT_CLASS string = "group"
//     LDAP_GROUP_SEARCH_ATTR  string = "member"
//     LDAP_GROUP_SEARCH_FULL  bool   = true
//     LDAP_AUTH_ATTR          string = "CN"
// )

func main() {
    username, groups, err := ldap.Auth(
        LDAP_SERVER,
        LDAP_PORT,
        LDAP_BIND_DN,
        LDAP_BIND_PASS,
        LDAP_USER_BASE_DN,
        LDAP_USER_OBJECT_CLASS,
        LDAP_USER_SEARCH_ATTR,
        LDAP_GROUP_BASE_DN,
        LDAP_GROUP_OBJECT_CLASS,
        LDAP_GROUP_SEARCH_ATTR,
        LDAP_AUTH_ATTR,
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