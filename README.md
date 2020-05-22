# TLDIPv6check

Application checking the TLD nameservers (de, com, to, ws, club, dating, etc) if their NS glue records resolve to an IPv6 address. It was quickly hacked together as you can see due to the missing error handling. Writing this README took longer than writing the code.   
The gathered result is stored in `results.txt`.

## Why

I wanted to create a new category for my side project https://ipv6-overview.xyz . I wanted to list TLD nameservers without IPv6 support. Adding over 1500 new check targets is of course too much so I wanted to pick only the black sheeps. But to find them this helper application had to be developed.

## How it works

1) Read the TLDs line by line from the file `tlds-alpha-by-domain.txt` I've downloaded from [IANA here](https://www.icann.org/resources/pages/tlds-2012-02-25-en)
2) Query Googles Public DNS resolver for the glue nameservers of the TLD
3) Resolve each returned nameserver for AAAA records. If 3 of them successfully resolve the rest is skipped
4) Print the gathered information

**Note: Just because they resolve with an IPv6 address doesn't mean that they also respond to DNS queries on it. Checking if they do is out of scope. But I guess they do.**

Sample output:

```
AFAMILYCOMPANY:
  ac1.nstld.com.                                -> true
  ac2.nstld.com.                                -> true
  ac4.nstld.com.                                -> true
  Skipping rest as 3 IPv6 enabled nameservers are sufficient
```

true = AAAA record returned, false = No AAAA record returned

So far I've only found a tiny fractiob of TLDs without IPv6 nameservers:

- CK
- DJ
- FK
- GF
- HM
- TO
- UZ
- WS

Subjectively > 95% of the TLD nameservers are reachable over IPv6. This honestly surprises me a bit. My estimate was only 60 to 70%. After thinking about it for a minute I guess the smart, technologically fit people at IANA made IPv6 mandatory. After searching the interwebz I quickly found an IANA document about [nameserver requirements](https://www.iana.org/help/nameserver-requirements). Quote from there:

> The minimal set of requisite glue records is considered to be:
>
>    One A record, if all authoritative name servers are in-bailiwick of the parent zone; and,
>    One AAAA record, if there are any IPv6-capable authoritative name servers and all IPv6- capable authoritative name servers are in-bailiwick of the parent zone.

So yeah it's kind of mandatory whis is of course nice. IPv6 is the current gen IP protocol. Not using it would be real dumb.
