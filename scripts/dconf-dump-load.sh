#!/usr/bin/bash

if [ -z "$1" ]; then
    echo "Usage: $0 [dump|load]"
    exit 1
fi

ACTION="$1"

declare -A CONFIGS=(
    ["/org/gnome/Ptyxis/"]="data/settings/ptyxis.ini"
    ["/org/gnome/desktop/background/"]="data/settings/background.ini"
    ["/org/gnome/desktop/interface/"]="data/settings/interface.ini"
    ["/org/gnome/desktop/input-sources/"]="data/settings/input.ini"
    ["/org/gnome/desktop/app-folders/"]="data/settings/app-folders.ini"
    ["/org/gnome/mutter/"]="data/settings/mutter.ini"
    ["/org/gnome/nautilus/preferences/"]="data/settings/nautilus.ini"
    ["/org/gnome/settings-daemon/plugins/power/"]="data/settings/power.ini"
    ["/org/gnome/shell/"]="data/settings/shell.ini"
)

case "$ACTION" in
    dump)
        echo "Dumping GNOME settings..."
        mkdir -p data/settings
        for path in "${!CONFIGS[@]}"; do
            file="${CONFIGS[$path]}"
            dconf dump "$path" > "$file"
            echo "  Saved: $path -> $file"
        done
        echo "Done!"
        ;;

    load)
        echo "Loading GNOME settings..."
        for path in "${!CONFIGS[@]}"; do
            file="${CONFIGS[$path]}"
            if [ -f "$file" ]; then
                dconf load "$path" < "$file"
                echo "  Loaded: $file -> $path"
            else
                echo "  Skipped: $file (file not found)"
            fi
        done
        echo "Done!"
        ;;

    *)
        echo "Error: Invalid argument '$ACTION'. Use 'dump' or 'load'."
        exit 1
        ;;
esac
