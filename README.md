# general

This utility let you search your list of shortcuts from `keybindings.xml`
in a more comfortable way.

# usage

In windows explorer, go to : `C:\Users\your_name\AppData\Renoise\v3.x.x`

Copy the path to `Keybindings.xml` (use "Copy as path", and remove the `"`)

You can drag/drop the file on the window too.

# build

Linux and others should run fine with just `go build`

With windows, you can use this command to avoid having it run along a console

```
go build -ldflags -H=windowsgui
```
