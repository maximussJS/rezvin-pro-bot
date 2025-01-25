package callback_data

const (
	MainPrefix      = "m"
	MainBackToMain  = "mbm"
	MainBackToStart = "mbts"

	BackPrefix             = "b"
	BackToProgramMenu      = "btpm"
	BackToProgramList      = "btpl"
	BackToPendingUsersList = "btpul"
	BackToClientList       = "btcl"

	ExercisePrefix     = "e"
	ExerciseList       = "el"
	ExerciseAdd        = "ea"
	ExerciseDelete     = "ed"
	ExerciseDeleteItem = "edi"

	ClientPrefix               = "c"
	ClientList                 = "cl"
	ClientSelected             = "cls"
	ClientProgramList          = "cpl"
	ClientProgramSelected      = "cps"
	ClientProgramAdd           = "cpad"
	ClientProgramAssign        = "cpan"
	ClientProgramDelete        = "cpd"
	ClientResultList           = "crl"
	ClientResultModifyList     = "crml"
	ClientResultModifySelected = "crms"

	ProgramPrefix   = "pr"
	ProgramSelected = "prs"
	ProgramRename   = "prr"
	ProgramDelete   = "prd"
	ProgramMenu     = "prm"
	ProgramList     = "prl"
	ProgramAdd      = "pra"

	PendingUsersPrefix   = "pu"
	PendingUsersList     = "pul"
	PendingUsersSelected = "pus"
	PendingUsersApprove  = "pua"
	PendingUsersDecline  = "pud"

	RegisterPrefix = "r"
	UserRegister   = "ru"

	UserPrefix               = "u"
	UserProgramList          = "upl"
	UserProgramSelected      = "ups"
	UserResultList           = "url"
	UserResultModifyList     = "urml"
	UserResultModifySelected = "urms"
)
