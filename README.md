# linuxpack

Linux packaging automation, used to publish public APT repositories using AWS S3 and CloudFront. Currently used to publish an APT repository for [Granted](https://granted.dev).

## Usage

Packaging:

```
go run cmd/main.go package -f granted_0.27.4_linux_amd64.deb -f granted_0.27.4_linux_386.deb -f granted_0.27.4_linux_arm64.deb --licence MIT --vendor "Common Fate" --channel stable --out dist --bucket example-bucket
```

This will download existing `Packages` files from the S3 bucket, merge the file with the new releases to be uploaded, and create a folder similar to the below, ready to be synced with an S3 bucket:

```
❯ tree dist
dist
├── dists
│   └── stable
│       ├── Release
│       └── main
│           ├── binary-amd64
│           │   ├── Packages
│           │   └── Packages.gz
│           ├── binary-arm64
│           │   ├── Packages
│           │   └── Packages.gz
│           └── binary-i386
│               ├── Packages
│               └── Packages.gz
└── pool
    ├── amd64
    │   └── stable
    │       └── granted_0.27.4_linux_amd64.deb
    ├── arm64
    │   └── stable
    │       └── granted_0.27.4_linux_arm64.deb
    └── i386
        └── stable
            └── granted_0.27.4_linux_386.deb
```

Prior to uploading you'll need to sign the `Release` file:

```bash
gpg -abs -u <signing key ID> -o dist/dists/stable/Release.gpg dist/dists/stable/Release

cat dist/dists/stable/Release | gpg -abs -u <signing key ID> --clearsign > dist/dists/stable/InRelease
```

Then, upload the release:

```bash
aws s3 cp --recursive dist s3://example-bucket
```

## Acknowledgements

Our APT implementation is inspired by [deb-s3](https://github.com/deb-s3/deb-s3).
