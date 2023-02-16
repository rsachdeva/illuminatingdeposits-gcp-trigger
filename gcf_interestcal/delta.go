// Calculations

//nolint:revive,stylecheck // using underscore in package name for clarity
package gcf_interestcal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	errCal     = errors.New("error in calculation")
	errMinDays = errors.New("NewDeposit period in years should not be less than 30 days")
)

const hundred = 100

type DeltaService struct{}

type NewDeposit struct {
	Account     string  `json:"account,omitempty"`
	AccountType string  `json:"account_type" validate:"required"` //nolint:tagliatelle
	Apy         float64 `json:"apy" validate:"gte=0"`
	Years       float64 `json:"years" validate:"required"`
	Amount      float64 `json:"amount" validate:"gte=0"`
}

type NewBank struct {
	Name        string        `json:"name" validate:"required"`
	NewDeposits []*NewDeposit `json:"new_deposits" validate:"required,dive"` //nolint:tagliatelle
}

type CreateInterestRequest struct {
	NewBanks []*NewBank `json:"new_banks" validate:"required,dive"` //nolint:tagliatelle
}

type CreateInterestResponse struct {
	Banks []*Bank `json:"banks,omitempty"`
	Delta float64 `json:"delta,omitempty"`
}

type Deposit struct {
	Account     string  `json:"account,omitempty"`
	AccountType string  `json:"account_type,omitempty"` //nolint:tagliatelle
	Apy         float64 `json:"apy,omitempty"`
	Years       float64 `json:"years,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Delta       float64 `json:"delta,omitempty"`
}

type Bank struct {
	Name     string     `json:"name,omitempty"`
	Deposits []*Deposit `json:"deposits,omitempty"`
	Delta    float64    `json:"delta,omitempty"`
}

const (
	// Saving for saving type.
	Saving = "Saving"
	// CertDep for cd type.
	CertDep = "CD"
	// Checking gor checking type.
	Checking = "Checking"
	// BrokerCD for Brokered cd type.
	BrokerCD = "Brokered CD"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

//nolint:gochecknoglobals // use a single instance of Validate, it caches struct info.
var validate *validator.Validate

func calculateDelta(ctx context.Context, reqData []byte) (*CreateInterestResponse, error) {
	var (
		deltaSvc DeltaService
		req      CreateInterestRequest
	)

	if err := json.Unmarshal(reqData, &req); err != nil {
		return nil, fmt.Errorf("decode request payload %w", err)
	}

	printDecodedReq(req)

	validate = validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("invalid request payload %w", err)
	}

	log.Printf("interestcalpb.NewBank nb is %v and type is %T", req, req)

	resp, err := deltaSvc.compute(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("compute err %w", err)
	}

	return resp, nil
}

// print nested struct.
func printDecodedReq(req CreateInterestRequest) {
	for _, nb := range req.NewBanks {
		log.Printf("new bank name is %#v", nb.Name)

		for _, nd := range nb.NewDeposits {
			log.Printf("for this bank new deposit is %#v", nd)
		}
	}
}

// CalculateDelta calculations for all banks
// computeBanksDelta uses concurrency more for example.
func (svc DeltaService) compute(ctx context.Context, cireq *CreateInterestRequest) (*CreateInterestResponse, error) {
	defer timeTrack(time.Now(), "Delta timed with withConcurrency for more I/O processing")

	// withConcurrency := true
	// if withConcurrency {
	// 	bks, delta, err = computeBanksDelta(cireq)
	// }
	// if !withConcurrency {
	// 	bks, delta, err = computeBanksDeltaSequentially(cireq)
	// }
	log.Println("Starting overall calculate delta")

	bks, delta, err := svc.computeBanksDelta(ctx, cireq)
	if err != nil {
		return nil, fmt.Errorf("error in computeBanksDelta %w", err)
	}

	ciresp := CreateInterestResponse{
		Banks: bks,
		Delta: roundToNearest(delta),
	}

	return &ciresp, nil
}

// func computeBanksDeltaSequentially(cireq *CreateInterestRequest) ([]*Bank, float64, error) {
// 	var bks []*Bank
// 	var delta float64
// 	for _, nb := range cireq.NewBanks {
// 		ds, bDelta, err := computeBankDeltaNoChannel(nb)
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 		bk := Bank{
// 			Name:     nb.Name,
// 			Deposits: ds,
// 			Delta:    roundToNearest(bDelta),
// 		}
// 		bks = append(bks, &bk)
// 		delta = delta + bk.Delta
// 	}
// 	return bks, delta, nil
// }
//
// func computeBankDeltaNoChannel(nb *NewBank) ([]*Deposit, float64, error) {
// 	// time.Sleep(5 * time.Second)
// 	var ds []*Deposit
// 	var bDelta float64
// 	for _, nd := range nb.NewDeposits {
// 		d := Deposit{
// 			Account:     nd.Account,
// 			AccountType: nd.AccountType,
// 			Apy:         nd.Apy,
// 			Years:       nd.Years,
// 			Amount:      nd.Amount,
// 		}
// 		err := computeDepositDelta(&d)
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 		ds = append(ds, &d)
// 		bDelta = bDelta + d.Delta
// 	}
// 	return ds, bDelta, nil
// }

type BankResult struct {
	name   string
	ds     []*Deposit
	bDelta float64
	err    error
}

func (svc DeltaService) computeBanksDelta(ctx context.Context, cireq *CreateInterestRequest) ([]*Bank, float64, error) {
	bks := make([]*Bank, 0, len(cireq.NewBanks))

	var (
		delta     float64
		waitGroup sync.WaitGroup
	)

	bkCh := make(chan BankResult)
	// to know when to close channel

	// for all banks just 1 go routine
	// go func() {
	// 	for _, nb := range cireq.NewBanks {
	// 		// ds, bDelta, err := computeBankDelta(nb)
	// 		// doing for each bank calculation concurrently
	// 		computeBankDelta(nb, bkCh)
	// 	}
	// 	close(bkCh)
	// }()

	// for each bank calculation with all deposits 1 goroutine
	log.Println("len(cireq.NewBanks) is", len(cireq.NewBanks))

	waitGroup.Add(len(cireq.NewBanks))

	for _, nb := range cireq.NewBanks {
		go func(nb *NewBank) {
			defer waitGroup.Done()
			svc.computeBankDelta(ctx, nb, bkCh)
		}(nb)
	}

	go func() {
		waitGroup.Wait()
		close(bkCh)
	}()

	for bkRes := range bkCh {
		if bkRes.err != nil {
			return nil, 0, bkRes.err
		}

		bank := Bank{
			Name:     bkRes.name,
			Deposits: bkRes.ds,
			Delta:    roundToNearest(bkRes.bDelta),
		}
		bks = append(bks, &bank)
		delta += bank.Delta
	}

	return bks, delta, nil
}

// return []*Deposit, float64, error now in channel.
func (svc DeltaService) computeBankDelta(_ context.Context, newBank *NewBank, bkCh chan<- BankResult) {
	// https://stackoverflow.com/questions/59734706/how-to-resolve-consider-preallocating-prealloc-lint
	deposits := make([]*Deposit, 0, len(newBank.NewDeposits))

	var bDelta float64

	for _, newDep := range newBank.NewDeposits {
		log.Printf("for newBank.name %v the newDep is %v", newBank.Name, newDep)

		dep := Deposit{
			Account:     newDep.Account,
			AccountType: newDep.AccountType,
			Apy:         newDep.Apy,
			Years:       newDep.Years,
			Amount:      newDep.Amount,
			Delta:       0,
		}
		err := computeDepositDelta(&dep)
		// if err != nil {
		// 	return nil, 0, err
		// }
		// Sending err result
		log.Println("any err in compute deposit delta", err)

		if err != nil {
			log.Printf("error in goroutine running computeDepositDelta %v\n", err)
			bkCh <- BankResult{
				name:   "",
				ds:     nil,
				bDelta: 0,
				err:    err,
			}

			return
		}

		deposits = append(deposits, &dep)
		log.Println("inside loop deposits is", deposits)

		bDelta += dep.Delta
		log.Println("inside loop bDelta is", bDelta)
	}

	log.Println("outside loop deposits is", deposits)
	log.Println("outside loop bDelta is", bDelta)
	// return deposits, bDelta, nil
	bkRes := BankResult{
		name:   newBank.Name,
		ds:     deposits,
		bDelta: bDelta,
		err:    nil,
	}
	bkCh <- bkRes
}

// Delta calculates interest for 30 days for output/response Deposit.
func computeDepositDelta(dep *Deposit) error {
	e := earned(dep)
	e30Days, err := earned30days(e, dep.Years)
	log.Println("err in computeDepositDelta is", err)

	if err != nil {
		log.Println("error in computeDepositDelta ", err)

		return fmt.Errorf("%w for accpunt %s", errCal, dep.Account)
	}

	dep.Delta = roundToNearest(e30Days)
	log.Println("no error returned computeDepositDelta dep.Delta is", dep.Delta)

	return nil
}

func earned(dep *Deposit) float64 {
	log.Println("dep.AccountType is", dep.AccountType)

	switch dep.AccountType {
	case Saving, CertDep:
		return compoundInterest(dep.Apy, dep.Years, dep.Amount)
	case BrokerCD:
		return simpleInterest(dep.Apy, dep.Years, dep.Amount)
	case Checking:
		return 0.0
	default:
		return -1.0
	}
}

func roundToNearest(n float64) float64 {
	return math.Round(n*hundred) / hundred
}

func simpleInterest(apy float64, years float64, amount float64) float64 {
	rateInDecimal := apy / hundred
	intEarned := amount * rateInDecimal * years

	return intEarned
}

func compoundInterest(apy float64, years float64, amount float64) float64 {
	rateInDecimal := apy / hundred
	calInProcess := math.Pow(1+rateInDecimal, years)
	intEarned := amount*calInProcess - amount

	return intEarned
}

func earned30days(iEarned float64, years float64) (float64, error) {
	const (
		monthDays = 30
		yearDays  = 365
	)

	if years*yearDays < monthDays {
		return 0, fmt.Errorf("%w for years %v", errMinDays, years)
	}

	i1Day := iEarned / (years * yearDays)
	i30 := i1Day * monthDays

	return math.Round(i30*hundred) / hundred, nil
}
