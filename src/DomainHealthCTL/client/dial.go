package client

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/longyuan/lib.v3/times"
	"github.com/samber/lo"
	"math"
	"net/url"
	"strconv"
	"time"
)

type X509Certificate struct {
	*x509.Certificate
}

// SSL 域名证书状态查询.
func SSL(host string) (*X509Certificate, error) {
	value, err := url.Parse("scheme://" + host)
	if err != nil {
		return nil, err
	}
	host = fmt.Sprintf("%s:%s", value.Hostname(), lo.If(value.Port() == "", "443").Else(value.Port()))
	dial, err := tls.Dial("tcp", host, nil)
	if err != nil {
		return nil, err
	}
	defer func(dial *tls.Conn) {
		err = dial.Close()
		if err != nil {
			panic(err)
		}
	}(dial)
	result := dial.ConnectionState()
	if len(result.PeerCertificates) <= 0 {
		return nil, errors.New("domain Not SSL")
	}
	return &X509Certificate{result.PeerCertificates[0]}, nil
}

// Print 打印证书信息.
func (cert *X509Certificate) Print() {
	colorPrint := color.New()
	colorPrint.Add(color.Bold)
	colorPrint.Add(color.FgWhite)
	_, _ = colorPrint.Println(fmt.Sprint("Subject:", cert.Subject.CommonName))
	_, _ = colorPrint.Println(fmt.Sprint("Issuer:", cert.Issuer.CommonName))
	_, _ = colorPrint.Println(fmt.Sprint("Serial Number:", cert.SerialNumber))

	// 开始日期
	color.Blue(fmt.Sprint("Not Before: ", times.In(cert.NotBefore).Format(time.DateTime)))
	// 截至日期
	_, level, text := cert.NotAfterDateParse()
	if level == 2 {
		color.Red("Not After: " + text)
	} else if level == 1 {
		color.Yellow("Not After: " + text)
	} else {
		color.Green("Not After: " + text)
	}

	color.White(fmt.Sprint("Signature Algorithm:", cert.SignatureAlgorithm))
	color.White(fmt.Sprint("Key Usage:", cert.KeyUsage))
	color.White(fmt.Sprint("Extended Key Usage:", cert.ExtKeyUsage))
	color.White(fmt.Sprint("Subject Alternative Names:"))
	for _, san := range cert.DNSNames {
		color.White(fmt.Sprint(" - ", san))
	}
	color.White(fmt.Sprint("Is CA:", cert.IsCA))
	color.Blue(fmt.Sprint("Authority Information Access (AIA):"))
	for _, aia := range cert.IssuingCertificateURL {
		color.White(fmt.Sprint(" - ", aia))
	}
	color.White(fmt.Sprint("Basic Constraints:"))
	color.White(fmt.Sprint("- Is CA:", cert.BasicConstraintsValid))
	color.White(fmt.Sprint("- Max Path Length:", cert.MaxPathLen))
	color.White(fmt.Sprint("- Is Certificate Authority:", cert.IsCA))
	color.White(fmt.Sprint("Subject Key ID:", cert.SubjectKeyId))
	color.White(fmt.Sprint("Authority Key ID:", cert.AuthorityKeyId))
	color.Blue(fmt.Sprint("Certificate Policies:"))
	for _, policy := range cert.PolicyIdentifiers {
		color.White(fmt.Sprint(" - ", policy.String()))
	}
	color.White(fmt.Sprint("Signature:", cert.Signature))
	color.White(fmt.Sprint("Public Key Algorithm:", cert.PublicKeyAlgorithm.String()))
	color.White(fmt.Sprint("Public Key Size (bits):", cert.PublicKey.(*rsa.PublicKey).N.BitLen()))
}

// NotAfterDateParse 解析证书到期时间; 过期天数, 危险级别 0.无危险 1.即将过期 2.已过期, 提示语.
func (cert *X509Certificate) NotAfterDateParse() (int, int, string) {
	var notAfter = times.In(cert.NotAfter)
	duration := notAfter.Sub(time.Now())
	days := duration.Hours() / 24
	var level int
	var text = notAfter.Format(time.DateTime)
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
