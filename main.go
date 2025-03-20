package main

type TimesData struct {
        Type         string
        WorkType     string
        StartTime    string
        EndTime      string
        Status       string
        IsYesterday  bool
        OrderNo      int
        DocNo        string
        TripTo       string
        TevenBegType string
        TevenEndType string
}

type DayData struct {
        Date      string
        Type      int
        IsWorkDay bool
        Times     []TimesData
}

type DonotJuan struct {
        Data  []DayData
        Code  int
        Msg   string
        Value int
}

func timeToNum(timeStr string) int {
        //"213800" -> 21*60+38
        if timeStr == "" || len(timeStr) != 6 {
                return 0
        }
        var hour, min int
        hour, _ = strconv.Atoi(timeStr[:2])
        min, _ = strconv.Atoi(timeStr[2:4])
        return hour*60 + min
}

func main() {
        var totalMinutes int
        var month [13]int
        var fileMonth int
        var temp int
        c := make(chan os.Signal)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

        dataBytes, err := os.ReadFile("./data.json")
        if err != nil {
                return
        }
        var juan DonotJuan
        err = json.Unmarshal(dataBytes, &juan)
        if err != nil {
                return
        }
        nowMonth := time.Now().Month()
        nowDay := time.Now().Day()

        for _, v := range juan.Data {
                dataMonth, _ := strconv.Atoi(v.Date[4:6])
                month[dataMonth]++
        }
        for k, v := range month {
                if v > temp {
                        temp = v
                        fileMonth = k
                }
        }

        for _, v := range juan.Data {
                dataMonth, _ := strconv.Atoi(v.Date[4:6])
                dataDay, _ := strconv.Atoi(v.Date[6:])
                if v.IsWorkDay && dataMonth == fileMonth {
                        if dataDay == nowDay && dataMonth == int(nowMonth) { //今天的不计算
                                continue
                        }
                        var actualStart, actualEnd int = 0, 0 //实际开始时间和结束时间
                        for _, vv := range v.Times {
                                if vv.StartTime == "000000" && vv.EndTime == "000000" { //全天请假/外出公干
                                        continue
                                }
                                if vv.StartTime != "" || vv.EndTime != "" {
                                        if actualStart > timeToNum(vv.StartTime) || actualStart == 0 {
                                                actualStart = timeToNum(vv.StartTime)
                                        }
                                        if actualEnd < timeToNum(vv.EndTime) || actualEnd == 0 {
                                                actualEnd = timeToNum(vv.EndTime)
                                        }
                                }
                        }
                        if actualEnd-actualStart > 0 {
                                totalMinutes += actualEnd - actualStart
                                totalMinutes -= 480 + 90
                        }
                }
        }

        fmt.Printf("剩余工时%d分钟\n", totalMinutes)
        if totalMinutes > 180 {
                fmt.Println("卷王求放过！")
        } else {
                fmt.Println("到点就走，继续保持！干的完就干，干不完就道歉！加班？不可能")
        }
        <-c
}
