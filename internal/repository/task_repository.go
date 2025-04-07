package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"todo-app/internal/model"
	"todo-app/internal/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type taskRepo struct {
	db        *gorm.DB
	cacheRepo model.CacheRepository
}

func NewTaskRepository(db *gorm.DB, cacheRepo model.CacheRepository) model.TaskRepository {
	return &taskRepo{
		db:        db,
		cacheRepo: cacheRepo,
	}
}

func (tr *taskRepo) Create(ctx context.Context, task *model.Task) error {

	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.Dump(ctx),
		"task": utils.Dump(task),
	})

	err := tr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(task).Error; err != nil {
			return err

		}
		return nil
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	cacheKeys := []string{
		tr.cacheHash(),
		tr.countAllCacheKey(),
	}

	if err := tr.cacheRepo.Delete(ctx, cacheKeys...); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (tr *taskRepo) DeleteByID(ctx context.Context, ID int64) error {

	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.Dump(ctx),
		"ID":  ID,
	})

	err := tr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Task{}, ID).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	cacheKeys := []string{
		tr.findByIDCacheKey(ID),
		tr.cacheHash(),
	}

	if err := tr.cacheRepo.Delete(ctx, cacheKeys...); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (tr *taskRepo) FindByID(ctx context.Context, ID int64) (*model.Task, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.Dump(ctx),
		"ID":  ID,
	})

	cacheKey := tr.findByIDCacheKey(ID)

	reply, err := tr.cacheRepo.Get(ctx, cacheKey)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if reply != "" {
		task := &model.Task{}
		if err := json.Unmarshal([]byte(reply), &task); err != nil {
			logger.Error(err)
			return nil, err
		}
		return task, nil
	}

	task := &model.Task{}

	err = tr.db.WithContext(ctx).Where("id = ?", ID).Take(&task).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	bytes, err := json.Marshal(task)
	if err != nil {
		logger.Error(err)
		return task, nil
	}

	if err := tr.cacheRepo.Set(ctx, cacheKey, string(bytes)); err != nil {
		logger.Error(err)
	}

	return task, nil
}

func (tr *taskRepo) FindAll(ctx context.Context, query model.GetTasksQueryParams) ([]*model.Task, error) {

	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.Dump(ctx),
		"query": utils.Dump(query),
	})

	cacheHash := tr.cacheHash()
	cacheKey := tr.findAllByQueryParams(query)
	reply, err := tr.cacheRepo.HashGet(ctx, cacheHash, cacheKey)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if reply != "" {
		tasks := []*model.Task{}
		if err := json.Unmarshal([]byte(reply), &tasks); err != nil {
			logger.Error(err)
			return nil, err
		}
		return tasks, nil
	}

	tasks := []*model.Task{}

	err = tr.db.WithContext(ctx).
		Order("id DESC").
		Offset(int(model.Offset(query.Page, query.Size))).
		Limit(int(query.Size)).
		Find(&tasks).
		Error

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	bytes, err := json.Marshal(tasks)
	if err != nil {
		logger.Error(err)
		return tasks, nil
	}

	if err := tr.cacheRepo.HashSet(ctx, cacheHash, cacheKey, string(bytes)); err != nil {
		logger.Error(err)
	}

	return tasks, nil
}

func (tr *taskRepo) CountAll(ctx context.Context) (int64, error) {
	logger := logrus.WithField("ctx", utils.Dump(ctx))

	cacheKey := tr.countAllCacheKey()
	reply, err := tr.cacheRepo.Get(ctx, cacheKey)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	if reply != "" {
		count := int64(0)
		if err := json.Unmarshal([]byte(reply), &count); err != nil {
			logger.Error(err)
			return 0, err
		}
		return count, nil
	}

	count := int64(0)
	err = tr.db.WithContext(ctx).
		Model(model.Task{}).
		Count(&count).
		Error
	if err != nil {
		logrus.WithField("ctx", utils.Dump(ctx)).Error(err)
		return int64(0), err
	}

	bytes, err := json.Marshal(count)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	if err := tr.cacheRepo.Set(ctx, cacheKey, string(bytes)); err != nil {
		logger.Error(err)
	}

	return count, nil
}

func (tr *taskRepo) Update(ctx context.Context, task *model.Task) (*model.Task, error) {

	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.Dump(ctx),
		"task": utils.Dump(task),
	})

	err := tr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(task).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	cacheKeys := []string{
		tr.cacheHash(),
		tr.countAllCacheKey(),
		tr.findByIDCacheKey(task.ID),
	}

	if err := tr.cacheRepo.Delete(ctx, cacheKeys...); err != nil {
		logger.Error(err)
		return nil, err
	}

	return tr.FindByID(ctx, task.ID)
}

func (tr *taskRepo) cacheHash() string {
	return "task"
}

func (tr *taskRepo) findByIDCacheKey(ID int64) string {
	return fmt.Sprintf("task:%d", ID)
}

func (tr *taskRepo) findAllByQueryParams(query model.GetTasksQueryParams) string {
	return fmt.Sprintf("task:page:%d:size:%d", query.Page, query.Size)
}

func (tr *taskRepo) countAllCacheKey() string {
	return "task:count"
}
