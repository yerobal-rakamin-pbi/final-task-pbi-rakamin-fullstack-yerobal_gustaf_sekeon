# Rakamin Final Task Go

## Pre-requisite
- Go v1.21
- PostgreSQL
- `make` command for running the Makefile

## How to Prepare the Environment
First of all, you need to set up the configuration in the `.env` file. You can copy the `.env.example` file and rename it to `.env`. Then, you can fill in the configuration in the `.env` file.

the configuration in the `.env` file is as follows:

```shell
export GCP_PRIVATE_KEY=""
export GCP_CLIENT_EMAIL=""
export GCP_PROJECT_ID=""
export GCP_SECRET_ID=""
export GCP_SECRET_NAME=""
```

I attached the value of the above configuration in private repository. After you have the value of the configuration, you can fill in the value and use this command to export the configuration to the environment variable.

```shell
source .env
```

Or if you do not want to run `source` command every time you start new shell session, you can add the content of `.env` file to the very bottom of your `.bashrc` or `.zshrc` file depending on your shell type.

If you set your environment variables correctly, the configuration file should be generated when the application run for the first time. The configuration file will be created as `./config/config.json`.

## How to Run the Application

Before running the application, you need to install the swagger command-line interface (CLI) to generate the API documentation. You can install the swagger CLI by running the following command:

```shell
make swag-install
```

After installing the swagger CLI, you can run the application by running the following command:

```shell
make run
```

This will run the application in the development environment. You might got an error if you run the application for the first time. If that happens, you can edit the `config/config.json` file and match the configuration with your local environment, and then try again. You can access the API documentation by opening the following URL in your browser:

```shell
http://localhost:8080/docs/index.html
```

## Tips
- If you want to access the protected API, you need to add the `Authorization` header with the value `Bearer <access_token>` at the top right of the API documentation page. You can get the access token in the register / login endpoint.
