package make

import (
	"bufio"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var targetRegexp = regexp.MustCompile(`^([[:alnum:]][^#:]*):.*$`)

type parser struct {
	targets []string
	scanner *bufio.Scanner
}

func ParseDataBase(output string) []string {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanLines)

	parser := &parser{
		scanner: scanner,
	}

	targets := parser.parse()
	return targets
}

func (p *parser) parse() []string {
	log.Info("start parse")
	for p.scanner.Scan() {
		if strings.HasPrefix(p.scanner.Text(), "# Make data base, printed on ") {
			p.parseDB()
		}
	}

	return p.targets
}

func (p *parser) parseDB() {
	log.Debug("start parseDB")
	for p.scanner.Scan() {
		if p.scanner.Text() == "# Files" {
			p.scanner.Scan() // skip the first empty line
			p.parseEntries()
			return
		}
	}
}

func (p *parser) parseEntries() {
	log.Debug("start parseEntries")
	for p.scanner.Scan() {
		line := p.scanner.Text()
		switch {
		case line == "# files hash-table stats:":
			return
		case line == "# Not a target:":
			p.skipUntilNextEntry()
		case targetRegexp.MatchString(line):
			target := targetRegexp.FindStringSubmatch(line)[1]
			log.Info("match target:", target)
			p.targets = append(p.targets, target)
			p.skipUntilNextEntry()
		}
	}
}


func (p *parser) skipUntilNextEntry() {
	for p.scanner.Scan() {
		if p.scanner.Text() == "" {
			return
		}
	}
}