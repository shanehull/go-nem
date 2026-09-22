package nem_test

import (
	"errors"
	"testing"

	"github.com/shanehull/go-nem"
)

func decodeDispatchIS(t *testing.T) []nem.Table {
	t.Helper()
	tables, err := nem.Parse(readFixture(t, "dispatchis.csv"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return tables
}

func TestDecodeDispatchPrice(t *testing.T) {
	prices, err := nem.DecodeDispatchPrice(decodeDispatchIS(t))
	if err != nil {
		t.Fatalf("DecodeDispatchPrice: %v", err)
	}
	if len(prices) != 2 {
		t.Fatalf("got %d prices, want 2", len(prices))
	}
	if prices[0].RegionID != "NSW1" || prices[0].RRP != 86.69 {
		t.Errorf("NSW price = %+v", prices[0])
	}
	if !prices[0].SettlementDate.Equal(mustTime(t, "2026-09-21T08:10:00+10:00")) {
		t.Errorf("settlement = %s", prices[0].SettlementDate)
	}
	if prices[1].RRP != -12.5 {
		t.Errorf("VIC price = %v, want -12.5", prices[1].RRP)
	}
}

func TestDecodeFCASPrice(t *testing.T) {
	prices, err := nem.DecodeFCASPrice(decodeDispatchIS(t))
	if err != nil {
		t.Fatalf("DecodeFCASPrice: %v", err)
	}
	if len(prices) != 2 {
		t.Fatalf("got %d prices, want 2", len(prices))
	}
	if prices[0].Raise6Sec != 12.34 || prices[0].LowerReg != 3.21 {
		t.Errorf("NSW FCAS = %+v", prices[0])
	}
}

func TestDecodeDispatchRegionSum(t *testing.T) {
	sums, err := nem.DecodeDispatchRegionSum(decodeDispatchIS(t))
	if err != nil {
		t.Fatalf("DecodeDispatchRegionSum: %v", err)
	}
	if len(sums) != 2 {
		t.Fatalf("got %d sums, want 2", len(sums))
	}
	if sums[0].RegionID != "NSW1" || sums[0].TotalDemand != 9123.45 {
		t.Errorf("NSW sum = %+v", sums[0])
	}
	if sums[0].LORSurplus != 450 || sums[0].NetInterchange != -321 {
		t.Errorf("NSW surplus/interchange = %+v", sums[0])
	}
}

func TestDecodeDispatchInterconnectorRes(t *testing.T) {
	flows, err := nem.DecodeDispatchInterconnectorRes(decodeDispatchIS(t))
	if err != nil {
		t.Fatalf("DecodeDispatchInterconnectorRes: %v", err)
	}
	if len(flows) != 1 {
		t.Fatalf("got %d flows, want 1", len(flows))
	}
	if flows[0].InterconnectorID != "VIC1-NSW1" || flows[0].MWFlow != 318.5 {
		t.Errorf("flow = %+v", flows[0])
	}
}

func TestDecodeDispatchUnitScada(t *testing.T) {
	tables, err := nem.Parse(readFixture(t, "unit_scada.csv"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	units, err := nem.DecodeDispatchUnitScada(tables)
	if err != nil {
		t.Fatalf("DecodeDispatchUnitScada: %v", err)
	}
	if len(units) != 3 {
		t.Fatalf("got %d units, want 3", len(units))
	}
	if units[0].DUID != "BARCSF1" || units[0].SCADAValue != 12.6 {
		t.Errorf("unit = %+v", units[0])
	}
}

func TestDecodeMissingTable(t *testing.T) {
	_, err := nem.DecodeDispatchInterconnectorRes(nil)
	if !errors.Is(err, nem.ErrTableNotFound) {
		t.Fatalf("error = %v, want ErrTableNotFound", err)
	}
}
