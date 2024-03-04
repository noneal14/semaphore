package bolt

import (
	"github.com/ansible-semaphore/semaphore/db"
)

/*
Integrations
*/
func (d *BoltDb) CreateIntegration(integration db.Integration) (db.Integration, error) {
	err := integration.Validate()

	if err != nil {
		return db.Integration{}, err
	}

	newIntegration, err := d.createObject(integration.ProjectID, db.IntegrationProps, integration)
	return newIntegration.(db.Integration), err
}

func (d *BoltDb) GetIntegrations(projectID int, params db.RetrieveQueryParams) (integrations []db.Integration, err error) {
	err = d.getObjects(projectID, db.IntegrationProps, params, nil, &integrations)
	return integrations, err
}

func (d *BoltDb) GetIntegration(projectID int, integrationID int) (integration db.Integration, err error) {
	err = d.getObject(projectID, db.IntegrationProps, intObjectID(integrationID), &integration)
	if err != nil {
		return
	}

	return
}

func (d *BoltDb) UpdateIntegration(integration db.Integration) error {
	err := integration.Validate()

	if err != nil {
		return err
	}

	return d.updateObject(integration.ProjectID, db.IntegrationProps, integration)

}

func (d *BoltDb) GetIntegrationRefs(projectID int, integrationID int) (db.IntegrationReferrers, error) {
	//return d.getObjectRefs(projectID, db.IntegrationProps, integrationID)
	return db.IntegrationReferrers{}, nil
}

/*
Integration Extractors
*/

func (d *BoltDb) CreateIntegrationExtractor(projectID int, integrationExtractor db.IntegrationExtractor) (db.IntegrationExtractor, error) {
	err := integrationExtractor.Validate()

	if err != nil {
		return db.IntegrationExtractor{}, err
	}

	newIntegrationExtractor, err := d.createObject(projectID, db.IntegrationExtractorProps, integrationExtractor)
	return newIntegrationExtractor.(db.IntegrationExtractor), err
}

func (d *BoltDb) GetIntegrationExtractors(projectID int, params db.RetrieveQueryParams, integrationID int) (extractors []db.IntegrationExtractor, err error) {
	var projectExtractors []db.IntegrationExtractor
	err = d.getObjects(projectID, db.IntegrationExtractorProps, params, nil, &projectExtractors)

	if err != nil {
		return
	}

	extractors = make([]db.IntegrationExtractor, 0)

	for _, extractor := range projectExtractors {
		if extractor.IntegrationID == integrationID {
			extractors = append(extractors, extractor)
		}
	}

	return extractors, err
}

func (d *BoltDb) GetIntegrationExtractor(projectID int, extractorID int, integrationID int) (db.IntegrationExtractor, error) {
	var extractor db.IntegrationExtractor
	err := d.getObject(projectID, db.IntegrationExtractorProps, intObjectID(extractorID), &extractor)

	return extractor, err

}

func (d *BoltDb) UpdateIntegrationExtractor(projectID int, integrationExtractor db.IntegrationExtractor) error {
	err := integrationExtractor.Validate()

	if err != nil {
		return err
	}

	return d.updateObject(projectID, db.IntegrationExtractorProps, integrationExtractor)
}

func (d *BoltDb) GetIntegrationExtractorRefs(projectID int, extractorID int, integrationID int) (db.IntegrationExtractorReferrers, error) {
	return d.getIntegrationExtractorRefs(projectID, db.IntegrationExtractorProps, extractorID)
}

/*
Integration ExtractValue
*/
func (d *BoltDb) GetIntegrationExtractValuesByExtractorID(extractorID int) (values []db.IntegrationExtractValue, err error) {
	err = d.getObjects(extractorID, db.IntegrationExtractValueProps, db.RetrieveQueryParams{}, nil, &values)
	return values, err
}

func (d *BoltDb) DeleteIntegrationExtractValue(projectID int, valueID int, extractorID int) error {
	return d.deleteObject(projectID, db.IntegrationExtractValueProps, intObjectID(valueID), nil)
}

func (d *BoltDb) GetIntegrationMatchersByExtractorID(extractorID int) (matchers []db.IntegrationMatcher, err error) {
	err = d.getObjects(extractorID, db.IntegrationMatcherProps, db.RetrieveQueryParams{}, nil, &matchers)

	return matchers, err
}

func (d *BoltDb) DeleteIntegrationExtractor(projectID int, extractorID int, integrationID int) error {
	values, err := d.GetIntegrationExtractValuesByExtractorID(extractorID)

	if err != nil {
		return err
	}

	for value := range values {
		d.DeleteIntegrationExtractValue(projectID, values[value].ID, extractorID)
	}

	matchers, err := d.GetIntegrationMatchersByExtractorID(extractorID)

	if err != nil {
		return err
	}

	for matcher := range matchers {
		d.DeleteIntegrationMatcher(projectID, matchers[matcher].ID, extractorID)
	}
	return d.deleteObject(integrationID, db.IntegrationExtractorProps, intObjectID(extractorID), nil)
}

