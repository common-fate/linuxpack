package packageset

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type Package struct {
	Package       string
	Version       string
	Licence       string
	Vendor        string
	Architecture  string
	Maintainer    string
	InstalledSize string
	Depends       string
	Priority      string
	Homepage      string
	Description   string
	Filename      string
	SHA1          string
	SHA256        string
	Size          int64
}

type packageKey struct {
	Package string
	Version string
}

type Set struct {
	Packages map[packageKey]Package
}

func (s *Set) Add(p Package) {
	if s.Packages == nil {
		s.Packages = map[packageKey]Package{}
	}

	key := packageKey{
		Package: p.Package,
		Version: p.Version,
	}
	s.Packages[key] = p
}

func (s *Set) Write(w io.Writer) error {
	var packages []Package

	for _, p := range s.Packages {
		packages = append(packages, p)
	}

	sortPackages(packages)

	for _, p := range s.Packages {
		_, err := fmt.Fprintf(w, "Package: %s\n", p.Package)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "Version: %s\n", p.Version)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "Licence: %s\n", p.Licence)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "Vendor: %s\n", p.Vendor)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "Architecture: %s\n", p.Architecture)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "Maintainer: %s\n", p.Maintainer)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "Installed-Size: %s\n", p.InstalledSize)
		if err != nil {
			return err
		}

		if p.Depends != "" {
			_, err = fmt.Fprintf(w, "Depends: %s\n", p.Depends)
			if err != nil {
				return err
			}
		}

		_, err = fmt.Fprintf(w, "Priority: %s\n", p.Priority)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "Homepage: %s\n", p.Homepage)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "Description: %s\n", p.Description)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "Filename: %s\n", p.Filename)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "SHA1: %s\n", p.SHA1)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "SHA256: %s\n", p.SHA256)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "Size: %v\n", p.Size)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "\n")
		if err != nil {
			return err
		}

	}

	return nil
}

// sortPackages sorts packages by Package, Version and Architecture.
func sortPackages(p []Package) {
	slices.SortFunc(p, func(a, b Package) int {
		if a.Package < b.Package {
			return -1
		}

		if a.Package > b.Package {
			return 1
		}

		if a.Version < b.Version {
			return -1
		}

		if a.Version > b.Version {
			return 1
		}

		if a.Architecture < b.Architecture {
			return -1
		}

		if a.Architecture > b.Architecture {
			return 1
		}

		return 0
	})
}

func ReadSet(r io.Reader) (Set, error) {
	sc := bufio.NewScanner(r)
	var packages []Package

	var currentPackage *Package

	var lineNum int
	for sc.Scan() {
		lineNum++
		line := sc.Text()

		if line == "" && currentPackage != nil {
			packages = append(packages, *currentPackage)
			currentPackage = nil
		}

		if line == "" {
			continue
		}

		before, after, found := strings.Cut(line, ": ")
		if !found {
			return Set{}, fmt.Errorf("invalid line %v: did not contain a \": \" separator: %q", lineNum, line)
		}

		if currentPackage == nil {
			currentPackage = &Package{}
		}

		switch before {
		case "Package":
			currentPackage.Package = after
		case "Version":
			currentPackage.Version = after
		case "Licence":
			currentPackage.Licence = after
		case "Vendor":
			currentPackage.Vendor = after
		case "Architecture":
			currentPackage.Architecture = after
		case "Maintainer":
			currentPackage.Maintainer = after
		case "Installed-Size":
			currentPackage.InstalledSize = after
		case "Priority":
			currentPackage.Priority = after
		case "Homepage":
			currentPackage.Homepage = after
		case "Description":
			currentPackage.Description = after
		case "Filename":
			currentPackage.Filename = after
		case "SHA1":
			currentPackage.SHA1 = after
		case "SHA256":
			currentPackage.SHA256 = after
		case "Size":
			sizeInt, err := strconv.ParseInt(after, 10, 0)
			if err != nil {
				return Set{}, fmt.Errorf("error parsing size %q: %w", after, err)
			}
			currentPackage.Size = sizeInt
		}
	}

	var set Set

	for _, p := range packages {
		set.Add(p)
	}

	return set, nil
}
