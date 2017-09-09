package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/inloop/go-transport-queue/model"
)

func createHTTPServer(q *Queue, t model.Transport) *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.POST("/push", func(c *gin.Context) {
		if messages, err := t.BindResponse(c); err == nil {
			for _, message := range messages {
				if validationErr := message.Validate(); validationErr != nil {
					c.String(http.StatusBadRequest, "invalid payload: %s", validationErr.Error())
					return
				}
			}
			priority := c.Query("priority")
			for _, message := range messages {
				p, _ := strconv.ParseInt(priority, 10, 32)
				q.Push(uint8(p), message)
			}
			c.String(http.StatusCreated, "OK")
		}
	})

	return router
	// srv := &http.Server{Addr: ":" + port}
	//
	// // newrelicLicenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	// // newrelicAppName := os.Getenv("NEW_RELIC_APP_NAME")
	// // var app newrelic.Application
	// // if newrelicLicenseKey != "" && newrelicAppName != "" {
	// // 	config := newrelic.NewConfig(newrelicAppName, newrelicLicenseKey)
	// // 	_app, err := newrelic.NewApplication(config)
	// // 	app = _app
	// // 	if err != nil {
	// // 		log.Println("error creating newrelic", err)
	// // 	}
	// // }
	//
	// http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case "POST":
	// 		decoder := json.NewDecoder(r.Body)
	// 		defer r.Body.Close()
	//
	// 		if messages, err := t.DecodeMessages(decoder); err == nil {
	// 			for _, message := range messages {
	// 				if validationErr := message.Validate(); validationErr != nil {
	// 					http.Error(w, "invalid payload: "+validationErr.Error(), 400)
	// 					return
	// 				}
	// 			}
	// 			for _, message := range messages {
	// 				q.Push(0, message)
	// 			}
	// 		} else {
	// 			http.Error(w, err.Error(), 400)
	// 			return
	// 		}
	//
	// 		io.WriteString(w, "OK")
	// 	default:
	// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	// 	}
	// })
	//
	// return srv
}
