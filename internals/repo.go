package internals

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Transaction struct {
	Txid     string
	Fee      int
	Weight   int
	Parents  []string
	Selected bool
}

func ParseCSV(filePath string) ([]Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var transactions []Transaction
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		// Trim any leading or trailing whitespace in the fields
		for i := range row {
			row[i] = strings.TrimSpace(row[i])
		}

		// Check if the row has at least three elements (txid, fee, weight)
		if len(row) < 3 {
			//fmt.Printf("Skipping invalid row: %v\n", row)
			continue
		}

		txid := row[0]
		fee, _ := strconv.Atoi(row[1])
		weight, _ := strconv.Atoi(row[2])

		var parentTxids []string
		if len(row) > 3 && row[3] != "" {
			parentTxids = strings.Split(row[3], ";")
		}

		transaction := Transaction{
			Txid:    txid,
			Fee:     fee,
			Weight:  weight,
			Parents: parentTxids,
		}

		fmt.Println("-----txn", transactions)

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func SelectTransactions(transactions []Transaction) (selectedTxids []string, selectedFee int) {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Fee > transactions[j].Fee
	})

	totalWeight := 0
	selectedTxidSet := make(map[string]bool)

	isValid := func(transaction Transaction) bool {
		for _, parent := range transaction.Parents {
			if !selectedTxidSet[parent] {
				return false
			}
		}
		return true
	}

	for _, transaction := range transactions {
		if !transaction.Selected && isValid(transaction) && totalWeight+transaction.Weight <= 4000000 {
			selectedTxids = append(selectedTxids, transaction.Txid)
			selectedFee += transaction.Fee
			totalWeight += transaction.Weight
			transaction.Selected = true
			selectedTxidSet[transaction.Txid] = true
		}
	}

	return selectedTxids, selectedFee
}
