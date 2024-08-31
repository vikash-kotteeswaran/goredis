package core

import (
	"syscall"
)

type Epoll struct {
	fd                 int
	subscribed_kevents []syscall.Kevent_t
	subscribed_events  []Event
}

type Event struct {
	// Fd denotes the file descriptor
	Fd int
	// Op denotes the operations on file descriptor that are to be monitored
	Op Operations
}

type Operations uint32

func GetMultiplexer(maxConnections int) (*Epoll, error) {
	fd, err := syscall.Kqueue()

	if err != nil {
		return nil, err
	}

	return &Epoll{
		fd:                 fd,
		subscribed_kevents: make([]syscall.Kevent_t, maxConnections),
		subscribed_events:  make([]Event, maxConnections),
	}, nil
}

func (ep *Epoll) Subscribe(fd int) error {
	event := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD,
	}

	subscribed, err := syscall.Kevent(ep.fd, []syscall.Kevent_t{event}, nil, nil)

	if err != nil || subscribed == -1 {
		return err
	}

	// fmt.Println("File Descriptor subscribed to Multiplexer :: " + strconv.Itoa(fd))

	return nil
}

func (ep *Epoll) UnSubscribe(fd int) error {
	event := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_DELETE,
	}

	subscribed, err := syscall.Kevent(ep.fd, []syscall.Kevent_t{event}, nil, nil)

	if err != nil || subscribed == -1 {
		return err
	}

	// fmt.Println("File Descriptor unsubscribed from Multiplexer :: " + strconv.Itoa(fd))

	return nil
}

func (ep *Epoll) Poll() ([]Event, error) {
	nEvents, err := syscall.Kevent(ep.fd, nil, ep.subscribed_kevents, nil)

	if err != nil {
		return nil, err
	}

	ep.subcribedKeventsToEvents(nEvents)

	return ep.subscribed_events[:nEvents], nil
}

// Helpers

func (ep *Epoll) subcribedKeventsToEvents(nEvents int) {
	for idx := 0; idx < nEvents; idx++ {
		ep.subscribed_events[idx] = Event{
			Fd: int(ep.subscribed_kevents[idx].Ident),
			Op: Operations(ep.subscribed_kevents[idx].Flags),
		}
	}
}
