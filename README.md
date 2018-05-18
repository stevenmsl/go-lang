# go-lang
How “go install” command work in Windows 10:
Let say you work space (where “.gitignore” file resides) is C:\Repos\Go\go-lang
1.	Add the following to your user variables:
Variable: GOPATH, Value: C:\Repos\Go\go-lang
2.	In the command prompt, change the directory to C:\Repos\Go\go-lang
3.	Make sure you have a src subfolder: C:\Repos\Go\go-lang\src
4.	Your source code files should live inside a subfolder under the src subfolder:
    For example: C:\Repos\Go\go-lang\src\hello\hello.go
5.	Run the following command in the command prompt. hello is the a subfolder
    go install hello

6.	You should see an exe file created in c:\Repos\Go\go-lang\bin folder. In this case it’s hello.exe 
“go build” command works the same way. The only difference is that the generated executable is resided in the folder from where the command was run.
