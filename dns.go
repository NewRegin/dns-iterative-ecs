package dnsiterativeecs

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"math/rand"
	"net"
	"time"
)

//顺序：本地——>根域名服务器——>顶级域名服务器——>权威域名服务器；我们维护权威域名服务器的信息
var (
	DnsRoots = []string{
		"a.root-servers.net", "b.root-servers.net",
		"c.root-servers.net", "d.root-servers.net",
		"e.root-servers.net", "f.root-servers.net",
		"g.root-servers.net", "h.root-servers.net",
		"i.root-servers.net", "j.root-servers.net",
		"k.root-servers.net", "l.root-servers.net",
		"m.root-servers.net",
	}
)

var (
	ErrNoNameservers = errors.New("No nameservers registered for that domain")
	ErrUnhandled     = errors.New("Unknown error")
)

func getECSOption(ip string) *dns.OPT {
	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_SUBNET)
	e.Code = dns.EDNS0SUBNET

	e.SourceScope = 0
	e.Address = net.ParseIP(ip)

	e.Family = 1 // IP4
	e.SourceNetmask = net.IPv4len * 8
	if e.Address.To4() == nil {
		e.Family = 2 // IP6
		e.SourceNetmask = net.IPv6len * 8
	}
	o.Option = append(o.Option, e)
	return o
}

func lookup(cl *dns.Client, ip, name, server string) error {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = false

	//合并
	// msg.Question = []dns.Question{
	// 	dns.Question{
	// 		name, dns.TypeA, dns.ClassINET,
	// 	},
	// }
	msg.SetQuestion(dns.Fqdn(name), dns.TypeA)
	opt := getECSOption(ip)
	msg.Extra = append(msg.Extra, opt)

	response, _, err := cl.Exchange(msg, server+":53")
	//记录解析失败时的Rcode(Rcode=16表示EDNS问题，可以
	//结合返回结果是否包含OPT EDNS字段来判断NS服务器是否支持ECS)
	if err != nil {
		return err
	}
	if len(response.Answer) == 0 {
		if len(response.Ns) == 0 {
			return ErrNoNameservers
		} else {
			ns, ok := response.Ns[rand.Intn(len(response.Ns))].(*dns.NS)
			if !ok {
				return ErrUnhandled
			}
			return lookup(cl, ip, name, string(ns.Ns[0:len(ns.Ns)-1]))
		}
	} else {
		//todo
		fmt.Println(response.Extra)
		fmt.Println(response.MsgHdr.Rcode)
		fmt.Println(response.Answer)
	}
	return nil
}

func Lookup(ip, name string) error {
	rand.Seed(time.Now().Unix())
	cl := new(dns.Client)
	cl.Timeout = time.Second * 10
	// cl.SingleInflight = true

	return lookup(cl, ip, name, DnsRoots[rand.Intn(len(DnsRoots))])

}
