package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type DataSourceConfig struct {
	Type   string `json:"type"`
	Name   string `json:"name,omitempty"`
	ID     string `json:"id,omitempty"`
	Format string `json:"format,omitempty"`
	Tags   []*Tag `json:"tags,omitempty"`

	Compression string `json:"compression,omitempty"`
	//Page string `json:"page,omitempty"`
	ReceiverName string `json:"receiverName,omitempty"`
	// raw, line, json
	ReceiverInput string `json:"receiverInput,omitempty"`

	//Redis  *RedisSourceConfig  `json:"redis,omitempty"`
	S3 *S3SourceConfig `json:"s3,omitempty"`

	S3Notification *S3NotificationSourceConfig `json:"s3Notification,omitempty"`
	Kinesis        *KinesisSourceConfig        `json:"kinesis,omitempty"`

	Prom    *PromSourceConfig    `json:"prom,omitempty"`
	Hec     *HecSourceConfig     `json:"hec,omitempty"`
	Webhook *WebhookSourceConfig `json:"webhook,omitempty"`
	Plugin  *PluginSourceConfig  `json:"plugin,omitempty"`
	Syslog  *SyslogSourceConfig  `json:"syslog,omitempty"`
	Secret  json.RawMessage      `json:"secret,omitempty"`

	IntegrationPull *IntegrationPullSourceConfig `json:"integrationPull,omitempty"`
}

type RedisConfig struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Queue string `json:"queue"`
}

type RedisSourceConfig struct {
	Redis *RedisConfig `json:"redis"`
}

type HecSourceConfig struct {
	Ack bool   `json:"ack"`
	URL string `json:"url"`
	// tenantConfig *TenantAWSConfig `json:"-"`
	// token        string           `json:"-"`

	secret *HecSecret `json:"-"`
}

type HecSecret struct {
	Token string `json:"token"`
}

type WebhookSourceConfig struct {
	URL             string `json:"url"`
	Token           string `json:"token"`
	SecurityToken   string `json:"securityToken"`
	SignatureHeader string `json:"signatureHeader"`
}

type PluginSourceConfig struct {
	Name        string `json:"name,omitempty"`
	ID          string `json:"id,omitempty"`
	Integration string `json:"integration,omitempty"`
}

type S3SourceConfig struct {
	Plugin *PluginSourceConfig `json:"plugin,omitempty"`
	Bucket string              `json:"bucket"`
	Prefix string              `json:"prefix,omitempty"`
	Begin  string              `json:"begin,omitempty"`
	End    string              `json:"end,omitempty"`
}

type SyslogSourceConfig struct {
	Path string `json:"path"`
}

type RedisSinkConfig struct {
	Redis *RedisConfig `json:"redis"`
}

type S3SinkConfig struct {
	IntegrationID string              `json:"integrationID"`
	Plugin        *PluginSourceConfig `json:"plugin,omitempty"`
	Bucket        string              `json:"bucket"`
	Compression   string              `json:"compression,omitempty"`
	ObjectPath    string              `json:"objectPath"`
}

type LavaDBSinkConfig struct {
	Tenant string `json:"tenant"`
	Index  string `json:"index"`
}

type DataSinkConfig struct {
	Type string `json:"type"`
	Name string `json:"name"`
	ID   string `json:"id"`

	GroupSupport bool   `json:"groupSupport,omitempty"`
	GroupLambda  string `json:"groupLambda,omitempty"`

	Redis *RedisSinkConfig `json:"redis,omitempty"`
	// LavaDB *LavaDBSinkConfig `json:"lavaDB,omitempty"`
	S3      *S3SinkConfig      `json:"s3,omitempty"`
	Hec     *HecSinkConfig     `json:"hec,omitempty"`
	Webhook *WebhookSinkConfig `json:"webhook,omitempty"`

	Kinesis  *KinesisSinkConfig   `json:"kinesis,omitempty"`
	Firehose *FirehoseSinkConfig  `json:"firehose,omitempty"`
	Lambda   *AWSLambdaSinkConfig `json:"lambda,omitempty"`

	DataLake *DataLakeSinkConfig `json:"dataLake,omitempty"`

	PROM *PromSinkConfig `json:"prom,omitempty"`
	Loki *LokiSinkConfig `json:"loki,omitempty"`

	Secret json.RawMessage `json:"secret,omitempty"`

	Tags []*Tag `json:"tags,omitempty"`

	PackerConfig *SinkPackerConfig `json:"packerConfig,omitempty"`
	// maximum event count before flush, default 1024
	FlushCount int64 `json:"flushCount,omitempty"`
	// maximum buffer size before flush, default 1MB
	FlushBuffer int64 `json:"flushBuffer,omitempty"`
	// if buffer is not empty, flush interval in seconds, default 60
	FlushInterval int64 `json:"flushInterval,omitempty"`
}

