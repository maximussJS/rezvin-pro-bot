package bot_types

import "fmt"

type Params struct {
	ProgramId            uint
	UserId               int64
	ExerciseId           uint
	UserProgramId        uint
	UserExerciseRecordId uint
	Limit                int
	Offset               int
}

func NewEmptyParams() *Params {
	return &Params{
		ProgramId:            0,
		UserId:               0,
		ExerciseId:           0,
		UserProgramId:        0,
		UserExerciseRecordId: 0,
		Limit:                5,
		Offset:               0,
	}
}

func (p *Params) String() string {
	return fmt.Sprintf("ProgramId: %d, UserID: %d, ExerciseID: %d, UserProgramID: %d, UserExerciseRecordID: %d, Limit: %d, Offset: %d",
		p.ProgramId, p.UserId, p.ExerciseId, p.UserProgramId, p.UserExerciseRecordId, p.Limit, p.Offset)
}
