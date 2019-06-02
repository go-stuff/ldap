package ldap_test

import (
	"os"
	"strings"
	"testing"

	"github.com/go-stuff/ldap"
)

func TestMain(m *testing.M) {
	// connect to database
	testsSetup()

	// run all tests
	ret := m.Run()

	// disconnect from database
	testsTeardown()

	// call flag.Parse() here if TestMain uses flags
	os.Exit(ret)
}

func testsSetup() {
	// if environment variable is does not exist or is empty set a default
	if os.Getenv("LDAP_SERVER") == "" {
		os.Setenv("LDAP_SERVER", "192.168.1.100")
	}
	if os.Getenv("LDAP_PORT") == "" {
		os.Setenv("LDAP_PORT", "636")
	}
	if os.Getenv("LDAP_BIND_DN") == "" {
		os.Setenv("LDAP_BIND_DN", "cn=admin,dc=go-stuff,dc=ca")
	}
	if os.Getenv("LDAP_BIND_PASS") == "" {
		os.Setenv("LDAP_BIND_PASS", "password")
	}
	if os.Getenv("LDAP_USER_BASE_DN") == "" {
		os.Setenv("LDAP_USER_BASE_DN", "ou=people,dc=go-stuff,dc=ca")
	}
	if os.Getenv("LDAP_USER_SEARCH_ATTR") == "" {
		os.Setenv("LDAP_USER_SEARCH_ATTR", "uid")
	}
	if os.Getenv("LDAP_GROUP_BASE_DN") == "" {
		os.Setenv("LDAP_GROUP_BASE_DN", "ou=group,dc=go-stuff,dc=ca")
	}
	if os.Getenv("LDAP_GROUP_OBJECT_CLASS") == "" {
		os.Setenv("LDAP_GROUP_OBJECT_CLASS", "posixGroup")
	}
	if os.Getenv("LDAP_GROUP_SEARCH_ATTR") == "" {
		os.Setenv("LDAP_GROUP_SEARCH_ATTR", "memberUid")
	}
	if os.Getenv("LDAP_GROUP_SEARCH_FULL") == "" {
		os.Setenv("LDAP_GROUP_SEARCH_FULL", "false")
	}

	if os.Getenv("LDAP_TEST_USER") == "" {
		os.Setenv("LDAP_TEST_USER", "web-admin")
	}
	if os.Getenv("LDAP_TEST_PASSWORD") == "" {
		os.Setenv("LDAP_TEST_PASSWORD", "password")
	}
}

func testsTeardown() {
}

