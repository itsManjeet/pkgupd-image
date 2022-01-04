package apt

import (
	"log"
)

var (
	APPLIST      []*Package
	ALREADY_DONE []string
)

func in_list(pkg *Package) bool {
	for _, p := range APPLIST {
		if p.Name == pkg.Name {
			return true
		}
	}
	return false
}

func is_done(pkg_id string) bool {
	if isExcluded(pkg_id) {
		return true
	}
	for _, p := range ALREADY_DONE {
		if p == pkg_id {
			return true
		}
	}
	return false
}

func (apt Apt) resolve(pkg *Package) bool {
	log.Println("found dependency", pkg.Name)
	for _, p := range pkg.GetDepends() {
		if is_done(p) {
			continue
		}

		ALREADY_DONE = append(ALREADY_DONE, p)

		pkg_dep := apt.Get(p)
		if pkg_dep == nil {
			log.Println(p + " is missing, required by " + pkg.Name)
			return false
		}

		if !apt.resolve(pkg_dep) {
			return false
		}

		if !in_list(pkg_dep) {
			APPLIST = append(APPLIST, pkg_dep)
		}
	}

	if !in_list(pkg) {
		APPLIST = append(APPLIST, pkg)
	}

	return true
}

func (apt Apt) packagesList(pkgs ...string) []*Package {
	APPLIST = make([]*Package, 0)
	ALREADY_DONE = make([]string, 0)

	for _, inc := range pkgs {
		inc_pkg := apt.Get(inc)
		if inc_pkg == nil {
			log.Println("missing include package " + inc)
			return nil
		}

		if !apt.resolve(inc_pkg) {
			return nil
		}
	}

	for _, a := range APPLIST {
		log.Println(a)
	}

	return APPLIST
}
