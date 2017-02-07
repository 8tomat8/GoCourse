package main

import (
	"time"
	"fmt"
)

func runMetrics(every *int) {
	stop := make(chan empty)
	go func(stop chan empty) {
		ticker := time.NewTicker(time.Duration(*every) * time.Second)
		var now time.Time
		var clients uint
		var rooms uint
		var messages uint
		for {
			clients = 0
			rooms = 0
			messages = 0
			select {
			case now = <-ticker.C:
				fmt.Printf("%d users in lobby\n", len(Chat.lobby.users))
				for _, r := range Chat.rooms {
					rooms += 1
					messages += r.sentMessages
					fmt.Println(now)
					fmt.Printf("Room %s has %d users online. List of users:\n", r.name, len(r.users))
					for name := range r.users {
						clients += 1
						fmt.Printf("\t-%s\n", name)
					}
					fmt.Printf("%d messages sent\n", r.sentMessages)
				}
				fmt.Printf("Total \n\tClients online: %d\n\tActive rooms: %d\n\tMessages sent: %d\n", clients, rooms, messages)
			case <-stop:
				return
			}
		}
	}(stop)
}
