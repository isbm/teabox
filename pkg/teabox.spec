#
# spec file for Teabox package
#
# Copyright (c) 2023 Bo Maryniuk
#
# All modifications and additions to the file contributed by third parties
# remain the property of their copyright owners, unless otherwise agreed
# upon. The license for this file, and modifications and additions to the
# file, is the same license as for the pristine package itself (unless the
# license for the pristine package is not an Open Source License, in which
# case the license is the MIT License). An "Open Source License" is a
# license that conforms to the Open Source Definition (Version 1.9)
# published by the Open Source Initiative.

Name:           teabox
Version:        0.3
Release:        0
Summary:        It is like a whiptail on steroids to make configuration consoles
License:        MIT
Group:          Tools
Url:            https://teabox.readthedocs.io/en/latest/
Source:         %{name}-%{version}.tar.gz
BuildRequires:  make
BuildRequires:  golang
BuildRequires:  debbuild

%description
If you need to resemble something like YaST2 or make a yet another configuration tool
or create an interface to your scripts, then Teabox is for you.

%prep
%setup -q

%build
go build -a -mod=vendor -buildmode=pie -ldflags="-s -w" -o %{name} ./cmd/*go

%install
install -D -m 0755 %{name} %{buildroot}%{_bindir}/%{name}

%files
%defattr(-,root,root)
%{_bindir}/%{name}

%changelog
