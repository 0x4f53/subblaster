[![Golang](https://img.shields.io/badge/Golang-fff.svg?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-purple?style=flat-square&logo=libreoffice)](LICENSE)
[![Latest Version](https://img.shields.io/github/v/tag/0x4f53/subblaster?label=Version&style=flat-square&logo=semver)](https://github.com/0x4f53/subblaster/releases)

<img src=logo-small.gif alt="subblaster logo">

# subblaster
Super-fast multi-source subdomain bruteforcer in Go.

<img src = preview.gif alt="dnscovery preview" width = "500dp">

**Note:** This is not a public subdomain enumerator and is not an efficient way to get pre-captured subdomains. If you need fast enumeration, please use pre-existing tools like [amass](https://github.com/owasp-amass/amass), [sublister](https://github.com/aboul3la/Sublist3r) etc.

### What is this then?

There are several domains whose subdomains are present in DNS records but aren't caught by popular enumeration services. These don't appear online due to them not being scraped by crawlers that these providers / security companies deploy. This tool is an attempt to maximize the speed of discovering them while minimizing the time taken.

## Features
- Customizable multi-source wordlists ([TheRook's subbrute](https://github.com/TheRook/subbrute), [Daniel Miessler's seclists](https://github.com/danielmiessler/SecLists) and more!)
- Multithreaded bruteforcing using Golang
- Multi-resolver subdomain resolution and port scanning in-built
- Multiple inputs, multiple outputs
<!-- - Terraform integration for powerful yet cheap bruteforcing-->

## Usage 

```bash
# to build the program
go build

./subblaster 0x4f.in
```

Examples:

### Generate paired JSON outputs ({"subdomain": "www.example.com", "domain":"example.com"})

```bash                S U B B L A S T E R  
        (https://github.com/0x4f53/subblaster)
        A fast subdomain bruteforcer in Golang.


[⟳] Generating batches for bruteforcing...


[✓] Batching complete! Generated 1192 batches

[+] Bruteforcing...

./subblaster -p 0x4f.in
...
# In 0x4f.in.json
{"subdomain":"blog.0x4f.in","domain":"0x4f.in"}
{"subdomain":"www.0x4f.in","domain":"0x4f.in"}
```

### Refresh all seclists and delete cache

```bash
./subblaster -r

                S U B B L A S T E R  
        (https://github.com/0x4f53/subblaster)
        A fast subdomain bruteforcer in Golang.

[✓] Deleted all cache data
[↓] Downloading wordlists mentioned in lists.yaml
 - onelistforallshort.txt [2.82MB / 0b] ╢░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░╟  10.84kB/s

 - onelistforallshort.txt [12.19MB / 0b] ╢░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░╟  4.28kB/s
 - dns-Jhaddix.txt [10.40kB / 0b] ╢░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░╟  0b/s
 - alexaTop1mAXFRcommonSubdomains.txt [378.92kB / 0b] 
...
```

## Credits

- [Assetnote Wordlists](https://wordlists.assetnote.io/)
- [six2dez/OneListForAll](https://github.com/six2dez/OneListForAll)
- [fuzzdb-project/fuzzdb](https://github.com/fuzzdb-project/fuzzdb)
- [TheRook/subbrute](https://github.com/TheRook/subbrute)
- [danielmiessler/seclists](https://github.com/danielmiessler/SecLists)

<!--[Terraform Exec](https://github.com/hashicorp/terraform-exec) is property of Hashicorp Terraform.-->

The animated logo is derived from work by [Ryan Whiteside](https://flickr.com/whytseyed/).

## License

Multimedia licensed under [![License: CC BY-NC-SA 4.0](https://licensebuttons.net/l/by-nc-sa/4.0/80x15.png)](https://creativecommons.org/licenses/by-nc-sa/4.0/) 

[Copyright © 2024 Owais Shaikh](LICENSE)

## Donate

[Click here to donate](https://github.com/sponsors/0x4f53). It incentivizes me to develop more.
