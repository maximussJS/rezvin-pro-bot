package types

import (
	"rezvin-pro-bot/src/constants"
)

type Params struct {
	ProgramId     uint
	UserId        int64
	ExerciseId    uint
	UserMeasureId uint
	UserProgramId uint
	UserResultId  uint
	MeasureId     uint
	Limit         int
	Offset        int
	Reps          constants.Reps
}

func NewEmptyParams() *Params {
	return &Params{
		ProgramId:     0,
		UserId:        0,
		ExerciseId:    0,
		UserMeasureId: 0,
		UserProgramId: 0,
		UserResultId:  0,
		MeasureId:     0,
		Limit:         constants.DefaultLimit,
		Offset:        constants.DefaultOffset,
		Reps:          constants.Zero,
	}
}
