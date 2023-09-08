package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"github.com/guptarohit/asciigraph"
)

func main() {
	start := time.Now()

	stockSymbol := flag.String("stock", "MSFT", "a string")
	flag.Parse()
	getHistoricData(*stockSymbol)
	fmt.Printf("Took %s seconds\n", time.Since(start))

}

func getHistoricData(stock string) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 13_3_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672.93 Safari/537.36"),
		colly.AllowedDomains("finance.yahoo.com"),
		colly.MaxBodySize(0),
		colly.AllowURLRevisit(),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		dataSlice := []float64{}
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			var cnt int = 0
			el.ForEach("td", func(_ int, el *colly.HTMLElement) {
				cnt += 1
				if cnt == 6 {
					if s, err := strconv.ParseFloat(el.Text, 32); err == nil {
						dataSlice = append(dataSlice, s)
					}
				}

			})
		})
		for i, j := 0, len(dataSlice)-1; i < j; i, j = i+1, j-1 {
			dataSlice[i], dataSlice[j] = dataSlice[j], dataSlice[i]
		}
		graph := asciigraph.Plot(dataSlice)
		fmt.Println(graph)
	})

	c.Visit("https://finance.yahoo.com/quote/" + stock + "/history/")

	c.Wait()
}
