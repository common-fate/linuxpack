package packageset

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadSet(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Set
		wantErr bool
	}{
		{
			name: "ok",
			input: `Package: granted
Version: 0.27.5
Licence: MIT
Vendor: Common Fate
Architecture: amd64
Maintainer: Chris Norman <chris@commonfate.io>
Installed-Size: 38697
Priority: optional
Homepage: https://granted.dev
Description: The easiest way to access your cloud.
Filename: pool/amd64/stable/granted_0.27.5_linux_amd64.deb
SHA1: ef07835809b153545ff323c2e903ae8647f5e849
SHA256: b88d280e2e94085503aa739142b35d4618692b13163d5771b1ab6fb7286113e9
Size: 14326932

`,
			want: Set{
				Packages: map[packageKey]Package{
					{Package: "granted", Version: "0.27.5"}: {
						Package:       "granted",
						Version:       "0.27.5",
						Licence:       "MIT",
						Vendor:        "Common Fate",
						Architecture:  "amd64",
						Maintainer:    "Chris Norman <chris@commonfate.io>",
						InstalledSize: "38697",
						Priority:      "optional",
						Homepage:      "https://granted.dev",
						Description:   "The easiest way to access your cloud.",
						Filename:      "pool/amd64/stable/granted_0.27.5_linux_amd64.deb",
						SHA1:          "ef07835809b153545ff323c2e903ae8647f5e849",
						SHA256:        "b88d280e2e94085503aa739142b35d4618692b13163d5771b1ab6fb7286113e9",
						Size:          14326932,
					},
				},
			},
		},

		{
			name: "multiple_packages",
			input: `Package: granted
Version: 0.27.5
Licence: MIT
Vendor: Common Fate
Architecture: amd64
Maintainer: Chris Norman <chris@commonfate.io>
Installed-Size: 38697
Priority: optional
Homepage: https://granted.dev
Description: The easiest way to access your cloud.
Filename: pool/amd64/stable/granted_0.27.5_linux_amd64.deb
SHA1: ef07835809b153545ff323c2e903ae8647f5e849
SHA256: b88d280e2e94085503aa739142b35d4618692b13163d5771b1ab6fb7286113e9
Size: 14326932



Package: granted
Version: 0.27.6
Licence: MIT
Vendor: Common Fate
Architecture: amd64
Maintainer: Chris Norman <chris@commonfate.io>
Installed-Size: 38697
Priority: optional
Homepage: https://granted.dev
Description: The easiest way to access your cloud.
Filename: pool/amd64/stable/granted_0.27.5_linux_amd64.deb
SHA1: ef07835809b153545ff323c2e903ae8647f5e849
SHA256: b88d280e2e94085503aa739142b35d4618692b13163d5771b1ab6fb7286113e9
Size: 14326932

`,
			want: Set{
				Packages: map[packageKey]Package{
					{Package: "granted", Version: "0.27.5"}: {
						Package:       "granted",
						Version:       "0.27.5",
						Licence:       "MIT",
						Vendor:        "Common Fate",
						Architecture:  "amd64",
						Maintainer:    "Chris Norman <chris@commonfate.io>",
						InstalledSize: "38697",
						Priority:      "optional",
						Homepage:      "https://granted.dev",
						Description:   "The easiest way to access your cloud.",
						Filename:      "pool/amd64/stable/granted_0.27.5_linux_amd64.deb",
						SHA1:          "ef07835809b153545ff323c2e903ae8647f5e849",
						SHA256:        "b88d280e2e94085503aa739142b35d4618692b13163d5771b1ab6fb7286113e9",
						Size:          14326932,
					},
					{Package: "granted", Version: "0.27.6"}: {
						Package:       "granted",
						Version:       "0.27.6",
						Licence:       "MIT",
						Vendor:        "Common Fate",
						Architecture:  "amd64",
						Maintainer:    "Chris Norman <chris@commonfate.io>",
						InstalledSize: "38697",
						Priority:      "optional",
						Homepage:      "https://granted.dev",
						Description:   "The easiest way to access your cloud.",
						Filename:      "pool/amd64/stable/granted_0.27.5_linux_amd64.deb",
						SHA1:          "ef07835809b153545ff323c2e903ae8647f5e849",
						SHA256:        "b88d280e2e94085503aa739142b35d4618692b13163d5771b1ab6fb7286113e9",
						Size:          14326932,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadSet(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ReadSet() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
