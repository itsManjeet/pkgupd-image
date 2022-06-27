package apt

import "strings"

const EXCLUDE = `
adduser
apt
apt-transport-https
dbus
dconf-gsettings-backend
dconf-service
dbus-user-session
debconf
dictionaries-common
dpkg
fontconfig
fontconfig-config
gcc-10-base
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
libblkid1
libc6
libc6-dev
libcairo2
libcap2
libcap2-bin
libcups2
libdbus-1-3
libdrm2
libegl1-mesa
libffi7
libfontconfig1
libgbm1
libgcc1
libgcc-s1
libgdk-pixbuf2.0-0
libgl1
libgl1-mesa
libgl1-mesa-dri
libgl1-mesa-glx
libglib2.0-0
libglu1-mesa
libgpg-error0
libgstreamer1.0-0
libgtk2.0-0
libgtk-3-0
libice6
libmount1
libnss3
libpango1.0-0
libpango-1.0-0
libpangocairo-1.0-0
libpangoft2-1.0-0
libstdc++6
libtasn1-6
libunistring2
libwayland-egl1-mesa
libxml2
lsb-base
libxcb1
login
mount
mime-support
passwd
perl-base
systemd
systemd-timesyncd
systemd-sysv
tzdata
udev
util-linux
uuid-runtime
x11-common
zlib1g
`

func isExcluded(id string) bool {
	for _, exc := range strings.Split(EXCLUDE, "\n") {
		if id == exc {
			return true
		}
	}
	return false
}
