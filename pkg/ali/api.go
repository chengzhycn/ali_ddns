package ali

import (
	"fmt"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type AliDNSClient struct {
	*alidns.Client

	runtimeOption *util.RuntimeOptions
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId string, accessKeySecret string) (Client, error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(accessKeyId),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(accessKeySecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Alidns
	config.Endpoint = tea.String("dns.aliyuncs.com")

	_result, err := alidns.NewClient(config)
	client := &AliDNSClient{
		_result,
		&util.RuntimeOptions{},
	}
	return client, err
}

/**
* 使用STS鉴权方式初始化账号Client，推荐此方式。
* @param accessKeyId
* @param accessKeySecret
* @param securityToken
* @return Client
* @throws Exception
 */
func CreateClientWithSTS(accessKeyId string, accessKeySecret string, securityToken string) (Client, error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(accessKeyId),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(accessKeySecret),
		// 必填，您的 Security Token
		SecurityToken: tea.String(securityToken),
		// 必填，表明使用 STS 方式
		Type: tea.String("sts"),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Alidns
	config.Endpoint = tea.String("dns.aliyuncs.com")
	_result, err := alidns.NewClient(config)
	client := &AliDNSClient{
		_result,
		&util.RuntimeOptions{},
	}
	return client, err
}

func (c *AliDNSClient) AddDNSRecord(record *DNSRecord) error {
	request := &alidns.AddDomainRecordRequest{
		DomainName: tea.String(record.Domain),
		RR:         tea.String(record.RR),
		Type:       tea.String(record.Type),
		Value:      tea.String(record.Value),
		TTL:        tea.Int64(record.TTL),
	}

	response, err := c.AddDomainRecordWithOptions(request, c.runtimeOption)
	if err != nil {
		return err
	}
	if tea.Int32Value(response.StatusCode) != 200 {
		return fmt.Errorf("wrong status code: %d", tea.Int32Value(response.StatusCode))
	}

	record.RecordId = tea.StringValue(response.Body.RecordId)

	return nil
}

func (c *AliDNSClient) UpdateDNSRecord(record *DNSRecord) error {
	request := &alidns.UpdateDomainRecordRequest{
		RecordId: tea.String(record.RecordId),
		RR:       tea.String(record.RR),
		Type:     tea.String(record.Type),
		Value:    tea.String(record.Value),
		TTL:      tea.Int64(record.TTL),
	}

	response, err := c.UpdateDomainRecordWithOptions(request, c.runtimeOption)
	if err != nil {
		return err
	}
	if tea.Int32Value(response.StatusCode) != 200 {
		return fmt.Errorf("wrong status code: %d", tea.Int32Value(response.StatusCode))
	}

	record.RecordId = tea.StringValue(response.Body.RecordId)

	return nil
}

func (c *AliDNSClient) DeleteDNSRecord(record *DNSRecord) error {
	request := &alidns.DeleteDomainRecordRequest{
		RecordId: tea.String(record.RecordId),
	}

	response, err := c.DeleteDomainRecordWithOptions(request, c.runtimeOption)
	if err != nil {
		return err
	}
	if tea.Int32Value(response.StatusCode) != 200 {
		return fmt.Errorf("wrong status code: %d", tea.Int32Value(response.StatusCode))
	}

	return nil
}

func (c *AliDNSClient) DescribeDNSRecord(domain string) ([]*DNSRecord, error) {
	reuqest := &alidns.DescribeDomainRecordsRequest{
		DomainName: tea.String(domain),
	}

	response, err := c.DescribeDomainRecordsWithOptions(reuqest, c.runtimeOption)
	if err != nil {
		return nil, err
	}
	if tea.Int32Value(response.StatusCode) != 200 {
		return nil, fmt.Errorf("wrong status code: %d", response.StatusCode)
	}

	result := make([]*DNSRecord, 0)
	for _, record := range response.Body.DomainRecords.Record {
		r := &DNSRecord{
			RecordId: tea.StringValue(record.RecordId),
			Domain:   tea.StringValue(record.DomainName),
			RR:       tea.StringValue(record.RR),
			Type:     tea.StringValue(record.Type),
			Value:    tea.StringValue(record.Value),
			TTL:      tea.Int64Value(record.TTL),
		}

		result = append(result, r)
	}

	return result, nil
}
