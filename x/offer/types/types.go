// This software is Copyright (c) 2019-2020 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package types

import (
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
	"strings"

	"github.com/emirpasic/gods/trees/btree"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Instrument struct {
	Source, Destination string

	Orders *btree.Tree
}

type Instruments []Instrument

func (is Instruments) String() string {
	sb := strings.Builder{}

	for _, instr := range is {
		sb.WriteString(fmt.Sprintf("%v/%v - %v\n", instr.Source, instr.Destination, instr.Orders.Size()))
	}

	return sb.String()
}

func (is *Instruments) InsertOrder(order *Order) {
	for _, i := range *is {
		if i.Destination == order.Destination.Denom && i.Source == order.Source.Denom {
			i.Orders.Put(order, nil)
			return
		}
	}

	i := Instrument{
		Source:      order.Source.Denom,
		Destination: order.Destination.Denom,
		Orders:      btree.NewWith(3, OrderPriorityComparator),
	}

	*is = append(*is, i)
	i.Orders.Put(order, nil)
}

func (is *Instruments) GetInstrument(source, destination string) *Instrument {
	for _, i := range *is {
		if i.Source == source && i.Destination == destination {
			return &i
		}
	}

	return nil
}

func (is *Instruments) RemoveInstrument(instr Instrument) {
	for index, v := range *is {
		if instr.Source == v.Source && instr.Destination == v.Destination {
			*is = append((*is)[:index], (*is)[index+1:]...)
			return
		}
	}
}

type Order struct {
	ID uint64

	Source, Destination sdk.Coin
	SourceRemaining     sdk.Int

	Owner         sdk.AccAddress
	ClientOrderID string

	price,
	invertedPrice sdk.Dec
}

// Should return a number:
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
func OrderPriorityComparator(a, b interface{}) int {
	aAsserted := a.(*Order)
	bAsserted := b.(*Order)

	// Price priority
	switch {
	case aAsserted.Price().LT(bAsserted.Price()):
		return -1
	case aAsserted.Price().GT(bAsserted.Price()):
		return 1
	}

	// Time priority
	return int(aAsserted.ID - bAsserted.ID)
}

func (o Order) InvertedPrice() sdk.Dec {
	return o.invertedPrice
}

func (o Order) Price() sdk.Dec {
	return o.price
}

func (o Order) String() string {
	return fmt.Sprintf("%d : %v -> %v @ %v/%v (%v remaining) %v", o.ID, o.Source, o.Destination, o.price, o.invertedPrice, o.SourceRemaining, o.Owner.String())
}

func NewOrder(src, dst sdk.Coin, seller sdk.AccAddress, clientOrderId string) *Order {
	return &Order{
		Owner:           seller,
		Source:          src,
		Destination:     dst,
		SourceRemaining: src.Amount,
		ClientOrderID:   clientOrderId,
		price:           dst.Amount.ToDec().Quo(src.Amount.ToDec()),
		invertedPrice:   src.Amount.ToDec().Quo(dst.Amount.ToDec()),
	}
}

type Orders struct {
	accountOrders map[string]*treeset.Set
}

func NewOrders() Orders {
	return Orders{make(map[string]*treeset.Set)}
}

func (o Orders) ContainsClientOrderId(owner sdk.AccAddress, clientOrderId string) bool {
	allOrders := o.GetAllOrders(owner)

	order := &Order{ClientOrderID: clientOrderId}
	return allOrders.Contains(order)
}

func (o Orders) GetOrder(owner sdk.AccAddress, clientOrderId string) (res *Order) {
	allOrders := o.GetAllOrders(owner)

	allOrders.Find(func(_ int, value interface{}) bool {
		order := value.(*Order)
		if order.ClientOrderID == clientOrderId {
			res = order
			return true
		}

		return false
	})

	return
}

func (o *Orders) GetAllOrders(owner sdk.AccAddress) *treeset.Set {
	allOrders, found := o.accountOrders[owner.String()]

	if !found {
		// Note that comparator only uses client order id.
		allOrders = treeset.NewWith(OrderClientIdComparator)
		o.accountOrders[owner.String()] = allOrders
	}

	return allOrders
}

func (o *Orders) AddOrder(order *Order) {
	orders := o.GetAllOrders(order.Owner)
	orders.Add(order)
}

func (o *Orders) RemoveOrder(order *Order) {
	orders := o.GetAllOrders(order.Owner)
	orders.Remove(order)
}

func OrderClientIdComparator(a, b interface{}) int {
	aAsserted := a.(*Order)
	bAsserted := b.(*Order)

	return utils.StringComparator(aAsserted.ClientOrderID, bAsserted.ClientOrderID)
}
