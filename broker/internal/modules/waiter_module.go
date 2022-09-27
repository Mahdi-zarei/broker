package modules

// WaiterManager struct is NOT THREAD SAFE!! you should run it using a sync.Mutex to prevent unexpected results
type WaiterManager struct {
	notifiers []chan struct{}
}

func NewWaiterManager() *WaiterManager {
	return &WaiterManager{
		notifiers: nil,
	}
}

// FreeWaiters closes all channels and empties the notifiers
func (w *WaiterManager) FreeWaiters() {
	for _, ch := range w.notifiers {
		close(ch)
	}
	w.notifiers = nil
}

// RegisterWaiter returns a channel that closes once the waiters are set free
func (w *WaiterManager) RegisterWaiter() chan struct{} {
	ch := make(chan struct{})
	w.notifiers = append(w.notifiers, ch)
	return ch
}
