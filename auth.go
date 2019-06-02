package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/ldap.v3"
)

func Auth(
	server string,
	port int,
	bindDN string,
	bindPass string,
	userBaseDN string,
	userObjectClass string,
	userSearchAttr string,
	groupBaseDN string,
	groupObjectClass string,
	groupSearchAttr string,
	authAttr string,
	user string,
	pass string,
) (string, []string, error) {
	var err error
	var username string
	var groups []string

	username = strings.ToLower(user)

	// connect to ldap
	con, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", server, port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return username, nil, err
	}
	defer con.Close()

	// use a binding account to collect information
	err = con.Bind(bindDN, bindPass)
	if err != nil {
		return username, nil, err
	}

	// search users request
	searchUsers := ldap.NewSearchRequest(
		userBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=%s)(%s=%s))", userObjectClass, userSearchAttr, user),
		[]string{"cn"},
		nil,
	)

	// search users
	usersResult, err := con.Search(searchUsers)
	if err != nil {
		return username, nil, err
	}

	// no user returned
	if len(usersResult.Entries) == 0 {
		return username, nil, errors.New("user does not exist")
	}

	// more than one user returned
	if len(usersResult.Entries) > 1 {
		return username, nil, errors.New("too many users returned")
	}

	// search groups request
	searchGroups := ldap.NewSearchRequest(
		groupBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=%s)(%s=%s))", groupObjectClass, groupSearchAttr, user),
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
			fmt.Printf("Group Attribute: %v, Values: %v\n", attribute.Name, attribute.Values)
			groups = append(groups, strings.ToLower(attribute.Values[0]))
		}
	}

	// bind the user account to authenticate
	err = con.Bind(fmt.Sprintf("%s=%s,%s", authAttr, user, userBaseDN), pass)
	if err != nil {
		return username, nil, err
	}

	return username, groups, nil
}
