package cli

var CLI *Client

func Setup(domain string, port int, https *bool, externalToken string) {
	CLI = NewClientWithNormalTimeout(&Client{
		Rest:          nil,
		Ip:            domain,
		Port:          port,
		HTTPS:         https,
		ExternalToken: externalToken,
	})
}
