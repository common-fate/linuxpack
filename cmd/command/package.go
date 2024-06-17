package command

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/common-fate/linuxpack/pkg/packager"
	"github.com/urfave/cli/v2"
)

var Package = cli.Command{
	Name: "package",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "file", Aliases: []string{"f"}, FilePath: "files to package", Required: true},
		&cli.StringFlag{Name: "licence", FilePath: "the licence to apply to the packages"},
		&cli.StringFlag{Name: "vendor", FilePath: "the vendor to apply to the packages"},
		&cli.StringFlag{Name: "bucket", FilePath: "the S3 bucket to store releases in"},
		&cli.StringFlag{Name: "channel", FilePath: "the release channel to use", Required: true},
		&cli.PathFlag{Name: "out", FilePath: "output directory", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}

		p := packager.Packager{
			OutputFolder: c.Path("out"),
			Licence:      c.String("licence"),
			Vendor:       c.String("vendor"),
			Channel:      c.String("channel"),
			Files:        c.StringSlice("file"),
			S3Client:     s3.NewFromConfig(cfg),
			Bucket:       c.String("bucket"),
			Description:  c.String("description"),
		}

		return p.Package(ctx)
	},
}
