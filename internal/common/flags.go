package common

type ASNFlags struct {
	IP       string
	ASN      string
	Include  []string
	Excludes []string
	Fields   []string
	Output   string
}

type AbuseFlags struct {
	IP       string
	Excludes []string
	Fields   []string
	Output   string
}

type TimezoneFlags struct {
	IP        string
	Tz        string
	Location  string
	Latitude  float64
	Longitude float64
	IataCode  string
	IcaoCode  string
	LoCode    string
	Language  string
	Output    string
}

type AstronomyFlags struct {
	IP        string
	Location  string
	Latitude  float64
	Longitude float64
	Language  string
	Tz        string
	Elevation float64
	Output    string
}

type AstronomyTimeSeriesFlags struct {
	IP        string
	Location  string
	Latitude  float64
	Longitude float64
	Language  string
	Output    string
	DateStart string
	DateEnd   string
}

type TimeConversionFlags struct {
	LatitudeFrom  float64
	LongitudeFrom float64
	LatitudeTo    float64
	LongitudeTo   float64
	Output        string
	TimezoneFrom  string
	TimezoneTo    string
	LocationFrom  string
	LocationTo    string
	IataCodeFrom  string
	IataCodeTo    string
	IcaoCodeFrom  string
	IcaoCodeTo    string
	LoCodeFrom    string
	LoCodeTo      string
	Time          string
}

type IPSecurityFlags struct {
	IP       string
	Include  []string
	Excludes []string
	Fields   []string
	Language string
	Output   string
}

type BulkIPSecurityFlags struct {
	IPs        []string
	Include    []string
	Excludes   []string
	Fields     []string
	Language   string
	Output     string
	OutputFile string
	File       string
}

type ParseUserAgentFlags struct {
	UserAgent string
	Output    string
}

type IpgeoFlags struct {
	IP       string
	Include  []string
	Excludes []string
	Fields   []string
	Language string
	Output   string
}

type BulkIpgeoFlags struct {
	IPs        []string
	Include    []string
	Excludes   []string
	Fields     []string
	Language   string
	Output     string
	File       string
	OutputFile string
}

type ParseBulkUserAgentFlags struct {
	UserAgents []string
	Output     string
}
