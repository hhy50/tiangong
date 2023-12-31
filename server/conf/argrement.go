package conf

import "tiangong/common"

var (
	SrvPortDef = common.Pair[string, int]{
		First:  "srvPort",
		Second: 2023,
	}

	HttpPortDef = common.Pair[string, int]{
		First:  "httpPort",
		Second: 2024,
	}

	UserNameDef = common.Pair[string, string]{
		First:  "userName",
		Second: "admin",
	}

	PasswdDef = common.Pair[string, string]{
		First:  "passwd",
		Second: "",
	}
)
