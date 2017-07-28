package gamer

import (
	"math/rand"
	"sort"
	"time"
)

type Cards []Card
type Card byte
type CardColor byte
type CardValue byte
type CardsType byte

const (
	HighCard      CardsType = iota + 1 //高牌
	Pair                               //对子
	Flush                              //同花
	Straight                           //顺子
	StraightFlush                      //同花顺
	ThreeOfAKind                       //豹子
)

const (
	Heart   CardColor = iota + 1 //红桃
	Spade                        //黑桃
	Club                         //梅花
	Diamond                      //方块
)

func (cds Cards) ToBytes() (cards []byte) {
	for _, card := range cds {
		cards = append(cards, byte(card))
	}
	return
}
func (cds Cards) Equal(cards Cards) bool {
	if cds.GetType() == cards.GetType() {
		//牌型相同
		cType := cds.GetType()
		if cType == ThreeOfAKind {
			//豹子的牌值需相等
			if cds[0].equal(cards[0]) {
				return true
			}
		} else if cType == StraightFlush || cType == Straight || cType == Flush || cType == HighCard {
			//三张牌依次相等
			bc := cds
			lc := cards
			sort.Sort(bc)
			sort.Sort(lc)
			if bc[0].equal(lc[0]) && bc[1].equal(lc[1]) && bc[2].equal(lc[2]) {
				return true
			}
		} else if cType == Pair {
			//对子牌和单牌需相等
			bp, bs := cds.splitPair()
			lp, ls := cards.splitPair()
			if bp.equal(lp) && bs.equal(ls) {
				return true
			}
		}
	}
	return false
}

func (cds Cards) BiggerThan(cards Cards) bool {
	if cds.GetType() > cards.GetType() {
		//牌型更大
		return true
	} else if cds.GetType() == cards.GetType() {
		//牌型相同
		cType := cds.GetType()
		if cType == ThreeOfAKind {
			//豹子
			if cds[0].biggerThan(cards[0]) {
				//牌值更大
				return true
			}
		} else if cType == StraightFlush || cType == Straight || cType == Flush || cType == HighCard {
			//依次比较已排序的牌
			bc := cds
			lc := cards
			sort.Sort(bc)
			sort.Sort(lc)
			if bc[0].biggerThan(lc[0]) {
				return true
			} else if bc[1].biggerThan(lc[1]) {
				return true
			} else if bc[2].biggerThan(lc[2]) {
				return true
			}
		} else if cType == Pair {
			bp, bs := cds.splitPair()
			lp, ls := cards.splitPair()
			if bp.biggerThan(lp) {
				return true
			} else if bp.equal(lp) {
				if bs.biggerThan(ls) {
					return true
				}
			}
		}
	}
	return false
}
func (cds Cards) GetType() CardsType {
	if cds.isThreeOfAKind() {
		return ThreeOfAKind
	} else if cds.isStraightFlush() {
		return StraightFlush
	} else if cds.isStraight() {
		return Straight
	} else if cds.isFlush() {
		return Flush
	} else if cds.isPair() {
		return Pair
	} else {
		return HighCard
	}
}

//splitPair 取出对子的对子牌和单牌
func (cds Cards) splitPair() (pairCard, singleCard Card) {
	p1 := cds[0:2]
	p2 := cds[1:3]
	if p1.isAllSameValue() || p2.isAllSameValue() {
		return cds[0], cds[2]
	}
	return cds[2], cds[0]
}

//isPair 是否对子
func (cds Cards) isPair() bool {
	p1 := cds[0:2]
	p2 := cds[1:3]
	if cds.isThreeOfAKind() {
		return false
	}
	if p1.isAllSameValue() || p2.isAllSameValue() {
		return true
	}
	return false
}

//isFlush 是否同花
func (cds Cards) isFlush() bool {
	if cds.isAllSeq() {
		return false
	}
	if !cds.isAllSameColor() {
		return false
	}

	return true
}

//isStraight 是否顺子
func (cds Cards) isStraight() bool {
	if cds.isAllSameColor() {
		return false
	}
	if cds.isA23() || cds.isAllSeq() {
		return true
	}

	return false
}

