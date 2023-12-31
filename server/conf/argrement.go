package conf

import "tiangong/common"

var (
	HostDef = common.Pair[string, string]{
		First:  "host",
		Second: "0.0.0.0",
	}

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
