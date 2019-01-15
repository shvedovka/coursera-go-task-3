package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	seenBrowsers := make(map[string]bool)
	uniqueBrowsers := 0
	var foundUsers []string

	var i = -1
	//line := make([]byte, 0, 1024)
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			break
		}

		user := make(map[string]interface{})
		err = json.Unmarshal(line, &user)

		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			// log.Println("cant cast browsers")
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}

			curIsAndroid := strings.Contains(browser, "Android")
			curIsMSIE := strings.Contains(browser, "MSIE")
			isMSIE = isMSIE || curIsMSIE
			isAndroid = isAndroid || curIsAndroid

			if _, found := seenBrowsers[browser]; found == true {
				continue
			}

			if curIsAndroid || curIsMSIE {
				seenBrowsers[browser] = true
				uniqueBrowsers++
			}
			isMSIE = isMSIE || curIsMSIE
			isAndroid = isAndroid || curIsAndroid

		}
		i++
		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user["email"].(string), "@", " [at] ", 1)
		foundUsers = append(foundUsers, fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email))
	}

	fmt.Fprintln(out, "found users:\n"+strings.Join(foundUsers, ""))
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
