package utils

import "io/ioutil"

func WriteAppRun(file, appid, destdir string) error {
	if len(file) == 0 {
		file = `#!/bin/sh
HERE="$(dirname "$(readlink -f "${0}")")"
export UNION_PRELOAD="${HERE}"
export LD_PRELOAD="${HERE}/libunionpreload.so"
export PATH="${HERE}"/usr/bin/:"${HERE}"/usr/sbin/:"${HERE}"/usr/games/:"${HERE}"/bin/:"${HERE}"/sbin/:"${PATH}"
export LD_LIBRARY_PATH="${HERE}"/usr/lib/:"${HERE}"/usr/lib/i386-linux-gnu/:"${HERE}"/usr/lib/x86_64-linux-gnu/:"${HERE}"/usr/lib32/:"${HERE}"/usr/lib64/:"${HERE}"/lib/:"${HERE}"/lib/i386-linux-gnu/:"${HERE}"/lib/x86_64-linux-gnu/:"${HERE}"/lib32/:"${HERE}"/lib64/:"${LD_LIBRARY_PATH}"
export PYTHONPATH="${HERE}"/usr/share/pyshared/:"${PYTHONPATH}"
export PYTHONHOME="${HERE}"/usr/
export XDG_DATA_DIRS="${HERE}"/usr/share/:"${XDG_DATA_DIRS}"
export PERLLIB="${HERE}"/usr/share/perl5/:"${HERE}"/usr/lib/perl5/:"${PERLLIB}"
export GSETTINGS_SCHEMA_DIR="${HERE}"/usr/share/glib-2.0/schemas/:"${GSETTINGS_SCHEMA_DIR}"
export QT_PLUGIN_PATH="${HERE}"/usr/lib/qt4/plugins/:"${HERE}"/usr/lib/i386-linux-gnu/qt4/plugins/:"${HERE}"/usr/lib/x86_64-linux-gnu/qt4/plugins/:"${HERE}"/usr/lib32/qt4/plugins/:"${HERE}"/usr/lib64/qt4/plugins/:"${HERE}"/usr/lib/qt5/plugins/:"${HERE}"/usr/lib/i386-linux-gnu/qt5/plugins/:"${HERE}"/usr/lib/x86_64-linux-gnu/qt5/plugins/:"${HERE}"/usr/lib32/qt5/plugins/:"${HERE}"/usr/lib64/qt5/plugins/:"${QT_PLUGIN_PATH}"
EXEC=$(grep -e '^Exec=.*' "${HERE}"/*.desktop | head -n 1 | cut -d "=" -f 2- | sed -e 's|%.||g')
exec ${EXEC} "$@"`
	}

	return ioutil.WriteFile(destdir+"/AppRun", []byte(file), 0755)

}

func WriteDesktop(file, appid, destdir string) error {
	if len(file) == 0 {
		file = `[Desktop Entry]
Name=` + appid + `
Exec=` + appid + `
Icon=` + appid + `
Type=Application
Comment=Some basic comment
Categories=GNOME;
`

	}

	return ioutil.WriteFile(destdir+"/"+appid+".desktop", []byte(file), 0755)

}
