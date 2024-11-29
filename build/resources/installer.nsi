!include "MUI.nsh"
!include x64.nsh

Name "Print2Fax Installer"
OutFile "../../bin/Print2Fax_Installer.exe"

Var EXECUTABLE_FILE
Var EXECUTABLE_ARGS
Var PRINTER_ICON_NAME

; Default installation directory
InstallDir $PROGRAMFILES\Print2Fax

; Request admin privileges
RequestExecutionLevel admin

!define MUI_ABORTWARNING

; =============================================
; Pages
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

; =============================================
; Languages
!insertmacro MUI_LANGUAGE English

; =============================================

; Installation directory page
Function .onInit
    StrCpy $PRINTER_ICON_NAME "printer_win.ico"
FunctionEnd

; =============================================

Section
    SetOutPath $INSTDIR
    WriteUninstaller "$INSTDIR\Uninstall.exe"

    File "printer_win.ico"
    
    ; Copy the appropriate file based on the operating system
    ${If} ${RunningX64}
        File "../../bin/fax_sender_ui_64.exe"
        StrCpy $EXECUTABLE_FILE "fax_sender_ui_64.exe"
    ${Else}
        MessageBox MB_ICONSTOP|MB_OK "this software is not working on 32 bit version."
        Abort
    ${EndIf}

    ; Create a shortcut for the installed file with an icon
    CreateShortCut "$DESKTOP\Print2Fax.lnk" "$INSTDIR\$EXECUTABLE_FILE" "" "$INSTDIR\$PRINTER_ICON_NAME" 0 SW_SHOWNORMAL
SectionEnd

; =============================================

Section "Start Menu Shortcuts"
    ; Create Start Menu directory
    CreateDirectory "$SMPROGRAMS"

    ; Create a shortcut in Start Menu for the installed file
    CreateShortCut "$SMPROGRAMS\Print2Fax.lnk" "$INSTDIR\$EXECUTABLE_FILE" "" "$INSTDIR\$PRINTER_ICON_NAME" 0 SW_SHOWNORMAL
SectionEnd

; =============================================

