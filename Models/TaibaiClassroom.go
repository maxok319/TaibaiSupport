package Models

type TaibaiClassroom struct {
	ClassroomId int
	StartTime   int
	StopTime    int

	Participants map[int]*TaibaiClassParticipant
}

// 这里要读取库才可以 传入一个json串
func NewTaibaiClassroom(classroomId int)  *TaibaiClassroom {

	classroom := &TaibaiClassroom{
		ClassroomId:classroomId,
		Participants:make(map[int]*TaibaiClassParticipant),
	}

	user1 := &TaibaiUser{UserId:111}
	user2 := &TaibaiUser{UserId:222}
	teacher := NewTaibaiClassParticipant(classroom, user1, TeacherRole)
	student := NewTaibaiClassParticipant(classroom, user2, StudentRole)
	classroom.Participants[user1.UserId] = teacher
	classroom.Participants[user2.UserId] = student

	return classroom
}
