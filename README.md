# Get Open data for Brazilian CNPJ

**goOpenCNPJ** takes advantage of "*Open Information Laws*" in Brazil, where public data about companies are available in a series of 
downloadable files from the **Federal Revenue** website.

**goOpenCNPJ** can download those files, which are in a special format (CSV), and import them to a **MongoDB** database and make them available through a REST API.

**goOpenCNPJ** also gets more information about companies to make the data more complete, extracting data from the Regulatory Standard number 4 - NR4 - which is a PDF. **goOpenCNPJ** does that in order to associate the risk level of a company - in terms of occupational safety and health - according to its *National Classification of Economic Activities*, *CNAE* in Portuguese.

**goOpenCNPJ** is a learning project. If you need something more professional, take a look at [Minha Receita](https://github.com/cuducos/minha-receita).

This project was inspired by [Brasil.io](https://brasil.io). Check it out for more public datasets and search options.

# Folder structure

```
goopencnpj\
  |- get-companies
  |- goOpenCNPJ
  |- .env
  |- config\
    |- companies-download.json
    |- cnpj-schema.json
```
Use the *.env-example* provided by this repo as example, fill it with your configuration and rename it to *.env*.

## .env file

The server and the *get-companies* command relies on the *.env* file for configurations like:

```
DATA_DOWNLOAD_PATH=/path/to/downlad/data
DBURI=mongodb://user:password@host:port/database
PORT=6543
```

## Config Folder

In *config* folder we need two JSON files:

- **companies-download.json**

URLs to download data. Don't change it unless you know what you are doing.
```json
{

   "companies.main.url": "https://www.gov.br/receitafederal/pt-br/assuntos/orientacao-tributaria/cadastros/consultas/dados-publicos-cnpj",
   "companies.mirrors.url": [],
   "nr04.url": "https://www.gov.br/trabalho/pt-br/inspecao/seguranca-e-saude-no-trabalho/normas-regulamentadoras/nr-04.pdf/@@download/file/NR-04.pdf"
}
```

- **cnpj-schema.json**

This file has the "schema" for the **CSV** files from **Federal Revenue**. Again, don't change it unless you know what you are doing.

Sample:

```json
[
   {
      "type": "0",
      "document": {
         "_id": {
               "field_type": "str",
               "position": 0
         },
         "razao_social": {
               "field_type": "str",
               "position": 1
         },
         "natureza_juridica": {
               "field_type": "int",
               "position": 2
         },
         "porte_empresa": {
               "field_type": "int",
               "position": 5
         }
      }
   }
]
```

# Database creation

The command *get-companies* is used to create and update data. Usage:
```
NAME:
   get-companies - download and parse Brazil's companies database

USAGE:
   get-companies [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  path to companies JSON config file
   --schema value, -s value  path to companies JSON schema file
   --aux, -a                 download and parse auxiliary tables only
   --force, -f               force download and parse
   --nfiles value, -n value  number of data files to download
   --help, -h                show help
```
Example:
```
$ ./get-companies -c /home/user/.config/urls.json -a
```

Will download and import only auxiliary tables (Company Status and NR-4) using URLs provided by */home/user/.config/urls.json* file.

Note: **urls.json** must has the same format as **companies-download.json** in the folder **config**.

# REST API

Just run:
```
./goOpenCNPJ

goOpenCNPJ server listening on localhost:6543
```

## Getting company data

Assuming that you set *6543* as *PORT* in *.env* configuration file:

```
curl --request GET \
  --url http://localhost:6543/cnpj/<CNPJ>
```

Change *\<CNPJ\>* with **CNPJ** code of the desired company.

**Example Response**:

```json
{
  "data": {
    "razao_social": "FULANO DA SILVA",
    "_id": "65747887000121",
    "empresa_base_id": "65747887",
    "id_matriz": 1,
    "nome_fantasia": "FULANO SA",
    "situacao_cadastral": 2,
    "data_situacao_cadastral": "2017-10-09",
    "codigo_situacao_cadastral": 0,
    "motivo_situacao_cadastral": "SEM MOTIVO",
    "nome_cidade_exterior": "",
    "codigo_pais": 0,
    "nome_pais": "",
    "data_inicio_atividade": "2017-10-09",
    "cnae_fiscal": "6120501",
    "cnaes_secundarios": [
      "7220700",
      "8412400",
      "8511200",
      "8513900",
      "8599604",
      "8630502",
      "8630503"
    ],
    "grau_risco": "2",
    "tipo_logradouro": "",
    "logradouro": "",
    "numero_logradouro": "",
    "complemento": "",
    "bairro": "",
    "cep": "",
    "uf": "SC",
    "codigo_municipio": 8327,
    "nome_municipio": "SAO JOSE",
    "ddd1": "48",
    "telefone1": "2345678",
    "ddd2": "",
    "telefone2": "",
    "ddd_fax": "",
    "fax": "",
    "email": "FULANO@FULANOSA.COM.BR",
    "situacao_especial": "",
    "data_situacao_especial": ""
  },
  "error": ""
}
```

```
curl --request GET \
  --url http://localhost:6543/nr04/<CNAE>
```

Change *\<CNAE\>* with **CNAE** code.

**Example Response**:

```json
{
  "data": {
    "_id": "61205",
    "grau_risco": "2"
  },
  "error": ""
}
```