Section "Context Menu Integration"
    
    StrCpy $EXECUTABLE_ARGS '"-show-sender" -working-dir="$INSTDIR" -file-path="%1"'

    WriteRegStr HKCR "SystemFileAssociations\.pdf\shell\Print2Fax" "" "Print2Fax"
    WriteRegStr HKCR "SystemFileAssociations\.pdf\shell\Print2Fax" "Icon" "$INSTDIR\$PRINTER_ICON_NAME"
    WriteRegStr HKCR "SystemFileAssociations\.pdf\shell\Print2Fax\command" "" '"$INSTDIR\$EXECUTABLE_FILE" $EXECUTABLE_ARGS'

    WriteRegStr HKCR "SystemFileAssociations\.txt\shell\Print2Fax" "" "Print2Fax"
    WriteRegStr HKCR "SystemFileAssociations\.txt\shell\Print2Fax" "Icon" "$INSTDIR\$PRINTER_ICON_NAME"
    WriteRegStr HKCR "SystemFileAssociations\.txt\shell\Print2Fax\command" "" '"$INSTDIR\$EXECUTABLE_FILE" $EXECUTABLE_ARGS'

    WriteRegStr HKCR "SystemFileAssociations\.tiff\shell\Print2Fax" "" "Print2Fax"
    WriteRegStr HKCR "SystemFileAssociations\.tiff\shell\Print2Fax" "Icon" "$INSTDIR\$PRINTER_ICON_NAME"
    WriteRegStr HKCR "SystemFileAssociations\.tiff\shell\Print2Fax\command" "" '"$INSTDIR\$EXECUTABLE_FILE" $EXECUTABLE_ARGS'

    WriteRegStr HKCR "SystemFileAssociations\.tif\shell\Print2Fax" "" "Print2Fax"
    WriteRegStr HKCR "SystemFileAssociations\.tif\shell\Print2Fax" "Icon" "$INSTDIR\$PRINTER_ICON_NAME"
    WriteRegStr HKCR "SystemFileAssociations\.tif\shell\Print2Fax\command" "" '"$INSTDIR\$EXECUTABLE_FILE" $EXECUTABLE_ARGS'

    WriteRegStr HKCR "SystemFileAssociations\.jpeg\shell\Print2Fax" "" "Print2Fax"
    WriteRegStr HKCR "SystemFileAssociations\.jpeg\shell\Print2Fax" "Icon" "$INSTDIR\$PRINTER_ICON_NAME"
    WriteRegStr HKCR "SystemFileAssociations\.jpeg\shell\Print2Fax\command" "" '"$INSTDIR\$EXECUTABLE_FILE" $EXECUTABLE_ARGS'

    WriteRegStr HKCR "SystemFileAssociations\.jpg\shell\Print2Fax" "" "Print2Fax"
    WriteRegStr HKCR "SystemFileAssociations\.ipg\shell\Print2Fax" "Icon" "$INSTDIR\$PRINTER_ICON_NAME"
    WriteRegStr HKCR "SystemFileAssociations\.jpg\shell\Print2Fax\command" "" '"$INSTDIR\$EXECUTABLE_FILE" $EXECUTABLE_ARGS'

    WriteRegStr HKCR "SystemFileAssociations\.png\shell\Print2Fax" "" "Print2Fax"
    WriteRegStr HKCR "SystemFileAssociations\.png\shell\Print2Fax" "Icon" "$INSTDIR\$PRINTER_ICON_NAME"
    WriteRegStr HKCR "SystemFileAssociations\.png\shell\Print2Fax\command" "" '"$INSTDIR\$EXECUTABLE_FILE" $EXECUTABLE_ARGS'

    WriteRegStr HKCR "SystemFileAssociations\.doc\shell\Print2Fax" "" "Print2Fax"
    WriteRegStr HKCR "SystemFileAssociations\.doc\shell\Print2Fax" "Icon" "$INSTDIR\$PRINTER_ICON_NAME"
    WriteRegStr HKCR "SystemFileAssociations\.doc\shell\Print2Fax\command" "" '"$INSTDIR\$EXECUTABLE_FILE" $EXECUTABLE_ARGS'

    WriteRegStr HKCR "SystemFileAssociations\.docx\shell\Print2Fax" "" "Print2Fax"
    WriteRegStr HKCR "SystemFileAssociations\.docx\shell\Print2Fax" "Icon" "$INSTDIR\$PRINTER_ICON_NAME"
    WriteRegStr HKCR "SystemFileAssociations\.docx\shell\Print2Fax\command" "" '"$INSTDIR\$EXECUTABLE_FILE" $EXECUTABLE_ARGS'

    WriteRegStr HKCR "SystemFileAssociations\.odt\shell\Print2Fax" "" "Print2Fax"
    WriteRegStr HKCR "SystemFileAssociations\.odt\shell\Print2Fax" "Icon" "$INSTDIR\$PRINTER_ICON_NAME"
    WriteRegStr HKCR "SystemFileAssociations\.odt\shell\Print2Fax\command" "" '"$INSTDIR\$EXECUTABLE_FILE" $EXECUTABLE_ARGS'
SectionEnd

; =============================================

Section "Uninstall"
    ${If} ${RunningX64}
        Delete "$INSTDIR\fax_sender_ui_64.exe"
    ${EndIf}
    Delete "$DESKTOP\Print2Fax.lnk"
    Delete "$INSTDIR\printer_win.ico"
    Delete "$INSTDIR\Uninstall.exe"

    RMDir /r "$INSTDIR\bin"
    RMDir $INSTDIR
    Delete "$SMPROGRAMS\Print2Fax.lnk"
    DeleteRegKey HKCR "SystemFileAssociations\.pdf\shell\Print2Fax"
    DeleteRegKey HKCR "SystemFileAssociations\.txt\shell\Print2Fax"
    DeleteRegKey HKCR "SystemFileAssociations\.tiff\shell\Print2Fax"
    DeleteRegKey HKCR "SystemFileAssociations\.tif\shell\Print2Fax"
    DeleteRegKey HKCR "SystemFileAssociations\.jpeg\shell\Print2Fax"
    DeleteRegKey HKCR "SystemFileAssociations\.jpg\shell\Print2Fax"
    DeleteRegKey HKCR "SystemFileAssociations\.png\shell\Print2Fax"
    DeleteRegKey HKCR "SystemFileAssociations\.doc\shell\Print2Fax"
    DeleteRegKey HKCR "SystemFileAssociations\.docx\shell\Print2Fax"
    DeleteRegKey HKCR "SystemFileAssociations\.odt\shell\Print2Fax"
  SectionEnd
