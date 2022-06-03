# Harvester

This Harvester is a service that consumes on-chain data from blockchain networks and offers them to a message queue.
Currently it provides functionality to consume data from Substrate based chains and publish this data on a MQTT instance. More particularly it is the data provider for the Kusamaverse, available on kusama.odyssey.xyz. 

The harvester uses an actors model that is extensible to other networks and data processing cases.

Next to publishing to MQTT is also stores data in a MySQL database.

## Usage

Use `make run` to run the application. 

The harvester's MQTT and database can be configured through environment variables or through the `config.yaml`. By default it will load the config.yaml, and override any values given by the environment variables.  
It is possible to harvest data from different chains concurrently by setting the `EnabledChains` config parameter.

## Prerequisites

1. Go >= v1.17.2
2. Ent - https://entgo.io/docs/tutorial-setup
3. golang-migrate/migrate - https://github.com/golang-migrate/migrate
4. Makefile
5. Docker
6. docker-compose

The harvester has MQTT and MySQL as dependencies.

## Run database migrations
`make prerequisites #Install all the CLI prerequisites for code generation e.t.c.`
`make db-migrations # Create the first initial migrations`
`make db-migrate-up name=${migration_name} # Create the updated migrations`
<br/>
or
<br/>
`make db-migrate-up --database=mysql://root@password@tcp(localhost:3306)/harvester_dev`

## Test

`make test`

## Contributors ✨

Thanks go to these wonderful people 😎

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>

  <tr>
  <td align="center"><a href="https://github.com/jellevdp"><img src="https://avatars.githubusercontent.com/jellevdp?v=3?s=100" width="100px;" alt=""/><br /><sub><b>Jelle van der Ploeg </b></sub></a><br />
    </td>
<td align="center"><a href="https://github.com/tech-sam"><img src="https://avatars.githubusercontent.com/tech-sam?v=3?s=100" width="100px;" alt=""/><br /><sub><b>Sumit</b></sub></a><br />
</td>
<td align="center"><a href="https://github.com/longyarnz"><img src="https://avatars.githubusercontent.com/longyarnz?v=3?s=100" width="100px;" alt=""/><br /><sub><b>Ayodele Olalekan</b></sub></a><br />
</td>
  </tr>

  <tr>
   <td align="center"><a href="https://github.com/e-nikolov"><img src="https://avatars.githubusercontent.com/e-nikolov" width="100px;" alt=""/><br /><sub><b>Emil Nikolov  </b></sub></a><br />
    </td>
  <td align="center"><a href="https://github.com/rwajon"><img src="https://avatars.githubusercontent.com/rwajon" width="100px;" alt=""/><br /><sub><b>Jonathan Rwabahizi  </b></sub></a><br />
    </td>
      <td align="center"><a href="https://github.com/nwasiqUC"><img src="https://avatars.githubusercontent.com/nwasiqUC" width="100px;" alt=""/><br /><sub><b>Wasiq  </b></sub></a><br />
    </td>

  </tr>

  <tr>
   <td align="center"><a href="https://github.com/antst"><img src="https://avatars.githubusercontent.com/antst" width="100px;" alt=""/><br /><sub><b>Anton Starikov</b></sub></a><br />
    </td>
  <td align="center"><a href="https://github.com/jor-rit"><img src="https://avatars.githubusercontent.com/jor-rit" width="100px;" alt=""/><br /><sub><b>Jorrit</b></sub></a><br />
    </td>

  </tr>
</table>
