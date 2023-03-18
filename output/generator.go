package output

import "net/http"

func MustGetTemplate(escFSPath string) string {
	// tack. decoupled.
	return _escFSMustString(false, escFSPath)
}

func MustGetByte(escFSPath string) []byte {
	return _escFSMustByte(false, escFSPath)
}

func MustFs() http.FileSystem {
	return _escFS(false)
}
