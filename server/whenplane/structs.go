package whenplane

import (
	"encoding/json"
	"time"
)

type Aggregate struct {
	Youtube struct {
		IsLive   bool    `json:"isLive"`
		Upcoming bool    `json:"upcoming"`
		VideoID  *string `json:"videoId"`
	} `json:"youtube"`
	Twitch struct {
		IsLive bool `json:"isLive"`
		IsWAN  bool `json:"isWAN"`
	} `json:"twitch"`
	SpecialStream bool `json:"specialStream"`
	Floatplane    struct {
		IsLive         bool   `json:"isLive"`
		IsWAN          bool   `json:"isWAN"`
		IsThumbnailNew bool   `json:"isThumbnailNew"`
		Thumbnail      string `json:"thumbnail"`
		Title          string `json:"title"`
	} `json:"floatplane"`
	NotablePeople struct {
		Bocabola struct {
			IsLive  bool      `json:"isLive"`
			Started time.Time `json:"started"`
			Title   string    `json:"title"`
			Name    string    `json:"name"`
			Channel string    `json:"channel"`
			Game    string    `json:"game"`
		} `json:"bocabola"`
		Buhdan struct {
			IsLive  bool   `json:"isLive"`
			Name    string `json:"name"`
			Channel string `json:"channel"`
		} `json:"buhdan"`
		LukeLafr struct {
			IsLive  bool   `json:"isLive"`
			Name    string `json:"name"`
			Channel string `json:"channel"`
		} `json:"luke_lafr"`
		Iitskasino struct {
			IsLive  bool   `json:"isLive"`
			Name    string `json:"name"`
			Channel string `json:"channel"`
		} `json:"iitskasino"`
	} `json:"notablePeople"`
	HasDone    bool `json:"hasDone"`
	IsThereWan struct {
		Text  interface{} `json:"text"`
		Image interface{} `json:"image"`
	} `json:"isThereWan"`
	Votes []struct {
		Name    string `json:"name"`
		Comment string `json:"comment,omitempty"`
		Votes   int    `json:"votes"`
		Time    int    `json:"time"`
	} `json:"votes"`
	ReloadNumber int `json:"reloadNumber"`
}

func AggregateFromJSON(data string) (Aggregate, error) {
	var aggregate Aggregate

	err := json.Unmarshal([]byte(data), &aggregate)

	return aggregate, err
}
