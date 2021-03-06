// Copyright © 2019 VMware, INC
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

package rm

import (
	"errors"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/spf13/cobra"
)

// NewCommand returns the rm command of type cobra.Command
var byAge bool

//TODO need revisit ! What are the supported flags ? Are they consistent with the other commands ?
func removeNotificationHandler(cmd *cobra.Command, args []string) (err error){
	// Checking for args
	if len(args) == 0 {
		return errors.New("no notification ID/Name provided")
	}
	url := config.Conf.Clients["Notification"].Url() + clients.ApiNotificationRoute
	if byAge {
		url += "/age/" + args[0]
	} else {
		url += "/slug/" + args[0]
	}
	return request.DeletePrt(url, args[0])
}

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm [notification slug or age]",
		Short: "Removes notification by slug or age",
		Long:  `Removes a notification given its slug or age timestamp.`,
		Args:  cobra.ExactValidArgs(1),
		RunE:   removeNotificationHandler,
	}
	cmd.Flags().BoolVar(&byAge, "age", false, "Remove by age")
	return cmd
}
