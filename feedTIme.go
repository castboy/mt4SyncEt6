package mt4SyncEt6

import (
	"github.com/go-xorm/xorm"
	"github.com/juju/errors"
	"strconv"
	"strings"
	"time"
)

var AM1 = "01:00"

func GetXROMProLocal() (engine *xorm.Engine) {
	engine, _ = NewET6EngineXorm()
	return
}

//Get the target Sources
//这周末的是所有的FX，所有Metal，所有Energy，所有Crypto，还有US30的，USA500的，NAS100的，JPN225的
func GetSourceByCatagory(sourceType SourceType) []Source {
	//Get Engine
	engine := GetXROMProLocal()
	//select source By category base on requests
	sourceSet := make([]Source, 0)
	err := engine.Table("source").Where("source_type=?", sourceType).Find(&sourceSet)
	if err != nil {
		return nil
	}
	//return the data
	return sourceSet
}
func GetSourceBySourceName(sourceName string) (source *Source) {
	//Get Engine
	source = new(Source)
	engine := GetXROMProLocal()
	//select source By category base on requests
	_, _ = engine.Table("source").Where("source=?", sourceName).NoAutoCondition(true).Get(source)
	//return the data
	return
}

//Get the fixed symbol sessions(Trade  quote)
func GetSessionBySource(sourceID int) (sessions []Session) {
	//Get Engine
	engine := GetXROMProLocal()
	engine.Table("session").Where("source_id=?", sourceID).Find(&sessions)
	return
}

//Modify the time by day-lighting time
//Input [sourceID]sourceSessions
//timeSpan "22:00-24:00"

func InsertSession(sess *Session) error {
	engine := GetEngine()
	_, err := engine.Table(Session{}).Omit("id").InsertOne(sess)
	return err
}

func UpdateTimeSpanByID(sess *Session) error {
	engine := GetEngine()
	_, err := engine.Table(Session{}).Cols("time_span").Where("id=?", sess.ID).Update(sess)
	return err
}

func DeRepeate(sessionMos, sessionIns []Session) (newSessionMos []Session, newSessionIns []Session) {
	//Copy data
	newSessionMos = append(newSessionMos, sessionMos...)    //22:05-23:00 00:00-1:00
	mediaSessionIns := append(newSessionIns, sessionIns...) //23:00-24:00

	for i, sessionIns := range sessionIns {
		timeIns := strings.Split(sessionIns.TimeSpan, "-")
		for k, sessionMod := range sessionMos {
			// Conditions
			if sessionIns.Weekday != sessionMod.Weekday || sessionIns.Type != sessionMod.Type {
				continue
			}
			//Operate
			timeMod := strings.Split(sessionMod.TimeSpan, "-")
			// merge sessionIns to sessionMod
			if timeIns[0] <= timeMod[1] {
				//Merge
				sessionMod.TimeSpan = timeMod[0] + "-" + timeIns[1]
				//kick the sessionIns from sessionsIns
				newSessionMos[k] = sessionMod
				mediaSessionIns[i].TimeSpan = ""
			}
		}
	}
	//Operation kick
	for _, sessionIn := range mediaSessionIns {
		if sessionIn.TimeSpan == "" {
			continue
		}
		newSessionIns = append(newSessionIns, sessionIn)
	}
	return
}

//Modify the time by day-lighting time
func modifyTheSession(sessions []Session, hour int, feedType FeedTimeType) (modSessions, extraSession []Session, err error) {
	// Extra container
	for _, v := range sessions {
		var normalSpan, preSpanInF string
		var weekday time.Weekday
		if feedType == EarlierThanCurrent {
			weekday = (v.Weekday + time.Weekday(EarlierDay) + WeekLength) % WeekLength
			normalSpan, preSpanInF, err = earlyTimeConv(v.TimeSpan, hour)
		}
		if feedType == LaterThanCurrent {
			weekday = (v.Weekday + WeekLength + time.Weekday(LaterDay)) % WeekLength
			normalSpan, preSpanInF, err = laterTimeConv(v.TimeSpan, hour)
		}
		if err != nil {
			return nil, nil, err
		}
		//normalSpan!=nil include -1 hour within same day or exceed one day
		if normalSpan != "" {
			v.TimeSpan = normalSpan
			modSessions = append(modSessions, v)
			//-1hour, and within previous day
			if preSpanInF != "" {
				sessionTemp := Session{
					SourceID: v.SourceID,
					Type:     v.Type,
					Weekday:  weekday,
					TimeSpan: preSpanInF,
				}
				extraSession = append(extraSession, sessionTemp)
			}
		}
		//If origin span is "xx:xx-hour:00", need update but not insert
		if normalSpan == "" {
			sessionTemp := Session{
				SourceID: v.SourceID,
				Type:     v.Type,
				Weekday:  weekday,
				TimeSpan: preSpanInF,
			}
			sessionTemp.ID = v.ID
			modSessions = append(modSessions, v)
		}
	}
	//remove the crossing time_span, for examplp "22:00-23:30" and "23:00-24:00",which can be combined to "22:00-24:00
	modSessions, extraSession = deRepeate(modSessions, extraSession, feedType)
	return
}

