package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/miekg/dns"
)

func main() {
	file, _ := os.Open("tlds-alpha-by-domain.txt")
	defer file.Close()

	client := new(dns.Client)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tld := scanner.Text()

		if strings.HasPrefix(tld, "#") || strings.HasPrefix(tld, "XN--") { // Comment or punycode
			continue
		}

		fmt.Println(tld + ":")
		aaaacounter := 0

		for _, rr := range QueryNSRecords(tld, client) {
			hasAAAA := HasAAAA(rr, client)
			fmt.Printf("  %-45s -> %t\n", rr, hasAAAA)

			if hasAAAA {
				aaaacounter++

				if aaaacounter == 3 {
					fmt.Println("  Skipping rest as 3 IPv6 enabled nameservers are sufficient")
					break
				}
			}
		}

	}
}

func QueryNSRecords(tld string, client *dns.Client) []string {
	message := new(dns.Msg)
	message.SetQuestion(tld+".", dns.TypeNS)

	result := make([]string, 0)

	for i := 0; i < 3; i++ {
		answer, _, _ := client.Exchange(message, "[2001:4860:4860::8888]:53")

		if answer == nil {
			continue
		}

		for _, record := range answer.Answer {
			if rr, ok := record.(*dns.NS); ok {
				result = append(result, rr.Ns)
			}
		}

		return result
	}

	return result
}

func HasAAAA(hostname string, client *dns.Client) bool {
	message := new(dns.Msg)
	message.SetQuestion(hostname, dns.TypeAAAA)

	for i := 0; i < 5; i++ {
		answer, _, _ := client.Exchange(message, "[2001:4860:4860::8888]:53")

		if answer == nil {
			fmt.Printf("timeout #%d for %s\n", i, hostname)
			continue
		}

		for _, record := range answer.Answer {
			if _, ok := record.(*dns.AAAA); ok {
				return true
			}
		}
	}

	fmt.Println("Giving up")

	return false
}
