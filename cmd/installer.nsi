!define APPNAME "Fehmi Cloud Connector"
!define VERSION "1.2.0"
!define COMPANY "Fehmi Corporation"

OutFile "${OUTPUT_FILE}"
InstallDir "$PROGRAMFILES64\FehmiCloudConnector"
RequestExecutionLevel admin

Page directory
Page instfiles

Section "Install"
    SetOutPath "$INSTDIR"
    
    File "${SOURCE_DIR}\connector.exe"
    File "${SOURCE_DIR}\app.exe"
    File "icon.ico"
    
    # 2. Create Desktop Shortcut with custom icon 
    # Syntax: CreateShortcut "link_file" "target_path" "parameters" "icon_file" icon_index_number
    CreateShortcut "$DESKTOP\Fehmi Cloud Connector.lnk" "$INSTDIR\connector.exe" "" "$INSTDIR\icon.ico" 0

    # 3. Create Start Menu entry with custom icon
    CreateDirectory "$SMPROGRAMS\Fehmi Corporation"
    CreateShortcut "$SMPROGRAMS\Fehmi Corporation\Fehmi Cloud Connector.lnk" "$INSTDIR\connector.exe" "" "$INSTDIR\icon.ico" 0
    
    # Write Uninstaller registry keys
    WriteUninstaller "$INSTDIR\uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\FehmiCloudConnector" "DisplayName" "${APPNAME}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\FehmiCloudConnector" "DisplayVersion" "${VERSION}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\FehmiCloudConnector" "Publisher" "${COMPANY}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\FehmiCloudConnector" "UninstallString" "$INSTDIR\uninstall.exe"
SectionEnd

Section "Uninstall"
    Delete "$INSTDIR\connector.exe"
    Delete "$INSTDIR\app.exe"
    Delete "$INSTDIR\icon.ico"
    Delete "$INSTDIR\uninstall.exe"
    
    Delete "$DESKTOP\Fehmi Cloud Connector.lnk"
    Delete "$SMPROGRAMS\Fehmi Corporation\Fehmi Cloud Connector.lnk"
    RMDir "$SMPROGRAMS\Fehmi Corporation"
    
    RMDir "$INSTDIR"
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\FehmiCloudConnector"
SectionEnd