package dummyemitter

import (
	"faceittt/domain"
	"fmt"
	"log"
)

func LogUserEvent(e domain.UserEvent) error {
	switch e.Type() {
	case domain.UserEventUnassigned:
		return fmt.Errorf("undefined event emitted")
	case domain.UserEventCreated:
		ce := e.(domain.UserCreatedEvent)
		log.Println("yaaaay user created ", ce.User.ID)
	case domain.UserEventUpdated:
		ue := e.(domain.UserUpdateEvent)
		log.Println("yaaaay user updated ", ue.NewUser.ID)
	case domain.UserEventDeleted:
		de := e.(domain.UserDeletedEvent)
		log.Println("yaaaay user deleted ", de.UserID)
	}
	return nil
}
