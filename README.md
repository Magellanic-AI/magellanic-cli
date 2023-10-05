# Magellanic CLI

Magellanic is the continuous Access Control to protect Non-Human identities and Workloads
Learn more at [the website(https://www.magellanic.ai).

Get started at (Magellanic)[https://admin.magellanic.ai].

## Overview
Magellanic CLI is a simple CLI tool for integrating with the Magellanic REST API

## Usage
`magellanic-cli --help`

Command: `magellanic-cli config get --help`<br>

Params:

| Environment variable | Flag               | Description                                                                   |
|----------------------|--------------------|-------------------------------------------------------------------------------|
| DOTENV_PATH          | --dotenv_path (-d) | .env configuration file path                                                  |
| API_KEY              | --api_key (-a)     | API key                                                                       |
| API_URL              | --api_url (-u)     | API base URL (default: https://api.magellanic.ai)                             |
| CONFIG_ID            | --config_id (-c)   | ID of the configuration you want to fetch (available in Magellanic web panel) |
| FORMAT               | --format (-f)      | output format (one of json, yaml, dotenv)                                     |
| OUTPUT               | --output (-o)      | output path                                                                   |

