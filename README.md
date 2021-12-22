# apphide

Apphide is a utility to hide applications from the application overview by creating a .desktop file in `~/.local/share/applications` with the `Hidden` set to `true`.

It can match applications by their name or their id (e.g. `org.gnome.gedit`).
Applications can also be unhidden using the `-uh` flag by removing the file created in `~/.local/share/applications`.

For now only applications with `.desktop` files stored in `/usr/share/applications` are supported.
Apphide will also not edit already existing `.desktop` files in `~/.local/share/applications` and set/unsets the `Hidden` attribute.
