package main

import (
	"fmt"
	"github.com/OniGbemiga/block-constructor/internals"
)

func main() {
	filePath := "internals/mempool.csv"
	transactions, err := internals.ParseCSV(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	selectedTxids, selectedFee := internals.SelectTransactions(transactions)

	// Print the selected transactions in the required format
	for _, txid := range selectedTxids {
		fmt.Println(txid)
	}

	fmt.Printf("Total Fee: %d\n", selectedFee)
}
