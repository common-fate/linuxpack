package packager

import (
	"compress/gzip"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/common-fate/linuxpack/pkg/control"
	"github.com/common-fate/linuxpack/pkg/packageset"
)

type Packager struct {
	S3Client     *s3.Client
	Bucket       string
	Description  string
	OutputFolder string
	Licence      string
	Vendor       string
	Channel      string
	Files        []string
}

func (p Packager) Package(ctx context.Context) error {
	// map of architecture -> package set
	sets := map[string]packageset.Set{}

	architectures := []string{"amd64", "arm64", "i386"}

	for _, arch := range architectures {
		// read the existing packages from the file in S3
		channelPath := filepath.Join("dists", p.Channel, "main", "binary-"+arch)
		releaseKey := filepath.Join(channelPath, "Packages")
		fmt.Printf("reading existing packages from s3://%s/%s\n", p.Bucket, releaseKey)
		existingReleaseResult, err := p.S3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &p.Bucket,
			Key:    &releaseKey,
		})
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			fmt.Printf("no packages found\n")
			sets[arch] = packageset.Set{}
			continue
		}

		if err != nil {
			return err
		}
		defer existingReleaseResult.Body.Close()

		sets[arch], err = packageset.ReadSet(existingReleaseResult.Body)
		if err != nil {
			return err
		}
	}

	err := os.RemoveAll(p.OutputFolder)
	if err != nil {
		return err
	}

	err = os.MkdirAll(p.OutputFolder, 0755)
	if err != nil {
		return err
	}

	for _, fileName := range p.Files {
		fileInfo, err := os.Stat(fileName)
		if err != nil {
			return err
		}
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		hash := sha1.New()
		if _, err := io.Copy(hash, file); err != nil {
			return err
		}
		sha1Hash := fmt.Sprintf("%x", hash.Sum(nil))

		file.Seek(0, io.SeekStart) // Reset file pointer to beginning for the next hash computation

		hash256 := sha256.New()
		if _, err := io.Copy(hash256, file); err != nil {
			return err
		}
		sha256Hash := fmt.Sprintf("%x", hash256.Sum(nil))

		tmpDir, err := os.MkdirTemp("", "package")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpDir)

		c1 := exec.Command("tar", "-z", "-xf", fileName, "--to-stdout", "control.tar.gz")
		c1.Stderr = os.Stderr

		c2 := exec.Command("tar", "-z", "-xf", "-", "-C", tmpDir)
		c2.Stderr = os.Stderr

		c2.Stdin, err = c1.StdoutPipe()
		if err != nil {
			return err
		}

		err = c2.Start()
		if err != nil {
			return err
		}
		err = c1.Run()
		if err != nil {
			return err
		}
		err = c2.Wait()
		if err != nil {
			return err
		}

		controlFile, err := os.Open(filepath.Join(tmpDir, "control"))
		if err != nil {
			return err
		}
		defer controlFile.Close()

		ctrl, err := control.Parse(controlFile)
		if err != nil {
			return err
		}

		pkg := packageset.Package{
			Package:       ctrl.Package,
			Version:       ctrl.Version,
			Licence:       p.Licence,
			Vendor:        p.Vendor,
			Architecture:  ctrl.Architecture,
			Maintainer:    ctrl.Maintainer,
			InstalledSize: ctrl.InstalledSize,
			Priority:      ctrl.Priority,
			Homepage:      ctrl.Homepage,
			Description:   ctrl.Description,
			Size:          fileInfo.Size(),
			SHA1:          sha1Hash,
			SHA256:        sha256Hash,
			Filename:      filepath.Join("pool", ctrl.Architecture, p.Channel, fileInfo.Name()),
		}

		pathToCopy := filepath.Join(p.OutputFolder, pkg.Filename)
		// copy the file over to the Filename path
		err = os.MkdirAll(filepath.Dir(pathToCopy), 0755)
		if err != nil {
			return err
		}

		destFile, err := os.Create(pathToCopy)
		if err != nil {
			return err
		}
		defer destFile.Close()

		file.Seek(0, io.SeekStart) // Reset file pointer to beginning for copying
		if _, err := io.Copy(destFile, file); err != nil {
			return err
		}

		set := sets[pkg.Architecture]
		set.Add(pkg)
		sets[pkg.Architecture] = set
	}

	var md5Checksums []Checksum
	var sha1Checksums []Checksum
	var sha256Checksums []Checksum

	for _, arch := range architectures {
		arch := arch
		channelPath := filepath.Join(p.OutputFolder, "dists", p.Channel, "main", "binary-"+arch)
		packagePath := filepath.Join(channelPath, "Packages")

		err = os.MkdirAll(channelPath, 0755)
		if err != nil {
			return err
		}

		packageFile, err := os.Create(packagePath)
		if err != nil {
			return err
		}
		defer packageFile.Close()

		// Create the gzipped version of the Packages file
		gzipPath := packagePath + ".gz"
		gzipFile, err := os.Create(gzipPath)
		if err != nil {
			return err
		}
		defer gzipFile.Close()

		gzipWriter := gzip.NewWriter(gzipFile)

		set := sets[arch]
		err = set.Write(io.MultiWriter(packageFile, gzipWriter))
		if err != nil {
			return err
		}

		err = gzipFile.Sync()
		if err != nil {
			return err
		}

		err = packageFile.Sync()
		if err != nil {
			return err
		}

		err = packageFile.Close()
		if err != nil {
			return fmt.Errorf("error closing package file: %w", err)
		}

		err = gzipWriter.Close()
		if err != nil {
			return err
		}

		err = gzipFile.Close()
		if err != nil {
			return err
		}

		// Calculate the md5, sha1, and sha256 sums of the packagePath and gzipPath

		paths := []string{packagePath, gzipPath}

		for _, path := range paths {

			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("error opening file %s: %w", path, err)
			}
			defer file.Close()

			hashMd5 := md5.New()
			hashSha1 := sha1.New()
			hashSha256 := sha256.New()

			if _, err := io.Copy(io.MultiWriter(hashMd5, hashSha1, hashSha256), file); err != nil {
				return err
			}

			fileInfo, err := file.Stat()
			if err != nil {
				return err
			}
			fmt.Printf("path: %s, size = %v\n", path, fileInfo.Size())

			md5Checksums = append(md5Checksums, Checksum{
				Sum:  fmt.Sprintf("%x", hashMd5.Sum(nil)),
				Size: fileInfo.Size(),
				Path: strings.TrimPrefix(path, "dist/dists/stable/"),
			})
			sha1Checksums = append(sha1Checksums, Checksum{
				Sum:  fmt.Sprintf("%x", hashSha1.Sum(nil)),
				Size: fileInfo.Size(),
				Path: strings.TrimPrefix(path, "dist/dists/stable/"),
			})
			sha256Checksums = append(sha256Checksums, Checksum{
				Sum:  fmt.Sprintf("%x", hashSha256.Sum(nil)),
				Size: fileInfo.Size(),
				Path: strings.TrimPrefix(path, "dist/dists/stable/"),
			})
		}
	}

	// create the Release file
	release := Release{
		Origin:        p.Vendor + " APT Repository",
		Label:         p.Vendor,
		Suite:         p.Channel,
		Codename:      p.Channel,
		Version:       "1.0",
		Architectures: architectures,
		Components:    "main",
		Description:   p.Description,
		Date:          time.Now().UTC(),
		MD5Sums:       md5Checksums,
		SHA1Sums:      sha1Checksums,
		SHA256Sums:    sha256Checksums,
	}

	releasePath := filepath.Join(p.OutputFolder, "dists", "stable", "Release")

	releaseFile, err := os.Create(releasePath)
	if err != nil {
		return err
	}
	defer releaseFile.Close()

	err = release.Write(releaseFile)
	if err != nil {
		return err
	}

	return nil
}
