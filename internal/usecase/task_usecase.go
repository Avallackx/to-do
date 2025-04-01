package usecase

import (
	"context"

	"github.com/sirupsen/logrus"

	"todo-app/internal/model"
	"todo-app/internal/utils"
)

type taskUsecase struct {
	taskRepo model.TaskRepository
}

func NewTaskUsecase(tr model.TaskRepository) model.TaskUsecase {
	return &taskUsecase{taskRepo: tr}
}

func (tu *taskUsecase) Create(ctx context.Context, task *model.Task) (*model.Task, error) {

	if err := tu.taskRepo.Create(ctx, task); err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":  utils.Dump(ctx),
			"task": utils.Dump(task),
		}).Error(err)
		return nil, err
	}

	return task, nil
}

func (tu *taskUsecase) DeleteByID(ctx context.Context, ID int64) error {
	if err := tu.taskRepo.DeleteByID(ctx, ID); err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"ID":  ID,
		}).Error(err)
		return err
	}

	return nil
}

func (tu *taskUsecase) FindByID(ctx context.Context, ID int64) (*model.Task, error) {
	task, err := tu.taskRepo.FindByID(ctx, ID)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"ID":  ID,
		}).Error(err)
		return nil, err
	}

	return task, nil
}

func (tu *taskUsecase) FindAll(ctx context.Context, params model.GetTasksQueryParams) ([]*model.Task, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.Dump(ctx),
		"params": utils.Dump(params),
	})

	tasks, err := tu.taskRepo.FindAll(ctx, params)
	if err != nil {
		logger.Error(err)
		return nil, int64(0), err
	}

	count, err := tu.taskRepo.CountAll(ctx)
	if err != nil {
		logger.Error(err)
		return nil, int64(0), err
	}

	return tasks, count, nil
}

func (tu *taskUsecase) Update(ctx context.Context, task *model.Task) (*model.Task, error) {
	task, err := tu.taskRepo.Update(ctx, task)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":  utils.Dump(ctx),
			"task": utils.Dump(task),
		}).Error(err)
		return nil, err
	}

	return task, nil
}
