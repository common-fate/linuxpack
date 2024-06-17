package packager

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Release struct {
	Origin        string
	Label         string
	Suite         string
	Codename      string
	Version       string
	Architectures []string
	Components    string
	Description   string
	Date          time.Time
	MD5Sums       []Checksum
	SHA1Sums      []Checksum
	SHA256Sums    []Checksum
}

func (r *Release) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "Origin: %s\n", r.Origin)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Label: %s\n", r.Label)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Suite: %s\n", r.Suite)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Codename: %s\n", r.Codename)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Version: %s\n", r.Version)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Architectures: %s\n", strings.Join(r.Architectures, " "))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Components: %s\n", r.Components)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Description: %s\n", r.Description)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Date: %s\n", r.Date.Format(time.RFC1123))
	if err != nil {
		return err
	}

	var md5sums []string
	var sha1sums []string
	var sha256sums []string

	for _, s := range r.MD5Sums {
		md5sums = append(md5sums, fmt.Sprintf(" %s %v %s", s.Sum, s.Size, s.Path))
	}

	for _, s := range r.SHA1Sums {
		sha1sums = append(sha1sums, fmt.Sprintf(" %s %v %s", s.Sum, s.Size, s.Path))
	}

	for _, s := range r.SHA256Sums {
		sha256sums = append(sha256sums, fmt.Sprintf(" %s %v %s", s.Sum, s.Size, s.Path))
	}

	_, err = fmt.Fprintf(w, "MD5Sum:\n%s\n", strings.Join(md5sums, "\n"))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "SHA1:\n%s\n", strings.Join(sha1sums, "\n"))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "SHA256:\n%s\n", strings.Join(sha256sums, "\n"))
	if err != nil {
		return err
	}

	return nil
}

type Checksum struct {
	Sum  string
	Size int64
	Path string
}
