package mt4SyncEt6

import (
	"testing"
	"time"
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
	quote, extraQuote, err := modifyTheSession(sessionsQuote[:2], 1, EarlierThanCurrent)
	if err != nil {
		t.Fatal("modifyTheSession err:", err)
	}
	//newTrade,extraTrade:=ModifyTheSession(sessionsTrade)
	t.Log("sessionsQuote", sessionsQuote)
	t.Log("newQuote", quote)
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

func Test_DeRepeate(t *testing.T) {
	/*sessions1:=[]Session{{241,48,0,time.Sunday,"23:05-24:00"},
	{241,48,0,time.Monday,"00:00-02:00"}}*/
	sessions := []Session{{241, 48, 0, time.Sunday, "23:05-24:00"},
		{241, 48, 0, time.Monday, "00:30-10:40"}}
	quote, extraQuote, err := modifyTheSession(sessions[:2], 1, EarlierThanCurrent)
	if err != nil {
		t.Fatal("modifyTheSession err:", err)
	}
	quoteNew, extraQuoteNew := DeRepeate(quote, extraQuote)
	t.Log("quote", quote)
	t.Log("extraQuote", extraQuote)

	t.Log("quoteNew", quoteNew)
	t.Log("extraQuoteNew", extraQuoteNew)
}

func Test_modifyDB(t *testing.T) {
	//preapare
	sources := make([]Source, 0)
	//GetSources
	//Target://这周末的是所有的FX，所有Metal，所有Energy，所有Crypto，还有US30的，USA500的，NAS100的，JPN225的
	sources = GetSourceByCatagory(SourceFx)
	sources = append(sources, GetSourceByCatagory(SourceMetal)...)
	sources = append(sources, GetSourceByCatagory(SourceEnergy)...)
	sources = append(sources, GetSourceByCatagory(SourceCrypto)...)
	sources = append(sources, *GetSourceBySourceName("US30"))
	sources = append(sources, *GetSourceBySourceName("USA500"))
	sources = append(sources, *GetSourceBySourceName("NAS100"))
	sources = append(sources, *GetSourceBySourceName("JPN225"))
	//GetSessions related to source
	for _, v := range sources {
		//QuoteSession
		sessions := GetSessionBySource(v.ID)
		modSessionsOld, extraSessionsOld, err := modifyTheSession(sessions, 1, EarlierThanCurrent)
		if err != nil {
			t.Fatal("modifyTheSession err:", err)
		}
		modSessions, extraSessions := DeRepeate(modSessionsOld, extraSessionsOld)
		//Remove and merge the repeated span
		//Modify for modSessions
		t.Log(modSessions)
		t.Log(extraSessions)
		for _, modSession := range modSessions {
			err := UpdateTimeSpanByID(&modSession)
			if err != nil {
				t.Fatalf("err:%+v", err)
			}
		}
		//Insert for extraSessions
		for _, extraSession := range extraSessions {
			err := InsertSession(&extraSession)
			if err != nil {
				t.Fatalf("err:%+v", err)
			}
		}
	}
}

func Test_earlyModifyTheSession(t *testing.T) {
	sessions := []Session{
		{241, 1, 0, time.Sunday, "22:05-24:00"},
		{241, 1, 0, time.Monday, "02:00-21:55"},
		{241, 1, 0, time.Monday, "22:05-24:00"},
		{241, 1, 0, time.Tuesday, "00:00-21:55"},
		{241, 1, 0, time.Tuesday, "22:05-24:00"},
		{241, 1, 0, time.Wednesday, "00:00-21:55"},
		{241, 1, 0, time.Wednesday, "22:05-24:00"},
		{241, 1, 0, time.Thursday, "00:00-21:55"},
		{241, 1, 0, time.Thursday, "22:05-24:00"},
		{241, 1, 0, time.Friday, "00:00-21:55"},
		{241, 1, 1, time.Sunday, "22:05-24:00"},
		{241, 1, 1, time.Monday, "02:00-21:55"},
		{241, 1, 1, time.Monday, "22:05-24:00"},
		{241, 1, 1, time.Tuesday, "00:00-21:55"},
		{241, 1, 1, time.Tuesday, "22:05-24:00"},
		{241, 1, 1, time.Wednesday, "00:00-21:55"},
		{241, 1, 1, time.Wednesday, "22:05-24:00"},
		{241, 1, 1, time.Thursday, "00:00-21:55"},
		{241, 1, 1, time.Thursday, "22:05-24:00"},
		{241, 1, 1, time.Friday, "00:00-21:55"},
	}
	mos, exts, err := modifyTheSession(sessions, 1, EarlierThanCurrent)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mos)
	t.Log(exts)
}

func Test_laterModifyTheSession(t *testing.T) {
	sessions := []Session{
		{241, 1, 0, time.Sunday, "22:05-24:00"},
		{241, 1, 0, time.Monday, "01:00-21:55"},
		{241, 1, 0, time.Monday, "22:05-24:00"},
		{241, 1, 0, time.Tuesday, "00:00-21:55"},
		{241, 1, 0, time.Tuesday, "22:05-24:00"},
		{241, 1, 0, time.Wednesday, "00:00-21:55"},
		{241, 1, 0, time.Wednesday, "22:05-24:00"},
		{241, 1, 0, time.Thursday, "00:00-21:55"},
		{241, 1, 0, time.Thursday, "22:05-24:00"},
		{241, 1, 0, time.Friday, "00:00-21:55"},
		{241, 1, 1, time.Sunday, "22:05-24:00"},
		{241, 1, 1, time.Monday, "01:00-21:55"},
		{241, 1, 1, time.Monday, "22:05-24:00"},
		{241, 1, 1, time.Tuesday, "00:00-21:55"},
		{241, 1, 1, time.Tuesday, "22:05-24:00"},
		{241, 1, 1, time.Wednesday, "00:00-21:55"},
		{241, 1, 1, time.Wednesday, "22:05-24:00"},
		{241, 1, 1, time.Thursday, "00:00-21:55"},
		{241, 1, 1, time.Thursday, "22:05-24:00"},
		{241, 1, 1, time.Friday, "00:00-21:55"},
	}

	mos, exts, err := modifyTheSession(sessions, 1, LaterThanCurrent)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mos)
	t.Log(exts)
}

func Test_timeTest(t *testing.T) {
	if "22:05" < "23:00" {
		t.Log(true)
	} else {
		t.Log(false)
	}
}
