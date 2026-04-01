Unicode true

####
## Alfa Remote Client NSIS Installer
## 
## Silent mode: Run with /S for silent installation
## Install to user directory (no admin required): $LOCALAPPDATA
####

!include "MUI.nsh"
!include "FileFunc.nsh"

!define INFO_PROJECTNAME    "ARC"
!define INFO_COMPANYNAME    "DlnKot"
!define INFO_PRODUCTNAME    "Alfa Remote Client"
!define INFO_PRODUCTVERSION "0.6.3"
!define INFO_COPYRIGHT      "Copyright 2025 DlnKot"
!define PRODUCT_EXECUTABLE  "ARC.exe"
!define UNINST_KEY_NAME     "DlnKotAlfaRemoteClient"
!define REQUEST_EXECUTION_LEVEL "user"

VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

ManifestDPIAware true
ManifestSupportedOS all

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
OutFile "..\..\bin\${INFO_PROJECTNAME}-${ARCH}-installer.exe"

InstallDir "$LOCALAPPDATA\Alfa Remote Client"

ShowInstDetails show
ShowUnInstDetails show

Function .onInit
    StrCpy $0 "" 
    Pop $0 
    StrCpy $isSilent 0
    ${GetParameters} $0
    ${GetOptions} "/S" $0
    StrCmp $0 "" 0 +2
        StrCpy $isSilent 1
FunctionEnd

Section "!${INFO_PRODUCTNAME}" SecMain
    SetOutPath $INSTDIR
    
    File /oname=ARC.exe "..\bin\ARC.exe"
    
    WriteUninstaller "$INSTDIR\Uninstall.exe"
    
    CreateDirectory "$SMPROGRAMS\Alfa Remote Client"
    CreateShortcut "$SMPROGRAMS\Alfa Remote Client\Alfa Remote Client.lnk" "$INSTDIR\ARC.exe"
    CreateShortcut "$SMPROGRAMS\Alfa Remote Client\Удалить.lnk" "$INSTDIR\Uninstall.exe"
    
    IfFileExists "$DESKTOP\Alfa Remote Client.lnk" skip_desktop
        CreateShortcut "$DESKTOP\Alfa Remote Client.lnk" "$INSTDIR\ARC.exe"
    skip_desktop:
    
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "DisplayName" "${INFO_PRODUCTNAME}"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "UninstallString" "$\"$INSTDIR\Uninstall.exe$\""
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "InstallLocation" "$INSTDIR"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "DisplayIcon" "$INSTDIR\ARC.exe"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "Publisher" "${INFO_COMPANYNAME}"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "DisplayVersion" "${INFO_PRODUCTVERSION}"
    WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "NoModify" 1
    WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "NoRepair" 1
    
    ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
    IntFmt $0 "0x%08X" $0
    WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}" "EstimatedSize" "$0"
    
    StrCmp $isSilent 1 0 +4
        Sleep 1500
        Exec '"$INSTDIR\ARC.exe"'
SectionEnd

Section "uninstall"
    RMDir /r "$INSTDIR"
    
    Delete "$SMPROGRAMS\Alfa Remote Client\Alfa Remote Client.lnk"
    Delete "$SMPROGRAMS\Alfa Remote Client\Удалить.lnk"
    RMDir "$SMPROGRAMS\Alfa Remote Client"
    
    Delete "$DESKTOP\Alfa Remote Client.lnk"
    
    DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${UNINST_KEY_NAME}"
SectionEnd
