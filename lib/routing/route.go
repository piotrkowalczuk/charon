package routing

// RoutePattern ...
type RoutePattern string

// RouteName ...
type RouteName string

// Routes ...
type Routes map[RouteName]RoutePattern

// String ...
func (rp RoutePattern) String() string {
	return string(rp)
}

// String ...
func (rn RouteName) String() string {
	return string(rn)
}

func NewRoutes(schema map[string]string) Routes {
	routes := Routes{}

	for name, pattern := range schema {
		routes.Add(RouteName(name), RoutePattern(pattern))
	}

	return routes
}

// Add ...
func (r *Routes) Add(name RouteName, pattern RoutePattern) {
	(*r)[name] = pattern
}

// GetPattern ...
func (r *Routes) GetPattern(name RouteName) RoutePattern {
	return (*r)[name]
}
