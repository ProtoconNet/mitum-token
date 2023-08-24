### mitum-token

*mitum-token* is a [mitum](https://github.com/ProtoconNet/mitum2)-based contract model and is a service that provides mint functions.

#### Installation

```sh
$ git clone https://github.com/ProtoconNet/mitum-token

$ cd mitum-token

$ go build -o ./mt ./main.go
```

#### Run

```sh
$ ./mt init --design=<config file> <genesis config file>

$ ./mt run --design=<config file>
```

[standalong.yml](standalone.yml) is a sample of `config file`.
[genesis-design.yml](genesis-design.yml) is a sample of `genesis config file`.