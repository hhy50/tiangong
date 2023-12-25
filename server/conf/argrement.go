package conf

import "tiangong/common"

var (
	TcpPort = common.Pair[string, int]{
		First:  "tcpPort",
		Second: 2024,
	}

	HttpPort = common.Pair[string, int]{
		First:  "httpPort",
		Second: 2023,
	}

	UserName = common.Pair[string, string]{
		First:  "userName",
		Second: "admin",
	}

	Passwd = common.Pair[string, string]{
		First:  "passwd",
		Second: "",
	}
)
