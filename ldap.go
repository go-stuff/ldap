package ldap

import (
	"crypto/tls"
	"fmt"
	"strings"

	"gopkg.in/ldap.v3"
)

// Auth will search with LDAP a given user, authenticate that user and return
// the username as a string in all lowercase, groups as a slice of strings
// all in lowercase or an error.
func Auth(
	server string,
	port string,
	bindDN string,
	bindPass string,
	userBaseDN string,
	userSearchAttr string,
	groupBaseDN string,
	groupObjectClass string,
	groupSearchAttr string,
	groupSearchFull string,
	user string,
	pass string,
) (string, []string, error) {
	var err error
	var username string
	var groups []string

	username = strings.ToLower(user)

	// connect to ldap
	con, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%s", server, port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return username, nil, err
	}
	defer con.Close()

	// use a binding account to collect information
	err = con.Bind(bindDN, bindPass)
	if err != nil {
		return username, nil, err
	}

	// bind the user account to authenticate
	err = con.Bind(fmt.Sprintf("%s=%s,%s", userSearchAttr, user, userBaseDN), pass)
	if err != nil {
		return username, nil, err
	}

	// full group search filter or simple search filter
	var groupFilter string
	if groupSearchFull == "true" {
		groupFilter = fmt.Sprintf("(&(objectClass=%s)(%s=%s=%s,%s))", groupObjectClass, groupSearchAttr, userSearchAttr, user, userBaseDN)
	} else {
		groupFilter = fmt.Sprintf("(&(objectClass=%s)(%s=%s))", groupObjectClass, groupSearchAttr, user)
	}

	// search groups request
	searchGroups := ldap.NewSearchRequest(
		groupBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		groupFilter,
		[]string{"cn"},
		nil,
	)

	// search groups
	groupsResult, err := con.Search(searchGroups)
	if err != nil {
		return username, nil, err
	}

	// get group values
	for _, group := range groupsResult.Entries {
		for _, attribute := range group.Attributes {
			groups = append(groups, strings.ToLower(attribute.Values[0]))
		}
	}

	return username, groups, nil
}