type SinkPackerConfig struct {
	Name string `json:"name,omitempty"`
	FPL  string `json:"-"`
}

type FirehoseSinkConfig struct {
	IntegrationID string              `json:"integrationID"`
	Plugin        *PluginSourceConfig `json:"plugin,omitempty"`

	Config *AWSFirehoseConfig `json:"config,omitempty"`
	//Secret  *integrationModel.AWSFirehoseSecret `json:"-"`
	//RoleObj *integrationModel.InstanceRole      `json:"-"`
	AuthInfo *AWSAuthInfo `json:"-"`
}

type AWSLambdaSinkConfig struct {
	IntegrationID string              `json:"integrationID"`
	Plugin        *PluginSourceConfig `json:"plugin,omitempty"`

	Config *AWSLambdaConfig `json:"config,omitempty"`
	//Secret  *integrationModel.AWSLambdaSecret `json:"-"`
	//RoleObj *integrationModel.InstanceRole    `json:"-"`
	AuthInfo *AWSAuthInfo `json:"-"`
}

type KinesisSinkConfig struct {
	IntegrationID string              `json:"integrationID"`
	Plugin        *PluginSourceConfig `json:"plugin,omitempty"`

	Config *AWSKinesisStreamConfig `json:"config,omitempty"`
	//Secret  *integrationModel.AWSKinesisStreamSecret `json:"-"`
	//RoleObj *integrationModel.InstanceRole           `json:"-"`
	AuthInfo *AWSAuthInfo `json:"-"`
}

//type HecSecret struct {
//	Token string `json:"token"`
//}

type HecSinkConfig struct {
	// Ack          bool             `json:"ack"`
	Token              string `json:"token"`
	URL                string `json:"url"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
	//tenantConfig       *TenantAWSConfig `json:"-"`
	//token        string           `json:"-"`
	secret *HecSecret `json:"-"`
}

//	func (r *HecSinkConfig) GetTenantConfig() *TenantAWSConfig {
//		return r.tenantConfig
//	}
func (r *HecSinkConfig) GetToken() string {
	if r.secret == nil {
		fmt.Printf("Hec Token is nil\n")
		return ""
	}
	return r.secret.Token
}

type RouterConfig struct {
	Name        string   `json:"name"`
	ID          string   `json:"id"`
	WorkerCount int      `json:"workerCount"`
	PipeIDs     []string `json:"pipeIDs"`
}

type StreamPipeConfig struct {
	Name           string   `json:"name"`
	ID             string   `json:"id"`
	RouterID       string   `json:"routerID"`
	MatchAll       bool     `json:"matchAll"`
	Selector       string   `json:"selector,omitempty"`
	ProcessorNames []string `json:"processorNames"`
	ChannelID      string   `json:"channelID,omitempty"`
	SinkIDs        []string `json:"sinkIDs,omitempty"`
}

type DataLakeSinkConfig struct {
	Datalake      string `json:"datalake,omitempty"`      // the data lake name, e.g. "managed"
	DatalakeIndex string `json:"datalakeIndex,omitempty"` // the data lake index, e.g. "cloudtrail", "aws-vpc"
	SchemaName    string `json:"schemaName" yaml:"schemaName"`
}

type PromSinkConfig struct {
	IntegrationID string              `json:"integrationID"`
	Plugin        *PluginSourceConfig `json:"plugin,omitempty"`
}

type LokiSinkConfig struct {
	IntegrationID string              `json:"integrationID"`
	Plugin        *PluginSourceConfig `json:"plugin,omitempty"`
}

type ChannelConfig struct {
	Name    string   `json:"name"`
	ID      string   `json:"id"`
	SinkIDs []string `json:"sinkIDs"`
}

type RouterInput struct {
	RouterID string `json:"routerID"`
	SourceID string `json:"sourceID"`
}

type PluginNotification struct {
	Severity  string `json:"severity"`
	Subject   string `json:"subject"`
	ID        string `json:"id,omitempty"`
	Source    string `json:"source,omitempty"`
	CreatedOn string `json:"createdOn,omitempty"`
	Message   string `json:"message"`
}

type ComponentErrorState struct {
	ID     string                `json:"id"`
	Name   string                `json:"name"`
	Errors []*PluginNotification `json:"errors"`
	Alerts []*PluginNotification `json:"alerts"`
}

type LogEvent struct {
	Source    string    `json:"source,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Msg       string    `json:"msg,omitempty"`
}

