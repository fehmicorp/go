package srvc

import (
	"fmt"
	"runtime"

	supabase "github.com/supabase-community/supabase-go"
)

func FetchFilteredServices(supabaseURL, supabaseKey, targetOS, targetArch string) ([]Services, error) {
	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize supabase client: %w", err)
	}

	var dbRecords []dbServiceResponse

	// Base query joining packages, developers, and environments
	query := client.From("services").
		Select("id, name, title, desc, tags, port, runtime_type, runtime_config, packages(*, developers(*)), environments(*)", "", false)

	_, err = query.ExecuteTo(&dbRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch services data: %w", err)
	}

	// Auto-detect local machine values if arguments are empty
	if targetOS == "" {
		targetOS = runtime.GOOS // e.g., "windows", "linux", "darwin"
	}
	if targetArch == "" {
		targetArch = runtime.GOARCH // e.g., "amd64", "arm64"
	}

	var resultList []Services
	for _, rec := range dbRecords {
		// Filter match validation for SupportedOS and SupportedArch arrays
		osMatch := false
		if len(rec.Packages.SupportedOS) == 0 {
			osMatch = true // Default allow if unrestricted
		} else {
			for _, osName := range rec.Packages.SupportedOS {
				if osName == targetOS {
					osMatch = true
					break
				}
			}
		}

		archMatch := false
		if len(rec.Packages.SupportedArch) == 0 {
			archMatch = true // Default allow if unrestricted
		} else {
			for _, archName := range rec.Packages.SupportedArch {
				if archName == targetArch {
					archMatch = true
					break
				}
			}
		}

		// Keep only compatible services
		if osMatch && archMatch {
			srv := Services{
				Id:            rec.Id,
				Name:          rec.Name,
				Title:         rec.Title,
				Desc:          rec.Desc,
				Tags:          rec.Tags,
				Port:          rec.Port,
				RuntimeType:   rec.RuntimeType,
				RuntimeConfig: rec.RuntimeConfig,
				Package: Packages{
					Id:            rec.Packages.Id,
					Name:          rec.Packages.Name,
					Title:         rec.Packages.Title,
					Desc:          rec.Packages.Desc,
					Version:       rec.Packages.Version,
					Repo:          rec.Packages.Repo,
					WorkDir:       rec.Packages.WorkDir,
					Developer:     rec.Packages.Developers,
					SupportedOS:   rec.Packages.SupportedOS,
					SupportedArch: rec.Packages.SupportedArch,
				},
				Environment: rec.Environments,
				Status: Status{
					Active:  false,
					Running: false,
					Startup: false,
				},
			}
			resultList = append(resultList, srv)
		}
	}

	return resultList, nil
}
