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
	"encoding/json"
	"fmt"
	"sync"

	"github.com/frozenpine/ngecli/logger"

	"go.uber.org/zap"

	"github.com/frozenpine/ngecli/common"

	"github.com/antihax/optional"
	"github.com/frozenpine/ngecli/models"
	"github.com/frozenpine/ngerest"

	"github.com/spf13/cobra"
)

const defaultGetOrderCount = 200

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

func printOrderResults(wait *sync.WaitGroup, results <-chan *models.Order) {
	var count int

	for ord := range results {
		jsonBytes, err := json.Marshal(ord)
		if err != nil {
			logger.Warn(err.Error())
		} else {
			fmt.Println(string(jsonBytes))
			count++
		}
	}

	logger.Info("All order results printed.", zap.Int("count", count))
	wait.Done()
}

// orderGetCmd represents the orderGet command
var orderGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get user's history orders.",
	Long:  `Get user's history orders.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !auths.HasSavedAuth(common.GetBaseHost()) &&
			!auths.HasDefaultAuth() {
			identity, password := CollectLoginInfo()

			auths.DefaultID, auths.DefaultPass = identity, *password
		}

		client, err := clientHub.GetClient(common.GetBaseHost())
		if err != nil {
			logger.Warn(err.Error())
			return
		}

		hisOrders, _, err := client.Order.OrderGetOrders(
			auths.NextAuth(nil), getOrderOpts(symbol, &orderGetVariables))

		if err != nil {
			common.PrintError("Get order failed", err)
			return
		} else if len(hisOrders) < 1 {
			logger.Warn("No history orders found.")
			return
		}

		waitOutput := sync.WaitGroup{}
		waitOutput.Add(1)

		go printOrderResults(&waitOutput, orderCache.GetResults())

		for _, order := range hisOrders {
			orderCache.PutResult(&order)
		}

		orderCache.CloseResults()

		waitOutput.Wait()
	},
}

func init() {
	orderCmd.AddCommand(orderGetCmd)

	orderGetCmd.Flags().StringVar(
		&orderGetVariables.filter, "filter", "", "Filter string applied in query result")
	orderGetCmd.Flags().StringVar(
		&orderGetVariables.columns, "columns", "", "Column names for query result.")

	orderGetCmd.Flags().BoolVarP(
		&orderGetVariables.reverse, "reverse", "r", false, "Getting query results in reversed order.")

	orderGetCmd.Flags().VarP(&orderGetVariables.start, "start", "s", "Start")
	orderGetCmd.Flags().VarP(&orderGetVariables.end, "end", "e", "End")

	orderGetCmd.Flags().IntVarP(&orderGetVariables.count, "count", "c", defaultGetOrderCount, "Order count in query result.")
}
