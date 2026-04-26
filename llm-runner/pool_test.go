package llm_runner

import (
	"context"
	"testing"
)

func TestNewPool(t *testing.T) {
	p := NewPool(nil)
	if p == nil {
		t.Fatal("NewPool(nil) не должен возвращать nil")
	}

	runners := p.GetRunners(context.Background())
	if len(runners) != 0 {
		t.Errorf("новый пул должен иметь 0 раннеров, получено %d", len(runners))
	}

	if p.HasActiveRunners() {
		t.Error("новый пул не должен иметь активных раннеров")
	}
}

func TestNewPool_withAddresses(t *testing.T) {
	p := NewPool([]string{"a:1", "b:2"})
	runners := p.GetRunners(context.Background())
	if len(runners) != 2 {
		t.Errorf("ожидалось 2 раннера, получено %d", len(runners))
	}

	if !p.HasActiveRunners() {
		t.Error("ожидалось наличие активных раннеров")
	}
}

func TestPool_Add_Remove(t *testing.T) {
	p := NewPool([]string{"a:1"})
	p.Add("b:2")

	runners := p.GetRunners(context.Background())
	if len(runners) != 2 {
		t.Errorf("после Add: ожидалось 2 раннера, получено %d", len(runners))
	}

	p.Remove("a:1")

	runners = p.GetRunners(context.Background())
	if len(runners) != 1 {
		t.Errorf("после Remove: ожидался 1 раннер, получено %d", len(runners))
	}

	if runners[0].Address != "b:2" {
		t.Errorf("оставшийся адрес: %s", runners[0].Address)
	}
}

func TestPool_Add_emptyIgnored(t *testing.T) {
	p := NewPool(nil)
	p.Add("")

	if len(p.GetRunners(context.Background())) != 0 {
		t.Error("Add с пустым адресом не должен добавлять")
	}
}

func TestPool_SetRunnerEnabled(t *testing.T) {
	p := NewPool([]string{"a:1"})
	if !p.HasActiveRunners() {
		t.Error("ожидался активный раннер до отключения")
	}

	p.SetRunnerEnabled("a:1", false)
	if p.HasActiveRunners() {
		t.Error("ожидалось отсутствие активных после отключения")
	}

	runners := p.GetRunners(context.Background())
	if len(runners) != 1 || runners[0].Enabled {
		t.Errorf("раннер должен быть отключён: %+v", runners)
	}

	p.SetRunnerEnabled("a:1", true)
	if !p.HasActiveRunners() {
		t.Error("ожидался активный раннер после включения")
	}
}

func TestPool_CheckConnection_noRunners(t *testing.T) {
	p := NewPool(nil)
	ok, err := p.CheckConnection(context.Background())
	if err == nil {
		t.Error("ожидалась ошибка при отсутствии раннеров")
	}

	if ok {
		t.Error("ожидалось ok == false при отсутствии раннеров")
	}
}

func TestPool_orderedAddresses(t *testing.T) {
	p := NewPool([]string{
		"a:1",
		"b:2",
		"c:3",
	})

	order := p.orderedAddresses("b:2")
	if len(order) != 3 {
		t.Fatalf("ожидалось 3 адреса, получено %d", len(order))
	}
	if order[0] != "b:2" {
		t.Errorf("первый адрес должен быть b:2, получено %s", order[0])
	}

	seen := make(map[string]bool)
	for _, a := range order {
		seen[a] = true
	}
	if !seen["a:1"] || !seen["b:2"] || !seen["c:3"] {
		t.Errorf("ожидались a:1, b:2, c:3, получено %v", order)
	}
}

func TestPool_orderedAddresses_emptyStart(t *testing.T) {
	p := NewPool([]string{"a:1", "b:2"})
	order := p.orderedAddresses("")
	if len(order) != 2 {
		t.Errorf("ожидалось 2 адреса при пустом startAddr, получено %d", len(order))
	}
}
