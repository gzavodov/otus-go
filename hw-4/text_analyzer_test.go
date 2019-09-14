package hw4

import (
	"strings"
	"testing"
)

type MostCommonWordTest struct {
	Title   string
	Text    string
	Limit   int
	Results []string
}

func (t *MostCommonWordTest) CheckResults(words []string) bool {
	if len(t.Results) != len(words) {
		return false
	}

	for i, word := range words {
		if !strings.EqualFold(t.Results[i], word) {
			return false
		}
	}

	return true
}

func TestGetMostCommonWords(t *testing.T) {
	tests := []MostCommonWordTest{
		MostCommonWordTest{
			Title: "This Is the House That Jack Built",
			Text: `
				This is the house that Jack built.
				This is the malt that lay in the house that Jack built.
		
				This is the rat that ate the malt
				That lay in the house that Jack built.
				This is the cat that killed the rat
				That ate the malt that lay in the house that Jack built.
		
				This is the dog that worried the cat
				That killed the rat that ate the malt
				That lay in the house that Jack built.
		
				This is the dog that worried the cat
				That killed the rat that ate the malt
				That lay in the house that Jack built.
				
				This is the cow with the crumpled horn
				That tossed the dog that worried the cat
				That killed the rat that ate the malt
				That lay in the house that Jack built.`,
			Limit: 10,
			Results: []string{
				"the",   /* 27 */
				"that",  /* 26 */
				"built", /*  7 */
				"house", /*  7 */
				"is",    /*  7 */
				"jack",  /*  7 */
				"this",  /*  7 */
				"in",    /*  6 */
				"lay",   /*  6 */
				"malt",  /*  6 */
			},
		},
		MostCommonWordTest{
			Title: "She Sells Seashells at the Seashore",
			Text: `
				She sells sea shells at the seashore,
				the shells she sells are the
				seashore shells, I am sure.`,
			Limit: 10,
			Results: []string{
				"shells",   /* 3 */
				"the",      /* 3 */
				"seashore", /* 2 */
				"sells",    /* 2 */
				"she",      /* 2 */
				"am",       /* 1 */
				"are",      /* 1 */
				"at",       /* 1 */
				"i",        /* 1 */
				"sea",      /* 1 */
			},
		},
		MostCommonWordTest{
			Title: "Говорил командир про полковника",
			Text: `
				Говорил командир про полковника и про полковницу, 
				про подполковника и про подполковницу, 
				про поручика и про поручицу, 
				про подпоручика и про подпоручицу, 
				про прапорщика и про прапорщицу, 
				про подпрапорщика, а про подпрапорщицу молчал.`,
			Limit: 10,
			Results: []string{
				"про",           /* 12 */
				"и",             /*  5 */
				"а",             /*  1 */
				"говорил",       /*  1 */
				"командир",      /*  1 */
				"молчал",        /*  1 */
				"подполковника", /*  1 */
				"подполковницу", /*  1 */
				"подпоручика",   /*  1 */
				"подпоручицу",   /*  1 */
			},
		},
		MostCommonWordTest{
			Title: "Десять негритят",
			Text: `
				Десять негритят отправились обедать,
				Один поперхнулся, их осталось девять.
				
				Девять негритят, поев, клевали носом,
				Один не смог проснуться, их осталось восемь.
				
				Восемь негритят в Девон ушли потом,
				Один не возвратился, остались всемером.
				
				Семь негритят дрова рубили вместе,
				Зарубил один себя — и осталось шесть их.
				
				Шесть негритят пошли на пасеку гулять,
				Одного ужалил шмель, их осталось пять.
				
				Пять негритят судейство учинили,
				Засудили одного, осталось их четыре.
				
				Четыре негритёнка пошли купаться в море,
				Один попался на приманку, их осталось трое.
				
				Трое негритят в зверинце оказались,
				Одного схватил медведь, и вдвоём остались.
				
				Двое негритят легли на солнцепёке,
				Один сгорел — и вот один, несчастный, одинокий.
				
				Последний негритёнок поглядел устало,
				Он пошёл повесился, и никого не стало.`,
			Limit: 10,
			Results: []string{
				"негритят", /* 8 */
				"один",     /* 7 */
				"их",       /* 6 */
				"осталось", /* 6 */
				"и",        /* 4 */
				"в",        /* 3 */
				"на",       /* 3 */
				"не",       /* 3 */
				"одного",   /* 3 */
				"восемь",   /* 2 */
			},
		},
	}
	analyzer := TextAnalyzer{}
	for _, test := range tests {
		words := analyzer.GetMostCommonWords(test.Text, test.Limit)
		if test.CheckResults(words) {
			t.Logf("SUCCESS: '%s'\n", test.Title)
			t.Logf("RECEIVED AS EXPECTED: %v\n", words)
		} else {
			t.Errorf("FAILURE: '%s'\n", test.Title)
			t.Errorf("EXPECTED: %v\n", test.Results)
			t.Errorf("RECEIVED: %v\n", words)
		}
	}
}
