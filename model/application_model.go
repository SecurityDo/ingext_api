package model

import (
	"encoding/json"
	"time"
)

type ApplicationTemplateConfig struct {
	Name           string            `json:"name" yaml:"name"`
	DisplayName    string            `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Description    string            `json:"description,omitempty" yaml:"description,omitempty"`
	Icon           string            `json:"icon,omitempty" yaml:"icon,omitempty"`
	Category       string            `json:"category,omitempty" yaml:"category,omitempty"`
	ResourceGroups []string          `json:"resourceGroups,omitempty" yaml:"resourceGroups,omitempty"`
	Parameters     []*InputParameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Output         []*InputParameter `json:"output,omitempty" yaml:"output,omitempty"`
}

type ApplicationConfigResource struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       ApplicationConfigSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	//Raw        string                `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type ApplicationConfigSpec struct {
	DisplayName    string            `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Description    string            `json:"description,omitempty" yaml:"description,omitempty"`
	AdminConsent   bool              `json:"adminConsent,omitempty" yaml:"adminConsent,omitempty"`
	Category       string            `json:"category,omitempty" yaml:"category,omitempty"`
	Icon           string            `json:"icon,omitempty" yaml:"icon,omitempty"`
	ResourceGroups []string          `json:"resourceGroups,omitempty" yaml:"resourceGroups,omitempty"`
	Parameters     []*InputParameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Output         []*InputParameter `json:"output,omitempty" yaml:"output,omitempty"`
}

type ApplicationTemplate struct {
	ApplicationConfig *ApplicationConfigResource `json:"application,omitempty" yaml:"application,omitempty"`
	//UserInput         *UserInputResource         `json:"userInput,omitempty" yaml:"userInput,omitempty"`
	Routers      map[string]*RouterResource      `json:"routers,omitempty" yaml:"routers,omitempty"`
	Processors   map[string]*ProcessorResource   `json:"processors,omitempty" yaml:"processors,omitempty"`
	DataSources  map[string]*DataSourceResource  `json:"dataSources,omitempty" yaml:"dataSources,omitempty"`
	DataSinks    map[string]*DataSinkResource    `json:"dataSinks,omitempty" yaml:"dataSinks,omitempty"`
	Pipes        map[string]*PipeResource        `json:"pipes,omitempty" yaml:"pipes,omitempty"`
	Integrations map[string]*IntegrationResource `json:"integrations,omitempty" yaml:"integrations,omitempty"`
	AWSRoles     map[string]*AWSRoleResource     `json:"awsRoles,omitempty" yaml:"awsRoles,omitempty"`
	Collectors   map[string]*CollectorResource   `json:"collectors,omitempty" yaml:"collectors,omitempty"`
}

type AppTemplateEntry struct {
	Config  *ApplicationTemplateConfig `json:"config"`
	Content string                     `json:"content"`
}

type CollectorSpec struct {
	Name        string            `json:"name,omitempty" yaml:"name,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Output      []*InputParameter `json:"output,omitempty" yaml:"output,omitempty"`
}

type CollectorResource struct {
	TypeMeta   `yaml:",inline" json:",inline"`
	ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       *CollectorSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	//Raw        string         `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type AWSRoleResource struct {
	TypeMeta   `yaml:",inline" json:",inline"`
	ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       *AWSRoleSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	//Raw        string       `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type AWSRoleSpec struct {
	IamRoleARN string `json:"iamRoleARN,omitempty" yaml:"iamRoleARN,omitempty"`
	ExternalID string `json:"externalID,omitempty" yaml:"externalID,omitempty"`
}

type IntegrationSpec struct {
	Type         string          `json:"type,omitempty" yaml:"type,omitempty"`
	Description  string          `json:"description,omitempty" yaml:"description,omitempty"`
	AdminConsent string          `json:"adminConsent,omitempty" yaml:"adminConsent,omitempty"`
	Config       json.RawMessage `json:"config,omitempty" yaml:"config,omitempty"`
	Secret       json.RawMessage `json:"secret,omitempty" yaml:"secret,omitempty"`

	//Config       *IntegrationConfig `json:"config,omitempty" yaml:"config,omitempty"`
	//Secret       *IntegrationSecret `json:"secret,omitempty" yaml:"secret,omitempty"`
	Output []*InputParameter `json:"output,omitempty" yaml:"output,omitempty"`
}

type IntegrationResource struct {
	TypeMeta   `yaml:",inline" json:",inline"`
	ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       *IntegrationSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	//Raw        string           `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type GitObject struct {
	Repo   string `json:"repo,omitempty" yaml:"repo,omitempty"`
	Path   string `json:"path,omitempty" yaml:"path,omitempty"`
	Branch string `json:"branch,omitempty" yaml:"branch,omitempty"`
}

type ProcessorSpec struct {
	Type   string     `json:"type,omitempty" yaml:"type,omitempty"`
	Local  *string    `json:"local,omitempty" yaml:"local,omitempty"`
	Remote *GitObject `json:"remote,omitempty" yaml:"remote,omitempty"`
}

type ProcessorResource struct {
	TypeMeta   `yaml:",inline" json:",inline"`
	ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       *ProcessorSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	//Raw        string         `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type DataSinkResource struct {
	TypeMeta   `yaml:",inline" json:",inline"`
	ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       *DataSinkSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	//Raw        string        `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type DataSinkSpec struct {
	Type   string          `json:"type,omitempty" yaml:"type,omitempty"`
	Config *DataSinkConfig `json:"config,omitempty" yaml:"config,omitempty"`
	Packer string          `json:"packer,omitempty" yaml:"packer,omitempty"`
}

type DataSourceResource struct {
	TypeMeta   `yaml:",inline" json:",inline"`
	ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       *DataSourceSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	//Raw        string          `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type DataSourceSpec struct {
	Type     string              `json:"type,omitempty" yaml:"type,omitempty"`
	Format   string              `json:"format,omitempty" yaml:"format,omitempty"`
	Config   *DataSourceConfig   `json:"config,omitempty" yaml:"config,omitempty"`
	Router   string              `json:"router,omitempty" yaml:"router,omitempty"`
	Receiver *DataSourceReceiver `json:"receiver,omitempty" yaml:"receiver,omitempty"`
	Output   []*InputParameter   `json:"output,omitempty" yaml:"output,omitempty"`
}

type DataSourceReceiver struct {
	Processor string `json:"processor,omitempty" yaml:"processor,omitempty"`
	Input     string `json:"input,omitempty" yaml:"input,omitempty"`
}

type PipeSpec struct {
	Name           string         `json:"name,omitempty" yaml:"name,omitempty"`
	Processors     []string       `json:"processors,omitempty" yaml:"processors,omitempty"`
	Sinks          []string       `json:"sinks,omitempty" yaml:"sinks,omitempty"`
	RouterSelector *LabelSelector `json:"routerSelector,omitempty" yaml:"routerSelector,omitempty"`
	Router         string         `json:"router,omitempty" yaml:"router,omitempty"`
	Priority       int            `json:"priority,omitempty" yaml:"priority,omitempty"`
	//RouterID       string         `json:"-" yaml:"-"`
}

type RouterSpec struct {
	ThreadCount int `json:"threadCount,omitempty" yaml:"threadCount,omitempty"`
	//Pipes       []*PipeSpec `json:"pipes,omitempty" yaml:"pipes,omitempty"`
}

type PipeResource struct {
	TypeMeta   `yaml:",inline" json:",inline"`
	ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       *PipeSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	Raw        string    `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type RouterResource struct {
	TypeMeta   `yaml:",inline" json:",inline"`
	ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       *RouterSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	Raw        string      `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type ValueReference struct {
	Kind     string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Resource string `json:"resource,omitempty" yaml:"resource,omitempty"`
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
}

type EnumValue struct {
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
	Label string `json:"label,omitempty" yaml:"label,omitempty"`
}
type InputParameter struct {
	Name         string       `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string       `json:"description,omitempty" yaml:"description,omitempty"`
	DefaultValue string       `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"`
	DataType     string       `json:"dataType,omitempty" yaml:"dataType,omitempty"`
	Sensitive    bool         `json:"sensitive,omitempty" yaml:"sensitive,omitempty"`
	Optional     bool         `json:"optional,omitempty" yaml:"optional,omitempty"`
	Value        string       `json:"value,omitempty" yaml:"value,omitempty"`
	IsList       bool         `json:"isList,omitempty" yaml:"isList,omitempty"`
	Enums        []*EnumValue `json:"enums,omitempty" yaml:"enums,omitempty"`

	ValueRef *ValueReference `json:"valueRef,omitempty" yaml:"valueRef,omitempty"`
}

type InstallAppInstanceRequest struct {
	Config *InstanceConfig `json:"config"`
}

type InstanceState struct {
	Application string               `json:"application,omitempty" yaml:"application,omitempty"`
	AppTemplate *ApplicationTemplate `json:"appTemplate" yaml:"appTemplate"`
	Config      *InstanceConfig      `json:"config,omitempty" yaml:"config,omitempty"`
	State       string               `json:"state,omitempty" yaml:"state,omitempty"`
	ErrorMsg    string               `json:"errorMsg,omitempty" yaml:"errorMsg,omitempty"`

	Actions []*AppInstanceAction `json:"-" yaml:"-"`
	Output  []*AppInstanceOutput `json:"-" yaml:"-"`
}

type ResourceMeta struct {
	ID       string `json:"id,omitempty" yaml:"id,omitempty"`
	Resource string `json:"resource,omitempty" yaml:"resource,omitempty"`
}

type AppInstanceAction struct {
	Action      string          `json:"action,omitempty" yaml:"action,omitempty"`
	Application string          `json:"application,omitempty" yaml:"application,omitempty"`
	Instance    string          `json:"instance,omitempty" yaml:"instance,omitempty"`
	Status      string          `json:"status,omitempty" yaml:"status,omitempty"`
	ErrorMsg    string          `json:"errorMsg,omitempty" yaml:"errorMsg,omitempty"`
	Resources   []*ResourceMeta `json:"resources,omitempty" yaml:"resources,omitempty"`
	Args        json.RawMessage `json:"args,omitempty" yaml:"args,omitempty"`
}

type AppInstanceOutput struct {
	Application string          `json:"application,omitempty" yaml:"application,omitempty"`
	Instance    string          `json:"instance,omitempty" yaml:"instance,omitempty"`
	Resource    string          `json:"resource,omitempty" yaml:"resource,omitempty"`
	Kind        string          `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name        string          `json:"name,omitempty" yaml:"name,omitempty"`
	Description string          `json:"description,omitempty" yaml:"description,omitempty"`
	Value       string          `json:"value,omitempty" yaml:"value,omitempty"`
	ValueRef    *ValueReference `json:"valueRef,omitempty" yaml:"valueRef,omitempty"`
}

// k8s TypeMeta and ObjectMeta are used to define the metadata for Kubernetes resources.
type TypeMeta struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	// Cannot be updated.
	// In CamelCase.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	// APIVersion defines the versioned schema of this representation of an object.
	// Servers should convert recognized schemas to the latest internal value, and
	// may reject unrecognized values.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	// +optional
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}
type ObjectMeta struct {

	// Name must be unique within a namespace. Is required when creating resources, although
	// some resources may allow a client to request the generation of an appropriate name
	// automatically. Name is primarily intended for creation idempotence and configuration
	// definition.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
	// +optional
	Name              string    `json:"name,omitempty" yaml:"name,omitempty"`
	UID               string    `json:"uid,omitempty" yaml:"uid,omitempty"`
	ResourceVersion   string    `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
	CreationTimestamp time.Time `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`

	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`

	// List of objects depended by this object. If ALL objects in the list have
	// been deleted, this object will be garbage collected. If this object is managed by a controller,
	// then an entry in this list will point to this controller, with the controller field set to true.
	// There cannot be more than one managing controller.
	// +optional
	// +patchMergeKey=uid
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=uid
	//OwnerReferences []OwnerReference `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
}

type LabelSelector struct {
	// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
	// map is equivalent to an element of matchExpressions, whose key field is "key", the
	// operator is "In", and the values array contains only "value". The requirements are ANDed.
	// +optional
	MatchLabels map[string]string `json:"matchLabels,omitempty" yaml:"matchLabels,omitempty"`
	// matchExpressions is a list of label selector requirements. The requirements are ANDed.
	// +optional
	// +listType=atomic
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty" yaml:"matchExpressions,omitempty"`
}

// A label selector requirement is a selector that contains values, a key, and an operator that
// relates the key and values.
type LabelSelectorRequirement struct {
	// key is the label key that the selector applies to.
	Key string `json:"key" yaml:"key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator LabelSelectorOperator `json:"operator" yaml:"operator"`
	// values is an array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. This array is replaced during a strategic
	// merge patch.
	// +optional
	// +listType=atomic
	Values []string `json:"values,omitempty" yaml:"values,omitempty"`
}

// A label selector operator is the set of operators that can be used in a selector requirement.
type LabelSelectorOperator string

type GenericResource struct {
	TypeMeta   `yaml:",inline" json:",inline"`
	ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       map[string]interface{} `yaml:"spec,omitempty" json:"spec,omitempty"`
	Raw        []byte                 `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type InstanceConfig struct {
	Application     string            `json:"application,omitempty" yaml:"application,omitempty"`
	Instance        string            `json:"instance,omitempty" yaml:"instance,omitempty"`
	DisplayName     string            `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	InputParameters []*InputParameter `json:"inputParameters,omitempty" yaml:"inputParameters,omitempty"`
}
