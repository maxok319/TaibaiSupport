package Models

type TaibaiClassroom struct {
	ClassroomId int
	StartTime   int
	StopTime    int

	Participants map[int]*TaibaiClassParticipant
}

// 这里要读取库才可以 传入一个json串
func NewTaibaiClassroom(classroomId int) *TaibaiClassroom {

	classroom := &TaibaiClassroom{
		ClassroomId:  classroomId,
		Participants: make(map[int]*TaibaiClassParticipant),
	}

	return classroom
}

func (this *TaibaiClassroom) addParticipant(userId int) *TaibaiClassParticipant {
	p, ok := this.Participants[userId]
	if !ok {
		p = NewTaibaiClassParticipant(this, &TaibaiUser{UserId: userId}, StudentRole)
		this.Participants[userId] = p
	}
	return p
}

func (this *TaibaiClassroom) broadcastMessage(message string) {
	for _, p := range this.Participants {
		p.SendMessage(message)
	}
}

func (this *TaibaiClassroom) singleMessage(userId int, message string) {
	if p, ok := this.Participants[userId]; ok {
		p.SendMessage(message)
	}
}
