package services

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
	"github.com/vtv-us/kahoot-backend/internal/constants"
)

type RoomContext struct {
	Username string
	RoomID   string
}

type Participant struct {
	Username  string
	IsTeacher bool
	Status    string
	SID       string
	// [Question index] -> Answer
	Answer map[int]int
}

// [ID] -> List of participants
type Room map[string][]Participant
type RoomState map[string]int

var room Room

func InitSocketServer(serverAPI *Server) *socketio.Server {

	server := socketio.NewServer(nil)

	room = make(Room)
	roomState := make(RoomState)

	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "getRoomActive", func(s socketio.Conn) {
		ids := make([]string, 0)
		for id, participants := range room {
			if len(participants) > 0 {
				ids = append(ids, id)
			}
		}
		s.Emit("getRoomActive", ids)
	})
	server.OnEvent("/", "getActiveParticipants", func(s socketio.Conn) {
		ctx := s.Context().(*RoomContext)
		activeParticipants := make([]Participant, 0)
		for _, participant := range room[ctx.RoomID] {
			if participant.Status == constants.SocketParticipantStatus_ACTIVE {
				activeParticipants = append(activeParticipants, participant)
			}
		}
		s.Emit("getActiveParticipants", activeParticipants)
	})

	server.OnEvent("/", "host", func(s socketio.Conn, username, roomID string) {
		ctx := &RoomContext{
			Username: username,
			RoomID:   roomID,
		}
		s.SetContext(ctx)
		fmt.Println("host:", roomID)
		roomState[roomID] = 1
		// check if room already has a teacher
		for _, participant := range room[roomID] {
			if participant.IsTeacher && participant.Username != username {
				s.Emit("error", "Room already has a teacher")
				return
			}
		}

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

	server.OnEvent("/", "getRoomState", func(s socketio.Conn) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		if _, ok := roomState[roomID]; !ok {
			roomState[roomID] = 0
		}
		s.Emit("getRoomState", roomState[roomID])
	})

	server.OnEvent("/", "next", func(s socketio.Conn) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		username := ctx.Username
		err := checkTeacherPermission(username, roomID)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}
		roomState[roomID]++
		server.BroadcastToRoom("/", roomID, "getRoomState", roomState[roomID])
	})

	server.OnEvent("/", "prev", func(s socketio.Conn) {
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
			server.BroadcastToRoom("/", roomID, "getRoomState", roomState[roomID])
		} else {
			s.Emit("error", "You are at the first question")
		}
	})

	server.OnEvent("/", "join", func(s socketio.Conn, username, roomID string) {
		fmt.Println(s.ID(), "join room", roomID)
		ctx := &RoomContext{
			Username: username,
			RoomID:   roomID,
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

	server.OnEvent("/", "submitAnswer", func(s socketio.Conn, question int, answer int) {
		ctx := s.Context().(*RoomContext)
		username := ctx.Username
		roomID := ctx.RoomID
		fmt.Println("submitAnswer:", username, roomID, answer)
		for i, participant := range room[roomID] {
			if participant.Username == username {
				if participant.Answer == nil {
					participant.Answer = make(map[int]int)
				}
				participant.Answer[question] = answer
				room[roomID][i] = participant
			}
		}
		s.Emit("notify", "Your answer has been submitted")
		count := make(map[int]int)
		for _, participant := range room[roomID] {
			if participant.Answer != nil {
				if participant.Answer[question] != 0 {
					count[participant.Answer[question]]++
				}
			}
		}
		// send to all participants
		server.BroadcastToRoom("/", roomID, "showStatistic", count)
	})

	server.OnEvent("/", "showStatistic", func(s socketio.Conn, question int) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		// count answer
		count := make(map[int]int)
		for _, participant := range room[roomID] {
			if participant.Answer != nil {
				if participant.Answer[question] != 0 {
					count[participant.Answer[question]]++
				}
			}
		}
		// send to all participants
		server.BroadcastToRoom("/", roomID, "showStatistic", count)
	})

	server.OnEvent("/", "saveSlideHistory", func(s socketio.Conn) {
		ctx := s.Context().(*RoomContext)
		roomID := ctx.RoomID
		// save slide history
		// question idx -> answer idx-> count
		question := make(map[int]map[int]int)
		for _, participant := range room[roomID] {
			if participant.Answer != nil {
				for q, a := range participant.Answer {
					if question[q] == nil {
						question[q] = make(map[int]int)
					}
					question[q][a]++
				}
			}
		}
		// save to db
		err := serverAPI.SlideService.CreateSlideHistory(CreateSlideHistoryRequest{
			SlideID:        roomID,
			QuestionResult: question,
		})
		if err != nil {
			s.Emit("error", fmt.Errorf("save slide history failed: %w", err))
			return
		}
		s.Emit("notify", "Your slide history has been saved")
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
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
			}
		}
	})

	server.OnConnect("/chat", func(s socketio.Conn) error {
		return nil
	})

	server.OnEvent("/chat", "join", func(s socketio.Conn, username, roomID string) {
		fmt.Println("chat join:", username, roomID)
		s.Join(roomID)
		s.SetContext(username)
	})

	server.OnEvent("/chat", "send", func(s socketio.Conn, msg string) {
		username := s.Context().(string)
		fmt.Println("chat send:", username, msg)
		s.Emit("reply", username, msg)
	})

	server.OnError("/chat", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
		s.Emit("error", e)
	})

	return server
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
