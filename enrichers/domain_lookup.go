package enrichers

import (
	"net"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/sheacloud/goflow-addons/utils"
)

type DomainLookup struct {
	SrcDomainName *string `json:",omitempty"`
	DstDomainName *string `json:",omitempty"`
}

type DomainLookupEnricher struct {
	dnsCache *cache.Cache
}

func (e *DomainLookupEnricher) Initialize() {
	e.dnsCache = cache.New(5*time.Minute, 10*time.Minute)
}

func reverseLookup(ip string) *string {
	names, err := net.LookupAddr(ip)
	if len(names) >= 1 && err == nil {
		return &names[0]
	}

	return nil
}

func (e *DomainLookupEnricher) Enrich(msgs []*utils.ExtendedFlowMessage) {
	for _, msg := range msgs {
		domainInfo := DomainLookup{}
		srcIP := net.IP(msg.SrcAddr).String()
		dstIP := net.IP(msg.DstAddr).String()

		var srcName, dstName *string

		srcAttempt, found := e.dnsCache.Get(srcIP)
		if found {
			srcName = srcAttempt.(*string)
		} else {
			srcName = reverseLookup(srcIP)
			e.dnsCache.Set(srcIP, srcName, cache.DefaultExpiration)
		}

		dstAttempt, found := e.dnsCache.Get(dstIP)
		if found {
			dstName = dstAttempt.(*string)
		} else {
			dstName = reverseLookup(dstIP)
			e.dnsCache.Set(dstIP, dstName, cache.DefaultExpiration)
		}

		domainInfo.SrcDomainName = srcName
		domainInfo.DstDomainName = dstName

		msg.Metadata["DomainLookup"] = domainInfo
	}
}
