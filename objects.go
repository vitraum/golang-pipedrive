package pipedrive

import "time"

type Pipeline struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	UrlTitle   string `json:"url_title"`
	OrderNr    int    `json:"order_nr"`
	Active     bool   `json:"active"`
	AddTime    Time   `json:"add_time"`
	UpdateTime Time   `json:"update_time"`
	Selected   bool   `json:"selected"`
}

type apiResult struct {
	Success        bool `json:"success"`
	AdditionalData struct {
		Pagination struct {
			Start                 int  `json:"start"`
			Limit                 int  `json:"limit"`
			MoreItemsInCollection bool `json:"more_items_in_collection"`
		} `json:"pagination"`
	} `json:"additional_data"`
}

type DealRef struct {
	Id int `json:"id"`
	//"user_id": 872124,
	Person struct {
		ID    int    `json:"value"`
		Name  string `json:"name"`
		Email []struct {
			Label   string `json:"label"`
			Value   string `json:"value"`
			Primary bool   `json:"primary"`
		} `json:"email"`
	} `json:"person_id"`
	Organization struct {
		ID   int    `json:"value"`
		Name string `json:"name"`
	} `json:"org_id"`
	Stage           int     `json:"stage_id"`
	StageChangetime *Time   `json:"stage_change_time"`
	Title           string  `json:"title"`
	Value           float64 `json:"value"`
	Added           Time    `json:"add_time"`
	Updated         *Time   `json:"update_time"`
	Status          string  `json:"status"`
	WonAt           *Time   `json:"won_time"`
	LostAt          *Time   `json:"lost_time"`
	LastActivity    *Date   `json:"last_activity_date"`
	Source          string  `json:"898dea9060ea3bb803e6a4f58c3c780b44e77cf7"`
}

type DealRefs []DealRef

type Deal struct {
	Id int `json:"id"`
	//"user_id": 872124,
	Person          int     `json:"person_id"`
	Organization    int     `json:"org_id"`
	Stage           int     `json:"stage_id"`
	StageChangetime *Time   `json:"stage_change_time"`
	Title           string  `json:"title"`
	Value           float64 `json:"value"`
	Added           Time    `json:"add_time"`
	Updated         *Time   `json:"update_time"`
	Status          string  `json:"status"`
	WonAt           *Time   `json:"won_time"`
	LostAt          *Time   `json:"lost_time"`
	LastActivity    *Date   `json:"last_activity_date"`
	Source          string  `json:"898dea9060ea3bb803e6a4f58c3c780b44e77cf7"`

	/*
	   "currency": "EUR",
	   "active": true,
	   "deleted": false,
	   "next_activity_date": "2015-12-14",
	   "next_activity_time": "10:15:00",
	   "next_activity_id": 5309,
	   "last_activity_id": 5303,
	   "lost_reason": null,
	   "visible_to": "3",
	   "close_time": null,
	   "pipeline_id": 1,
	   "products_count": null,
	   "files_count": 3,
	   "notes_count": 5,
	   "followers_count": 1,
	   "email_messages_count": 6,
	   "activities_count": 24,
	   "done_activities_count": 23,
	   "undone_activities_count": 1,
	   "reference_activities_count": 6,
	   "participants_count": 1,
	   "expected_close_date": null,
	   "a1c687608ebada2429e8d2df245571a4e38a6774": null,
	   "7d3fdbd55ff9f4fd0fac2d8933fa93f07a52bbca": null,
	   "778351f1b8b4010e770fc658aaf2e5fb2863628b": "60",
	   "6a89023701ead149ac22ac01029a9af9978fc9d0": null,
	   "c232115a7bc200de54dede15b32b5e2226ee5233": null,
	   "24522f29110e41f2a40de9ef98d490d6183ffa5e": null,
	   "9a4c96f6ab4f39cfbc791788a77c6bee93d8dba2": null,
	   "08f2639a112a0389520e84aa093bdcd0ddfdaafe": null,
	   "fdd4beef0370d4acc60ac43811e563422f440b81": null,
	   "19cef73a8ff77b70bf05736552155a7b9f97a36f": "162",
	   "76e9384783a680b1cd5e5d6de7ecb52b08a8fef3": null,
	   "898dea9060ea3bb803e6a4f58c3c780b44e77cf7": "150",
	   "e09493a4cf121a8a7a64ba9544863328bc5a85fe": null,
	   "2f762ed412b570a329ffb44845334503634caa57": null,
	   "stage_order_nr": 7,
	   "person_name": "Kai Verwayen",
	   "org_name": "Kai Verwayen",
	   "next_activity_subject": "Angebotsnachfass",
	   "next_activity_type": "call",
	   "next_activity_duration": null,
	   "next_activity_note": "Uhrzeit: einfach probieren",
	   "formatted_value": "1.384,70 €",
	   "weighted_value": 1384.7,
	   "formatted_weighted_value": "1.384,70 €",
	   "rotten_time": null,
	   "owner_name": "Timo Selent",
	   "cc_email": "vitraumgmbh+deal674@pipedrivemail.com",
	   "org_hidden": false,
	   "person_hidden": false
	*/
}

