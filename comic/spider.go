package comic

import (
	"encoding/json"
	"net/http"
	"otakumate/comic/comicdao"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

type result struct {
	Code string        `json:"code"`
	Data []interface{} `json:"data"`
}

const (
	hhcomic = "http://hhzzee.com/"
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
	//解析参数
	vars := mux.Vars(r)
	popType := vars["type"]
	pageNum := vars["pageNum"]

	//创建Dao
	dao := comicdao.NewComicDao()
	//查询数据库是否存在
	res := dao.QueryComic(popType + pageNum)
	//创建slice
	list := make([]interface{}, 0)

	//如果数据库中存在数据直接返回库中的内容
	if len(res) == 20 {
		for _, item := range res {
			item := Comic{Title: item.Name, ComicID: strconv.FormatInt(item.ID, 10), Cover: item.Cover}
			list = append(list, item)
		}
	} else { //否则去爬接口
		list = fetchNewestComic(popType, pageNum)
	}
	creatResJson(&w, list)
}

func fetchNewestComic(popType string, pageNum string) []interface{} {
	list := make([]interface{}, 0)
	dao := comicdao.NewComicDao()
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
		panic(err)
	} else {
		doc.Find("body").Find("li").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Find("a").Attr("href")
			path := strings.Split(link, "/")
			idString := path[len(path)-1]

			title, _ := s.Find("img").Attr("alt")
			cover, _ := s.Find("img").Attr("src")

			titles := strings.Split(title, " ")
			item := Comic{Title: titles[0], ComicID: idString, Cover: cover}
			list = append(list, item)

			idint, _ := strconv.ParseInt(idString, 10, 64)
			dao.InsertComicData(&comicdao.ComicTable{ID: idint, Name: titles[0], Cover: cover, TP: (popType + pageNum)})
		})
	}
	return list
}

///解析Vols
func ComicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	comicID := vars["comicID"]

	cid, _ := strconv.ParseInt(comicID, 10, 64)

	list := make([]interface{}, 0)
	dao := comicdao.NewComicDao()

	res := dao.QueryVols(cid)

	if len(res) == 20 {
		for _, item := range res {
			vol := Volume{Tittle: item.Name, VolID: item.Vols}
			list = append(list, vol)
		}
	} else {
		list = fetchNewestVols(comicID)
	}
	creatResJson(&w, list)
}

func fetchNewestVols(comicID string) []interface{} {
	target := hhcomic + "comic/" + comicID
	cid, _ := strconv.ParseInt(comicID, 10, 64)
	dao := comicdao.NewComicDao()
	doc, _ := goquery.NewDocument(target)
	list := make([]interface{}, 0)

	doc.Find(".list_href").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		volstrings := strings.Split(link, "/")
		volstr := volstrings[len(volstrings)-2]
		text := s.Text()
		vol := Volume{Tittle: text, VolID: volstr}
		dao.InsertVolsData(&comicdao.VolsTable{ComicID: cid, Name: text, Vols: volstr})
		list = append(list, vol)
	})
	return list
}

//返回漫画列表
func VolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volID := vars["volID"]
	volss := strings.Split(volID, "_")
	vid, _ := strconv.ParseInt(volss[len(volss)-1], 10, 64)

	list := make([]interface{}, 0)
	dao := comicdao.NewComicDao()
	res := dao.QueryVolDetail(vid)
	if len(res) == 20 {
		for _, item := range res {
			list = append(list, item.Image)
		}
	} else {
		list = fetchNewestVolDetail(volID)
	}
	creatResJson(&w, list)
}

func fetchNewestVolDetail(volID string) []interface{} {
	volss := strings.Split(volID, "_")
	vid, _ := strconv.ParseInt(volss[len(volss)-1], 10, 64)

	target := hhcomic + "vols/" + volID

	doc, _ := goquery.NewDocument(target)
	list := make([]interface{}, 0)

	dao := comicdao.NewComicDao()

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
					image := imgURLPrefix + item
					dao.InsertVolDetailData(&comicdao.VolTable{VolID: vid, Image: image})
					list = append(list, image)

				}
			}

		}
	})
	return list
}

func UpdateWholeData() {
	types := [3]string{"daily"}
	println("start to check " + time.Now().Format("2006-01-02 15:04:05"))
	for _, item := range types {
		for i := 1; i < 6; i++ {
			list := fetchNewestComic(item, strconv.Itoa(i))
			comics := make([]Comic, len(list))
			for i, item := range list {
				comics[i] = item.(Comic)
			}

			for _, item := range comics {
				println("check comic %v", item.Title)

				list := fetchNewestVols(item.ComicID)
				vols := make([]Volume, len(list))
				for i, item := range list {
					vols[i] = item.(Volume)
				}

				for _, item := range vols {
					println("check vols %v", item.Tittle)
					fetchNewestVolDetail(item.VolID)
				}

			}

		}
	}
	println("complete %v", time.Now().Format("2006-01-02 15:04:05"))

}
