//
// Copyright 2019 Chef Software, Inc.
// Author: Salim Afiune <afiune@chef.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package credentials

//
// The intend of this file is to have a single place where we can easily
// visualize the list of all error messages that we present to users.
//

const (
	ProfileNotFoundErr = `
  profile not found in credentials file.

  verify the format of the credentials file by following this documentation:
    - https://docs.chef.io/knife_setup.html#knife-profiles
`
	CredentialsNotFoundErr = `
  credentials file not found. (default: $HOME/.chef/credentials)

  setup your local credentials config by following this documentation:
    - https://docs.chef.io/knife_setup.html#knife-profiles
`
	MalformedCredentialsFileErr = `
  unable to parse credentials file.

  verify the format of the credentials file by following this documentation:
    - https://docs.chef.io/knife_setup.html#knife-profiles
`
)
