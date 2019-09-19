package model

// OnmsAlarm OpenNMS alarm entity
type OnmsAlarm struct {
	// Inherit from Events
	ID                   int              `json:"id" yaml:"id"`
	UEI                  string           `json:"uei" yaml:"uei"`
	EventTime            *Time            `json:"time,omitempty" yaml:"time,omitempty"`
	EventHost            string           `json:"host,omitempty" yaml:"host,omitempty"`
	EventSource          string           `json:"source,omitempty" yaml:"source,omitempty"`
	CreateTime           *Time            `json:"createTime,omitempty" yaml:"createTime,omitempty"`
	SnmpHost             string           `json:"snmpHost,omitempty" yaml:"snmpHost,omitempty"`
	Snmp                 string           `json:"snmp,omitempty" yaml:"snmp,omitempty"`
	NodeID               int              `json:"nodeId,omitempty" yaml:"nodeId,omitempty"`
	NodeLabel            string           `json:"nodeLabel,omitempty" yaml:"nodeLabel,omitempty"`
	IPAddress            string           `json:"ipAddress,omitempty" yaml:"ipAddress,omitempty"`
	ServiceType          OnmsServiceType  `json:"serviceType,omitempty" yaml:"serviceType,omitempty"`
	IfIndex              int              `json:"ifIndex,omitempty" yaml:"ifIndex,omitempty"`
	Severity             string           `json:"severity,omitempty" yaml:"severity,omitempty"`
	Log                  string           `json:"log,omitempty" yaml:"log,omitempty"`
	LogGroup             string           `json:"logGroup,omitempty" yaml:"logGroup,omitempty"`
	LogMessage           string           `json:"logMessage,omitempty" yaml:"logMessage,omitempty"`
	Display              string           `json:"display,omitempty" yaml:"display,omitempty"`
	Description          string           `json:"description,omitempty" yaml:"description,omitempty"`
	PathOutage           string           `json:"pathOutage,omitempty" yaml:"pathOutage,omitempty"`
	Correlation          string           `json:"correlation,omitempty" yaml:"correlation,omitempty"`
	SuppressedCount      int              `json:"suppressedCount,omitempty" yaml:"suppressedCount,omitempty"`
	OperatorInstructions string           `json:"operatorInstructions,omitempty" yaml:"operatorInstructions,omitempty"`
	OperatorAction       string           `json:"operatorAction,omitempty" yaml:"operatorAction,omitempty"`
	AutoAction           string           `json:"autoAction,omitempty" yaml:"autoAction,omitempty"`
	Parameters           []OnmsEventParam `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// Alarm fields
	ReductionKey          string     `json:"reductionKey,omitempty" yaml:"reductionKey,omitempty"`
	ClearKey              string     `json:"clearKey,omitempty" yaml:"clearKey,omitempty"`
	Type                  int        `json:"type,omitempty" yaml:"type,omitempty"`
	Count                 int        `json:"count,omitempty" yaml:"count,omitempty"`
	TroubleTicketID       string     `json:"troubleTicket,omitempty" yaml:"troubleTicket,omitempty"`
	TroubleTicketState    string     `json:"troubleTicketState,omitempty" yaml:"troubleTicketState,omitempty"`
	SuppressedUntil       *Time      `json:"suppressedUntil,omitempty" yaml:"suppressedUntil,omitempty"`
	SuppressedBy          string     `json:"suppressedBy,omitempty" yaml:"suppressedBy,omitempty"`
	SuppressedTime        *Time      `json:"suppressedTime,omitempty" yaml:"suppressedTime,omitempty"`
	AckID                 int        `json:"ackId,omitempty" yaml:"ackId,omitempty"`
	AckUser               string     `json:"ackUser,omitempty" yaml:"ackUser,omitempty"`
	AckTime               *Time      `json:"ackTime,omitempty" yaml:"ackTime,omitempty"`
	ApplicationDN         string     `json:"applicationDN,omitempty" yaml:"applicationDN,omitempty"`
	ManagedObjectInstance string     `json:"managedObjectInstance,omitempty" yaml:"managedObjectInstance,omitempty"`
	ManagedObjectType     string     `json:"managedObjectType,omitempty" yaml:"managedObjectType,omitempty"`
	OssPrimaryKey         string     `json:"ossPrimaryKey,omitempty" yaml:"ossPrimaryKey,omitempty"`
	X733AlarmType         string     `json:"x733AlarmType,omitempty" yaml:"x733AlarmType,omitempty"`
	X733ProbableCause     int        `json:"x733ProbableCause,omitempty" yaml:"x733ProbableCause,omitempty"`
	QosAlarmState         string     `json:"qosAlarmState,omitempty" yaml:"qosAlarmState,omitempty"`
	FirstAutomationTime   *Time      `json:"firstAutomationTime,omitempty" yaml:"firstAutomationTime,omitempty"`
	LastAutomationTime    *Time      `json:"lastAutomationTime,omitempty" yaml:"lastAutomationTime,omitempty"`
	FirstEventTime        *Time      `json:"firstEventTime,omitempty" yaml:"firstEventTime,omitempty"`
	LastEventTime         *Time      `json:"lastEventTime,omitempty" yaml:"lastEventTime,omitempty"`
	LastEvent             *OnmsEvent `json:"lastEvent,omitempty" yaml:"-"`
}

// OnmsAlarmList a list of alarms
type OnmsAlarmList struct {
	Count      int         `json:"count" yaml:"count"`
	TotalCount int         `json:"totalCount" yaml:"totalCount"`
	Offset     int         `json:"offset" yaml:"offset"`
	Alarms     []OnmsAlarm `json:"alarm" yaml:"alarms"`
}
