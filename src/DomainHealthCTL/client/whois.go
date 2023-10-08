package client

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/longyuan/lib.v3/times"
	"github.com/samber/lo"
	"math"
	"net"
	"strconv"
	"strings"
	"time"
)

// whois 服务器信息.
var whois = map[string]string{
	".br.com": "whois.centralnic.com",
	".cn.com": "whois.internic.net",
	".de.com": "whois.cnnic.cn",
	".eu.com": "whois.nic.top",
	".gb.com": "whois.nic.org",
	".gb.net": "whois.centralnic.com",
	".hu.com": "whois.centralnic.com",
	".no.com": "whois.centralnic.com",
	".qc.com": "whois.centralnic.com",
	".ru.com": "whois.centralnic.com",
	".sa.com": "whois.centralnic.com",
	".se.com": "whois.centralnic.com",
	".se.net": "whois.centralnic.com",
	".uk.com": "whois.centralnic.com",
	".uk.net": "whois.centralnic.com",
	".us.com": "whois.centralnic.com",
	".uy.com": "whois.centralnic.com",
	".za.com": "whois.centralnic.com",
	".com.au": "whois.ausregistry.net.au",
	".net.au": "whois.ausregistry.net.au",
	".org.au": "whois.ausregistry.net.au",
	".asn.au": "whois.ausregistry.net.au",
	".id.au":  "hois.ausregistry.net.au",
	".ac.uk":  "hois.ja.net",
	".gov.uk": "whois.ja.net",
	".museum": "whois.nic.museum",
	".asia":   "whois.internic.net",
	".info":   "whois.afilias.net",
	".name":   "whois.nic.name",
	".aero":   "whois.aero",
	".coop":   "whois.nic.coop",
	".com":    "whois.internic.net",
	".net":    "whois.internic.net",
	".org":    "whois.publicinterestregistry.net",
	".edu":    "whois.educause.net",
	".gov":    "whois.nic.gov",
	".int":    "whois.iana.org",
	".mil":    "whois.nic.mil",
	".biz":    "whois.neulevel.biz",
	".as":     "whois.nic.as",
	".ac":     "whois.nic.ac",
	".al":     "whois.ripe.net",
	".am":     "whois.amnic.net",
	".at":     "whois.nic.at",
	".au":     "whois.aunic.net",
	".az":     "whois.ripe.net",
	".ba":     "whois.ripe.net",
	".be":     "whois.dns.be",
	".bg":     "whois.ripe.net",
	".br":     "whois.nic.br",
	".by":     "whois.ripe.net",
	".ca":     "whois.cira.ca",
	".cc":     "whois.nic.cc",
	".cd":     "whois.nic.cd",
	".ch":     "whois.nic.ch",
	".cl":     "whois.nic.cl",
	".cn":     "whois.cnnic.cn",
	".cx":     "whois.nic.cx",
	".cy":     "whois.ripe.net",
	".cz":     "whois.ripe.net",
	".de":     "whois.denic.de",
	".dk":     "whois.dk-hostmaster.dk",
	".dz":     "whois.ripe.net",
	".ee":     "whois.eenet.ee",
	".eg":     "whois.ripe.net",
	".es":     "whois.ripe.net",
	".eu":     "whois.eu",
	".fi":     "whois.ripe.net",
	".fo":     "whois.ripe.net",
	".fr":     "whois.nic.fr",
	".gb":     "whois.ripe.net",
	".ge":     "whois.ripe.net",
	".gr":     "whois.ripe.net",
	".gs":     "whois.nic.gs",
	".hk":     "whois.hkirc.hk",
	".hr":     "whois.ripe.net",
	".hu":     "whois.ripe.net",
	".ie":     "whois.domainregistry.ie",
	".il":     "whois.isoc.org.il",
	".in":     "whois.inregistry.net",
	".ir":     "whois.nic.ir",
	".is":     "whois.ripe.net",
	".it":     "whois.nic.it",
	".jp":     "whois.jp",
	".kh":     "whois.nic.net.kh",
	".kr":     "whois.kr",
	".li":     "whois.nic.ch",
	".lt":     "whois.ripe.net",
	".lu":     "whois.dns.lu",
	".lv":     "whois.ripe.net",
	".ma":     "whois.ripe.net",
	".md":     "whois.ripe.net",
	".mk":     "whois.ripe.net",
	".ms":     "whois.nic.ms",
	".mt":     "whois.ripe.net",
	".mx":     "whois.nic.mx",
	".nl":     "whois.domain-registry.nl",
	".no":     "whois.norid.no",
	".nu":     "whois.nic.nu",
	".nz":     "whois.srs.net.nz",
	".pl":     "whois.dns.pl",
	".pt":     "whois.ripe.net",
	".ro":     "whois.ripe.net",
	".ru":     "whois.tcinet.ru",
	".se":     "whois.nic-se.se",
	".sg":     "whois.nic.net.sg",
	".si":     "whois.ripe.net",
	".sh":     "whois.nic.sh",
	".sk":     "whois.ripe.net",
	".sm":     "whois.ripe.net",
	".su":     "whois.ripn.net",
	".tc":     "whois.nic.tc",
	".tf":     "whois.nic.tf",
	".th":     "whois.thnic.net",
	".tj":     "whois.nic.tj",
	".tn":     "whois.ripe.net",
	".to":     "whois.tonic.to",
	".tr":     "whois.ripe.net",
	".tv":     "tvwhois.verisign-grs.com",
	".tw":     "whois.twnic.net",
	".ua":     "whois.ripe.net",
	".uk":     "whois.nic.uk",
	".us":     "whois.nic.us",
	".va":     "whois.ripe.net",
	".vg":     "whois.nic.vg",
	".ws":     "whois.website.ws",
	".vip":    "whois.nic.vip",
	".co":     "whois.nic.co",
	".top":    "whois.nic.top",
}

