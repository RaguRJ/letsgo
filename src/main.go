/*
This cli app does the following
1. Creates bin/, src/ and pkg/ directories if not present in current directory if flag is present
2. exports GOPATH and PATH commands to use
3. Update the zshrc file if flag is present
*/

package main

import "letsgo"

func main() {
	letsgo.Execute()
}
