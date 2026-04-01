!define INFO_PROJECTNAME "ARC"
!define INFO_COMPANYNAME "DlnKot"
!define INFO_PRODUCTNAME "Alfa Remote Client"
!define INFO_PRODUCTVERSION "VERSION"
!define INFO_COPYRIGHT "Copyright 2025 DlnKot"
!define PRODUCT_EXECUTABLE "ARC.exe"
!define UNINST_KEY_NAME "DlnKotAlfaRemoteClient"

Unicode true
!include "MUI.nsh"

VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion "${INFO_PRODUCTVERSION}.0"
VIAddVersionKey "CompanyName" "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion" "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion" "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright" "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName" "${INFO_PRODUCTNAME}"

ManifestDPIAware true
!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
!define MUI_ABORTWARNING
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "Russian"

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\ARC-VERSION-installer.exe"
InstallDir "$LOCALAPPDATA\Alfa Remote Client"
ShowInstDetails show

Section "!${INFO_PRODUCTNAME}" SecMain
SetOutPath $INSTDIR
File "..\..\bin\ARC.exe"
WriteUninstaller "$INSTDIR\Uninstall.exe"
CreateDirectory "$SMPROGRAMS\Alfa Remote Client"
CreateShortcut "$SMPROGRAMS\Alfa Remote Client\Alfa Remote Client.lnk" "$INSTDIR\ARC.exe"
CreateShortcut "$SMPROGRAMS\Alfa Remote Client\Удалить.lnk" "$INSTDIR\Uninstall.exe"
IfFileExists "$DESKTOP\Alfa Remote Client.lnk" skip_desktop
CreateShortcut "$DESKTOP\Alfa Remote Client.lnk" "$INSTDIR\ARC.exe"
skip_desktop:
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "DisplayName" "${INFO_PRODUCTNAME}"
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "UninstallString" "msiexec /x$INSTDIR\Uninstall.exe"
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "InstallLocation" "$INSTDIR"
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "DisplayIcon" "$INSTDIR\ARC.exe"
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "Publisher" "${INFO_COMPANYNAME}"
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "DisplayVersion" "${INFO_PRODUCTVERSION}"
WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "NoModify" 1
WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "NoRepair" 1
SectionEnd

Section "uninstall"
RMDir /r "$INSTDIR"
Delete "$SMPROGRAMS\Alfa Remote Client\Alfa Remote Client.lnk"
Delete "$SMPROGRAMS\Alfa Remote Client\Удалить.lnk"
RMDir "$SMPROGRAMS\Alfa Remote Client"
Delete "$DESKTOP\Alfa Remote Client.lnk"
DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}"
SectionEnd
