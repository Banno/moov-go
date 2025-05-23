package moov_test

import (
	"testing"

	"github.com/moovfinancial/moov-go/internal/testtools"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_Receipts(t *testing.T) {
	mc := NewTestClient(t)

	customer := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())
	customerCard := createTemporaryCard(t, mc, customer.AccountID)

	var transfer *moov.TransferStarted = nil
	var receipt *moov.Receipt = nil

	t.Run("make async transfer", func(t *testing.T) {
		started, err := mc.CreateTransfer(BgCtx(), testtools.PARTNER_ID, moov.CreateTransfer{
			Source: moov.CreateTransfer_Source{
				PaymentMethodID: customerCard.PaymentMethods[0].PaymentMethodID,
			},
			Destination: moov.CreateTransfer_Destination{
				PaymentMethodID: testtools.MERCHANT_WALLET_PM_ID,
			},
			Amount: moov.Amount{
				Currency: "usd",
				Value:    1,
			},
		}).Started()
		NoResponseError(t, err)

		// We made an async transfer, so completed should be nil, while started not nil
		require.NotNil(t, started)

		transfer = started
	})

	t.Run("create receipt", func(t *testing.T) {
		require.NotNil(t, transfer)

		receipts, err := mc.CreateReceipt(BgCtx(), moov.CreateReceipt{
			Kind:          "sale.customer.v1",
			ForTransferID: &transfer.TransferID,
			Email:         moov.PtrOf("noreply@moov.io"),
		})
		require.NoError(t, err)
		require.Len(t, receipts, 1)
		require.NotEmpty(t, receipts[0].ID)

		receipt = &receipts[0]
		t.Logf("receipt: %+v\n", receipt)
	})

	t.Run("list receipts", func(t *testing.T) {
		require.NotNil(t, transfer)

		receipts, err := mc.ListReceipts(BgCtx(), moov.ReceiptByTransferID(transfer.TransferID))
		require.NoError(t, err)
		require.Len(t, receipts, 1)

		// The receipt may be sent or not by this point so we need to clear out the sentFor list
		for i := range receipts {
			receipts[i].SentFor = nil
		}

		require.Contains(t, receipts, *receipt)
	})
}
