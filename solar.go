package tesla

type Solar struct {
	EnergySiteId         int64  `json:"energy_site_id"`
	ResourceType         string `json:"resource_type"`
	Id                   string `json:"id"`
	AssetSiteId          string `json:"asset_site_id"`
	SolarPower           int    `json:"solar_power"`
	SolarType            string `json:"solar_type"`
	SyncGridAlertEnabled bool   `json:"sync_grid_alert_enabled"`
	BreakerAlertEnabled  bool   `json:"breaker_alert_enabled"`
	Components           struct {
		Battery    bool   `json:"battery"`
		Solar      bool   `json:"solar"`
		SolarType  string `json:"solar_type"`
		Grid       bool   `json:"grid"`
		LoadMeter  bool   `json:"load_meter"`
		MarketType string `json:"market_type"`
	} `json:"components"`
}
