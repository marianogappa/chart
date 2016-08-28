package main

import (
	"sort"
)

func preprocessFreq(isss [][]string) ([][]float64, [][]string) {
	fss := [][]float64{}
	sss := [][]string{}
	fqs := make(map[string]int)
	for _, ss := range isss {
		if len(ss) == 0 {
			break //TODO this probably shouldn't happen
		}
		fqs[ss[0]] = fqs[ss[0]] + 1
	}

	keys := sortedKeys(fqs)
	topN := make(map[string]int)

	for i := 0; i < 9; i++ {
		if i >= len(keys) {
			break
		}
		topN[keys[i]] = fqs[keys[i]]
	}

	if len(keys) > 10 {
		sum := 0
		for i := 9; i < len(keys); i++ {
			sum += fqs[keys[i]]
		}
		topN["Other"] = sum
	}

	for s, f := range topN {
		fss = append(fss, []float64{float64(f)})
		sss = append(sss, []string{s})
	}

	return fss, sss
}

type sortedMap struct {
	m map[string]int
	s []string
}

func (sm *sortedMap) Len() int           { return len(sm.m) }
func (sm *sortedMap) Less(i, j int) bool { return sm.m[sm.s[i]] > sm.m[sm.s[j]] }
func (sm *sortedMap) Swap(i, j int)      { sm.s[i], sm.s[j] = sm.s[j], sm.s[i] }

func sortedKeys(m map[string]int) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}