// DomainWhois 域名Whois 实体.
type DomainWhois struct {
	DomainName                 string
	RegistryDomainID           string
	RegistrarURL               string
	UpdatedDate                time.Time
	CreationDate               time.Time
	RegistryExpiryDate         time.Time
	Registrar                  string
	RegistrarIANAID            string
	RegistrarAbuseContactEmail string
	RegistrarAbuseContactPhone string
	DomainStatus               []string
	NameServer                 []string
	DNSSEC                     string
}

func (whois *DomainWhois) Println() {
	titleColor := color.New(color.FgWhite)
	titleColor = titleColor.Add(color.Bold)
	_, _ = titleColor.Println("DomainName: " + whois.DomainName)
	color.White("RegistryDomainID: " + whois.RegistryDomainID)
	color.White("RegistrarURL: " + whois.RegistrarURL)
	_, level, message := whois.UpdatedDateParse()
	color.Blue("UpdatedDate: " + message)
	_, level, message = whois.CreationDateParse()
	color.Blue("CreationDate: " + message)
	_, level, message = whois.RegistryExpiryDateParse()
	if level == 2 {
		color.Red("RegistryExpiryDate: " + message)
	} else if level == 1 {
		color.Yellow("RegistryExpiryDate: " + message)
	} else {
		color.Green("RegistryExpiryDate: " + message)
	}
	color.White("Registrar: " + whois.Registrar)
	color.White("RegistrarIANAID: " + whois.RegistrarIANAID)
	color.White("RegistrarAbuseContactEmail: " + whois.RegistrarAbuseContactEmail)
	color.White("RegistrarAbuseContactPhone: " + whois.RegistrarAbuseContactPhone)
	color.White("DomainStatus: \n")
	color.Blue(" - " + strings.Join(whois.DomainStatus, "\n - "))
	color.White("NameServer: \n")
	color.Blue(" - " + strings.Join(whois.NameServer, "\n - "))
	color.White("DNSSEC: " + whois.DNSSEC)
}

