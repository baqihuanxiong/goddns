package dns

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"goddns/config"
)

var (
	ErrRecordNotFound = errors.New("domain record not found")
)

type AliDNSService struct {
	conf *config.AliDNS
	cli  *alidns.Client
}

// SDK reference: https://help.aliyun.com/document_detail/29821.html?spm=a2c4g.11186623.6.612.75f31096EdmS94
func NewAliDNSService(config *config.AliDNS) (*AliDNSService, error) {
	svc := &AliDNSService{conf: config}

	var err error
	if svc.cli, err = alidns.NewClientWithAccessKey(
		svc.conf.Region, svc.conf.AccessKey, svc.conf.AccessKeySecret); err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *AliDNSService) AddDomain(domain *Domain) error {
	req := alidns.CreateAddDomainRequest()
	req.Scheme = "https"
	req.DomainName = domain.DomainName

	_, err := s.cli.AddDomain(req)
	return err
}

func (s *AliDNSService) HasDomain(domain *Domain) (bool, error) {
	req := alidns.CreateDescribeDomainInfoRequest()
	req.Scheme = "https"
	req.DomainName = domain.DomainName

	_, err := s.cli.DescribeDomainInfo(req)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *AliDNSService) AddRecord(record *Record) error {
	req := alidns.CreateAddDomainRecordRequest()
	req.Scheme = "https"
	req.DomainName = record.DomainName
	req.RR = record.RR
	req.Type = record.Type
	req.Value = record.Value

	_, err := s.cli.AddDomainRecord(req)
	return err
}

func (s *AliDNSService) UpdateRecord(record *Record) error {
	req := alidns.CreateUpdateDomainRecordRequest()
	req.Scheme = "https"
	if info, err := s.getRecordInfo(record); err != nil {
		return err
	} else {
		req.RecordId = info.RecordId
		req.RR = info.RR
	}
	req.Type = record.Type
	req.Value = record.Value

	_, err := s.cli.UpdateDomainRecord(req)
	return err
}

func (s *AliDNSService) getRecordInfo(record *Record) (*alidns.Record, error) {
	req := alidns.CreateDescribeDomainRecordsRequest()
	req.Scheme = "https"
	req.DomainName = record.DomainName
	req.PageSize = "500"

	res, err := s.cli.DescribeDomainRecords(req)
	if err != nil {
		return nil, err
	}

	for _, r := range res.DomainRecords.Record {
		if r.DomainName == record.DomainName && r.RR == record.RR {
			return &r, nil
		}
	}
	return nil, ErrRecordNotFound
}

func (s *AliDNSService) HasRecord(record *Record) (bool, error) {
	if _, err := s.getRecordInfo(record); err != nil {
		return false, nil
	}
	return true, nil
}
