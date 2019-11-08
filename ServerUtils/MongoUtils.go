package ServerUtils

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"sync"
	"fmt"
	"github.com/bogem/id3v2"
	"strings"
	// "unicode"
	// "study/conf"
)

var once_mongo sync.Once
var session *mgo.Session
var database *mgo.Database

func init() {
	var err error
	 
	dialInfo := &mgo.DialInfo {
		Addrs: []string{"localhost:27017"},
		// Addrs: []string{"localhost:27017"},
		Database: "admin",
		Username: "torch",
		Password: "torch@mongo",
		Direct: false,
		Timeout: time.Second * time.Duration(30),
		// Session.SetPoolLimit
		PoolLimit: 1024,
	}
	//创建一个维护套接字池的session
	session, err = mgo.DialWithInfo(dialInfo)

	if err != nil {
		fmt.Println(err.Error())
	}
	// session.SetMode(mgo.Monotonic, true)
	//使用指定数据库
	dbName := "torch"
	database = session.DB(dbName)

}

func SignUserMongo(openid string, playlist_id string) {
	var result map[string]interface{}
	var playlist map[string]interface{}
	var pl = make(map[string]interface{})
	b := database.C("playlist")
	b.Find(bson.M{"playlist_id": playlist_id}).Select(bson.M{"_id": 0, "playlist_id": 1, "playlist_name": 1, "album_url": 1}).One(&playlist)

	c := database.C("user")
	err := c.Find(bson.M{"openid": openid}).Select(bson.M{"_id": 0}).One(&result)
	if err == nil {
		pl = result["playlist"].(map[string]interface{})
	}
	pl[playlist["playlist_id"].(string)] = playlist
	c.Upsert(bson.M{"openid": openid}, bson.M{"openid": openid, "playlist": pl})
}

func UpsertPlayListMongo(playlist_id string, playlist_name string, album_url string, midList []string) {
	c := database.C("playlist")
	c.Upsert(bson.M{"playlist_id": playlist_id}, bson.M{"playlist_id": playlist_id, "playlist_name": playlist_name, "album_url": album_url, "musicList": midList})
}

func InsertMusic(mid, filename, filetitle string) {
	c := database.C("music")
	tag, err := id3v2.Open(filename, id3v2.Options{Parse: true})
	if err != nil {
		return
	}
	// f := func(c rune) bool {
	// 	return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	// }
	title := strings.Replace(strings.Trim(tag.Title(), " "), "\u0000", "", -1)
	artist := strings.Replace(strings.Trim(tag.Artist(), " "), "\u0000", "", -1)
	album := strings.Replace(strings.Trim(tag.Album(), " "), "\u0000", "", -1)
	if title == "" {
		title = strings.Replace(strings.Trim(filetitle, " "), "\u0000", "", -1)
	}
	c.Insert(bson.M{"mid": mid, "title": title, "artist": artist, "album": album, "filename": filename})
}

func FindUser(openid string) interface{} {
	var playlist []interface{}
	var result map[string]interface{}
	c := database.C("user")
	c.Find(bson.M{"openid": openid}).Select(bson.M{"playlist": 1, "_id": 0}).One(&result)
	for _, v := range result["playlist"].(map[string]interface{}) {
		playlist = append(playlist, v)
	}
	return playlist
}

func FindPlaylist(playlist_id string) []interface{} {
	var result map[string]interface{}
	c := database.C("playlist")
	c.Find(bson.M{"playlist_id": playlist_id}).Select(bson.M{"musicList": 1, "_id": 0}).One(&result)
	musicList := result["musicList"].([]interface{})
	var ret []interface{}
	for _, mid := range musicList {
		v := FindMusic(mid.(string))
		ret = append(ret, v)
	}
	return ret
}

func FindMusic(mid string) interface{} {
	var result map[string]interface{}
	c := database.C("music")
	c.Find(bson.M{"mid": mid}).Select(bson.M{"_id": 0, "title": 1, "mid": 1, "encode_1": 1, "encode_2": 1}).One(&result)
	return result
}
