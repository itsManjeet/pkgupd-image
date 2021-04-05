package archlinux

import "strings"

type ArchPkg struct {
	Filename    string
	Name        string
	Version     string
	Description string
	Provides    []string
	Depends     []string
}

func (a *ArchLinux) parseDatabase(content string) (map[string]ArchPkg, error) {
	archdb := make(map[string]ArchPkg, 0)

	striper := func(f string) string {
		if strings.Contains(f, ">=") {
			return strings.Split(f, ">=")[0]
		}
		if strings.Contains(f, "=") {
			return strings.Split(f, "=")[0]
		}
		return f
	}

	arr := strings.Split(content, "\n")
	i := 0
	for arr[i] == "%FILENAME%" {
		i++

		var t ArchPkg
		t.Filename = arr[i]
		i++

		for arr[i] != "%FILENAME%" {
			switch arr[i] {
			case "%NAME%":
				i++
				t.Name = arr[i]

			case "%VERSION%":
				i++
				t.Version = arr[i]

			case "%DEPENDS%":
				i++
				t.Depends = make([]string, 0)
				for arr[i] != "" && arr[i][0] != '%' {
					t.Depends = append(t.Depends, striper(arr[i]))
					i++
				}

			case "%PROVIDES%":
				i++
				t.Provides = make([]string, 0)
				for arr[i] != "" && arr[i][0] != '%' {
					t.Provides = append(t.Provides, striper(arr[i]))
					i++
				}
			}
			i++
			if i >= len(arr) {
				break
			}
		}
		archdb[t.Name] = t
		if i >= len(arr) {
			break
		}

	}

	return archdb, nil
}
