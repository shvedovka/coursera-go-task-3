package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

var dataPool = sync.Pool{
	New: func() interface{} {
		return new(User)
	},
}

//easyjson:json
type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

func FastSearch(out io.Writer) {
	var i = -1
	var unique = 0
	var seen = make(map[string]bool, 256)
	var users = make([]string, 0, 100)

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(out, "found users:")

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			break
		}

		var user = dataPool.Get().(*User)

		err = user.UnmarshalJSON(line)
		if err != nil {
			println(err)
		}
		dataPool.Put(user)

		var userAndroid = false
		var userMSIE = false

		for _, browser := range user.Browsers {
			isBrowserSeen, found := seen[browser]

			var curMSIE = false
			var curAndroid = false

			if !found {
				curMSIE = strings.Contains(browser, "MSIE")
				curAndroid = strings.Contains(browser, "Android")
				seen[browser] = curMSIE || curAndroid
			} else {
				if isBrowserSeen {
					curMSIE = strings.Contains(browser, "MSIE")
					curAndroid = strings.Contains(browser, "Android")
				} else {
					curMSIE = false
					curAndroid = false
				}
			}

			userMSIE = userMSIE || curAndroid
			userAndroid = userAndroid || curMSIE

			if !found && (curMSIE || curAndroid) {
				unique++
			}
		}

		i++
		if !(userMSIE && userAndroid) {
			continue
		}

		var email = strings.Split(user.Email, "@")
		users = append(users, "["+strconv.Itoa(i)+"] "+user.Name+" <"+email[0]+" [at] "+email[1]+">")
	}

	fmt.Fprintln(out, strings.Join(users, "\n"))
	fmt.Fprintln(out, "\nTotal unique browsers", unique)
}