func (d *BoltDb) CreateIntegrationExtractValue(projectId int, value db.IntegrationExtractValue) (db.IntegrationExtractValue, error) {
	err := value.Validate()

	if err != nil {
		return db.IntegrationExtractValue{}, err
	}

	newValue, err := d.createObject(projectId, db.IntegrationExtractValueProps, value)
	return newValue.(db.IntegrationExtractValue), err

}

func (d *BoltDb) GetIntegrationExtractValues(projectID int, params db.RetrieveQueryParams, extractorID int) (values []db.IntegrationExtractValue, err error) {
	values = make([]db.IntegrationExtractValue, 0)
	var allValues []db.IntegrationExtractValue

	err = d.getObjects(projectID, db.IntegrationExtractValueProps, params, nil, &allValues)

	if err != nil {
		return
	}

	for _, v := range allValues {
		if v.ExtractorID == extractorID {
			values = append(values, v)
		}
	}

	return
}

func (d *BoltDb) GetIntegrationExtractValue(projectID int, valueID int, extractorID int) (value db.IntegrationExtractValue, err error) {
	err = d.getObject(projectID, db.IntegrationExtractValueProps, intObjectID(valueID), &value)
	return value, err
}

func (d *BoltDb) UpdateIntegrationExtractValue(projectID int, integrationExtractValue db.IntegrationExtractValue) error {
	err := integrationExtractValue.Validate()

	if err != nil {
		return err
	}

	return d.updateObject(projectID, db.IntegrationExtractValueProps, integrationExtractValue)
}

func (d *BoltDb) GetIntegrationExtractValueRefs(projectID int, valueID int, extractorID int) (db.IntegrationExtractorChildReferrers, error) {
	return d.getIntegrationExtractorChildrenRefs(projectID, db.IntegrationExtractValueProps, valueID)
}

/*
Integration Matcher
*/
func (d *BoltDb) CreateIntegrationMatcher(projectID int, matcher db.IntegrationMatcher) (db.IntegrationMatcher, error) {
	err := matcher.Validate()

	if err != nil {
		return db.IntegrationMatcher{}, err
	}
	newMatcher, err := d.createObject(projectID, db.IntegrationMatcherProps, matcher)
	return newMatcher.(db.IntegrationMatcher), err
}

func (d *BoltDb) GetIntegrationMatchers(projectID int, params db.RetrieveQueryParams, extractorID int) (matchers []db.IntegrationMatcher, err error) {
	matchers = make([]db.IntegrationMatcher, 0)
	var allMatchers []db.IntegrationMatcher

	err = d.getObjects(projectID, db.IntegrationMatcherProps, db.RetrieveQueryParams{}, nil, &allMatchers)

	if err != nil {
		return
	}

	for _, v := range allMatchers {
		if v.ExtractorID == extractorID {
			matchers = append(matchers, v)
		}
	}

	return
}

func (d *BoltDb) GetIntegrationMatcher(projectID int, matcherID int, extractorID int) (matcher db.IntegrationMatcher, err error) {
	var matchers []db.IntegrationMatcher
	matchers, err = d.GetIntegrationMatchers(projectID, db.RetrieveQueryParams{}, extractorID)

	for _, v := range matchers {
		if v.ID == matcherID {
			matcher = v
		}
	}

	return
}

func (d *BoltDb) UpdateIntegrationMatcher(projectID int, integrationMatcher db.IntegrationMatcher) error {
	err := integrationMatcher.Validate()

	if err != nil {
		return err
	}

	return d.updateObject(projectID, db.IntegrationMatcherProps, integrationMatcher)
}

func (d *BoltDb) DeleteIntegrationMatcher(projectID int, matcherID int, extractorID int) error {
	return d.deleteObject(projectID, db.IntegrationMatcherProps, intObjectID(matcherID), nil)
}
func (d *BoltDb) DeleteIntegration(projectID int, integrationID int) error {
	extractors, err := d.GetIntegrationExtractors(projectID, db.RetrieveQueryParams{}, integrationID)

	if err != nil {
		return err
	}

	for extractor := range extractors {
		d.DeleteIntegrationExtractor(projectID, extractors[extractor].ID, integrationID)
	}

	return d.deleteObject(projectID, db.IntegrationProps, intObjectID(integrationID), nil)
}

func (d *BoltDb) GetIntegrationMatcherRefs(projectID int, matcherID int, extractorID int) (db.IntegrationExtractorChildReferrers, error) {
	return d.getIntegrationExtractorChildrenRefs(projectID, db.IntegrationMatcherProps, matcherID)
}
