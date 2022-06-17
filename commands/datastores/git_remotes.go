package datastores

import (
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

var GitRemotes = gitRemotes{
	remotes: map[string]*gitRemote{},
}

type gitRemote struct {
	uri      string
	initOnce sync.Once
	refs     []*plumbing.Reference
}

type gitRemotes struct {
	mutex   sync.Mutex
	remotes map[string]*gitRemote
}

func (g *gitRemotes) List(uri string) ([]*plumbing.Reference, error) {
	g.mutex.Lock()
	remote, ok := g.remotes[uri]
	if !ok {
		remote = &gitRemote{
			uri: uri,
		}
		g.remotes[uri] = remote
	}
	g.mutex.Unlock()
	return remote.list()
}

func (g *gitRemote) list() ([]*plumbing.Reference, error) {
	var err error
	g.initOnce.Do(func() { err = g.init() })
	return g.refs, err
}

func (g *gitRemotes) Reset() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.remotes = map[string]*gitRemote{}
}

func (g *gitRemote) init() (err error) {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{g.uri},
	})
	g.refs, err = rem.List(&git.ListOptions{})
	return
}
