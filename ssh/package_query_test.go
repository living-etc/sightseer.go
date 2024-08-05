package ssh_test

import (
	"log"
	"testing"

	sightseer "github.com/living-etc/sightseer.go/ssh"
)

func Test_packageFromDpkgOutput(t *testing.T) {
	tests := []struct {
		testname    string
		dpkgoutput  string
		packageWant *sightseer.Package
	}{
		{
			testname: "openssh-server is installed",
			dpkgoutput: `Package: openssh-server
Status: install ok installed
Priority: optional
Section: net
Installed-Size: 2140
Maintainer: Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>
Architecture: arm64
Multi-Arch: foreign
Source: openssh
Version: 1:9.6p1-3ubuntu13
Replaces: openssh-client (<< 1:7.9p1-8), ssh, ssh-krb5
Provides: ssh-server
Depends: adduser, libpam-modules, libpam-runtime, lsb-base, openssh-client (= 1:9.6p1-3ubuntu13), openssh-sftp-server, procps, ucf, debconf (>= 0.5) | debconf-2.0, libaudit1 (>= 1:2.2.1), libc6 (>= 2.38), libcom-err2 (>= 1.43.9), libcrypt1 (>= 1:4.1.0), libgssapi-krb5-2 (>= 1.17), libkrb5-3 (>= 1.13~alpha1+dfsg), libpam0g (>= 0.99.7.1), libselinux1 (>= 3.1~), libssl3t64 (>= 3.0.13), libwrap0 (>= 7.6-4~), zlib1g (>= 1:1.1.4)
Pre-Depends: init-system-helpers (>= 1.54~)
Recommends: default-logind | logind | libpam-systemd, ncurses-term, xauth, ssh-import-id
Suggests: molly-guard, monkeysphere, ssh-askpass, ufw
Conflicts: sftp, ssh-socks, ssh2
Conffiles:
 /etc/default/ssh 500e3cf069fe9a7b9936108eb9d9c035
 /etc/init.d/ssh 3649a6fe8c18ad1d5245fd91737de507
 /etc/pam.d/sshd 8b4c7a12b031424b2a9946881da59812
 /etc/ssh/moduli 366395e79244c54223455e5f83dafba3
 /etc/ufw/applications.d/openssh-server 486b78d54b93cc9fdc950c1d52ff479e
Description: secure shell (SSH) server, for secure access from remote machines
 This is the portable version of OpenSSH, a free implementation of
 the Secure Shell protocol as specified by the IETF secsh working
 group.
 .
 Ssh (Secure Shell) is a program for logging into a remote machine
 and for executing commands on a remote machine.
 It provides secure encrypted communications between two untrusted
 hosts over an insecure network. X11 connections and arbitrary TCP/IP
 ports can also be forwarded over the secure channel.
 It can be used to provide applications with a secure communication
 channel.
 .
 This package provides the sshd server.
 .
 In some countries it may be illegal to use any encryption at all
 without a special permit.
 .
 sshd replaces the insecure rshd program, which is obsolete for most
 purposes.
Homepage: https://www.openssh.com/
Original-Maintainer: Debian OpenSSH Maintainers <debian-ssh@lists.debian.org>
`,
			packageWant: &sightseer.Package{
				Name:          "openssh-server",
				Status:        "install ok installed",
				Priority:      "optional",
				Section:       "net",
				InstalledSize: "2140",
				Maintainer:    "Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>",
				Architecture:  "arm64",
				MultiArch:     "foreign",
				Source:        "openssh",
				Version:       "1:9.6p1-3ubuntu13",
				Replaces:      "openssh-client (<< 1:7.9p1-8), ssh, ssh-krb5",
				Provides:      "ssh-server",
				Depends:       "adduser, libpam-modules, libpam-runtime, lsb-base, openssh-client (= 1:9.6p1-3ubuntu13), openssh-sftp-server, procps, ucf, debconf (>= 0.5) | debconf-2.0, libaudit1 (>= 1:2.2.1), libc6 (>= 2.38), libcom-err2 (>= 1.43.9), libcrypt1 (>= 1:4.1.0), libgssapi-krb5-2 (>= 1.17), libkrb5-3 (>= 1.13~alpha1+dfsg), libpam0g (>= 0.99.7.1), libselinux1 (>= 3.1~), libssl3t64 (>= 3.0.13), libwrap0 (>= 7.6-4~), zlib1g (>= 1:1.1.4)",
				PreDepends:    "init-system-helpers (>= 1.54~)",
				Recommends:    "default-logind | logind | libpam-systemd, ncurses-term, xauth, ssh-import-id",
				Suggests:      "molly-guard, monkeysphere, ssh-askpass, ufw",
				Conflicts:     "sftp, ssh-socks, ssh2",
				Conffiles: `
/etc/default/ssh 500e3cf069fe9a7b9936108eb9d9c035
/etc/init.d/ssh 3649a6fe8c18ad1d5245fd91737de507
/etc/pam.d/sshd 8b4c7a12b031424b2a9946881da59812
/etc/ssh/moduli 366395e79244c54223455e5f83dafba3
/etc/ufw/applications.d/openssh-server 486b78d54b93cc9fdc950c1d52ff479e`,
				Description: `secure shell (SSH) server, for secure access from remote machines
This is the portable version of OpenSSH, a free implementation of
the Secure Shell protocol as specified by the IETF secsh working
group.
.
Ssh (Secure Shell) is a program for logging into a remote machine
and for executing commands on a remote machine.
It provides secure encrypted communications between two untrusted
hosts over an insecure network. X11 connections and arbitrary TCP/IP
ports can also be forwarded over the secure channel.
It can be used to provide applications with a secure communication
channel.
.
This package provides the sshd server.
.
In some countries it may be illegal to use any encryption at all
without a special permit.
.
sshd replaces the insecure rshd program, which is obsolete for most
purposes.`,
				Homepage:           "https://www.openssh.com/",
				OriginalMaintainer: "Debian OpenSSH Maintainers <debian-ssh@lists.debian.org>",
			},
		},
	}

	for _, testcase := range tests {
		var query sightseer.PackageQuery
		packageGot, err := query.ParseOutput(testcase.dpkgoutput)
		if err != nil {
			log.Fatalf("Error in %v: %v", testcase.testname, err)
		}

		EvaluateStructTypesAreEqual(packageGot, testcase.packageWant, testcase.testname, t)
	}
}
