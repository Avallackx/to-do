package repository

import (
	"context"
	"todo-app/internal/model"
	"todo-app/internal/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type taskRepo struct {
	db        *gorm.DB
	cacheRepo model.CacheRepository
}

func NewTaskRepository(db *gorm.DB) model.TaskRepository {
	return &taskRepo{db: db}
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
	err := tr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Task{}, ID).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"ID":  ID,
		}).Error(err)
		return err
	}

	return nil
}

func (tr *taskRepo) FindByID(ctx context.Context, ID int64) (*model.Task, error) {
	task := &model.Task{}

	err := tr.db.WithContext(ctx).Where("id = ?", ID).Take(&task).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"ID":  ID,
		}).Error(err)
		return nil, err
	}

	return task, nil
}

func (tr *taskRepo) FindAll(ctx context.Context) ([]*model.Task, error) {
	tasks := []*model.Task{}

	err := tr.db.WithContext(ctx).Order("id DESC").Find(&tasks).Error
	if err != nil {
		logrus.WithField("ctx", utils.Dump(ctx)).Error(err)
		return nil, err
	}

	return tasks, nil
}

func (tr *taskRepo) Update(ctx context.Context, task *model.Task) (*model.Task, error) {
	err := tr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(task).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":  utils.Dump(ctx),
			"book": utils.Dump(task),
		}).Error(err)
		return nil, err
	}

	return tr.FindByID(ctx, task.ID)
}
