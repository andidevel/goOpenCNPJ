//    Copyright 2021 Anderson Rodrigues do Livramento

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package model

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/catfishlabs/goOpenCNPJ/consts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Errors
var ErrNoRows = errors.New("no rows returned")

type DateTime time.Time

// Partner export partner data
// Deprecated - Not used anymore
// type Partner struct {
// 	Identificacao      int64    `bson:"identificacao" json:"identificacao"`
// 	Nome               string   `bson:"nome" json:"nome"`
// 	Documento          string   `bson:"documento" json:"documento"`
// 	CodigoQualificacao string   `bson:"codigo_qualificacao" json:"codigo_qualificacao"`
// 	DataEntrada        DateTime `bson:"data_entrada" json:"data_entrada"`
// 	// CpfReprLegal           string   `bson:"cpf_repr_legal" json:"cpf_repr_legal"`
// 	NomeReprLegal          string `bson:"nome_repr_legal" json:"nome_repr_legal"`
// 	CodigoQualificacaoRepr string `bson:"codigo_qualificacao_rpr" json:"codigo_qualificacao_rpr"`
// }

type BaseCompany struct {
	ID                      string  `bson:"_id" json:"_id"`
	RazaoSocial             string  `bson:"razao_social" json:"razao_social"`
	CodigoNaturezaJuridica  int64   `bson:"codigo_natureza_juridica" json:"codigo_natureza_juridica"`
	QualificacaoResponsavel int64   `bson:"qualificacao_responsavel" json:"qualificacao_responsavel"`
	CapitalSocial           float64 `bson:"capital_social" json:"capital_social"`
	PorteEmpresa            int64   `bson:"porte_empresa" json:"porte_empresa"`
	EnteFederativo          string  `bson:"ente_federativo" json:"ente_federativo"`
}

// Company exports a cnpj document
type Company struct {
	ID                      string   `bson:"_id" json:"_id"`
	BaseID                  string   `bson:"empresa_base_id" json:"empresa_base_id"`
	IDMatriz                int64    `bson:"id_matriz" json:"id_matriz"`
	NomeFantasia            string   `bson:"nome_fantasia" json:"nome_fantasia"`
	SituacaoCadastral       int64    `bson:"situacao_cadastral" json:"situacao_cadastral"`
	DataSituacaoCadastral   DateTime `bson:"data_situacao_cadastral" json:"data_situacao_cadastral"`
	CodigoSituacaoCadastral int64    `bson:"codigo_situacao_cadastral" json:"codigo_situacao_cadastral"`
	MotivoSituacaoCadastral string   `bson:"motivo_situacao_cadastral" json:"motivo_situacao_cadastral"`
	NomeCidadeExterior      string   `bson:"nome_cidade_exterior" json:"nome_cidade_exterior"`
	CodigoPais              int64    `bson:"codigo_pais" json:"codigo_pais"`
	NomePais                string   `bson:"nome_pais" json:"nome_pais"`
	DataInicioAtividade     DateTime `bson:"data_inicio_atividade" json:"data_inicio_atividade"`
	CNAEFiscal              string   `bson:"cnae_fiscal" json:"cnae_fiscal"`
	CNAEsSecundarios        []string `bson:"cnaes_secundarios" json:"cnaes_secundarios"`
	GrauRisco               string   `bson:"grau_risco" json:"grau_risco"`
	TipoLogradouro          string   `bson:"tipo_logradouro" json:"tipo_logradouro"`
	Logradouro              string   `bson:"logradouro" json:"logradouro"`
	NumeroLogradouro        string   `bson:"numero_logradouro" json:"numero_logradouro"`
	Complemento             string   `bson:"complemento" json:"complemento"`
	Bairro                  string   `bson:"bairro" json:"bairro"`
	CEP                     string   `bson:"cep" json:"cep"`
	UF                      string   `bson:"uf" json:"uf"`
	CodigoMunicipio         int64    `bson:"codigo_municipio" json:"codigo_municipio"`
	NomeMunicipio           string   `bson:"nome_municipio" json:"nome_municipio"`
	DDD1                    string   `bson:"ddd1" json:"ddd1"`
	Telefone1               string   `bson:"telefone1" json:"telefone1"`
	DDD2                    string   `bson:"ddd2" json:"ddd2"`
	Telefone2               string   `bson:"telefone2" json:"telefone2"`
	DDDFax                  string   `bson:"ddd_fax" json:"ddd_fax"`
	Fax                     string   `bson:"fax" json:"fax"`
	Email                   string   `bson:"email" json:"email"`
	// OptanteSimples          string    `bson:"optante_simples" json:"optante_simples"`
	// DataOpcaoSimples        DateTime  `bson:"data_opcao_simples" json:"data_opcao_simples"`
	// DataExclusaoSimples     DateTime  `bson:"data_exclusao_simples" json:"data_exclusao_simples"`
	// OptanteMEI              string    `bson:"optante_mei" json:"optante_mei"`
	SituacaoEspecial     string   `bson:"situacao_especial" json:"situacao_especial"`
	DataSituacaoEspecial DateTime `bson:"data_situacao_especial" json:"data_situacao_especial"`
	// Socios                  []Partner `bson:"socios" json:"socios"`
}

