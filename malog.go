package malog

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const url = "http://www.metal-archives.com"

type Response struct {
	Title string
	Type  string
	Name  string
	URL   string
	Date  string
}

var (
	bandsAdded   = make(map[string]*Response)
	bandsUpdated = make(map[string]*Response)

	labelsAdded   = make(map[string]*Response)
	labelsUpdated = make(map[string]*Response)

	artistsAdded   = make(map[string]*Response)
	artistsUpdated = make(map[string]*Response)

	reviewsAdded = make(map[string]*Response)
)

func getData(doc *goquery.Document, cache map[string]*Response, id string) map[string]*Response {
	bands := make(map[string]*Response)
	doc.Find(id + " tr").Each(func(i int, s *goquery.Selection) {
		td := s.Find("td")
		url, _ := td.Eq(0).Find("a").First().Attr("href")

		var info string
		if td.Length() == 3 {
			url, _ = td.Eq(1).Find("a").First().Attr("href")
			info += `"` + td.Eq(1).Text() + `", for `
		}
		info += td.Eq(0).Find("a").First().Text()

		b := &Response{
			Name: info,
			Date: td.Last().Text(),
			URL:  url,
		}
		bands[b.URL] = b
	})

	for url1, b := range bands {
		for url2 := range cache {
			if url1 == url2 {
				delete(bands, url2)
			}
		}
		cache[url1] = b
	}

	return bands
}

func job(title string, r chan Response, doc *goquery.Document, cache map[string]*Response, divID, status string) {
	data := getData(doc, cache, divID)
	for _, b := range data {
		res := Response{
			Title: strings.Title(title),
			Name:  b.Name,
			Type:  status,
			URL:   b.URL,
		}
		r <- res
	}
}

func fetch(cb func(chan Response, chan error)) (r chan Response, er chan error) {
	r = make(chan Response)
	er = make(chan error)
	go cb(r, er)
	return
}

func Fetch() (chan Response, chan error) {
	return fetch(func(r chan Response, er chan error) {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			er <- err
			return
		}
		job("band", r, doc, bandsAdded, "#additionBands", "added")
		job("band", r, doc, bandsUpdated, "#updatedBands", "updated")
		job("review", r, doc, reviewsAdded, "#lastReviews", "added")
		go fetchLabels(r, er, labelsAdded, url+"/index/latest-labels", "added")
		go fetchLabels(r, er, labelsUpdated, url+"/index/latest-labels/by/modified", "updated")
		go fetchArtists(r, er, artistsAdded, url+"/index/latest-artists", "added")
		go fetchArtists(r, er, artistsUpdated, url+"/index/latest-artists/by/modified", "updated")
	})
}

func FetchBands() (chan Response, chan error) {
	return fetch(func(r chan Response, er chan error) {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			er <- err
			return
		}
		job("band", r, doc, bandsAdded, "#additionBands", "added")
		job("band", r, doc, bandsUpdated, "#updatedBands", "updated")
	})
}

func fetchLabels(r chan Response, er chan error, cache map[string]*Response, url, status string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		er <- err
		return
	}
	job("label", r, doc, cache, "#additionLabels", status)
}

func FetchLabels() (chan Response, chan error) {
	return fetch(func(r chan Response, er chan error) {
		go fetchLabels(r, er, labelsAdded, url+"/index/latest-labels", "added")
		go fetchLabels(r, er, labelsUpdated, url+"/index/latest-labels/by/modified", "updated")
	})
}

func fetchArtists(r chan Response, er chan error, cache map[string]*Response, url, status string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		er <- err
		return
	}
	job("artist", r, doc, cache, "#additionArtists", status)
}

func FetchArtists() (chan Response, chan error) {
	return fetch(func(r chan Response, er chan error) {
		go fetchArtists(r, er, artistsAdded, url+"/index/latest-artists", "added")
		go fetchArtists(r, er, artistsUpdated, url+"/index/latest-artists/by/modified", "updated")
	})
}