type ComponentInfo struct {
	ID     string                `json:"id"`
	Name   string                `json:"name"`
	Errors []*PluginNotification `json:"errors"`
	Alerts []*PluginNotification `json:"alerts"`
	State  json.RawMessage       `json:"state"`
	Logs   []*LogEvent           `json:"logs"`
}

type PromSourceConfig struct {
	URL string `json:"url"`
	// tenantConfig *TenantAWSConfig `json:"-"`
	// token        string           `json:"-"`

	// Basic, Bearer, None
	Authorization string `json:"authorization"` //Basic or Bearer or None

	// only valid when Authorization is Basic
	Username string `json:"username"`

	secret *PromSecret `json:"-"`
}

func (r *PromSourceConfig) GetSecret() *PromSecret {
	return r.secret
}
func (r *PromSourceConfig) SetSecret(raw json.RawMessage) {
	var secret PromSecret
	err := json.Unmarshal(raw, &secret)
	if err != nil {
		fmt.Printf("failed to unmarshal token: %s\n", err.Error())
	} else {
		r.secret = &secret
	}
}

type PromSecret struct {
	Token    string `json:"token,omitempty"`
	Password string `json:"password,omitempty"`
}

type S3NotificationSourceConfig struct {
	Plugin *PluginSourceConfig `json:"plugin,omitempty"`

	Config *S3NotificationConfig `json:"config,omitempty"`
	//Secret  *integrationModel.S3NotificationSecret `json:"-"`
	//RoleObj *integrationModel.InstanceRole         `json:"-"`
	AuthInfo *AWSAuthInfo `json:"-"`
}

type KinesisSourceConfig struct {
	IntegrationID string              `json:"integrationID"`
	Plugin        *PluginSourceConfig `json:"plugin,omitempty"`

	Config *AWSKinesisStreamConfig `json:"config,omitempty"`
	//Secret  *integrationModel.AWSKinesisStreamSecret `json:"-"`
	//RoleObj *integrationModel.InstanceRole           `json:"-"`

	AuthInfo *AWSAuthInfo `json:"-"`
}

func (r *PluginSourceConfig) IsResourceWatch() bool {
	return r.Integration == "Office365ResourceWatch" || r.Integration == "AzureAudit" || r.Integration == "LDAP" || r.Integration == "SentinelOne" || r.Integration == "Tenable" || r.Integration == "Qualys"
}

type IntegrationPullSourceConfig struct {
	Plugin    *PluginSourceConfig `json:"plugin,omitempty"`
	Processor string              `json:"processor"`
	FPL       string              `json:"-"`
	// default 300
	PollInterval int64 `json:"pollInterval"`
}

//func (r *PluginSourceConfig) GetTenantConfig() *TenantAWSConfig {
//	return r.tenantConfig
//}
//func (r *PluginSourceConfig) SetTenantConfig(c *TenantAWSConfig) {
//	r.tenantConfig = c
//}

type DelaySinkConfig struct {
	Delay int64 `json:"delay"`
}

type HttpHeader struct {
	Header string `json:"header"`
	Value  string `json:"value"`
}
type WebhookSinkConfig struct {
	URL                string        `json:"url"`
	Headers            []*HttpHeader `json:"headers"`
	InsecureSkipVerify bool          `json:"insecureSkipVerify"`
}

type ElasticSinkConfig struct {
	URL   string `json:"url"`
	Index string `json:"index"`
	// optional
	DataType string `json:"dataType"`
}

type ChangeEvent struct {
	Timestamp time.Time   `json:"timestamp"`
	EventType string      `json:"eventType"`
	Entry     interface{} `json:"entry"`
	Revision  int64       `json:"revision"`
}

type FPLScript struct {
	ID         int64  `json:"id"`
	Repository string `json:"repository,omitempty"`
	Group      string `json:"group"`
	GitPath    string `json:"gitPath,omitempty"`

	Name string `json:"name"`
	Type string `json:"type"`
	Tags []*Tag `json:"tags,omitempty"`

	//Channel      string           `json:"channel,omitempty"`
	//ActionConfig *FplActionConfig `json:"actionConfig,omitempty"`
	//RuleConfig   *FplRuleConfig   `json:"ruleConfig,omitempty"`

	//Arguments    []*fplmodel.ProgramArgument      `json:"arguments,omitempty"`
	//ReportConfig *fplreportModel.TaskReportConfig `json:"reportConfig,omitempty"`

	Description string `json:"description"`
	//ScriptLang  string    `json:"scriptLang"`
	ScriptText string    `json:"scriptText"`
	CreatedOn  time.Time `json:"createdOn"`
	UpdatedOn  time.Time `json:"updatedOn"`

	//Tenant string `json:"tenant,omitempty"`
	//Source string `json:"source,omitempty"`
}
