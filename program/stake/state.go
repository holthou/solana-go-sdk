package stake

import (
	"encoding/binary"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"math"
)

const StakeAccountSize = 200

type StakeAccountType uint32

const (
	Uninitialized StakeAccountType = iota
	Initialized
	Delegated
	RewardsPool
)

type Meta struct {
	RentExemptReserve uint64
	Authorized        Authorized
	Lockup            Lockup
}

type Delegation struct {
	Voter              common.PublicKey
	Stake              uint64
	ActivationEpoch    uint64
	DeactivationEpoch  uint64
	WarmupCooldownRate float64
}

type Stake struct {
	Delegation      Delegation
	creditsObserved uint64
}

type StakeAccount struct {
	Meta  Meta
	Stake *Stake
}

func StakeAccountDeserialize(data []byte) (StakeAccount, error) {
	if len(data) < 120 {
		return StakeAccount{}, fmt.Errorf("stake account data size is not enough")
	}

	// Initialized 只有meta
	// Delegated 都有
	// 其他状态为空
	stakeType := binary.LittleEndian.Uint32(data[:4])
	_ = stakeType

	account := StakeAccount{Meta: Meta{
		RentExemptReserve: binary.LittleEndian.Uint64(data[4:12]),
		Authorized: Authorized{
			Staker:     common.PublicKeyFromBytes(data[12:44]),
			Withdrawer: common.PublicKeyFromBytes(data[44:76]),
		},
		//TODO 没有足够数据验证，这个不确定
		Lockup: Lockup{
			UnixTimestamp: int64(binary.LittleEndian.Uint32(data[76:80])),
			Epoch:         binary.LittleEndian.Uint64(data[80:88]),
			Cusodian:      common.PublicKeyFromBytes(data[88:120]),
		},
	}}

	//中间有4个字段没用
	if len(data) >= StakeAccountSize {
		account.Stake = &Stake{
			Delegation: Delegation{
				Voter:              common.PublicKeyFromBytes(data[124:156]),
				Stake:              binary.LittleEndian.Uint64(data[156:164]),
				ActivationEpoch:    binary.LittleEndian.Uint64(data[164:172]),
				DeactivationEpoch:  binary.LittleEndian.Uint64(data[172:180]),
				WarmupCooldownRate: math.Float64frombits(binary.LittleEndian.Uint64(data[180:188])),
			},
			creditsObserved: binary.LittleEndian.Uint64(data[188:196]),
		}
	}
	//最后有4个字段没用

	return account, nil
}

const VoteAccountSize = 3762

type Lockout struct {
	Slot              uint64
	ConfirmationCount uint32
}

type AuthorizedVoters struct {
	Epoch           uint64
	AuthorizedVoter common.PublicKey
}

type PriorVoters struct {
	AuthorizedPubkey            common.PublicKey
	EpochOfLastAuthorizedSwitch uint64
	TargetEpoch                 uint64
}

type EpochCredits struct {
	Epoch            uint64
	Credits          uint64
	Previous_credits uint64
}

type BlockTimestamp struct {
	Slot      uint64
	Timestamp int64
}

type VoteAccount struct {
	NodePubkey           common.PublicKey
	AuthorizedWithdrawer common.PublicKey
	Commission           uint8
	Votes                []Lockout
	RootSlot             *uint64
	AuthorizedVoters     []AuthorizedVoters
	priorVoters          []PriorVoters
	EpochCredits         []EpochCredits
	LastTimestamp        BlockTimestamp
}

func VoteAccountDeserialize(data []byte) (VoteAccount, error) {
	if len(data) < VoteAccountSize {
		return VoteAccount{}, fmt.Errorf("vote account data size is not enough")
	}
	index := 0
	stakeType := binary.LittleEndian.Uint32(data[:index+4])
	index += 4
	_ = stakeType

	account := VoteAccount{
		NodePubkey:           common.PublicKeyFromBytes(data[index : index+32]),
		AuthorizedWithdrawer: common.PublicKeyFromBytes(data[index+32 : index+64]),
		Commission:           data[index+64],
	}
	index += 65

	// Votes
	{
		voteCount := int(binary.LittleEndian.Uint64(data[index : index+8]))
		index += 8
		index += 1 //TODO
		votes := make([]Lockout, 0, voteCount)
		for i := 0; i < voteCount; i++ {
			votes = append(votes, Lockout{
				Slot:              binary.LittleEndian.Uint64(data[index : index+8]),
				ConfirmationCount: binary.LittleEndian.Uint32(data[index+8 : index+12]),
			})
			index += 12
			index += 1 //TODO
		}
		account.Votes = votes
	}

	//RootSlot
	{
		rootslot := binary.LittleEndian.Uint64(data[index : index+8])
		account.RootSlot = &rootslot
		index += 8
	}

	// AuthorizedVoters
	{
		size := binary.LittleEndian.Uint64(data[index : index+8])
		index += 8
		authorizedVotes := make([]AuthorizedVoters, 0, size)
		for i := 0; i < int(size); i++ {
			authorizedVotes = append(authorizedVotes, AuthorizedVoters{
				Epoch:           binary.LittleEndian.Uint64(data[index : index+8]),
				AuthorizedVoter: common.PublicKeyFromBytes(data[index+8 : index+40]),
			})
			index += 40
		}
		account.AuthorizedVoters = authorizedVotes
	}

	//priorVoters
	index += 1545

	//EpochCredits
	{
		size := binary.LittleEndian.Uint64(data[index : index+8])
		index += 8
		epochCredits := make([]EpochCredits, 0, size)
		for i := 0; i < int(size); i++ {
			epochCredits = append(epochCredits, EpochCredits{
				Epoch:            binary.LittleEndian.Uint64(data[index : index+8]),
				Credits:          binary.LittleEndian.Uint64(data[index+8 : index+16]),
				Previous_credits: binary.LittleEndian.Uint64(data[index+16 : index+24]),
			})
			index += 24
		}
		account.EpochCredits = epochCredits
	}

	//LastTimestamp
	account.LastTimestamp = BlockTimestamp{
		Slot:      binary.LittleEndian.Uint64(data[index : index+8]),
		Timestamp: int64(binary.LittleEndian.Uint32(data[index+8 : index+16])),
	}
	index += 16

	//其他内容

	return account, nil
}
