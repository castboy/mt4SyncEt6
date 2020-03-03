package mt4SyncEt6

import (
	"testing"
)

func Test_GetSourceByCatagory(t *testing.T) {
	sourceFx := GetSourceByCatagory(SourceFx)
	sourceMetal := GetSourceByCatagory(SourceMetal)
	sourceEnergy := GetSourceByCatagory(SourceEnergy)
	sourceIndex := GetSourceByCatagory(SourceIndex)
	sourceCrypto := GetSourceByCatagory(SourceCrypto)
	t.Log("sources", sourceFx)
	t.Log("sources", sourceMetal)
	t.Log("sources", sourceEnergy)
	t.Log("sources", sourceIndex)
	t.Log("sources", sourceCrypto)
}

func Test_GetSessionBySource(t *testing.T) {
	sessions := GetSessionBySource(1)
	t.Log("sessionsQuote", sessions)
}

func Test_ModifyTheSession(t *testing.T) {
	sessionsQuote := GetSessionBySource(1)
	//sessionsTrade:=GetSessionBySource(1,1,1)
	newQuote, extraQuote := ModifyTheSession(sessionsQuote)
	//newTrade,extraTrade:=ModifyTheSession(sessionsTrade)
	t.Log("sessionsQuote", sessionsQuote)
	t.Log("newQuote", newQuote)
	t.Log("extraQuote", extraQuote)
}
func Test_GetSourceBySourceName(t *testing.T) {
	sourceUS30 := GetSourceBySourceName("US30")
	sourceUSA500 := GetSourceBySourceName("USA500")
	sourceNAS100 := GetSourceBySourceName("NAS100")
	sourceJPN225 := GetSourceBySourceName("JPN225")
	t.Log("sourceUS30", sourceUS30)
	t.Log("sourceUSA500", sourceUSA500)
	t.Log("sourceNAS100", sourceNAS100)
	t.Log("sourceJPN225", sourceJPN225)
}
func Test_modifyDB(t *testing.T) {
	//preapare
	sources := make([]Source, 0)
	//GetSources
	//Target://这周末的是所有的FX，所有Metal，所有Energy，所有Crypto，还有US30的，USA500的，NAS100的，JPN225的
	sources = GetSourceByCatagory(SourceFx)
	sources = append(sources, GetSourceByCatagory(SourceMetal)...)
	sources = append(sources, GetSourceByCatagory(SourceCrypto)...)
	sources = append(sources, GetSourceByCatagory(SourceCrypto)...)
	sources = append(sources, *GetSourceBySourceName("US30"))
	sources = append(sources, *GetSourceBySourceName("USA500"))
	sources = append(sources, *GetSourceBySourceName("NAS100"))
	sources = append(sources, *GetSourceBySourceName("JPN225"))
	//GetSessions related to source
	for _, v := range sources {
		//QuoteSession
		sessions := GetSessionBySource(v.ID)
		modSessions, extraSessions := ModifyTheSession(sessions)
		//Modify for modSessions
		for _, modSession := range modSessions {
			UpdateTimeSpanByID(&modSession)
		}
		//Insert for extraSessions
		for _, extraSession := range extraSessions {
			InsertSession(&extraSession)
		}
	}
	//Modify timespan
	//reset DB and write new sessions
}

func Test_TimeStringTest(t *testing.T) {
	if "12:45" > "13:00" {
		t.Log("true")
	} else {
		t.Log("false")
	}
}
