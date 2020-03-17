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

func DeleteSession(sess *Session) error {
	engine := GetEngine()
	_, err := engine.Where("id=？", sess.ID).Delete(sess)
	return err
}

func timeUpdate(sessions []Session, hour int, feedType FeedTimeType) (err error) {
	modSessions, insSessions, delSeessions, err := modifyTheSession(sessions, hour, feedType)
	if err != nil {
		return err
	}
	//insert
	for _, insSession := range insSessions {
		timeNods := strings.Split(insSession.TimeSpan, "-")
		if len(timeNods) != TimeNodesLength || timeNods[0] == timeNods[1] {
			continue
		}
		err = InsertSession(&insSession)
		if err != nil {
			return err
		}
	}
	//Update
	for _, modSession := range modSessions {
		timeNods := strings.Split(modSession.TimeSpan, "-")
		if len(timeNods) != TimeNodesLength || timeNods[0] == timeNods[1] {
			continue
		}
		err = UpdateTimeSpanByID(&modSession)
		if err != nil {
			return err
		}
	}
	//Del session
	for _, delsession := range delSeessions {
		timeNods := strings.Split(delsession.TimeSpan, "-")
		if len(timeNods) != TimeNodesLength || timeNods[0] == timeNods[1] {
			continue
		}
		err = DeleteSession(&delsession)
		if err != nil {
			return err
		}
	}

	return
}

func sessionConv(source *Source, sessionsMap map[time.Weekday]map[int]string, sessionType SessionType) (sessions []Session) {
	//sessionMap check
	if sessionsMap == nil || len(sessionsMap) == 0 {
		return nil
	}
	for weekDay, v := range sessionsMap {
		for id, ss := range v {
			session := Session{
				id,
				source.ID,
				sessionType,
				weekDay,
				ss,
			}
			sessions = append(sessions, session)
		}
	}
	return
}

//Modify the time by day-lighting time
func modifyTheSession(sessions []Session, hour int, feedType FeedTimeType) (modSessions, extraSession, delSession []Session, err error) {
	// Extra container
	modSessionsTemp := []Session{}
	extraSessionsTemp := []Session{}
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
			return nil, nil, nil, err
		}
		//normalSpan!=nil include -1 hour within same day or exceed one day 00:30-01:00====23:30-24:00
		if normalSpan != "" {
			v.TimeSpan = normalSpan
			modSessionsTemp = append(modSessionsTemp, v)
			//-1hour, and within previous day
			if preSpanInF != "" {
				sessionTemp := Session{
					SourceID: v.SourceID,
					Type:     v.Type,
					Weekday:  weekday,
					TimeSpan: preSpanInF,
				}
				extraSessionsTemp = append(extraSessionsTemp, sessionTemp)
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
			modSessionsTemp = append(modSessionsTemp, sessionTemp)
		}
	}
	//remove the crossing time_span, for examplp "22:00-23:30" and "23:00-24:00",which can be combined to "22:00-24:00

	modSessionsRe, extraSessionRe := deRepeate(modSessionsTemp, extraSessionsTemp, feedType)
	extraSession = append(extraSession, extraSessionRe...)
	modSessionsR, delSessionR, err := modDeRepead(modSessionsRe)

	modSessions = append(modSessions, modSessionsR...)
	delSession = append(delSession, delSessionR...)

	return
}

