package comic

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type result struct {
	Code string        `json:"code"`
	Data []interface{} `json:"data"`
}

const (
	hhcomic = "http://3gmanhua.com/"
)

func creatResJson(w *http.ResponseWriter, list []interface{}) {
	res := result{}
	if len(list) != 0 {
		res.Code = "200"
		res.Data = list
	} else {
		res.Code = "101"
		res.Data = nil
	}
	js, jserr := json.Marshal(res)
	if jserr != nil {
	}
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(js)
}

///处理top类型的请求
func TopHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	popType := vars["type"]
	pageNum := vars["pageNum"]
	s := ""
	switch popType {
	case "daily":
		s = "a-"
		break
	case "week":
		s = "b-"
		break
	case "pop":
		s = "c-"
		break
	default:
		s = "a-"
		break
	}
	target := hhcomic + "top/" + s + pageNum + ".htm"
	doc, err := goquery.NewDocument(target)
	if err != nil {
		// fmt.Println(err)
		return
	}
	list := make([]interface{}, 0)
	doc.Find("body").Find("li").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a").Attr("href")
		title, _ := s.Find("img").Attr("alt")
		cover, _ := s.Find("img").Attr("src")
		item := Comic{Title: title, ComicID: link, Cover: cover}
		list = append(list, item)
	})
	creatResJson(&w, list)
}

///解析Vols
func ComicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	comicID := vars["comicID"]

	target := hhcomic + "comic/" + comicID

	doc, err := goquery.NewDocument(target)
	if err != nil {
		///
		return
	}
	list := make([]interface{}, 0)

	doc.Find(".list_href").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		text := s.Text()
		vol := Volume{Tittle: text, VolID: link}
		list = append(list, vol)
	})
	creatResJson(&w, list)
}

//返回漫画列表
func VolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volID := vars["volID"]
	target := hhcomic + "vols/" + volID
	doc, err := goquery.NewDocument(target)
	if err != nil {
	}
	list := make([]interface{}, 0)
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		var texts = s.Text()
		if len(texts) != 0 {
			params := strings.Split(texts, "\"")
			if len(params) >= 4 {
				var sfiles = params[1]
				var spath = params[3]
				var xstring = sfiles[len(sfiles)-1:]
				var xs = "abcdefghijklmnopqrstuvwxyz"
				var xi = strings.Index(xs, xstring) + 1
				var skstring = sfiles[len(sfiles)-xi-12 : len(sfiles)-xi-1]
				sfiles = sfiles[0 : len(sfiles)-xi-12]
				var kString = skstring[0 : len(skstring)-1]
				var fString = skstring[len(skstring)-1:]

				for i := 0; i < len(kString); i++ {
					old := kString[i : i+1]
					new := strconv.Itoa(i) //直接 string（i） 并不会转成功
					sfiles = strings.Replace(sfiles, old, new, -1)
				}

				var subs = strings.Split(sfiles, fString)
				var result = ""

				for _, item := range subs {
					b, _ := strconv.Atoi(item)
					result = result + string(rune(b))
				}
				var imgURLPrefix = ""

				if len(spath) >= 2 {
					imgURLPrefix = "http://3g.1112223333.com:9393/dm" + spath
				} else {
					imgURLPrefix = "http://3g.1112223333.com:9393/dm0" + spath
				}
				for _, item := range strings.Split(result, "|") {
					list = append(list, (imgURLPrefix + item))
				}
			}

		}
	})
	creatResJson(&w, list)
}
