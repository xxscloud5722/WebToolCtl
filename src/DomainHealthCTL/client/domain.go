package client

import (
	"errors"
	"github.com/samber/lo"
	"strings"
)

// Domain 域名结构
type Domain struct {
	Name    string
	Child   *[]*DomainChild
	Whois   *DomainWhois
	Message *string
}

// DomainChild 子域名结构
type DomainChild struct {
	Name    string
	SSL     *X509Certificate
	Message *string
}

// ParseDomain 解析域名
func ParseDomain(domain string) (*Domain, error) {
	var domains = ParseDomains([]string{domain})
	if len(domains) <= 0 {
		return nil, errors.New("domain error")
	}
	return domains[0], nil
}

// ParseDomains 解析域名
func ParseDomains(rows []string) []*Domain {
	var domainContext []string
	for _, row := range rows {
		row = strings.ToLower(strings.TrimSpace(row))
		if strings.HasPrefix(row, "#") {
			continue
		}
		domainContext = append(domainContext, row)
	}

	var domains []*Domain
	for _, context := range domainContext {
		var suffix = MatchSuffix(context)
		if suffix == "" {
			continue
		}
		var domainPrefix = strings.TrimSuffix(context, suffix)
		var prefix = strings.Split(domainPrefix, ".")
		if len(prefix) <= 0 {
			continue
		}
		var baseDomain = prefix[len(prefix)-1] + suffix
		domain, ok := lo.Find(domains, func(item *Domain) bool {
			return item.Name == baseDomain
		})
		if ok {
			*domain.Child = append(*domain.Child, &DomainChild{Name: context})
		} else {
			domains = append(domains, &Domain{
				Name:  baseDomain,
				Child: &[]*DomainChild{{Name: context}},
			})
		}
	}

	return domains
}

func Analysis(domains []*Domain) []*Domain {
	for _, domain := range domains {
		// Whois 解析
		whoisRows, err := Whois(*domain)
		if err != nil {
			var message = err.Error()
			domain.Message = &message
		} else {
			domain.Whois, err = ParseWhoisInfo(whoisRows)
			if err != nil {
				var message = err.Error()
				domain.Message = &message
			}
		}

		for _, child := range *domain.Child {
			// SSL
			ssl, err := SSL(child.Name)
			if err != nil {
				var message = err.Error()
				child.Message = &message
			}
			child.SSL = ssl
		}
	}
	return domains
}
