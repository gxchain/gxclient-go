package types

type TransactionResult struct {
	SignedTransaction *SignedTransaction `json:"signed_transaction"`
	BroadcastResponse *BroadcastResponse `json:"broadcast_response"`
}
