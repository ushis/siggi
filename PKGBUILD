# Maintainer: ushi <ushi@honkgong.info>
pkgname='siggi-git'
pkgver='0.0.0'
pkgrel=1
pkgdesc='WebRTC Signaling Server'
arch=('x86_64' 'i686' 'armv6h')
url='https://github.com/ushis/siggi'
license=('GPL2')
conflicts=('siggi')
provides=('siggi')
makedepends=('go')
options=(!strip)
source=('siggi::git+https://github.com/ushis/siggi.git#branch=master')
sha256sums=('SKIP')

pkgver() {
  cd siggi
  echo "$(git rev-list --count master).$(git rev-parse --short master)"
}

build() {
  cd siggi
  make
}

package() {
  cd siggi
  install -Dm0755 siggi                 "${pkgdir}/usr/bin/siggi"
  install -Dm0644 systemd/siggi.conf    "${pkgdir}/usr/lib/tmpfiles.d/siggi.conf"
  install -Dm0644 systemd/siggi.service "${pkgdir}/usr/lib/systemd/system/siggi.service"
}

# vim:set ts=2 sw=2 et:
