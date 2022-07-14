# apphide

Apphide is a utility to hide applications from the application overview of your Linux desktop environment by creating a `.desktop` file in `~/.local/share/applications` with the `Hidden` attribute set to `true`.

It can match applications through regular expressions by the name they appear in the appplication menu in or by their id (e.g. `org.gnome.gedit`).
Applications can also be unhidden using the `-uh` flag by removing the file created in `~/.local/share/applications`.

Only applications with `.desktop` files stored in `/usr/share/applications` are supported, as apphide isn't able to edit already existing `.desktop` files in `~/.local/share/applications`.
