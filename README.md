# tf

Less verbose and more shell friendly Terraform.

## Install

```sh
curl https://raw.githubusercontent.com/dex4er/tf/main/tf.sh | sudo tee /usr/local/bin/tf
chmod +x /usr/local/bin/tf
```

## Usage

```sh
tf init
tf plan
tf apply
tf list
tf show
```

etc...

You can combine commands, ie.:

```sh
tf apply $(tf list | grep aws_vpc)
tf list | grep data.aws_region.current | xargs tf refresh
tf list | grep random_password | xargs tf rm
tf list | grep aws_subnet | xargs tf show
```

Be careful: you cannot use `xargs` with `tf apply`, `tf destroy` or `tf refresh`
because these commands are interactive. Still you can use `$()`.

### `tf apply`

The same as `terraform apply` with less verbose output.

Instead of Reading/Creating/Destroying... messages it will show a short progress
indicator.

It will skip `(known after apply)` lines from the output.

An additional option is `-compact` which will skip the content of the resources
completely.

The command accepts resource name as an argument without `-target=` option.

The command will generate temporarily the `terraform.tfplan` file.

The command will log to the file named in `TF_LOG_FILE` environment variable.
The file name is solved by `date` command for `%` sequences.

### `tf destroy`

The same as `terraform destroy` with less verbose output.

Instead of Reading/Creating/Destroying... messages it will show a short progress
indicator.

It will skip `(known after apply)` lines from the output.

An additional option is `-compact` which will skip the content of the resources
completely.

The command accepts resource name as an argument without `-target=` option.

The command will generate temporarily the `terraform.tfplan` file.

The command will log to the file named in `TF_LOG_FILE` environment variable.
The file name is solved by `date` command for `%` sequences.

### `tf init`

The same as `terraform init` with less verbose output.

The command will log to the file named in `TF_LOG_FILE` environment variable.
The file name is solved by `date` command for `%` sequences.

### `tf list`

The same as `terraform state list` with less verbose output and ANSI stripped.

### `tf mv`

The same as `terraform state mv` with less verbose output.

### `tf plan`

The same as `terraform plan` with less verbose output.

Instead of Reading... messages it will show a short progress indicator.

It will skip `(known after apply)` lines from the output.

An additional option is `-compact` which will skip the content of the resources
completely.

The command will log to the file named in `TF_LOG_FILE` environment variable.
The file name is solved by `date` command for `%` sequences.

### `tf refresh`

The same as `terraform apply -refresh-only` with less verbose output.

The command accepts resource name as an argument without `-target=` option.

The command will log to the file named in `TF_LOG_FILE` environment variable.
The file name is solved by `date` command for `%` sequences.

### `tf rm`

The same as `terraform state rm` with less verbose output.

The command accepts multiple arguments.

### `tf show`

The same as `terraform show` and `terraform state show` with less verbose output
and ANSI stripped.

`terraform show` is used when the command is run without arguments and
`terraform state show` when arguments are used.

The command accepts multiple arguments.

### `tf taint`

The same as `terraform taint` and it accepts multiple arguments.

### `tf untaint`

The same as `terraform untaint` and it accepts multiple arguments.

### License

Copyright (c) 2020-2022 Piotr Roszatycki <piotr.roszatycki@gmail.com>

[MIT](https://opensource.org/licenses/MIT)
