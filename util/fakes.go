package util

//**Fake Mutex**

//go:generate counterfeiter -o ../fakes/locker/locker.go . Locker

type Locker interface {
	Lock()
	Unlock()
}
