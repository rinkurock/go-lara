package sample

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"

	c "app/config"
	"app/conn"
	m "app/models"
)

func PostData(userId uint64) (bool, string, error) {
	var message = "Failed to get users"
	var data m.SamplePostReq
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		log.Errorln(err)
		return false, message, err
	}
	con := c.GetConfig().Services.Sample
	spew.Dump(con)
	url := fmt.Sprintf("%s/api/v1/users/%d/users", con.Host, userId)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		log.Warn(err)
		return false, message, errors.New("failed to create get user PostData sample request")
	}
	req.Header.Set("Content-Type", "application/json")
	client := conn.GetSampleCon()
	res, err := client.Do(req)
	if err != nil {
		log.Println("failed to sample service PostData request:", err)
		return false, message, err
	}
	if res == nil {
		log.Println("failed to sample service PostData request response nil")
		return false, message, err
	} else {
		if res.StatusCode == 201 || res.StatusCode == 200 {
			return true, message, nil
		} else {
			var r m.SamplePostRes
			err := json.NewDecoder(res.Body).Decode(&r)
			if err == nil && r.Data.Name != "" {
				message = r.Data.Name
			} else {
				message = "Not get any name"
			}
			log.Println("failed to sample service PostData request status:", res.StatusCode)
			return false, message, nil
		}
	}
}
