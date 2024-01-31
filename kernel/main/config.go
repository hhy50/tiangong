package main

type config struct {
	// proxy
	ProxyHost string `prop:"proxyHost" default:"127.0.0.1"`
	ProxyPort int    `prop:"proxyPort" default:"8090"`

	// server
	ServerHost string `prop:"serverHost" default:"127.0.0.1"`
	ServerPort int    `prop:"serverPort" default:"2024"`
	Token      string `prop:"token"`
	Subhost    string `prop:"subhost"`
}
