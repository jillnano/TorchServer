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

func SignUserMongo(openid string, playlist []map[string]string) {
	c := database.C("user")
	c.Insert(bson.M{"openid": openid, "playlist": playlist})
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
	var result map[string]interface{}
	c := database.C("user")
	c.Find(bson.M{"openid": openid}).Select(bson.M{"playlist": 1, "_id": 0}).One(&result)
	return result["playlist"]
}

func FindPlaylist(playlist_id string) interface{} {
	var result map[string]interface{}
	c := database.C("playlist")
	c.Find(bson.M{"playlist_id": playlist_id}).Select(bson.M{"musicList": 1, "_id": 0}).One(&result)
	return result["musicList"]
}

func main() {
	var playlist []map[string]string
	playlist = append(playlist, map[string]string{"pl_id": "112545205", "pl_name": "自强番茄喜欢的音乐", "pl_url": "http://p1.music.126.net/UCuTwLDEDVWSeMIG2DBAoQ==/19094118928164389.jpg"})
	// playlist["112545205"] = "自强番茄喜欢的音乐"
	SignUserMongo("oZxug4uNQw8Xnxe8k6FJsfGvdqlQ", playlist)
}

// func GetConfigSettings() map[string]interface{} {
// 	config := make(map[string]interface{})
// 	var result []map[string]interface{}
// 	c := settingsDatabase.C("config")
// 	c.Find(bson.M{}).Select(bson.M{"_id": 0}).All(&result)
// 	for _, i := range result {
// 		config[i["name"].(string)] = i
// 	}
// 	return config
// }
