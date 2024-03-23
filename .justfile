set dotenv-filename := "just.env"
set dotenv-load

import? 'development.justfile'
import? 'production.justfile'

env := env_var_or_default('WERF_ENV', "Development")

checkout env:
  echo '{{ if env == "prod" { "WERF_ENV=\"Production\"" } else if env == "dev" { "WERF_ENV=\"Development\"" } else { error("There is no such environment") } }}' > just.env

synth:
  @just {{if env == "Production" { '_synth-production' } else if env == "Development" { '_synth-development' } else { error("Invalid environment value") } }}


encrypt:
  @just {{if env == "Production" { '_encrypt-production' } else if env == "Development" { '_encrypt-development' } else { error("Invalid environment value") } }}
decrypt:
  @just {{if env == "Production" { '_decrypt-production' } else if env == "Development" { '_decrypt-development' } else { error("Invalid environment value") } }}


up *FLAGS:
  @just {{if env == "Production" { '_up-production' } else if env == "Development" { '_up-development' } else { error("Invalid environment value") } }} {{FLAGS}}

down *FLAGS:
  @just {{if env == "Production" { '_down-production' } else if env == "Development" { '_down-development' } else { error("Invalid environment value") } }} {{FLAGS}}
