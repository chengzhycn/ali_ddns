package ali

type DNSRecord struct {
	RecordId string
	Domain   string
	RR       string
	Type     string
	Value    string
	TTL      int64
}

func (r *DNSRecord) String() string {
	return r.Domain + " " + r.RR + " " + r.Type + " " + r.Value + " " + r.RecordId
}

func newDNSRecord(domain, rr, type_, value string, ttl int64) *DNSRecord {
	return &DNSRecord{
		Domain: domain,
		RR:     rr,
		Type:   type_,
		Value:  value,
		TTL:    ttl,
	}
}

func NewDNSRecordWithDefaults(domain, rr, type_ string, value string) *DNSRecord {
	return newDNSRecord(domain, rr, type_, value, 600)
}

type Client interface {
	DescribeDNSRecord(domain string) ([]*DNSRecord, error)
	AddDNSRecord(record *DNSRecord) error
	UpdateDNSRecord(record *DNSRecord) error
	DeleteDNSRecord(record *DNSRecord) error
}
