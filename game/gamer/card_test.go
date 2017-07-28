package gamer

import (
	"fmt"
	"sort"
	"testing"
)

var t2s = map[CardsType]string{
	HighCard:      "高牌",
	Pair:          "对子",
	Flush:         "同花",
	Straight:      "顺子",
	StraightFlush: "同花顺",
	ThreeOfAKind:  "豹子",
}

func TestComebineCard(t *testing.T) {
	de := NewCardDealer()
	fmt.Println(len(de.deck))
	for _, v := range de.deck {
		fmt.Printf("%x ", v)
	}
}

func TestNextSuit(t *testing.T) {
	de := NewCardDealer()
	cards := de.NextSuit()
	for _, v := range cards {
		fmt.Printf("%x ", v)
	}
	fmt.Println("len", len(de.deck))
}
func TestIsThreeOfAKind(t *testing.T) {
	c1 := Cards{0x11, 0x21, 0x31}
	fmt.Println("c1: ", c1.isThreeOfAKind())
	c2 := Cards{0x13, 0x21, 0x31}
	fmt.Println("c2: ", c2.isThreeOfAKind())
}
func TestSort(t *testing.T) {
	de := NewCardDealer()
	cards := de.NextSuit()
	for _, v := range cards {
		fmt.Printf("%x ", v)
	}
	fmt.Println()
	sort.Sort(cards)
	for _, v := range cards {
		fmt.Printf("%x ", v)
	}
}
func TestIsStraight(t *testing.T) {
	c1 := Cards{0x16, 0x27, 0x38}
	fmt.Println("c1: ", c1.isStraight())
	c2 := Cards{0x13, 0x22, 0x3e}
	fmt.Println("c2: ", c2.isStraight())
}
func TestIsFlush(t *testing.T) {
	c1 := Cards{0x14, 0x15, 0x16}
	fmt.Println("c1: ", c1.isFlush())
	c2 := Cards{0x13, 0x11, 0x11}
	fmt.Println("c2: ", c2.isFlush())
}
func TestIsStraightFlush(t *testing.T) {
	c1 := Cards{0x13, 0x12, 0x1e}
	fmt.Println("c1: ", c1.isStraightFlush())
	c2 := Cards{0x13, 0x14, 0x15}
	fmt.Println("c2: ", c2.isStraightFlush())
}

func TestIsPair(t *testing.T) {
	c1 := Cards{0x16, 0x26, 0x34}
	fmt.Println("c1: ", c1.isPair())
	c2 := Cards{0x13, 0x21, 0x31}
	fmt.Println("c2: ", c2.isPair())
}
func TestIsA23(t *testing.T) {
	c1 := Cards{0x11, 0x21, 0x31}
	fmt.Println("c1: ", c1.isA23())
	c2 := Cards{0x1e, 0x23, 0x42}
	fmt.Println("c2: ", c2.isA23())
}
func TestGetType(t *testing.T) {

	c1 := Cards{0x11, 0x21, 0x31}
	fmt.Println("c1: ", t2s[c1.GetType()])
	c2 := Cards{0x1e, 0x13, 0x12}
	fmt.Println("c2: ", t2s[c2.GetType()])
	c3 := Cards{0x1e, 0x23, 0x42}
	fmt.Println("c2: ", t2s[c3.GetType()])
	c4 := Cards{0x1e, 0x13, 0x15}
	fmt.Println("c2: ", t2s[c4.GetType()])
	c5 := Cards{0x1a, 0x2e, 0x3e}
	fmt.Println("c2: ", t2s[c5.GetType()])
	c6 := Cards{0x17, 0x2e, 0x15}
	fmt.Println("c2: ", t2s[c6.GetType()])
}
func TestBiggerThan(t *testing.T) {
	//豹子比较
	a1 := Cards{0x12, 0x22, 0x32}
	a2 := Cards{0x13, 0x23, 0x33}
	fmt.Printf("a1>a2: %v, a1=a2: %v, a1<a2: %v\n======\n", a1.BiggerThan(a2), a1.Equal(a2), a2.BiggerThan(a1))
	//同花顺比较
	a1 = Cards{0x22, 0x23, 0x24}
	a2 = Cards{0x35, 0x36, 0x37}
	fmt.Printf("a1>a2: %v, a1=a2: %v, a1<a2: %v\n======\n", a1.BiggerThan(a2), a1.Equal(a2), a2.BiggerThan(a1))
	//顺子比较
	a1 = Cards{0x12, 0x23, 0x34}
	a2 = Cards{0x15, 0x26, 0x37}
	fmt.Printf("a1>a2: %v, a1=a2: %v, a1<a2: %v\n======\n", a1.BiggerThan(a2), a1.Equal(a2), a2.BiggerThan(a1))
	//同花比较
	a1 = Cards{0x12, 0x23, 0x35}
	a2 = Cards{0x35, 0x36, 0x47}
	fmt.Printf("a1>a2: %v, a1=a2: %v, a1<a2: %v\n======\n", a1.BiggerThan(a2), a1.Equal(a2), a2.BiggerThan(a1))
	//对子比较
	a1 = Cards{0x13, 0x23, 0x34}
	a2 = Cards{0x36, 0x36, 0x47}
	fmt.Printf("a1>a2: %v, a1=a2: %v, a1<a2: %v\n======\n", a1.BiggerThan(a2), a1.Equal(a2), a2.BiggerThan(a1))
	//高牌比较
	a1 = Cards{0x13, 0x2a, 0x34}
	a2 = Cards{0x36, 0x36, 0x4e}
	fmt.Printf("a1>a2: %v, a1=a2: %v, a1<a2: %v\n======\n", a1.BiggerThan(a2), a1.Equal(a2), a2.BiggerThan(a1))
}
