// SPDX-License-Identifier: Apache-2.0
package mutate

type Config struct {
	TLSCertPath  string   `yaml:"tlscertpath"`
	TLSKeyPath   string   `yaml:"tlskeypath"`
	Port         string   `yaml:"port"`
	IgnoreImages []string `yaml:"ignoreimages"`
}

var Configuration Config
