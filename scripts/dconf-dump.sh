dconf dump /org/gnome/Ptyxis/ \
    > settings/ptyxis.ini

dconf dump /org/gnome/desktop/background/ \
    > settings/background.ini

dconf dump /org/gnome/desktop/interface/ \
    > settings/interface.ini

dconf dump /org/gnome/desktop/input-sources/ \
    > settings/input.ini

dconf dump /org/gnome/desktop/app-folders/ \
    > settings/app-folders.ini

dconf dump /org/gnome/mutter/ \
    > settings/mutter.ini

dconf dump /org/gnome/nautilus/preferences/ \
    > settings/nautilus.ini

dconf dump /org/gnome/settings-daemon/plugins/power/ \
    > settings/power.ini

dconf dump /org/gnome/shell/ \
    > settings/shell.ini
