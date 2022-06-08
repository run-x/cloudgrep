package datastore

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/model"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type SQLiteStore struct {
	logger  *zap.Logger
	db      *gorm.DB
	indexer resourceIndexer
	//fetchedAt is the last time the resources were fetched
	fetchedAt       time.Time
	muResource      sync.Mutex
	muResourceEvent sync.Mutex
}

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
	if err = s.db.AutoMigrate(&model.Resource{}, &model.Tag{}, &model.ResourceEvent{}); err != nil {
		return nil, fmt.Errorf("can't create the SQLite data model: %w", err)
	}

	//create the indexer
	s.indexer, err = newResourceIndexer(ctx, s.logger, s.db)
	if err != nil {
		return nil, fmt.Errorf("can't create the query builder: %w", err)
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

func (s *SQLiteStore) getResourcesById(ctx context.Context, ids []model.ResourceId) ([]*model.Resource, error) {
	var resources []*model.Resource
	if len(ids) == 0 {
		return resources, nil
	}
	db := s.db.Preload("Tags").Find(&resources, ids)

	if db.Error != nil {
		return nil, db.Error
	}
	return resources, nil

}
func (s *SQLiteStore) GetResource(ctx context.Context, id string) (*model.Resource, error) {
	resources, err := s.getResourcesById(ctx, []model.ResourceId{model.ResourceId(id)})
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

func (s *SQLiteStore) WriteResources(ctx context.Context, resources model.Resources) error {
	if len(resources) == 0 {
		//nothing to write
		return nil
	}
	s.muResource.Lock()
	defer s.muResource.Unlock()

	var count int64
	err := s.db.Transaction(func(tx *gorm.DB) error {
		//delete all the previously stored tags if any
		ids := resources.Ids()
		if err := deleteTags(tx, ids); err != nil {
			return err
		}

		// Create or Update all columns
		result := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(resources)
		count = result.RowsAffected

		// Create or Update the resource indexes
		if err := s.indexer.writeResourceIndexes(ctx, tx, resources); err != nil {
			return err
		}
		return result.Error
	})

	s.logger.Sugar().Infow("Writting resources: ", zap.Int64("count", count))

	if err != nil {
		return fmt.Errorf("can't write resources to database: %w", err)
	}
	return nil
}

func resourceCount(db *gorm.DB) (int, error) {
	var count int64
	return int(count), db.Table("resources").Count(&count).Error
}

func (s *SQLiteStore) Stats(context.Context) (model.Stats, error) {
	count, err := resourceCount(s.db)
	if err != nil {
		return model.Stats{}, fmt.Errorf("can't read stats: %w", err)
	}
	return model.Stats{ResourcesCount: count}, nil
}

func (s *SQLiteStore) getResourceField(name string) (model.Field, error) {
	/*
		SELECT DISTINCT `type` , count(*) as count
		FROM `resources`
		group by  `type`
		order by `type`
		sort by `count` desc
	*/
	rows, err := s.db.Model(&model.Resource{}).Select(name, "count() as count").
		Distinct().
		Group(name).
		Order("count desc").
		Rows()
	if err != nil {
		return model.Field{}, fmt.Errorf("can't get resource field '%v' from database: %w", name, err)
	}
	defer rows.Close()
	field := model.Field{
		Name: name,
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

func (s *SQLiteStore) GetFields(context.Context) (model.FieldGroups, error) {
	var fieldGroups model.FieldGroups

	//get core fields
	coreGroup := model.FieldGroup{
		Name: "core",
	}
	for _, name := range []string{"region", "type"} {
		field, err := s.getResourceField(name)
		if err != nil {
			return nil, err
		}
		coreGroup.Fields = append(coreGroup.Fields, field)
	}
	fieldGroups = append(fieldGroups, coreGroup)

	//get tag fields
	tagFields, err := s.getTagFields()
	if err != nil {
		return nil, err
	}
	tagsGroup := model.FieldGroup{
		Name:   "tags",
		Fields: tagFields,
	}
	fieldGroups = append(fieldGroups, tagsGroup)

	return fieldGroups.AddNullValues(), nil
}

func (s *SQLiteStore) GetResources(ctx context.Context, jsonQuery []byte) (model.Resources, error) {
	ids, err := s.indexer.findResourceIds(*s.db, s.logger, jsonQuery)
	if err != nil {
		return nil, err
	}
	return s.getResourcesById(ctx, ids)
}

func (s *SQLiteStore) WriteResourceEvent(ctx context.Context, resourceEvent model.ResourceEvent) error {
	s.muResourceEvent.Lock()
	defer s.muResourceEvent.Unlock()
	result := s.db.Create(&resourceEvent)
	if result.Error != nil {
		return fmt.Errorf("can't write resource event to database: %w", result.Error)
	}
	return nil
}

func (s *SQLiteStore) CaptureEngineStart(ctx context.Context) {
	s.fetchedAt = time.Now()
}

func (s *SQLiteStore) CaptureEngineEnd(ctx context.Context) error {
	//once engine is complete, we delete all the resources that no longer exist
	_, err := s.deleteResourcesBefore(s.fetchedAt)
	if err != nil {
		return err
	}
	return nil
}

func deleteTags(db *gorm.DB, ids []model.ResourceId) error {
	return db.Table("tags").Where("resource_id in ?", ids).Delete(ids).Error
}

func (s *SQLiteStore) deleteResourcesBefore(before time.Time) (int, error) {

	var rowsAffected int64
	err := s.db.Transaction(func(tx *gorm.DB) error {
		//get the resource ids to delete
		var ids []model.ResourceId
		if err := tx.Table("resources").Select("id").Where("updated_at < ?", before).Find(&ids).Error; err != nil {
			return err
		}

		if len(ids) == 0 {
			//nothing to delete
			return nil
		}

		totalCount, err := resourceCount(tx)
		if err != nil {
			return err
		}
		if totalCount == len(ids) {
			s.logger.Sugar().Warnf("deleting resources before last run would delete all resources (count: %v), ignoring delete.", totalCount)
			// the most common case for this would be the engine encountered a global error,
			// and we are now writting the resource event in the DB to record the error.
			// we woud rather keep the previous resources than deleting them all.
			return nil
		}

		//delete all the tags
		if err := deleteTags(tx, ids); err != nil {
			return err
		}

		//delete the resource indexes and purge the unused columns
		if err := s.indexer.deleteResourceIndexes(tx, ids); err != nil {
			return err
		}
		if err := s.indexer.purgeUnusedColumns(tx); err != nil {
			return err
		}

		//delete the resources
		result := tx.Table("resources").Where("id in ?", ids).Delete(ids)
		if result.Error != nil {
			return result.Error
		}
		rowsAffected = result.RowsAffected
		return nil
	})

	s.logger.Sugar().Infow("Deleting resources: ", zap.Int64("rowsAffected", rowsAffected))

	if err != nil {
		return 0, fmt.Errorf("can't delete resources from database: %w", err)
	}
	return int(rowsAffected), nil
}

func (s *SQLiteStore) GetEngineStatus(context.Context) (model.EngineStatus, error) {
	var resourceEvents model.ResourceEvents
	result := s.db.
		Model(&model.ResourceEvent{}).
		Where("id in (?)",
			s.db.
				Model(&model.ResourceEvents{}).
				Select("max(id)").
				Group("resource_type")).
		Find(&resourceEvents)

	if result.Error != nil {
		return model.NewEngineStatus(nil, true),
			fmt.Errorf("unable to fetch the resource events: %w", result.Error)
	}

	return model.NewEngineStatus(resourceEvents, false), nil
}
