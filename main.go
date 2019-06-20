package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)

func main() {
    logHandler := func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Content-Type", "application/javascript")

        resp, _ := http.Get("http://ip-api.com/json/" + strings.Split(req.RemoteAddr, ":")[0])
        defer resp.Body.Close()

        apiString := ""

        if resp.StatusCode == http.StatusOK {
            apiBytes, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                log.Fatal(err)
            }
            apiString = string(apiBytes)
        }

        loc, err := time.LoadLocation("Europe/Riga")
        if err != nil {
            log.Println(err)
        }

        data := &struct {
            UserAgent  string
            Method     string
            Host       string
            RemoteAddr string
            OtherData  string
            Time       string
        }{
            req.Header.Get("User-Agent"),
            req.Method,
            req.Host,
            req.RemoteAddr,
            apiString,
            time.Now().In(loc).Format("Mon Jan 2 15:04:05"),
        }

        intelligence, err := json.Marshal(data)
        if err != nil {
            log.Println(err)
        }

        f, err := os.OpenFile("visitors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        if err != nil {
            log.Println(err)
            return
        }

        if _, err = f.WriteString("," + string(intelligence)); err != nil {
            panic(err)
        }

        if err := f.Close(); err != nil {
            log.Println(err)
        }

        //if err := ioutil.WriteFile("visitors.log", []byte(intelligence), 0644); err != nil {
        //    log.Println(err)
        //}

        if _, err := fmt.Fprint(w, "let data = {\"_status\":\"ok\"}"); err != nil {
            log.Println(err)
            return
        }
    }

    getLogHandler := func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Content-Type", "application/json")

        read, err := ioutil.ReadFile("visitors.log")
        if err != nil {
            if _, err := fmt.Fprint(w,  err.Error()); err != nil {
                log.Println(err)
                return
            }
            log.Println(err)
            return
        }

        if _, err := fmt.Fprint(w, "[" + string(read) + "]"); err != nil {
            if _, err := fmt.Fprint(w, err.Error()); err != nil {
                log.Println(err)
                return
            }
            log.Println(err)
        }
    }

    statsHandler := func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")

        read, err := ioutil.ReadFile("index.html")
        if err != nil {
            log.Println(err)
        }

        if _, err := fmt.Fprint(w, string(read)); err != nil {
            if _, err := fmt.Fprint(w, err.Error()); err != nil {
                log.Println(err)
                return
            }
            log.Println(err)
        }
    }

    http.HandleFunc("/get", getLogHandler)
    http.HandleFunc("/stats", statsHandler)
    http.HandleFunc("/sc.js", logHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
