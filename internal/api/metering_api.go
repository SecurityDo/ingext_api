package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/ingext_api/api"
	"github.com/SecurityDo/ingext_api/model"
)

// DailyUsage returns the daily billing ledger for [from, to] inclusive.
//
// Elapsed days with no row come back as State "missing" rather than omitted,
// and StoreAvailable false means the ledger was unreachable -- neither is a day
// of zero usage.
func (c *Client) DailyUsage(from, to string, opts *fluencyAPI.UsageOptions) (*model.UsageResponse, error) {
	svc := fluencyAPI.NewMeteringService(c.ingextClient)
	resp, err := svc.DailyUsage(from, to, opts)
	if err != nil {
		c.Logger.Error("failed to read daily usage", "from", from, "to", to, "error", err)
		return nil, fmt.Errorf("failed to read daily usage %s..%s: %w", from, to, err)
	}
	return resp, nil
}

// UsageAttempts returns the collection attempt ledger for [from, to].
func (c *Client) UsageAttempts(from, to string, limit int, opts *fluencyAPI.UsageOptions) (*model.AttemptsResponse, error) {
	svc := fluencyAPI.NewMeteringService(c.ingextClient)
	resp, err := svc.Attempts(from, to, limit, opts)
	if err != nil {
		c.Logger.Error("failed to read usage attempts", "from", from, "to", to, "error", err)
		return nil, fmt.Errorf("failed to read usage attempts %s..%s: %w", from, to, err)
	}
	return resp, nil
}

// CollectUsageDay collects and closes one past UTC day on demand.
func (c *Client) CollectUsageDay(date string, opts *fluencyAPI.UsageOptions) (*model.CollectDayResponse, error) {
	svc := fluencyAPI.NewMeteringService(c.ingextClient)
	resp, err := svc.CollectDay(date, opts)
	if err != nil {
		c.Logger.Error("failed to collect usage day", "date", date, "error", err)
		return nil, fmt.Errorf("failed to collect usage day %s: %w", date, err)
	}
	return resp, nil
}

// UsageSinkClassification reports how the collector sorts this account's
// datasinks into meters.
func (c *Client) UsageSinkClassification() (*model.SinkClassification, error) {
	svc := fluencyAPI.NewMeteringService(c.ingextClient)
	resp, err := svc.SinkClassification()
	if err != nil {
		c.Logger.Error("failed to read sink classification", "error", err)
		return nil, fmt.Errorf("failed to read sink classification: %w", err)
	}
	return resp, nil
}

// GridUsage collects the daily ledger for every tenant of a provider in one
// call. Must be issued against a provider site.
func (c *Client) GridUsage(from, to string, accounts []string, opts *fluencyAPI.UsageOptions) (*model.GridUsageResponse, error) {
	svc := fluencyAPI.NewMeteringService(c.ingextClient)
	resp, err := svc.GridUsage(from, to, accounts, opts)
	if err != nil {
		c.Logger.Error("failed to read grid usage", "from", from, "to", to, "error", err)
		return nil, fmt.Errorf("failed to read grid usage %s..%s: %w", from, to, err)
	}
	return resp, nil
}
