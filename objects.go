package pipedrive

import "time"

// Pipeline models the pipeline API object
type Pipeline struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	URLTitle   string `json:"url_title"`
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

// DealRef models a Pipedrive Deal
type DealRef struct {
	ID int `json:"id"`
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
	LostReason      string  `json:"lost_reason"`
	LastActivity    *Date   `json:"last_activity_date"`
	Source          string  `json:"898dea9060ea3bb803e6a4f58c3c780b44e77cf7"`
	LeadDate        string  `json:"19cef73a8ff77b70bf05736552155a7b9f97a36f"`
	User            struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		HasPicture bool   `json:"has_pic"`
		Active     bool   `json:"active_flag"`
		Value      int    `json:"value"`
		// PicHash
	} `json:"user_id"`
}

type Activity struct {
	ID                 int         `json:"id"`
	CompanyID          int         `json:"company_id"`
	UserID             int         `json:"user_id"`
	Done               bool        `json:"done"`
	Type               string      `json:"type"`
	ReferenceType      string      `json:"reference_type"`
	ReferenceID        interface{} `json:"reference_id"`
	DueDate            string      `json:"due_date"`
	DueTime            string      `json:"due_time"`
	Duration           string      `json:"duration"`
	AddTime            string      `json:"add_time"`
	MarkedAsDoneTime   string      `json:"marked_as_done_time"`
	Subject            string      `json:"subject"`
	DealID             int         `json:"deal_id"`
	OrgID              int         `json:"org_id"`
	PersonID           int         `json:"person_id"`
	ActiveFlag         bool        `json:"active_flag"`
	UpdateTime         string      `json:"update_time"`
	GcalEventID        interface{} `json:"gcal_event_id"`
	GoogleCalendarID   interface{} `json:"google_calendar_id"`
	GoogleCalendarEtag interface{} `json:"google_calendar_etag"`
	Note               string      `json:"note"`
	NoteClean          string      `json:"note_clean"`
	Participants       []struct {
		PersonID    int  `json:"person_id"`
		PrimaryFlag bool `json:"primary_flag"`
	} `json:"participants"`
	PersonName       string `json:"person_name"`
	OrgName          string `json:"org_name"`
	DealTitle        string `json:"deal_title"`
	AssignedToUserID int    `json:"assigned_to_user_id"`
	CreatedByUserID  int    `json:"created_by_user_id"`
	OwnerName        string `json:"owner_name"`
	PersonDropboxBcc string `json:"person_dropbox_bcc"`
	DealDropboxBcc   string `json:"deal_dropbox_bcc"`
}

// DealRefs is a list of Deals
type DealRefs []DealRef

