# icon-detect

Icon-detect detects the order of the icon-overlay of Windows, and performs auto-adjustment.

## Background

Windows supports icon-overlay for files and folders. But Windows also has a  limited number for icon-overlay, so the order of icon-overlay is important. It's called 'icon-overlay battle'.

Find more in [List of shell icon overlay identifiers](https://en.wikipedia.org/wiki/List_of_shell_icon_overlay_identifiers) from Wikipedia. It also lists the icon of popular software.

## Usage

1. Build the binary.

    ```powershell
        PS > go build .
    ```

    OR (to hide cmd windows):

    ```powershell
        PS > go build -ldflags -H=windowsgui .
    ```

2. Run it

    - ⚠️icon-detect need to be run as administrator, as it need to access the registry.

    - icon-detect will output log file (with rotation) under the execute directory.

    - `-b` to backup the registry before auto-adjustment.
