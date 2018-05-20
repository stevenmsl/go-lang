# go-lang
Install Go 
- Go to https://golang.org/doc/install
- Run the msi to install go
- Close the command prompt and re-open it. Check the path to see if it includes the path to the bin folder where   the GO is installed something like C:\Go\bin 
- Close and re-launch the VS Code after the GO installation so the updated PATH environment variable can be        picked up by the VS Code

How the “go install” command works in Windows 10:
Let say you work space (where “.gitignore” file resides) folder is C:\Repos\Go\go-lang
-   Add the following to your user variables:
    Variable: GOPATH, Value: C:\Repos\Go\go-lang
-	In the command prompt, change the directory to C:\Repos\Go\go-lang
-	Create a folder named src under the workspace folder. The is where your source codes reside.        
    C:\Repos\Go\go-lang\src
-	Create package folders under the src folder. Put you source code files under package folder
    For example: C:\Repos\Go\go-lang\src\hello\hello.go
-	Run the following command in the command prompt. hello is a package folder under the src folder
    go install hello
-	You should see an exe file created in c:\Repos\Go\go-lang\bin folder. In this case it will be hello.exe 

The “go build” command works the same way. The only difference is that the generated executable is resided in the folder from where the command was run.