// UpdatedDateParse 修改时间; 过期天数, 危险级别 0.无危险 1.即将过期 2.已过期, 提示语.
func (whois *DomainWhois) UpdatedDateParse() (int, int, string) {
	return 0, 0, whois.UpdatedDate.Format("2006-01-02 15:04")
}

// CreationDateParse 创建时间; 过期天数, 危险级别 0.无危险 1.即将过期 2.已过期, 提示语.
func (whois *DomainWhois) CreationDateParse() (int, int, string) {
	return 0, 0, whois.CreationDate.Format("2006-01-02 15:04")
}

// RegistryExpiryDateParse 过期时间; 过期天数, 危险级别 0.无危险 1.即将过期 2.已过期, 提示语.
func (whois *DomainWhois) RegistryExpiryDateParse() (int, int, string) {
	duration := whois.RegistryExpiryDate.Sub(time.Now())
	days := duration.Hours() / 24
	var text = whois.RegistryExpiryDate.Format("2006-01-02 15:04")
	var level int
	if days < 0 {
		level = 2
		text = text + " ( Expired : " + strconv.Itoa(int(math.Abs(days))) + " Day )"
	} else if days < 15 {
		level = 1
		text = text + " ( Expiring : " + strconv.Itoa(int(days)) + " Day )"
	} else {
		level = 0
		text = text + " ( Remaining: " + strconv.Itoa(int(days)) + " Day )"
	}
	return int(days), level, text
}

// WhoisServer 获取域名对应的服务器.
func WhoisServer(host string) (*string, error) {
	for key, value := range whois {
		if strings.HasSuffix(host, key) {
			return &value, nil
		}
	}
	return nil, errors.New("domain not found")
}

// WhoisDomainFormat 解析域名让其适配Whois 格式.
func WhoisDomainFormat(host string) string {
	for key := range whois {
		if strings.HasSuffix(host, key) {
			var temp = strings.TrimSuffix(host, key)
			var hosts = strings.Split(temp, ".")
			return hosts[len(hosts)-1] + key
		}
	}
	var hosts = strings.Split(host, ".")
	if len(hosts) < 2 {
		return host
	}
	return hosts[len(hosts)-2] + "." + hosts[len(hosts)-1]
}

// MatchSuffix 匹配域名后缀.
func MatchSuffix(domain string) string {
	for suffix := range whois {
		if strings.HasSuffix(domain, suffix) {
			return suffix
		}
	}
	return ""
}

// Whois 域名Whois 文本信息.
func Whois(domain Domain) ([]string, error) {
	var host = domain.Name
	server, err := WhoisServer(host)
	if err != nil {
		return nil, err
	}
	var conn net.Conn
	conn, err = net.Dial("tcp", *server+":43")
	if err != nil {
		time.Sleep(time.Second / 4)
		conn, err = net.Dial("tcp", *server+":43")
		if err != nil {
			time.Sleep(time.Second / 4)
			conn, err = net.Dial("tcp", *server+":43")
			if err != nil {
				return nil, err
			}
		}
	}
	defer func(dial net.Conn) {
		err := dial.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	// Send the domain name query
	_, err = fmt.Fprintf(conn, "%s\r\n", host)
	if err != nil {
		time.Sleep(time.Second / 4)
		_, err = fmt.Fprintf(conn, "%s\r\n", host)
		if err != nil {
			time.Sleep(time.Second / 4)
			_, err = fmt.Fprintf(conn, "%s\r\n", host)
			if err != nil {
				return nil, err
			}
		}
	}

	// Read the response from the server
	var rows []string
	for i := 0; i < 10; i++ {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, ">>> ") {
				break
			}
			rows = append(rows, line)
		}
		if err = scanner.Err(); err != nil {
			time.Sleep(time.Second / 4)
			continue
		}
		break
	}

	if len(rows) <= 1 {
		var message = lo.IfF(len(rows) > 0, func() string { return rows[0] }).Else("")
		// 是否需要重试任务
		if isRetry(message) {
			time.Sleep(time.Second / 2)
			return Whois(domain)
		}
		return nil, errors.New(message)
	}
	return rows, nil
}

