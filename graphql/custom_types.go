package graphql

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-json"
)

type GroupingContext string

const (
	GroupingContextMalwareAnalysis    GroupingContext = "malware-analysis"
	GroupingContextSuspiciousActivity GroupingContext = "suspicious-activity"
	GroupingContextUnspecified        GroupingContext = "unspecified"
)

type IndicatorType string

const (
	IndicatorTypeAnomalousActivity IndicatorType = "anomalous-activity"
	IndicatorTypeAnonymization     IndicatorType = "anonymization"
	IndicatorTypeAttribution       IndicatorType = "attribution"
	IndicatorTypeBenign            IndicatorType = "benign"
	IndicatorTypeCompromised       IndicatorType = "compromised"
	IndicatorTypeMaliciousActivity IndicatorType = "malicious-activity"
	IndicatorTypeUnknown           IndicatorType = "unknown"
)

type OpinionType string

const (
	OpinionTypeStronglyDisagree OpinionType = "strongly-disagree"
	OpinionTypeDisagree         OpinionType = "disagree"
	OpinionTypeNeutral          OpinionType = "neutral"
	OpinionTypeAgree            OpinionType = "agree"
	OpinionTypeStronglyAgree    OpinionType = "strongly-agree"
)

type PatternType string

const (
	PatternTypeEql          PatternType = "eql"
	PatternTypePcre         PatternType = "pcre"
	PatternTypeShodan       PatternType = "shodan"
	PatternTypeSigma        PatternType = "sigma"
	PatternTypeSnort        PatternType = "snort"
	PatternTypeSpl          PatternType = "spl"
	PatternTypeStix         PatternType = "stix"
	PatternTypeSuricata     PatternType = "suricata"
	PatternTypeTaniumSignal PatternType = "tanium-signal"
	PatternTypeYara         PatternType = "yara"
)

type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformLinux   Platform = "linux"
	PlatformMacos   Platform = "macos"
	PlatformWindows Platform = "windows"
)

type StixCyberObservableType string

const (
	StixCyberObservableTypeArtifact                 StixCyberObservableType = "Artifact"
	StixCyberObservableTypeAutonomousSystem         StixCyberObservableType = "Autonomous-System"
	StixCyberObservableTypeBankAccount              StixCyberObservableType = "Bank-Account"
	StixCyberObservableTypeCryptocurrencyWallet     StixCyberObservableType = "Cryptocurrency-Wallet"
	StixCyberObservableTypeCryptographicKey         StixCyberObservableType = "Cryptographic-Key"
	StixCyberObservableTypeCredential               StixCyberObservableType = "Credential"
	StixCyberObservableTypeDirectory                StixCyberObservableType = "Directory"
	StixCyberObservableTypeDomainName               StixCyberObservableType = "Domain-Name"
	StixCyberObservableTypeEmailAddr                StixCyberObservableType = "Email-Addr"
	StixCyberObservableTypeEmailMessage             StixCyberObservableType = "Email-Message"
	StixCyberObservableTypeEmailMimePartType        StixCyberObservableType = "Email-Mime-Part-Type"
	StixCyberObservableTypeFile                     StixCyberObservableType = "File"
	StixCyberObservableTypeHostname                 StixCyberObservableType = "Hostname"
	StixCyberObservableTypeIPV4Addr                 StixCyberObservableType = "IPv4-Addr"
	StixCyberObservableTypeIPV6Addr                 StixCyberObservableType = "IPv6-Addr"
	StixCyberObservableTypeMacAddr                  StixCyberObservableType = "Mac-Addr"
	StixCyberObservableTypeMediaContent             StixCyberObservableType = "Media-Content"
	StixCyberObservableTypeMutex                    StixCyberObservableType = "Mutex"
	StixCyberObservableTypeNetworkTraffic           StixCyberObservableType = "Network-Traffic"
	StixCyberObservableTypePaymentCard              StixCyberObservableType = "Payment-Card"
	StixCyberObservableTypePersona                  StixCyberObservableType = "Persona"
	StixCyberObservableTypePhoneNumber              StixCyberObservableType = "Phone-Number"
	StixCyberObservableTypeProcess                  StixCyberObservableType = "Process"
	StixCyberObservableTypeSimpleObservable         StixCyberObservableType = "Simple-Observable"
	StixCyberObservableTypeSoftware                 StixCyberObservableType = "Software"
	StixCyberObservableTypeStixFile                 StixCyberObservableType = "StixFile"
	StixCyberObservableTypeText                     StixCyberObservableType = "Text"
	StixCyberObservableTypeTrackingNumber           StixCyberObservableType = "TrackingNumber"
	StixCyberObservableTypeURL                      StixCyberObservableType = "Url"
	StixCyberObservableTypeUserAccount              StixCyberObservableType = "User-Account"
	StixCyberObservableTypeUserAgent                StixCyberObservableType = "User-Agent"
	StixCyberObservableTypeWindowsRegistryKey       StixCyberObservableType = "Windows-Registry-Key"
	StixCyberObservableTypeWindowsRegistryValueType StixCyberObservableType = "Windows-Registry-Value-Type"
	StixCyberObservableTypeX509Certificate          StixCyberObservableType = "X509-Certificate"
)

type ReportType string

const (
	ReportTypeInternal ReportType = "internal-report"
	ReportTypeThreat   ReportType = "threat-report"
)

type ConfidenceLevelInput struct {
	MaxConfidence int                            `json:"max_confidence"`
	Overrides     []ConfidenceLevelOverrideInput `json:"overrides,omitempty"`
}

func (c ConfidenceLevelInput) MarshalJSON() ([]byte, error) {
	if len(c.Overrides) == 0 {
		return []byte(fmt.Sprintf(`{"max_confidence":%d, "overrides":[]}`, c.MaxConfidence)), nil
	}

	type tempConfidenceLevelInput ConfidenceLevelInput

	//nolint:wrapcheck // To avoid complete implementation
	return json.Marshal((*tempConfidenceLevelInput)(&c))
}

type ConfidenceLevelOverrideInput struct {
	EntityType    string `json:"entity_type,omitempty"`
	MaxConfidence int    `json:"max_confidence,omitempty"`
}

func (c ConfidenceLevelOverrideInput) MarshalJSON() ([]byte, error) {
	if reflect.ValueOf(c).IsZero() {
		return []byte("null"), nil
	}

	type tempConfidenceLevelOverrideInput ConfidenceLevelOverrideInput

	//nolint:wrapcheck // To avoid complete implementation
	return json.Marshal((*tempConfidenceLevelOverrideInput)(&c))
}
