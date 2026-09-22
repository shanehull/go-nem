package nem

import (
	"fmt"
	"strconv"
	"time"
)

type columnIndex map[string]int

func indexColumns(columns []string) columnIndex {
	idx := make(columnIndex, len(columns))
	for i, name := range columns {
		idx[name] = i
	}
	return idx
}

func (idx columnIndex) get(row []string, name string) string {
	i, ok := idx[name]
	if !ok || i >= len(row) {
		return ""
	}
	return row[i]
}

func parseInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

func decode[T any](tables []Table, group, name string, rowFn func(columnIndex, []string) (T, error)) ([]T, error) {
	var out []T
	found := false
	for _, t := range tables {
		if t.Group != group || t.Name != name {
			continue
		}
		found = true
		idx := indexColumns(t.Columns)
		for i, r := range t.Rows {
			v, err := rowFn(idx, r)
			if err != nil {
				return nil, fmt.Errorf("nem: decode %s/%s row %d: %w", group, name, i, err)
			}
			out = append(out, v)
		}
	}
	if !found {
		return nil, fmt.Errorf("%w: %s/%s", ErrTableNotFound, group, name)
	}
	return out, nil
}

func decodeTime(idx columnIndex, row []string, column string) (time.Time, error) {
	t, err := ParseNEMTime(idx.get(row, column))
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// DispatchPrice is a row of the DISPATCH/PRICE table, the regional energy price.
type DispatchPrice struct {
	SettlementDate time.Time
	RunNo          int
	RegionID       string
	RRP            float64
}

// DecodeDispatchPrice decodes the DISPATCH/PRICE table.
func DecodeDispatchPrice(tables []Table) ([]DispatchPrice, error) {
	return decode(tables, "DISPATCH", "PRICE", func(idx columnIndex, row []string) (DispatchPrice, error) {
		sd, err := decodeTime(idx, row, "SETTLEMENTDATE")
		if err != nil {
			return DispatchPrice{}, err
		}
		runNo, err := parseInt(idx.get(row, "RUNNO"))
		if err != nil {
			return DispatchPrice{}, fmt.Errorf("runno: %w", err)
		}
		rrp, err := parseFloat(idx.get(row, "RRP"))
		if err != nil {
			return DispatchPrice{}, fmt.Errorf("rrp: %w", err)
		}
		return DispatchPrice{SettlementDate: sd, RunNo: runNo, RegionID: idx.get(row, "REGIONID"), RRP: rrp}, nil
	})
}

// DispatchRegionSum is a row of the DISPATCH/REGIONSUM table.
type DispatchRegionSum struct {
	SettlementDate      time.Time
	RunNo               int
	RegionID            string
	TotalDemand         float64
	AvailableGeneration float64
	NetInterchange      float64
	LORSurplus          float64
	LRCSurplus          float64
}

// DecodeDispatchRegionSum decodes the DISPATCH/REGIONSUM table.
func DecodeDispatchRegionSum(tables []Table) ([]DispatchRegionSum, error) {
	return decode(tables, "DISPATCH", "REGIONSUM", func(idx columnIndex, row []string) (DispatchRegionSum, error) {
		sd, err := decodeTime(idx, row, "SETTLEMENTDATE")
		if err != nil {
			return DispatchRegionSum{}, err
		}
		runNo, err := parseInt(idx.get(row, "RUNNO"))
		if err != nil {
			return DispatchRegionSum{}, fmt.Errorf("runno: %w", err)
		}
		totalDemand, err := parseFloat(idx.get(row, "TOTALDEMAND"))
		if err != nil {
			return DispatchRegionSum{}, fmt.Errorf("totaldemand: %w", err)
		}
		availableGeneration, err := parseFloat(idx.get(row, "AVAILABLEGENERATION"))
		if err != nil {
			return DispatchRegionSum{}, fmt.Errorf("availablegeneration: %w", err)
		}
		netInterchange, err := parseFloat(idx.get(row, "NETINTERCHANGE"))
		if err != nil {
			return DispatchRegionSum{}, fmt.Errorf("netinterchange: %w", err)
		}
		lorSurplus, err := parseFloat(idx.get(row, "LORSURPLUS"))
		if err != nil {
			return DispatchRegionSum{}, fmt.Errorf("lorsurplus: %w", err)
		}
		lrcSurplus, err := parseFloat(idx.get(row, "LRCSURPLUS"))
		if err != nil {
			return DispatchRegionSum{}, fmt.Errorf("lrcsurplus: %w", err)
		}
		return DispatchRegionSum{
			SettlementDate:      sd,
			RunNo:               runNo,
			RegionID:            idx.get(row, "REGIONID"),
			TotalDemand:         totalDemand,
			AvailableGeneration: availableGeneration,
			NetInterchange:      netInterchange,
			LORSurplus:          lorSurplus,
			LRCSurplus:          lrcSurplus,
		}, nil
	})
}

// DispatchInterconnectorRes is a row of the DISPATCH/INTERCONNECTORRES table.
type DispatchInterconnectorRes struct {
	SettlementDate   time.Time
	RunNo            int
	InterconnectorID string
	MeteredMWFlow    float64
	MWFlow           float64
	ExportLimit      float64
	ImportLimit      float64
}

// DecodeDispatchInterconnectorRes decodes the DISPATCH/INTERCONNECTORRES table.
func DecodeDispatchInterconnectorRes(tables []Table) ([]DispatchInterconnectorRes, error) {
	return decode(tables, "DISPATCH", "INTERCONNECTORRES", func(idx columnIndex, row []string) (DispatchInterconnectorRes, error) {
		sd, err := decodeTime(idx, row, "SETTLEMENTDATE")
		if err != nil {
			return DispatchInterconnectorRes{}, err
		}
		runNo, err := parseInt(idx.get(row, "RUNNO"))
		if err != nil {
			return DispatchInterconnectorRes{}, fmt.Errorf("runno: %w", err)
		}
		meteredMWFlow, err := parseFloat(idx.get(row, "METEREDMWFLOW"))
		if err != nil {
			return DispatchInterconnectorRes{}, fmt.Errorf("meteredmwflow: %w", err)
		}
		mwFlow, err := parseFloat(idx.get(row, "MWFLOW"))
		if err != nil {
			return DispatchInterconnectorRes{}, fmt.Errorf("mwflow: %w", err)
		}
		exportLimit, err := parseFloat(idx.get(row, "EXPORTLIMIT"))
		if err != nil {
			return DispatchInterconnectorRes{}, fmt.Errorf("exportlimit: %w", err)
		}
		importLimit, err := parseFloat(idx.get(row, "IMPORTLIMIT"))
		if err != nil {
			return DispatchInterconnectorRes{}, fmt.Errorf("importlimit: %w", err)
		}
		return DispatchInterconnectorRes{
			SettlementDate:   sd,
			RunNo:            runNo,
			InterconnectorID: idx.get(row, "INTERCONNECTORID"),
			MeteredMWFlow:    meteredMWFlow,
			MWFlow:           mwFlow,
			ExportLimit:      exportLimit,
			ImportLimit:      importLimit,
		}, nil
	})
}

// DispatchUnitScada is a row of the DISPATCH/UNIT_SCADA table.
type DispatchUnitScada struct {
	SettlementDate time.Time
	DUID           string
	SCADAValue     float64
}

// DecodeDispatchUnitScada decodes the DISPATCH/UNIT_SCADA table.
func DecodeDispatchUnitScada(tables []Table) ([]DispatchUnitScada, error) {
	return decode(tables, "DISPATCH", "UNIT_SCADA", func(idx columnIndex, row []string) (DispatchUnitScada, error) {
		sd, err := decodeTime(idx, row, "SETTLEMENTDATE")
		if err != nil {
			return DispatchUnitScada{}, err
		}
		value, err := parseFloat(idx.get(row, "SCADAVALUE"))
		if err != nil {
			return DispatchUnitScada{}, fmt.Errorf("scadavalue: %w", err)
		}
		return DispatchUnitScada{SettlementDate: sd, DUID: idx.get(row, "DUID"), SCADAValue: value}, nil
	})
}

// FCASPrice is a row of the DISPATCH/PRICE table, restricted to ancillary prices.
type FCASPrice struct {
	SettlementDate time.Time
	RunNo          int
	RegionID       string
	Raise6Sec      float64
	Raise60Sec     float64
	Raise5Min      float64
	RaiseReg       float64
	Lower6Sec      float64
	Lower60Sec     float64
	Lower5Min      float64
	LowerReg       float64
}

// DecodeFCASPrice decodes the ancillary service prices from the DISPATCH/PRICE table.
func DecodeFCASPrice(tables []Table) ([]FCASPrice, error) {
	return decode(tables, "DISPATCH", "PRICE", func(idx columnIndex, row []string) (FCASPrice, error) {
		sd, err := decodeTime(idx, row, "SETTLEMENTDATE")
		if err != nil {
			return FCASPrice{}, err
		}
		runNo, err := parseInt(idx.get(row, "RUNNO"))
		if err != nil {
			return FCASPrice{}, fmt.Errorf("runno: %w", err)
		}
		p := FCASPrice{SettlementDate: sd, RunNo: runNo, RegionID: idx.get(row, "REGIONID")}
		fields := []struct {
			column string
			dst    *float64
		}{
			{"RAISE6SECRRP", &p.Raise6Sec},
			{"RAISE60SECRRP", &p.Raise60Sec},
			{"RAISE5MINRRP", &p.Raise5Min},
			{"RAISEREGRRP", &p.RaiseReg},
			{"LOWER6SECRRP", &p.Lower6Sec},
			{"LOWER60SECRRP", &p.Lower60Sec},
			{"LOWER5MINRRP", &p.Lower5Min},
			{"LOWERREGRRP", &p.LowerReg},
		}
		for _, f := range fields {
			v, err := parseFloat(idx.get(row, f.column))
			if err != nil {
				return FCASPrice{}, fmt.Errorf("%s: %w", f.column, err)
			}
			*f.dst = v
		}
		return p, nil
	})
}
