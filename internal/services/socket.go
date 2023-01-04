package services

import (
	"context"
	"fmt"

	socketio "github.com/googollee/go-socket.io"
	"github.com/vtv-us/kahoot-backend/internal/constants"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
)

type RoomContext struct {
	Username  string
	RoomID    string
	IsTeacher bool
}

type Participant struct {
	Username  string
	IsTeacher bool
	Status    string
	SID       string
	// [Question index] -> Answer
	Answer map[int]int
}

type ChatMessage struct {
	Username string
	Message  string
}

// [ID] -> List of participants
type Room map[string][]Participant
type IsRoomGroup map[string]bool
type RoomGroup map[string]string
type RoomState map[string]int
type GroupSlidePresent map[string]string

type PresentationNotification struct {
	SlideID string
	GroupID string
}

var room Room

func InitSocketServer(server *Server) *socketio.Server {

	socket := socketio.NewServer(nil)

	room = make(Room)
	roomState := make(RoomState)
	isRoomGroup := make(IsRoomGroup)
	roomGroup := make(RoomGroup)
	groupSlidePresent := make(GroupSlidePresent)

	socket.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		return nil
	})

	socket.OnEvent("/", "getRoomActive", func(s socketio.Conn) {
		ids := make([]string, 0)
		for id, participants := range room {
			if len(participants) > 0 {
				ids = append(ids, id)
			}
		}
		s.Emit("getRoomActive", ids)
	})
	socket.OnEvent("/", "getActiveParticipants", func(s socketio.Conn) {
		ctx := s.Context().(*RoomContext)
		activeParticipants := make([]Participant, 0)
		for _, participant := range room[ctx.RoomID] {
			if participant.Status == constants.SocketParticipantStatus_ACTIVE {
				activeParticipants = append(activeParticipants, participant)
			}
		}
		s.Emit("getActiveParticipants", activeParticipants)
	})

	socket.OnEvent("/", "manualDisconnect", func(s socketio.Conn) {
		s.Close()
	})

	socket.OnEvent("/", "host", func(s socketio.Conn, username, roomID string, isGroup bool, groupID string, token string) {
		isRoomGroup[roomID] = isGroup
		if isGroup {
			err := checkUserInGroup(server, groupID, token)
			if err != nil {
				s.Emit("error", err.Error())
				return
			}
			roomGroup[roomID] = groupID
			groupSlidePresent[groupID] = roomID
			socket.BroadcastToRoom("/notification", groupID, "notify", PresentationNotification{
				SlideID: roomID,
				GroupID: groupID,
			})
		}
		ctx := &RoomContext{
			Username:  username,
			RoomID:    roomID,
			IsTeacher: true,
		}
		s.SetContext(ctx)
		fmt.Println(s.ID(), username, "host:", roomID)
		if _, ok := room[roomID]; !ok {
			roomState[roomID] = 1
		}
		// check if room already has a teacher
		// for _, participant := range room[roomID] {
		// 	if participant.IsTeacher && participant.Username != username {
		// 		s.Emit("error", "Room already has a teacher")
		// 		return
		// 	}
		// }

		exist := checkExistInRoom(username, roomID)
		if exist {
			for i, participant := range room[roomID] {
				if participant.Username == username {
					room[roomID][i].Status = constants.SocketParticipantStatus_ACTIVE
					room[roomID][i].SID = s.ID()
				}
			}
		} else {
			room[roomID] = append(room[roomID], Participant{
				IsTeacher: true,
				Username:  username,
				Status:    constants.SocketParticipantStatus_ACTIVE,
				SID:       s.ID(),
			})
		}
		s.Join(roomID)
	})

	socket.OnEvent("/", "getRoomState", func(s socketio.Conn) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		if _, ok := roomState[roomID]; !ok {
			roomState[roomID] = 0
		}
		s.Emit("getRoomState", roomState[roomID])
	})
	socket.OnEvent("/", "setRoomState", func(s socketio.Conn, state int) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		username := ctx.Username
		err := checkTeacherPermission(username, roomID)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}
		roomState[roomID] = state
		socket.BroadcastToRoom("/", roomID, "getRoomState", roomState[roomID])
	})

	socket.OnEvent("/", "next", func(s socketio.Conn) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		username := ctx.Username
		err := checkTeacherPermission(username, roomID)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}
		roomState[roomID]++
		socket.BroadcastToRoom("/", roomID, "getRoomState", roomState[roomID])
	})

	socket.OnEvent("/", "prev", func(s socketio.Conn) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		username := ctx.Username
		err := checkTeacherPermission(username, roomID)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}
		if roomState[roomID] > 1 {
			roomState[roomID]--
			socket.BroadcastToRoom("/", roomID, "getRoomState", roomState[roomID])
		} else {
			s.Emit("error", "You are at the first question")
		}
	})

	socket.OnEvent("/", "join", func(s socketio.Conn, username, roomID, token string) {
		fmt.Println(s.ID(), "join room", roomID)
		if isRoomGroup[roomID] {
			err := checkUserInGroup(server, roomGroup[roomID], token)
			if err != nil {
				s.Emit("error", err.Error())
				return
			}
		}
		ctx := &RoomContext{
			Username:  username,
			RoomID:    roomID,
			IsTeacher: false,
		}
		s.SetContext(ctx)
		exist := checkExistInRoom(username, roomID)
		if exist {
			for i, participant := range room[roomID] {
				if participant.Username == username {
					if participant.Status == constants.SocketParticipantStatus_ACTIVE {
						s.Emit("error", "You are already in the room")
						return
					}
					room[roomID][i].Status = constants.SocketParticipantStatus_ACTIVE
					room[roomID][i].SID = s.ID()
				}
			}
		} else {
			room[roomID] = append(room[roomID], Participant{
				IsTeacher: false,
				Username:  username,
				Status:    constants.SocketParticipantStatus_ACTIVE,
				SID:       s.ID(),
			})
		}
		s.Join(roomID)
	})

	socket.OnEvent("/", "cancelPresentation", func(s socketio.Conn, roomID string) {
		if _, ok := room[roomID]; !ok {
			s.Emit("notify", "Room does not active, skip cancel")
			return
		}
		ctx := s.Context().(*RoomContext)
		username := ctx.Username
		err := checkTeacherPermission(username, roomID)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}
		delete(room, roomID)
		delete(roomState, roomID)
		delete(groupSlidePresent, roomGroup[roomID])
		delete(roomGroup, roomID)
		socket.BroadcastToRoom("/", roomID, "cancelPresentation", roomID)
	})

	socket.OnEvent("/", "getSlidePresentation", func(s socketio.Conn, groupID string) {
		_, ok := groupSlidePresent[groupID]
		if !ok {
			s.Emit("error", "Group does not have any slide presentation")
			return
		}
		s.Emit("getSlidePresentation", groupSlidePresent[groupID])
	})

	socket.OnEvent("/", "submitAnswer", func(s socketio.Conn, question string, answer string) {
		ctx := s.Context().(*RoomContext)
		username := ctx.Username
		roomID := ctx.RoomID
		fmt.Println("submitAnswer:", username, roomID, answer)
		err := server.SlideService.SaveAnswerHistory(username, roomID, question, answer)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}
		s.Emit("notify", "Your answer has been submitted")
		count, err := server.SlideService.CountAnswerByQuestionID(question)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}
		// send to all participants
		socket.BroadcastToRoom("/", roomID, "showStatistic", count)
		result, err := server.SlideService.ListAnswerHistoryByQuestionID(question)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}
		socket.BroadcastToRoom("/", roomID, "resultList", result)
	})

	socket.OnEvent("/", "showStatistic", func(s socketio.Conn, question string) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		// count answer
		count, err := server.SlideService.CountAnswerByQuestionID(question)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}
		// send to all participants
		socket.BroadcastToRoom("/", roomID, "showStatistic", count)
		result, err := server.SlideService.ListAnswerHistoryByQuestionID(question)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}
		socket.BroadcastToRoom("/", roomID, "resultList", result)
	})

	// socket.OnEvent("/", "saveSlideHistory", func(s socketio.Conn) {
	// 	ctx := s.Context().(*RoomContext)
	// 	roomID := ctx.RoomID
	// 	// save slide history
	// 	// question idx -> answer idx-> count
	// 	question := make(map[int]map[int]int)
	// 	for _, participant := range room[roomID] {
	// 		if participant.Answer != nil {
	// 			for q, a := range participant.Answer {
	// 				if question[q] == nil {
	// 					question[q] = make(map[int]int)
	// 				}
	// 				question[q][a]++
	// 			}
	// 		}
	// 	}
	// 	// save to db
	// 	err := server.SlideService.CreateSlideHistory(CreateSlideHistoryRequest{
	// 		SlideID:        roomID,
	// 		QuestionResult: question,
	// 	})
	// 	if err != nil {
	// 		s.Emit("error", fmt.Errorf("save slide history failed: %w", err))
	// 		return
	// 	}
	// 	s.Emit("notify", "Your slide history has been saved")
	// })

	socket.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	socket.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason, s.ID())
		for id, participants := range room {
			for i, participant := range participants {
				if participant.SID == s.ID() {
					participant.Status = constants.SocketParticipantStatus_LEFT
					room[id][i] = participant
				}
			}
		}
		// check all left
		for id, participants := range room {
			allLeft := true
			for _, participant := range participants {
				if participant.Status == constants.SocketParticipantStatus_ACTIVE {
					allLeft = false
					break
				}
			}
			if allLeft {
				delete(room, id)
				delete(roomState, id)
				delete(groupSlidePresent, roomGroup[id])
				delete(roomGroup, id)
			}
		}
	})

	// chat
	socket.OnEvent("/", "chat", func(s socketio.Conn, msg string) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		username := ctx.Username
		if err := server.SlideService.SaveChatMsg(roomID, username, msg); err != nil {
			s.Emit("error", fmt.Errorf("save chat message failed: %w", err))
			return
		}
		// send to all participants
		socket.BroadcastToRoom("/", roomID, "chat", username, msg)
	})

	socket.OnEvent("/", "getChatHistory", func(s socketio.Conn) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		chatMsgs, err := server.SlideService.GetChatMsgs(roomID)
		if err != nil {
			s.Emit("error", fmt.Errorf("get chat history failed: %w", err))
			return
		}
		s.Emit("chatHistory", chatMsgs)
	})

	// user question
	socket.OnEvent("/", "postQuestion", func(s socketio.Conn, msg string) {
		ctx := s.Context().(*RoomContext)
		cctx := context.Background()
		roomID := ctx.RoomID
		username := ctx.Username
		question, err := server.UserQuestionService.PostQuestion(cctx, PostQuestionRequest{
			SlideID:  roomID,
			Username: username,
			Content:  msg,
		})
		if err != nil {
			s.Emit("error", fmt.Errorf("post question failed: %w", err).Error())
			return
		}
		// send to all participants
		socket.BroadcastToRoom("/", roomID, "postQuestion", question)
	})

	socket.OnEvent("/", "listUserQuestion", func(s socketio.Conn) {
		ctx := s.Context().(*RoomContext)
		cctx := context.Background()
		roomID := ctx.RoomID
		questions, err := server.UserQuestionService.ListQuestionBySlideID(cctx, roomID)
		if err != nil {
			s.Emit("error", fmt.Errorf("list user question failed: %w", err).Error())
			return
		}
		s.Emit("listUserQuestion", questions)
	})
	socket.OnEvent("/", "upvoteQuestion", func(s socketio.Conn, questionID string) {
		ctx := s.Context().(*RoomContext)
		cctx := context.Background()
		roomID := ctx.RoomID
		question, err := server.UserQuestionService.UpvoteQuestion(cctx, questionID)
		if err != nil {
			s.Emit("error", fmt.Errorf("upvote question failed: %w", err).Error())
			return
		}
		// send to all participants
		socket.BroadcastToRoom("/", roomID, "upvoteQuestion", question)
	})
	socket.OnEvent("/", "toggleUserQuestionAnswered", func(s socketio.Conn, questionID string) {
		ctx := s.Context().(*RoomContext)
		cctx := context.Background()
		roomID := ctx.RoomID
		isTeacher := ctx.IsTeacher
		if !isTeacher {
			s.Emit("error", "only owner and co-owner can toggle user question answered")
			return
		}
		question, err := server.UserQuestionService.ToggleUserQuestionAnswered(cctx, questionID)
		if err != nil {
			s.Emit("error", fmt.Errorf("toggle user question answered failed: %w", err).Error())
			return
		}
		// send to all participants
		socket.BroadcastToRoom("/", roomID, "toggleUserQuestionAnswered", question)
	})

	// server notification
	socket.OnEvent("/notification", "join", func(s socketio.Conn, token string) {
		res, err := server.AuthService.JWT.ValidateToken(token)
		if err != nil {
			s.Emit("error", fmt.Errorf("invalid token: %w", err).Error())
			return
		}
		groups, err := server.AuthService.DB.GetGroupByUser(context.Background(), res.UserID)
		if err != nil {
			s.Emit("error", fmt.Errorf("get group by user failed: %w", err).Error())
			return
		}
		for _, group := range groups {
			s.Join(group.GroupID)
		}
	})

	return socket
}

func checkExistInRoom(username, roomID string) bool {
	for _, participant := range room[roomID] {
		if participant.Username == username {
			return true
		}
	}
	return false
}

func checkTeacherPermission(username, roomID string) error {
	for _, participant := range room[roomID] {
		if participant.Username == username {
			if participant.IsTeacher {
				return nil
			}
			return fmt.Errorf("you are not a teacher in the room")
		}
	}
	return fmt.Errorf("you are not in the room")
}

func checkUserInGroup(server *Server, groupID, token string) error {
	// check token
	res, err := server.AuthService.JWT.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	isUserInGroup, err := server.GroupService.DB.CheckUserInGroup(context.Background(), repositories.CheckUserInGroupParams{
		GroupID: groupID,
		UserID:  res.UserID,
	})
	if err != nil {
		return fmt.Errorf("check user in group failed: %w", err)
	}
	if !isUserInGroup {
		return fmt.Errorf("you are not in the group")
	}
	return nil
}
