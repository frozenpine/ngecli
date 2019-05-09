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
	"os"

	"github.com/frozenpine/ngecli/models"

	"github.com/spf13/cobra"
)

const (
	defaultPrice     = float64(5050)
	defaultTick      = float64(0.01)
	defaultVolume    = int64(1)
	defaultMaxVolume = int64(10)
	defaultCount     = 1
)

type orderNewArgs struct {
	price  float64
	volume int64
	side   models.OrderSide

	basePrice  float64
	priceTick  float64
	baseVolume int64
	maxVolume  int64
	random     bool
	bothSide   bool
	count      int
}

var orderNewVariables orderNewArgs

func checkArgs(vars *orderNewArgs) bool {
	if err := CheckSymbol(symbol); err != nil {
		fmt.Println(err)
		return false
	}

	if err := CheckPrice(vars.price); err != nil {
		fmt.Println(err)
		return false
	}

	if err := CheckQuantity(vars.volume); err != nil {
		fmt.Println(err)
		return false
	}

	if err := MatchSide(&vars.side, vars.volume); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

// orderNewCmd represents the orderGet command
var orderNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Make new order for user.",
	Long:  `Make new orders either by args inputed or a order source file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("orderNew called")

		if !checkArgs(&orderNewVariables) {
			fmt.Println("variables check failed.")
			os.Exit(1)
		}
	},
}

func init() {
	orderCmd.AddCommand(orderNewCmd)

	orderNewCmd.Flags().Float64VarP(&orderNewVariables.price, "price", "p", 0, "Price for new order.")
	orderNewCmd.Flags().Int64VarP(&orderNewVariables.volume, "volume", "v", 0, "Volume for new order.")
	orderNewCmd.Flags().Var(&orderNewVariables.side, "side", "Side for new order.")

	orderNewCmd.Flags().Float64Var(
		&orderNewVariables.basePrice, "base-price", defaultPrice, "Base price for random order.")
	orderNewCmd.Flags().Float64Var(
		&orderNewVariables.priceTick, "tick", defaultTick, "Price tick for new order.")

	orderNewCmd.Flags().Int64Var(
		&orderNewVariables.baseVolume, "base-volume", defaultVolume, "Base volume for random order.")
	orderNewCmd.Flags().Int64Var(
		&orderNewVariables.maxVolume, "max-volume", defaultMaxVolume, "Max volume for random order.")

	orderNewCmd.Flags().BoolVar(&orderNewVariables.random, "random", false, "Random price/volume if not specified.")
	orderNewCmd.Flags().BoolVar(
		&orderNewVariables.bothSide, "both-side", false, "Make new orders in both side on same volume@price.")
	orderNewCmd.Flags().IntVarP(&orderNewVariables.count, "count", "c", defaultCount, "Count of new orders.")
}
