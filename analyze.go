package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var top int

type Site struct {
	Name   string
	Domain string
	Rank   int
}
type Extension struct {
	Name  string
	Count int
}
type ByRank []*Extension

func (e ByRank) Len() int      { return len(e) }
func (e ByRank) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e ByRank) Less(i, j int) bool {
	return e[i].Count < e[j].Count
}

var sites []Site
var extensionMap map[string]int = map[string]int{}

func read() {
	file, _ := os.Open("top1m.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	k := 1

	for scanner.Scan() && k < top {
		parts := strings.Split(scanner.Text(), "\t")
		domain := strings.Split(parts[1], "/")[0]
		sites = append(sites, Site{parts[1], domain, k})
		k += 1

		domain_parts := strings.Split(domain, ".")
		ext := domain_parts[len(domain_parts)-1]
		if val, ok := extensionMap[ext]; ok {
			extensionMap[ext] = val + 1
		} else {
			extensionMap[ext] = 1
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func measureTimes() {
	for _, site := range sites {
		t0 := time.Now()
		resp, err := http.Get("https://" + site.Domain)
		if err != nil {
			fmt.Println(">>>>>>>>>", err)
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		t1 := time.Now()
		fmt.Printf("The call took %v to run and %v big. %s\n", t1.Sub(t0), len(body), site.Domain)

	}
}

func countExtensions() {
	// create list of extensions
	extensions := []*Extension{}
	for k, v := range extensionMap {
		extensions = append(extensions, &Extension{k, v})
	}

	sort.Sort(sort.Reverse(ByRank(extensions)))

	for e := range extensions[0:int(math.Min(float64(len(extensions)), float64(50)))] {
		fmt.Printf("%d.\t%s\t%d\t%.2f\n", e+1, extensions[e].Name, extensions[e].Count, float32(extensions[e].Count)/float32(top)*100.0)
	}
}

func main() {
	flag.IntVar(&top, "top", 1000, "top number of sites to count")
	flag.Parse()

	read()

	measureTimes()
}
