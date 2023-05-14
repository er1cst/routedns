package rdns

import (
	"bufio"
	"os"

	"github.com/miekg/dns"
)

type labelTable map[string]labelTable

func (lm labelTable) add(labels []string) {
	if len(labels) == 0 {
		return
	}
	m := lm
	for i, v := range labels {
		next, ok := m[v]
		if !ok {
			next = make(labelTable)
			m[v] = next
		} else if len(next) == 0 {
			break
		} else if i == len(labels)-1 {
			m[v] = make(labelTable)
			break
		}
		m = next
	}
}

func (lm labelTable) test(labels []string) bool {
	if len(lm) == 0 {
		return false
	}
	m := lm
	for _, v := range labels {
		if len(m) == 0 {
			return true
		}
		next, ok := m[v]
		if ok {
			m = next
			continue
		}
		return false
	}
	return len(m) == 0
}

type Matcher struct {
	labels labelTable
}

func reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func NewMatcherFromFile(path string) (*Matcher, error) {
	m := new(Matcher)
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m.labels = make(labelTable)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m.add(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Matcher) add(domain string) {
	if _, ok := dns.IsDomainName(domain); !ok {
		return
	}
	labels := dns.SplitDomainName(domain)
	reverse(labels)
	m.labels.add(labels)
}

func (m *Matcher) Match(name string) bool {
	labels := dns.SplitDomainName(name)
	reverse(labels)
	return m.labels.test(labels)
}