func isRetry(message string) bool {
	return strings.HasPrefix(message, "Queried interval is too short")
}

// ParseWhoisInfo 域名Whois实体信息.
func ParseWhoisInfo(rows []string) (*DomainWhois, error) {
	var domain = DomainWhois{
		DomainStatus: []string{},
		NameServer:   []string{},
	}
	for _, row := range rows {
		value := strings.TrimSpace(row)
		if strings.HasPrefix(value, "Domain Name:") {
			domain.DomainName = strings.TrimSpace(value[12:])
			continue
		}
		if strings.HasPrefix(value, "Registry Domain ID:") {
			domain.RegistryDomainID = strings.TrimSpace(value[19:])
			continue
		}
		if strings.HasPrefix(value, "Registrar URL:") {
			domain.RegistrarURL = strings.TrimSpace(value[14:])
			continue
		}
		if strings.HasPrefix(value, "Updated Date:") {
			var date time.Time
			date, err := times.Parse(time.RFC3339, strings.TrimSpace(value[13:]))
			if err != nil {
				return nil, err
			}
			domain.UpdatedDate = date
			continue
		}
		if strings.HasPrefix(value, "Creation Date:") {
			var date time.Time
			date, err := times.Parse(time.RFC3339, strings.TrimSpace(value[14:]))
			if err != nil {
				return nil, err
			}
			domain.CreationDate = date
			continue
		}
		if strings.HasPrefix(value, "Registration Time:") {
			var date time.Time
			date, err := times.Parse(time.DateTime, strings.TrimSpace(value[18:]))
			if err != nil {
				return nil, err
			}
			domain.CreationDate = date
			continue
		}
		if strings.HasPrefix(value, "Registry Expiry Date:") {
			var date time.Time
			date, err := times.Parse(time.RFC3339, strings.TrimSpace(value[21:]))
			if err != nil {
				return nil, err
			}
			domain.RegistryExpiryDate = date
			continue
		}
		if strings.HasPrefix(value, "Expiration Time:") {
			var date time.Time
			date, err := times.Parse(time.DateTime, strings.TrimSpace(value[16:]))
			if err != nil {
				return nil, err
			}
			domain.RegistryExpiryDate = date
			continue
		}
		if strings.HasPrefix(value, "Registrar:") {
			domain.Registrar = strings.TrimSpace(value[10:])
			continue
		}
		if strings.HasPrefix(value, "Registrar IANA ID:") {
			domain.RegistrarIANAID = strings.TrimSpace(value[18:])
			continue
		}
		if strings.HasPrefix(value, "Registrar Abuse Contact Email:") {
			domain.RegistrarAbuseContactEmail = strings.TrimSpace(value[30:])
			continue
		}
		if strings.HasPrefix(value, "Registrar Abuse Contact Phone:") {
			domain.RegistrarAbuseContactPhone = strings.TrimSpace(value[30:])
			continue
		}
		if strings.HasPrefix(value, "Domain Status:") {
			domain.DomainStatus = append(domain.DomainStatus, strings.TrimSpace(value[14:]))
			continue
		}
		if strings.HasPrefix(value, "Name Server:") {
			domain.NameServer = append(domain.NameServer, strings.TrimSpace(value[12:]))
			continue
		}
		if strings.HasPrefix(value, "DNSSEC:") {
			domain.DNSSEC = strings.TrimSpace(value[7:])
			continue
		}
	}
	if domain.DomainName == "" {
		return nil, errors.New("domain error")
	}
	return &domain, nil
}