//Remove and merge the crossing time_span
func deRepeate(sessionMos, sessionIns []Session, feedType FeedTimeType) (newSessionMos []Session, newSessionIns []Session) {
	//Copy data
	newSessionMos = append(newSessionMos, sessionMos...)    //22:05-23:00 00:00-1:00
	mediaSessionIns := append(newSessionIns, sessionIns...) //23:00-24:00
	//Prepare
	for i, sessionIns := range sessionIns {
		timeIns := strings.Split(sessionIns.TimeSpan, "-")
		for k, sessionMod := range sessionMos {
			// Conditions
			if sessionIns.Weekday != sessionMod.Weekday || sessionIns.Type != sessionMod.Type {
				continue
			}
			//Operate
			timeMod := strings.Split(sessionMod.TimeSpan, "-")
			// merge sessionIns to sessionMod
			if feedType == EarlierThanCurrent {
				if timeIns[0] <= timeMod[1] {
					//Merge
					sessionMod.TimeSpan = timeMod[0] + "-" + timeIns[1]
					//kick the sessionIns from sessionsIns
					newSessionMos[k] = sessionMod
					mediaSessionIns[i].TimeSpan = ""
				}
			}
			if feedType == LaterThanCurrent {
				if timeIns[1] >= timeMod[0] {
					//Merge
					sessionMod.TimeSpan = timeIns[0] + "-" + timeMod[1]
					//kick the sessionIns from sessionsIns
					newSessionMos[k] = sessionMod
					mediaSessionIns[i].TimeSpan = ""
				}
			}

		}
	}
	//Operation kick
	for _, sessionIn := range mediaSessionIns {
		if sessionIn.TimeSpan == "" {
			continue
		}
		newSessionIns = append(newSessionIns, sessionIn)
	}
	return
}

// convert time
func earlyTimeConv(span string, hour int) (normalSpan string, sessionInPreDay string, err error) {
	//Check
	if span == "" {
		return "", "", nil
	}
	if hour == 0 {
		return span, "", nil
	}
	//Prepare
	var AM1 string
	AM1 = strconv.Itoa(24-hour) + ":00"
	if hour < 10 {
		AM1 = "0" + strconv.Itoa(24-hour) + ":00"
	}

	timeSlices := strings.Split(span, "-") //2:35-4:25
	start := timeSlices[0]                 //2:35
	end := timeSlices[1]                   //4:25
	startHourBit := strings.Split(timeSlices[0], ":")[0]
	startMinBit := strings.Split(timeSlices[0], ":")[1]
	endHourBit := strings.Split(timeSlices[1], ":")[0]
	endMinBit := strings.Split(timeSlices[1], ":")[1]
	hourStart, err := strconv.Atoi(startHourBit)
	if err != nil {
		return "", "", errors.Errorf("strconv.Atoi:%+v", err)
	}
	hourEnd, err := strconv.Atoi(endHourBit)
	if err != nil {
		return "", "", errors.Errorf("strconv.Atoi:%+v", err)
	}
	//Process
	if start >= AM1 { //2:35-4:25 1:00-2:00
		hourStart = hourStart - hour //Backward one hour
		hourEnd = hourEnd - hour     //Backward one hour
		normalStartHour := timeFormatToStr(int(hourStart))
		normalEndHour := timeFormatToStr(int(hourEnd))
		normalStartTime := normalStartHour + ":" + startMinBit
		normalEndTime := normalEndHour + ":" + endMinBit
		normalSpan = normalStartTime + "-" + normalEndTime
		sessionInPreDay = ""
	} else if start < AM1 && end > AM1 { //00:35-4:25
		hourStart = hourStart + Daylength - 1 //Backward one hour
		hourEnd = hourEnd - 1                 //Backward one hour
		normalStartHour := timeFormatToStr(int(hourStart))
		normalEndHour := timeFormatToStr(int(hourEnd))
		normalStartTime := "00" + ":" + "00"
		normalEndTime := normalEndHour + ":" + endMinBit
		normalSpan = normalStartTime + "-" + normalEndTime
		sessionInPreDay = normalStartHour + ":" + startMinBit + "-24:00"
	} else if end <= AM1 { //00:35-01:00
		hourStart = hourStart + Daylength - 1 //Backward one hour
		hourEnd = hourEnd + Daylength - 1     //Backward one hour
		normalStartHour := timeFormatToStr(int(hourStart))
		normalEndHour := timeFormatToStr(int(hourEnd))
		normalStartTime := normalStartHour + ":" + startMinBit
		normalEndTime := normalEndHour + ":" + endMinBit
		normalSpan = ""
		sessionInPreDay = normalStartTime + "-" + normalEndTime
	}
	return
}

