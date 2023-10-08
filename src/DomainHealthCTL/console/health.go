package console

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/longyuan/domain.v3/client"
	"github.com/longyuan/lib.v3/ctl"
	"github.com/longyuan/lib.v3/message"
	"github.com/robfig/cron/v3"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func Whois(host string, original bool) error {
	domain, err := client.ParseDomain(host)
	if err != nil {
		return err
	}
	context, err := client.Whois(*domain)
	if err != nil {
		return err
	}
	if original {
		color.White(strings.Join(context, "\n"))
	} else {
		whoisInfo, err := client.ParseWhoisInfo(context)
		if err != nil {
			return err
		}
		whoisInfo.Println()
	}
	return nil
}

func SSL(host string) error {
	cert, err := client.SSL(host)
	if err != nil {
		return err
	}
	cert.Print()
	return nil
}

type DomainScan struct {
	domain    string
	message   string
	sslBefore time.Time
	sslAfter  time.Time

	whoisCreationDate       time.Time
	whoisUpdatedDate        time.Time
	whoisRegistryExpiryDate time.Time
}

func (domain *DomainScan) sslDays() float64 {
	sslDuration := domain.sslAfter.Sub(time.Now())
	return sslDuration.Hours() / 24
}

func (domain *DomainScan) whoisDays() float64 {
	whoisDuration := domain.whoisRegistryExpiryDate.Sub(time.Now())
	return whoisDuration.Hours() / 24
}

func (domain *DomainScan) sslDaysContext() string {
	sslDays := domain.sslDays()
	if sslDays <= -106751 {
		return "查询失败: " + domain.message
	} else if sslDays < 0 {
		return "已过期: " + strconv.Itoa(int(math.Abs(sslDays))) + "天"
	} else if sslDays < 15 {
		return "即将过期: " + strconv.Itoa(int(math.Abs(sslDays))) + "天"
	} else {
		return "剩余: " + strconv.Itoa(int(math.Abs(sslDays))) + "天"
	}
}

func (domain *DomainScan) whoisDaysContext() string {
	whoisDays := domain.whoisDays()
	if whoisDays <= -106751 {
		return "查询失败: " + domain.message
	} else if whoisDays < 0 {
		return "已过期: " + strconv.Itoa(int(math.Abs(whoisDays))) + "天"
	} else if whoisDays < 15 {
		return "即将过期: " + strconv.Itoa(int(math.Abs(whoisDays))) + "天"
	} else {
		return "剩余: " + strconv.Itoa(int(math.Abs(whoisDays))) + "天"
	}
}

func Scan(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var rows = strings.Split(strings.ReplaceAll(string(file), "\r\n", "\n"), "\n")
	color.Green("Scan Domain ....")
	var domain = client.Analysis(client.ParseDomains(rows))
	if err != nil {
		return err
	}
	var table [][]string
	for index, item := range domain {
		if item.Message != nil {
			table = append(table, []string{
				strconv.Itoa(index + 1), item.Name,
				"",
				"",
				"0", *item.Message,
			})
		} else {
			day, _, _ := item.Whois.RegistryExpiryDateParse()
			table = append(table, []string{
				strconv.Itoa(index + 1), item.Name,
				item.Whois.CreationDate.Format("2006-01-02"),
				item.Whois.RegistryExpiryDate.Format("2006-01-02"),
				strconv.Itoa(day), "",
			})
		}
	}
	ctl.PrintTable([]string{"序号", "域名", "Whois 创建日期", "Whois 过期日期", "Whois 剩余天数", "错误消息"}, table)

	table = [][]string{}
	var sslIndex = 0
	for _, item := range domain {
		for _, child := range *item.Child {
			if child.Message != nil {
				table = append(table, []string{
					strconv.Itoa(sslIndex + 1), child.Name,
					"", "", "0", *child.Message,
				})
			} else {
				day, _, _ := child.SSL.NotAfterDateParse()
				table = append(table, []string{
					strconv.Itoa(sslIndex + 1), child.Name,
					child.SSL.NotBefore.Format("2006-01-02"),
					child.SSL.NotAfter.Format("2006-01-02"),
					strconv.Itoa(day), "",
				})
			}
			sslIndex += 1
		}
	}
	ctl.PrintTable([]string{"序号", "域名", "SSL 创建日期", "SSL 过期日期", "SSL 剩余天数", "错误消息"}, table)

	return nil
}

