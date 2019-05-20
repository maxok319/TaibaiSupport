package TaibaiDBHelper

import (
	"TaibaiSupport/Models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var redisClientInstance *redis.Client

func GetInstance() *redis.Client {
	return redisClientInstance
}

func init() {
	redisClientInstance = redis.NewClient(&redis.Options{
		Addr:     "taibai-redis-service:6379",
		Password: "", // no password set
		DB:       2,  // use default DB
	})
}

func AddClassroom(classroomId int)  {
	already, err:= redisClientInstance.SIsMember("classroomSet", classroomId).Result()
	if err!=nil {
		return
	}

	if already {
		return
	}

	redisClientInstance.SAdd("classroomSet", classroomId)
	redisClientInstance.HMSet("classroomStatus:"+strconv.Itoa(classroomId), map[string]interface{}{
		"startTime": time.Now().Unix(),
		"stopTime": time.Now().Unix() + 25 * 60 * 1000,
		"currentPage": 0,
	})
}

func AddUserIntoClassroom( classroomId, userId int)  {
	already, err:= redisClientInstance.SIsMember("studentSet:"+strconv.Itoa(classroomId), userId).Result()
	if err !=nil{
		return
	}
	if already{
		return
	}

	index := redisClientInstance.SCard("studentSet:"+strconv.Itoa(classroomId)).Val()
	rect := &Models.TaibaiRect{
			X: int(20 + (20+200) * index),
			Y: 810-200,
			Width : 200,
			Height:200,
		}
	redisClientInstance.SAdd("studentSet:"+strconv.Itoa(classroomId), userId)
	redisClientInstance.HMSet("studentStatus:" + strconv.Itoa(userId), map[string]interface{}{
		"index":index,
		"rect":rect,
	})
}

func UpdateUserOnlineStatus(userId int, online bool)  {
	redisClientInstance.HSet("studentStatus:" + strconv.Itoa(userId), "online", online)
}

func UpdateUserRect(userId int, rect string)  {
	redisClientInstance.HSet("studentStatus:" + strconv.Itoa(userId), "rect", rect)
}

func GetClassroomStatus(classroomId int) map[string]interface{} {
	startTime,_ :=redisClientInstance.HGet("classroomStatus:"+strconv.Itoa(classroomId), "startTime").Int()
	stopTime,_ :=redisClientInstance.HGet("classroomStatus:"+strconv.Itoa(classroomId), "stopTime").Int()
	currentPage,_ :=redisClientInstance.HGet("classroomStatus:"+strconv.Itoa(classroomId), "currentPage").Int()

	return map[string]interface{}{
		"startTime": startTime,
		"stopTime": stopTime,
		"currentPage": currentPage,
	}
}

func GetUserList(classroomId int) []int {
	result := []int{}
	students, _ := redisClientInstance.SMembers("studentSet:"+strconv.Itoa(classroomId)).Result()
	for _, student := range students{
		studentId,_:= strconv.Atoi(student)
		result = append(result, studentId)
	}
	return result
}

func GetUserStatus(userId int) map[string] interface{} {
	index, _ :=redisClientInstance.HGet("studentStatus:"+strconv.Itoa(userId), "index").Int()
	onlineInt,_ :=redisClientInstance.HGet("studentStatus:"+strconv.Itoa(userId), "online").Int()
	online := onlineInt == 1
	rect := Models.TaibaiRect{}
	_ = redisClientInstance.HGet("studentStatus:"+strconv.Itoa(userId), "rect").Scan(&rect)

	return map[string] interface{}{
		"index": index,
		"online": online,
		"rect": rect,
	}
}