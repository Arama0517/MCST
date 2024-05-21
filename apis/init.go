package apis

import "github.com/Arama-Vanarana/MCSCS-Go/lib"

func InitApis() {
	configs, err :=	lib.LoadConfigs()
	if err != nil {
		panic(err)
	}
	switch configs.API {
	case 0:
		initFastMirror()
	case 1:
		initPolars()
	default:
		panic("Invalid API selected")
	}
}