func CronJob(configPath, backupCron, noticeConfig string) error {

	c := cron.New()
	_, err := c.AddFunc(backupCron, func() {
		err := func() error {
			// 读取文件
			file, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}
			var domainSSL []DomainScan
			var domainWhois []DomainScan
			var rows = strings.Split(strings.ReplaceAll(string(file), "\r\n", "\n"), "\n")
			var scanDomain []string
			var scanRootDomain []string
			for _, item := range rows {
				item = strings.TrimSpace(item)
				if item == "" {
					continue
				}
				if strings.HasPrefix(item, "#") {
					continue
				}
				scanDomain = append(scanDomain, item)
				var rootItem = client.WhoisDomainFormat(item)
				var exist = func() bool {
					for _, domain := range scanRootDomain {
						if domain == rootItem {
							return true
						}
					}
					return false
				}()
				if !exist {
					scanRootDomain = append(scanRootDomain, rootItem)
				}
			}
			// SSL
			var resultSSL []string
			for _, item := range scanDomain {
				color.Green(fmt.Sprintf("[域名] 执行SSL检查: %s", item))
				var domainScan = DomainScan{domain: item}

				func() {
					// SSL
					certificate, err := client.SSL(domainScan.domain)
					if err != nil {
						domainScan.message += fmt.Sprint(err)
						return
					}
					domainScan.sslAfter = certificate.NotAfter
				}()

				domainSSL = append(domainSSL, domainScan)
			}
			for _, item := range domainSSL {
				resultSSL = append(resultSSL, item.domain+" **SSL ( "+item.sslDaysContext()+" )** ")
			}

			// Whois
			var resultWhois []string
			for _, item := range scanRootDomain {
				color.Green(fmt.Sprintf("[域名] 执行Whois检查: %s", item))
				var domainScan = DomainScan{domain: item}

				//func() {
				//	// Whois
				//	whois, err := client.DomainWhoisInfo(domainScan.domain)
				//	if err != nil {
				//		domainScan.message += fmt.Sprint(err)
				//		return
				//	}
				//	domainScan.whoisRegistryExpiryDate = whois.RegistryExpiryDate
				//}()

				domainWhois = append(domainWhois, domainScan)
			}
			for _, item := range domainWhois {
				resultWhois = append(resultWhois, item.domain+" **Whois ( "+item.whoisDaysContext()+" )** ")
			}

			if len(resultSSL) > 0 {
				// 消息通知
				var content string
				for _, item := range resultSSL {
					content += DomainTemplateCPWeChat(item)
				}
				err = message.Push(message.DomainType, noticeConfig, "域名证书SSL 检查", content)
				if err != nil {
					return err
				}
			}
			if len(resultWhois) > 0 {
				// 消息通知
				var content string
				for _, item := range resultWhois {
					content += DomainTemplateCPWeChat(item)
				}
				err = message.Push(message.DomainType, noticeConfig, "域名Whois 检查", content)
				if err != nil {
					return err
				}
			}

			return nil
		}()
		if err != nil {
			color.Red(fmt.Sprint(err))
		}
	})
	if err != nil {
		return err
	}
	color.Blue(fmt.Sprintf("Cron (%s) Start Success ...", backupCron))
	color.Blue(fmt.Sprintf("ConfigPath: %s", configPath))
	color.Blue(fmt.Sprintf("Notice Config: %s", noticeConfig))
	c.Start()
	select {}
}

func DomainTemplateCPWeChat(value string) string {
	if strings.Index(value, "查询失败") > -1 {
		return "> <font color=\"red\">" + value + "</font>\n"
	}
	if strings.Index(value, "已过期") > -1 {
		return "> <font color=\"red\">" + value + "</font>\n"
	}
	if strings.Index(value, "即将过期") > -1 {
		return "> <font color=\"warning\">" + value + "</font>\n"
	}
	return "> " + value + "\n"
}
