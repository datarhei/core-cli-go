package api

type IAMUser struct {
	Name      string      `json:"name"`
	Superuser bool        `json:"superuser"`
	Auth      IAMUserAuth `json:"auth"`
	Policies  []IAMPolicy `json:"policies"`
}

type IAMUserAuth struct {
	API      IAMUserAuthAPI      `json:"api"`
	Services IAMUserAuthServices `json:"services"`
}

type IAMUserAuthAPI struct {
	Password string              `json:"userpass"`
	Auth0    IAMUserAuthAPIAuth0 `json:"auth0"`
}

type IAMUserAuthAPIAuth0 struct {
	User   string         `json:"user"`
	Tenant IAMAuth0Tenant `json:"tenant"`
}

type IAMUserAuthServices struct {
	Basic   []string `json:"basic"`
	Token   []string `json:"token"`
	Session []string `json:"session"`
}

type IAMAuth0Tenant struct {
	Domain   string `json:"domain"`
	Audience string `json:"audience"`
	ClientID string `json:"client_id"`
}

type IAMPolicy struct {
	Name     string   `json:"name,omitempty"`
	Domain   string   `json:"domain"`
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}