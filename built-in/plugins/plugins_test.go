package plugins

import (
	"minimal/minimal-core/built-in/messaging"
	"testing"
)

func TestEmpty(t *testing.T) {
    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    m.AddOutput(to)

    r := NewRegistry(m)
    r.Setup()

    m.Close()
    to.CheckMessages(t, []messaging.Message{})
}

type okPlugin struct {
    ok bool
}

func (o *okPlugin) Init(*plugins) {
    o.ok = true
}

func TestAdd(t *testing.T) {
    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    m.AddOutput(to)

    r := NewRegistry(m)

    o := &okPlugin{}
    r.Add(o)

    r.Setup()

    if !o.ok {
        t.Error("Expected ok to be true")
    }

    m.Close()
    to.CheckMessages(t, []messaging.Message{})
}

type firstGetPlugin struct {
    t *testing.T
    ok bool
}

func (f *firstGetPlugin) Init(p *plugins) {
    p2, ok := Get[*secondGetPlugin](p)

    if !ok {
        f.t.Fatal("Could not retrieve *secondGetPlugin")
    }

    p2.ok = true
}

type secondGetPlugin struct {
    t *testing.T
    ok bool
}

func (s *secondGetPlugin) Init(p *plugins) {
    p1, ok := Get[*firstGetPlugin](p)

    if !ok {
        s.t.Fatal("Could not retrieve *firstGetPlugin")
    }

    p1.ok = true
}

func TestGet(t *testing.T) {
    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    m.AddOutput(to)

    r := NewRegistry(m)

    p1 := &firstGetPlugin{t: t}
    r.Add(p1)

    p2 := &secondGetPlugin{t: t}
    r.Add(p2)

    r.Setup()

    if !p1.ok {
        t.Error("Expected p1.ok to be true")
    }

    if !p2.ok {
        t.Error("Expected p2.ok to be true")
    }

    m.Close()
    to.CheckMessages(t, []messaging.Message{})
}

type selfRetrievingPlugin struct {
    t *testing.T
    ok bool
}

func (s *selfRetrievingPlugin) Init(p *plugins) {
    me, ok := Get[*selfRetrievingPlugin](p)

    if !ok {
        s.t.Fatal("Could not retrieve *selfRetrievingPlugin")
    }

    if s != me {
        s.t.Fatal("s does not reference the same memory as me")
    }

    me.ok = true
}

func TestGetItself(t *testing.T) {
    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    m.AddOutput(to)

    r := NewRegistry(m)

    s := &selfRetrievingPlugin{t: t}
    r.Add(s)

    r.Setup()

    if !s.ok {
        t.Error("Expected ok to be true")
    }

    m.Close()
    to.CheckMessages(t, []messaging.Message{})
}

type missingPlugin struct {
    t *testing.T
}

func (m *missingPlugin) Init(p *plugins) {
    _, ok := Get[*okPlugin](p)

    if ok {
        m.t.Fatal("Expected ok to be false")
    }
}

func TestMissing(t *testing.T) {
    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    m.AddOutput(to)

    r := NewRegistry(m)

    mp := &missingPlugin{t: t}
    r.Add(mp)

    r.Setup()

    m.Close()
    to.CheckMessages(t, []messaging.Message{})
}

func TestDuplicate(t *testing.T) {
    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    m.AddOutput(to)

    r := NewRegistry(m)

    o1 := &okPlugin{}
    r.Add(o1)

    o2 := &okPlugin{}
    r.Add(o2)

    r.Setup()

    if o1.ok {
        t.Error("Expected o1.ok to be false")
    }

    if !o2.ok {
        t.Error("Expected o2.ok to be true")
    }

    m.Close()
    to.CheckMessages(
        t,
        []messaging.Message{
            {
                Message: "Duplicate plugin declaration",
                Severity: messaging.Error,
                Notes: []string{
                    "Plugin type: *plugins.okPlugin",
                    "The plugin in the registry will be overwritten",
                },
            },
        },
    )
}

type valuePlugin struct {}

func (valuePlugin) Init(*plugins) {}

func TestNonReference(t *testing.T) {
    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    m.AddOutput(to)

    r := NewRegistry(m)

    v := valuePlugin{}
    r.Add(v)

    r.Setup()

    m.Close()
    to.CheckMessages(
        t,
        []messaging.Message{
            {
                Message: "Plugin should be a pointer",
                Severity: messaging.Error,
                Notes: []string{
                    "Plugin type: plugins.valuePlugin",
                    "The plugin wont be added to the registry",
                },
            },
        },
    )
}

type pointerPlugin struct {
    t *testing.T
}

func (p  *pointerPlugin) Init(plugins *plugins) {
    _, ok := Get[pointerPlugin](plugins)

    if ok {
        p.t.Error("Expected ok to be false")
    }
}

func TestRetrieveNonReference(t *testing.T) {
    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    m.AddOutput(to)

    r := NewRegistry(m)

    p := pointerPlugin{}
    r.Add(&p)

    r.Setup()

    m.Close()
    to.CheckMessages(
        t,
        []messaging.Message{
            {
                Message: "Cannot retrieve a non pointer plugin since all plugins are pointers",
                Severity: messaging.Error,
                Notes: []string{"Requested type: plugins.pointerPlugin"},
            },
        },
    )
}
