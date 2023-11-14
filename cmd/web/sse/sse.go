package sse

import "fmt"



type ServerEvent struct {
	Event string
	Data string 
}

func (s *ServerEvent) ToString() (d string) {
	d = fmt.Sprintf("event: %s\ndata: %s\n\n",s.Event,s.Data)
	return
}

type EventList []ServerEvent