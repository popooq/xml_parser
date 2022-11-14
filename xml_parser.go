package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type ContinentLogs struct {
	XMLName         xml.Name        `xml:"-"`
	Cryptogateways  Cryptogateways  `xml:"-"`
	FilterRules     FilterRules     `xml:"-"`
	AccessServerLog AccessServerLog `xml:"AccessServerLog"`
}

type Cryptogateways struct {
	XMLName xml.Name `xml:"-"`
	Cgw     []Cgw    `xml:"-"`
}

type Cgw struct {
	XMLName xml.Name `xml:"-"`
	Id      string   `xml:"-"`
	Cssid   string   `xml:"-"`
	Ip      string   `xml:"-"`
	Tz      string   `xml:"-"`
}

type FilterRules struct {
	XMLName xml.Name `xml:"-"`
	Rule    []Rule   `xml:"-"`
}

type Rule struct {
	XMLName xml.Name `xml:"-"`
	Id      string   `xml:"-"`
	Deleted string   `xml:"-"`
}
type AccessServerLog struct {
	XMLName xml.Name `xml:"AccessServerLog"`
	Text    string   `xml:",chardata"`
	Record  []Record `xml:"record"`
}

type Record struct {
	XMLName xml.Name `xml:"record"`
	Cgwid   string   `xml:"-"`
	DtLocal string   `xml:"dt_local,attr"`
	DtCgw   string   `xml:"-"`
	User    string   `xml:"user,attr"`
	Event   string   `xml:"event,attr"`
	Ip      string   `xml:"ip,attr"`
	Desc    string   `xml:"desc,attr"`
}

type RecordList []Record

func metricTime(start time.Time) {
	// функция Now() возвращает текущее время, а функция Sub возвращает разницу между двумя временными метками
	fmt.Println(time.Since(start))
}

func (e RecordList) Len() int {
	return len(e)
}

func (e RecordList) Less(i, j int) bool {
	if e[i].User != e[j].User {
		return e[i].User < e[j].User
	}
	return e[i].DtLocal < e[j].DtLocal
}

func (e RecordList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func main() {
	defer metricTime(time.Now())
	var path string
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &path)

	xmlFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Файл считан успешно")
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var Log ContinentLogs
	Kok := AccessServerLog{}
	xml.Unmarshal(byteValue, &Log)
	File, err := os.Create("log_month.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer File.Close()

	for i := 0; i < len(Log.AccessServerLog.Record); i++ {
		if Log.AccessServerLog.Record[i].Event == "Аутентификация АП" && Log.AccessServerLog.Record[i].User != "" || Log.AccessServerLog.Record[i].Event == "Клиент АП отключен" && Log.AccessServerLog.Record[i].User != "" {
			Kok.Record = append(Kok.Record, Log.AccessServerLog.Record[i])
		}
	}

	for i := 0; i < len(Kok.Record); i++ {
		if Kok.Record[i].Event == "Аутентификация АП" {
			desc := strings.Split(Kok.Record[i].Desc, ";")
			Kok.Record[i].Desc = desc[2]
		}
	}

	sort.Sort(RecordList(Kok.Record))

	data, _ := xml.MarshalIndent(Kok, "", "\t")
	File.Write(data)

	pathwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Лог сохранен в", pathwd)

	/**	fmt.Println("Press 'q' to quit")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		exit := scanner.Text()
		if exit == "q" {
			break
		} else {
			fmt.Println("Press 'q' to quit")
		}
	} **/
}
