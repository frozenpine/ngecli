// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"fmt"

	"github.com/frozenpine/ngecli/channels"

	"github.com/antihax/optional"
	"github.com/frozenpine/ngecli/models"
	"github.com/frozenpine/ngerest"

	"github.com/spf13/cobra"
)

type orderGetArgs struct {
	filter  string
	columns string
	start   models.FlagTime
	end     models.FlagTime
	count   int
	reverse bool
}

var orderGetVariables orderGetArgs

func getOrderOpts(symbol string, args *orderGetArgs) *ngerest.OrderGetOrdersOpts {
	options := ngerest.OrderGetOrdersOpts{}

	if symbol != "" {
		options.Symbol = optional.NewString(symbol)
	}

	if args.filter != "" {
		options.Filter = optional.NewString(args.filter)
	}

	if args.columns != "" {
		options.Columns = optional.NewString(args.columns)
	}

	if args.reverse {
		options.Reverse = optional.NewBool(args.reverse)
	}

	if args.start != models.EmptyTime {
		options.StartTime = optional.NewTime(args.start.GetTime())
	}

	if args.end != models.EmptyTime {
		options.EndTime = optional.NewTime(args.end.GetTime())
	}

	if args.count > 0 {
		options.Count = optional.NewFloat32(float32(args.count))
	}

	return &options
}

// orderGetCmd represents the orderGet command
var orderGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get user's history orders.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("orderGet called")

		hisOrders, _, err := client.Order.OrderGetOrders(rootCtx, getOrderOpts(symbol, &orderGetVariables))

		if err != nil {
			if swErr, ok := err.(ngerest.GenericSwaggerError); ok {
				channels.ErrChan <- fmt.Errorf("Get order failed: %s\n%s", swErr.Error(), string(swErr.Body()))
			} else {
				channels.ErrChan <- fmt.Errorf("Get order failed: %s", err.Error())
			}
		}
	},
}

func init() {
	orderCmd.AddCommand(orderGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// orderGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// orderGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	orderGetCmd.Flags().StringVar(
		&orderGetVariables.filter, "filter", "", "Filter string applied in query result")
	orderGetCmd.Flags().StringVar(
		&orderGetVariables.columns, "columns", "", "Column names for query result.")

	orderGetCmd.Flags().BoolVarP(
		&orderGetVariables.reverse, "reverse", "r", false, "Getting query results in reversed order.")

	orderGetCmd.Flags().VarP(&orderGetVariables.start, "start", "s", "Start")
	orderGetCmd.Flags().VarP(&orderGetVariables.end, "end", "e", "End")

	orderGetCmd.Flags().IntVarP(&orderGetVariables.count, "count", "c", 200, "Order count in query result.")
}
