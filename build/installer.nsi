!define INFO_PROJECTNAME "ARC"
!define INFO_COMPANYNAME "DlnKot"
!define INFO_PRODUCTNAME "Alfa Remote Client"
!define INFO_COPYRIGHT "Copyright 2025 DlnKot"
!define PRODUCT_EXECUTABLE "ARC.exe"
!define UNINST_KEY_NAME "DlnKotAlfaRemoteClient"

Unicode true
!include "MUI.nsh"

VIProductVersion "__VERSION__.0"
VIFileVersion "__VERSION__.0"
VIAddVersionKey "CompanyName" "DlnKot"
VIAddVersionKey "FileDescription" "Alfa Remote Client Installer"
VIAddVersionKey "ProductVersion" "__VERSION__"
VIAddVersionKey "FileVersion" "__VERSION__"
VIAddVersionKey "LegalCopyright" "Copyright 2025 DlnKot"
VIAddVersionKey "ProductName" "Alfa Remote Client"

ManifestDPIAware true
!define MUI_ICON "windows\icon.ico"
!define MUI_UNICON "windows\icon.ico"
!define MUI_ABORTWARNING
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "Russian"

Name "Alfa Remote Client"
OutFile "bin\ARC-__VERSION__-installer.exe"
InstallDir "$LOCALAPPDATA\Alfa Remote Client"
ShowInstDetails show

Section "!Alfa Remote Client" SecMain
SetOutPath $INSTDIR
File "bin\ARC.exe"
WriteUninstaller "$INSTDIR\Uninstall.exe"
CreateDirectory "$SMPROGRAMS\Alfa Remote Client"
CreateShortcut "$SMPROGRAMS\Alfa Remote Client\Alfa Remote Client.lnk" "$INSTDIR\ARC.exe"
CreateShortcut "$SMPROGRAMS\Alfa Remote Client\Удалить.lnk" "$INSTDIR\Uninstall.exe"
IfFileExists "$DESKTOP\Alfa Remote Client.lnk" skip_desktop
CreateShortcut "$DESKTOP\Alfa Remote Client.lnk" "$INSTDIR\ARC.exe"
skip_desktop:
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DlnKotAlfaRemoteClient" "DisplayName" "Alfa Remote Client"
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DlnKotAlfaRemoteClient" "UninstallString" "$INSTDIR\Uninstall.exe"
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DlnKotAlfaRemoteClient" "InstallLocation" "$INSTDIR"
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DlnKotAlfaRemoteClient" "DisplayIcon" "$INSTDIR\ARC.exe"
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DlnKotAlfaRemoteClient" "Publisher" "DlnKot"
WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DlnKotAlfaRemoteClient" "DisplayVersion" "__VERSION__"
WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DlnKotAlfaRemoteClient" "NoModify" 1
WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DlnKotAlfaRemoteClient" "NoRepair" 1
SectionEnd

Section "uninstall"
RMDir /r "$INSTDIR"
Delete "$SMPROGRAMS\Alfa Remote Client\Alfa Remote Client.lnk"
Delete "$SMPROGRAMS\Alfa Remote Client\Удалить.lnk"
RMDir "$SMPROGRAMS\Alfa Remote Client"
Delete "$DESKTOP\Alfa Remote Client.lnk"
DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\DlnKotAlfaRemoteClient"
SectionEnd