//Remove and merge the crossing time_span
func deRepeate(sessionMos, sessionIns []Session, feedType FeedTimeType) (newSessionMos []Session, newSessionIns []Session) {
	//Copy data
	newSessionMos = append([]Session{}, sessionMos...)    //22:05-23:00 00:00-1:00
	mediaSessionIns := append([]Session{}, sessionIns...) //23:00-24:00
	//Prepare
	for i, sessionIns := range sessionIns {
		timeIns := strings.Split(sessionIns.TimeSpan, "-")
		for k, sessionMod := range sessionMos {
			// Conditions
			if sessionIns.Weekday != sessionMod.Weekday || sessionIns.Type != sessionMod.Type || sessionIns.SourceID != sessionMod.SourceID {
				continue
			}
			//Operate
			timeMod := strings.Split(sessionMod.TimeSpan, "-")
			// merge sessionIns to sessionMod
			if feedType == EarlierThanCurrent { //
				if timeIns[0] <= timeMod[1] {
					//Merge
					sessionMod.TimeSpan = timeMod[0] + "-" + timeIns[1]
					//kick the sessionIns from sessionsIns
					newSessionMos[k] = sessionMod
					mediaSessionIns[i].TimeSpan = ""
				}
			}
			if feedType == LaterThanCurrent { //22:00-24:00 00:00-2:00
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
	//sessionMos

	for i := range mediaSessionIns {
		if mediaSessionIns[i].TimeSpan == "" {
			continue
		}
		newSessionIns = append(newSessionIns, mediaSessionIns[i])
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
	AM1 = strconv.Itoa(hour) + ":00"
	if hour < 10 {
		AM1 = "0" + strconv.Itoa(hour) + ":00"
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
		hourStart = hourStart + Daylength - 1                  //Backward one hour   23
		hourEnd = hourEnd + Daylength - 1                      //Backward one hour   23
		normalStartHour := timeFormatToStr(int(hourStart))     //23
		normalEndHour := timeFormatToStr(int(hourEnd))         //23
		normalStartTime := normalStartHour + ":" + startMinBit //23:35
		normalEndTime := normalEndHour + ":" + endMinBit       //24:00
		normalSpan = ""
		sessionInPreDay = normalStartTime + "-" + normalEndTime //23:35-24:00
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

func modDeRepead(sessions []Session) (modSS, delSS []Session, err error) {
	//Return sss
	sss := append([]Session{}, sessions...)
	sssCopy := append([]Session{}, sessions...)
	for i := range sss {
		nodsSss := strings.Split(sss[i].TimeSpan, "-")
		if len(nodsSss) != TimeNodesLength || nodsSss[0] == nodsSss[1] {
			return nil, nil, errors.Errorf("Nil or invalid time span!:%+v", nodsSss)
		}
		for j := range sessions {
			//Symbol diff and should same
			if sss[i].SourceID != sessions[j].SourceID {
				continue
			}
			//ID different should be diff
			if sss[i].ID == sessions[j].ID {
				continue
			}
			//weekday should be same
			if sss[i].Weekday != sessions[j].Weekday {
				continue
			}
			nodsSessions := strings.Split(sss[j].TimeSpan, "-")
			if len(nodsSessions) != TimeNodesLength || nodsSessions[0] == nodsSessions[1] {
				return nil, nil, errors.Errorf("Nil or invalid time span!:%+v", nodsSss)
			}
			//Merge backward
			// case1: self check 23:00-24:00 && 00:00-3:00======> 00:00-1:00 && 01:00-4:00  later
			// case1: self check 23:00-24:00 && 00:00-2:00======> 22:00-23:00  23:00-24:00 && 00:00-1:00 early
			if nodsSss[1] == nodsSessions[0] {
				span := nodsSss[0] + "-" + nodsSessions[1]
				//keep smaller SessionID
				if sss[i].ID < sessions[j].ID {
					sssCopy[i].TimeSpan = span
					//delete j
					sssCopy[j].TimeSpan = ""
				} else {
					sssCopy[j].TimeSpan = span
					//delete j
					sssCopy[i].TimeSpan = ""
				}
			}
		}
	}
	//Arrangement

	for i := range sssCopy {
		if sssCopy[i].TimeSpan == "" {
			delSS = append(delSS, sss[i])
		} else {
			modSS = append(modSS, sssCopy[i])
		}
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

type feedTimeOperator struct{}
type FeedTimeType int

const (
	EarlierThanCurrent FeedTimeType = iota
	LaterThanCurrent
)

const (
	WeekLength      = 7
	Daylength       = 24
	TimeNodesLength = 2
)

type ChangeDay int

const (
	EarlierDay ChangeDay = -1
	LaterDay   ChangeDay = 1
)
