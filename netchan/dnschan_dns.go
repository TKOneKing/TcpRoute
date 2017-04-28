package netchan

import (
	"log"
	"net"
	"time"

	"github.com/miekg/dns"
)

type dnsServer struct {
	Credit int //dns服务器信誉，根据信誉不同，查询得到的dns记录信誉也不同。
	Addr   net.UDPAddr
}

type ConfigDnsServer struct {
	Addr   string
	Credit int
}

type dnsDns struct {
	servers map[string]*dnsServer
}

func NewDnsDns() *dnsDns {
	return &dnsDns{
		servers: make(map[string]*dnsServer),
	}
}

func (d *dnsDns) setServers(servers []ConfigDnsServer) {

}

func (d *dnsDns) query(domain string, RecordChan chan *DnsRecord, ExitChan chan int) {
	dnschanDnsQuery(d.servers, domain, RecordChan, ExitChan)
}

func dnschanDnsQuery(servers map[string]*dnsServer, domain string, RecordChan chan *DnsRecord, ExitChan chan int) {
	select {
	case <-ExitChan:
		return
	default:
	}
	myexitChan := make(chan int)
	defer func() { close(myexitChan) }()

	// 打开端口
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Printf("为dns请求打开udp失败，%v", err)
		return
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// dns 请求
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true
	mData, err := m.Pack()
	if err != nil {
		log.Printf("生成dns请求包失败，%v", err)
		return
	}

	// 另开一个线程发出查询
	go func() {
		for k, v := range servers {
			if _, err := conn.WriteToUDP(mData, &v.Addr); err != nil {
				log.Printf("向%v发送dns请求失败，%v", k, err)
			}
		}
	}()

	// 等待关闭
	go func() {
		select {
		case <-ExitChan:
		case <-myexitChan:
		}
		conn.Close()
	}()

	// 接收查询结果并输入到信道
	buf := make([]byte, 1500)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		nbuf := buf[:n]

		r := new(dns.Msg)
		if err := r.Unpack(nbuf); err != nil {
			log.Printf("解 DNS 包失败,%v", err)
			continue
		}

		if r.Id != m.Id {
			log.Printf("错误的dns id。")
			continue
		}

		v := servers[addr.IP.String()]
		if v == nil {
			log.Printf("未知的服务器 %v 回应，忽略。", *addr)
			continue
		}

		for _, a := range r.Answer {
			dnsA, ok := a.(*dns.A)
			if ok != true || dnsA == nil {
				log.Printf("内部错误，a=%v,err=%v", dnsA, err)
			}
			select {
			case RecordChan <- &DnsRecord{
				Ip:     dnsA.A.String(),
				Credit: v.Credit,
			}:
			case <-ExitChan:
				return
			}
		}
	}
}
