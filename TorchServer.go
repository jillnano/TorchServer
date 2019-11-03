package main

import (
	"net/http"
	"fmt"
	// "database/sql"
	"time"
	"math/rand"
	"strconv"
	// "math"
	// "bytes"
	// "strings"
	"os"
	"os/exec"
	"io/ioutil"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"./ServerUtils"
)

var appid string
var secret string

func init() {
	fileObj, _ := os.OpenFile("config.cfg", os.O_RDONLY, 0644)
	defer fileObj.Close()
	content, _ := ioutil.ReadAll(fileObj)
	cfg, _ = JsonToMap(string(content))
	appid = cfg["appid"].(string)
	secret = cfg["secret"].(string)
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/addtorch", AddTorch)
	e.POST("/gettorch", GetTorch)
	e.GET("/test", ServerTestGet)  // ok
	e.GET("/session", LoginSession)
	e.GET("/login", LoginUser)

	// Start server
	// e.Start(":53249")
	e.Logger.Fatal(e.Start(":5959"))
	// e.Logger.Fatal(e.StartTLS(":5959", "cert.pem", "key.pem"))
	// fmt.Println("Start")
}

func LoginSession(c echo.Context) error {
	js_code := c.QueryParam("code")
	wxApi := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%+v&secret=%+v&js_code=%+v&grant_type=authorization_code", appid, secret, js_code)
	client := &http.Client{}
	client.Timeout = 5 * time.Second
	req, _ := http.NewRequest("GET", wxApi, nil)
	req.Header.Set("User-Agent", "wx")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)  // 读取Body
	content := string(body)
	fmt.Println(content)
	data, _ := JsonToMap(content)
	return c.String(http.StatusOK, data["session_key"].(string))
}

func LoginUser(c echo.Context) error {
	ed := c.QueryParam("ed")
	iv := c.QueryParam("iv")
	sk := c.QueryParam("sk")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tmpFile := "Temp/" + strconv.FormatInt(int64(r.Intn(99999)), 16) + ".txt"
	os.Remove(tmpFile)
	content := sk + "\r\n" + ed + "\r\n" + iv
	// fmt.Println(content)
	fileObj, _ := os.OpenFile(tmpFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fileObj.WriteString(content)
	fileObj.Close()
	cmd := exec.Command("python", "decode.py", tmpFile)
	buf, _ := cmd.Output()
	data, _ := JsonToMap(string(buf))
	openid := data["openId"].(string)
	fmt.Printf("output: %+v: %+v\n", tmpFile, openid)
	os.Remove(tmpFile)
	playlist := ServerUtils.FindUser(openid)
	// fmt.Printf("err: %v",err)
	return c.String(http.StatusOK, MapToJson(playlist))
}

// post
func AddTorch(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	content := string(body)
	data, _ := JsonToMap(content)
	fmt.Println(data)
	return c.String(http.StatusOK, "200")
}

func GetTorch(c echo.Context) error {
	playlist_id := c.QueryParam("playlist_id")
	musicList := ServerUtils.FindPlaylist(playlist_id)
	fmt.Println(musicList)
	return c.String(http.StatusOK, MapToJson(musicList))
}


func ServerTestGet(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func JsonToMap(jsonStr string) (map[string]interface{}, []map[string]interface{}) {
	var mapResult map[string]interface{}
	var mapListResult []map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	errList := json.Unmarshal([]byte(jsonStr), &mapListResult)
	if err != nil && errList != nil {
		fmt.Println(fmt.Sprintf("JsonToMap err: %v", err))
		fmt.Println(fmt.Sprintf("JsonToMap err: %v", errList))
		fmt.Println(fmt.Sprintf("JsonToMap err: %v", jsonStr))
		// return mapResult, mapListResult
	}
	return mapResult, mapListResult
}

func MapToJson(mapData interface{}) string {
	jsonStr, err := json.Marshal(mapData)
	if err != nil {
		fmt.Println(fmt.Sprintf("MapToJson err: %v", err))
		fmt.Println(fmt.Sprintf("MapToJson err: %v", mapData))
		return ""
	}
	return string(jsonStr)
}