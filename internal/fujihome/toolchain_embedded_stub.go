//go:build !release

package fujihome

func embeddedToolchain() (*Toolchain, error) {
	return nil, nil
}
