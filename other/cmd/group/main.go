package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"
)

const (
	SexMale   = "male"
	SexFemale = "female"
)

type Sex string

type people struct {
	name  string
	sex   Sex
	i64   int64
	rank  int
	group int
}

func (s people) String() string {
	return fmt.Sprintf("%s_%s", s.name, s.sex)
}

type peoples []people

func (s peoples) Len() int {
	return len(s)
}

func (s peoples) Less(i, j int) bool {
	return s[i].i64 < s[j].i64
}

func (s peoples) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *peoples) setRank() {
	for i := 0; i < len(*s); i++ {
		(*s)[i].rank = i
	}
}

func (s peoples) String() string {
	var ms []byte
	for _, p := range s {
		ms = append(ms, []byte(p.name)...)
		ms = append(ms, []byte(",")...)
	}

	return string(ms)
}

func (s *peoples) Mod() {
	for i, p := range *s {
		(*s)[i].group = p.rank % 12
	}
}

//main go run main.go a,b,c,e f,g,h
func main() {
	males := peoples(group(strings.Split(os.Args[1], ","), SexMale))
	females := peoples(group(strings.Split(os.Args[2], ","), SexFemale))

	sort.Sort(males)
	males.setRank()
	males.Mod()

	sort.Sort(females)
	females.setRank()
	females.Mod()

	pm := make(map[int][]people)
	for i := 0; i < len(males); i++ {
		pm[males[i].group] = append(pm[males[i].group], males[i])
	}
	for j := 0; j < len(females); j++ {
		pm[11-females[j].group] = append(pm[11-females[j].group], females[j])
	}

	for _, ps := range pm {
		fmt.Println(ps)
	}

}

func group(slice []string, s Sex) []people {
	var arr []people
	for _, name := range slice {
		arr = append(arr, people{
			name: name,
			sex:  s,
			i64:  randInt(),
			rank: 0,
		})
	}
	return arr
}

func randInt() int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(100))
	return n.Int64()
}
