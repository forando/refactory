package parser

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
	"strconv"
)

type Attributes struct {
	Map        map[string]*hclwrite.Attribute
	ModuleName string
}

func (attrs *Attributes) getAttr(key string) (*hclwrite.Attribute, error) {
	attr := attrs.Map[key]
	if attr == nil {
		return nil, errors.Errorf("[%s] property not found in %s module", key, attrs.ModuleName)
	}
	return attr, nil
}

func (attrs *Attributes) getStr(key string) (string, error) {
	attr := attrs.Map[key]
	if attr == nil {
		return "", errors.Errorf("[%s] property not found in %s module", key, attrs.ModuleName)
	}
	return getExpressionAsString(attr), nil
}

func (attrs *Attributes) getInt(key string) (int, error, error) {
	attr := attrs.Map[key]
	if attr == nil {
		return 0, errors.Errorf("[%s] property not found in %s module", key, attrs.ModuleName), nil
	}
	str := getExpressionAsString(attr)
	intVal, intErr := strconv.Atoi(str)
	if intErr != nil {
		return 0, nil, errors.Errorf("cannot parse %s vlaue of [%s] into int", key, str)
	}
	return intVal, nil, nil
}

func (attrs *Attributes) getBool(key string) (bool, error, error) {
	attr := attrs.Map[key]
	if attr == nil {
		return false, errors.Errorf("[%s] property not found in %s module", key, attrs.ModuleName), nil
	}
	str := getExpressionAsString(attr)
	boolVal, boolErr := strconv.ParseBool(str)
	if boolErr != nil {
		return false, nil, errors.Errorf("cannot parse %s vlaue of [%s] into bool", key, str)
	}
	return boolVal, nil, nil
}
