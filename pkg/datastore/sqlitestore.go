package datastore

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/model"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLiteStore struct {
	logger *zap.Logger
	db     *gorm.DB
}

type resourceId string

func NewSQLiteStore(ctx context.Context, cfg config.Config, zapLogger *zap.Logger) (*SQLiteStore, error) {
	s := SQLiteStore{}
	logLevel := logger.Error
	if zapLogger.Core().Enabled(zap.DebugLevel) {
		//log all SQL queries
		logLevel = logger.Info
	}
	//gormLogger has it's own logger for SQL queries - better than zaplog for that purpose
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	s.logger = zapLogger
	//create the DB client
	var err error
	s.db, err = gorm.Open(sqlite.Open(cfg.Datastore.DataSourceName),
		&gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("can't create the SQLite database: %w", err)
	}

	// Migrate the schema
	if err = s.db.AutoMigrate(&model.Resource{}, &model.Tag{}); err != nil {
		return nil, fmt.Errorf("can't create the SQLite data model: %w", err)
	}

	return &s, nil
}

func (s *SQLiteStore) Ping() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

func (s *SQLiteStore) getAllResourceIds(ctx context.Context) ([]resourceId, error) {
	var resourceIds []resourceId
	result := s.db.Model(&model.Resource{}).Select("id").Distinct().Find(&resourceIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return resourceIds, nil
}

func (s *SQLiteStore) getResourceIdsByTag(ctx context.Context, tags model.Tags) ([]resourceId, error) {
	var resourceIds []resourceId
	db := s.db.Model(&model.Tag{}).Select("resource_id").Distinct()
	if !tags.Empty() {
		for _, tag := range tags {
			if tag.Value == "*" {
				//wild card - check for any value
				db = db.Or("key=? and value like ?", tag.Key, "%")
			} else {
				//exact value
				db = db.Or("key=? and value=?", tag.Key, tag.Value)
			}
		}

		// our data model has one row per tag, we need a group by to only keep the results matching all tags
		/* table tags
		//original data, tag filter = team=infra and cluster=staging
		1. i-123	cluster		staging
		2. i-123	team		infra
		3. i-124	cluster		staging
		4. i-124	team		dev

		//where team=infra OR cluster=staging
		1. i-123	cluster		staging
		2. i-123	team		infra
		3. i-124	cluster		staging

		// group by + count on id
		1. i-123	2 <- we only want to keep this row
		2. i-124	1
		*/
		db = db.Group("resource_id").
			Having("count(resource_id)=?", tags.DistinctKeys())
	}
	result := db.Find(&resourceIds)
	if result.Error != nil {
		return nil, result.Error
	}

	return resourceIds, nil
}

func (s *SQLiteStore) getResourcesById(ctx context.Context, idsInclude []resourceId, idsExclude []resourceId) ([]*model.Resource, error) {
	var resources []*model.Resource
	if len(idsInclude) == 0 {
		return resources, nil
	}
	db := s.db
	if len(idsExclude) != 0 {
		// exlucde these ids
		db = db.Where("id not in ?", idsExclude)
	}
	db = db.
		Preload("Tags").Find(&resources, idsInclude)

	if db.Error != nil {
		return nil, db.Error
	}
	return resources, nil

}
func (s *SQLiteStore) GetResource(ctx context.Context, id string) (*model.Resource, error) {
	resources, err := s.getResourcesById(ctx, []resourceId{resourceId(id)}, nil)
	if err != nil {
		return nil, fmt.Errorf("can't get resource from database: %w", err)
	}
	if len(resources) == 0 {
		//not found
		return nil, nil
	}
	s.logger.Sugar().Infow("Getting resource: ",
		zap.String("id", id),
	)
	return resources[0], nil
}
func (s *SQLiteStore) GetResources(ctx context.Context, filter model.Filter) ([]*model.Resource, error) {
	var resources []*model.Resource

	//apply filter tags:include
	includeTags := filter.TagsInclude()
	var err error
	var idsInclude []resourceId
	if includeTags.Empty() {
		//no include filter set - include all
		idsInclude, err = s.getAllResourceIds(ctx)
	} else {
		idsInclude, err = s.getResourceIdsByTag(ctx, includeTags)
	}
	if err != nil {
		return nil, fmt.Errorf("can't get resources from database: %w", err)
	}

	//apply filter tags:exclude
	var idsExclude []resourceId
	if tagsExclude := filter.TagsExclude(); !tagsExclude.Empty() {
		idsExclude, err = s.getResourceIdsByTag(ctx, filter.TagsExclude())
		if err != nil {
			return nil, fmt.Errorf("can't get resources from database: %w", err)
		}
	}

	resources, err = s.getResourcesById(ctx, idsInclude, idsExclude)
	if err != nil {
		return nil, fmt.Errorf("can't get resources from database: %w", err)
	}
	s.logger.Sugar().Infow("Getting resources: ",
		zap.Object("filter", filter),
		zap.Int("count", len(resources)),
	)
	return resources, nil
}

func (s *SQLiteStore) WriteResources(ctx context.Context, resources []*model.Resource) error {
	if len(resources) == 0 {
		//nothing to write
		return nil
	}
	result := s.db.Create(resources)
	s.logger.Sugar().Infow("Writting resources: ",
		zap.Int("count", int(result.RowsAffected)),
	)
	if result.Error != nil {
		return fmt.Errorf("can't write resources to database: %w", result.Error)
	}
	return nil
}

func (s *SQLiteStore) Stats(context.Context) (model.Stats, error) {
	var count int64
	result := s.db.Table("resources").Count(&count)
	if result.Error != nil {
		return model.Stats{}, fmt.Errorf("can't read resources count from database: %w", result.Error)
	}
	return model.Stats{ResourcesCount: int(count)}, nil
}

func (s *SQLiteStore) getResourceField(name string) (model.Field, error) {
	/*
		SELECT DISTINCT `type` , count(*)
		FROM `resources`
		group by  `type`
		order by `type`
	*/
	rows, err := s.db.Model(&model.Resource{}).Select(name, "count() as count").
		Distinct().
		Group(name).
		Order(name).
		Rows()
	if err != nil {
		return model.Field{}, fmt.Errorf("can't get resource field '%v' from database: %w", name, err)
	}
	defer rows.Close()
	field := model.Field{
		Name:  name,
		Group: "core",
	}
	var totalCount int
	for rows.Next() {
		var value string
		var count int
		err = rows.Scan(&value, &count)
		if err != nil {
			return model.Field{}, fmt.Errorf("can't get resource field '%v' from database: %w", name, err)
		}
		field.Values = append(field.Values, model.FieldValue{
			Value: value,
			Count: count,
		})
		totalCount = totalCount + count
	}
	field.Count = totalCount
	return field, nil
}

func (s *SQLiteStore) getTagFields() (model.Fields, error) {
	fields, err := s.getTagKeys()
	if err != nil {
		return model.Fields{}, fmt.Errorf("can't get tag keys from database: %w", err)
	}
	var result model.Fields
	for _, f := range fields {
		values, err := s.getTagValues(f.Name)
		if err != nil {
			return nil, err
		}
		f.Values = values
		result = append(result, f)
	}
	return result, nil
}

func (s *SQLiteStore) getTagKeys() (model.Fields, error) {
	/*
		SELECT distinct(`key`), count() as count
		FROM `tags`
		group by key
		order by count desc
	*/
	rows, err := s.db.Model(&model.Tag{}).Select("key", "count() as count").
		Distinct().
		Group("key").
		Order("count desc").
		Rows()
	if err != nil {
		return model.Fields{}, fmt.Errorf("can't get tag keys from database: %w", err)
	}
	var fields model.Fields
	defer rows.Close()
	for rows.Next() {
		var key string
		var count int
		err = rows.Scan(&key, &count)
		if err != nil {
			return nil, err
		}
		field := model.Field{
			Group: "tags",
			Name:  key,
			Count: count,
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func (s *SQLiteStore) getTagValues(key string) (model.FieldValues, error) {
	/*
		SELECT distinct(`value`), count() as count
		FROM `tags`
		where key=?
		group by value
		order by count desc
	*/
	var values model.FieldValues
	db := s.db.Model(&model.Tag{}).Select("value", "count() as count").
		Distinct().
		Where("key=?", key).
		Group("value").
		Order("count desc").
		Find(&values)

	if db.Error != nil {
		return nil, fmt.Errorf("can't get tag value for key '%v' : %w", key, db.Error)
	}
	return values, nil
}

func (s *SQLiteStore) GetFields(context.Context) (model.Fields, error) {
	var fields model.Fields

	//get fields on resources
	for _, name := range []string{"region", "type"} {
		field, err := s.getResourceField(name)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	//get fields on tags
	tagFields, err := s.getTagFields()
	if err != nil {
		return nil, err
	}
	fields = append(fields, tagFields...)

	return fields, nil
}
