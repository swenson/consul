//go:build !consulent
// +build !consulent

package state

import (
	"fmt"
	"strings"

	memdb "github.com/hashicorp/go-memdb"

	"github.com/hashicorp/consul/acl"
	"github.com/hashicorp/consul/agent/configentry"
	"github.com/hashicorp/consul/agent/structs"
)

func indexFromConfigEntryKindName(arg interface{}) ([]byte, error) {
	var b indexBuilder

	switch n := arg.(type) {
	case *acl.EnterpriseMeta:
		return nil, nil
	case acl.EnterpriseMeta:
		return b.Bytes(), nil
	case ConfigEntryKindQuery:
		b.String(strings.ToLower(n.Kind))
		return b.Bytes(), nil
	case configentry.KindName:
		b.String(strings.ToLower(n.Kind))
		b.String(strings.ToLower(n.Name))
		return b.Bytes(), nil
	}

	return nil, fmt.Errorf("invalid type for ConfigEntryKindName query: %T", arg)
}

func validateConfigEntryEnterprise(_ ReadTxn, _ structs.ConfigEntry) error {
	return nil
}

func getAllConfigEntriesWithTxn(tx ReadTxn, _ *acl.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get(tableConfigEntries, indexID)
}

func getAllConfigEntriesByKindWithTxn(tx ReadTxn, kind string) (memdb.ResultIterator, error) {
	return getConfigEntryKindsWithTxn(tx, kind, nil)
}

func getConfigEntryKindsWithTxn(tx ReadTxn, kind string, _ *acl.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get(tableConfigEntries, indexID+"_prefix", ConfigEntryKindQuery{Kind: kind})
}

func configIntentionsConvertToList(iter memdb.ResultIterator, _ *acl.EnterpriseMeta) structs.Intentions {
	var results structs.Intentions
	for v := iter.Next(); v != nil; v = iter.Next() {
		entry := v.(*structs.ServiceIntentionsConfigEntry)
		for _, src := range entry.Sources {
			results = append(results, entry.ToIntention(src))
		}
	}
	return results
}
