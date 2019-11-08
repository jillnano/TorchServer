package main

import (
	"./ServerUtils"
	"net/http"
	// "net/url"
	"strings"
	"time"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	// "encoding/json"
	"regexp"
	"flag"
	"os/exec"
	"path/filepath"
)

var playlist string
var user string
var bind bool

func init() {
	flag.StringVar(&playlist, "playlist", "", "playlist id")
	flag.StringVar(&user, "user", "", "user openid")
	flag.BoolVar(&bind, "bind", false, "bind")
}

func getPlayList(playListUrl string) (string, string, [][]string) {
	client := &http.Client{}
	client.Timeout = 5 * time.Second
	req, _ := http.NewRequest("GET", playListUrl, nil)
	req.Header.Set("authority", "music.163.com")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "nested-navigate")
	req.Header.Set("referer", "https://music.163.com/")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()  // 函数结束时关闭Body
	body, err := ioutil.ReadAll(resp.Body)  // 读取Body
	content := string(body)
	exp := regexp.MustCompile(`<li><a href="/song\?id=(\d+?)">(.+?)</a></li>`)
	playList := exp.FindAllStringSubmatch(content, -1)
	titleExp := regexp.MustCompile(`<meta property="og:title" content="(.+?)" />`)
	albumExp := regexp.MustCompile(`<meta property="og:image" content="(.+?)" />`)
	playlist_name := titleExp.FindStringSubmatch(content)
	album_url := albumExp.FindStringSubmatch(content)
	return playlist_name[1], album_url[1], playList
}

func FileExists(path string) bool {
	isExist := false
	_, err := os.Stat(path)
	if err == nil {
		isExist = true
	}
	if os.IsNotExist(err) {
		isExist = false
	}
	return isExist
}

func download(mUrl string, path string, name string, ext string) string {
	buf := make([]byte, 32*1024)
	var written int64
	client := &http.Client{}
	client.Timeout = 60 * time.Second
	req, _ := http.NewRequest("GET", mUrl, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if strings.HasSuffix(resp.Request.URL.String(), "404") {
		return ""
	}
	tmpFilePath := fmt.Sprintf("%+v/%+v.%+v.temp", path, name, ext)
	realFilePath := fmt.Sprintf("%+v/%+v.%+v", path, name, ext)
	if FileExists(tmpFilePath) {
		os.Remove(tmpFilePath)
	}
	if FileExists(realFilePath) {
		return realFilePath
	}
	file, err := os.Create(tmpFilePath)
	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		//没有错误了快使用 callback
		// fb(fsize, written)
	}
	if err != nil {
		os.Remove(tmpFilePath)
		return ""
	} else {
		os.Rename(tmpFilePath, realFilePath)
	}
	return realFilePath
}

func GetMusicFromUrl(url string) {
	var midList []string
	exp := regexp.MustCompile(`playlist\?id=(\d+)`)
	playlist_id := exp.FindStringSubmatch(url)[1]
	playlist_name, album_url, playList := getPlayList(url)
	// download(album_url, "album", playlist_id, "jpg")
	fmt.Println(playlist_name + ": " + playlist_id)
	for _, v := range playList {
		mUrl := fmt.Sprintf("http://music.163.com/song/media/outer/url?id=%+v.mp3", v[1])
		fn := download(mUrl, "music", v[1], "mp3")
		if strings.HasSuffix(fn, ".mp3") {
			fn, _ = filepath.Abs(fn)
			fmt.Println(v[2] + ": " + fn)
			midList = append(midList, v[1])
			ServerUtils.InsertMusic(v[1], fn, v[2])
		}
	}
	ServerUtils.UpsertPlayListMongo(playlist_id, playlist_name, album_url, midList)
}

func main() {
	flag.Parse()
	if bind && playlist != "" && user != "" {
		ServerUtils.SignUserMongo(user, playlist)
	} else if playlist != "" {
		url := "https://music.163.com/playlist?id=" + playlist
		fmt.Println(url)
		GetMusicFromUrl(url)
		cmd := exec.Command("python", "-u", "Encoder_main.py", playlist)
		stdout, err := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout
		if err != nil {
			fmt.Println(err)
		}
		if err = cmd.Start(); err != nil {
			fmt.Println(err)
		}
		// 从管道中实时获取输出并打印到终端
		for {
			tmp := make([]byte, 1024)
			n, err := stdout.Read(tmp)
			fmt.Print(string(tmp[0:n]))
			if err != nil {
				break
			}
		}
		if err = cmd.Wait(); err != nil {
			fmt.Println(err)
		}
		// buf, _ := cmd.Output()
		// fmt.Println(string(buf))
	}
}