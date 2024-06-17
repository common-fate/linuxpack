package control

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Control struct {
	Package       string
	Version       string
	Section       string
	Priority      string
	Architecture  string
	Maintainer    string
	InstalledSize string
	Homepage      string
	Description   string
}

func Parse(r io.Reader) (Control, error) {
	values := map[string]string{}

	sc := bufio.NewScanner(r)

	for sc.Scan() {
		line := sc.Text()
		before, after, found := strings.Cut(line, ": ")
		if !found {
			return Control{}, fmt.Errorf("invalid line: did not contain a \": \" separator: %s", line)
		}

		values[before] = after
	}

	res := Control{
		Package:       values["Package"],
		Version:       values["Version"],
		Section:       values["Section"],
		Priority:      values["Priority"],
		Architecture:  values["Architecture"],
		Maintainer:    values["Maintainer"],
		InstalledSize: values["Installed-Size"],
		Homepage:      values["Homepage"],
		Description:   values["Description"],
	}

	return res, nil
}
