package types

import (
	"fmt"
	"gxclient-go/transaction"

	"github.com/pquerna/ffjson/ffjson"

	"strconv"
	"strings"

	sort "github.com/emirpasic/gods/utils"
	"github.com/juju/errors"
)

type Votes []VoteID

//TODO: define this
func (p Votes) MarshalTransaction(enc *transaction.Encoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	//TODO: remove duplicates
	//copy votes and sort
	votes := make([]interface{}, len(p))
	for idx, id := range p {
		votes[idx] = id
	}

	sort.Sort(votes, VoteIDComparator)
	for _, v := range votes {
		if err := enc.Encode(v); err != nil {
			return errors.Annotate(err, "encode VoteID")
		}
	}

	return nil
}

type VoteID struct {
	typ      int
	instance int
}

func (p *VoteID) UnmarshalJSON(data []byte) error {
	var str string

	if err := ffjson.Unmarshal(data, &str); err != nil {
		return errors.Annotate(err, "Unmarshal")
	}

	tk := strings.Split(str, ":")
	if len(tk) != 2 {
		return errors.Errorf("unable to unmarshal Vote from %s", str)
	}

	t, err := strconv.Atoi(tk[0])
	if err != nil {
		return errors.Annotate(err, "Atoi VoteID [type]")
	}
	p.typ = t

	in, err := strconv.Atoi(tk[1])
	if err != nil {
		return errors.Annotate(err, "Atoi VoteID [instance]")
	}
	p.instance = in

	return nil
}

func (p VoteID) GetType() int {
	return p.typ
}

func (p VoteID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d:%d"`, p.typ, p.instance)), nil
}

func (p VoteID) MarshalTransaction(enc *transaction.Encoder) error {
	bin := (p.typ & 0xff) | (p.instance << 8)
	if err := enc.Encode(uint32(bin)); err != nil {
		return errors.Annotate(err, "encode ID")
	}

	return nil
}

func NewVoteID(id string) *VoteID {
	v := VoteID{}
	if err := v.UnmarshalJSON([]byte(id)); err != nil {
		panic(errors.Annotatef(err, "unmarshal VoteID from %v", id))
	}

	return &v
}

func NewVoteIDV2(id string) *VoteID {
	v := VoteID{}
	kov := strings.Split(id, ":")
	if len(kov) != 2 {
		return nil
	}

	var err error
	if v.typ, err = strconv.Atoi(kov[0]); err != nil {
		return nil
	}

	if v.instance, err = strconv.Atoi(kov[1]); err != nil {
		return nil
	}
	return &v
}

func VoteIDComparator(a, b interface{}) int {
	aID := a.(VoteID)
	bID := b.(VoteID)

	switch {
	case aID.instance > bID.instance:
		return 1
	case aID.instance < bID.instance:
		return -1
	default:
		return 0
	}
}
