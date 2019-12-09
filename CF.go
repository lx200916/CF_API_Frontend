package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
	"strings"
	"time"
)
import "github.com/gin-contrib/cors"
import _ "github.com/jinzhu/gorm/dialects/mysql"
type Mess struct {
	Name string
	Nick string
	Toprating int
	Lowrating int
	Lastone json.RawMessage
	Status string
	Num int
	Rating int
	Rank int
	LastFive []string
	Avg float32
}
type Person struct {
	ID        uint `gorm:"primary_key"`
	LastFive string
	Name string
	Nick string
	Toprating int
	Lowrating int
	Num int
	Rating int
	Rank int
	Avg float32
	LastR []byte

}
type Contestlist struct {
	ID int `gorm:"primary_key"`
	ContestName string


}
type Contest struct {
	ID        uint `gorm:"primary_key"`

	ContestId int
	ContestName string
	Handle string
	RatingUpdateTimeSeconds string
	Rating int



}
var names = [...]string{}
var lists = []string{}

var db, err = gorm.Open("mysql", "ACM:123456@/ACM?charset=utf8&parseTime=True&loc=Local")


func main() {
	//


fmt.Println(err)
	app:=gin.Default()


	app.LoadHTMLGlob("./templates/*")
	app.StaticFS("/static",http.Dir("./static"))
	app.GET("/", func(c *gin.Context) {
		c.HTML(200,"index.html",nil)
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH","POST","GET",},
		AllowHeaders:     []string{"Origin","X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	app.GET("/info", func(c *gin.Context) {
		var returns []Mess;
		var infos []Person
		var lastfive []string
		db.Order("avg desc").Find(&infos)
		for _, value := range infos {
			lastfive=strings.Split(value.LastFive,",")
			if(len(lastfive)>1){
				lastfive=lastfive[0:len(lastfive)-2]
			}else {
				lastfive=[]string{"0"}
			}

			returns = append(returns, Mess{LastFive:lastfive,Nick:value.Nick,Name:value.Name,Num:value.Num,Avg:value.Avg,Toprating:value.Toprating,Lowrating:value.Lowrating,Lastone:value.LastR,Rank:value.Rank,Rating:value.Rating})


		}



	//
	//	infos:=[]Mess{}
	//for key, value := range lists {
	//info, _ :=http.Get("http://codeforces.com/api/user.rating?handle="+value)
	//defer info.Body.Close()
	//jsons,_:=ioutil.ReadAll(info.Body)
	//status:=gjson.Get(string(jsons),"status")
	//	var toprating int;
	//	var lowrating=100000;
	//ratings:=gjson.Get(string(jsons),"result").Array()
	//var rating int
	//var lastone []byte
	//	if len(ratings)!=0 {
	//		lastone = []byte(ratings[len(ratings)-1].Raw)
	//		rating = int(gjson.Get(ratings[len(ratings)-1].Raw, "newRating").Int())
	//
	//	}else {lastone=[]byte("{\"contestName\":\"No Contest\"}")}
	//
	//
	//	for _, value := range ratings {
	//		ratings:=gjson.Get(value.String(),"newRating").Int()
	//		if int(ratings)>int(toprating) {
	//			//fmt.Println(rating)
	//			toprating=int(rating)
	//
	//		}
	//		if int(ratings)<lowrating {
	//			//fmt.Println(rating)
	//			lowrating=int(rating)
	//		}
	//
	//
	//
	//
	//	}
	//
	//
	//
	//	infos=append(infos,Mess{Name:names[key],Nick:value,Toprating:toprating,Lowrating:lowrating,Lastone:lastone,Status:status.Str,Num:len(ratings),Rating:rating})
	//
	//
	//	//fmt.Println(status,infos)
	//
	//}


	//json, _ := json.Marshal(infos)


		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		c.JSON(200,returns)
	})

	app.POST("/detail" ,func(c *gin.Context) {
		//ids:=c.PostFormArray("ID")
		fmt.Println()
		json, _ :=c.GetRawData()
		ids:=		gjson.Get(string(json),"ID").Array()

		fmt.Println(ids)
		contests:= []Contestlist{}
		db.Find(&contests)


		 infos:= make(map[string][]interface{})
		var contestinfo []map[string]string
		for _, value := range contests {
			var flag bool=false


			for _, id := range ids {
				var contest Contest
				db.Where(Contest{ContestId:value.ID,Handle:id.Str}).Attrs("Rating",0).FirstOrInit(&contest)


				if contest.Rating!=0 && !flag {
					flag=true
					contestinfo = append(contestinfo, map[string]string{"ID":strconv.Itoa( contest.ContestId),"Time":contest.RatingUpdateTimeSeconds,"Name":contest.ContestName})


				}

			}
			if flag{
			for _, id := range ids {
				var contest Contest
				db.Where(Contest{ContestId:value.ID,Handle:id.Str}).Attrs("Rating",0).FirstOrInit(&contest)
				var p Person
				db.Where(Person{Nick:id.Str}).First(&p)
				if contest.Rating!=0 {

					infos[p.Name] = append(infos[p.Name], contest.Rating)
				}
				if contest.Rating==0 {
					infos[p.Name] = append(infos[p.Name], "-")
				}

			}}

		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")

		c.JSON(200,map[string]interface{}{"infos":infos,"contests":contestinfo} )


	})

	_ = app.Run("0.0.0.0:8880")
}
