package apt

import "strings"

const EXCLUDE = `
apt
apt-transport-https
dbus
debconf
dictionaries-common
dpkg
fontconfig
fontconfig-config
gvfs-backends
gksu
glib-networking
gstreamer1.0-plugins-base
gstreamer1.0-plugins-good
gstreamer1.0-plugins-ugly
gstreamer1.0-pulseaudio
gtk2-engines-pixbuf
kde-runtime
libasound2
libatk1.0-0
libc6
libc6-dev
libcairo2
libcups2
libdbus-1-3
libdrm2
libegl1-mesa
libfontconfig1
libgbm1
libgcc1
libgdk-pixbuf2.0-0
libgl1
libgl1-mesa
libgl1-mesa-dri
libgl1-mesa-glx
libglu1-mesa
libgpg-error0
libgtk2.0-0
libgtk-3-0
libnss3
libpango1.0-0
libpango-1.0-0
libpangocairo-1.0-0
libpangoft2-1.0-0
libstdc++6
libtasn1-6
libwayland-egl1-mesa
lsb-base
libxcb1
mime-support
passwd
udev
uuid-runtime
`

func isExcluded(id string) bool {
	for _, exc := range strings.Split(EXCLUDE, "\n") {
		if id == exc {
			return true
		}
	}
	return false
}