//isStraightFlush 是否同花顺
func (cds Cards) isStraightFlush() bool {
	if (cds.isAllSeq() || cds.isA23()) && cds.isAllSameColor() {
		return true
	}
	return false
}

//isThreeOfAKind 是否豹子
func (cds Cards) isThreeOfAKind() bool {
	return cds.isAllSameValue()
}
func (cds Cards) isAllSameValue() bool {
	tv, _ := cds[0].splitCard()
	for _, card := range cds {
		cv, _ := card.splitCard()
		if tv != cv {
			return false
		}
	}
	return true
}

//isAllSameSeq 是否所有连牌
func (cds Cards) isAllSeq() bool {
	cards := cds
	sort.Sort(cards)
	for i := 0; i < len(cards)-1; i++ {
		b, _ := cards[i].splitCard()
		l, _ := cards[i+1].splitCard()
		if (l + 1) != b {
			return false
		}
	}
	return true
}

//isAllSameColor 是否花色相同
func (cds Cards) isAllSameColor() bool {

	_, c := cds[0].splitCard()
	for _, card := range cds {
		_, cc := card.splitCard()
		if c != cc {
			return false
		}
	}

	return true
}
func (cds Cards) isA23() bool {
	cards := cds
	sort.Sort(cards)
	var cvs []CardValue
	for _, card := range cards {
		cv, _ := card.splitCard()
		cvs = append(cvs, cv)
	}
	if cvs[0] == 0x2 && cvs[1] == 0x3 && cvs[2] == 0xe {
		return true
	}
	return false
}

//inCards 判断单牌是否在牌组中
func (cds Cards) inCards(c Card) bool {
	for _, card := range cds {
		if c == card {
			return true
		}
	}
	return false
}

//Len
func (cds Cards) Len() int {
	return len(cds)
}

//Swap
func (cds Cards) Swap(i, j int) {
	cds[i], cds[j] = cds[j], cds[i]
}

//Less
func (cds Cards) Less(i, j int) bool {
	vi, _ := cds[i].splitCard()
	vj, _ := cds[j].splitCard()
	return vi > vj
}

//filter 返回从牌组中去掉指定牌后的牌组
func (cds Cards) diff(cards Cards) (newCards Cards) {
	for _, card := range cds {
		if !cards.inCards(card) {
			newCards = append(newCards, card)
		}
	}
	return
}
func comebineCard(color CardColor, value CardValue) (card Card) {
	x := byte(color<<4) | 0
	x = byte(value) | x
	card = Card(x)
	return
}

type CardDealer struct {
	deck Cards
}

func NewCardDealer() *CardDealer {
	deck := Cards{}
	for color := Heart; color <= Diamond; color++ {
		for value := CardValue(2); value <= CardValue(0xe); value++ {
			deck = append(deck, comebineCard(color, value))
		}
	}
	de := &CardDealer{}
	de.deck = deck
	return de
}

//Shuffle 洗牌
func (de *CardDealer) Shuffle() {
	de.deck = nil
	for color := Heart; color <= Diamond; color++ {
		for value := CardValue(2); value <= CardValue(0xe); value++ {
			de.deck = append(de.deck, comebineCard(color, value))
		}
	}
}

//NextSuit 获取下一组牌
func (de *CardDealer) NextSuit() (cards Cards) {
	rand.Seed(time.Now().UnixNano())
	for {
		p := rand.Intn(len(de.deck))
		if cards.inCards(de.deck[p]) {
			continue
		}
		cards = append(cards, de.deck[p])
		if len(cards) >= 3 {
			break
		}
	}
	newDeck := de.deck.diff(cards)
	de.deck = newDeck
	return
}

//splitCard 获取牌的牌值和花色
func (c Card) splitCard() (cv CardValue, cc CardColor) {
	cv = CardValue(byte(c) & 0x0F)
	cc = CardColor(byte(c) >> 4)
	return
}

//biggerThan 比单牌牌值大小
func (c Card) biggerThan(card Card) bool {
	bv, _ := c.splitCard()
	lv, _ := card.splitCard()
	if bv > lv {
		return true
	}
	return false
}

//equal 判断两单牌牌值是否相同
func (c Card) equal(card Card) bool {
	bv, _ := c.splitCard()
	lv, _ := card.splitCard()
	if bv == lv {
		return true
	}
	return false
}
