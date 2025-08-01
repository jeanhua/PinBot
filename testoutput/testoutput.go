package testoutput

import "os"

func Output(s string) {
	os.WriteFile("./temp.txt", []byte(s), 0755)
}
