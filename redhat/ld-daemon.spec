%global commit aec8718aea6cec2103705c3b629adca296d16477
%global shortcommit %(c=%{commit}; echo ${c:0:7})

%define debug_package %{nil}

Name: ld-daemon
Version: 2.0.0
Release: 1.%{shortcommit}%{?dist}
Summary: LaunchDarkly Redis Daemon

Group: Development/Tools
License: Apache License, Version 2.0
URL: https://github.com/launchdarkly/ld-daemon/
Source: https://github.com/launchdarkly/%{name}/archive/%{commit}.tar.gz#/%{name}-%{shortcommit}.tar.gz

BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires: git, golang

%description
The LaunchDarkly Redis daemon establishes a connection to the LaunchDarkly streaming API, then pushes feature updates to a Redis store.

The daemon can be used to offload the task of maintaining a stream and writing to Redis from our SDKs. This can give platforms that do not support SSE (e.g. PHP) the benefits of LaunchDarkly's streaming model.

The daemon can be configured to synchronize multiple environments, even across multiple projects.

%prep
%setup -q -n %{name}-%{commit}

%build
export GOPATH=$(pwd)
go get github.com/tools/godep
bin/godep go build

%install
rm -rf %{buildroot}
install -p -D -m 0755 ld-daemon-%{commit} %{buildroot}%{_bindir}/ld-daemon
install -p -D -m 0644 deb-contents/etc/ld-daemon.conf %{buildroot}%{_sysconfdir}/ld-daemon.conf

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%{_bindir}/ld-daemon
%config(noreplace) %{_sysconfdir}/ld-daemon.conf

%changelog
* Tue May 31 2016 kevin.pankonen - 2.0.0-1.aec8718
- build ld-daemon 2.0.0
