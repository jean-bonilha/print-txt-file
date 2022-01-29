# Print TXT file with Go Language

# Installation

Clone project

```bash
$ git clone https://github.com/jean-bonilha/print-txt-file.git && cd print-txt-file
```

Install dependencies

```bash
$ go mod tidy
```

Real-time test command

```bash
$ go run main.go
```

Then you can make a test by putting a file inside "C:/PrinterLabel" file

# Setup on the client printer service

Building a production .exe file

```bash
$ go build -ldflags -H=windowsgui -o PrintTXTFile.exe main.go
```

**Open Windows 7/10 Startup Folder with Explorer**

To access the “*All Users”* Startup folder in Windows 7/10, open the Run dialog box (**Windows Key + R**), type **`shell:common startup`**, and click **OK**.

![https://i1.wp.com/www.alphr.com/wp-content/uploads/2020/10/windows-10-run-shell-common-startup.jpg?w=690&ssl=1](https://i1.wp.com/www.alphr.com/wp-content/uploads/2020/10/windows-10-run-shell-common-startup.jpg?w=690&ssl=1)

Put your PrintTXTFile.exe file created before into the Startup folder in Windows 7/10.

![https://jeanbonilhawebdev100418.files.wordpress.com/2021/12/printlabeldoc.png](https://jeanbonilhawebdev100418.files.wordpress.com/2021/12/printlabeldoc.png)

To finish setup you need to run the PrintTXTFile.exe with a double click **or** restart the Windows.