type Deals []Deal

type DealUpdate struct {
	// "id": 20215,
	//"company_id": 452851,
	StoryData struct {
		ActionType string `json:"action_type"`
		ChangeLog  []struct {
			OldValue  interface{} `json:"old_value"`
			NewValue  interface{} `json:"new_value"`
			FieldName string      `json:"field_name"`
			FieldKey  string      `json:"field_key"`
		} `json:"change_log"`
		AddTime Time `json:"add_time"`
	} `json:"story_data"`
	AddTime Time `json:"add_time"`
	/*
	   "story_data": {
	       "action_type": "edit",
	       "deal_id": 674,
	       "deal_public_id": 674,
	       "deal_title": "Herr Kai Verwayen / 56412 Heiligenroth",
	       "person_id": 518,
	       "person_name": "Kai Verwayen",
	       "org_id": 2703,
	       "org_name": "Kai Verwayen",
	       "user_id": 872124,
	       "user_name": "Timo Selent",
	       "company_id": 452851,
	       "company_name": "Vitraum GmbH",
	       "add_time": "2015-10-29 10:59:37",
	   },
	   "item_type": "deal",
	   "item_id": 674,
	   "writer_id": 872124,
	   "story_type": "text",
	   "active_flag": true,
	   "like_story": null,
	   "last_update": "2015-10-29 10:59:37",
	   "writer_name": "Timo Selent",
	   "writer_email": "timo.selent@vitraum.de",
	   "has_pic": "",
	   "pic_hash": ""
	*/
}

type DealUpdates []DealUpdate

type Stage struct {
	Id              int    `json:"id"`
	OrderNr         int    `json:"order_nr"`
	Name            string `json:"name"`
	Active          bool   `json:"active_flag"`
	AddTime         Time   `json:"add_time"`
	UpdateTime      Time   `json:"update_time"`
	Rotten          bool   `json:"rotten_flag"`
	RottenDays      int    `json:"rotten_days"`
	PipelineID      int    `json:"pipeline_id"`
	PipelineName    string `json:"pipeline_name"`
	DealProbability int    `json:"deal_probability"`
}

type Stages []Stage

type StagesResult struct {
	apiResult
	Data Stages `json:"data"`
}

type Filter struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Active     bool   `json:"active_flag"`
	Type       string `json:"type"`
	AddTime    *Time  `json:"add_time"`
	UpdateTime *Time  `json:"update_time"`
	VisibleTo  string `json:"visible_to"`
	UserID     int    `json:"user_id"`
}

type Filters []Filter

type DealFlowUpdate struct {
	Phase           string
	PiT             Time
	Duration        float64
	PhaseTouchdowns int
}

type DealFlowUpdates []DealFlowUpdate

func (slice DealFlowUpdates) Len() int {
	return len(slice)
}

func (slice DealFlowUpdates) Less(i, j int) bool {
	return slice[i].PiT.Unix() < slice[j].PiT.Unix()
}

func (slice DealFlowUpdates) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type PipelineChangeResult struct {
	Deal    Deal
	Updates DealFlowUpdates
}

type PipelineChangeResults []PipelineChangeResult

func (cr PipelineChangeResult) DecisionTime() time.Time {
	end := time.Now()
	if cr.Deal.Status == "won" {
		end = cr.Deal.WonAt.Time
	} else if cr.Deal.Status == "lost" {
		end = cr.Deal.LostAt.Time
	}
	return end
}

type DealField struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	key        string `json:"key"`
	OrderNr    int    `json:"order_nr"`
	AddTime    Time   `json:"add_time"`
	UpdateTime Time   `json:"update_time"`
	Options    *[]struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	} `json:"options"`

	/*
	   "field_type": "enum",
	   "active_flag": true,
	   "edit_flag": true,
	   "index_visible_flag": true,
	   "details_visible_flag": true,
	   "add_visible_flag": true,
	   "important_flag": true,
	   "bulk_edit_allowed": true,
	*/
}
