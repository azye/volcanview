package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/gocolly/colly/v2"
)

type Volcano struct {
	number       int
	name         string
	description  string
	lastEruption string
	elevation    string
	latitude     float32
	longitude    float32
}

func newVolcano(number int, name, lastEruption, elevation string, latitude, longitude float32) *Volcano {
	v := new(Volcano)

	v.number = number
	v.name = name
	v.lastEruption = lastEruption
	v.elevation = elevation
	v.latitude = latitude
	v.longitude = longitude
	return v
}

func (v *Volcano) toString() string {
	return fmt.Sprintf("%d: %s %f %f", v.number, v.name, v.latitude, v.longitude)
}

func Split(r rune) bool {
	return unicode.IsUpper(r)
}

func main() {
	c := colly.NewCollector(
		colly.CacheDir("./volcano_cache"),
		colly.MaxBodySize(0),
		// colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
	)

	c.SetRequestTimeout(time.Minute * 10)

	v := new(Volcano)

	VOLCANO_NUMBER := 332010
	// VOLCANO_NUMBERS := []int{342080, 332010}

	v.number = VOLCANO_NUMBER

	c.OnHTML("div.tabbed-content:nth-child(6) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > p:nth-child(1)", func(e *colly.HTMLElement) {
		v.description = e.Text
	})

	// get volcano basic data
	c.OnHTML(".DivTable35 > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(2)", func(e *colly.HTMLElement) {
		// todo: gets value only. maybe get label values as validation in future
		props := strings.Split(strings.TrimSpace(e.Text), "\n")
		for i, t := range props {
			fmt.Println(i, strings.TrimSpace(t))
		}
	})

	// get the volcano types
	c.OnHTML(".DivTable35 > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(4) > td:nth-child(1)", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)

		a := strings.FieldsFunc(strings.TrimSpace(e.Text), Split)

		fmt.Println(a)

		// fmt.Println(strings.Split(e.Text, " "))
		// b, _ := json.Marshal(strings.Split(e.Text, " "))
		// fmt.Printf("%v", string(b))
		// props := strings.Split(strings.TrimSpace(e.Text), "<br>")
		// for i, t := range props {
		// 	fmt.Println(i, strings.TrimSpace(t))
		// }
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(fmt.Sprintf("https://volcano.si.edu/volcano.cfm?vn=%d&vtab=GeneralInfo", VOLCANO_NUMBER))

	c.Wait()
	fmt.Printf("%+v", v)

}
