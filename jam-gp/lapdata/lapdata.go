package lapdata

import "time"

// Event is the top-level data structure containing all event information.
type Event struct {
	EventDate string `json:"event_date"`
	EventName string `json:"event_name"`

	Officials []interface{} `json:"officials"`

	Session Session `json:"session"`

	Provisional bool `json:"provisional"`
	UnderAppeal bool `json:"underAppeal"`
}

// Session is one of the highest-level containers for data.
type Session struct {
	ID ID `json:"_id"`

	EventID EventID `json:"event_id"`

	StartTime time.Time `json:"start_time"`

	Com               interface{} `json:"com"`
	GridLastGenerated interface{} `json:"grid_last_generated"`
	GridVer           interface{} `json:"grid_ver"`
	ResLastIssued     interface{} `json:"res_last_issued"`
	ResVer            interface{} `json:"res_ver"`
	Tc                interface{} `json:"tc"`
	Temp              interface{} `json:"temp"`
	TimingSystem      interface{} `json:"timing_system"`
	Weather           interface{} `json:"weather"`

	Name          string `json:"name"`
	NumberAndName string `json:"number_and_name"`
	SessionType   string `json:"session_type"`
	Status        string `json:"status"`

	BestTimes []BestTimes `json:"best_times"`

	Cls []Class `json:"cls"`

	Competitors []Competitors `json:"competitors"`

	Pens []interface{} `json:"pens"`

	Results []Results `json:"results"`

	Dp      int `json:"dp"`
	Number  int `json:"number"`
	Sectors int `json:"sectors"`
	Tid     int `json:"tid"`

	Sg bool `json:"sg"`
	St bool `json:"st"`
}

// ID is a unique identifier for various data types.
type ID struct {
	Oid string `json:"$oid"`
}

// EventID is a unique ID for the event.
type EventID struct {
	Oid string `json:"$oid"`
}

// BestTimes holds best lap times.
type BestTimes struct {
	Scid string `json:"scid,omitempty"`
	Type string `json:"type"`

	L int `json:"l,omitempty"`
	T int `json:"t"`
}

// Class refers to the kart type.
type Class struct {
	// Code of the kart class.
	Code string `json:"c"`

	// Name of the kart class.
	Name string `json:"nm"`

	// ID of the kart class.
	ID int `json:"id"`
}

// Competitors describe teams of drivers at the event.
type Competitors struct {
	ID ID `json:"_id"`

	Ch  interface{} `json:"ch"`
	Cid interface{} `json:"cid"`
	E   interface{} `json:"e"`
	N   interface{} `json:"n"`
	Rt  interface{} `json:"rt"`

	// KartType is the type of kart driven by the team.
	KartType string `json:"c"`

	Cc string `json:"cc"`
	Cm string `json:"cm"`
	Gq string `json:"gq"`

	// Name is the name of the team.
	Name string `json:"na"`

	// Number is the kart number of the team.
	Number string `json:"nm"`

	Rcid string `json:"rcid"`
	Sc   string `json:"sc"`
	Scid string `json:"scid"`
	Tdn  string `json:"tdn"`

	Laps SliceOfLaps `json:"laps"`

	Gp int `json:"gp"`

	Redact bool `json:"redact"`
}

// Laps hold times and sectors for the various drivers.
type Laps struct {
	Sectors []Sectors `json:"sectors,omitempty"`

	// LapTime is the individual lap time.
	LapTime int `json:"lt"`

	N int `json:"n"`
	P int `json:"p"`

	// TotalTime is the cumulative total lap time so far.
	TotalTime int `json:"tt"`

	Blor bool `json:"blor,omitempty"`
	Pb   bool `json:"pb,omitempty"`
}

// Sectors number three on the race track.
type Sectors struct {
	// SectorID will refer to one of the three sectors on the track.
	SectorID string `json:"sid"`

	// SectorTime is the elapsed time in this sector on this lap.
	SectorTime int `json:"st"`
}

// Results contains various timings, averages, best sectors, and so on.
type Results struct {
	ID ID `json:"_id"`

	Bs2Or interface{} `json:"bs2or"`
	Bs1Or interface{} `json:"bs1or"`
	Pen   interface{} `json:"pen"`

	Avg   string `json:"avg"`
	B     string `json:"b"`
	Blavg string `json:"blavg"`
	Blt   string `json:"blt"`
	Fp    string `json:"fp"`
	G     string `json:"g"`
	P     string `json:"p"`
	Pls   string `json:"pls"`
	Scid  string `json:"scid"`

	// Total time elapsed.
	T string `json:"t"`

	Ty []interface{} `json:"ty"`

	// Best lap number.
	Bln int `json:"bln"`

	// Best lap time in milliseconds.
	BltMs int `json:"blt_ms"`

	// Best time for sector one.
	Bs1 int `json:"bs1"`

	// Best time for sector one was set on this lap.
	Bs1L int `json:"bs1l"`

	// Best time for sector two.
	Bs2 int `json:"bs2"`

	// Best time for sector two was set on this lap.
	Bs2L int `json:"bs2l"`

	// Best time for sector three.
	Bs3 int `json:"bs3"`

	// Best time for sector three was set on this lap.
	Bs3L int `json:"bs3l"`

	// The number of laps completed.
	Nl int `json:"nl"`

	// Total time elapsed in milliseconds.
	TMs int `json:"t_ms"`

	// Ultimate is the theoretical best lap based on the three best sector times.
	Ultimate int `json:"ultimate"`

	Gpts float64 `json:"gpts"`
	Pts  float64 `json:"pts"`

	Blor     bool `json:"blor"`
	Bs3Or    bool `json:"bs3or"`
	Finisher bool `json:"finisher"`
	Starter  bool `json:"starter"`
}
