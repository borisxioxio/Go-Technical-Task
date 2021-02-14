package app

import (
	"srvguide/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
)


func GetFunction(c *gin.Context) {
	method, _ := c.GetQuery("method")
	if method == "nextraces-category-group"{
		sdata, _ := GetRace(c)
		c.Set("res", gin.H{
			"status":  1,
			"message": "ok",
			"data":    sdata})
		return
	}
	c.Set("res", gin.H{
		"status":  1,
		"message": "failed",
		"data":    ""})
}

// GetFunction ...
func GetRace(c *gin.Context) (ssdata map[string]interface{}, err error){
	var sdata = make(map[string]interface{})
	cmap, _ := c.GetQueryArray("include_categories[]")
	var mdata = make(map[string]interface{})
	for _, v := range cmap {
		// var stmp = make(map[string]interface{})
		var tp = make(map[string]interface{})
		cdata, err := model.GetCatRaces(v)
		if err != nil {
			return sdata,err
		}
		tp["race_ids"] = []interface{}{}
		var rtmp []string
		for _, c := range cdata {
			rtmp = append(rtmp, c.RaceID)
		}
		tp["race_ids"] = rtmp
		mdata[v] = tp
		// mdata = append(mdata, stmp)
	}
	sdata["category_race_map"] = mdata

	count, _ := c.GetQuery("count")
	var page = 5
	if count != "" {
		page, _ = strconv.Atoi(count)
	}
	data, err := model.GetRace(1, page)
	if err != nil {
		log.Error("error in operate db:" + err.Error())
		return sdata,nil
	}
	var res []map[string]interface{}
	for _, v := range data {
		var tmp = make(map[string]interface{})
		tmp["race_id"] = v.RaceID
		tmp["race_name"] = v.RaceName
		tmp["race_number"] = v.RaceNumber
		tmp["meeting_id"] = v.MeetingID
		mdata, err := model.GetMetInfo(v.MeetingID)
		if err != nil {
			tmp["meeting_name"] = ""
		} else {
			tmp["meeting_name"] = mdata.MeetingName
		}
		cdata, err := model.GetCatInfo(v.RaceID)
		if err != nil {
			tmp["category_id"] = ""
		} else {
			tmp["category_id"] = cdata.CategoryID
		}

		res = append(res, tmp)
	}
	sdata["race_summaries"] = res
	return sdata,nil
}

// func (app *App) postFunction(w http.ResponseWriter, r *http.Request) {
// 	_, err := app.Database.Exec("INSERT INTO `test` (name) VALUES ('myname')")
// 	if err != nil {
// 		log.Fatal("Database INSERT failed")
// 	}

// 	log.Println("You called a thing!")
// 	w.WriteHeader(http.StatusOK)
// }

// 内网使用，没有使用token验证
