package webpaysdk

import (
	"time"

	"github.com/pepelias/webpay-sdk/request"
	"gopkg.in/mgo.v2/bson"
)

const (
	productionHost          = "https://webpay3g.transbank.cl"
	integrationHost         = "https://webpay3gint.transbank.cl"
	integrationCommerceCode = "597055555532"
	integrationAPIKeySecret = "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C"
)

// Transaction para generar el token
type Transaction struct {
	BuyOrder  string  `json:"buy_order"`  // Orden1234 (Generado por el comercio)
	SessionID string  `json:"session_id"` // Sesion1234 (Generado por el comercio)
	Amount    float64 `json:"amount"`     // Monto a pagar (CLP)
	ReturnURL string  `json:"return_url"` // Recibirá el pago
}

// InitTransaction es la transacción ya inicializada
type InitTransaction struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

// configuration
type configuration struct {
	Request *request.Request
}

// TransactionResult es la información de la transacción
type TransactionResult struct {
	VCI                string                 `json:"vci"`
	Amount             float64                `json:"amount"`
	Status             string                 `json:"status"`
	BuyOrder           string                 `json:"buy_order" bson:"buy_order"`
	SessionID          string                 `json:"session_id" bson:"session_id"`
	CardDetail         map[string]interface{} `json:"card_detail" bson:"card_detail"`
	AccountingDate     string                 `json:"accouting_date" bson:"accouting_date"`
	TransactionDate    time.Time              `json:"transaction_date" bson:"transaction_date"`
	AuthorizationCode  string                 `json:"authorization_code" bson:"authorization_code"`
	PaymentTypeCode    string                 `json:"payment_type_mode" bson:"payment_type_mode"`
	ResponseCode       int                    `json:"response_code" bson:"response_code"`
	InstallmentsAmount float64                `json:"installments_amount" bson:"installments_amount"`
	InstallmentsNumber float64                `json:"installments_number" bson:"installments_number"`
	Balance            float64                `json:"balance"`
}

// Refund es un reembolso
type Refund struct {
	Type              string  `json:"type"`
	AuthorizationCode string  `json:"authorization_code" bson:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date" bson:"authorization_date"`
	NullifiedAmount   float64 `json:"nullified_amount" bson:"nullified_amount"`
	Balance           float64 `json:"balance"`
	ResponseCode      int     `json:"response_code" bson:"reponse_code"`
}

// NewIntegrationPlusNormal transacciones en entorno de desarrollo
func NewIntegrationPlusNormal() *configuration {
	return &configuration{
		Request: request.New(integrationHost+"/rswebpaytransaction/api/webpay/v1.0/transactions/", request.Headers{
			"Tbk-Api-Key-Id":     integrationCommerceCode,
			"Tbk-Api-Key-Secret": integrationAPIKeySecret,
		}),
	}
}

// NewPlusNormal transacciones en entorno de producción
func NewPlusNormal(commeceCode, apiKeySecret string) *configuration {
	return &configuration{
		Request: request.New(productionHost+"/rswebpaytransaction/api/webpay/v1.0/transactions/", request.Headers{
			"Tbk-Api-Key-Id":     commeceCode,
			"Tbk-Api-Key-Secret": apiKeySecret,
		}),
	}
}

// Init inicializa la transacción
func (c *configuration) Init(t *Transaction) (*InitTransaction, error) {
	transaction := &InitTransaction{}
	err := c.Request.POST("", nil, t, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// Confirm avisa que recibimos la información del pago
func (c *configuration) Confirm(token string) (*TransactionResult, error) {
	transaction := &TransactionResult{}
	err := c.Request.PUT(token, nil, nil, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// Refund (Reembolso) revierte o anula una transacción
func (c *configuration) Refund(token string, amount float64) (*Refund, error) {
	refund := &Refund{}
	err := c.Request.POST(token+"/refunds", nil, bson.M{"amount": amount}, refund)
	if err != nil {
		return nil, err
	}
	return refund, nil
}
