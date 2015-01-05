package log4go

import (
	"os"
	"path"
)

func (log Logger) LoadConfig(filename string) {
	if len(filename) <= 0 {
		return
	}
	ext := path.Ext(filename)
	ext = ext[1:]

	switch ext {
	case "xml":
		log.LoadXMLConfig(filename)
		break
	case "json":
		//log.LoadJSONConfig(filename)
		break
	default:
		os.Stderr.WriteString("warning: unknow config file type, only XML and JSON are supported\n")
	}
}
