package repository

import "context"

func (r *ActivityRepository) UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error {
	return r.updateTaskState(ctx, userID, messageID, "completed")
}

func (r *ActivityRepository) UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error {
	return r.updateTaskState(ctx, userID, messageID, "fail")
}

func (r *ActivityRepository) updateTaskState(ctx context.Context, userID string, messageID string, state string) error {
	return setTaskState(ctx, r.db.Shard(r.sharder.DBKey(userID)), userID, messageID, state)
}
