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

	ID ID `json:"_id"`

	EventID EventID `json:"event_id"`

	BestTimes []BestTimes `json:"best_times"`

	Cls []Cls `json:"cls"`

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

// Cls seems to refer to the kart type.
type Cls struct {
	C  string `json:"c"`
	Nm string `json:"nm"`

	ID int `json:"id"`
}

// Competitors describe teams of drivers at the event.
type Competitors struct {
	Ch  interface{} `json:"ch"`
	Cid interface{} `json:"cid"`
	E   interface{} `json:"e"`
	N   interface{} `json:"n"`
	Rt  interface{} `json:"rt"`

	ID ID `json:"_id"`

	// C is the type of kart driven by the team.
	C string `json:"c"`

	Cc string `json:"cc"`
	Cm string `json:"cm"`
	Gq string `json:"gq"`

	// Na is the name of the team.
	Na string `json:"na"`

	// Nm is the kart number of the team.
	Nm string `json:"nm"`

	Rcid string `json:"rcid"`
	Sc   string `json:"sc"`
	Scid string `json:"scid"`
	Tdn  string `json:"tdn"`

	Laps []Laps `json:"laps"`

	Gp int `json:"gp"`

	Redact bool `json:"redact"`
}

// Laps hold times and sectors for the various drivers.
type Laps struct {
	Sectors []Sectors `json:"sectors,omitempty"`

	Lt int `json:"lt"`
	N  int `json:"n"`
	P  int `json:"p"`
	Tt int `json:"tt"`

	Blor bool `json:"blor,omitempty"`
	Pb   bool `json:"pb,omitempty"`
}

// Sectors number three on the race track.
type Sectors struct {
	Sid string `json:"sid"`

	St int `json:"st"`
}

// Results contains various timings, averages, best sectors, and so on.
type Results struct {
	Bs2Or interface{} `json:"bs2or"`
	Bs1Or interface{} `json:"bs1or"`
	Pen   interface{} `json:"pen"`

	ID ID `json:"_id"`

	Avg   string `json:"avg"`
	B     string `json:"b"`
	Blavg string `json:"blavg"`
	Blt   string `json:"blt"`
	Fp    string `json:"fp"`
	G     string `json:"g"`
	P     string `json:"p"`
	Pls   string `json:"pls"`
	Scid  string `json:"scid"`
	T     string `json:"t"`

	Ty []interface{} `json:"ty"`

	Bln      int `json:"bln"`
	BltMs    int `json:"blt_ms"`
	Bs1      int `json:"bs1"`
	Bs1L     int `json:"bs1l"`
	Bs2      int `json:"bs2"`
	Bs2L     int `json:"bs2l"`
	Bs3      int `json:"bs3"`
	Bs3L     int `json:"bs3l"`
	Nl       int `json:"nl"`
	TMs      int `json:"t_ms"`
	Ultimate int `json:"ultimate"`

	Gpts float64 `json:"gpts"`
	Pts  float64 `json:"pts"`

	Blor     bool `json:"blor"`
	Bs3Or    bool `json:"bs3or"`
	Finisher bool `json:"finisher"`
	Starter  bool `json:"starter"`
}