func laterTimeConv(span string, hour int) (normalSpan string, sessionInNextDay string, err error) {
	//Check
	if span == "" {
		return "", "", nil
	}
	if hour == 0 {
		return span, "", nil
	}
	//Prepare
	AM1 := strconv.Itoa(24-hour) + ":00"
	if 24-hour < 10 {
		AM1 = "0" + strconv.Itoa(24-hour) + ":00"
	}
	timeSlices := strings.Split(span, "-") //2:35-4:25
	start := timeSlices[0]                 //2:35
	end := timeSlices[1]                   //4:25
	startHourBit := strings.Split(timeSlices[0], ":")[0]
	startMinBit := strings.Split(timeSlices[0], ":")[1]
	endHourBit := strings.Split(timeSlices[1], ":")[0]
	endMinBit := strings.Split(timeSlices[1], ":")[1]
	hourStart, err := strconv.Atoi(startHourBit)
	if err != nil {
		return "", "", errors.Errorf("strconv.Atoi:%+v", err)
	}
	hourEnd, err := strconv.Atoi(endHourBit)
	if err != nil {
		return "", "", errors.Errorf("strconv.Atoi:%+v", err)
	}
	//Process
	if end <= AM1 { //22:00-22:30 22:00-23:00
		hourStart = hourStart + hour //Forward one hour
		hourEnd = hourEnd + hour     //Forward one hour
		normalStartHour := timeFormatToStr(int(hourStart))
		normalEndHour := timeFormatToStr(int(hourEnd))
		normalStartTime := normalStartHour + ":" + startMinBit
		normalEndTime := normalEndHour + ":" + endMinBit
		normalSpan = normalStartTime + "-" + normalEndTime
		sessionInNextDay = ""
	} else if start < AM1 && end > AM1 { //22:05-24:00
		hourStart = hourStart + hour                           //23
		hourEnd = hourEnd + hour - Daylength                   //1
		normalStartHour := timeFormatToStr(int(hourStart))     //23
		normalEndHour := timeFormatToStr(int(hourEnd))         //1
		normalStartTime := normalStartHour + ":" + startMinBit //23:00
		normalEndTime := normalEndHour + ":" + endMinBit       //00:30
		normalSpan = normalStartTime + "-24:00"                //23:00-24:00
		sessionInNextDay = "00:00-" + normalEndTime            //00:00-00:30
	} else if start >= AM1 { //23:00-24:00
		hourStart = hourStart + hour - Daylength //Backward one hour
		hourEnd = hourEnd + hour - Daylength     //Backward one hour
		normalStartHour := timeFormatToStr(int(hourStart))
		normalEndHour := timeFormatToStr(int(hourEnd))
		normalStartTime := normalStartHour + ":" + startMinBit
		normalEndTime := normalEndHour + ":" + endMinBit
		normalSpan = ""
		sessionInNextDay = normalStartTime + "-" + normalEndTime
	}
	return
}

//Time format
func timeFormatToStr(hour int) (time string) {
	time = strconv.Itoa(hour)
	if hour < 10 {
		time = "0" + strconv.Itoa(hour)
	}
	return time
}

type FeedTimeType int

const (
	EarlierThanCurrent FeedTimeType = iota
	LaterThanCurrent
)

const (
	WeekLength = 7
	Daylength  = 24
)

type ChangeDay int

const (
	EarlierDay ChangeDay = -1
	LaterDay   ChangeDay = 1
)
