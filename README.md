# SimpleChain_preview_19.0
The Darren Chain 19.0 Preview Version

__Windows System:__

You can clone this repo and run <code>main.exe</code> directly.

![Server_EXE](assets\Server_EXE.png)

__All System:__

Since we’re going to run this chain, you should <a href="https://golang.org/dl/">installing</a> and configuring Golang first.

And we’ll also want to grab the following packages:

```shell
go get github.com/davecgh/go-spew/spew
```

```shell
go get github.com/joho/godotenv
```

If you all done, type the following commands to run the chain:

```shell
git clone https://github.com/darrenchain/SimpleChain_preview_19.0.git
cd SimpleChain_preview_19.0
```

```shell
go run main.go
```

As expected, we see the same genesis block.

![Server Compiled](assets\Server_compiled.png)

Finally, open your Terminal (It called "Command Prompt" on Windows System), and type the following command to run the client service:

Windows: <code>ncat localhost 9000</code>

Linux: <code>nc localhost 9000</code>

Then, we will see the same picture.

![Client_ncat](assets\Client_ncat.png)

Side Note: You can open multiple command line interface to activate multiple clients. I suggest you use a different color interface to make sure these are different clients like this:

![Server Clients Broadcast](assets\server_client_broadcast.png)
