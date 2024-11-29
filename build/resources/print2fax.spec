%define _binaries_in_noarch_packages_terminate_build   0

Name: ##APP_NAME##
Version: ##VERSION##
Release: ##RELEASE##
Summary: online fax sender

BuildArch: ##ARCH##

License: GPL
URL: https://home.local

%description
create and load online fax sender

%install
mkdir -p %{buildroot}/usr/bin
cp -p %{_builddir}/##POSTINST_FILE##  %{buildroot}/usr/bin/##POSTINST_FILE##
cp -p %{_builddir}/##EXEC_NAME##     %{buildroot}/usr/bin/##EXEC_NAME##

%files
/usr/bin/##EXEC_NAME##
/usr/bin/##POSTINST_FILE##


%post
chmod a+x /usr/bin/##EXEC_NAME##
mkdir -p /etc/##APP_NAME##/bin/logs
chmod 777 -R /etc/##APP_NAME##/
chmod a+x /usr/bin/##POSTINST_FILE##
source /usr/bin/##POSTINST_FILE##

%postun
rm -rf /etc/##APP_NAME##


%changelog
* Fri Nov 15 2023 John Doe <john.doe@example.com> - 1.0-1
%autochangelog

