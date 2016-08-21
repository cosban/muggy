package messages

import (
	"log"
	"time"

	"github.com/sorcix/irc"
)

var (
	quit     = make(chan struct{})
	queueing = false
	msgQueue queue
)

//BuildMessages builds one or more PRIVMSGs to be queued
func BuildMessages(s irc.Sender, command string, params []string, messages ...string) {
	for _, m := range messages {
		log.Printf("ADDING: %s", m)
		QueueMessages(s,
			&irc.Message{
				Command:  command,
				Params:   params,
				Trailing: m,
			},
		)
	}
}

// QueueMessages adds one or more messages to the message queue
func QueueMessages(s irc.Sender, messages ...*irc.Message) {
	for _, m := range messages {
		msgQueue = append(msgQueue, queuedMessage{
			Sender:  s,
			Message: m,
		})
	}
}

// WriteLoop is a timed task that writes irc messages
func WriteLoop() {
	if !queueing {
		queueing = true
		ticker := time.NewTicker(200 * time.Millisecond)
		go func() {
			for {
				select {
				case <-ticker.C:
					{
						if !msgQueue.IsEmpty() {
							result := msgQueue.Remove(0)
							result.Sender.Send(result.Message)
						}
					}
				case <-quit:
					{
						ticker.Stop()
						return
					}
				}
			}
		}()
	}
}

// StopLoop causes the WriteLoop to quit
func StopLoop() {
	close(quit)
}
