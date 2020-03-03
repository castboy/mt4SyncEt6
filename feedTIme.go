package mt4SyncEt6

import (
	"github.com/go-xorm/xorm"
	"strconv"
	"strings"
)

var AM1 = "1:00"

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
//map[int][]Session	key1:sourceID  key2:weekday session string
func ModifyTheSession(sessions []Session) (newSessions, extraSession []Session) {
	// Extra container
	for k, v := range sessions {
		normalSession, presessionInF := timeConv(v.TimeSpan)
		//-1hour, and within same day
		if normalSession != "" {
			sessions[k].TimeSpan = normalSession
		} else {
			//删掉切片span=""的元素
			sessions = append(sessions[:k], sessions[k+1:]...)
		}
		//-1hour, and within same day
		if presessionInF != "" {
			sessionTemp := Session{
				SourceID: v.SourceID,
				Type:     v.Type,
				Weekday:  (v.Weekday - 1),
				TimeSpan: presessionInF,
			}
			extraSession = append(extraSession, sessionTemp)
		}
	}
	return sessions, extraSession
}

func timeConv(span string) (normalSpan string, sessionInPreDay string) {
	timeSlices := strings.Split(span, "-") //2:35-4:25
	start := timeSlices[0]                 //2:35
	end := timeSlices[1]                   //4:25
	startHourBit := strings.Split(timeSlices[0], ":")[0]
	startMinBit := strings.Split(timeSlices[0], ":")[1]
	endHourBit := strings.Split(timeSlices[1], ":")[0]
	endMinBit := strings.Split(timeSlices[1], ":")[1]
	backHourStart, err := strconv.ParseInt(startHourBit, 10, 64)
	backHourEnd, err := strconv.ParseInt(endHourBit, 10, 64)
	if err != nil {
		return
	}
	if start >= AM1 { //2:35-4:25
		backHourStart = backHourStart - 1 //Backward one hour
		backHourEnd = backHourEnd - 1     //Backward one hour
		normalStartHour := strconv.Itoa(int(backHourStart))
		normalEndHour := strconv.Itoa(int(backHourEnd))
		normalStartTime := normalStartHour + ":" + startMinBit
		normalEndTime := normalEndHour + ":" + endMinBit
		normalSpan = normalStartTime + "-" + normalEndTime
		sessionInPreDay = ""
	} else if start < AM1 && end > AM1 { //00:35-4:25
		backHourStart = backHourStart + 23 //Backward one hour
		backHourEnd = backHourEnd - 1      //Backward one hour
		normalStartHour := strconv.Itoa(int(backHourStart))
		normalEndHour := strconv.Itoa(int(backHourEnd))
		normalStartTime := "00" + ":" + "00"
		normalEndTime := normalEndHour + ":" + endMinBit
		normalSpan = normalStartTime + "-" + normalEndTime
		sessionInPreDay = normalStartHour + ":" + startMinBit + "-24:00"
	} else if start < AM1 && end < AM1 {
		backHourStart = backHourStart - 1 //Backward one hour
		backHourEnd = backHourStart - 1   //Backward one hour
		normalStartHour := strconv.Itoa(int(backHourStart))
		normalEndHour := strconv.Itoa(int(backHourEnd))
		normalStartTime := normalStartHour + ":" + startMinBit
		normalEndTime := normalEndHour + ":" + endMinBit
		normalSpan = ""
		sessionInPreDay = normalStartTime + "-" + normalEndTime
	}
	return
}

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