// StatusDescription exports an status applied to a company
type StatusDescription struct {
	ID     int64  `bson:"_id" json:"_id"`
	Motivo string `bson:"motivo" json:"motivo"`
}

type RiskLevel struct {
	ID        string `bson:"_id" json:"_id"`
	GrauRisco string `bson:"grau_risco" json:"grau_risco"`
}

type Parameter struct {
	ID    string      `bson:"_id" json:"_id"`
	Value interface{} `bson:"value" json:"value"`
}

type City struct {
	ID            int64  `bson:"_id" json:"_id"`
	NomeMunicipio string `bson:"nome_municipio" json:"nome_municipio"`
}

type IDataStorage interface {
	Connect() error
	Close()

	FindOneUpsertParameter(Parameter) (Parameter, error)
	// SaveParameter(Parameter) error

	FindOneUpsertStatusDescription(StatusDescription) (StatusDescription, error)
	FindOneStatusDescriptionById(int64) (StatusDescription, error)
	// SaveStatusDescription(StatusDescription) error

	FindOneUpsertBaseCompany(BaseCompany) (BaseCompany, error)
	FindOneBaseCompanyById(string) (BaseCompany, error)
	// SaveBaseCompany(BaseCompany) error

	FindOneUpsertCompany(Company) (Company, error)
	FindOneCompanyById(string) (Company, error)
	// SaveCompany(Company) error

	FindOneUpsertRiskLevel(RiskLevel) (RiskLevel, error)
	FindOneRiskLevelById(string) (RiskLevel, error)
	// SaveRiskLevel(RiskLevel) error

	FindOneUpsertCity(City) (City, error)
	FindOneCityById(int64) (City, error)
	// SaveCity(City) error
}

// DB interface to be used in controllers
var DB IDataStorage

// DecodeFromMap is a simple struct decoder operating only at first level
func DecodeFromMap(m map[string]interface{}, v interface{}) {
	s := reflect.ValueOf(v).Elem()
	t := s.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Using JSON struct tag
		if tag, ok := field.Tag.Lookup("json"); ok {
			if tag != "" {
				mapKeyName := strings.Split(tag, ",")[0]
				if vmap, ok := m[mapKeyName]; ok && vmap != nil {
					s.FieldByName(field.Name).Set(reflect.ValueOf(vmap))
				}
			}
		}
	}
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	var t string
	var err error
	if err = json.Unmarshal(b, &t); err != nil {
		return err
	}
	dt, err := time.Parse(consts.DateLayoutJSON, t)
	if err != nil {
		return err
	}
	*d = DateTime(dt.UTC())
	return nil
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	strTime := ""
	if !t.IsZero() {
		strTime = t.Format(consts.DateLayoutJSON)
	}
	return json.Marshal(strTime)
}

func (d *DateTime) UnmarshalBSONValue(bType bsontype.Type, b []byte) error {
	var t time.Time
	rv := bson.RawValue{
		Type:  bType,
		Value: b,
	}
	// unixTime is in milliseconds
	// unixTime, ok := rv.DateTimeOK()
	t, ok := rv.TimeOK()
	if !ok {
		return errors.New("not a unix timestamp")
	}
	// t = time.Unix(unixTime/1000, 0)
	*d = DateTime(t.UTC())
	return nil
}

func (d DateTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	t := time.Time(d)
	_, b, err := bson.MarshalValue(t)
	return bsontype.DateTime, b, err
}
