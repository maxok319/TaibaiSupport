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

func (this *TaibaiClassroom) broadcastMessage(message string)  {
	for _, p := range this.Participants {
		p.SendMessage(message)
	}
}

func (this *TaibaiClassroom) singleMessage(userId int, message string)  {
	if p, ok := this.Participants[userId]; ok{
		p.SendMessage(message)
	}
}