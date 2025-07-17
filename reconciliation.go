package pelucio

import (
	"crypto/sha256"
	"fmt"
	"time"
)

const (
	Missing    ReconciliationEntryStatus = "missing"
	ExactMatch ReconciliationEntryStatus = "exact_match"
)

type (
	ReconciliationEntryStatus string

	Reconciliation struct {
		DateTolerance   time.Duration
		InternalEntries []*ReconciliationEntry
		ExternalEntries []*ReconciliationEntry
	}

	ReconciliationEntry struct {
		Entry
		ExternalID             string
		RelatedEntryExternalID *string
		Status                 ReconciliationEntryStatus

		CreatedAt time.Time
	}
)

func (p ReconciliationEntry) Hash(dateInterval time.Duration) (string, error) {
	normalizedCreatedAt := p.CreatedAt.Truncate(dateInterval)
	amountStr := p.Amount.String()
	toHash := fmt.Sprintf("%s-%s-%s-%s-%s", p.ExternalID, normalizedCreatedAt, amountStr, p.EntrySide, p.Currency)
	hash := sha256.New()
	_, err := hash.Write([]byte(toHash))
	if err != nil {
		return "", err
	}

	hexHash := fmt.Sprintf("%x", hash.Sum(nil))

	return hexHash, nil
}

func (p *Reconciliation) Reconcile() error {
	internalEntries, err := p.InternalHashes()
	if err != nil {
		return err
	}

	externalEntries, err := p.ExternalHashes()
	if err != nil {
		return err
	}

	for k, internalEntry := range internalEntries {
		if externalEntry, ok := externalEntries[k]; ok {
			internalEntry.Status = ExactMatch
			internalEntry.RelatedEntryExternalID = &externalEntry.ExternalID
			delete(externalEntries, k)
		}
	}

	return nil
}

func (p *Reconciliation) InternalHashes() (map[string]*ReconciliationEntry, error) {
	return p.HashEntries(p.InternalEntries)
}

func (p *Reconciliation) ExternalHashes() (map[string]*ReconciliationEntry, error) {
	return p.HashEntries(p.ExternalEntries)
}

func (p *Reconciliation) HashEntries(toHash []*ReconciliationEntry) (map[string]*ReconciliationEntry, error) {
	hashes := make(map[string]*ReconciliationEntry, len(p.InternalEntries))
	for _, e := range p.InternalEntries {
		hash, err := e.Hash(p.DateTolerance)
		if err != nil {
			return nil, err
		}
		hashes[hash] = e
	}

	return hashes, nil
}
