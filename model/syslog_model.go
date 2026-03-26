package model

type SyslogConfig struct {
	Domain     string `json:"domain"`
	PortBegin  int    `json:"port_begin"`
	PortEnd    int    `json:"port_end"`
	SyslogTLS  bool   `json:"syslog_tls,omitempty"`
	TLSRfc6587 bool   `json:"tls_rfc6587,omitempty"`
	SyslogUDP  bool   `json:"syslog_udp,omitempty"`
	SyslogTCP  bool   `json:"syslog_tcp,omitempty"`

	SyslogTLSPort  int `json:"syslog_tls_port,omitempty"`
	TLSRfc6587Port int `json:"tls_rfc6587_port,omitempty"`
	SyslogUDPPort  int `json:"syslog_udp_port,omitempty"`
	SyslogTCPPort  int `json:"syslog_tcp_port,omitempty"`
}

type SyslogPortRequest struct {
	SyslogTLS  bool `json:"syslog_tls"`
	TLSRfc6587 bool `json:"tls_rfc6587"`
	SyslogUDP  bool `json:"syslog_udp"`
	SyslogTCP  bool `json:"syslog_tcp"`
}
type GetSyslogConfigResponse struct {
	Config *SyslogConfig `json:"config,omitempty"`
}
