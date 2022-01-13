# K3241.K03200Y0.D10710.EMPRECSV

| #  | Field                       | Description                                                                                                  |
| -- | --------------------------- | ------------------------------------------------------------------------------------------------------------ |
| 00 | CNPJ BÁSICO                 | NÚMERO BASE DE INSCRIÇÃO NO CNPJ (OITO PRIMEIROS DÍGITOS DO CNPJ)                                            |
| 01 | RAZÃO SOCIAL                | NOME EMPRESARIAL DA PESSOA JURÍDICA                                                                          | 	
| 02 | NATUREZA JURÍDICA           | CÓDIGO DA NATUREZA JURÍDICA                                                                                  |
| 03 | QUALIFICAÇÃO DO RESPONSÁVEL | QUALIFICAÇÃO DA PESSOA FÍSICA RESPONSÁVEL PELA EMPRESA                                                       |
| 04 | CAPITAL SOCIAL DA EMPRESA   | CAPITAL SOCIAL DA EMPRESA                                                                                    |
| 05 | PORTE DA EMPRESA            | CÓDIGO DO PORTE DA EMPRESA: 1 – NÃO INFORMADO, 2 - MICRO EMPRESA, 03 - EMPRESA DE PEQUENO PORTE, 05 - DEMAIS |
| 06 | ENTE FEDERATIVO RESPONSÁVEL | O ENTE FEDERATIVO RESPONSÁVEL É PREENCHIDO PARA OS CASOS DE ÓRGÃOS E ENTIDADES DO GRUPO DE NATUREZA JURÍDICA 1XXX. PARA AS DEMAIS NATUREZAS, ESTE ATRIBUTO FICA EM BRANCO. |


# K3241.K03200Y0.D10710.ESTABELE

| #  | Field                       | Description                                                                                       |
| -- | --------------------------- | ------------------------------------------------------------------------------------------------- |
| 00 | CNPJ BÁSICO                 | NÚMERO BASE DE INSCRIÇÃO NO CNPJ (OITO PRIMEIROS DÍGITOS DO CNPJ)                                 | 
| 01 | CNPJ ORDEM                  | NÚMERO DO ESTABELECIMENTO DE INSCRIÇÃO NO CNPJ (DO NONO ATÉ O DÉCIMO SEGUNDO DÍGITO DO CNPJ)      |
| 02 | CNPJ DV                     | DÍGITO VERIFICADOR DO NÚMERO DE INSCRIÇÃO NO CNPJ (DOIS ÚLTIMOS DÍGITOS DO CNPJ)                  |
| 03 | IDENTIFICADOR MATRIZ/FILIAL | CÓDIGO DO IDENTIFICADOR MATRIZ/FILIAL: 1 – MATRIZ, 2 – FILIAL                                     |
| 04 | NOME FANTASIA               | CORRESPONDE AO NOME FANTASIA                                                                      |
| 05 | SITUAÇÃO CADASTRAL          | CÓDIGO DA SITUAÇÃO CADASTRAL: 01 – NULA, 2 – ATIVA, 3 – SUSPENSA, 4 – INAPTA, 08 – BAIXADA        |
| 06 | DATA SITUAÇÃO CADASTRAL     | DATA DO EVENTO DA SITUAÇÃO CADASTRAL                                                              |
| 07 | MOTIVO SITUAÇÃO CADASTRAL   | CÓDIGO DO MOTIVO DA SITUAÇÃO CADASTRAL                                                            |  	
| 08 | NOME DA CIDADE NO EXTERIOR  | NOME DA CIDADE NO EXTERIOR                                                                        |  	
| 09 | PAIS                        | CÓDIGO DO PAIS                                                                                    |
| 10 | DATA DE INÍCIO ATIVIDADE    | DATA DE INÍCIO DA ATIVIDADE                                                                       |
| 11 | CNAE FISCAL PRINCIPAL       | CÓDIGO DA ATIVIDADE ECONÔMICA PRINCIPAL DO ESTABELECIMENTO                                        |
| 12 | CNAE FISCAL SECUNDÁRIA      | CÓDIGO DA(S) ATIVIDADE(S) ECONÔMICA(S) SECUNDÁRIA(S) DO ESTABELECIMENTO                           |
| 13 | TIPO DE LOGRADOURO          | DESCRIÇÃO DO TIPO DE LOGRADOURO                                                                   |
| 14 | LOGRADOURO                  | NOME DO LOGRADOURO ONDE SE LOCALIZA O ESTABELECIMENTO                                             |
| 15 | NÚMERO                      | NÚMERO ONDE SE LOCALIZA O ESTABELECIMENTO. QUANDO NÃO HOUVER PREENCHIMENTO DO NÚMERO HAVERÁ ‘S/N’ |
| 16 | COMPLEMENTO                 | COMPLEMENTO PARA O ENDEREÇO DE LOCALIZAÇÃO DO ESTABELECIMENTO                                     |
| 17 | BAIRRO                      | BAIRRO ONDE SE LOCALIZA O ESTABELECIMENTO                                                         |
| 18 | CEP                         | CÓDIGO DE ENDEREÇAMENTO POSTAL REFERENTE AO LOGRADOURO NO QUAL O ESTABELECIMENTO ESTA LOCALIZADO  |
| 19 | UF                          | SIGLA DA UNIDADE DA FEDERAÇÃO EM QUE SE ENCONTRA O ESTABELECIMENTO                                |
| 20 | MUNICÍPIO                   | CÓDIGO DO MUNICÍPIO DE JURISDIÇÃO ONDE SE ENCONTRA O ESTABELECIMENTO                              |
| 21 | DDD 1                       | CONTÉM O DDD 1                                                                                    |
| 22 | TELEFONE 1                  | CONTÉM O NÚMERO DO TELEFONE 1                                                                     |
| 23 | DDD 2                       | CONTÉM O DDD 2                                                                                    |
| 24 | TELEFONE 2                  | CONTÉM O NÚMERO DO TELEFONE 2                                                                     |
| 25 | DDD DO FAX                  | CONTÉM O DDD DO FAX                                                                               |
| 26 | FAX                         | CONTÉM O NÚMERO DO FAX                                                                            |
| 27 | CORREIO ELETRÔNICO          | CONTÉM O E-MAIL DO CONTRIBUINTE                                                                   |
| 28 | SITUAÇÃO ESPECIAL           | SITUAÇÃO ESPECIAL DA EMPRESA                                                                      |
| 29 | DATA DA SITUAÇÃO ESPECIAL   | DATA EM QUE A EMPRESA ENTROU EM SITUAÇÃO ESPECIAL                                                 | 	

# F.K03200$Z.D10710.MOTICSV

| #  | Field                       | Description                                                                                       |
| -- | --------------------------- | ------------------------------------------------------------------------------------------------- |
| 00 | CÓDIGO MOTIVO               | CÓDIGO DO MOTIVO DA SITUAÇÃO CADASTRAL                                                            | 
| 01 | DESCRIÇÃO                   | DESCRIÇÃO DO MOTIVO DA SITUAÇÃO CADASTRAL                                                         |

# F.K03200$Z.D10710.CNAECSV

| #  | Field                       | Description                                                                                       |
| -- | --------------------------- | ------------------------------------------------------------------------------------------------- |
| 00 | CÓDIGO                      | CÓDIGO DA ATIVIDADE ECONÔMICA                                                                     | 
| 01 | DESCRIÇÃO                   | NOME DA ATIVIDADE ECONÔMICA                                                                       |
