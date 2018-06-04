package messaging

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/magiconair/properties/assert"
	"errors"
)

func TestMessageDispatcher_SendLeaderTransactions_TransactionsEmpty(t *testing.T) {
	// Given
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {
		return nil
	})
	md := MessageDispatcher{
		publisher: publisher,
	}
	transactions := []*txpool.Transaction{}
	leader := txpool.Leader{
		LeaderId: txpool.LeaderId{
			Id: "888",
		},
	}

	// When
	err := md.SendLeaderTransactions(transactions, leader)

	// Then
	assert.Equal(t, errors.New("Empty transaction list proposed"), err)

}

func TestMessageDispatcher_SendLeaderTransactions(t *testing.T) {
	// Given
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {
		return nil
	})
	md := MessageDispatcher{
		publisher: publisher,
	}
	transactions := []*txpool.Transaction{}
	leader := txpool.Leader{
		LeaderId: txpool.LeaderId{
			Id: "888",
		},
	}

	// When
	t1 := txpool.Transaction{
		TxId: "111",
	}
	transactions = append(transactions, &t1)
	err := md.SendLeaderTransactions(transactions, leader)

	// Then
	assert.Equal(t, nil, err);
}

func TestMessageDispatcher_ProposeBlock_TransactionsEmpty(t *testing.T) {
	// Given
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {
		return nil
	})
	md := MessageDispatcher{
		publisher: publisher,
	}
	var transactions = []txpool.Transaction{}

	// When
	err := md.ProposeBlock(transactions)

	// Then
	assert.Equal(t, errors.New("Empty transaction list proposed"), err)
}

func TestMessageDispatcher_ProposeBlock(t *testing.T) {
	// Given
	publisher := Publisher(func(exchange string, topic string, data interface{}) error {
		return nil
	})
	md := MessageDispatcher{
		publisher: publisher,
	}
	var transactions = []txpool.Transaction{}

	// When
	t1 := txpool.Transaction{
		TxId: "888",
	}
	transactions = append(transactions, t1)
	err := md.ProposeBlock(transactions)

	// Then
	assert.Equal(t, nil, err)
}