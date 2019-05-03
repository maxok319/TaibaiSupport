package Models

type TaibaiUser struct {
	UserId int
}

type TaibaiStudent struct {
	TaibaiUser
	age int
}

type TaibaiTeacher struct {
	TaibaiUser
	email string
}

type TaibaiObserver struct {
	TaibaiUser
}
