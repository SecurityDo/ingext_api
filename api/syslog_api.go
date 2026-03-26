package api

import (
	"github.com/SecurityDo/ingext_api/client"
	"github.com/SecurityDo/ingext_api/model"
)

// PlatformService provides helpers for calling platform_* endpoints.
type SyslogService struct {
	client *client.IngextClient
}

func NewSyslogService(client *client.IngextClient) *SyslogService {
	return &SyslogService{client: client}
}

func (s *SyslogService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

func (s *SyslogService) Register(ports []string) (resp *model.GetSyslogConfigResponse, err error) {
	request := &model.SyslogPortRequest{}

	for _, port := range ports {
		if port == "udp" {
			request.SyslogUDP = true
		} else if port == "tcp" {
			request.SyslogTCP = true
		} else if port == "tls" {
			request.SyslogTLS = true
		} else if port == "tls-rfc6587" {
			request.TLSRfc6587 = true
		}
	}
	//var resp model.LakeSearchResponse
	if err := s.call("ingext_syslog_register_config", request, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SyslogService) Update(ports []string) (resp *model.GetSyslogConfigResponse, err error) {
	request := &model.SyslogPortRequest{}

	for _, port := range ports {
		if port == "udp" {
			request.SyslogUDP = true
		} else if port == "tcp" {
			request.SyslogTCP = true
		} else if port == "tls" {
			request.SyslogTLS = true
		} else if port == "tls-rfc6587" {
			request.TLSRfc6587 = true
		}
	}
	//var resp model.LakeSearchResponse
	if err := s.call("ingext_syslog_update_config", request, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SyslogService) Get() (resp *model.GetSyslogConfigResponse, err error) {
	//var resp model.LakeSearchResponse
	if err := s.call("ingext_syslog_get_config", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SyslogService) Delete() (err error) {
	//var resp model.LakeSearchResponse
	if err := s.call("ingext_syslog_delete_config", nil, nil); err != nil {
		return err
	}
	return nil
}
