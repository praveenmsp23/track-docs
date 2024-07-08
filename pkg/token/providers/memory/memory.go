package memory

import (
	"container/list"
	"sync"
	"time"

	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/token/base"
)

var pder = &MemoryProvider{list: list.New()}

func GetProvider(cfg *config.Config) (*MemoryProvider, error) {
	pder.Tokens = make(map[string]*list.Element, 0)
	return pder, nil
}

type MemoryTokenStore struct {
	tid          string
	timeAccessed time.Time
	value        map[string]string
}

func (st *MemoryTokenStore) Set(key, value string) error {
	st.value[key] = value
	pder.TokenUpdate(st.tid)
	return nil
}

func (st *MemoryTokenStore) GetAll() (map[string]string, error) {
	pder.TokenUpdate(st.tid)
	return st.value, nil
}

func (st *MemoryTokenStore) Get(key string) (string, bool) {
	pder.TokenUpdate(st.tid)
	if v, ok := st.value[key]; ok {
		return v, true
	}
	return "", false
}

func (st *MemoryTokenStore) Delete(key string) error {
	delete(st.value, key)
	pder.TokenUpdate(st.tid)
	return nil
}

func (st *MemoryTokenStore) TokenID() string {
	return st.tid
}

type MemoryProvider struct {
	lock   sync.Mutex
	Tokens map[string]*list.Element
	list   *list.List
}

func (pder *MemoryProvider) TokenInit(tid string) (base.Token, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[string]string, 0)
	newsess := &MemoryTokenStore{tid: tid, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newsess)
	pder.Tokens[tid] = element
	return newsess, nil
}

func (pder *MemoryProvider) TokenRead(tid string) (base.Token, error) {
	if element, ok := pder.Tokens[tid]; ok {
		return element.Value.(*MemoryTokenStore), nil
	} else {
		return nil, nil
	}
}

func (pder *MemoryProvider) TokenDestroy(tid string) error {
	if element, ok := pder.Tokens[tid]; ok {
		delete(pder.Tokens, tid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *MemoryProvider) TokenGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*MemoryTokenStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.Tokens, element.Value.(*MemoryTokenStore).tid)
		} else {
			break
		}
	}
}

func (pder *MemoryProvider) TokenUpdate(tid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.Tokens[tid]; ok {
		element.Value.(*MemoryTokenStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}
