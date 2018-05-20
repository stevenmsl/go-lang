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
    (It seems you should only have one workspace for the Go on the machine for a specific user. By doing so you will have only one entry in GOPATH variable, which is important as later you will start using go get command to download additional packages and want to control where to store those packages. It seems the go get command will take the first entry in the GOPATH as the default workspace folder and put the downloaded packages in its src subfolder.)

-	In the command prompt, change the directory to C:\Repos\Go\go-lang
-	Create a folder named src under the workspace folder. The is where your source codes reside.        
    C:\Repos\Go\go-lang\src
-	Create package folders under the src folder. Put you source code files under package folder
    For example: C:\Repos\Go\go-lang\src\hello\hello.go
-	Run the following command in the command prompt. hello is a package folder under the src folder
    go install hello
-	You should see an exe file created in c:\Repos\Go\go-lang\bin folder. In this case it will be hello.exe 

The “go build” command works the same way. The only difference is that the generated executable is resided in the folder from where the command was run.

Install go-ethereum

- Go to http://tdm-gcc.tdragon.net/ to download the gcc compiler. After the installation, you should be able to    locate gcc command in the bin subfolder under your TDM-GCC installation folder something like: C:\TDM-GCC-64\bin

- run "go get -v github.com/ethereum/go-ethereum" to download go-ethereum. You should be able to see the downloaded packages in C:\repos\go-lang\src\github.com folder

- run "go install github.com/ethereum/go-ethereum/cmd/geth" to build geth exe file and put it in the bin folder under the workspace 