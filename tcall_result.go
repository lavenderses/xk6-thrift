package thrift

type TCallResult struct {
	body TValue
	err error
}

func NewTCallResult(body *TValue, err error) *TCallResult {
	switch {
	case body == nil:
		return &TCallResult{body: nil, err: err}
	case err != nil:
		return &TCallResult{body: nil, err: err}
	default:
		return &TCallResult{body: *body, err: err}
	}
}

func (r *TCallResult) IsSuccess() bool {
	return r.err == nil
}
