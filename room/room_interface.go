package room

import (
	b "barrage-server/base"
	m "barrage-server/message"
	"barrage-server/user"
	"time"
)

var logger = b.Log

const (
	// hallID id of hall
	hallID = 0
)

const (
	// room status
	roomClose = uint8(iota)
	roomOpen
)

var commonHall *Hall

func init() {
	commonHall = NewHall()
	Open(commonHall, time.Minute)
}

// JoinHall join a user into common hall.
func JoinHall(u user.User) {
	if err := commonHall.UserJoin(u); err != nil {
		logger.Errorln(err)
	}
}

// LeftHall ...
func LeftHall(userID b.UserID) {
	if err := commonHall.UserLeft(userID); err != nil {
		logger.Errorln(err)
		return
	}
}

// Tiggler is a interface for Open and Close Room.
type Tiggler interface {
	// CompareAndSetStatus compare status of Tiggler with oldStatus, if oldStatus
	// is the same as that then set newStatus for Tiggler and return ture, otherwise
	// return false.
	CompareAndSetStatus(oldStatus, newStatus uint8) (isSet bool)
	ID() b.RoomID
	Status() uint8
	HandleInfoPkg(m.InfoPkg)
	InfoChan() <-chan m.InfoPkg

	// LoopOperation will be call periodically.
	LoopOperation()
}

// Open Tiggler.
func Open(r Tiggler, loopDuration time.Duration) {
	if isSet := r.CompareAndSetStatus(roomClose, roomOpen); isSet == false {
		return
	}

	// check status every 1 second.
	// if Room has been closed, stop ticker, break from loop and over the fucntion
	// else wait for ticker or infopkg.
	go func() {
		closeCheckTicker := time.NewTicker(1 * time.Second)
		broadCastTicker := time.NewTicker(loopDuration)
		var ipkg m.InfoPkg

	CLOSEROOM:
		for {
			select {
			case <-closeCheckTicker.C:
				if r.Status() != roomOpen {
					closeCheckTicker.Stop()
					broadCastTicker.Stop()
					break CLOSEROOM
				}
			case <-broadCastTicker.C:
				r.LoopOperation()
			case ipkg = <-r.InfoChan():
				r.HandleInfoPkg(ipkg)
			}
		}

		logger.Infof("InfoPkg handler of Room %d closed. \n", r.ID())
	}()

}

// Close Tiggler.
// It should send GameOverInfo to client first.
func Close(r Tiggler) {
	r.CompareAndSetStatus(roomOpen, roomClose)
}