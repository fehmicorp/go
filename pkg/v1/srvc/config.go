package srvc

type Developer struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
}

type Environment struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	GroupName string `json:"group_name"`
	Desc      string `json:"description"`
}

type Packages struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Title         string    `json:"title"`
	Desc          string    `json:"desc"`
	Version       string    `json:"version"`
	Repo          string    `json:"repo"`
	WorkDir       string    `json:"work_dir"`
	Developer     Developer `json:"developers"`
	SupportedOS   []string  `json:"supported_os"`
	SupportedArch []string  `json:"supported_arch"`
}

type Services struct {
	Id            int            `json:"id"`
	Name          string         `json:"name"`
	Title         string         `json:"title"`
	Desc          string         `json:"desc"`
	Tags          []string       `json:"tags"`
	Port          int            `json:"port"`
	Package       Packages       `json:"packages"`
	Environment   Environment    `json:"environments"`
	RuntimeType   string         `json:"runtime_type"`
	RuntimeConfig map[string]any `json:"runtime_config"`
	Status        Status         `json:"status"`
}

type Status struct {
	Active  bool `json:"active"`
	Running bool `json:"running"`
	Startup bool `json:"startup"`
}

// DBService mirrors the joined Supabase response schema for unmarshalling
type dbServiceResponse struct {
	Id            int            `json:"id"`
	Name          string         `json:"name"`
	Title         string         `json:"title"`
	Desc          string         `json:"desc"`
	Tags          []string       `json:"tags"`
	Port          int            `json:"port"`
	RuntimeType   string         `json:"runtime_type"`
	RuntimeConfig map[string]any `json:"runtime_config"`
	Packages      struct {
		Id            int       `json:"id"`
		Name          string    `json:"name"`
		Title         string    `json:"title"`
		Desc          string    `json:"desc"`
		Version       string    `json:"version"`
		Repo          string    `json:"repo"`
		WorkDir       string    `json:"work_dir"`
		SupportedOS   []string  `json:"supported_os"`
		SupportedArch []string  `json:"supported_arch"`
		Developers    Developer `json:"developers"`
	} `json:"packages"`
	Environments Environment `json:"environments"`
}