// Deal models a Pipedrive Deal
type Deal struct {
	ID int `json:"id"`
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
	LostReason      *string `json:"lost_reason"`
	Source          string  `json:"898dea9060ea3bb803e6a4f58c3c780b44e77cf7"`
	LeadDate        string  `json:"19cef73a8ff77b70bf05736552155a7b9f97a36f"`

	/*
	   "currency": "EUR",
	   "active": true,
	   "deleted": false,
	   "next_activity_date": "2015-12-14",
	   "next_activity_time": "10:15:00",
	   "next_activity_id": 5309,
	   "last_activity_id": 5303,
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
	Deal            Deal
	PipelineUpdates DealFlowUpdates
	Updates         DealUpdates
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
	Key        string `json:"key"`
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

type Organization struct {
	ID        int    `json:"id"`
	CompanyID int    `json:"company_id"`
	Name      string `json:"name"`
	OwnerID   struct {
		ID         int         `json:"id"`
		Name       string      `json:"name"`
		Email      string      `json:"email"`
		HasPic     bool        `json:"has_pic"`
		PicHash    interface{} `json:"pic_hash"`
		ActiveFlag bool        `json:"active_flag"`
		Value      int         `json:"value"`
	} `json:"owner_id"`
	OpenDealsCount           int         `json:"open_deals_count"`
	RelatedOpenDealsCount    int         `json:"related_open_deals_count"`
	ClosedDealsCount         int         `json:"closed_deals_count"`
	RelatedClosedDealsCount  int         `json:"related_closed_deals_count"`
	EmailMessagesCount       int         `json:"email_messages_count"`
	PeopleCount              int         `json:"people_count"`
	ActivitiesCount          int         `json:"activities_count"`
	DoneActivitiesCount      int         `json:"done_activities_count"`
	UndoneActivitiesCount    int         `json:"undone_activities_count"`
	ReferenceActivitiesCount int         `json:"reference_activities_count"`
	FilesCount               int         `json:"files_count"`
	NotesCount               int         `json:"notes_count"`
	FollowersCount           int         `json:"followers_count"`
	WonDealsCount            int         `json:"won_deals_count"`
	RelatedWonDealsCount     int         `json:"related_won_deals_count"`
	LostDealsCount           int         `json:"lost_deals_count"`
	RelatedLostDealsCount    int         `json:"related_lost_deals_count"`
	ActiveFlag               bool        `json:"active_flag"`
	CategoryID               interface{} `json:"category_id"`
	PictureID                interface{} `json:"picture_id"`
	CountryCode              interface{} `json:"country_code"`
	FirstChar                string      `json:"first_char"`
	UpdateTime               string      `json:"update_time"`
	AddTime                  string      `json:"add_time"`
	VisibleTo                string      `json:"visible_to"`
	NextActivityDate         interface{} `json:"next_activity_date"`
	NextActivityTime         interface{} `json:"next_activity_time"`
	NextActivityID           interface{} `json:"next_activity_id"`
	LastActivityID           int         `json:"last_activity_id"`
	LastActivityDate         string      `json:"last_activity_date"`
	Address                  string      `json:"address"`
	AddressSubpremise        string      `json:"address_subpremise"`
	AddressStreetNumber      string      `json:"address_street_number"`
	AddressRoute             string      `json:"address_route"`
	AddressSublocality       string      `json:"address_sublocality"`
	AddressLocality          string      `json:"address_locality"`
	AddressAdminAreaLevel1   string      `json:"address_admin_area_level_1"`
	AddressAdminAreaLevel2   string      `json:"address_admin_area_level_2"`
	AddressCountry           string      `json:"address_country"`
	AddressPostalCode        string      `json:"address_postal_code"`
	AddressFormattedAddress  string      `json:"address_formatted_address"`
	Label                    interface{} `json:"label"`
	CcEmail                  string      `json:"cc_email"`
	OwnerName                string      `json:"owner_name"`
	EditName                 bool        `json:"edit_name"`
	LastActivity             struct {
		ID                         int         `json:"id"`
		CompanyID                  int         `json:"company_id"`
		UserID                     int         `json:"user_id"`
		Done                       bool        `json:"done"`
		Type                       string      `json:"type"`
		ReferenceType              string      `json:"reference_type"`
		ReferenceID                interface{} `json:"reference_id"`
		DueDate                    string      `json:"due_date"`
		DueTime                    string      `json:"due_time"`
		Duration                   string      `json:"duration"`
		BusyFlag                   interface{} `json:"busy_flag"`
		AddTime                    string      `json:"add_time"`
		MarkedAsDoneTime           string      `json:"marked_as_done_time"`
		LastNotificationTime       interface{} `json:"last_notification_time"`
		LastNotificationUserID     interface{} `json:"last_notification_user_id"`
		NotificationLanguageID     interface{} `json:"notification_language_id"`
		Subject                    string      `json:"subject"`
		PublicDescription          interface{} `json:"public_description"`
		CalendarSyncIncludeContext interface{} `json:"calendar_sync_include_context"`
		Location                   interface{} `json:"location"`
		OrgID                      int         `json:"org_id"`
		PersonID                   int         `json:"person_id"`
		DealID                     int         `json:"deal_id"`
		ActiveFlag                 bool        `json:"active_flag"`
		UpdateTime                 string      `json:"update_time"`
		UpdateUserID               interface{} `json:"update_user_id"`
		GcalEventID                interface{} `json:"gcal_event_id"`
		GoogleCalendarID           interface{} `json:"google_calendar_id"`
		GoogleCalendarEtag         interface{} `json:"google_calendar_etag"`
		SourceTimezone             interface{} `json:"source_timezone"`
		RecRule                    interface{} `json:"rec_rule"`
		RecRuleExtension           interface{} `json:"rec_rule_extension"`
		RecMasterActivityID        interface{} `json:"rec_master_activity_id"`
		Note                       string      `json:"note"`
		CreatedByUserID            int         `json:"created_by_user_id"`
		LocationSubpremise         interface{} `json:"location_subpremise"`
		LocationStreetNumber       interface{} `json:"location_street_number"`
		LocationRoute              interface{} `json:"location_route"`
		LocationSublocality        interface{} `json:"location_sublocality"`
		LocationLocality           interface{} `json:"location_locality"`
		LocationAdminAreaLevel1    interface{} `json:"location_admin_area_level_1"`
		LocationAdminAreaLevel2    interface{} `json:"location_admin_area_level_2"`
		LocationCountry            interface{} `json:"location_country"`
		LocationPostalCode         interface{} `json:"location_postal_code"`
		LocationFormattedAddress   interface{} `json:"location_formatted_address"`
		Attendees                  interface{} `json:"attendees"`
		Participants               []struct {
			PersonID    int  `json:"person_id"`
			PrimaryFlag bool `json:"primary_flag"`
		} `json:"participants"`
		Series           interface{} `json:"series"`
		OrgName          string      `json:"org_name"`
		PersonName       string      `json:"person_name"`
		DealTitle        string      `json:"deal_title"`
		OwnerName        string      `json:"owner_name"`
		PersonDropboxBcc string      `json:"person_dropbox_bcc"`
		DealDropboxBcc   string      `json:"deal_dropbox_bcc"`
		AssignedToUserID int         `json:"assigned_to_user_id"`
		File             interface{} `json:"file"`
	} `json:"last_activity"`
	NextActivity interface{} `json:"next_activity"`
}
