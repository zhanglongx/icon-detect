# icon-detect

Icon-detect detects the order of the icon-overlay of Windows, and performs auto-adjustment.

## Background

Windows supports icon-overlay for files and folders. But Windows also has a limited number for icon-overlay, so the order of icon-overlay is important. It's called 'icon-overlay battle'.

Find more in [List of shell icon overlay identifiers](https://en.wikipedia.org/wiki/List_of_shell_icon_overlay_identifiers) from Wikipedia. It also lists the icon of popular software.

## Usage

1. (FIXME) Modify source code.: 

    - Change `BOOST` list in `pkg/detect/detect.go` to the software you want to adjust. 

    - Change `PROGRAMTOKILL` in `main.go` to the application you want to restart, after the icon-overlay adjustment.

2. Build the binary.

    ```powershell
        PS > go build .
    ```

    OR (to hide cmd windows):

    ```powershell
        PS > go build -ldflags -H=windowsgui .
    ```

    ⚠️`-H=windowsgui` will also silence the stdout and stderr. e.g, `-h` and `-v` will not have any output.

3. Run it

    - ⚠️icon-detect need to be run as administrator, as it need to access the registry. See [UAC](#UAC) for more information.

    - icon-detect will output log file (with rotation) under the execute directory.

    - `-b` to backup the registry before auto-adjustment.

4. (Optional) Restarting Application via [Windows URI Scheme](https://learn.microsoft.com/en-us/windows/uwp/app-resources/uri-schemes).

    1. icon-detect needs to be registered as a protocol handler.

        ```powershell
            PS > .\icon-detect.exe -r
        ```

        ⚠️run as administrator.

    2. once registered, icon-detect can pop notification to restart application.

    3. unregister the protocol handler, if needed.

        ```powershell
            PS > .\icon-detect.exe -u
        ```

        ⚠️run as administrator.

## UAC

Windows has a UAC(User Account Control) mechanism. You can build icon-detect try to detect if UAC is enabled, and if it is, it will prompt a UAC dialog to ask for administrator permission.

```powershell
    PS > \path\to\mt.exe -manifest app.manifest -outputresource:icon-detect.exe;1
```

️`mt.exe` is included in [Windows SDK](https://developer.microsoft.com/en-us/windows/downloads/windows-10-sdk/).
