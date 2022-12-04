package services

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
	"github.com/vtv-us/kahoot-backend/internal/constants"
)

type Participant struct {
	Username  string
	IsTeacher bool
	Status    string
	SID       string
	Score     []int
}

// [ID] -> List of participants
type Room map[string][]Participant

var room Room

func InitSocketServer() *socketio.Server {

	server := socketio.NewServer(nil)

	room = make(Room)

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
	server.OnEvent("/", "getActiveParticipants", func(s socketio.Conn, id string) {
		activeParticipants := make([]Participant, 0)
		for _, participant := range room[id] {
			if participant.Status == constants.SocketParticipantStatus_ACTIVE {
				activeParticipants = append(activeParticipants, participant)
			}
		}
		s.Emit("getActiveParticipants", activeParticipants)
	})

	server.OnEvent("/", "host", func(s socketio.Conn, username, roomID string) {
		fmt.Println("host:", roomID)
		room[roomID] = append(room[roomID], Participant{
			IsTeacher: true,
			Username:  username,
			Status:    constants.SocketParticipantStatus_ACTIVE,
			SID:       s.ID(),
		})
		fmt.Println(room)
		s.Join(roomID)
	})

	server.OnEvent("/", "join", func(s socketio.Conn, username, roomID string) {
		fmt.Println(s.ID(), "join room", roomID)
		s.Join(roomID)
		room[roomID] = append(room[roomID], Participant{
			IsTeacher: false,
			Username:  username,
			Status:    constants.SocketParticipantStatus_ACTIVE,
			SID:       s.ID(),
		})
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
		// reply to another connection
		s.Join("chat")
		server.BroadcastToRoom("/", "chat", "notice", "notice from server")
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
		s.Emit("error", e)
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
