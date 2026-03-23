// Appends geosite list RU-WHITELIST from a plain-domain whitelist file (one host per line).
package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

const listCode = "RU-WHITELIST"

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "usage: %s <geosite.dat> <ru-whitelist.txt> <output.dat>\n", os.Args[0])
		os.Exit(2)
	}
	geositePath, whitelistPath, outPath := os.Args[1], os.Args[2], os.Args[3]

	raw, err := os.ReadFile(geositePath)
	if err != nil {
		fatal(err)
	}
	var list router.GeoSiteList
	if err := proto.Unmarshal(raw, &list); err != nil {
		fatal(fmt.Errorf("unmarshal geosite: %w", err))
	}

	newDomains, err := loadWhitelistDomains(whitelistPath)
	if err != nil {
		fatal(err)
	}

	idx := slices.IndexFunc(list.Entry, func(s *router.GeoSite) bool {
		return strings.EqualFold(s.GetCountryCode(), listCode)
	})
	if idx >= 0 {
		mergeDomains(list.Entry[idx], newDomains)
	} else {
		site := &router.GeoSite{CountryCode: listCode, Domain: newDomains}
		list.Entry = append(list.Entry, site)
	}

	slices.SortFunc(list.Entry, func(a, b *router.GeoSite) int {
		return strings.Compare(strings.ToLower(a.GetCountryCode()), strings.ToLower(b.GetCountryCode()))
	})

	out, err := proto.Marshal(&list)
	if err != nil {
		fatal(fmt.Errorf("marshal geosite: %w", err))
	}
	if err := os.WriteFile(outPath, out, 0o644); err != nil {
		fatal(err)
	}
}

func loadWhitelistDomains(path string) ([]*router.Domain, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	seen := make(map[string]struct{})
	var out []*router.Domain
	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line, _, _ := strings.Cut(sc.Text(), "#")
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "" {
			continue
		}
		if _, ok := seen[line]; ok {
			continue
		}
		seen[line] = struct{}{}
		out = append(out, &router.Domain{
			Type:  router.Domain_RootDomain,
			Value: line,
		})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%s: empty whitelist", path)
	}
	return out, nil
}

func mergeDomains(site *router.GeoSite, add []*router.Domain) {
	seen := make(map[string]struct{})
	for _, d := range site.Domain {
		if d == nil {
			continue
		}
		seen[domainKey(d)] = struct{}{}
	}
	for _, d := range add {
		k := domainKey(d)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		site.Domain = append(site.Domain, d)
	}
}

func domainKey(d *router.Domain) string {
	return fmt.Sprintf("%d:%s", d.GetType(), d.GetValue())
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
