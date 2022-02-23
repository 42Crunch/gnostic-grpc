package generator

import (
	openapiv3 "github.com/google/gnostic/openapiv3"
)

//NamedSchemaOrReference
func (c *GrpcChecker) getSchemaFromSchemaOrReference(or *openapiv3.SchemaOrReference) *openapiv3.Schema {
	schema := or.GetSchema()
	if schema != nil {
		return schema
	}

	ref := or.GetReference().XRef
	refName := getRefName(ref)
	return c.getSchemaByRefName(refName)
}

func (c *GrpcChecker) getSchemaByRefName(name string) *openapiv3.Schema {
	for _, schemaOrReference := range c.document.Components.Schemas.GetAdditionalProperties() {
		if schemaOrReference.Name == name {
			return schemaOrReference.Value.GetSchema()
		}
	}

	return nil
}

func (c *GrpcChecker) getSchemaFromParametersOrReference(or *openapiv3.ParameterOrReference) *openapiv3.Schema {
	schema := or.GetParameter().GetSchema()
	if schema != nil {
		return c.getSchemaFromSchemaOrReference(schema)
	}

	ref := or.GetReference().XRef
	return c.getParamSchemaByRefName(getRefName(ref))
}

func (c *GrpcChecker) getParamSchemaByRefName(name string) *openapiv3.Schema {
	for _, namedParameterOrRef := range c.document.Components.Parameters.AdditionalProperties {
		if namedParameterOrRef.GetName() == name {
			param := namedParameterOrRef.Value.GetParameter()
			schema := c.getSchemaFromSchemaOrReference(param.Schema)
			name = param.GetName()

			return schema
		}
	}

	return nil
}
