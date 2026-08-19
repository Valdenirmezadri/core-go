package helperecho

import (
	"sync"

	"github.com/labstack/echo/v4"
)

// Mux permite que handlers instanciados por tenant registrem rotas com API
// echo-like (r.GET("/x", h.Fetch)) sem colidir entre si: cada rota entra no
// echo uma única vez, e a cada request o resolve extrai o ID da instância dona
// (ex: tenant do contexto) e o request é despachado pro handler daquela
// instância.
//
// Rotas podem ser registradas antes do Attach (tenants sobem antes do echo):
// ficam em buffer e são conectadas quando o Attach rodar.
type Mux struct {
	mu      sync.RWMutex
	g       *echo.Group
	mid     []echo.MiddlewareFunc
	resolve func(echo.Context) (uint, error)
	// routes: "GET /path" -> instanceID -> handler
	routes map[string]map[uint]echo.HandlerFunc
	// pending: rotas coletadas antes do Attach
	pending []pendingRoute
}

type pendingRoute struct {
	method, path string
	mid          []echo.MiddlewareFunc
}

func NewMux() *Mux {
	return &Mux{routes: map[string]map[uint]echo.HandlerFunc{}}
}

// Attach conecta o mux num grupo echo. resolve extrai o ID da instância do
// request; mid roda em toda rota do mux (ex: resolução do tenant), antes dos
// middlewares por rota. Rotas já em buffer são registradas agora.
func (m *Mux) Attach(g *echo.Group, resolve func(echo.Context) (uint, error), mid ...echo.MiddlewareFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.g = g
	m.resolve = resolve
	m.mid = mid
	for _, p := range m.pending {
		m.register(p.method, p.path, p.mid)
	}
	m.pending = nil
}

// For retorna a visão do mux pra uma instância: API echo-like cujos registros
// pertencem àquela instância. Chamar Close() quando a instância parar.
func (m *Mux) For(id uint) *Router {
	return &Router{mux: m, id: id}
}

// register precisa do lock; registra a rota no echo com um dispatcher que
// resolve a instância por request.
func (m *Mux) register(method, path string, mid []echo.MiddlewareFunc) {
	key := method + " " + path
	m.g.Add(method, path, func(c echo.Context) error {
		id, err := m.resolve(c)
		if err != nil {
			return err
		}
		m.mu.RLock()
		h := m.routes[key][id]
		m.mu.RUnlock()
		if h == nil {
			return echo.ErrNotFound
		}
		return h(c)
	}, append(append([]echo.MiddlewareFunc{}, m.mid...), mid...)...)
}

// add registra o handler da instância id; a rota entra no echo só na primeira
// vez que o path aparece — os middlewares por rota valem os do primeiro
// registrante (são idênticos entre instâncias, vêm do mesmo código).
func (m *Mux) add(method, path string, id uint, h echo.HandlerFunc, mid []echo.MiddlewareFunc) {
	key := method + " " + path
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.routes[key] == nil {
		m.routes[key] = map[uint]echo.HandlerFunc{}
		if m.g != nil {
			m.register(method, path, mid)
		} else {
			m.pending = append(m.pending, pendingRoute{method: method, path: path, mid: mid})
		}
	}
	m.routes[key][id] = h
}

// Remove descarta todos os handlers da instância (ex: tenant parou).
func (m *Mux) Remove(id uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, byID := range m.routes {
		delete(byID, id)
	}
}

// Router é a visão do Mux pra uma instância, com API echo-like.
type Router struct {
	mux *Mux
	id  uint
}

func (r *Router) GET(path string, h echo.HandlerFunc, mid ...echo.MiddlewareFunc) {
	r.mux.add(echo.GET, path, r.id, h, mid)
}

func (r *Router) POST(path string, h echo.HandlerFunc, mid ...echo.MiddlewareFunc) {
	r.mux.add(echo.POST, path, r.id, h, mid)
}

func (r *Router) PUT(path string, h echo.HandlerFunc, mid ...echo.MiddlewareFunc) {
	r.mux.add(echo.PUT, path, r.id, h, mid)
}

func (r *Router) PATCH(path string, h echo.HandlerFunc, mid ...echo.MiddlewareFunc) {
	r.mux.add(echo.PATCH, path, r.id, h, mid)
}

func (r *Router) DELETE(path string, h echo.HandlerFunc, mid ...echo.MiddlewareFunc) {
	r.mux.add(echo.DELETE, path, r.id, h, mid)
}

// Close remove os handlers da instância do mux.
func (r *Router) Close() {
	r.mux.Remove(r.id)
}
