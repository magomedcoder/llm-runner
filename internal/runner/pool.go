package runner

import (
	"context"
	"sync"

	"github.com/magomedcoder/gen/internal/domain"
	"github.com/magomedcoder/gen/internal/repository"
	"github.com/magomedcoder/gen/pkg/logger"
)

type Pool struct {
	registry     *Registry
	defaultModel string
	repos        map[string]*repository.LLMRunnerRepository
	mu           sync.RWMutex
}

func NewPool(registry *Registry, defaultModel string) *Pool {
	return &Pool{
		registry:     registry,
		defaultModel: defaultModel,
		repos:        make(map[string]*repository.LLMRunnerRepository),
	}
}

func (p *Pool) getRepo(ctx context.Context, addr string) (domain.LLMRepository, error) {
	p.mu.RLock()
	repo, ok := p.repos[addr]
	p.mu.RUnlock()
	if ok && repo != nil {
		return repo, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if repo, ok := p.repos[addr]; ok && repo != nil {
		return repo, nil
	}

	repo, err := repository.NewLLMRunnerRepository(addr, p.defaultModel)
	if err != nil {
		logger.W("Pool: не удалось подключиться к раннеру %s: %v", addr, err)
		return nil, err
	}
	logger.I("Pool: подключение к раннеру %s установлено", addr)
	p.repos[addr] = repo

	return repo, nil
}

func (p *Pool) HasConnection(addr string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.repos[addr]
	return ok
}

func (p *Pool) CloseAddr(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if repo, ok := p.repos[addr]; ok && repo != nil {
		_ = repo.Close()
		delete(p.repos, addr)
	}
}

func (p *Pool) getOneEnabled(ctx context.Context) (domain.LLMRepository, error) {
	addrs := p.registry.GetEnabledAddresses()
	if len(addrs) == 0 {
		logger.W("Pool: нет доступных раннеров")
		return nil, domain.ErrNoRunners
	}

	return p.getRepo(ctx, addrs[0])
}

var _ domain.LLMRepository = (*Pool)(nil)

func (p *Pool) CheckConnection(ctx context.Context) (bool, error) {
	repo, err := p.getOneEnabled(ctx)
	if err != nil {
		return false, nil
	}

	return repo.CheckConnection(ctx)
}

func (p *Pool) GetModels(ctx context.Context) ([]string, error) {
	repo, err := p.getOneEnabled(ctx)
	if err != nil {
		return nil, err
	}

	return repo.GetModels(ctx)
}

func (p *Pool) SendMessage(ctx context.Context, sessionID string, model string, messages []*domain.Message) (chan string, error) {
	repo, err := p.getOneEnabled(ctx)
	if err != nil {
		return nil, err
	}

	return repo.SendMessage(ctx, sessionID, model, messages)
}

func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, repo := range p.repos {
		_ = repo.Close()
	}

	p.repos = make(map[string]*repository.LLMRunnerRepository)

	return nil
}
