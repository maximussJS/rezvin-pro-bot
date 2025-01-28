package types

import (
	"fmt"
	"rezvin-pro-bot/src/constants"
)

type Params struct {
	ProgramId            uint
	UserId               int64
	ExerciseId           uint
	UserProgramId        uint
	UserExerciseRecordId uint
	Limit                int
	Offset               int
	Reps                 constants.Reps
}

func NewEmptyParams() *Params {
	return &Params{
		ProgramId:            0,
		UserId:               0,
		ExerciseId:           0,
		UserProgramId:        0,
		UserExerciseRecordId: 0,
		Limit:                constants.DefaultLimit,
		Offset:               constants.DefaultOffset,
		Reps:                 constants.Zero,
	}
}

func (p *Params) String() string {
	return fmt.Sprintf("ProgramId: %d, UserID: %d, ExerciseID: %d, UserProgramID: %d, UserExerciseRecordID: %d, Limit: %d, Offset: %d",
		p.ProgramId, p.UserId, p.ExerciseId, p.UserProgramId, p.UserExerciseRecordId, p.Limit, p.Offset)
}
