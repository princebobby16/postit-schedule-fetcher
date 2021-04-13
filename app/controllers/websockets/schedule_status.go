package websockets

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"gitlab.com/pbobby001/postit-schedule-status/db"
	"gitlab.com/pbobby001/postit-schedule-status/pkg"
	"gitlab.com/pbobby001/postit-schedule-status/pkg/logs"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func Writer(conn *websocket.Conn, tenantNamespace string, connection *sql.DB) {
	for {
		ticker := time.NewTicker(5 * time.Second)

		for t := range ticker.C {
			logs.Logger.Info("Updating status: ", t)
			statuses, err := FetchStatuses(connection, tenantNamespace)
			if err != nil {
				_ = logs.Logger.Error(err)
				return
			}

			jsonBytes, err := json.Marshal(statuses)
			if err != nil {
				_ = logs.Logger.Error(err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, jsonBytes)
			if err != nil {
				_ = logs.Logger.Error(err)
				return
			}

		}
	}
}

func FetchStatuses(connection *sql.DB, tenantNamespace string) ([]pkg.ScheduleStatus, error) {

	//  Prepare the query
	query := fmt.Sprintf("SELECT schedule_id, schedule_title, schedule_from, schedule_to, post_ids FROM %s.schedule WHERE is_due = $1", tenantNamespace)

	// run the query
	rows, err := connection.Query(query, true)
	if err != nil {
		return nil, err
	}

	var schedules []pkg.PostSchedule
	for rows.Next() {
		var scheduleData pkg.PostSchedule
		err = rows.Scan(
			&scheduleData.ScheduleId,
			&scheduleData.ScheduleTitle,
			&scheduleData.From,
			&scheduleData.To,
			pq.Array(&scheduleData.PostIds),
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, scheduleData)
	}

	// 3. create a global post list of posted posts to store the posts that are due
	var postedPosts []pkg.ScheduledPost
	// 4. create a global schedule list to accommodate the posts
	var scheduleStatuses []pkg.ScheduleStatus

	if schedules != nil {
		// 1. iterate over the schedules
		for _, schedule := range schedules {
			// 2. iterate over the post ids in the schedule
			for _, id := range schedule.PostIds {
				query := fmt.Sprintf("SELECT "+
					"post_id, "+
					"facebook_post_id, "+
					"post_message, "+
					"hash_tags, "+
					"post_images, "+
					"image_paths, "+
					"scheduled, "+
					"post_fb_status "+
					"FROM %s.post "+
					"WHERE "+
					"post_id = $1 "+
					"AND post_fb_status = $2 AND "+
					"scheduled = $3",
					tenantNamespace,
				)
				rows, err := db.Connection.Query(query, id, true, true)
				if err != nil {
					return nil, err
				}

				for rows.Next() {
					var post pkg.ScheduledPost
					err = rows.Scan(
						&post.PostId,
						&post.FacebookPostId,
						&post.PostMessage,
						pq.Array(&post.HashTags),
						pq.Array(&post.PostImages),
						pq.Array(&post.ImagePaths),
						&post.Scheduled,
						&post.PostFbStatus,
					)
					if err != nil {
						return nil, err
					}

					postedPosts = append(postedPosts, post)
				}
			}
			// 5. add the post list to the particular schedule it belongs to
			scheduleStatuses = append(scheduleStatuses, pkg.ScheduleStatus{
				ScheduleId:    schedule.ScheduleId,
				ScheduleTitle: schedule.ScheduleTitle,
				From:          schedule.From,
				To:            schedule.To,
				TotalPost:     len(schedule.PostIds),
				Posts:         postedPosts,
				PostCount:     len(postedPosts),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			})
		}
	}

	return scheduleStatuses, nil
}

func ScheduleStatus(w http.ResponseWriter, r *http.Request) {

	logs.Logger.Info("connecting to websocket")

	ws, err := upgrade(w, r)
	if err != nil {
		_ = logs.Logger.Error(err)
		return
	}

	var webSocketHandshake pkg.WebSocketHandShakeData
	err = ws.ReadJSON(&webSocketHandshake)
	if err != nil {
		_ = logs.Logger.Error(err)
		return
	}
	logs.Logger.Info(webSocketHandshake)

	// validate token
	err = pkg.WebSocketTokenValidateToken(webSocketHandshake.Token, webSocketHandshake.TenantNamespace)
	if err != nil {
		_ = logs.Logger.Error(err)
		err = ws.Close()
		if err != nil {
			_ = logs.Logger.Error(err)
		}
		return
	}
	//q := `UPDATE postit.scheduled_post SET post_status = true WHERE post_id=$1`
	//_, err = db.Connection.Exec(q, "298bccf8-c103-474b-b708-a8797860feb0")
	//if err != nil {
	//	logs.Logger.Error(err)
	//	return
	//}

	logs.Logger.Info("connection upgraded")
	go Writer(ws, webSocketHandshake.TenantNamespace, db.Connection)
}
