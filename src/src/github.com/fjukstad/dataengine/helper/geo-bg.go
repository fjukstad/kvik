package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

// Will read the entire matrix thing from geo and ouput it in csv format which
// looks a bit nicer

const (
	smokingStatus int = iota + 66
	hormoneTherapy
	medicationUse
	fastingStatus
	daysColStore
	menopausalStatus
	bmiGroup
	ageGroup
	sampleId = ageGroup + 9
)

var smoke, hormone, medication, fasting, days, menopause, bmi, age, ids []string

type Questionnaire struct {
	Results map[string]Result
}

type Result struct {
	Id               string
	SmokingStatus    string
	HomoneTherapy    string
	MeicationUse     string
	FastingStatus    string
	DaysColStore     string
	MenopausalStatus string
	BMIGroup         string
	AgeGroup         string
}

func main() {

	filename := "geobg.tsv"

	bgfile, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}

	defer bgfile.Close()

	reader := bufio.NewReader(bgfile)
	var lineNum = 1
	for {

		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}

		items := strings.Split(string(line), "\t")
		data := items[1:]

		// NOTE THAT SINCE SOME OF THE COLUMNS IN THE ORIGINAL DATASET SOME OF
		// THE ROWS CAN HAVE BEEN MIXED UP. IF ONE SAMPLE MISSES EITHER
		// BMI/HORMONE ETC ETC THE REST OF THE ROWS ARE MOVED ONE ROW UP.
		// THIS IS UNFORTUNATE.

		switch lineNum {
		case smokingStatus:
			s := stripFromArray("\"smoking status: ", "\"", data)
			smoke = s
		case hormoneTherapy:
			h := stripFromArray("\"hormone therapy use: ", "\"", data)
			hormone = h
		case medicationUse:
			m := stripFromArray("\"medication use: ", "\"", data)
			medication = m
		case fastingStatus:
			f := stripFromArray("\"fasting status: ", "\"", data)
			fasting = f
		case daysColStore:
			d := stripFromArray("\"days between blood collection and storage: ",
				"\"", data)
			days = d
		case menopausalStatus:
			m := stripFromArray("\"menopausal status: ", "\"", data)
			menopause = m
		case bmiGroup:
			b := stripFromArray("\"bmi group: ", "\"", data)
			bmi = b
		case ageGroup:
			a := stripFromArray("\"agegroup: ", "\"", data)
			age = a
		case sampleId:
			i := stripFromArray("\"", "\"", data)
			ids = i
		}

		lineNum = lineNum + 1

	}

	q := Questionnaire{}
	q.Results = make(map[string]Result)
	for i, id := range ids {
		s := smoke[i]
		h := hormone[i]
		med := medication[i]
		f := fasting[i]
		d := days[i]
		men := menopause[i]
		b := bmi[i]
		a := age[i]
		p := Result{id, s, h, med, f, d, men, b, a}

		q.Results[id] = p
	}

	return q
}

func stripFromArray(prefix, suffix string, array []string) []string {

	var res []string
	for _, v := range array {
		if strings.Contains(v, prefix) {
			a := strings.TrimPrefix(v, prefix)
			b := strings.TrimSuffix(a, suffix)
			c := strings.TrimSuffix(b, "\"\n") // trailing newline
			res = append(res, c)

		} else {
			res = append(res, "NA")
		}

	}

	return res

}
