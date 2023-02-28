package service

import (
	"CompanySystemsMonitoring/internal/domain"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const (
	CreateCustomerMask int64 = 1 << iota
	PurchaseMask
	PayoutMask
	RecurringMask
	FraudControlMask
	CheckoutPageMask
)

type BillingService struct {
}

func NewBillingService() *BillingService {
	return &BillingService{}
}

// BillingRead read billing data
func (b BillingService) BillingRead(ctx context.Context, path string) domain.BillingData {
	billingDataResult := domain.BillingData{}
	select {
	case <-ctx.Done():
		log.Printf("cansel: BillingRead")
	default:
		file, err := os.Open(path)
		if err != nil {
			log.Println("Cannot open billingData file:", err)
		}
		defer file.Close()
		reader, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal("Cannot read billingData file:", err)
		}
		mask, err := strconv.ParseInt(string(reader), 2, 0)
		bCreateCustomer := mask&CreateCustomerMask != 0
		bPurchase := mask&PurchaseMask != 0
		bPayout := mask&PayoutMask != 0
		bRecurring := mask&RecurringMask != 0
		bFraudControl := mask&FraudControlMask != 0
		bCheckoutPage := mask&CheckoutPageMask != 0
		billingDataResult = domain.BillingData{
			CreateCustomer: bCreateCustomer,
			Purchase:       bPurchase,
			Payout:         bPayout,
			Recurring:      bRecurring,
			FraudControl:   bFraudControl,
			CheckoutPage:   bCheckoutPage,
		}
	}
	return billingDataResult
}
