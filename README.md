# GoEnv
The virtual enviroments manager for Go!.

## Installation

### Get Binaries

See to [RELEASES PAGE](./releases) and get your binary file for various operational systems and architecturies.

NOTE: If gets alone binary, extract it (if archived) and PUT the directory into your PATH enviroment variable.

Check your $PATH enviroment variable:

```bash
export PATH=/dir/of/my/goenv_binary:$PATH
```

### Build

Get sources:
```bash
go get -u github.com/moisespsena-go/goenv/goenv
```

Check your $PATH enviroment variable:

```bash
export PATH=$GOPATH/bin:$PATH
```

### Add shortcut commands and Auto Completion

```bash
goenv setup | tee -a ~/.bashrc
source ~/.bashrc
```

## Usage

### The help command

```bash
goenv -h
```

Output:

```
Manage GoLang virtual enviroments.

Usage:
  goenv [command]

Available Commands:
  activate    Activate the virtualenv with NAME.
  backup      Create backup .tar.gz file for the virtualenv with have NAME.
  completion  Generates bash completion scripts
  db          Returns the current database path.
  help        Help about any command
  init        Init new virtual enviroment.
  ls          List all virtual enviroments on current database.
  path        Print env path
  restore     Restore backup.tar.gz file to the virtualenv with have NAME.
  rm          Remove the virtualenv with have NAME.
  setup       Generate sources for custom prompt commands.
  update      Update activation scripts.
  versions    Manage golang binary versions

Flags:
  -d, --db string   Database directory (default is $HOME/.goenv). (default "~/.goenv")
  -h, --help        help for goenv

Use "goenv [command] --help" for more information about a command.
```

### Init new enviroment

```bash
goenv init env1
```

### List enviroments

```bash
goenv ls
```

### Activate repository:

```bash
eval $(goenv activate env1)
```

or (see to [Database](#database) section)

```bash
source $(goenv db)/env1/activate
```

#### Deactivate it

```bash
goenv-deactivate
```

### Remove repository:

Move to trash directory (`DB_DIR/.trash`, see to [Database](#database) section):
```bash
goenv rm env1
```

Remove permanently:
```bash
goenv rm -p env1
```

### Database

Get database path:

```bash
goenv db
```

#### Custom database path

##### With command paramenter
For use custom database path, set the `-d` flag:

```bash
goenv -d ~/custom-db args...
```

Examples:

```bash
goenv -d ~/custom-db init env2
goenv -d ~/custom-db ls
```

##### With enviroment variable

Set the enviroment variable:
 
```bash
export GOENVDB=~/custom-db
```

or

```bash
GOENVDB=~/custom-db goenv args...
```

## Thank's!

By [Moises P. Sena](https://github.com/moisespsena).