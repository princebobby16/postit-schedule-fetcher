package pkg

import (
	"time"
)

type (
	WebSocketHandShakeData struct {
		TenantNamespace string `json:"tenant_namespace"`
		AuthToken       string `json:"auth_token"`
	}

	ScheduleStatus struct {
		ScheduleId    string          `json:"schedule_id"`
		ScheduleTitle string          `json:"schedule_title"`
		From          time.Time       `json:"from"`
		To            time.Time       `json:"to"`
		TotalPost     int             `json:"total_post"`
		Posts         []ScheduledPost `json:"posts"`
		PostCount     int             `json:"post_count"`
		CreatedAt     time.Time       `json:"created_at"`
		UpdatedAt     time.Time       `json:"updated_at"`
	}

	SocialMediaProfiles struct {
		Facebook []string `json:"facebook"`
		Twitter  []string `json:"twitter"`
		LinkedIn []string `json:"linked_in"`
	}

	PostSchedule struct {
		ScheduleId    string              `json:"schedule_id"`
		ScheduleTitle string              `json:"schedule_title"`
		PostToFeed    bool                `json:"post_to_feed"`
		From          time.Time           `json:"from"`
		To            time.Time           `json:"to"`
		PostIds       []string            `json:"post_ids"`
		Duration      float64             `json:"duration"`
		IsDue         bool                `json:"is_due"`
		Profiles      SocialMediaProfiles `json:"profiles"`
		CreatedOn     time.Time           `json:"created_on"`
		UpdatedOn     time.Time           `json:"updated_on"`
	}

	ScheduledPost struct {
		PostId         string    `json:"post_id"`
		FacebookPostId string    `json:"facebook_post_id"`
		PostMessage    string    `json:"post_message"`
		PostImages     [][]byte  `json:"post_image"`
		ImagePaths     []string  `json:"image_paths"`
		HashTags       []string  `json:"hash_tags"`
		PostFbStatus   bool      `json:"post_fb_status"`
		Scheduled      bool      `json:"scheduled"`
		PostPriority   bool      `json:"post_priority"`
		CreatedOn      time.Time `json:"created_on"`
		UpdatedOn      time.Time `json:"updated_on"`
	}

	//Post struct {
	//	PostId       string    `json:"post_id"`
	//	PostMessage  string    `json:"post_message"`
	//	PostImages   [][]byte  `json:"post_images"`
	//	ImagePaths   []string  `json:"image_paths"`
	//	HashTags     []string  `json:"hash_tags"`
	//	PostPriority bool      `json:"post_priority"`
	//	CreatedOn    time.Time `json:"created_on"`
	//	UpdatedOn    time.Time `json:"updated_on"`
	//}

	StandardResponse struct {
		Data Data `json:"data"`
		Meta Meta `json:"meta"`
	}

	DbPost struct {
		PostId         string    `json:"post_id"`
		FacebookPostId string    `json:"facebook_post_id"`
		PostMessage    string    `json:"post_message"`
		PostImages     [][]byte  `json:"post_images"`
		ImagePaths     []string  `json:"image_paths"`
		HashTags       []string  `json:"hash_tags"`
		Scheduled      bool      `json:"scheduled"`
		PostFbStatus   bool      `json:"post_fb_status"`
		PostTwStatus   bool      `json:"post_tw_status"`
		PostLiStatus   bool      `json:"post_li_status"`
		PostPriority   bool      `json:"post_priority"`
		CreatedOn      time.Time `json:"created_on"`
		UpdatedOn      time.Time `json:"updated_on"`
	}

	Data struct {
		Id        string `json:"id"`
		UiMessage string `json:"ui_message"`
	}

	Meta struct {
		Timestamp     time.Time `json:"timestamp"`
		TransactionId string    `json:"transaction_id"`
		TraceId       string    `json:"trace_id"`
		Status        string    `json:"status"`
	}
)
