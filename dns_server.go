package main

import (
	"net"
	"os"
	"strings"

	"fmt"

	"github.com/miekg/dns"
)

func (app *App) startDNSServer() {
	fmt.Println("Starting dns server")

	dns.HandleFunc(".", app.serveDNS)

	go serve("tcp", *app.dnsBind)
	serve("udp", *app.dnsBind)
}

func serve(net string, bindAddr string) {
	err := dns.ListenAndServe(bindAddr, net, nil)

	if err != nil {
		panic(fmt.Sprintf("Failed to set "+net+" listener %s\n", err.Error()))
		os.Exit(1)
	}
}

func (app *App) serveDNS(w dns.ResponseWriter, r *dns.Msg) {
	var answer []dns.RR

	if dns.TypeA != r.Question[0].Qtype {
		m := new(dns.Msg)
		m.Authoritative = false
		m.SetRcode(r, dns.RcodeNotImplemented)
		w.WriteMsg(m)
	}

	ip, exists := getRecord(app, r.Question[0].Name)

	if false == exists {
		m := new(dns.Msg)
		m.Authoritative = false
		m.SetRcode(r, dns.RcodeNotImplemented)
		w.WriteMsg(m)
	}

	record := new(dns.A)
	record.Hdr = dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60}
	record.A = net.ParseIP(ip)

	answer = append(answer, record)
	setAnswer(w, r, answer)
}

func getRecord(app *App, name string) (string, bool) {
	domain := strings.ToLower(name)
	domain = strings.TrimSuffix(domain, ".")

	if ip, ok := app.records[domain]; ok {
		return ip, true
	}

	return "", false
}

func setAnswer(w dns.ResponseWriter, r *dns.Msg, data []dns.RR) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = data
	w.WriteMsg(m)
}