func TestAuthWithOpenLDAP(t *testing.T) {

	// test an invalid connection
	username, groups, err := ldap.Auth(
		"server-doesnotexit",
		os.Getenv("LDAP_PORT"),
		os.Getenv("LDAP_BIND_DN"),
		os.Getenv("LDAP_BIND_PASS"),
		os.Getenv("LDAP_USER_BASE_DN"),
		os.Getenv("LDAP_USER_SEARCH_ATTR"),
		os.Getenv("LDAP_GROUP_BASE_DN"),
		os.Getenv("LDAP_GROUP_OBJECT_CLASS"),
		os.Getenv("LDAP_GROUP_SEARCH_ATTR"),
		os.Getenv("LDAP_GROUP_SEARCH_FULL"),
		os.Getenv("LDAP_TEST_USER"),
		os.Getenv("LDAP_TEST_PASSWORD"),
	)
	if err == nil {
		t.Fatal("expected to fail server does not exist")
	}

	// test an invalid bind account
	username, groups, err = ldap.Auth(
		os.Getenv("LDAP_SERVER"),
		os.Getenv("LDAP_PORT"),
		"binddn-doesnotexist",
		os.Getenv("LDAP_BIND_PASS"),
		os.Getenv("LDAP_USER_BASE_DN"),
		os.Getenv("LDAP_USER_SEARCH_ATTR"),
		os.Getenv("LDAP_GROUP_BASE_DN"),
		os.Getenv("LDAP_GROUP_OBJECT_CLASS"),
		os.Getenv("LDAP_GROUP_SEARCH_ATTR"),
		os.Getenv("LDAP_GROUP_SEARCH_FULL"),
		os.Getenv("LDAP_TEST_USER"),
		os.Getenv("LDAP_TEST_PASSWORD"),
	)
	if err == nil {
		t.Fatal("expected to fail binddn does not exist")
	}

	// test an invalid user search
	username, groups, err = ldap.Auth(
		os.Getenv("LDAP_SERVER"),
		os.Getenv("LDAP_PORT"),
		os.Getenv("LDAP_BIND_DN"),
		os.Getenv("LDAP_BIND_PASS"),
		"baseDN-doesnotexist",
		"attribute-doesnotexist",
		os.Getenv("LDAP_GROUP_BASE_DN"),
		os.Getenv("LDAP_GROUP_OBJECT_CLASS"),
		os.Getenv("LDAP_GROUP_SEARCH_ATTR"),
		os.Getenv("LDAP_GROUP_SEARCH_FULL"),
		os.Getenv("LDAP_TEST_USER"),
		os.Getenv("LDAP_TEST_PASSWORD"),
	)
	if err == nil {
		t.Fatal("expected to fail invalid user search")
	}

	// test an invalid group search
	username, groups, err = ldap.Auth(
		os.Getenv("LDAP_SERVER"),
		os.Getenv("LDAP_PORT"),
		os.Getenv("LDAP_BIND_DN"),
		os.Getenv("LDAP_BIND_PASS"),
		os.Getenv("LDAP_USER_BASE_DN"),
		os.Getenv("LDAP_USER_SEARCH_ATTR"),
		"baseDN-doesnotexist",
		"objectClass-doesnotexist",
		"attribute-doesnotexist",
		os.Getenv("LDAP_GROUP_SEARCH_FULL"),
		os.Getenv("LDAP_TEST_USER"),
		os.Getenv("LDAP_TEST_PASSWORD"),
	)
	if err == nil {
		t.Fatal("expected to fail invalid group search")
	}

	// test lowercase username
	if username != strings.ToLower(os.Getenv("LDAP_TEST_USER")) {
		t.Fatal("unexpected username")
	}

	// test group full search
	username, groups, err = ldap.Auth(
		os.Getenv("LDAP_SERVER"),
		os.Getenv("LDAP_PORT"),
		os.Getenv("LDAP_BIND_DN"),
		os.Getenv("LDAP_BIND_PASS"),
		os.Getenv("LDAP_USER_BASE_DN"),
		os.Getenv("LDAP_USER_SEARCH_ATTR"),
		os.Getenv("LDAP_GROUP_BASE_DN"),
		os.Getenv("LDAP_GROUP_OBJECT_CLASS"),
		os.Getenv("LDAP_GROUP_SEARCH_ATTR"),
		"true",
		os.Getenv("LDAP_TEST_USER"),
		os.Getenv("LDAP_TEST_PASSWORD"),
	)
	if len(groups) != 0 {
		t.Fatal("expected to return 0 groups in OpenLDAP")
	}

	// // test a valid user
	// username, groups, err = ldap.Auth(
	// 	os.Getenv("LDAP_SERVER"),
	// 	os.Getenv("LDAP_PORT"),
	// 	os.Getenv("LDAP_BIND_DN"),
	// 	os.Getenv("LDAP_BIND_PASS"),
	// 	os.Getenv("LDAP_USER_BASE_DN"),
	// 	os.Getenv("LDAP_USER_SEARCH_ATTR"),
	// 	os.Getenv("LDAP_GROUP_BASE_DN"),
	// 	os.Getenv("LDAP_GROUP_OBJECT_CLASS"),
	// 	os.Getenv("LDAP_GROUP_SEARCH_ATTR"),
	// 	os.Getenv("LDAP_GROUP_SEARCH_FULL"),
	// 	os.Getenv("LDAP_TEST_USER"),
	// 	os.Getenv("LDAP_TEST_PASSWORD"),
	// )
	// if username != os.Getenv("LDAP_TEST_USER") ||
	// 	len(groups) != 2 ||
	// 	groups[0] != "domain users" ||
	// 	groups[1] != "group-admin" {
	// 	t.Fatal("did not get expected results")
	// }

	// test a user that does not exist
	username, groups, err = ldap.Auth(
		os.Getenv("LDAP_SERVER"),
		os.Getenv("LDAP_PORT"),
		os.Getenv("LDAP_BIND_DN"),
		os.Getenv("LDAP_BIND_PASS"),
		os.Getenv("LDAP_USER_BASE_DN"),
		os.Getenv("LDAP_USER_SEARCH_ATTR"),
		os.Getenv("LDAP_GROUP_BASE_DN"),
		os.Getenv("LDAP_GROUP_OBJECT_CLASS"),
		os.Getenv("LDAP_GROUP_SEARCH_ATTR"),
		os.Getenv("LDAP_GROUP_SEARCH_FULL"),
		"web-doesnotexist",
		os.Getenv("LDAP_TEST_PASSWORD"),
	)
	if err == nil {
		t.Fatal("expected to fail user does not exist")
	}

}
