// Copyright 2017 EcoSystem Software LLP

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utilities

import (
	"errors"
)

type Bundles []string

//Config is the basic structure of the config.json file
type Config struct {
	PgSuperUser              string  `json:"pgSuperUser"`
	PgDBName                 string  `json:"pgDBName"`
	PgPort                   string  `json:"pgPort"`
	PgServer                 string  `json:"pgServer"`
	PgDisableSSL             bool    `json:"pgDisableSSL"`
	ApiPort                  string  `json:"apiPort"`
	WebsitePort              string  `json:"websitePort"`
	AdminPanelPort           string  `json:"adminPanelPort"`
	AdminPanelServeDirectory string  `json:"adminPanelServeDirectory"`
	PublicSiteSlug           string  `json:"publicSiteSlug"`
	PrivateSiteSlug          string  `json:"privateSiteSlug"`
	SmtpHost                 string  `json:"smtpHost"`
	SmtpPort                 string  `json:"smtpPort"`
	SmtpUserName             string  `json:"smtpUserName"`
	SmtpFrom                 string  `json:"smtpFrom"`
	EmailFrom                string  `json:"emailFrom"`
	JWTRealm                 string  `json:"jwtRealm"`
	AdminPrimaryColor        string  `json:"adminPrimaryColor"`
	AdminSecondaryColor      string  `json:"adminSecondaryColor"`
	AdminTextColor           string  `json:"adminTextColor"`
	AdminErrorColor          string  `json:"adminErrorColor"`
	AdminTitle               string  `json:"adminTitle"`
	AdminLogoFile            string  `json:"adminLogoHorizontal"`
	AdminLogoBundle          string  `json:"adminLogoVertical"`
	BundlesInstalled         Bundles `json:"bundlesInstalled"`
	Host                     string  `json:"host"`
	Protocol                 string  `json:"protocol"`
}

//installBundle adds the name of the new bundle to the slice of bundles
func (b Bundles) InstallBundle(newBundle string) (Bundles, error) {

	//Check if the bundle is already installed (should only happen if user has messed with config.json)
	//If the name of the bundle being installed coincides with any of the names already in the bundle slice,
	//then just return the original bundle slice
	for _, a := range b {
		if a == newBundle {
			return b, errors.New("Bundle is already installed")
		}
	}
	//Otherwise append
	return append(b, newBundle), nil
}

func (b Bundles) UnInstallBundle(bundle string) (Bundles, error) {

	//Search for the bundle to be uninstalled
	for index, a := range b {
		if a == bundle {
			//If found, splice it out
			return append(b[:index], b[index+1:]...), nil
		}
	}

	//Otherwise just return the original
	return b, errors.New("Bundle is not installed")

}

func compareBundles(b1, b2 Bundles) bool {
	//If lengths are not equal
	if len(b1) != len(b2) {
		return false
	}

	//If any of the elements are not the same
	for k := range b1 {
		if b1[k] != b2[k] {
			return false
		}
	}

	return true
}
