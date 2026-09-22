package nem

// NEMWEB report definitions. Directory and file prefix are taken from the
// live NEMWEB listings.
var (
	// ReportMarketNotice is the human-formatted market notice report.
	ReportMarketNotice = Report{Dir: "Market_Notice", Prefix: "NEMITWEB1_MKTNOTICE_", Kind: KindText}
	// ReportDispatchIS is the five-minute dispatch, price, and regional summary report.
	ReportDispatchIS = Report{Dir: "DispatchIS_Reports", Prefix: "PUBLIC_DISPATCHIS_", Kind: KindNested}
	// ReportDispatch is the legacy dispatch report.
	ReportDispatch = Report{Dir: "Dispatch_Reports", Prefix: "PUBLIC_DISPATCH_", Kind: KindNested}
	// ReportDispatchSCADA is the unit-level SCADA report.
	ReportDispatchSCADA = Report{Dir: "Dispatch_SCADA", Prefix: "PUBLIC_DISPATCHSCADA_", Kind: KindNested}
	// ReportP5 is the five-minute predispatch report.
	ReportP5 = Report{Dir: "P5_Reports", Prefix: "PUBLIC_P5MIN_", Kind: KindNested}
	// ReportNetwork is the network constraint report.
	ReportNetwork = Report{Dir: "Network", Prefix: "PUBLIC_NETWORK_", Kind: KindNested}
	// ReportMTPASADUIDAvailability is the medium-term unit availability report.
	ReportMTPASADUIDAvailability = Report{Dir: "MTPASA_DUIDAvailability", Prefix: "PUBLIC_MTPASADUIDAVAILABILITY_", Kind: KindNested}
	// ReportPDPASADUIDAvailability is the short-term unit availability report.
	ReportPDPASADUIDAvailability = Report{Dir: "PDPASA_DUIDAvailability", Prefix: "PUBLIC_PDPASA_DUIDAVAILABILITY_", Kind: KindNested}
	// ReportSevenDayOutlookFull is the seven-day supply and demand outlook report.
	ReportSevenDayOutlookFull = Report{Dir: "SEVENDAYOUTLOOK_FULL", Prefix: "PUBLIC_SEVENDAYOUTLOOK_FULL_", Kind: KindNested}
	// ReportTradingIS is the thirty-minute trading report.
	ReportTradingIS = Report{Dir: "TradingIS_Reports", Prefix: "PUBLIC_TRADINGIS_", Kind: KindNested}
)
