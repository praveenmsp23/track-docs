package base

type Token interface {
	Set(key, value string) error        //set token value
	Get(key string) (string, bool)      //get token value
	Delete(key string) error            //delete token value
	TokenID() string                    //back current tokenID
	GetAll() (map[string]string, error) //get all values of token
}

type Provider interface {
	TokenInit(tid string) (Token, error)
	TokenRead(tid string) (Token, error)
	TokenDestroy(tid string) error
	TokenGC(maxlifetime int64)
}
