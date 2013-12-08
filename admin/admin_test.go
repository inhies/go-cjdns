package admin

import "testing"

/*
import (
	"fmt"
	"reflect"
)
*/

var c *Conn

func TestConnect(t *testing.T) {
	var err error
	c, err = Connect(nil)
	if err != nil {
		t.Fatal("Failed to connect,", err)
	}
}

func TestPing(t *testing.T) {
	if err := c.Ping(); err != nil {
		t.Error("Failed to ping,", err)
	}
}

func TestMemory(t *testing.T) {
	_, err := c.Memory()
	if err != nil {
		t.Error(err)
	}
}

func TestCookie(t *testing.T) {
	_, err := c.cookie()
	if err != nil {
		t.Error(err)
	}
}

func TestAuth(t *testing.T) {
	err := c.AuthedPing()
	if err != nil {
		t.Error(err)
	}
}

func TestAvailableFunctions(t *testing.T) {
	_, err := c.Admin_availableFunctions()
	if err != nil {
		t.Error(err)
	}
	/*
		fmt.Println(reflect.Value(testing))
		for cmd, args := range funcs {
			fmt.Printf("// %s(", cmd)
			if len(args) > 0 {
				for k, arg := range args {
					fmt.Printf(" %s %s, ", k, arg.Type)
					//if arg.Required {
					//	fmt.Print("(required)")
					//}
				}
			}
			fmt.Println(")")
		}
	*/
}

func TestAdmin_asyncEnabled(t *testing.T) {
	_, err := c.Admin_asyncEnabled()
	if err != nil {
		t.Error("Admin_asyncEnabled failed,", err)
	}
}

func TestNodeStore_dumpTable(t *testing.T) {
	_, err := c.NodeStore_dumpTable()
	if err != nil {
		t.Error("NodeStore_dumpTable failed,", err)
	}
}

func TestInterfaceController_peerStats(t *testing.T) {
	_, err := c.InterfaceController_peerStats()
	if err != nil {
		t.Error("InterfaceController_peerStats failed,", err)
	}
}

func TestSecurity_noFiles(t *testing.T) {
	_, err := c.Security_noFiles()
	if err != nil {
		t.Error("Security_noFiles failed,", err)
	}
}

func TestAuthorizedPasswords(t *testing.T) {
	user := "dickweed"
	pass := "hackme"

	if err := c.AuthorizedPasswords_add(user, pass, 0); err != nil {
		t.Error("failed to add password to cjdns,", err)
		return
	}

	users, err := c.AuthorizedPasswords_list()
	if err != nil {
		t.Error("failed to get list of password users from cjdns,", err)
		return
	}

	var found bool
	for _, u := range users {
		if u == user {
			found = true
			break
		}
	}
	if !found {
		t.Error("previously added user not found in users list")
		return
	}

	err = c.AuthorizedPasswords_remove(user)
	if err != nil {
		t.Error("failed to remove password for user,", err)
		return
	}

	users, err = c.AuthorizedPasswords_list()
	if err != nil {
		t.Error("failed to get list of password users from cjdns", err)
		return
	}

	found = false
	for _, u := range users {
		if u == user {
			found = true
			break
		}
	}
	if found {
		t.Error("previously removed user still found in users list")
	}
}
