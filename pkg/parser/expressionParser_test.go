package parser

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
	"testing"
)

const groupPermissions = `group_permissions = {
"Cloud Shuttle"           = ["AWSAdministratorAccess"]
"pt-after_sales-order-processing" = [
  module.komueb_1260_extended_read_only_access_permission_set.permission_set_name
]
}`

const wrongTypeAttr = `group_permissions = ["hello World!!!"]`

const expression = `[
	"AWSAdministratorAccess",
	module.komueb_1260_extended_read_only_access_permission_set.permission_set_name
]`

var testAttrs = map[string]string{"correct-type": groupPermissions, "wrong-type": wrongTypeAttr}

func TestShouldParseObject(t *testing.T) {
	attr, err := getGroupPermissionsAttr("correct-type")
	if err != nil {
		t.Error(err)
		return
	}
	p := ExpressionParser{}
	res, err := p.ParseObjectExpr("group_permissions", attr.Expr().BuildTokens(nil).Bytes())
	if err != nil {
		t.Error(err)
		return
	}
	if len(*res) != 2 {
		t.Errorf("p.ParseObjectExpr: expexted 2 items in result but was %d", len(*res))
	}
	if val, ok := (*res)["Cloud Shuttle"]; ok {
		expected := `["AWSAdministratorAccess"]`
		actual := string(val)
		if actual != expected {
			t.Errorf("p.ParseObjectExpr:\nexpexted (%d chars):\n%s\nactual (%d chars):\n%s\n", len(expected), expected, len(actual), actual)
		}
	} else {
		t.Errorf("p.ParseObjectExpr: expexted to contain 'Cloud Shuttle' key but didn't")
	}
	if val, ok := (*res)["pt-after_sales-order-processing"]; ok {
		expected := `[
  module.komueb_1260_extended_read_only_access_permission_set.permission_set_name
]`
		actual := string(val)
		if actual != expected {
			t.Errorf("p.ParseObjectExpr:\nexpexted (%d chars):\n%s\nactual (%d chars):\n%s\n", len(expected), expected, len(actual), actual)
		}
	} else {
		t.Errorf("p.ParseObjectExpr: expexted to contain 'pt-after_sales-order-processing' key but didn't")
	}
}

func TestShouldFailIfWrongType(t *testing.T) {
	attr, err := getGroupPermissionsAttr("wrong-type")
	if err != nil {
		t.Error(err)
		return
	}
	p := ExpressionParser{}
	if _, err := p.ParseObjectExpr("group_permissions", attr.Expr().BuildTokens(nil).Bytes()); err == nil {
		t.Errorf("p.ParseObjectExpr: expected to return error but didn't")
	}
}

func TestShouldParseArray(t *testing.T) {
	bytes := []byte(expression)
	pSetNames := map[string]string{
		"komueb_1260_extended_read_only_access_permission_set": "ExtendedReadOnlyAccess",
	}
	ctx := buildContext(&pSetNames)
	p := ExpressionParser{Context: ctx}
	res, err := p.ParseArrayExpr("Cloud Shuttle", bytes)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res) != 2 {
		t.Errorf("p.ParseArrayExpr: expexted 2 items in result but was %d", len(res))
	}
	if item := res[0]; item != "AWSAdministratorAccess" {
		t.Errorf("p.ParseArrayExpr: expexted item to be 'AWSAdministratorAccess' but was %s", item)
	}
	if item := res[1]; item != "ExtendedReadOnlyAccess" {
		t.Errorf("p.ParseArrayExpr: expexted item to be 'ExtendedReadOnlyAccess' but was %s", item)
	}
}

func getGroupPermissionsAttr(key string) (*hclwrite.Attribute, error) {
	testAttr, ok := testAttrs[key]
	if !ok {
		return nil, errors.Errorf("key: %s not found in testAttrs", key)
	}
	bytes := []byte(testAttr)
	file, configDiags := hclwrite.ParseConfig(bytes, "dummy.hcl", hcl.Pos{Line: 1, Column: 1})
	if configDiags.HasErrors() {
		return nil, configDiags
	}
	attrs := file.Body().Attributes()
	attr, ok := attrs["group_permissions"]
	if !ok {
		return nil, errors.New("Cannot find group_permissions attribute")
	}
	return attr, nil
}
